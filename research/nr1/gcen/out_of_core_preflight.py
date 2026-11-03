"""G24 pinned-index and storage preflight: no weights, inference or runner activation.

Only two upstream metadata files may be fetched (config.json and the Safetensors
index). Their raw bytes MUST match the already-frozen 18-file SHA-256 inventory
before parsing. Host scratch space is diagnostic, not runtime authorization.
"""
import argparse
import hashlib
import json
import pathlib
import shutil
import urllib.request

from runtime_feasibility import (
    MODEL, REVISION, PINNED_ORIGINAL_SHARD_BYTES, check_manifest,
)

BASE = f"https://huggingface.co/{MODEL}/resolve/{REVISION}"
MAX_BYTES = {"config.json": 2 * 1024 * 1024,
             "model.safetensors.index.json": 24 * 1024 * 1024}
SHARDS = tuple(f"model-{i:05d}-of-00016.safetensors" for i in range(1, 17))
SINGLE_SHARD_SCRATCH_RESERVE = 6 * 1024**3
# Conservative planning reserve only, NOT an Accelerate API requirement.
ADDITIONAL_OFFLOAD_WORKSPACE_RESERVE = PINNED_ORIGINAL_SHARD_BYTES + 8 * 1024**3


def pinned_metadata(name, manifest, fetch=None):
    check_manifest(manifest)
    if name not in MAX_BYTES:
        raise ValueError("metadata filename not allowlisted")
    if fetch is None:
        def fetch(url, limit):
            with urllib.request.urlopen(url, timeout=40) as response:
                data = response.read(limit + 1)
                if len(data) > limit:
                    raise ValueError("metadata response exceeded byte limit")
                return data
    raw = fetch(f"{BASE}/{name}", MAX_BYTES[name])
    if not isinstance(raw, bytes) or not 0 < len(raw) <= MAX_BYTES[name]:
        raise ValueError("invalid or oversized metadata bytes")
    if hashlib.sha256(raw).hexdigest() != manifest["files"][name]:
        raise ValueError("pinned metadata checksum mismatch")
    decoded = json.loads(raw)
    if not isinstance(decoded, dict):
        raise ValueError("metadata is not a JSON object")
    return decoded


def analyze(config, index, manifest, scratch_free_bytes):
    check_manifest(manifest)
    if type(scratch_free_bytes) is not int or scratch_free_bytes < 0:
        raise ValueError("scratch free space is invalid")
    if not isinstance(config, dict) or config.get("model_type") != "qwen3_moe":
        raise ValueError("pinned Qwen3-MoE architecture missing")
    if config.get("num_hidden_layers") != 48 or config.get("num_experts") != 128 or config.get("num_experts_per_tok") != 8:
        raise ValueError("pinned 48-layer/128-expert/8-active topology mismatch")
    if not isinstance(index, dict):
        raise ValueError("invalid checkpoint index")
    weight_map = index.get("weight_map")
    if not isinstance(weight_map, dict) or not weight_map:
        raise ValueError("empty or invalid weight map")
    if any(not isinstance(k, str) or not k or v not in SHARDS for k, v in weight_map.items()):
        raise ValueError("unknown checkpoint tensor shard")
    if set(weight_map.values()) != set(SHARDS):
        raise ValueError("index does not reference exactly 16 original shards")
    for key in (
        "model.layers.0.mlp.experts.0.down_proj.weight",
        "model.layers.47.mlp.experts.127.down_proj.weight",
    ):
        if key not in weight_map:
            raise ValueError("original expert tensor mapping missing: " + key)
    metadata = index.get("metadata")
    total_size = metadata.get("total_size") if isinstance(metadata, dict) else None
    if type(total_size) is not int or not 0 < total_size <= PINNED_ORIGINAL_SHARD_BYTES:
        raise ValueError("invalid original index tensor-byte total")
    if total_size < PINNED_ORIGINAL_SHARD_BYTES * 0.9:
        raise ValueError("indexed tensors unexpectedly small versus original shard bytes")
    return {
        "schema": "wingless.gcen.g24.index-storage-preflight.v1",
        "model": MODEL,
        "revision": REVISION,
        "checkpoint_format": "original_safetensors",
        "metadata_sha256": {
            n: manifest["files"][n] for n in MAX_BYTES
        },
        "expected_shards": len(SHARDS),
        "indexed_tensor_count": len(weight_map),
        "index_payload_bytes": total_size,
        "original_shard_bytes": PINNED_ORIGINAL_SHARD_BYTES,
        "scratch_free_bytes_at_probe": scratch_free_bytes,
        "scratch_single_shard_planning_bytes": SINGLE_SHARD_SCRATCH_RESERVE,
        "scratch_original_checkpoint_floor_bytes": PINNED_ORIGINAL_SHARD_BYTES,
        "scratch_plus_separate_offload_planning_bytes":
            PINNED_ORIGINAL_SHARD_BYTES + ADDITIONAL_OFFLOAD_WORKSPACE_RESERVE,
        "single_shard_staging_capacity_present": scratch_free_bytes >= SINGLE_SHARD_SCRATCH_RESERVE,
        "full_checkpoint_staging_capacity_present": scratch_free_bytes >= PINNED_ORIGINAL_SHARD_BYTES,
        "full_checkpoint_plus_separate_offload_planning_capacity_present":
            scratch_free_bytes >= PINNED_ORIGINAL_SHARD_BYTES + ADDITIONAL_OFFLOAD_WORKSPACE_RESERVE,
        "result": "PINNED_METADATA_AND_DISK_PREFLIGHT_ONLY",
        "checkpoint_bytes_on_probe_host": "NOT_VERIFIED",
        "compatible_offload_engine": "NOT_VERIFIED",
        "expert_forward_hook": "NOT_VERIFIED",
        "source_runner_authorization": "NOT_VERIFIED",
        "full_model_forward": "NOT_EXECUTED",
        "execution_authorized_by_this_probe": False,
    }


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--manifest", required=True, type=pathlib.Path)
    ap.add_argument("--scratch-root", required=True, type=pathlib.Path)
    ap.add_argument("--output", required=True, type=pathlib.Path)
    args = ap.parse_args()
    manifest = json.loads(args.manifest.read_text(encoding="utf-8"))
    if not args.scratch_root.is_dir():
        raise ValueError("scratch-root must be an existing directory")
    config = pinned_metadata("config.json", manifest)
    index = pinned_metadata("model.safetensors.index.json", manifest)
    data = analyze(config, index, manifest, shutil.disk_usage(args.scratch_root).free)
    if args.output.exists():
        raise FileExistsError("refusing to overwrite existing G24 receipt")
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print("GCEN_G24_METADATA_PREFLIGHT_COMPLETE", data["indexed_tensor_count"],
          data["scratch_free_bytes_at_probe"])


if __name__ == "__main__":
    main()
