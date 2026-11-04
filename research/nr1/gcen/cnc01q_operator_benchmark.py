"""CNC-01Q synthetic operator comparison, development only.

All representations originate from the same independently factorized SVD.
Compare inference kernels after setup; report preparation separately.
Not a Qwen checkpoint test, GPU measurement or scientific acceptance.
"""
import json
import time

import numpy as np

from direct_expert import apply_direct, factorize, factor_storage_bytes

SHAPE = (128, 256)
RANK = 14
BATCH = 32
SEED = 23


def measure(repeats=100, warmups=10):
    if type(repeats) is not int or not 1 <= repeats <= 10000:
        raise ValueError("invalid repeats")
    if type(warmups) is not int or not 0 <= warmups <= 1000:
        raise ValueError("invalid warmups")
    rng = np.random.default_rng(SEED)
    w = rng.standard_normal(SHAPE).astype(np.float32)
    x = rng.standard_normal((BATCH, SHAPE[1])).astype(np.float32)
    start = time.perf_counter()
    a, b = factorize(w, RANK)
    prep_seconds = time.perf_counter() - start
    dense_reconstructed = a @ b
    reference = x @ w.T
    candidates = {
        "dense_original": lambda: x @ w.T,
        "dense_reconstructed": lambda: x @ dense_reconstructed.T,
        "direct_factors": lambda: apply_direct(x, a, b),
    }
    results = {}
    for name, op in candidates.items():
        for _ in range(warmups):
            op()
        timings = []
        for _ in range(repeats):
            start = time.perf_counter()
            y = op()
            timings.append(time.perf_counter() - start)
        delta = y.astype(np.float64) - reference.astype(np.float64)
        energy = float(np.sum(reference.astype(np.float64)**2))
        results[name] = {
            "median_seconds": float(np.median(timings)),
            "relative_output_squared_error": float(np.sum(delta**2) / energy),
        }
    np.testing.assert_allclose(
        candidates["direct_factors"](), candidates["dense_reconstructed"](),
        atol=2e-5, rtol=2e-5,
    )
    return {
        "status": "SYNTHETIC_DEVELOPMENT_DIAGNOSTIC_ONLY",
        "seed": SEED, "source_shape": list(SHAPE), "rank": RANK,
        "batch": BATCH, "repeats": repeats, "warmups": warmups,
        "factor_preparation_seconds": float(prep_seconds),
        "dense_fp16_storage_bytes": int(w.size * 2),
        "factor_fp16_storage_bytes": factor_storage_bytes(a, b),
        "runtime_peak_rss_bytes": "UNMEASURED",
        "runtime_peak_vram_bytes": "UNMEASURED",
        "numpy_runtime_dtype": "float32",
        "timings_are_cpu_only": True,
        "results": results,
        "limitation": "No original Qwen tensors or model inference; timing excludes SVD setup but includes Python dispatch; factor storage assumes FP16 serialization, not resident float32 arrays.",
    }


if __name__ == "__main__":
    print(json.dumps(measure(), sort_keys=True))
