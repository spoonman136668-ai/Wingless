"""GCEN test-only reconstruction controls; no model inference."""
import numpy as np

def original(weights):
    return np.asarray(weights, dtype=np.float32).copy()

def svd(weights, rank):
    w = np.asarray(weights, dtype=np.float32)
    if w.ndim != 2 or not 1 <= rank <= min(w.shape):
        raise ValueError('invalid shape or rank')
    u, s, vh = np.linalg.svd(w, full_matrices=False)
    return (u[:, :rank] * s[:rank]) @ vh[:rank, :]

def shared_basis(weights, rank):
    w = np.asarray(weights, dtype=np.float32)
    if w.ndim != 3:
        raise ValueError('expected [experts,out,in]')
    e, m, n = w.shape
    if not 1 <= rank <= min(e*m,n):
        raise ValueError('invalid rank')
    _, _, vh = np.linalg.svd(w.reshape(e*m,n),full_matrices=False)
    basis = vh[:rank]
    return (w @ basis.T) @ basis

def storage_bytes(shape, rank, kind, itemsize=2):
    e,m,n = shape
    if min(shape) <= 0 or itemsize <= 0:
        raise ValueError('invalid shape')
    if kind == 'original': return e*m*n*itemsize
    if kind == 'svd' and 1 <= rank <= min(m,n): return e*rank*(m+n+1)*itemsize
    if kind == 'shared' and 1 <= rank <= min(e*m,n): return (e*m*rank+rank*n)*itemsize
    raise ValueError('invalid kind or rank')
