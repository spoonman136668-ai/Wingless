"""G9N: fail-closed, locally sourced original checkpoint shard admission.

Does not download weights, construct or dispatch a model, or execute inference.
Tensor sinks and partial-state cleanup are supplied by an authorized caller.
Torch/Safetensors imports occur only after source and resource preflight.
"""
import hashlib
from pathlib import Path
import shutil
import tempfile

from g9b_tensor_budget import MODEL, REVISION
from runtime_feasibility import check_manifest

SHARDS = tuple(f"model-{i:05d}-of-00016.safetensors" for i in range(1, 17))
MAX_SHARD_BYTES = 6 * 1024**3
MIN_RESERVE_BYTES = 1024**3
CHUNK_BYTES = 1024**2


def preflight(manifest, index, shapes, source_dir):
    check_manifest(manifest)
    if (manifest["model"], manifest["revision"]) != (MODEL, REVISION):
        raise ValueError("original source identity mismatch")
    if not isinstance(shapes, dict) or not isinstance(index, dict):
        raise ValueError("source tensor metadata missing")
    mapping = index.get("weight_map")
    if (not isinstance(mapping, dict) or set(mapping) != set(shapes)
        or set(mapping.values()) != set(SHARDS)):
        raise ValueError("original tensor/shard inventory mismatch")
    if any(not isinstance(shape, (tuple, list)) or not shape
           or any(type(dim) is not int or dim <= 0 for dim in shape)
           for shape in shapes.values()):
        raise ValueError("invalid original source tensor shape")
    root = Path(source_dir)
    if root.is_symlink() or not root.is_dir():
        raise ValueError("source root missing or symlinked")
    root = root.resolve(strict=True)
    sizes = {}
    for shard in SHARDS:
        p = root / shard
        if p.is_symlink() or not p.is_file() or p.resolve(strict=True).parent != root:
            raise ValueError("source shard missing, symlinked, or escaped: " + shard)
        size = p.stat().st_size
        if not 0 < size <= MAX_SHARD_BYTES:
            raise ValueError("original shard size outside bound: " + shard)
        sizes[shard] = size
    return root, sizes


def admit_local_shards(*, manifest, index, expected_shapes, source_dir,
                       scratch_dir, assign_tensor, discard_partial,
                       min_reserve_bytes=MIN_RESERVE_BYTES):
    """Stage one shard at a time; discard all partial sink state on failure."""
    if not callable(assign_tensor) or not callable(discard_partial):
        raise ValueError("assignment and mandatory discard callbacks required")
    if type(min_reserve_bytes) is not int or min_reserve_bytes < MIN_RESERVE_BYTES:
        raise ValueError("scratch reserve cannot be weakened")
    root, sizes = preflight(manifest, index, expected_shapes, source_dir)
    scratch = Path(scratch_dir)
    if scratch.is_symlink() or not scratch.is_dir():
        raise ValueError("scratch root missing or symlinked")
    scratch = scratch.resolve(strict=True)
    if scratch == root:
        raise ValueError("source and scratch must differ")
    if shutil.disk_usage(scratch).free < MAX_SHARD_BYTES + min_reserve_bytes:
        raise RuntimeError("insufficient staging disk free space")
    from safetensors import safe_open
    import torch
    mapping = index["weight_map"]
    members = {name: {key for key, shard in mapping.items() if shard == name}
               for name in SHARDS}
    completed = []
    assigned = set()
    try:
        with tempfile.TemporaryDirectory(prefix="g9n-stage-", dir=scratch) as temp:
            stage = Path(temp)
            for shard in SHARDS:
                p = root / shard
                if p.is_symlink() or not p.is_file() or p.stat().st_size != sizes[shard]:
                    raise ValueError("source changed since preflight")
                if shutil.disk_usage(scratch).free < sizes[shard] + min_reserve_bytes:
                    raise RuntimeError("insufficient disk for shard and reserve")
                target = stage / shard
                digest = hashlib.sha256()
                copied = 0
                try:
                    with p.open("rb") as reader, target.open("xb") as writer:
                        while True:
                            block = reader.read(CHUNK_BYTES)
                            if not block:
                                break
                            copied += len(block)
                            if copied > sizes[shard]:
                                raise ValueError("source expanded during staging")
                            digest.update(block)
                            writer.write(block)
                    if copied != sizes[shard] or digest.hexdigest() != manifest["files"][shard]:
                        raise ValueError("original shard digest mismatch")
                    if p.stat().st_size != sizes[shard]:
                        raise ValueError("source changed during staging")
                    with safe_open(str(target), framework="pt", device="cpu") as reader:
                        if set(reader.keys()) != members[shard]:
                            raise ValueError("shard tensor keys mismatch")
                        for key in sorted(members[shard]):
                            view = reader.get_slice(key)
                            if tuple(view.get_shape()) != tuple(expected_shapes[key]):
                                raise ValueError("source tensor shape mismatch")
                            tensor = reader.get_tensor(key)
                            if tensor.dtype != torch.bfloat16:
                                raise ValueError("source tensor dtype mismatch")
                            if key in assigned:
                                raise ValueError("duplicate tensor admission")
                            assign_tensor(key, tensor)
                            assigned.add(key)
                    completed.append(shard)
                finally:
                    target.unlink(missing_ok=True)
                    if list(stage.iterdir()):
                        raise RuntimeError("staging cleanup invariant violated")
        if len(completed) != 16 or assigned != set(expected_shapes):
            raise RuntimeError("incomplete original shard admission")
    except BaseException:
        discard_partial()
        raise
    return {"schema": "wingless.gcen.g9n.shard-admission.v1",
            "model": MODEL, "revision": REVISION,
            "verified_shards": len(completed),
            "verified_tensor_count": len(assigned),
            "full_model_forward_executed": False,
            "model_constructed": False,
            "gpu_residency_verified": False,
            "scientific_acceptance": "NOT_GRANTED"}
