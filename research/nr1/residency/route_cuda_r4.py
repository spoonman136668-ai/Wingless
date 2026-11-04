"""R4 optional CUDA CPU->GPU transfer calibration; no model weights or inference.

The qualified native workflow runs only offline contract tests. Hardware execution
requires a separately authorized GPU runner and an explicit --require-cuda call.
Measured CUDA events and synchronized host wall time are not interchangeable.
"""
import argparse
import hashlib
import json
import math
import os
from pathlib import Path
import statistics
import subprocess
import time

SCHEMA = "wingless.nr1.residency.r4.cuda-transfer-development.v1"
SIZES_MIB = (1, 4, 16)
TRIALS = 12
WARMUPS = 3
MODES = ("pageable", "pinned")
SEED = 20261103


def plan():
    """Frozen 72-trial matrix; no hardware discovered or resources allocated."""
    return tuple((size, mode, rep)
                 for rep in range(TRIALS)
                 for size in SIZES_MIB
                 for mode in (MODES if rep % 2 == 0 else tuple(reversed(MODES))))


def validate_plan(rows):
    if len(rows) != TRIALS * len(SIZES_MIB) * len(MODES):
        raise ValueError("wrong trial count")
    for size in SIZES_MIB:
        for mode in MODES:
            reps = sorted(rep for n, m, rep in rows if n == size and m == mode)
            if reps != list(range(TRIALS)):
                raise ValueError("missing/repeated trial")
    expected = plan()
    if rows != expected:
        raise ValueError("nonfrozen or noncounterbalanced order")
    return True


def validate_receipt(receipt):
    if receipt.get("schema") != SCHEMA or receipt.get("status") != "R4_SYNTHETIC_CUDA_COPY_MEASURED":
        raise ValueError("receipt schema/status mismatch")
    if receipt.get("original_weights_used") is not False or receipt.get("model_inference_executed") is not False:
        raise ValueError("invalid scope")
    if receipt.get("source_custody", {}).get("source_identity_verified") is not True or receipt.get("source_custody", {}).get("checkout_clean") is not True or receipt.get("source_custody", {}).get("source_sha") != receipt.get("source_sha"):
        raise ValueError("source provenance unverified")
    if receipt.get("scientific_acceptance") != "NOT_GRANTED":
        raise ValueError("scientific acceptance forbidden")
    expected = plan()
    items = receipt.get("trials")
    if not isinstance(items, list) or len(items) != len(expected):
        raise ValueError("missing trials")
    for item, (size, mode, rep) in zip(items, expected):
        if (item.get("size_mib"), item.get("source_mode"), item.get("repetition")) != (size, mode, rep):
            raise ValueError("trial identity mismatch")
        if item.get("correct_copy") is not True:
            raise ValueError("source/device parity failed")
        for key in ("host_synchronized_wall_ms", "cuda_event_elapsed_ms"):
            val = item.get(key)
            if type(val) not in (float, int) or not math.isfinite(val) or val <= 0:
                raise ValueError("invalid timing")
        if item.get("device_bytes") != size * 1024**2:
            raise ValueError("allocation mismatch")
    if not receipt.get("cuda_available") or not receipt.get("cuda_device_name"):
        raise ValueError("no verified CUDA device")
    return True


def verify_source_identity(source_sha, repository_root=None):
    """Reject checkout drift and dirty/untracked source before GPU allocation."""
    if not isinstance(source_sha, str) or len(source_sha) != 40 or any(
            c not in "0123456789abcdef" for c in source_sha):
        raise ValueError("invalid exact source SHA")
    root = Path(repository_root) if repository_root is not None else Path(__file__).resolve().parents[3]
    if not root.is_dir():
        raise ValueError("repository root absent")
    def git(*args):
        try:
            run = subprocess.run(
                ["git", *args], cwd=str(root), text=True,
                capture_output=True, timeout=15, check=True)
        except (OSError, subprocess.SubprocessError) as exc:
            raise RuntimeError("R4_SOURCE_CUSTODY_UNVERIFIED") from exc
        return run.stdout.strip()
    if git("rev-parse", "HEAD") != source_sha:
        raise RuntimeError("R4_EXACT_SOURCE_SHA_MISMATCH")
    if git("status", "--porcelain", "--untracked-files=all"):
        raise RuntimeError("R4_DIRTY_OR_UNTRACKED_WORKSPACE")
    return {"source_sha": source_sha, "checkout_clean": True,
            "source_identity_verified": True}


def require_cuda():
    try:
        import torch
    except ImportError as exc:
        raise RuntimeError("R4_GPU_UNAVAILABLE: torch is not installed") from exc
    if not torch.cuda.is_available():
        raise RuntimeError("R4_GPU_UNAVAILABLE: no authorized CUDA device visible")
    return torch


