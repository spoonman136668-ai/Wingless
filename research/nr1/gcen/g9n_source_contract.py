"""G9N repository-native original-Qwen source/placement admission preflight.

Pure source/index validation. No weight bytes, model construction, downloads,
tensor assignment, GPU allocations, or scientific acceptance.
"""
from collections import Counter
from g9b_tensor_budget import MODEL, REVISION, expected_tensors, reconcile_index
from g28_original_device_plan import resolve_assignment
from runtime_feasibility import check_manifest

SOURCE_TENSORS = 18867
MAX_DISK_BYTES = 256 * 1024**3


def freeze_original_assignment(*, manifest, config, index, device_map,
                               cpu_budget_bytes, disk_budget_bytes):
    """Reconstruct expected tensor identities and G28 placement before loading."""
    check_manifest(manifest)
    if (manifest["model"], manifest["revision"]) != (MODEL, REVISION):
        raise ValueError("original checkpoint identity mismatch")
    if (type(cpu_budget_bytes) is not int or cpu_budget_bytes <= 0
        or type(disk_budget_bytes) is not int
        or not 0 < disk_budget_bytes <= MAX_DISK_BYTES):
        raise ValueError("invalid CPU/disk resource budget")
    if (not isinstance(device_map, dict) or not device_map
        or any(type(k) is not str or type(v) is not str
               or v not in ("cpu", "disk") for k, v in device_map.items())):
        raise ValueError("unsupported device placement")
    tensors = expected_tensors(config)
    if len(tensors) != SOURCE_TENSORS:
        raise ValueError("source tensor inventory changed")
    source_bytes = reconcile_index(tensors, index)
    locations, shapes, totals = {}, {}, Counter()
    for key, (shape, _) in tensors.items():
        placement = resolve_assignment(key, device_map)
        locations[key] = placement
        shapes[key] = shape
        count = 1
        for n in shape:
            count *= n
        totals[placement] += 2 * count
    if totals["cpu"] + totals["disk"] != source_bytes:
        raise ValueError("original source byte accounting mismatch")
    if totals["cpu"] > cpu_budget_bytes or totals["disk"] > disk_budget_bytes:
        raise ValueError("device map exceeds resource budget")
    if not totals["disk"]:
        raise ValueError("no out-of-core disk placement")
    return {
        "schema": "wingless.gcen.g9n.repository-source-placement.v1",
        "model": MODEL, "revision": REVISION,
        "shapes": shapes, "target_devices": locations,
        "indexed_tensors": len(shapes),
        "original_tensor_bytes": source_bytes,
        "cpu_bytes": totals["cpu"], "disk_bytes": totals["disk"],
        "weight_shards_read": 0, "full_model_forward_executed": False,
        "gpu_runtime_verified": False, "scientific_acceptance": "NOT_GRANTED",
    }
