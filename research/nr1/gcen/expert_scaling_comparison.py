"""GCEN development-only expert-sharing scaling screen (not model inference).

The frozen test varies only the number of genuinely distinct expert matrices;
factor-byte budgets are matched to within 1.5 percent for each group size.
"""
import time
import numpy as np

SEEDS = (11, 23, 37, 53, 71)
COUNTS = (4, 8, 16)
ROWS, COLS = 128, 256
SVD_RANK = 14


def rank_for_matched_bytes(experts):
    if experts not in COUNTS:
        raise ValueError("unfrozen expert count")
    independent = experts * SVD_RANK * (ROWS + COLS + 1) * 2
    basis_units = experts * ROWS + COLS
    # Freeze the closest integer rank by exact storage; no error-based tuning.
    return min(range(1, COLS + 1), key=lambda rank: (abs(2 * basis_units * rank - independent), rank))


def evaluate_group(weights):
    w = np.asarray(weights, dtype=np.float32)
    if w.ndim != 3 or w.shape[0] not in COUNTS or w.shape[1:] != (ROWS, COLS) or not np.isfinite(w).all():
        raise ValueError("invalid genuine expert group")
    experts = w.shape[0]
    shared_rank = rank_for_matched_bytes(experts)
    original_bytes = int(w.size * 2)
    independent_bytes = experts * SVD_RANK * (ROWS + COLS + 1) * 2
    shared_bytes = (experts * ROWS + COLS) * shared_rank * 2
    if not (independent_bytes <= original_bytes // 4 and shared_bytes <= original_bytes // 4):
        raise ValueError("4x target failed")
    if abs(independent_bytes - shared_bytes) > max(independent_bytes, shared_bytes) * .015:
        raise ValueError("factor storage not matched within 1.5 percent")
    denominator = float(np.sum(w.astype(np.float64)**2))
    if denominator <= 0: raise ValueError("empty energy")

    start=time.perf_counter()
    independent=np.empty_like(w)
    for i in range(experts):
        u,s,vh=np.linalg.svd(w[i],full_matrices=False)
        independent[i]=(u[:,:SVD_RANK]*s[:SVD_RANK])@vh[:SVD_RANK]
    independent_sec=time.perf_counter()-start

    start=time.perf_counter()
    _,_,vh=np.linalg.svd(w.reshape(experts*ROWS,COLS),full_matrices=False)
    basis=vh[:shared_rank]
    shared=(w@basis.T)@basis
    shared_sec=time.perf_counter()-start

    sse=lambda estimate: float(np.sum((w.astype(np.float64)-estimate.astype(np.float64))**2)/denominator)
    return dict(
        experts=experts,
        source_window_shape=list(w.shape),
        independent_rank=SVD_RANK,
        shared_rank=shared_rank,
        original_fp16_bytes=original_bytes,
        independent_factor_fp16_bytes=independent_bytes,
        shared_factor_fp16_bytes=shared_bytes,
        independent_normalized_squared_error=sse(independent),
        shared_normalized_squared_error=sse(shared),
        independent_cpu_seconds=float(independent_sec),
        shared_cpu_seconds=float(shared_sec),
        peak_memory="UNMEASURED",
        purpose="EXPERT_SHARING_SCALING_DEVELOPMENT_ONLY",
        limitation="128x256 sampled windows, factor-only budget; no full-model, task-quality, peak-VRAM, or throughput claim",
    )
