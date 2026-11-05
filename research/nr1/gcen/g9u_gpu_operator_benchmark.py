"""G9U bounded CUDA dense-vs-packed operator comparison (opt-in only).

Caller supplies authenticated original BF16 expert weight tensor. This is
NOT a CUDA fused kernel: packed INT2 codes decode a bounded 256-column tile
on GPU before GEMM. Peak VRAM is joint-residency, not full-model residency.
No model launch, actual activations, or automatic runner dispatch.
"""
import json
import time
import numpy as np
from g9t_gpu_admission import inspect_current_cuda

BATCH=8
COLUMNS=768
ROWS=2048
BLOCK=256
SEEDS=(751,811)
WARMUP=3
REPEATS=7
LEVELS=(-1.5104,-0.4528,0.4528,1.5104)
SCHEMA="wingless.gcen.g9u.authenticated-original-expert-gpu-operator.v1"


def validate_source(weight, source_sha):
    w=np.asarray(weight)
    if (w.shape!=(ROWS,COLUMNS) or not np.issubdtype(w.dtype,np.floating)
        or not np.isfinite(w).all()):
        raise ValueError("invalid original BF16 expert source")
    if (not isinstance(source_sha,str) or len(source_sha)!=40
        or any(c not in "0123456789abcdef" for c in source_sha)):
        raise ValueError("invalid exact source SHA")
    return w.astype(np.float32,copy=False)


def validate_packed(data, scales):
    d=np.asarray(data)
    s=np.asarray(scales)
    if (d.shape!=(ROWS,COLUMNS//4) or d.dtype!=np.uint8
        or s.shape!=(ROWS,COLUMNS//BLOCK)
        or s.dtype!=np.float16 or not np.isfinite(s).all()):
        raise ValueError("invalid original-expert packed INT2 format")
    return d,s


def run_authenticated_expert(*, shard_path, manifest_path, source_sha, torch_module=None):
    from g9u_verified_loader import load_verified_original_expert
    original_weight, source_custody = load_verified_original_expert(shard_path, manifest_path)
    w=validate_source(original_weight,source_sha)
    from packed_companded_int2 import pack_companded_int2
    d,s=validate_packed(*pack_companded_int2(w,block=BLOCK))
    if torch_module is None:
        import torch
    else:
        torch=torch_module
    admission=inspect_current_cuda(torch,expected_source_sha=source_sha)
    device="cuda:0"
    # Only a tiny authenticated expert, never a full checkpoint.
    original=torch.as_tensor(w,device=device,dtype=torch.float32)
    codes=torch.as_tensor(d,device=device,dtype=torch.uint8)
    scales=torch.as_tensor(s,device=device,dtype=torch.float32)
    levels=torch.as_tensor(LEVELS,device=device,dtype=torch.float32)
    torch.cuda.synchronize()
    torch.cuda.reset_peak_memory_stats()
    resident_before=torch.cuda.memory_allocated()
    def dense(x):
        return x @ original.T
    def packed_tiled(x):
        result=torch.zeros((BATCH,ROWS),device=device,dtype=torch.float32)
        for block in range(COLUMNS//BLOCK):
            tile=codes[:,block*(BLOCK//4):(block+1)*(BLOCK//4)]
            unpacked=torch.stack(tuple((tile>>(2*q))&3 for q in range(4)),dim=-1)
            decoded=levels[unpacked.long()].reshape(ROWS,BLOCK)
            result += (x[:,block*BLOCK:(block+1)*BLOCK] @ decoded.T)*scales[:,block]
        return result
    measurements=[]
    for seed in SEEDS:
        generator=torch.Generator(device="cpu")
        generator.manual_seed(seed)
        x=torch.randn((BATCH,COLUMNS),generator=generator,dtype=torch.float32).to(device)
        ref=dense(x)
        observed=packed_tiled(x)
        numerator=torch.sum((observed.to(torch.float64)-ref.to(torch.float64))**2).item()
        denominator=torch.sum(ref.to(torch.float64)**2).item()
        if not denominator>0 or not np.isfinite(numerator/denominator):
            raise RuntimeError("invalid GPU output error")
        times={}
        for name,fn in (("dense_fp32",dense),("packed_int2_tiled_decode",packed_tiled)):
            for _ in range(WARMUP):
                fn(x)
            torch.cuda.synchronize()
            samples=[]
            for _ in range(REPEATS):
                start=torch.cuda.Event(enable_timing=True)
                end=torch.cuda.Event(enable_timing=True)
                start.record()
                y=fn(x)
                end.record()
                torch.cuda.synchronize()
                if tuple(y.shape)!=(BATCH,ROWS) or not bool(torch.isfinite(y).all()):
                    raise RuntimeError("invalid GPU operator output")
                samples.append(float(start.elapsed_time(end)))
            times[name]=float(np.median(samples))
        measurements.append({
            "seed":seed,"dense_fp32_median_gpu_ms":times["dense_fp32"],
            "packed_int2_median_gpu_ms":times["packed_int2_tiled_decode"],
            "packed_speedup_vs_dense":times["dense_fp32"]/times["packed_int2_tiled_decode"],
            "relative_output_squared_error_vs_dense":numerator/denominator,
        })
    torch.cuda.synchronize()
    peak=torch.cuda.max_memory_allocated()
    return {
        "schema":SCHEMA,"exact_source_sha":source_sha,
        "hardware_admission":admission,"input_type":"FROZEN_SYNTHETIC_GAUSSIAN",
        "seeds":list(SEEDS),"shape":[ROWS,COLUMNS],"batch":BATCH,
        "warmup":WARMUP,"repeats":REPEATS,
        "dense_operator":"FP32_CUDA_GEMM",
        "packed_operator":"INT2_GPU_TILED_DECODE_PLUS_FP32_GEMM_NOT_FUSED",
        "packed_storage_bytes":int(d.nbytes+s.nbytes),
        "original_bf16_storage_bytes":int(ROWS*COLUMNS*2),
        "gpu_resident_baseline_bytes":int(resident_before),
        "gpu_joint_residency_peak_allocated_bytes":int(peak),
        "measurements":measurements,
        "original_shard_authenticated_in_process":True,
        "original_shard_custody":source_custody,
        "original_activation_traces":False,"full_model_forward":False,
        "scientific_acceptance":"NOT_GRANTED",
    }
