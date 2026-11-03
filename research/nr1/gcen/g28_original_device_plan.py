"""G28 original 30.53B meta model CPU/disk device-map and capacity plan.

Uses exact pinned metadata and the original pinned Transformers/Accelerate
sources; NO checkpoint weight bytes, live inference, GPU, or activation capture.
This is a static placement estimate, not a qualified full-model loader.
"""
import argparse
import json
from pathlib import Path
import shutil
import time

from g27_original_meta_mapping import check_full_key_mapping
from out_of_core_preflight import pinned_metadata, analyze
from runtime_feasibility import (
    MODEL, REVISION, PINNED_ORIGINAL_SHARD_BYTES, check_manifest,
)

GIB = 1024**3
CPU_BUDGET_GIB = (4, 8, 12)
DISK_BUDGET_BYTES = 256 * GIB
SEPARATE_WORKSPACE_RESERVE = 8 * GIB
NO_SPLIT = ("Qwen3MoeDecoderLayer",)


def resolve_assignment(parameter_name, device_map):
    """Match longest parent device map key for a checkpoint tensor name."""
    if not isinstance(parameter_name, str) or not parameter_name:
        raise ValueError("invalid parameter name")
    if not isinstance(device_map, dict) or not device_map:
        raise ValueError("empty device map")
    parent = parameter_name
    while True:
        if parent in device_map:
            target = device_map[parent]
            if target not in ("cpu", "disk"):
                raise ValueError("G28 host-only plan unexpectedly targets GPU/other device")
            return target
        if "." not in parent:
            break
        parent = parent.rsplit(".", 1)[0]
    if "" in device_map:
        target = device_map[""]
        if target in ("cpu", "disk"):
            return target
    raise ValueError("unmapped checkpoint tensor: " + parameter_name)


