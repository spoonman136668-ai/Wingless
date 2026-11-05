"""CNC-01Q-G9C: packed 2-bit, blockwise CPU development operator.

Frozen symmetric mid-rise codebook (-3,-1,1,3). Four 2-bit codes per byte;
FP16 per-row, per-block scales. Computes a block at a time without constructing
a full dense dequantized weight. Not a Qwen GPU kernel or inference.
"""
import numpy as np

LEVELS = np.array([-3, -1, 1, 3], dtype=np.int8)
ALLOWED_BLOCKS = (128, 256)


def pack_int2(weight, *, block=128):
    w = np.asarray(weight, dtype=np.float32)
    if (type(block) is not int or block not in ALLOWED_BLOCKS or
            w.ndim != 2 or min(w.shape) <= 0 or w.shape[1] % block or
            not np.isfinite(w).all()):
        raise ValueError("invalid INT2 expert shape or block")
    rows, cols = w.shape
    groups = cols // block
    packed = np.empty((rows, cols // 4), dtype=np.uint8)
    scales = np.empty((rows, groups), dtype=np.float16)
    for k in range(groups):
        tile = w[:, k*block:(k+1)*block]
        peak = np.max(np.abs(tile), axis=1)
        candidate = (peak / 3.).astype(np.float32)
        # Use the actually stored scale during quantization; fail closed
        # on underflow/overflow rather than silently changing representation.
        stored = candidate.astype(np.float16)
        if not np.isfinite(stored).all() or np.any((peak > 0) & (stored <= 0)):
            raise ValueError("FP16 quantization scale overflow/underflow")
        scales[:, k] = stored
        divisor = np.where(stored > 0, stored, np.float16(1.)).astype(np.float32)
        unit = tile / divisor[:, None]
        # Boundaries -2, 0, 2 map to one of the four fixed mid-rise values.
        codes = np.clip(np.floor((unit + 4.) / 2.), 0, 3).astype(np.uint8)
        out = (codes[:, 0::4] |
               (codes[:, 1::4] << 2) |
               (codes[:, 2::4] << 4) |
               (codes[:, 3::4] << 6))
        packed[:, (k*block)//4:((k+1)*block)//4] = out
    return packed, scales


def _validate(packed, scales, cols, block):
    if (not isinstance(packed, np.ndarray) or packed.dtype != np.uint8 or
            packed.ndim != 2 or not isinstance(scales, np.ndarray) or
            scales.dtype != np.float16 or scales.ndim != 2 or
            type(cols) is not int or cols <= 0 or
            type(block) is not int or block not in ALLOWED_BLOCKS or
            cols % block or packed.shape[0] <= 0 or packed.shape[1]*4 != cols or
            scales.shape != (packed.shape[0], cols // block) or
            not np.isfinite(scales).all() or np.any(scales < 0)):
        raise ValueError("invalid INT2 operator state")


def apply_packed_int2(inputs, packed, scales, *, block=128):
    x = np.asarray(inputs, dtype=np.float32)
    if x.ndim != 2 or min(x.shape) <= 0 or not np.isfinite(x).all():
        raise ValueError("invalid INT2 activation")
    _validate(packed, scales, x.shape[1], block)
    result = np.zeros((x.shape[0], packed.shape[0]), dtype=np.float32)
    for k in range(x.shape[1] // block):
        tile = packed[:, k*block//4:(k+1)*block//4]
        codes = np.empty((tile.shape[0], block), dtype=np.uint8)
        for offset in range(4):
            codes[:, offset::4] = (tile >> (offset*2)) & 3
        weights = LEVELS[codes].astype(np.float32)
        result += (x[:, k*block:(k+1)*block] @ weights.T) * scales[:, k][None, :].astype(np.float32)
    if not np.isfinite(result).all():
        raise ValueError("nonfinite INT2 result")
    return result


def stored_bytes(packed, scales, cols, *, block=128):
    _validate(packed, scales, cols, block)
    return int(packed.nbytes + scales.nbytes + 3*8)
