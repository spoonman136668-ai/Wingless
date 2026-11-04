"""Small real-weight GCEN reconstruction comparison; no model inference.

Development-only: compare independent SVD rank 5 and shared right-basis
rank 10 on four expert matrices, 32 x 64 tiles, at closely matched FP16
factor-storage budgets. No scientific or full-model quality claims.
"""
import numpy as np

SEEDS = (11, 23, 37, 53, 71)
EXPERTS = 4
ROWS = 32
COLS = 64
SVD_RANK = 5
SHARED_RANK = 10


def compare_tile(weights):
    w = np.asarray(weights, dtype=np.float32)
    if w.shape != (EXPERTS, ROWS, COLS) or not np.isfinite(w).all():
        raise ValueError("expected four finite 32x64 expert tiles")
    # Independent factor storage counts U, singular values, and Vh.
    svd_bytes = EXPERTS * SVD_RANK * (ROWS + COLS + 1) * 2
    # Shared factor storage counts per-expert coefficients and the shared basis.
    shared_bytes = (EXPERTS * ROWS * SHARED_RANK + SHARED_RANK * COLS) * 2
    original_bytes = w.size * 2
    if not (svd_bytes <= original_bytes // 4 and shared_bytes <= original_bytes // 4):
        raise ValueError("the frozen 4x expert-only compression target is not met")

    independent = np.empty_like(w)
    for i in range(EXPERTS):
        u, singular, vh = np.linalg.svd(w[i], full_matrices=False)
        independent[i] = (u[:, :SVD_RANK] * singular[:SVD_RANK]) @ vh[:SVD_RANK]
    _, _, vh = np.linalg.svd(w.reshape(EXPERTS * ROWS, COLS), full_matrices=False)
    basis = vh[:SHARED_RANK]
    shared = (w @ basis.T) @ basis
    denominator = float(np.sum(w.astype(np.float64) ** 2))
    if denominator <= 0:
        raise ValueError("zero-energy tile")
    err_ind = float(np.sum((w.astype(np.float64) - independent.astype(np.float64)) ** 2) / denominator)
    err_shared = float(np.sum((w.astype(np.float64) - shared.astype(np.float64)) ** 2) / denominator)
    return {
        "original_fp16_bytes": original_bytes,
        "svd_factor_fp16_bytes": svd_bytes,
        "shared_factor_fp16_bytes": shared_bytes,
        "svd_rank": SVD_RANK,
        "shared_rank": SHARED_RANK,
        "svd_relative_squared_error": err_ind,
        "shared_relative_squared_error": err_shared,
        "minimum_compression_ratio": min(original_bytes / svd_bytes, original_bytes / shared_bytes),
        "qualifier": "TILE_ONLY_FACTOR_STORAGE_NO_RUNTIME_MEMORY_OR_MODEL_QUALITY_CLAIM",
    }
