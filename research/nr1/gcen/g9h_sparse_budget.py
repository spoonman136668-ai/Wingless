"""CNC-G9H: exact key census and fixed sparse correction storage ceiling.

Analytical-only; earlier G9B source index authentication is independent
evidence. This procedure builds a reconstructed index fixture for tests and
does NOT assert it just refetched original model metadata.
"""
import math
from g9b_tensor_budget import (
    FROZEN_CONFIG, VRAM_CEILING, expected_tensors, make_report
)
from sparse_residual_int2 import BUDGET_K, correction_bytes_for_shape

SOURCE_SHA = "e"*40  # Fixture marker, not original source identity.


def fixture_index(shapes):
    return {
        "weight_map": {
            name: f"model-{i%16+1:05d}-of-00016.safetensors"
            for i, name in enumerate(sorted(shapes))},
        "metadata": {
            "total_size": 2 * sum(math.prod(shape) for shape, _ in shapes.values())},
    }


def sparse_budget():
    shapes = expected_tensors(dict(FROZEN_CONFIG))
    base = make_report(dict(FROZEN_CONFIG), fixture_index(shapes),
                       source_sha=SOURCE_SHA)
    candidate = next(x for x in base["scenarios"]
                     if x["name"] == "int2_block256_fp16")
    cases = []
    for k in BUDGET_K:
        delta = sum(correction_bytes_for_shape(shape, k=k)
                    for shape, kind in shapes.values() if kind == "experts")
        remain = candidate["remaining_bytes"] - delta
        cases.append({
            "corrections_per_row": k,
            "extra_sparse_index_value_offset_bytes": delta,
            "model_total_planned_bytes": VRAM_CEILING - remain,
            "remaining_after_reserve_bytes": remain,
            "analytical_fit": remain >= 0,
            "source": "RECONSTRUCTED_INDEX_FIXTURE_FROM_PRIOR_G9B_AUTHENTICATED_ARCH",
            "gpu_memory_measured": False,
            "model_behavior_verified": False,
        })
    return {
        "schema": "wingless.cnc.g9h.frozen-sparse-budget.v1",
        "scope": "DEVELOPMENT_ANALYTICAL_ONLY",
        "parameter_total": base["parameter_total"],
        "indexed_tensor_count": base["indexed_tensor_count"],
        "baseline_remaining_bytes": candidate["remaining_bytes"],
        "12_gib_vram_bytes": VRAM_CEILING,
        "context_tokens": 8192,
        "hypothetical_unmeasured_workspace_bytes": 1<<30,
        "cases": cases,
        "actual_original_model_index_fetched": False,
        "original_qwen_activation_probes": False,
        "scientific_acceptance": "NOT_GRANTED",
    }


if __name__ == "__main__":
    import json
    print("CNC_G9H_SPARSE_BUDGET " + json.dumps(sparse_budget(), sort_keys=True))
