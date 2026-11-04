"""G9L: independent original BF16 low-rank corrected INT2 versus original INT4 control.

Reuses the original G8 SHA256-verification, frozen shard, 16 experts,
five fit seeds and 128x256 windows. Test inputs are synthetic Gaussian
probes, NOT native model activations. No Qwen forward or GPU allocation.
"""
import argparse
import json
import os
from pathlib import Path
import re
import shutil
import tempfile
import time
from urllib.request import Request, urlopen

import numpy as np
from collect_upstream_manifest import build_manifest
from expert_scaling_comparison import COUNTS, COLS, ROWS, SEEDS as PRIOR_SEEDS
from g31_original_shard_stage import (
    MAX_SHARD_BYTES, RESERVE_BYTES, SHARD, URL, check_upstream,
    inspect_header, parse_unique_json, peak_rss_bytes, stage_verified,
)
from g32_real_expert_scaling import read_bf16_window, select_experts
from out_of_core_preflight import pinned_metadata
from packed_int2 import pack_int2, apply_packed_int2, stored_bytes as stored2
from packed_companded_int2 import pack_companded_int2, apply_companded_int2, stored_bytes as compbytes
from lowrank_residual_int2 import RANKS, pack_lowrank_residual, apply_lowrank_residual, stored_bytes as lowrankbytes
from g9k_lowrank_budget import report as lowrank_budget
from packed_int4 import pack_int4, apply_packed_int4, stored_bytes as stored4
from runtime_feasibility import MODEL, REVISION, check_manifest

SOURCE_SHA = re.compile(r"[0-9a-f]{40}\Z")
BATCH = 32
SEEDS = (541, 563, 587, 613, 641)  # Independent of G9E/G9F/G9I development
assert not (set(SEEDS) & set(PRIOR_SEEDS))
BLOCKS = (128, 256)
RESULT = "G9L_ORIGINAL_BF16_LOWRANK_RESIDUAL_COMPARISON_COMPLETE"


def compare_expert(weight, x):
    w=np.asarray(weight,dtype=np.float32)
    x=np.asarray(x,dtype=np.float32)
    if (w.shape!=(ROWS,COLS) or x.shape!=(BATCH,COLS) or
        not np.isfinite(w).all() or not np.isfinite(x).all() or
        not np.any(w) or not np.any(x)):
        raise ValueError("invalid frozen original expert weight and inputs")
    reference=x@w.T
    energy=float(np.sum(reference.astype(np.float64)**2))
    if energy<=0 or not np.isfinite(energy):
        raise ValueError("invalid original output energy")
    p4,s4=pack_int4(w,block=64)
    int4=apply_packed_int4(x,p4,s4,block=64)
    pc,sc=pack_companded_int2(w,block=256)
    baseline=apply_companded_int2(x,pc,sc,block=256)
    out={
        "output_energy":energy,"original_bf16_bytes":int(w.size*2),
        "int4_bytes":stored4(p4,s4,COLS,block=64),
        "int4_squared_error":float(np.sum((int4.astype(np.float64)-reference.astype(np.float64))**2)),
        "companded_int2_bytes":compbytes(pc,sc,COLS,block=256),
        "companded_squared_error":float(np.sum((baseline.astype(np.float64)-reference.astype(np.float64))**2)),
        "lowrank":{}
    }
    for rank in RANKS:
        packed=pack_lowrank_residual(w,rank=rank,block=256)
        predicted=apply_lowrank_residual(x,*packed,rank=rank,block=256)
        size=lowrankbytes(*packed,COLS,rank=rank,block=256)
        if size<=out["companded_int2_bytes"] or not np.isfinite(predicted).all():
            raise ValueError("lowrank corrected source invalid")
        out["lowrank"][str(rank)]={
            "bytes":size,
            "output_squared_error":float(np.sum((predicted.astype(np.float64)-reference.astype(np.float64))**2))
        }
    return out


def compare_windows(windows,seed):
    if type(seed) is not int or seed not in SEEDS:
        raise ValueError("unfrozen original selector")
    w=np.asarray(windows,dtype=np.float32)
    if w.shape!=(max(COUNTS),ROWS,COLS) or not np.isfinite(w).all():
        raise ValueError("invalid 16-expert original inventory")
    x=np.random.default_rng(71000+seed).standard_normal((BATCH,COLS)).astype(np.float32)
    rows=[compare_expert(expert,x) for expert in w]
    results=[]
    for n in COUNTS:
        selected=rows[:n]
        output_energy=sum(r["output_energy"] for r in selected)
        original_bytes=sum(r["original_bf16_bytes"] for r in selected)
        cases=[]
        for rank in RANKS:
            records=[r["lowrank"][str(rank)] for r in selected]
            total=sum(q["bytes"] for q in records)
            cases.append({
                "rank":rank,
                "stored_bytes_including_fp16_factors_and_headers":total,
                "compression_ratio_vs_bf16":original_bytes/total,
                "relative_output_squared_error":
                    sum(q["output_squared_error"] for q in records)/output_energy,
            })
        results.append({
            "seed":seed,"experts":n,
            "original_bf16_bytes":original_bytes,
            "int4_bytes":sum(r["int4_bytes"] for r in selected),
            "int4_relative_output_squared_error":
                sum(r["int4_squared_error"] for r in selected)/output_energy,
            "companded_int2_bytes":sum(r["companded_int2_bytes"] for r in selected),
            "companded_relative_output_squared_error":
                sum(r["companded_squared_error"] for r in selected)/output_energy,
            "lowrank_cases":cases,
            "inputs":"SYNTHETIC_GAUSSIAN_NOT_ORIGINAL_QWEN_ACTIVATIONS"
        })
    return results


