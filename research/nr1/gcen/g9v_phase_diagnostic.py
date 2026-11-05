"""G9V: one bounded original-expert CUDA phase attribution experiment.

Reuses the existing G9U representation/seed protocol. No checkpoint
download, full-model inference, authentic activations, or model deployment.
"""
import argparse
import json
from pathlib import Path
import numpy as np
from g9u_run_authenticated import check_checkout
from g9u_gpu_operator_benchmark import (
    BATCH, ROWS, COLUMNS, BLOCK, LEVELS, SEEDS, WARMUP, REPEATS,
    validate_source, validate_packed, verify_packed_operator_parity,
)
from g9t_gpu_admission import inspect_current_cuda

PHASES = ("allocate", "unpack", "decode", "gemm", "scale_accumulate")
SCHEMA = "wingless.gcen.g9v.authenticated-gpu-phase-diagnostic.v1"
MAX_SYNTHETIC_ERROR = 0.02
MIN_DECODE_SHARE = 0.60
MAX_TRACE_INFLATION = 1.50
MAX_SEED_LATENCY_RATIO = 1.20


def decide(rows):
    """Predeclared screening, never automatic scientific acceptance."""
    if len(rows) != 2 or {r.get("seed") for r in rows} != set(SEEDS):
        raise ValueError("frozen seed pair required")
    for row in rows:
        vals = [row.get(k) for k in ("dense_median_ms", "packed_median_ms",
                "trace_median_ms", "relative_squared_error",
                "unpack_decode_fraction")]
        if any(type(v) not in (float, int) or not np.isfinite(v) for v in vals):
            raise ValueError("invalid diagnostic metric")
        if min(vals[:3]) <= 0 or not 0 <= vals[3] <= 1 or not 0 <= vals[4] <= 1:
            raise ValueError("invalid diagnostic metric")
    if max(r["relative_squared_error"] for r in rows) > MAX_SYNTHETIC_ERROR:
        return "STOP_INT2_FIDELITY_SCREEN"
    speeds = [r["packed_median_ms"] for r in rows]
    if max(speeds) / min(speeds) > MAX_SEED_LATENCY_RATIO:
        return "NO_GO_UNSTABLE_GPU_LATENCY"
    if any(r["trace_median_ms"] / r["packed_median_ms"] > MAX_TRACE_INFLATION for r in rows):
        return "NO_GO_TRACE_OVERHEAD"
    if min(r["unpack_decode_fraction"] for r in rows) < MIN_DECODE_SHARE:
        return "NO_GO_NO_DOMINANT_DECODE"
    return "POTENTIAL_FUSION_STUDY_NOT_SCIENTIFIC_PASS"


def event_bench(torch, fn, x):
    for _ in range(WARMUP):
        fn(x)
    torch.cuda.synchronize()
    times = []
    for _ in range(REPEATS):
        start = torch.cuda.Event(enable_timing=True)
        end = torch.cuda.Event(enable_timing=True)
        start.record()
        output = fn(x)
        end.record()
        torch.cuda.synchronize()
        if tuple(output.shape) != (BATCH, ROWS) or not bool(torch.isfinite(output).all()):
            raise ValueError("invalid GPU operator output")
        times.append(float(start.elapsed_time(end)))
    return times