def run(source_sha):
    if len(source_sha) != 40 or any(c not in "0123456789abcdef" for c in source_sha):
        raise ValueError("expected immutable 40-character source SHA")
    validate_plan(plan())
    provenance = verify_source_identity(source_sha)
    torch = require_cuda()
    torch.manual_seed(SEED)
    torch.cuda.set_device(0)
    device = torch.device("cuda:0")
    name = torch.cuda.get_device_name(device)
    free, total = torch.cuda.mem_get_info(device)
    if free < 512 * 1024**2:
        raise RuntimeError("R4_GPU_CAPACITY_INSUFFICIENT: 512 MiB free VRAM required")
    torch.cuda.reset_peak_memory_stats(device)
    samples = {}
    trial_records = []
    try:
        for size in SIZES_MIB:
            count = size * 1024**2 // 4
            # Small bounded deterministic synthetic data, not checkpoint tensors.
            pageable = torch.arange(count, dtype=torch.int32, device="cpu")
            pinned = pageable.pin_memory()
            target = torch.empty_like(pageable, device=device)
            samples[size] = (pageable, pinned, target)
        torch.cuda.synchronize(device)
        for size, mode, repetition in plan():
            pageable, pinned, target = samples[size]
            origin = pinned if mode == "pinned" else pageable
            for _ in range(WARMUPS):
                target.copy_(origin, non_blocking=(mode == "pinned"))
            torch.cuda.synchronize(device)
            start_event = torch.cuda.Event(enable_timing=True)
            end_event = torch.cuda.Event(enable_timing=True)
            torch.cuda.synchronize(device)
            start_wall = time.perf_counter_ns()
            start_event.record()
            target.copy_(origin, non_blocking=(mode == "pinned"))
            end_event.record()
            torch.cuda.synchronize(device)
            wall_ms = (time.perf_counter_ns() - start_wall) / 1e6
            event_ms = start_event.elapsed_time(end_event)
            # Verify outside the timed interval; parity must be complete.
            parity = bool(torch.equal(target.cpu(), pageable))
            trial_records.append({
                "size_mib": size, "source_mode": mode, "repetition": repetition,
                "device_bytes": size * 1024**2,
                "correct_copy": parity,
                "host_synchronized_wall_ms": wall_ms,
                "cuda_event_elapsed_ms": event_ms,
            })
            if not parity:
                raise RuntimeError("R4_CUDA_TRANSFER_PARITY_FAIL")
        summaries = [{
            "size_mib": size, "mode": mode,
            "median_cuda_event_ms": statistics.median(
                x["cuda_event_elapsed_ms"] for x in trial_records
                if x["size_mib"] == size and x["source_mode"] == mode),
            "median_host_wall_ms": statistics.median(
                x["host_synchronized_wall_ms"] for x in trial_records
                if x["size_mib"] == size and x["source_mode"] == mode),
        } for size in SIZES_MIB for mode in MODES]
        receipt = {
            "schema": SCHEMA, "status": "R4_SYNTHETIC_CUDA_COPY_MEASURED",
            "source_sha": source_sha,
            "source_custody": provenance,
            "torch_version": str(torch.__version__),
            "cuda_runtime_version": str(torch.version.cuda),
            "cuda_available": True,
            "cuda_device_name": name,
            "gpu_total_memory_bytes": total,
            "gpu_free_memory_before_bytes": free,
            "gpu_peak_allocated_bytes": torch.cuda.max_memory_allocated(device),
            "gpu_peak_reserved_bytes": torch.cuda.max_memory_reserved(device),
            "synthetic_source": "DETERMINISTIC_INT32_CPU_BUFFERS_ONLY",
            "pinned_and_pageable_modes": list(MODES),
            "synchronization": "DEVICE_SYNC_BEFORE_AND_AFTER_EACH_COPY",
            "gpu_cache_policy_executed": False,
            "real_router_traces_used": False,
            "original_weights_used": False,
            "model_inference_executed": False,
            "protected_heldout_used": False,
            "scientific_acceptance": "NOT_GRANTED",
            "trials": trial_records,
            "summaries": summaries,
        }
        validate_receipt(receipt)
        return receipt
    finally:
        samples.clear()
        torch.cuda.empty_cache()


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--require-cuda", action="store_true")
    parser.add_argument("--source-sha", required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    if not args.require_cuda:
        raise SystemExit("R4_NOT_EXECUTED: explicit --require-cuda is mandatory")
    if args.output.exists():
        raise FileExistsError("refusing to overwrite immutable receipt")
    outcome = run(args.source_sha)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("x", encoding="utf-8") as stream:
        json.dump(outcome, stream, sort_keys=True, indent=2)
        stream.write("\n")
    print("R4_CUDA_TRANSFER_DEVELOPMENT_PASS", len(outcome["trials"]))


if __name__ == "__main__":
    main()
