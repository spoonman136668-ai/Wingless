"""G27 full original Qwen3-MoE meta skeleton and index-to-module audit.

No original weight files are downloaded, loaded, inferred over, or materialized.
Uses only pinned SHA-verified config.json and Safetensors index metadata.
CPU-only, zero-weight model construction via init_empty_weights(include_buffers=True).
"""
import argparse
import gc
import json
import pathlib
import time

from out_of_core_preflight import pinned_metadata, analyze, SHARDS
from runtime_feasibility import MODEL, REVISION, check_manifest

EXPECTED_INDEXED_TENSORS = 18867
EXPECTED_EXPERTS = 48 * 128
EXPECTED_EXPERT_PROJECTIONS = EXPECTED_EXPERTS * 3
SELECTED_TENSORS = (
    "model.layers.0.mlp.experts.0.down_proj.weight",
    "model.layers.47.mlp.experts.127.down_proj.weight",
)


def check_full_key_mapping(parameter_names, index_weight_map):
    """Pure validation, independent of Torch and actual weight byte access."""
    if not isinstance(index_weight_map, dict) or len(index_weight_map) != EXPECTED_INDEXED_TENSORS:
        raise ValueError("unexpected original indexed tensor count")
    if not isinstance(parameter_names, (list, tuple)):
        raise ValueError("model parameter names missing")
    if len(parameter_names) != len(set(parameter_names)):
        raise ValueError("duplicate registered parameter names")
    if len(parameter_names) != EXPECTED_INDEXED_TENSORS:
        raise ValueError(f"meta parameter count mismatch: {len(parameter_names)}")
    model_keys, source_keys = set(parameter_names), set(index_weight_map)
    missing, unexpected = source_keys - model_keys, model_keys - source_keys
    if missing or unexpected:
        raise ValueError(
            "checkpoint/module name mismatch: "
            f"missing={sorted(missing)[:12]} unexpected={sorted(unexpected)[:12]}"
        )
    if set(index_weight_map.values()) != set(SHARDS):
        raise ValueError("incomplete 16-shard source mapping")
    for key in SELECTED_TENSORS:
        if key not in model_keys:
            raise ValueError("required expert module tensor missing")
    experts = set()
    count = 0
    for name in parameter_names:
        parts = name.split(".")
        if (len(parts) == 8 and parts[0] == "model" and parts[1] == "layers"
            and parts[3:5] == ["mlp", "experts"] and parts[6] in ("gate_proj", "up_proj", "down_proj")
            and parts[7] == "weight" and parts[2].isdigit() and parts[5].isdigit()):
            layer, expert = int(parts[2]), int(parts[5])
            if not 0 <= layer < 48 or not 0 <= expert < 128:
                raise ValueError("expert topology out of frozen range")
            experts.add((layer, expert))
            count += 1
    if len(experts) != EXPECTED_EXPERTS or count != EXPECTED_EXPERT_PROJECTIONS:
        raise ValueError(
            f"expert projection coverage missing: experts={len(experts)} projections={count}"
        )
    return {
        "indexed_tensors": len(model_keys),
        "expert_modules": len(experts),
        "expert_projection_tensors": count,
        "source_shards": len(SHARDS),
        "missing_tensors": 0,
        "unexpected_tensors": 0,
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", type=pathlib.Path, required=True)
    parser.add_argument("--output", type=pathlib.Path, required=True)
    args = parser.parse_args()

    manifest = json.loads(args.manifest.read_text(encoding="utf-8"))
    check_manifest(manifest)
    config_data = pinned_metadata("config.json", manifest)
    index_data = pinned_metadata("model.safetensors.index.json", manifest)
    # G24 validates the pinned topology, index, shard coverage, and aggregate
    # byte accounting without downloading even one real tensor.
    layout = analyze(config_data, index_data, manifest, 0)
    start = time.perf_counter()

    import torch
    from accelerate import init_empty_weights
    from transformers.models.qwen3_moe.configuration_qwen3_moe import Qwen3MoeConfig
    from transformers.models.qwen3_moe.modeling_qwen3_moe import (
        Qwen3MoeForCausalLM,
        Qwen3MoeMLP,
    )
    from g26_sharded_checkpoint_smoke import test_runtime_identity
    test_runtime_identity()

    with init_empty_weights(include_buffers=True):
        model = Qwen3MoeForCausalLM(Qwen3MoeConfig.from_dict(config_data))
    if len(model.model.layers) != 48:
        raise ValueError("model has wrong decoder layer count")
    for layer in model.model.layers:
        if len(layer.mlp.experts) != 128:
            raise ValueError("original expert module count mismatch")
    if not isinstance(model.model.layers[0].mlp.experts[0], Qwen3MoeMLP):
        raise ValueError("wrong original expert implementation")
    params = list(model.named_parameters(remove_duplicate=False))
    if not params or any(tensor.device.type != "meta" for _, tensor in params):
        raise ValueError("original skeleton unexpectedly allocated real parameter storage")
    model_names = [name for name, _ in params]
    mapping = check_full_key_mapping(model_names, index_data["weight_map"])
    total_parameters = sum(tensor.numel() for _, tensor in params)
    if total_parameters <= 20_000_000_000:
        raise ValueError("skeleton unexpectedly small for unquantized 30B checkpoint")

    # Peak RSS is an OS measurement, NOT total model runtime memory or peak VRAM.
    import resource
    peak_rss_bytes = resource.getrusage(resource.RUSAGE_SELF).ru_maxrss * 1024
    elapsed_seconds = time.perf_counter() - start
    result = {
        "schema": "wingless.gcen.g27.meta-checkpoint-key-audit.v1",
        "result": "PINNED_ORIGINAL_ARCHITECTURE_META_KEY_MAPPING_PASS",
        "checkpoint_model": MODEL,
        "checkpoint_revision": REVISION,
        "model_type": config_data["model_type"],
        "model_parameters": total_parameters,
        "meta_parameter_tensors": len(params),
        "model_parameter_storage_allocated_bytes": 0,
        "source_index_tensor_payload_bytes": layout["index_payload_bytes"],
        "checkpoint_metadata_sha256": layout["metadata_sha256"],
        "source_index_mapping": mapping,
        "meta_construction_seconds": elapsed_seconds,
        "observed_process_peak_rss_bytes": peak_rss_bytes,
        "gpu_vram_measured": False,
        "checkpoint_weights_downloaded": False,
        "original_model_forward_executed": False,
        "activation_source_authenticity": "UNVERIFIED_NO_ORIGINAL_MODEL_EXECUTION",
        "original_model_disk_dispatch": "NOT_TESTED",
        "scientific_acceptance": "NOT_GRANTED",
        "operational_authorization_by_this_probe": False,
    }
    out = args.output
    if out.exists():
        raise FileExistsError("refusing to overwrite G27 receipt")
    out.parent.mkdir(parents=True, exist_ok=True)
    with out.open("x", encoding="utf-8") as stream:
        json.dump(result, stream, sort_keys=True, indent=2)
        stream.write("\n")
    del model, params
    gc.collect()
    print("GCEN_G27_ORIGINAL_META_MAPPING_PASS", total_parameters, mapping["indexed_tensors"])


if __name__ == "__main__":
    main()
