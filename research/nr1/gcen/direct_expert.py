"""CNC-01Q development-only direct factorized expert operator.

This is not full Qwen inference, a GPU kernel, or scientific acceptance.
Evaluate y = x @ (U diag(S) Vh).T without materializing the dense matrix.
"""
import numpy as np


def factorize(weight, rank):
    w = np.asarray(weight, dtype=np.float32)
    if w.ndim != 2 or not w.size or not np.isfinite(w).all():
        raise ValueError("finite nonempty rank-2 weight required")
    if type(rank) is not int or not 1 <= rank <= min(w.shape):
        raise ValueError("invalid rank")
    u, s, vh = np.linalg.svd(w, full_matrices=False)
    return (np.ascontiguousarray(u[:, :rank] * s[:rank], dtype=np.float32),
            np.ascontiguousarray(vh[:rank, :], dtype=np.float32))


def apply_direct(inputs, left, right):
    x = np.asarray(inputs, dtype=np.float32)
    a = np.asarray(left, dtype=np.float32)
    b = np.asarray(right, dtype=np.float32)
    if (x.ndim != 2 or a.ndim != 2 or b.ndim != 2 or
            x.shape[1] != b.shape[1] or a.shape[1] != b.shape[0] or
            min(*x.shape, *a.shape, *b.shape) <= 0 or
            not all(np.isfinite(z).all() for z in (x, a, b))):
        raise ValueError("invalid factorized expert operands")
    return (x @ b.T) @ a.T


def factor_storage_bytes(left, right, itemsize=2):
    a, b = np.asarray(left), np.asarray(right)
    if (a.ndim != 2 or b.ndim != 2 or a.shape[1] != b.shape[0] or
            min(*a.shape, *b.shape) <= 0 or type(itemsize) is not int or itemsize <= 0):
        raise ValueError("invalid factor storage")
    return (a.size + b.size) * itemsize
