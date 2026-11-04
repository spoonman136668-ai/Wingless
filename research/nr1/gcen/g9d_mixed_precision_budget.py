"""CNC-G9D: exact-index-accounted mixed-precision Qwen residency screen.

Analytical only: no original weights, router observations, model activations,
kernel timing or GPU allocations. Preserve G9B's BF16 nonexperts, KV cache,
12 GiB ceiling and illustrative unmeasured 1 GiB workspace reserve.
Frozen deterministic 4-bit expert masks are not sensitivity rankings.
"""
from collections import defaultdict
import math
from g9b_tensor_budget import (
    EXPECTED_TENSORS, CONTEXT, FROZEN_CONFIG, VRAM_CEILING,
    WORKSPACE_RESERVE, expected_tensors, make_report,
)

LAYERS = 48
EXPERTS_PER_LAYER = 128
PROJECTIONS = ("gate_proj", "up_proj", "down_proj")
BASE = (2, 256, 2)  # 2-bit, block-256, FP16 scales
PROTECTED = (4, 64, 4)  # G8 style 4-bit, block-64, FP32 scales
MASKS = {
    "all_int2_baseline": (),
    "first_layer_int4": (0,),
    "last_layer_int4": (47,),
    "both_boundary_layers_int4": (0, 47),
}
EXPERT_ID_MASKS = {
    "first_expert_each_layer_int4": 1,
    "first_two_experts_each_layer_int4": 2,
}


def _tensor_bytes(shape, spec):
    bits, block, scale_bytes = spec
    if len(shape) != 2 or shape[1] % block:
        raise ValueError("incompatible expert tensor shape")
    n = math.prod(shape)
    return (n * bits + 7)//8 + shape[0]*(shape[1]//block)*scale_bytes + 24


def mixed_report(config, index, *, source_sha):
    # G9B establishes identity, indexed tensor count, canonical source bytes,
    # and budget assumptions. Abort on source drift rather than estimating.
    baseline = make_report(config, index, source_sha=source_sha)
    if baseline["indexed_tensor_count"] != EXPECTED_TENSORS:
        raise ValueError("drift in original model inventory")
    shapes = expected_tensors(config)
    reference = next(c for c in baseline["scenarios"]
                     if c["name"] == "int2_block256_fp16")
    all_int2 = reference["total_expert_bytes"]
    nonexpert = reference["nonexpert_bf16_bytes"]
    kv = reference["bf16_kv_bytes"]
    reserve = reference["hypothetical_workspace_reserve_bytes"]

    # Precompute per-tensor increment. This is a storage delta, not a
    # quality-based decision; no expert activations or router traces are read.
    by_layer = defaultdict(int)
    by_expert = defaultdict(int)
    count = 0
    base_total = 0
    for name, (shape, kind) in shapes.items():
        if kind != "experts":
            continue
        parts = name.split(".")
        if len(parts) != 7 or parts[0:2] != ["model", "layers"] or parts[3:5] != ["mlp", "experts"] or parts[6] not in {p+".weight" for p in PROJECTIONS}:
            raise ValueError("unexpected authentic expert key structure")
        layer, expert = int(parts[2]), int(parts[5])
        if not (0 <= layer < LAYERS and 0 <= expert < EXPERTS_PER_LAYER):
            raise ValueError("unfrozen expert identity")
        plain = _tensor_bytes(shape, BASE)
        extra = _tensor_bytes(shape, PROTECTED) - plain
        if extra <= 0:
            raise ValueError("protected representation did not increase storage")
        base_total += plain
        by_layer[layer] += extra
        by_expert[(layer, expert)] += extra
        count += 1
    if count != 18432 or base_total != all_int2 or len(by_layer) != LAYERS or len(by_expert) != LAYERS*EXPERTS_PER_LAYER:
        raise ValueError("mixed precision coverage/reconciliation failed")

    cases = []
    for name, layers in MASKS.items():
        additional = sum(by_layer[layer] for layer in layers)
        protected_experts = len(layers)*EXPERTS_PER_LAYER
        cases.append((name, additional, protected_experts, "FROZEN_LAYER_MASK_NOT_SENSITIVITY"))
    for name, n in EXPERT_ID_MASKS.items():
        additional = sum(by_expert[(layer, expert)]
                         for layer in range(LAYERS) for expert in range(n))
        cases.append((name, additional, LAYERS*n, "FROZEN_EXPERT_ID_MASK_NOT_ROUTER_RANK"))
    results = []
    for name, delta, protected_experts, selection in cases:
        total_expert_bytes = all_int2 + delta
        planned = nonexpert + total_expert_bytes + kv + reserve
        results.append({
            "name": name,
            "protected_expert_count": protected_experts,
            "total_expert_count": LAYERS*EXPERTS_PER_LAYER,
            "protected_expert_precision": "INT4_BLOCK64_FP32_SCALES",
            "remaining_expert_precision": "INT2_BLOCK256_FP16_SCALES",
            "selection": selection,
            "extra_expert_bytes_vs_full_int2": delta,
            "total_expert_bytes": total_expert_bytes,
            "nonexpert_bf16_bytes": nonexpert,
            "kv_cache_bytes": kv,
            "hypothetical_workspace_reserve_bytes": reserve,
            "total_planned_bytes": planned,
            "remaining_bytes": VRAM_CEILING - planned,
            "analytical_fit_under_assumptions": planned <= VRAM_CEILING,
            "quality_or_router_sensitivity": "UNVERIFIED",
            "gpu_runtime_or_residency": "UNVERIFIED",
        })
    if results[0]["remaining_bytes"] != reference["remaining_bytes"]:
        raise ValueError("reference G9B budget mismatch")
    return {
        "schema": "wingless.cnc.g9d.mixed-precision-analytical.v1",
        "status": "DEVELOPMENT_ONLY_FROZEN_MIXED_STORAGE_SCREEN",
        "source_sha": source_sha,
        "model": baseline["checkpoint_model"],
        "revision": baseline["checkpoint_revision"],
        "source_index_tensor_count": baseline["indexed_tensor_count"],
        "source_index_bf16_bytes": baseline["indexed_bf16_tensor_bytes"],
        "vram_ceiling_bytes": VRAM_CEILING,
        "context_tokens": CONTEXT,
        "source_assumptions": baseline["assumptions"],
        "scenarios": results,
        "full_model_forward": False,
        "actual_gpu_measurement": False,
        "real_expert_sensitivity_study": False,
        "protected_heldout_used": False,
        "scientific_acceptance": "NOT_GRANTED",
        "deployment": False,
    }
