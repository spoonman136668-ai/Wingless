"""CNC-01Q-G7: development-only blockwise packed-INT4 expert operation.

Each weight is a signed four-bit code. The operator decodes one input-column
block at a time and multiplies codes before applying per-row block scales.
It never constructs a full dense dequantized weight matrix. NumPy may still
create temporary arrays for one block. CPU-only: NOT Qwen inference, a fused
GPU kernel, real-weight quality, or scientific acceptance.
"""
import numpy as np

MAX_Q = 7


def pack_int4(weight, block=64):
    w = np.asarray(weight, dtype=np.float32)
    if (w.ndim != 2 or min(w.shape) <= 0 or
            type(block) is not int or block <= 0 or block % 2 or
            w.shape[1] % block or not np.isfinite(w).all()):
        raise ValueError("invalid blockwise INT4 input")
    rows, cols = w.shape
    groups = cols // block
    quant = np.empty((rows, cols), dtype=np.int8)
    scales = np.empty((rows, groups), dtype=np.float32)
    for k in range(groups):
        tile = w[:, k*block:(k+1)*block]
        max_abs = np.max(np.abs(tile), axis=1)
        scale = np.where(max_abs > 0, max_abs / MAX_Q, 1.0).astype(np.float32)
        if not np.isfinite(scale).all() or np.any(scale <= 0):
            raise ValueError("invalid blockwise INT4 scales")
        scales[:, k] = scale
        quant[:, k*block:(k+1)*block] = np.clip(
            np.rint(tile / scale[:, None]), -MAX_Q, MAX_Q
        ).astype(np.int8)
    lo = quant[:, 0::2].astype(np.uint8) & 15
    hi = quant[:, 1::2].astype(np.uint8) & 15
    packed = (lo | (hi << 4)).astype(np.uint8)
    return packed, scales


def _validate(packed, scales, input_cols, block):
    if (not isinstance(packed, np.ndarray) or packed.dtype != np.uint8 or
            packed.ndim != 2 or not isinstance(scales, np.ndarray) or
            scales.dtype != np.float32 or scales.ndim != 2 or
            type(input_cols) is not int or input_cols <= 0 or
            type(block) is not int or block <= 0 or block % 2 or
            input_cols % block or packed.shape[1] * 2 != input_cols or
            packed.shape[0] <= 0 or
            scales.shape != (packed.shape[0], input_cols // block) or
            not np.isfinite(scales).all() or np.any(scales <= 0)):
        raise ValueError("invalid packed INT4 operator state")


def _block_codes(packed_tile):
    lo = (packed_tile & 15).astype(np.int16)
    hi = (packed_tile >> 4).astype(np.int16)
    lo = np.where(lo >= 8, lo - 16, lo)
    hi = np.where(hi >= 8, hi - 16, hi)
    codes = np.empty((packed_tile.shape[0], packed_tile.shape[1] * 2), dtype=np.int16)
    codes[:, 0::2] = lo
    codes[:, 1::2] = hi
    if np.any(codes < -MAX_Q) or np.any(codes > MAX_Q):
        raise ValueError("invalid signed four-bit weight code")
    return codes


def apply_packed_int4(inputs, packed, scales, block=64):
    x = np.asarray(inputs, dtype=np.float32)
    if x.ndim != 2 or min(x.shape) <= 0 or not np.isfinite(x).all():
        raise ValueError("invalid inputs")
    _validate(packed, scales, x.shape[1], block)
    out = np.zeros((x.shape[0], packed.shape[0]), dtype=np.float32)
    for k in range(x.shape[1] // block):
        codes = _block_codes(packed[:, (k*block)//2:((k+1)*block)//2])
        # Integer codes are multiplied before per-block scaling; no full
        # dense dequantized weight tensor is ever reconstructed.
        out += (x[:, k*block:(k+1)*block] @ codes.astype(np.float32).T) * scales[:, k][None, :]
    if not np.isfinite(out).all():
        raise ValueError("nonfinite packed INT4 output")
    return out


def stored_bytes(packed, scales, input_cols, block=64):
    _validate(packed, scales, input_cols, block)
    # Three fixed-width shape metadata fields, in addition to actual buffers.
    return int(packed.nbytes + scales.nbytes + 3 * 8)
