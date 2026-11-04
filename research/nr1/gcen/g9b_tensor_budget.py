"""G9B: exact-index-reconciled, weight-free Qwen3-MoE VRAM budget.

Development-only analytical screen. Fetches ONLY SHA256-pinned config/index
metadata through existing GCEN custody tooling. Never fetches weight shards,
constructs a model, runs inference, or measures GPU VRAM. Alternative 2/3-bit
storage formats below are hypothetical budgets, NOT implemented operators.
"""
import argparse
from collections import defaultdict
import json
import math
from pathlib import Path
import re

MODEL = "Qwen/Qwen3-30B-A3B"
REVISION = "ad44e777bcd18fa416d9da3bd8f70d33ebb85d39"

EXPECTED_TENSORS = 18867  # Previously checked against original meta model by G27.
VRAM_CEILING = 12 * (1 << 30)
CONTEXT = 8192
WORKSPACE_RESERVE = 1 << 30  # Planning-only, NOT measured runtime peak.
SOURCE_SHA = re.compile(r"[a-f0-9]{40}\Z")

# Shape-changing config changes must be requalified, not silently absorbed.
FROZEN_CONFIG = {
    "model_type": "qwen3_moe", "hidden_size": 2048,
    "vocab_size": 151936, "head_dim": 128,
    "num_hidden_layers": 48, "num_attention_heads": 32,
    "num_key_value_heads": 4, "num_experts": 128,
    "num_experts_per_tok": 8, "moe_intermediate_size": 768,
    "torch_dtype": "bfloat16", "tie_word_embeddings": False,
    "attention_bias": False, "mlp_only_layers": [],
}


def expected_tensors(config):
    if not isinstance(config, dict) or any(
        config.get(key) != expected or type(config.get(key)) is not type(expected)
        for key, expected in FROZEN_CONFIG.items()
    ):
        raise ValueError("pinned original Qwen architecture/config mismatch")
    h = config["hidden_size"]
    head = config["head_dim"]
    qdim = config["num_attention_heads"] * head
    kvdim = config["num_key_value_heads"] * head
    intermediate = config["moe_intermediate_size"]
    out = {}

    def put(name, shape, kind):
        if name in out or any(type(d) is not int or d <= 0 for d in shape):
            raise ValueError("invalid or duplicate expected tensor")
        out[name] = (tuple(shape), kind)

    put("model.embed_tokens.weight", (config["vocab_size"], h), "shared")
    put("model.norm.weight", (h,), "shared")
    put("lm_head.weight", (config["vocab_size"], h), "shared")
    for layer in range(config["num_hidden_layers"]):
        p = f"model.layers.{layer}."
        put(p + "input_layernorm.weight", (h,), "shared")
        put(p + "post_attention_layernorm.weight", (h,), "shared")
        att = p + "self_attn."
        put(att + "q_proj.weight", (qdim, h), "attention")
        put(att + "k_proj.weight", (kvdim, h), "attention")
        put(att + "v_proj.weight", (kvdim, h), "attention")
        put(att + "o_proj.weight", (h, qdim), "attention")
        put(att + "q_norm.weight", (head,), "attention")
        put(att + "k_norm.weight", (head,), "attention")
        put(p + "mlp.gate.weight", (config["num_experts"], h), "router")
        for expert in range(config["num_experts"]):
            e = p + f"mlp.experts.{expert}."
            put(e + "gate_proj.weight", (intermediate, h), "experts")
            put(e + "up_proj.weight", (intermediate, h), "experts")
            put(e + "down_proj.weight", (h, intermediate), "experts")
    if len(out) != EXPECTED_TENSORS:
        raise ValueError("expected Qwen source skeleton count drift")
    return out


def reconcile_index(shapes, index):
    if not isinstance(index, dict) or not isinstance(index.get("weight_map"), dict):
        raise ValueError("source index missing weight map")
    mapping = index["weight_map"]
    if set(shapes) != set(mapping):
        raise ValueError("source index tensor keys differ from frozen G27 skeleton")
    allowed_shards = {f"model-{i:05d}-of-00016.safetensors" for i in range(1, 17)}
    if set(mapping.values()) != allowed_shards:
        raise ValueError("source index shard coverage mismatch")
    total = sum(math.prod(shape) * 2 for shape, _ in shapes.values())
    metadata = index.get("metadata")
    if not isinstance(metadata, dict) or type(metadata.get("total_size")) is not int or metadata["total_size"] != total:
        raise ValueError("pinned index BF16 total bytes do not reconcile")
    return total


def split_tensors(shapes):
    totals = defaultdict(lambda: {"tensor_count": 0, "parameters": 0, "bf16_bytes": 0})
    for shape, kind in shapes.values():
        count = math.prod(shape)
        totals[kind]["tensor_count"] += 1
        totals[kind]["parameters"] += count
        totals[kind]["bf16_bytes"] += count * 2
    return {k: totals[k] for k in ("experts", "router", "attention", "shared")}


