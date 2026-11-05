"""CNC-G9G: outlier robustness and frozen mixed-precision capacity screen.

Development-only synthetic diagnostic; no original weights, real activations,
GPU model, or scientific acceptance. Reuses qualified G9F operators and
G9D exact-index-compatible fixed masks. This evaluates stress, not task quality.
"""
import json
import math
import numpy as np
from g9d_mixed_precision_budget import mixed_report
from g9b_tensor_budget import FROZEN_CONFIG, expected_tensors
from packed_companded_int2 import pack_companded_int2, apply_companded_int2, stored_bytes as compbytes
from packed_int2 import pack_int2, apply_packed_int2, stored_bytes as uniformbytes
from packed_int4 import pack_int4, apply_packed_int4, stored_bytes as int4bytes

SEEDS = (101, 211, 307)  # Independent of G9E/G9F original-weight seeds.
FAMILIES = ("gaussian", "student_t3", "sparse_outliers")
BLOCKS = (128, 256)
SHAPE = (128, 256)
BATCH = 32
SOURCE_SHA = "d" * 40  # Offline synthetic-only fixture identity, NOT a repo receipt.


def indexed_fixture():
    shapes = expected_tensors(dict(FROZEN_CONFIG))
    return {
        "weight_map": {key: f"model-{(i % 16)+1:05d}-of-00016.safetensors"
                       for i, key in enumerate(sorted(shapes))},
        "metadata": {"total_size": 2 * sum(math.prod(shape) for shape, kind in shapes.values())},
    }


def synthetic_weights(seed, family):
    if type(seed) is not int or seed not in SEEDS or family not in FAMILIES:
        raise ValueError("unfrozen G9G development conditions")
    rng = np.random.default_rng(seed)
    weights = rng.standard_normal(SHAPE).astype(np.float32)
    if family == "student_t3":
        # Unit variance in expectation, heavy tails.
        weights = (rng.standard_t(3, SHAPE) / np.sqrt(3)).astype(np.float32)
    elif family == "sparse_outliers":
        # Fixed probability/scale; no empirical calibration.
        hits = rng.random(SHAPE) < 0.01
        weights[hits] *= 16
    if not np.isfinite(weights).all() or not np.any(weights):
        raise ValueError("nonfinite or empty synthetic weights")
    return weights


def compare_case(seed, family, block):
    if type(block) is not int or block not in BLOCKS:
        raise ValueError("unfrozen G9G block size")
    w = synthetic_weights(seed, family)
    # Different stream from weight RNG; no native model activations.
    x = np.random.default_rng(seed + 60000).standard_normal((BATCH, SHAPE[1])).astype(np.float32)
    ref = x @ w.T
    energy = float(np.sum(ref.astype(np.float64) ** 2))
    if not np.isfinite(energy) or energy <= 0:
        raise ValueError("zero output energy")
    p4, s4 = pack_int4(w, block=64)
    pu, su = pack_int2(w, block=block)
    pc, sc = pack_companded_int2(w, block=block)
    ubytes = uniformbytes(pu, su, SHAPE[1], block=block)
    cbytes = compbytes(pc, sc, SHAPE[1], block=block)
    if ubytes != cbytes:
        raise ValueError("INT2 competitor storage is not matched")
    def error(pred):
        arr = np.asarray(pred, dtype=np.float64) - ref.astype(np.float64)
        return float(np.sum(arr**2) / energy)
    results = {
        "seed": seed, "synthetic_family": family, "block": block,
        "bf16_reference_bytes": w.size * 2, "uniform_int2_bytes": ubytes,
        "companded_int2_bytes": cbytes,
        "int4_control_bytes": int4bytes(p4, s4, SHAPE[1], block=64),
        "uniform_int2_relative_output_squared_error": error(apply_packed_int2(x, pu, su, block=block)),
        "companded_int2_relative_output_squared_error": error(apply_companded_int2(x, pc, sc, block=block)),
        "int4_relative_output_squared_error": error(apply_packed_int4(x, p4, s4, block=64)),
    }
    if not all(math.isfinite(v) for k,v in results.items() if k.endswith("_error")):
        raise ValueError("nonfinite diagnostic")
    return results


def analyze():
    rows = [compare_case(seed, fam, block)
            for fam in FAMILIES for seed in SEEDS for block in BLOCKS]
    scenarios = mixed_report(dict(FROZEN_CONFIG), indexed_fixture(), source_sha=SOURCE_SHA)["scenarios"]
    budget = [
        {"name": v["name"], "protected_experts": v["protected_expert_count"],
         "remaining_bytes": v["remaining_bytes"],
         "analytical_fit": v["analytical_fit_under_assumptions"]}
        for v in scenarios
    ]
    family_metrics = []
    for family in FAMILIES:
        for block in BLOCKS:
            subset = [r for r in rows if r["synthetic_family"] == family and r["block"] == block]
            family_metrics.append({
                "family": family, "block": block,
                "companded_mean_relative_output_squared_error":
                    sum(r["companded_int2_relative_output_squared_error"] for r in subset) / len(subset),
                "companded_worst_relative_output_squared_error":
                    max(r["companded_int2_relative_output_squared_error"] for r in subset),
                "uniform_mean_relative_output_squared_error":
                    sum(r["uniform_int2_relative_output_squared_error"] for r in subset) / len(subset),
                "int4_mean_relative_output_squared_error":
                    sum(r["int4_relative_output_squared_error"] for r in subset) / len(subset),
                "companded_better_than_uniform_count":
                    sum(r["companded_int2_relative_output_squared_error"] <
                        r["uniform_int2_relative_output_squared_error"] for r in subset),
            })
    return {
        "schema": "wingless.cnc.g9g.synthetic-stress-mixed-capacity.v1",
        "status": "DEVELOPMENT_SYNTHETIC_ONLY",
        "seeds": list(SEEDS), "families": list(FAMILIES),
        "blocks": list(BLOCKS), "case_count": len(rows),
        "cases": rows, "family_metrics": family_metrics,
        "mixed_precision_capacity": budget,
        "source_identity": "SYNTHETIC_FIXTURE_NOT_ORIGINAL_INDEX",
        "original_qwen_weights": False, "authentic_qwen_activations": False,
        "model_forward": False, "gpu_residency": "UNVERIFIED",
        "scientific_acceptance": "NOT_GRANTED",
        "protected_heldout_used": False,
    }


if __name__ == "__main__":
    print("CNC_G9G_SYNTHETIC_SCREEN " + json.dumps(analyze(), sort_keys=True))
