"""G9N full-original host admission disk headroom, preflight only.

Source tensor inventory must be present locally. No weight contents read,
full-checkpoint transfer, Torch import, model creation, inference or promotion.
The separate CPU RAM safety and actual weight custody remain unverified.
"""
from pathlib import Path
import shutil

from g9n_source_contract import freeze_original_assignment
from g9n_shard_admission import preflight

GIB = 1024 ** 3
WORKSPACE_RESERVE = 8 * GIB


def full_checkpoint_disk_preflight(*, manifest, config, index, device_map,
                                   cpu_budget_bytes, disk_budget_bytes,
                                   source_dir, scratch_dir, disk_usage=None):
    """Reject obviously impossible full-source staging before any sink mutation.

    Once all local source shards exist, scratch must have space for a second
    disk-resident copy of planned offloaded BF16 tensors, a temporary copy of
    the largest source shard, and an independent 8 GiB workspace reserve.
    Not a runtime peak-RSS measurement or a permit to execute original weights.
    """
    plan = freeze_original_assignment(
        manifest=manifest, config=config, index=index,
        device_map=device_map, cpu_budget_bytes=cpu_budget_bytes,
        disk_budget_bytes=disk_budget_bytes)
    source, sizes = preflight(manifest, index, plan["shapes"], source_dir)
    scratch = Path(scratch_dir)
    if scratch.is_symlink() or not scratch.is_dir():
        raise ValueError("scratch root absent or symlinked")
    scratch = scratch.resolve(strict=True)
    if scratch == source.resolve(strict=True):
        raise ValueError("scratch root cannot be the original shard directory")
    usage = (disk_usage or shutil.disk_usage)(scratch)
    if type(usage.free) is not int or usage.free < 0:
        raise ValueError("invalid measured scratch free space")
    needed = plan["disk_bytes"] + max(sizes.values()) + WORKSPACE_RESERVE
    if usage.free < needed:
        raise RuntimeError("G9N original-source scratch capacity insufficient")
    return {
        "schema": "wingless.gcen.g9n.full-checkpoint-disk-preflight.v1",
        "model": plan["model"], "revision": plan["revision"],
        "original_shards_present": len(sizes),
        "indexed_tensors": plan["indexed_tensors"],
        "source_shard_logical_bytes": sum(sizes.values()),
        "largest_staged_source_shard_bytes": max(sizes.values()),
        "planned_disk_tensors_bytes": plan["disk_bytes"],
        "workspace_reserve_bytes": WORKSPACE_RESERVE,
        "scratch_free_bytes_at_probe": usage.free,
        "minimum_scratch_free_bytes": needed,
        "scratch_floor_met": True,
        "source_shard_hashes_verified": False,
        "cpu_available_ram_measured": False,
        "full_original_admission_authorized": False,
        "full_model_forward_executed": False,
        "scientific_acceptance": "NOT_GRANTED",
    }
