"""G9F: original Qwen BF16 RMS-companded INT2 versus frozen uniform/INT4 controls.

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
from expert_scaling_comparison import COUNTS, COLS, ROWS, SEEDS
from g31_original_shard_stage import (
    MAX_SHARD_BYTES, RESERVE_BYTES, SHARD, URL, check_upstream,
    inspect_header, parse_unique_json, peak_rss_bytes, stage_verified,
)
from g32_real_expert_scaling import read_bf16_window, select_experts
from out_of_core_preflight import pinned_metadata
from packed_int2 import pack_int2, apply_packed_int2, stored_bytes as stored2
from packed_companded_int2 import pack_companded_int2, apply_companded_int2, stored_bytes as compbytes
from packed_int4 import pack_int4, apply_packed_int4, stored_bytes as stored4
from runtime_feasibility import MODEL, REVISION, check_manifest

SOURCE_SHA = re.compile(r"[0-9a-f]{40}\Z")
BATCH = 32
BLOCKS = (128, 256)
RESULT = "G9F_ORIGINAL_BF16_COMPANDED_WINDOW_COMPARISON_COMPLETE"


def compare_expert(weight, x):
    w = np.asarray(weight, dtype=np.float32)
    x = np.asarray(x, dtype=np.float32)
    if (w.shape != (ROWS, COLS) or x.shape != (BATCH, COLS)
        or not np.isfinite(w).all() or not np.isfinite(x).all()
        or not np.any(w) or not np.any(x)):
        raise ValueError("invalid frozen original window and synthetic inputs")
    ref = x @ w.T
    we = float(np.sum(w.astype(np.float64) ** 2))
    oe = float(np.sum(ref.astype(np.float64) ** 2))
    if we <= 0 or oe <= 0:
        raise ValueError("zero reference energy")
    p4, s4 = pack_int4(w, block=64)
    q4 = apply_packed_int4(x, p4, s4, block=64)
    i4 = float(np.sum((q4.astype(np.float64) - ref.astype(np.float64)) ** 2))
    row = {
        "weight_energy": we, "output_energy": oe,
        "original_bf16_bytes": int(w.size * 2),
        "int4_storage_bytes": stored4(p4, s4, COLS, block=64),
        "int4_output_squared_error": i4, "int2": {}, "companded": {},
    }
    for block in BLOCKS:
        packed, scales = pack_int2(w, block=block)
        y = apply_packed_int2(x, packed, scales, block=block)
        recovered = apply_packed_int2(np.eye(COLS, dtype=np.float32),
                                      packed, scales, block=block).T
        if not np.isfinite(y).all() or not np.isfinite(recovered).all():
            raise ValueError("non-finite INT2 comparison")
        row["int2"][str(block)] = {
            "bytes": stored2(packed, scales, COLS, block=block),
            "weight_squared_error": float(np.sum(
                (recovered.astype(np.float64) - w.astype(np.float64)) ** 2)),
            "output_squared_error": float(np.sum(
                (y.astype(np.float64) - ref.astype(np.float64)) ** 2)),
        }
        cp, cs = pack_companded_int2(w, block=block)
        cy = apply_companded_int2(x, cp, cs, block=block)
        cw = apply_companded_int2(np.eye(COLS, dtype=np.float32),
                                  cp, cs, block=block).T
        if not np.isfinite(cy).all() or not np.isfinite(cw).all():
            raise ValueError("nonfinite companded comparison")
        row["companded"][str(block)] = {
            "bytes": compbytes(cp, cs, COLS, block=block),
            "weight_squared_error": float(np.sum(
                (cw.astype(np.float64) - w.astype(np.float64)) ** 2)),
            "output_squared_error": float(np.sum(
                (cy.astype(np.float64) - ref.astype(np.float64)) ** 2)),
        }
        if row["companded"][str(block)]["bytes"] != row["int2"][str(block)]["bytes"]:
            raise ValueError("G9F requires exactly matched storage")

    return row


def compare_windows(windows, seed):
    if type(seed) is not int or seed not in SEEDS:
        raise ValueError("unfrozen original sampling seed")
    w = np.asarray(windows, dtype=np.float32)
    if w.shape != (max(COUNTS), ROWS, COLS) or not np.isfinite(w).all():
        raise ValueError("invalid 16-expert BF16 window inventory")
    x = np.random.default_rng(26000 + seed).standard_normal(
        (BATCH, COLS)).astype(np.float32)
    rows = [compare_expert(expert, x) for expert in w]
    out = []
    for n in COUNTS:
        selected = rows[:n]
        energy = sum(v["output_energy"] for v in selected)
        weight_energy = sum(v["weight_energy"] for v in selected)
        int4 = sum(v["int4_output_squared_error"] for v in selected) / energy
        cases = []
        for block in BLOCKS:
            values = [v["int2"][str(block)] for v in selected]
            comp = [v["companded"][str(block)] for v in selected]
            stored = sum(v["bytes"] for v in values)
            if sum(v["bytes"] for v in comp) != stored:
                raise ValueError("G9F aggregate storage mismatch")
            originals = sum(v["original_bf16_bytes"] for v in selected)
            cases.append({
                "block": block,
                "packed_storage_bytes_including_metadata": stored,
                "compression_ratio_vs_bf16": originals / stored,
                "relative_output_squared_error": sum(
                    v["output_squared_error"] for v in values) / energy,
                "companded_relative_output_squared_error": sum(
                    v["output_squared_error"] for v in comp) / energy,
                "companded_relative_weight_squared_error": sum(
                    v["weight_squared_error"] for v in comp) / weight_energy,
                "relative_weight_squared_error": sum(
                    v["weight_squared_error"] for v in values) / weight_energy,
                "per_expert_output_relative_error": [
                    v["int2"][str(block)]["output_squared_error"] / v["output_energy"]
                    for v in selected],
            })
        out.append({
            "seed": seed, "experts": n,
            "original_bf16_bytes": sum(v["original_bf16_bytes"] for v in selected),
            "int4_storage_bytes_including_metadata": sum(v["int4_storage_bytes"] for v in selected),
            "int4_relative_output_squared_error": int4,
            "int2_cases": cases,
            "inputs": "SYNTHETIC_GAUSSIAN_NOT_NATIVE_QWEN_ACTIVATIONS",
        })
    return out


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
    with tempfile.TemporaryDirectory(prefix="cnc-g9e-original-", dir=scratch_root) as td:
        staged = Path(td) / SHARD
        request = Request(URL, headers={"Accept-Encoding": "identity",
                          "User-Agent": "Wingless-CNC-G9E-Development/1"})
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
        "schema": "wingless.cnc.g9f.real-original-companded-development.v1",
        "result": RESULT, "model": MODEL, "revision": REVISION,
        "source_sha": source_sha, "original_shard": SHARD,
        "original_shard_sha256": manifest["files"][SHARD],
        "source_custody": custody, "original_header": header_receipt,
        "frozen_seeds": list(SEEDS), "frozen_expert_groups": list(COUNTS),
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
    print("CNC_G9F_VERIFIED_ORIGINAL_COMPANDED_WINDOW_COMPARISON", len(results))


if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--manifest", type=Path, required=True)
    p.add_argument("--scratch-root", type=Path, required=True)
    p.add_argument("--source-sha", required=True)
    p.add_argument("--output", type=Path, required=True)
    a = p.parse_args()
    run(a.manifest, a.scratch_root, a.output, a.source_sha)
