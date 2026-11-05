"""Fail-closed original Qwen BF16 expert custody for the optional G9U benchmark.

No download, model forward, GPU allocation, or authority granted here.
The byte hash and original pinned index are verified in this process, not
inferred from a caller-supplied boolean or an unrelated receipt.
"""
import hashlib
import json
from pathlib import Path

from runtime_feasibility import MODEL, REVISION, check_manifest
from g31_original_shard_stage import SHARD, inspect_header, parse_unique_json
from g9n_one_original_tensor import choose_target
from out_of_core_preflight import pinned_metadata

PINNED_SHARD_SHA256 = "454e77b346a61bfb201d54df60e15158838cf930617ee135113556204f2802b5"
EXPECTED_SHAPE = (2048, 768)
EXPECTED_LAYER = 0
MAX_BYTES = 6 * 1024**3


def load_verified_original_expert(shard_path, manifest_path):
    """Return (weight_fp32, custody_receipt) after verifying original bytes."""
    path = Path(shard_path)
    manifest_path = Path(manifest_path)
    if path.is_symlink() or not path.is_file() or path.name != SHARD:
        raise ValueError("original shard path identity invalid")
    size = path.stat().st_size
    if size <= 8 or size > MAX_BYTES:
        raise ValueError("original shard size outside safe bound")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    expected = manifest["files"][SHARD]
    if expected != PINNED_SHARD_SHA256:
        raise ValueError("manifest shard digest differs from authenticated pin")
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(4 * 1024**2), b""):
            digest.update(chunk)
    if digest.hexdigest() != expected:
        raise ValueError("original shard SHA256 mismatch")
    index = pinned_metadata("model.safetensors.index.json", manifest)
    inspect_header(path, index)
    with path.open("rb") as stream:
        length = int.from_bytes(stream.read(8), "little")
        header = parse_unique_json(stream.read(length))
    name, shape, layer = choose_target(header, index["weight_map"])
    if layer != EXPECTED_LAYER or tuple(shape) != EXPECTED_SHAPE:
        raise ValueError("unexpected original expert tensor selection")
    from safetensors import safe_open
    import torch
    with safe_open(str(path), framework="pt", device="cpu") as reader:
        tensor = reader.get_tensor(name)
        if tuple(tensor.shape) != EXPECTED_SHAPE or tensor.dtype != torch.bfloat16:
            raise ValueError("original expert tensor dtype/shape mismatch")
        weight = tensor.to(dtype=torch.float32).numpy()
    after = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(4 * 1024**2), b""):
            after.update(chunk)
    if after.hexdigest() != digest.hexdigest():
        raise ValueError("original shard changed during authenticated tensor read")
    return weight, {
        "model": MODEL, "revision": REVISION,
        "shard": SHARD, "shard_sha256": digest.hexdigest(),
        "tensor_name": name, "tensor_shape": list(shape),
        "tensor_original_dtype": "BF16", "original_shard_authenticated": True,
        "source_authentication_method": "IN_PROCESS_SHA256_PINNED_INDEX_BF16",
    }
