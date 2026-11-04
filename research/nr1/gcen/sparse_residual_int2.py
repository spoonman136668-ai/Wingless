"""CNC-G9H: fixed per-row sparse residual atop G9F compressed INT2.

Development-only CPU operator. Stores uint16 local-column indexes, float16
residuals, uint32 row offsets. No separate learned codebooks, reranking from
evaluation activations, full-matrix runtime decode, GPU, or model forward.
All correction selection uses ORIGINAL WEIGHT residuals at encoding time only.
"""
import numpy as np
from packed_companded_int2 import (
    pack_companded_int2, apply_companded_int2, stored_bytes as base_bytes
)

BUDGET_K = (1, 2, 3)  # Corrections per row; immutable research variants.


def _check_k(k):
    if type(k) is not int or k not in BUDGET_K:
        raise ValueError("unfrozen sparse residual corrections per row")


def pack_sparse_residual(weight, *, k=2, block=256):
    _check_k(k)
    w = np.asarray(weight, dtype=np.float32)
    if w.ndim != 2 or min(w.shape) <= 0 or w.shape[1] > 65535 or k > w.shape[1] or not np.isfinite(w).all():
        raise ValueError("invalid sparse residual weight shape or magnitude")
    packed, scales = pack_companded_int2(w, block=block)
    # Diagnostic encoder uses a bounded row basis to reconstruct at encode
    # time only; the inference operator below never reconstructs full W.
    rows, cols = w.shape
    recovered = apply_companded_int2(
        np.eye(cols, dtype=np.float32), packed, scales, block=block).T
    residual = w - recovered
    selected = np.argsort(-np.abs(residual), axis=1, kind="stable")[:, :k]
    selected = np.sort(selected, axis=1).astype(np.uint16)
    delta = np.take_along_axis(residual, selected.astype(np.intp), axis=1)
    values = delta.astype(np.float16)
    if not np.isfinite(values).all() or np.any((delta != 0) & (values == 0)):
        raise ValueError("sparse FP16 residual overflow or underflow")
    offsets = (np.arange(rows + 1, dtype=np.uint64) * k)
    if offsets[-1] > np.iinfo(np.uint32).max:
        raise ValueError("row offset width exceeded")
    return packed, scales, offsets.astype(np.uint32), selected, values


def _validate(packed, scales, offsets, indexes, values, cols, k, block):
    _check_k(k)
    if (type(cols) is not int or cols <= 0 or cols > 65535 or k > cols or
        not isinstance(indexes, np.ndarray) or indexes.dtype != np.uint16 or
        indexes.ndim != 2 or indexes.shape != (packed.shape[0], k) or
        not isinstance(values, np.ndarray) or values.dtype != np.float16 or
        values.shape != indexes.shape or not np.isfinite(values).all() or
        not isinstance(offsets, np.ndarray) or offsets.dtype != np.uint32 or
        offsets.shape != (packed.shape[0] + 1,) or
        not np.array_equal(offsets, np.arange(packed.shape[0] + 1, dtype=np.uint64) * k) or
        np.any(indexes >= cols) or
        np.any(np.diff(indexes.astype(np.int32), axis=1) <= 0)):
        raise ValueError("invalid sparse correction representation")
    base_bytes(packed, scales, cols, block=block)


def apply_sparse_residual(inputs, packed, scales, offsets, indexes, values, *, k=2, block=256):
    x = np.asarray(inputs, dtype=np.float32)
    if x.ndim != 2 or min(x.shape) <= 0 or not np.isfinite(x).all():
        raise ValueError("invalid input activation shape")
    _validate(packed, scales, offsets, indexes, values, x.shape[1], k, block)
    out = apply_companded_int2(x, packed, scales, block=block)
    for row in range(indexes.shape[0]):
        out[:, row] += x[:, indexes[row].astype(np.intp)] @ values[row].astype(np.float32)
    if not np.isfinite(out).all():
        raise ValueError("sparse direct operator nonfinite")
    return out


def stored_bytes(packed, scales, offsets, indexes, values, cols, *, k=2, block=256):
    _validate(packed, scales, offsets, indexes, values, cols, k, block)
    return (base_bytes(packed, scales, cols, block=block) +
            offsets.nbytes + indexes.nbytes + values.nbytes)


def correction_bytes_for_shape(shape, *, k=2):
    _check_k(k)
    if (not isinstance(shape, tuple) or len(shape) != 2 or
        any(type(x) is not int or x <= 0 for x in shape) or
        shape[1] > 65535 or k > shape[1] or shape[0] * k > 0xffffffff):
        raise ValueError("incompatible full-model correction shape")
    rows, _ = shape
    return 4 * (rows + 1) + 4 * rows * k
