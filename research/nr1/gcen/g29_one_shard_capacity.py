"""G29 original-model single-shard staging feasibility, META / METADATA ONLY.

This does NOT implement streamed sharded loading. In pinned Accelerate v1.10.1,
load_checkpoint_and_dispatch expects the source checkpoint/index; it is NOT
admitted as a sequential verified-shard streaming engine by this estimate.

Given a qualified G28 meta device map and the original 16-shard Safetensors
index, calculate an optimistic per-shard sequential-staging capacity proposal
without fetching any of the 61GB original checkpoint weights.
"""
import argparse
import json
from pathlib import Path
import shutil
import time

from g28_original_device_plan import (
    GIB, CPU_BUDGET_GIB, DISK_BUDGET_BYTES, NO_SPLIT,
    resolve_assignment, summarize_plan, SEPARATE_WORKSPACE_RESERVE,
)
from g27_original_meta_mapping import check_full_key_mapping
from out_of_core_preflight import analyze, pinned_metadata, SHARDS
from runtime_feasibility import MODEL, REVISION, PINNED_ORIGINAL_SHARD_BYTES, check_manifest


def shard_payload_footprints(named_parameters, index_weight_map, total_file_bytes):
    """Compute per-shard BF16 tensor payload from genuine meta parameter shapes.

    Aggregate storage overhead is assigned entirely to the largest shard as
    a conservative one-shard physical byte upper bound. No per-shard actual
    file size/custody is asserted.
    """
    if not isinstance(named_parameters, (list, tuple)):
        raise ValueError("parameter collection absent")
    if not isinstance(index_weight_map, dict):
        raise ValueError("source index missing")
    if type(total_file_bytes) is not int or total_file_bytes <= 0:
        raise ValueError("source file storage total invalid")
    check_full_key_mapping([name for name, _ in named_parameters], index_weight_map)
    payload = {name: 0 for name in SHARDS}
    for name, param in named_parameters:
        if getattr(param, "device", None) is None or param.device.type != "meta":
            raise ValueError("all checkpoint model parameters must be meta")
        numel = param.numel()
        if type(numel) is not int or numel <= 0:
            raise ValueError("invalid meta parameter size")
        payload[index_weight_map[name]] += numel * 2
    if any(value <= 0 for value in payload.values()):
        raise ValueError("source index has empty shard")
    total_tensor_bytes = sum(payload.values())
    overhead_bytes = total_file_bytes - total_tensor_bytes
    if overhead_bytes < 16 * 8 or overhead_bytes >= 32 * 1024**2:
        raise ValueError("original file/header byte reconciliation failed")
    return {
        "per_shard_bf16_payload_bytes": payload,
        "total_tensor_payload_bytes": total_tensor_bytes,
        "aggregate_file_overhead_bytes": overhead_bytes,
        "largest_shard_tensor_payload_bytes": max(payload.values()),
        "single_shard_physical_upper_bound_bytes": max(payload.values()) + overhead_bytes,
        "actual_per_shard_file_bytes": "UNVERIFIED_NOT_IN_INDEX",
    }