def run(manifest_path, scratch_root, output, source_sha):
    if os.name != "posix" or type(source_sha) is not str or not SOURCE_SHA.fullmatch(source_sha):
        raise ValueError("requires a pinned source SHA and Linux custody probes")
    if not scratch_root.is_dir() or output.exists():
        raise ValueError("missing scratch directory or pre-existing receipt")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    if shutil.disk_usage(scratch_root).free < MAX_SHARD_BYTES + RESERVE_BYTES:
        raise RuntimeError("insufficient one-shard scratch space")
    check_upstream(manifest, build_manifest())
    index = pinned_metadata("model.safetensors.index.json", manifest)
    start = time.monotonic()
    rss_start = peak_rss_bytes()
    with tempfile.TemporaryDirectory(prefix="cnc-g9l-original-", dir=scratch_root) as td:
        staged = Path(td) / SHARD
        request = Request(URL, headers={"Accept-Encoding": "identity",
                          "User-Agent": "Wingless-CNC-G9L-Development/1"})
        with urlopen(request, timeout=90) as response:
            custody = stage_verified(response, staged, manifest["files"][SHARD])
        header_receipt = inspect_header(staged, index)
        with staged.open("rb") as stream:
            header_bytes = int.from_bytes(stream.read(8), "little")
            if header_bytes != header_receipt["header_bytes"]:
                raise ValueError("verified header size changed")
            header = parse_unique_json(stream.read(header_bytes))
            selected = select_experts(header, index["weight_map"])
            results = []
            for seed in SEEDS:
                rng = np.random.default_rng(seed)
                row = int(rng.integers(0, selected["shape"][0] - ROWS + 1))
                col = int(rng.integers(0, selected["shape"][1] - COLS + 1))
                windows = np.stack([
                    read_bf16_window(stream, header_bytes, off, selected["shape"], row, col)
                    for off in selected["offsets"]])
                for result in compare_windows(windows, seed):
                    result.update(source_window_start=[row, col],
                                  original_selected_layer=selected["layer"],
                                  distinct_original_expert_names=selected["names"][:result["experts"]])
                    results.append(result)
        if len(results) != len(SEEDS) * len(COUNTS):
            raise ValueError("incomplete original-seed comparisons")
        if custody["observed_peak_staged_file_count"] != 1 or list(Path(td).iterdir()) != [staged]:
            raise ValueError("original shard custody invariant failed")
        rss_peak = peak_rss_bytes()
    receipt = {
        "schema": "wingless.cnc.g9l.real-original-lowrank-development.v1",
        "result": RESULT, "model": MODEL, "revision": REVISION,
        "source_sha": source_sha, "original_shard": SHARD,
        "original_shard_sha256": manifest["files"][SHARD],
        "source_custody": custody, "original_header": header_receipt,
        "frozen_seeds": list(SEEDS), "frozen_expert_groups": list(COUNTS),
        "synthetic_activation_seed_offset": 71000,
        "original_seed_overlap_with_prior": list(set(SEEDS) & set(PRIOR_SEEDS)),
        "frozen_lowrank_ranks": list(RANKS),
        "analytical_model_budget": lowrank_budget(),
        "window_shape": [ROWS, COLS], "synthetic_batch": BATCH,
        "comparison_count": len(results), "comparisons": results,
        "int2_blocks": list(BLOCKS), "int2_codebook": [-3, -1, 1, 3],
        "int2_scales_dtype": "float16",
        "companded_codebook": [-1.5104, -0.4528, 0.4528, 1.5104],
        "companded_scale_kind": "FP16_BLOCK_RMS",
        "int4_control": "block64_float32_scales",
        "rss_baseline_bytes": rss_start, "rss_peak_bytes": rss_peak,
        "elapsed_seconds": time.monotonic() - start,
        "staged_original_files_after_cleanup": 0,
        "authentic_model_activations": "NOT_CAPTURED",
        "full_model_inference": False, "gpu_measurement": False,
        "protected_heldout_used": False, "scientific_acceptance": "NOT_GRANTED",
        "production_activation": False,
    }
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8") as stream:
        json.dump(receipt, stream, sort_keys=True, indent=2)
        stream.write("\n")
    print("CNC_G9L_VERIFIED_ORIGINAL_LOWRANK_WINDOW_COMPARISON", len(results))


if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--manifest", type=Path, required=True)
    p.add_argument("--scratch-root", type=Path, required=True)
    p.add_argument("--source-sha", required=True)
    p.add_argument("--output", type=Path, required=True)
    a = p.parse_args()
    run(a.manifest, a.scratch_root, a.output, a.source_sha)
