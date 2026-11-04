"""G9K: dense low-rank residual factors over existing packed G9F INT2.

An alternative correction family to G9H's sparse per-row corrections.
Encode-time residual SVD only; direct forward never reconstructs dense W.
No GPU/model-forward, original Qwen weights or native activations consumed.
"""
import numpy as np
from packed_companded_int2 import pack_companded_int2, apply_companded_int2, stored_bytes as base_bytes

RANKS = (1, 2, 3)


def pack_lowrank_residual(weight, *, rank=1, block=256):
    if type(rank) is not int or rank not in RANKS:
        raise ValueError("unfrozen low rank")
    w = np.asarray(weight, dtype=np.float32)
    if w.ndim != 2 or min(w.shape) < rank or not np.isfinite(w).all():
        raise ValueError("invalid weight shape or magnitude")
    packed, scales = pack_companded_int2(w, block=block)
    # Dense residual and thin SVD exist only in CPU ENCODER.
    # Compressed inference below uses only G9F + two FP16 low-rank factors.
    residual = w - apply_companded_int2(
        np.eye(w.shape[1], dtype=np.float32), packed, scales, block=block).T
    u, sigma, vt = np.linalg.svd(residual, full_matrices=False)
    a = (u[:, :rank] * np.sqrt(sigma[:rank])[None, :]).astype(np.float16)
    b = (vt[:rank, :].T * np.sqrt(sigma[:rank])[None, :]).astype(np.float16)
    if not np.isfinite(a).all() or not np.isfinite(b).all():
        raise ValueError("FP16 factor overflow")
    return packed, scales, a, b


def _check(packed, scales, a, b, cols, rank, block):
    if (type(rank) is not int or rank not in RANKS or
        not isinstance(a, np.ndarray) or a.dtype != np.float16 or a.ndim != 2 or
        not isinstance(b, np.ndarray) or b.dtype != np.float16 or b.ndim != 2 or
        a.shape != (packed.shape[0], rank) or b.shape != (cols, rank) or
        not np.isfinite(a).all() or not np.isfinite(b).all()):
        raise ValueError("invalid FP16 correction factors")
    base_bytes(packed, scales, cols, block=block)


def apply_lowrank_residual(inputs, packed, scales, a, b, *, rank=1, block=256):
    x = np.asarray(inputs, dtype=np.float32)
    if x.ndim != 2 or min(x.shape) < 1 or not np.isfinite(x).all():
        raise ValueError("invalid input activation")
    _check(packed, scales, a, b, x.shape[1], rank, block)
    base = apply_companded_int2(x, packed, scales, block=block)
    out = base + (x @ b.astype(np.float32)) @ a.astype(np.float32).T
    if not np.isfinite(out).all():
        raise ValueError("nonfinite output")
    return out


def stored_bytes(packed, scales, a, b, cols, *, rank=1, block=256):
    _check(packed, scales, a, b, cols, rank, block)
    return int(base_bytes(packed, scales, cols, block=block) + a.nbytes + b.nbytes + 24)


def correction_bytes_for_shape(shape, *, rank=1):
    if (type(rank) is not int or rank not in RANKS or
        not isinstance(shape, tuple) or len(shape) != 2 or
        any(type(n) is not int or n < rank for n in shape)):
        raise ValueError("unsupported low-rank source tensor")
    return 2 * rank * sum(shape) + 24