def trace_once(torch, x, codes, scales, levels):
    """Same tile and in-place accumulation as G9U, with CUDA event brackets."""
    spans = {name: [] for name in PHASES}
    def measure(name, fn):
        a = torch.cuda.Event(enable_timing=True)
        b = torch.cuda.Event(enable_timing=True)
        a.record()
        result = fn()
        b.record()
        spans[name].append((a, b))
        return result
    a = torch.cuda.Event(enable_timing=True)
    b = torch.cuda.Event(enable_timing=True)
    a.record()
    result = measure("allocate", lambda: torch.zeros(
        (BATCH, ROWS), device="cuda:0", dtype=torch.float32))
    for block in range(COLUMNS // BLOCK):
        tile = codes[:, block*(BLOCK//4):(block+1)*(BLOCK//4)]
        unpacked = measure("unpack", lambda: torch.stack(
            tuple((tile >> (2*q)) & 3 for q in range(4)), dim=-1))
        decoded = measure("decode", lambda: levels[unpacked.long()].reshape(ROWS, BLOCK))
        product = measure("gemm", lambda: x[:, block*BLOCK:(block+1)*BLOCK] @ decoded.T)
        measure("scale_accumulate", lambda: result.add_(product * scales[:,block]))
    b.record()
    torch.cuda.synchronize()
    if tuple(result.shape) != (BATCH, ROWS) or not bool(torch.isfinite(result).all()):
        raise ValueError("invalid phase trace output")
    phase = {name: sum(float(x.elapsed_time(y)) for x,y in spans[name])
             for name in PHASES}
    return result, float(a.elapsed_time(b)), phase


def run(*, root, source_sha, shard, manifest, output):
    check_checkout(root, source_sha)
    output = Path(output)
    if output.exists() or output.is_symlink() or not output.parent.is_dir():
        raise ValueError("refuse missing or existing diagnostic receipt")
    from g9u_verified_loader import load_verified_original_expert
    w, custody = load_verified_original_expert(shard, manifest)
    w = validate_source(w, source_sha)
    from packed_companded_int2 import pack_companded_int2, apply_companded_int2
    d, s = validate_packed(*pack_companded_int2(w, block=BLOCK))
    import torch
    admission = inspect_current_cuda(torch, expected_source_sha=source_sha)
    original = torch.as_tensor(w, device="cuda:0", dtype=torch.float32)
    codes = torch.as_tensor(d, device="cuda:0", dtype=torch.uint8)
    scales = torch.as_tensor(s, device="cuda:0", dtype=torch.float32)
    levels = torch.as_tensor(LEVELS, device="cuda:0", dtype=torch.float32)
    def dense(x):
        return x @ original.T
    def packed(x):
        result = torch.zeros((BATCH, ROWS), device="cuda:0", dtype=torch.float32)
        for block in range(COLUMNS // BLOCK):
            tile = codes[:,block*(BLOCK//4):(block+1)*(BLOCK//4)]
            unpacked = torch.stack(tuple((tile >> (2*q)) & 3 for q in range(4)), dim=-1)
            decoded = levels[unpacked.long()].reshape(ROWS,BLOCK)
            result += (x[:,block*BLOCK:(block+1)*BLOCK] @ decoded.T) * scales[:,block]
        return result
    torch.cuda.synchronize()
    torch.cuda.reset_peak_memory_stats()
    rows = []
    for seed in SEEDS:
        generator = torch.Generator(device="cpu")
        generator.manual_seed(seed)
        x = torch.randn((BATCH, COLUMNS), generator=generator, dtype=torch.float32).to("cuda:0")
        dense_output = dense(x)
        packed_output = packed(x)
        cpu = apply_companded_int2(x.cpu().numpy(), d, s, block=BLOCK)
        parity = verify_packed_operator_parity(cpu, packed_output.detach().cpu().numpy())
        numerator = torch.sum((packed_output.double() - dense_output.double())**2).item()
        denominator = torch.sum(dense_output.double()**2).item()
        if denominator <= 0:
            raise RuntimeError("invalid original reference")
        dense_samples = event_bench(torch, dense, x)
        packed_samples = event_bench(torch, packed, x)
        negative = []
        for _ in range(REPEATS):
            a = torch.cuda.Event(enable_timing=True)
            b = torch.cuda.Event(enable_timing=True)
            a.record()
            b.record()
            torch.cuda.synchronize()
            negative.append(float(a.elapsed_time(b)))
        for _ in range(WARMUP):
            trace_once(torch, x, codes, scales, levels)
        trace_totals = []
        phases = {key: [] for key in PHASES}
        for _ in range(REPEATS):
            traced, elapsed, phase = trace_once(torch, x, codes, scales, levels)
            verify_packed_operator_parity(packed_output.detach().cpu().numpy(),
                                          traced.detach().cpu().numpy())
            trace_totals.append(elapsed)
            for name in PHASES:
                phases[name].append(phase[name])
        medians = {name: float(np.median(phases[name])) for name in PHASES}
        denom = sum(medians.values())
        if denom <= 0:
            raise RuntimeError("no measured phase costs")
        rows.append({
            "seed": seed, "dense_gpu_ms": dense_samples, "packed_gpu_ms": packed_samples,
            "event_only_gpu_ms": negative, "trace_gpu_ms": trace_totals,
            "phase_samples_gpu_ms": phases, "phase_median_gpu_ms": medians,
            "dense_median_ms": float(np.median(dense_samples)),
            "packed_median_ms": float(np.median(packed_samples)),
            "trace_median_ms": float(np.median(trace_totals)),
            "unpack_decode_fraction": (medians["unpack"] + medians["decode"]) / denom,
            "relative_squared_error": float(numerator / denominator),
            "cpu_packed_parity": parity,
        })
    receipt = {
        "schema": SCHEMA, "source_repo": "spoonman136668-ai/Wingless",
        "source_sha": source_sha, "source_custody": custody, "hardware_admission": admission,
        "protocol": {"input": "FROZEN_SYNTHETIC_GAUSSIAN_NOT_ORIGINAL_ACTIVATIONS",
                     "shape": [ROWS,COLUMNS], "batch": BATCH, "seeds": list(SEEDS),
                     "warmup": WARMUP, "repeats": REPEATS, "phases": list(PHASES),
                     "max_synthetic_error": MAX_SYNTHETIC_ERROR,
                     "min_decode_fraction": MIN_DECODE_SHARE,
                     "max_trace_inflation": MAX_TRACE_INFLATION,
                     "max_seed_latency_ratio": MAX_SEED_LATENCY_RATIO},
        "measurements": rows, "decision": decide(rows),
        "gpu_joint_residency_peak_bytes": int(torch.cuda.max_memory_allocated()),
        "original_activations": False, "full_model_forward": False,
        "scientific_acceptance": "NOT_GRANTED", "deployment": False,
    }
    with output.open("x", encoding="utf-8") as sink:
        json.dump(receipt, sink, sort_keys=True, indent=2)
        sink.write("\n")
    print("G9V_PHASE_RESULT " + receipt["decision"])
    return receipt


if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--root", type=Path, required=True)
    p.add_argument("--source-sha", required=True)
    p.add_argument("--shard", type=Path, required=True)
    p.add_argument("--manifest", type=Path, required=True)
    p.add_argument("--output", type=Path, required=True)
    args = p.parse_args()
    run(root=args.root, source_sha=args.source_sha,
        shard=args.shard, manifest=args.manifest, output=args.output)
