"""G9K: exact Qwen tensor census based *analytical* low-rank factor budget."""
import math
from g9b_tensor_budget import FROZEN_CONFIG, expected_tensors, make_report
from lowrank_residual_int2 import RANKS, correction_bytes_for_shape


def fixture_index(tensors):
    return {
        "weight_map": {name: f"model-{i%16+1:05d}-of-00016.safetensors"
                       for i, name in enumerate(sorted(tensors))},
        "metadata": {"total_size": sum(2*math.prod(shape)
                   for shape, _ in tensors.values())},
    }


def report():
    tensors = expected_tensors(dict(FROZEN_CONFIG))
    base = make_report(dict(FROZEN_CONFIG), fixture_index(tensors),
                       source_sha="f"*40)
    candidate = next(x for x in base["scenarios"]
                     if x["name"] == "int2_block256_fp16")
    cases = []
    for rank in RANKS:
        delta = sum(correction_bytes_for_shape(shape, rank=rank)
                    for shape, kind in tensors.values() if kind == "experts")
        remaining = candidate["remaining_bytes"] - delta
        cases.append({
            "rank": rank, "additional_fp16_factors_and_headers_bytes": delta,
            "remaining_after_hypothetical_workspace_bytes": remaining,
            "analytical_12_gib_fit": remaining >= 0,
            "actual_gpu_memory": "UNVERIFIED",
        })
    return {
        "schema": "wingless.cnc.g9k.lowrank-analytical.v1",
        "scope": "DEVELOPMENT_ONLY",
        "model": "Qwen/Qwen3-30B-A3B",
        "parameter_total": base["parameter_total"],
        "indexed_tensor_count": base["indexed_tensor_count"],
        "baseline_remaining_bytes": candidate["remaining_bytes"],
        "cases": cases,
        "source_index_custody": "RECONSTRUCTED_FROM_PRIOR_G9B_AUTHENTICATED_INDEX",
        "model_forward": False,
        "original_qwen_weight_windows_used": False,
        "real_model_activations": False,
        "gpu_residency": "UNVERIFIED",
        "protected_heldout_used": False,
        "scientific_acceptance": "NOT_GRANTED",
    }


if __name__ == "__main__":
    import json
    print("CNC_G9K_FROZEN_ANALYTICAL "+json.dumps(report(), sort_keys=True))
