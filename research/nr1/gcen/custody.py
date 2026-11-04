"""Fail-closed SHA256 custody and selected-tensor extraction; research only."""
import hashlib
import json
from pathlib import Path

REQUIRED = [f"model-{i:05d}-of-00016.safetensors" for i in range(1, 17)] + ["model.safetensors.index.json", "config.json"]

def verify_manifest(root, manifest):
    root = Path(root).resolve()
    if not isinstance(manifest, dict) or set(manifest) != set(REQUIRED):
        raise ValueError("manifest must list exactly all 16 shards, index and config")
    for name in REQUIRED:
        expected = manifest[name]
        if not isinstance(expected, str) or len(expected) != 64 or any(c not in "0123456789abcdef" for c in expected):
            raise ValueError("invalid SHA256: " + name)
        path = root / name
        if not path.is_file() or path.is_symlink():
            raise ValueError("missing or linked file: " + name)
        digest = hashlib.sha256()
        with path.open("rb") as stream:
            for chunk in iter(lambda: stream.read(4 * 1024 * 1024), b""):
                digest.update(chunk)
        if digest.hexdigest() != expected:
            raise ValueError("hash mismatch: " + name)
    return True

def extract_tensor(root, manifest, tensor_name):
    verify_manifest(root, manifest)
    index = json.loads((Path(root) / "model.safetensors.index.json").read_text(encoding="utf-8"))
    shard = index.get("weight_map", {}).get(tensor_name)
    if shard not in REQUIRED[:16]:
        raise ValueError("tensor missing from index or invalid shard")
    from safetensors import safe_open
    with safe_open(str(Path(root) / shard), framework="np", device="cpu") as reader:
        return reader.get_tensor(tensor_name)
