"""CNC-02: bounded training directly on packed INT4 expert codes.

The learner stores ONLY packed uint8 codes and frozen float32 block scales.
No float32/fp16 master weight, optimizer moments, gradients, or residuals are
retained between steps. One exact coordinate update is selected per step from
the training batch's squared error. Development-only; synthetic input/output
fixtures, no Qwen model forward or scientific admission.

The source INT4 tensor is encoded as signed nibbles -7..7; -8 is forbidden.
"""
import numpy as np
from packed_int4 import apply_packed_int4, pack_int4, stored_bytes

BLOCK = 64


def _code_at(packed, row, col):
    raw = int(packed[row, col // 2])
    nibble = (raw >> (4 * (col % 2))) & 15
    return nibble - 16 if nibble >= 8 else nibble


def _set_code(packed, row, col, code):
    if not -7 <= code <= 7:
        raise ValueError("signed INT4 code overflow")
    shift = 4 * (col % 2)
    mask = 15 << shift
    old = int(packed[row, col // 2])
    packed[row, col // 2] = (old & (~mask & 255)) | ((code & 15) << shift)


class PackedCodeLearner:
    """Immutable-scale, one-coordinate-per-step compressed-state learner."""

    def __init__(self, packed, scales):
        if not isinstance(packed, np.ndarray) or not isinstance(scales, np.ndarray):
            raise ValueError("missing packed state")
        self.packed = np.array(packed, copy=True)
        self.scales = np.array(scales, copy=True)
        # Validate state with an empty-probe-free call; no full weights allocated.
        apply_packed_int4(
            np.zeros((1, self.packed.shape[1] * 2), dtype=np.float32),
            self.packed, self.scales, block=BLOCK,
        )
        self.steps = 0

    @classmethod
    def from_weight(cls, weight):
        codes, scales = pack_int4(weight, block=BLOCK)
        return cls(codes, scales)

    @property
    def resident_weight_bytes(self):
        return stored_bytes(self.packed, self.scales, self.packed.shape[1] * 2, block=BLOCK)

    def predict(self, x):
        return apply_packed_int4(x, self.packed, self.scales, block=BLOCK)

    def train_step(self, x, target):
        """Choose strict-loss-decreasing +/-1 nibble edit, no full-weight update."""
        x = np.asarray(x, dtype=np.float32)
        target = np.asarray(target, dtype=np.float32)
        if (x.ndim != 2 or x.shape[1] != self.packed.shape[1] * 2 or
                x.shape[0] < 1 or target.shape != (x.shape[0], self.packed.shape[0]) or
                not np.isfinite(x).all() or not np.isfinite(target).all()):
            raise ValueError("invalid fixed-shape learning batch")
        # Prediction/residual are transient activation workspace, not weight
        # or optimizer master copies; count separately from persistent state.
        pred = self.predict(x)
        residual = pred.astype(np.float64) - target.astype(np.float64)
        x64 = x.astype(np.float64)
        norms = np.sum(x64 * x64, axis=0)
        dot = residual.T @ x64
        scale_for_col = np.repeat(self.scales.astype(np.float64), BLOCK, axis=1)
        best = (0.0, None, None, None)
        for row in range(self.packed.shape[0]):
            for col in range(x.shape[1]):
                current = _code_at(self.packed, row, col)
                s = scale_for_col[row, col]
                for direction in (-1, 1):
                    candidate = current + direction
                    if not -7 <= candidate <= 7:
                        continue
                    delta = s * direction
                    change = 2.0 * delta * dot[row, col] + delta * delta * norms[col]
                    if change < best[0]:
                        best = (float(change), row, col, candidate)
        self.steps += 1
        if best[1] is None:
            return {"changed": False, "loss_decrease_sum_squared": 0., "steps": self.steps}
        _, row, col, code = best
        _set_code(self.packed, row, col, code)
        return {"changed": True, "loss_decrease_sum_squared": -best[0], "steps": self.steps,
                "updated_row": row, "updated_column": col}


def fit_control(weight, x, target, steps=8):
    """Full FP32 master-weight SGD control; diagnostic, not deployment."""
    w = np.array(weight, dtype=np.float32, copy=True)
    x = np.asarray(x, dtype=np.float32)
    y = np.asarray(target, dtype=np.float32)
    if (type(steps) is not int or steps < 0 or steps > 64 or
            w.ndim != 2 or x.ndim != 2 or x.shape[1] != w.shape[1] or
            y.shape != (x.shape[0], w.shape[0])):
        raise ValueError("invalid baseline")
    for _ in range(steps):
        gradient = ((x @ w.T - y).T @ x) / x.shape[0]
        w -= 0.03 * gradient / max(1., float(np.sum(x * x) / x.shape[0]))
    return w
