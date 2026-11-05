"""G9T GPU admission: current CUDA evidence only; fail closed, no model activation.

The archived G9P GPU receipt is *not* a live admission token. This probe
must run in the intended isolated GPU process immediately before a separately
authorized benchmark. It never downloads source, transfers weights, or
activates a Qwen runtime.
"""
import json
import time

MIB=1<<20
MIN_FREE_MIB=1024
REQUIRED_CUDA_DEVICE="NVIDIA"
SOURCE_REPOSITORY="spoonman136668-ai/Wingless"
SCHEMA="wingless.gcen.g9t.current-gpu-admission.v1"


def _exact_int(value, label):
    if type(value) is not int or value <= 0:
        raise ValueError("invalid GPU "+label)
    return value


def inspect_current_cuda(torch_module, *, expected_source_sha, minimum_free_mib=MIN_FREE_MIB):
    """Measure selected device, never automatically select another device.

    This captures only a point-in-time capacity/identity receipt. It cannot
    grant workstation ownership, GPU execution authority or model acceptance.
    """
    if (not isinstance(expected_source_sha,str)
        or len(expected_source_sha)!=40
        or any(c not in "0123456789abcdef" for c in expected_source_sha)):
        raise ValueError("exact source SHA required")
    if type(minimum_free_mib) is not int or minimum_free_mib<MIN_FREE_MIB:
        raise ValueError("minimum free VRAM floor cannot be weakened")
    cuda=getattr(torch_module,"cuda",None)
    if cuda is None or not cuda.is_available():
        raise RuntimeError("CUDA unavailable: reject CPU fallback")
    if _exact_int(cuda.device_count(),"device count")!=1:
        raise RuntimeError("expected isolated single-GPU runner")
    device=cuda.current_device()
    if type(device) is not int or device!=0:
        raise RuntimeError("unexpected active GPU device")
    properties=cuda.get_device_properties(device)
    name=str(properties.name)
    if not name.startswith(REQUIRED_CUDA_DEVICE):
        raise RuntimeError("non-NVIDIA device is outside frozen CUDA scope")
    total=_exact_int(int(properties.total_memory),"total VRAM")
    free,reported_total=cuda.mem_get_info(device)
    free=_exact_int(free,"free VRAM")
    reported_total=_exact_int(reported_total,"reported total VRAM")
    if total!=reported_total or free>total:
        raise RuntimeError("CUDA physical/current memory counters disagree")
    floor=minimum_free_mib*MIB
    if free<floor:
        raise RuntimeError("measured GPU free VRAM below isolated benchmark floor")
    capability=cuda.get_device_capability(device)
    if (not isinstance(capability,(tuple,list)) or len(capability)!=2
        or any(type(v) is not int or v<0 for v in capability)):
        raise RuntimeError("invalid CUDA capability")
    torch_version=str(torch_module.__version__)
    cuda_version=getattr(torch_module.version,"cuda",None)
    if not cuda_version:
        raise RuntimeError("CPU-only or incomplete Torch CUDA build")
    return {
        "schema":SCHEMA,
        "source_repo":SOURCE_REPOSITORY,
        "exact_source_sha":expected_source_sha,
        "gpu_name":name,
        "device_index":device,
        "compute_capability":list(capability),
        "torch_version":torch_version,
        "torch_cuda_version":str(cuda_version),
        "physical_total_vram_bytes":total,
        "currently_free_vram_bytes":free,
        "minimum_free_floor_bytes":floor,
        "probe_clock_utc_epoch_seconds":time.time(),
        "probe_is_current_in_this_process":True,
        "historical_gpu_receipt_used_for_current_admission":False,
        "original_shards_loaded":False,
        "cuda_kernel_executed":False,
        "peak_vram_measured":False,
        "runner_use_authorized_by_probe":False,
        "full_model_forward":False,
        "scientific_acceptance":"NOT_GRANTED",
    }


def command_line():
    import argparse
    p=argparse.ArgumentParser()
    p.add_argument("--source-sha",required=True)
    p.add_argument("--minimum-free-mib",type=int,default=MIN_FREE_MIB)
    a=p.parse_args()
    import torch
    print("G9T_CURRENT_GPU_ADMISSION "+json.dumps(
        inspect_current_cuda(torch,expected_source_sha=a.source_sha,
                             minimum_free_mib=a.minimum_free_mib),
        sort_keys=True))


if __name__=="__main__":
    command_line()
