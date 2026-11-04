"""G23 diagnostic host-capacity census for the pinned original Qwen checkpoint.

No weights fetched, model loaded, runtime activated, or execution authority granted.
RAM + NVIDIA VRAM total is an optimistic upper bound for simultaneously resident
original weights. Disk/paged/offloaded inference remains unverified, not rejected.
"""
import argparse
import ctypes
import json
import os
import pathlib
import platform
import shutil
import subprocess

MODEL = "Qwen/Qwen3-30B-A3B"
REVISION = "ad44e777bcd18fa416d9da3bd8f70d33ebb85d39"
# Source: verified 16-shard G14/G15 custody reconciliations, PR #440.
# Storage size is an optimistic weight-residency floor, not peak inference RSS/VRAM.
PINNED_ORIGINAL_SHARD_BYTES = 61066575648


def check_manifest(manifest):
    if not isinstance(manifest, dict) or manifest.get("model") != MODEL or manifest.get("revision") != REVISION:
        raise ValueError("wrong original-model identity")
    files = manifest.get("files")
    shards = [f"model-{i:05d}-of-00016.safetensors" for i in range(1, 17)]
    required = set(shards + ["config.json", "model.safetensors.index.json"])
    if not isinstance(files, dict) or set(files) != required:
        raise ValueError("incomplete original checkpoint manifest")
    if any(not isinstance(x, str) or len(x) != 64 or any(c not in "0123456789abcdef" for c in x) for x in files.values()):
        raise ValueError("invalid checkpoint digests")
    return True


def capacity_assessment(ram_total_bytes, gpu_total_bytes, *, runner_qualified=False):
    if type(ram_total_bytes) is not int or ram_total_bytes <= 0:
        raise ValueError("RAM total absent or invalid")
    if gpu_total_bytes is not None and (type(gpu_total_bytes) is not int or gpu_total_bytes < 0):
        raise ValueError("GPU total invalid")
    optimistic_total = None if gpu_total_bytes is None else ram_total_bytes + gpu_total_bytes
    if gpu_total_bytes is None:
        status = "GPU_CAPACITY_UNKNOWN_NO_RESIDENCY_ADMISSION"
    elif optimistic_total < PINNED_ORIGINAL_SHARD_BYTES:
        status = "INSUFFICIENT_FOR_SIMULTANEOUS_ORIGINAL_WEIGHT_RESIDENCY"
    else:
        status = "CAPACITY_LOWER_BOUND_MET_RUNTIME_UNVERIFIED"
    return {
        "status": status,
        "ram_total_bytes": ram_total_bytes,
        "nvidia_vram_total_bytes": gpu_total_bytes,
        "optimistic_combined_bytes": optimistic_total,
        "original_shard_bytes_lower_bound": PINNED_ORIGINAL_SHARD_BYTES,
        "runner_qualified_claim": False,
        "runner_qualification_independently_verified": False,
        "original_model_forward_verified": False,
        "execution_authorized_by_this_probe": False,
        "offload_or_paging_feasibility": "UNVERIFIED",
        "interpretation": "Capacity necessary only for simultaneous full-weight residency; neither a sufficient runtime test nor an offload impossibility result.",
    }


def physical_ram():
    system = platform.system()
    if system == "Linux":
        meminfo = pathlib.Path("/proc/meminfo").read_text()
        import re
        match = re.search(r"^MemTotal:\s+(\d+)\s+kB$", meminfo, re.M)
        if match is None:
            raise RuntimeError("Linux memory total unavailable")
        total = int(match.group(1)) * 1024
        cgroup = pathlib.Path("/sys/fs/cgroup/memory.max")
        if cgroup.is_file():
            limit = cgroup.read_text().strip()
            if limit.isdigit() and int(limit) > 0:
                total = min(total, int(limit))
        return total
    if system == "Windows":
        class MEMORYSTATUSEX(ctypes.Structure):
            _fields_ = [("dwLength",ctypes.c_ulong), ("dwMemoryLoad",ctypes.c_ulong),
                ("ullTotalPhys",ctypes.c_ulonglong), ("ullAvailPhys",ctypes.c_ulonglong),
                ("ullTotalPageFile",ctypes.c_ulonglong), ("ullAvailPageFile",ctypes.c_ulonglong),
                ("ullTotalVirtual",ctypes.c_ulonglong), ("ullAvailVirtual",ctypes.c_ulonglong),
                ("ullAvailExtendedVirtual",ctypes.c_ulonglong)]
        status = MEMORYSTATUSEX()
        status.dwLength = ctypes.sizeof(status)
        if not ctypes.windll.kernel32.GlobalMemoryStatusEx(ctypes.byref(status)):
            raise RuntimeError("Windows physical RAM query failed")
        return int(status.ullTotalPhys)
    raise RuntimeError("Unsupported host for physical RAM measurement")


def nvidia_vram():
    command = shutil.which("nvidia-smi")
    if not command:
        return None
    try:
        result = subprocess.run([command, "--query-gpu=memory.total", "--format=csv,noheader,nounits"],
            text=True, capture_output=True, timeout=8, check=True)
        amounts = [int(line.strip()) for line in result.stdout.splitlines() if line.strip()]
        if not amounts or any(value <= 0 for value in amounts):
            return None
        # NVIDIA-SMI reports MiB, not bytes. Includes every visible device.
        return sum(amounts) * 1024 * 1024
    except (ValueError, subprocess.SubprocessError, OSError):
        return None


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", type=pathlib.Path, required=True)
    parser.add_argument("--output", type=pathlib.Path, required=True)
    args = parser.parse_args()
    check_manifest(json.loads(args.manifest.read_text(encoding="utf-8")))
    result = capacity_assessment(physical_ram(), nvidia_vram())
    result.update({"schema": "wingless.gcen.g23.capacity-diagnostic.v1",
        "platform": platform.system(), "checkpoint_model": MODEL, "checkpoint_revision": REVISION,
        "qualification_source": "G14/G15 custody original-shard receipt; G22 exact-head code qualification",
        "source_authenticity": "PREVIOUS_SHARD_CUSTODY_ONLY_NO_CURRENT_SESSION_CHECKPOINT",
    })
    if args.output.exists():
        raise FileExistsError("Refusing to overwrite existing capacity receipt")
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print("GCEN_G23_CAPACITY_DIAGNOSTIC " + result["status"])


if __name__ == "__main__":
    main()