def summarize_plan(named_parameters, device_map, source_index, cpu_budget_bytes, scratch_free_bytes, largest_layer_bytes):
    """Fail closed on all parameter/index IDs and every placement/byte budget."""
    if type(cpu_budget_bytes) is not int or cpu_budget_bytes <= 0:
        raise ValueError("CPU budget invalid")
    if type(scratch_free_bytes) is not int or scratch_free_bytes < 0:
        raise ValueError("scratch free invalid")
    if type(largest_layer_bytes) is not int or largest_layer_bytes <= 0:
        raise ValueError("largest-layer size invalid")
    if not isinstance(named_parameters, (list, tuple)):
        raise ValueError("parameter collection invalid")
    if not isinstance(source_index, dict):
        raise ValueError("checkpoint index invalid")
    names = [name for name, _ in named_parameters]
    check_full_key_mapping(names, source_index)
    totals = {"cpu": 0, "disk": 0}
    counts = {"cpu": 0, "disk": 0}
    for name, param in named_parameters:
        if getattr(param, "device", None) is None or param.device.type != "meta":
            raise ValueError("model parameter is not meta")
        dev = resolve_assignment(name, device_map)
        size = int(param.numel()) * 2  # pinned original BF16 effective storage
        totals[dev] += size
        counts[dev] += 1
    if totals["cpu"] + totals["disk"] <= 0:
        raise ValueError("empty plan")
    if totals["cpu"] > cpu_budget_bytes:
        raise ValueError("allocated CPU model-weight bytes exceed proposed budget")
    if totals["disk"] > DISK_BUDGET_BYTES:
        raise ValueError("disk-map parameter size exceeds planning ceiling")
    if totals["disk"] == 0:
        raise ValueError("no disk offload in a purported out-of-core plan")
    # Conservative, deliberately separate copy assumption. This is a
    # PLANNING bound, not guaranteed required bytes or actual peak RSS.
    staged_plus_offload = PINNED_ORIGINAL_SHARD_BYTES + totals["disk"] + SEPARATE_WORKSPACE_RESERVE
    return {
        "cpu_budget_bytes": cpu_budget_bytes,
        "largest_unsplittable_decoder_layer_bytes": largest_layer_bytes,
        "cpu_weight_bytes": totals["cpu"],
        "disk_weight_bytes": totals["disk"],
        "cpu_tensors": counts["cpu"],
        "disk_tensors": counts["disk"],
        "model_parameter_bytes": sum(totals.values()),
        "map_entries": len(device_map),
        "disk_staging_plan_bytes": staged_plus_offload,
        "scratch_free_bytes": scratch_free_bytes,
        "can_stage_full_checkpoint_plus_separate_offload_copy_and_workspace":
            scratch_free_bytes >= staged_plus_offload,
        "placement_level": "BF16_PARAMETER_STATIC_MODEL_PLAN_ONLY",
        "no_split_decoder_layer": True,
        "forward_and_disk_loading_verified": False,
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", type=Path, required=True)
    parser.add_argument("--scratch-root", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    if not args.scratch_root.is_dir():
        raise ValueError("scratch root not a directory")
    manifest = json.loads(args.manifest.read_text(encoding="utf-8"))
    check_manifest(manifest)
    config_data = pinned_metadata("config.json", manifest)
    index = pinned_metadata("model.safetensors.index.json", manifest)
    layout = analyze(config_data, index, manifest, shutil.disk_usage(args.scratch_root).free)

    import torch
    from accelerate import init_empty_weights, infer_auto_device_map
    from accelerate.utils.modeling import compute_module_sizes
    from transformers.models.qwen3_moe.configuration_qwen3_moe import Qwen3MoeConfig
    from transformers.models.qwen3_moe.modeling_qwen3_moe import Qwen3MoeForCausalLM
    from g26_sharded_checkpoint_smoke import test_runtime_identity

    test_runtime_identity()
    t0 = time.perf_counter()
    with init_empty_weights(include_buffers=True):
        model = Qwen3MoeForCausalLM(Qwen3MoeConfig.from_dict(config_data))
    named_parameters = list(model.named_parameters(remove_duplicate=False))
    check_full_key_mapping([name for name, _ in named_parameters], index["weight_map"])
    if any(param.device.type != "meta" for _, param in named_parameters):
        raise ValueError("full original architecture materialized real weights")
    sizes = compute_module_sizes(model, dtype=torch.bfloat16)
    payload_bytes = sum(int(param.numel()) * 2 for _, param in named_parameters)
    if payload_bytes != layout["index_payload_bytes"]:
        raise ValueError("full original BF16 parameter count differs from pinned source index")
    if not payload_bytes <= sizes[""] <= payload_bytes + 64 * 1024**2:
        raise ValueError("full model size disagrees with pinned source weights")
    largest_layer = max(int(sizes[f"model.layers.{i}"]) for i in range(48))
    if largest_layer <= 0:
        raise ValueError("invalid unsplittable layer size")

    plans = []
    for gib in CPU_BUDGET_GIB:
        budget = gib * GIB
        auto_map = infer_auto_device_map(
            model,
            max_memory={"cpu": budget, "disk": DISK_BUDGET_BYTES},
            no_split_module_classes=list(NO_SPLIT),
            dtype=torch.bfloat16,
            offload_buffers=True,
            clean_result=True,
        )
        result = summarize_plan(
            named_parameters, auto_map, index["weight_map"], budget,
            layout["scratch_free_bytes_at_probe"], largest_layer,
        )
        result["cpu_budget_gib"] = gib
        result["device_map"] = dict(auto_map)
        plans.append(result)

    receipt = {
        "schema": "wingless.gcen.g28.meta-bf16-device-plan.v1",
        "result": "PINNED_ORIGINAL_META_CPU_DISK_DEVICE_PLAN_PASS",
        "model": MODEL, "revision": REVISION,
        "source_metadata_sha256": layout["metadata_sha256"],
        "original_shard_bytes": PINNED_ORIGINAL_SHARD_BYTES,
        "model_weight_bytes": payload_bytes,
        "full_model_static_module_bytes": int(sizes[""]),
        "largest_unsplittable_decoder_layer_bytes": largest_layer,
        "no_split": list(NO_SPLIT),
        "cpu_only_scenarios": plans,
        "source_model_weights_loaded": False,
        "original_model_forward_executed": False,
        "actual_full_model_peak_cpu_rss": "UNMEASURED",
        "actual_full_model_peak_vram": "UNMEASURED",
        "actual_full_model_inference_latency": "UNMEASURED",
        "original_model_loader_operational": "UNVERIFIED",
        "offload_plan_feasibility": "STATIC_ESTIMATE_ONLY",
        "execution_authorized_by_this_preflight": False,
        "scientific_acceptance": "NOT_GRANTED",
        "elapsed_seconds": time.perf_counter() - t0,
    }
    if args.output.exists():
        raise FileExistsError("G28 receipt already exists")
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as handle:
        json.dump(receipt, handle, sort_keys=True, indent=2)
        handle.write("\n")
    print("GCEN_G28_META_DEVICE_MAP_PASS", payload_bytes, len(plans))


if __name__ == "__main__":
    main()