def sequential_staging_proposal(disk_weight_bytes, shard_upper_bytes, scratch_free_bytes,
                                workspace_reserve_bytes=SEPARATE_WORKSPACE_RESERVE):
    """One-original-shard-at-a-time hypothetical local scratch requirement.

    Assumes no retained full checkpoint or duplicate disk offload outputs; that
    implementation has NOT been demonstrated or authorized.
    """
    for val in (disk_weight_bytes, shard_upper_bytes, scratch_free_bytes, workspace_reserve_bytes):
        if type(val) is not int or val < 0:
            raise ValueError("invalid staging byte input")
    if disk_weight_bytes <= 0 or shard_upper_bytes <= 0 or workspace_reserve_bytes < GIB:
        raise ValueError("invalid staging or reserve requirement")
    required = disk_weight_bytes + shard_upper_bytes + workspace_reserve_bytes
    return {
        "model_disk_weight_payload_bytes": disk_weight_bytes,
        "worst_one_shard_physical_upper_bytes": shard_upper_bytes,
        "extra_workspace_reserve_bytes": workspace_reserve_bytes,
        "optimistic_sequential_scratch_plan_bytes": required,
        "scratch_free_bytes": scratch_free_bytes,
        "fits_this_planning_estimate": scratch_free_bytes >= required,
        "original_checkpoint_full_copy_retained": False,
        "loader_implements_this_streaming_plan": False,
        "all_16_original_shard_hashes_verified_on_runner": False,
        "actual_disk_allocation_and_cpu_peak": "UNMEASURED",
        "original_model_execution_authorized": False,
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", type=Path, required=True)
    parser.add_argument("--scratch-root", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    if not args.scratch_root.is_dir():
        raise ValueError("scratch root is not a directory")
    if args.output.exists():
        raise FileExistsError("G29 receipt already exists")
    manifest = json.loads(args.manifest.read_text(encoding="utf-8"))
    check_manifest(manifest)
    config_data = pinned_metadata("config.json", manifest)
    index_data = pinned_metadata("model.safetensors.index.json", manifest)
    free = shutil.disk_usage(args.scratch_root).free
    original_layout = analyze(config_data, index_data, manifest, free)

    import torch
    from accelerate import init_empty_weights, infer_auto_device_map
    from accelerate.utils.modeling import compute_module_sizes
    from transformers.models.qwen3_moe.configuration_qwen3_moe import Qwen3MoeConfig
    from transformers.models.qwen3_moe.modeling_qwen3_moe import Qwen3MoeForCausalLM
    from g26_sharded_checkpoint_smoke import test_runtime_identity
    test_runtime_identity()
    started = time.perf_counter()
    with init_empty_weights(include_buffers=True):
        model = Qwen3MoeForCausalLM(Qwen3MoeConfig.from_dict(config_data))
    params = list(model.named_parameters(remove_duplicate=False))
    check_full_key_mapping([name for name, _ in params], index_data["weight_map"])
    footprints = shard_payload_footprints(
        params, index_data["weight_map"], PINNED_ORIGINAL_SHARD_BYTES
    )
    if footprints["total_tensor_payload_bytes"] != original_layout["index_payload_bytes"]:
        raise ValueError("checkpoint index tensor bytes disagree with meta architecture")
    sizes = compute_module_sizes(model, dtype=torch.bfloat16)
    largest_layer = max(int(sizes[f"model.layers.{i}"]) for i in range(48))
    if largest_layer <= 0:
        raise ValueError("missing unsplittable layer")

    results = []
    for gib in CPU_BUDGET_GIB:
        memory = gib * GIB
        plan_map = infer_auto_device_map(
            model, max_memory={"cpu": memory, "disk": DISK_BUDGET_BYTES},
            no_split_module_classes=list(NO_SPLIT), dtype=torch.bfloat16,
            offload_buffers=True, clean_result=True,
        )
        map_result = summarize_plan(
            params, plan_map, index_data["weight_map"], memory, free, largest_layer,
        )
        seq = sequential_staging_proposal(
            map_result["disk_weight_bytes"],
            footprints["single_shard_physical_upper_bound_bytes"], free,
        )
        results.append({
            "cpu_budget_gib": gib,
            "cpu_weight_bytes": map_result["cpu_weight_bytes"],
            "disk_weight_bytes": map_result["disk_weight_bytes"],
            "conservative_two_copy_planning_bytes": map_result["disk_staging_plan_bytes"],
            "one_shard_streaming_hypothesis": seq,
        })

    receipt = {
        "schema": "wingless.gcen.g29.single-shard-capacity-proposal.v1",
        "result": "PINNED_SINGLE_SHARD_STORAGE_HYPOTHESIS_CALCULATED",
        "model": MODEL, "revision": REVISION,
        "source_metadata_sha256": original_layout["metadata_sha256"],
        "original_checkpoint_full_file_bytes": PINNED_ORIGINAL_SHARD_BYTES,
        "shard_footprints": footprints,
        "scenarios": results,
        "observed_scratch_free_bytes": free,
        "existing_accelerate_load_checkpoint_and_dispatch_sequential_streaming_supported": "NOT_VERIFIED",
        "full_original_checkpoint_bytes_accessed": False,
        "shard_staging_custody_tested": False,
        "original_model_forward_executed": False,
        "scientific_acceptance": "NOT_GRANTED",
        "execution_authorized_by_this_preflight": False,
        "elapsed_seconds": time.perf_counter() - started,
    }
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as stream:
        json.dump(receipt, stream, indent=2, sort_keys=True)
        stream.write("\n")
    print("GCEN_G29_ONE_SHARD_CAPACITY_HYPOTHESIS_COMPLETE",
          footprints["single_shard_physical_upper_bound_bytes"],
          len(results))


if __name__ == "__main__":
    main()