def quantized_expert_bytes(shapes, bits, block, scale_bytes):
    if bits not in (2, 3, 4) or scale_bytes not in (2, 4) or type(block) is not int or block <= 0:
        raise ValueError("unsupported hypothetical storage configuration")
    codes = scales = headers = 0
    for shape, kind in shapes.values():
        if kind != "experts":
            continue
        if len(shape) != 2 or shape[1] % block:
            raise ValueError("expert shape incompatible with compression block")
        n = math.prod(shape)
        codes += (n * bits + 7) // 8
        scales += shape[0] * (shape[1] // block) * scale_bytes
        headers += 24  # Same 3x uint64 accounting used by G8, per matrix.
    return {"packed_codes_bytes": codes, "scales_bytes": scales,
            "per_tensor_shape_metadata_bytes": headers,
            "total_expert_bytes": codes + scales + headers}


def make_report(config, index, *, source_sha):
    if type(source_sha) is not str or not SOURCE_SHA.fullmatch(source_sha):
        raise ValueError("exact source commit SHA required")
    tensors = expected_tensors(config)
    source_bytes = reconcile_index(tensors, index)
    kinds = split_tensors(tensors)
    nonexpert_bf16_bytes = sum(v["bf16_bytes"] for k, v in kinds.items() if k != "experts")
    token_kv_bytes = (2 * config["num_hidden_layers"] * config["num_key_value_heads"] *
                      config["head_dim"] * 2)
    kv_bytes = CONTEXT * token_kv_bytes
    options = (("int4_block64_fp32", 4, 64, 4),
               ("int3_block64_fp32", 3, 64, 4),
               ("int2_block64_fp32", 2, 64, 4),
               ("int2_block128_fp16", 2, 128, 2),
               ("int2_block256_fp16", 2, 256, 2))
    cases = []
    for name, bits, block, scale_bytes in options:
        expert = quantized_expert_bytes(tensors, bits, block, scale_bytes)
        weights = nonexpert_bf16_bytes + expert["total_expert_bytes"]
        planned = weights + kv_bytes + WORKSPACE_RESERVE
        cases.append({"name": name, "expert_bits": bits, "block_columns": block,
                      "scale_bytes": scale_bytes, **expert,
                      "nonexpert_bf16_bytes": nonexpert_bf16_bytes,
                      "total_weight_storage_bytes": weights,
                      "bf16_kv_bytes": kv_bytes,
                      "hypothetical_workspace_reserve_bytes": WORKSPACE_RESERVE,
                      "planned_total_bytes": planned,
                      "remaining_bytes": VRAM_CEILING - planned,
                      "analytical_fit_under_assumptions": planned <= VRAM_CEILING,
                      "real_gpu_residency_verified": False,
                      "operator_available": name == "int4_block64_fp32",
                      "operator_scope": "G8_WEIGHT_WINDOWS_CPU_ONLY" if name == "int4_block64_fp32" else "NOT_IMPLEMENTED"})
    return {
        "schema": "wingless.cnc.g9b.pinned-index-analytical-budget.v1",
        "result": "EXACT_INDEX_RECONCILED_STORAGE_MODEL_ONLY",
        "source_sha": source_sha, "checkpoint_model": MODEL,
        "checkpoint_revision": REVISION, "indexed_bf16_tensor_bytes": source_bytes,
        "indexed_tensor_count": len(tensors),
        "parameter_total": sum(v["parameters"] for v in kinds.values()),
        "parameter_classes": kinds,
        "vram_target_bytes": VRAM_CEILING,
        "context_tokens": CONTEXT, "batch_size": 1,
        "kv_cache_dtype": "BF16", "kv_cache_bytes_per_token": token_kv_bytes,
        "assumptions": ["all nonexpert parameters BF16 resident",
                        "all expert projections compressed in place and fully resident",
                        "uniform blockwise per-row scales; no padding/fragmentation",
                        "12 GiB nominal available GPU storage; real usable capacity unknown",
                        "1 GiB arbitrary workspace reserve, not a measured runtime peak"],
        "scenarios": cases,
        "full_model_forward": False, "actual_gpu_peak_vram_measured": False,
        "real_model_activation_custody": "NOT_CAPTURED",
        "protected_heldout_used": False, "scientific_acceptance": "NOT_GRANTED",
        "deployment_or_activation": False,
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", required=True, type=Path)
    parser.add_argument("--source-sha", required=True)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    if args.output.exists():
        raise FileExistsError("refuse to overwrite prior immutable budget receipt")
    from out_of_core_preflight import pinned_metadata
    from runtime_feasibility import MODEL as PINNED_MODEL, REVISION as PINNED_REVISION, check_manifest
    if (PINNED_MODEL, PINNED_REVISION) != (MODEL, REVISION):
        raise ValueError("original GCEN custody constants changed")
    manifest = json.loads(args.manifest.read_text(encoding="utf-8"))
    check_manifest(manifest)
    config = pinned_metadata("config.json", manifest)
    index = pinned_metadata("model.safetensors.index.json", manifest)
    result = make_report(config, index, source_sha=args.source_sha)
    result["source_metadata_sha256"] = {
        "config.json": manifest["files"]["config.json"],
        "model.safetensors.index.json": manifest["files"]["model.safetensors.index.json"],
    }
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as stream:
        json.dump(result, stream, sort_keys=True, indent=2)
        stream.write("\n")
    print("G9B_PINNED_INDEX_STORAGE_BUDGET_COMPLETE", result["parameter_total"])


if __name__ == "__main__":
    main()
