"""Larger real-weight expert reconstruction development screen.

Only measures reconstruction error on authentic selected expert weight windows.
No inference, no task-quality claim, no protected held-out execution.
"""
import time
import numpy as np

SEEDS = (11, 23, 37, 53, 71)
EXPERTS, ROWS, COLS = 4, 128, 256
SVD_RANK, SHARED_RANK = 14, 28


def compare_window(weights, clock=time.perf_counter):
    w = np.asarray(weights, dtype=np.float32)
    if w.shape != (EXPERTS, ROWS, COLS) or not np.isfinite(w).all():
        raise ValueError("expected four finite 128x256 expert windows")
    energy = float(np.sum(w.astype(np.float64) ** 2))
    if energy <= 0:
        raise ValueError("zero-energy weights")
    original_bytes = int(w.size * 2)
    svd_bytes = EXPERTS * SVD_RANK * (ROWS + COLS + 1) * 2
    shared_bytes = (EXPERTS * ROWS * SHARED_RANK + SHARED_RANK * COLS) * 2
    # Frozen near-equal stored FP16 factor budget; all factors counted.
    if max(svd_bytes, shared_bytes) > original_bytes // 4:
        raise ValueError("4x storage threshold not met")
    if abs(svd_bytes - shared_bytes) > 0.01 * max(svd_bytes, shared_bytes):
        raise ValueError("factor budgets are not matched within 1%")

    start = clock()
    independent = np.empty_like(w)
    for i in range(EXPERTS):
        u, singular, vh = np.linalg.svd(w[i], full_matrices=False)
        independent[i] = (u[:, :SVD_RANK] * singular[:SVD_RANK]) @ vh[:SVD_RANK]
    svd_seconds = max(0.0, float(clock() - start))

    start = clock()
    _, _, vh = np.linalg.svd(w.reshape(EXPERTS * ROWS, COLS), full_matrices=False)
    basis = vh[:SHARED_RANK]
    shared = (w @ basis.T) @ basis
    shared_seconds = max(0.0, float(clock() - start))
    errors = []
    per_expert = []
    for i in range(EXPERTS):
        denom = float(np.sum(w[i].astype(np.float64) ** 2))
        if denom <= 0:
            raise ValueError("zero-energy expert")
        a = float(np.sum((w[i].astype(np.float64) - independent[i].astype(np.float64)) ** 2) / denom)
        b = float(np.sum((w[i].astype(np.float64) - shared[i].astype(np.float64)) ** 2) / denom)
        per_expert.append({"expert": i, "svd_error": a, "shared_error": b})
    for approximation in (independent, shared):
        errors.append(float(np.sum((w.astype(np.float64) - approximation.astype(np.float64)) ** 2) / energy))
    # These are explicit, conservative accounting components, NOT measured GPU VRAM:
    # input float32 + reconstructed float32 output + full SVD float32 U/S/Vh.
    input_f32 = int(w.nbytes)
    output_f32 = input_f32
    independent_svd_full_workspace_lower_bound = int(
        EXPERTS * (ROWS * ROWS + ROWS + ROWS * COLS) * 4
    )
    shared_svd_full_workspace_lower_bound = int(
        (EXPERTS * ROWS * COLS + COLS + COLS * COLS) * 4
    )
    return {
        "shape": list(w.shape),
        "original_fp16_bytes": original_bytes,
        "svd_factor_fp16_bytes": svd_bytes,
        "shared_factor_fp16_bytes": shared_bytes,
        "svd_rank": SVD_RANK,
        "shared_rank": SHARED_RANK,
        "svd_normalized_squared_error": errors[0],
        "shared_normalized_squared_error": errors[1],
        "per_expert": per_expert,
        "svd_reconstruction_seconds": svd_seconds,
        "shared_reconstruction_seconds": shared_seconds,
        "input_float32_bytes": input_f32,
        "output_float32_bytes": output_f32,
        "svd_workspace_lower_bound_bytes": independent_svd_full_workspace_lower_bound,
        "shared_workspace_lower_bound_bytes": shared_svd_full_workspace_lower_bound,
        "scope": "SELECTED_128x256_WEIGHT_WINDOWS_NOT_FULL_MODEL_OR_INFERENCE",
        "memory_caveat": "CPU workspace lower bounds omit BLAS/LAPACK temporaries; do not report as peak RSS or VRAM",
    }
