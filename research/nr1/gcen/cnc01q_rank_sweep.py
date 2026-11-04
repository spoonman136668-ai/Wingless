"""CNC-01Q-G6 controlled synthetic rank/batch crossover diagnostic.

This isolates storage, reconstruction accuracy, and CPU operator costs. No Qwen
weights, model forward, GPU operation, or scientific acceptance is implied.
"""
import json
import platform
import time

import numpy as np

from direct_expert import apply_direct, factor_storage_bytes

ROWS, COLS = 128, 256
RANKS = (8, 14, 32, 64, 96)
BATCHES = (1, 32, 128)
SEED = 3101


def _median_seconds(fn, repeats):
    for _ in range(2):
        fn()
    times = []
    for _ in range(repeats):
        t0 = time.perf_counter()
        fn()
        times.append(time.perf_counter() - t0)
    return float(np.median(times))


def compare(repeats=3):
    if type(repeats) is not int or not 1 <= repeats <= 30:
        raise ValueError("invalid repeat budget")
    rng = np.random.default_rng(SEED)
    families = {
        "full_rank_random": rng.standard_normal((ROWS, COLS)).astype(np.float32),
        "rank8_positive_control": None,
    }
    left = rng.standard_normal((ROWS, 8)).astype(np.float32)
    right = rng.standard_normal((8, COLS)).astype(np.float32)
    families["rank8_positive_control"] = (left @ right).astype(np.float32)
    results = []
    for family, w in families.items():
        u, singular, vh = np.linalg.svd(w, full_matrices=False)
        weight_energy = float(np.sum(w.astype(np.float64)**2))
        if weight_energy <= 0:
            raise ValueError("zero energy reference")
        x_by_batch = {batch: rng.standard_normal((batch, COLS)).astype(np.float32) for batch in BATCHES}
        for rank in RANKS:
            a = np.ascontiguousarray(u[:, :rank] * singular[:rank], dtype=np.float32)
            b = np.ascontiguousarray(vh[:rank, :], dtype=np.float32)
            dense = a @ b
            weight_error = float(np.sum((w.astype(np.float64) - dense.astype(np.float64))**2) / weight_energy)
            for batch in BATCHES:
                x = x_by_batch[batch]
                y = x @ w.T
                y_direct = apply_direct(x, a, b)
                y_dense = x @ dense.T
                # Independent float64 algebraic oracle separates associativity
                # rounding from genuine factor/operator divergence.
                x64, a64, b64 = x.astype(np.float64), a.astype(np.float64), b.astype(np.float64)
                oracle_direct = (x64 @ b64.T) @ a64.T
                oracle_reconstructed = x64 @ (a64 @ b64).T
                np.testing.assert_allclose(oracle_direct, oracle_reconstructed, rtol=1e-12, atol=1e-11)
                scale = float(np.linalg.norm(oracle_direct))
                if scale <= 0:
                    raise ValueError("zero energy factorized output")
                f32_direct_error = float(np.linalg.norm(y_direct.astype(np.float64) - oracle_direct) / scale)
                f32_dense_error = float(np.linalg.norm(y_dense.astype(np.float64) - oracle_direct) / scale)
                if f32_direct_error > 1e-5 or f32_dense_error > 1e-5:
                    raise AssertionError("float32 kernel diverges from float64 factorized oracle")
                energy = float(np.sum(y.astype(np.float64)**2))
                if energy <= 0:
                    raise ValueError("zero energy output")
                output_error = float(np.sum((y_direct.astype(np.float64) - y.astype(np.float64))**2) / energy)
                results.append({
                    "family": family, "rank": rank, "batch": batch,
                    "factor_fp16_bytes": factor_storage_bytes(a, b),
                    "dense_fp16_bytes": ROWS * COLS * 2,
                    "relative_weight_squared_error": weight_error,
                    "relative_output_squared_error": output_error,
                    "direct_median_cpu_seconds": _median_seconds(lambda: apply_direct(x, a, b), repeats),
                    "reconstructed_median_cpu_seconds": _median_seconds(lambda: x @ dense.T, repeats),
                })
    return {
        "status": "DEVELOPMENT_ONLY_SYNTHETIC_RANK_SWEEP",
        "scope": "NO_REAL_QWEN_NO_GPU_NO_INFERENCE_NO_ACCEPTANCE",
        "numpy_version": np.__version__, "platform": platform.platform(),
        "seed": SEED, "ranks": list(RANKS), "batches": list(BATCHES),
        "repeats": repeats, "float32_runtime": True,
        "serialized_factor_bytes_are_theoretical_fp16": True,
        "peak_rss": "UNMEASURED", "peak_vram": "UNMEASURED",
        "results": results,
    }


if __name__ == "__main__":
    print("CNC01Q_RANK_SWEEP_JSON " + json.dumps(compare(), sort_keys=True))
