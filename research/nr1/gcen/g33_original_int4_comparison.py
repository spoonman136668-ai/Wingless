"""CNC-01Q-G8: original Qwen BF16 weight-window INT4 comparison.

Development-only. Reuses G31 original-shard SHA-256 custody and G32
frozen layer/window sampling. Native model inference is NOT performed.
The 32x256 activations here are SYNTHETIC probes, NOT captured Qwen activations.
This code never loads a complete model or consumes protected held-out prompts.
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
from direct_expert import apply_direct, factorize
from expert_scaling_comparison import COUNTS, COLS, ROWS, SEEDS, SVD_RANK
from g31_original_shard_stage import (
    MAX_SHARD_BYTES, RESERVE_BYTES, SHARD, URL, check_upstream,
    inspect_header, parse_unique_json, peak_rss_bytes, stage_verified,
)
from g32_real_expert_scaling import read_bf16_window, select_experts
from out_of_core_preflight import pinned_metadata
from packed_int4 import apply_packed_int4, pack_int4, stored_bytes
from runtime_feasibility import MODEL, REVISION, check_manifest

SOURCE_SHA = re.compile(r"[0-9a-f]{40}\Z")
BATCH = 32
BLOCK = 64
RESULT = "G8_ORIGINAL_BF16_WINDOW_INT4_DEVELOPMENT_COMPLETED"


def compare_expert(weight, activations):
    """Compare one authentic-window-compatible expert under matched inputs.

    Weight comparison uses the operator on an identity matrix, not an expanded
    dequantized full-weight allocation. SVD is the preexisting rank-14 control.
    Storage is estimated FP16 factor serialization vs packed bytes, not RSS.
    """
    w = np.asarray(weight, dtype=np.float32)
    x = np.asarray(activations, dtype=np.float32)
    if (w.shape != (ROWS, COLS) or x.shape != (BATCH, COLS) or
            not np.isfinite(w).all() or not np.isfinite(x).all()):
        raise ValueError("invalid frozen expert and diagnostic input shape")
    if not np.any(w) or not np.any(x):
        raise ValueError("zero-energy sample forbidden")
    packed, scales = pack_int4(w, block=BLOCK)
    x_ref = x @ w.T
    quant = apply_packed_int4(x, packed, scales, block=BLOCK)
    left, right = factorize(w, SVD_RANK)
    svd = apply_direct(x, left, right)
    original_weight_energy = float(np.sum(w.astype(np.float64)**2))
    ref_energy = float(np.sum(x_ref.astype(np.float64)**2))
    if ref_energy <= 0:
        raise ValueError("zero-energy original output")
    # Identity probe measures the full matrix error without producing a
    # dequantized original weight tensor for the inference operator.
    identity = np.eye(COLS, dtype=np.float32)
    reconstructed_by_operator = apply_packed_int4(identity, packed, scales, block=BLOCK).T
    delta_weight = reconstructed_by_operator.astype(np.float64) - w.astype(np.float64)
    delta_int4 = quant.astype(np.float64) - x_ref.astype(np.float64)
    delta_svd = svd.astype(np.float64) - x_ref.astype(np.float64)
    return {
        "original_weight_energy": original_weight_energy,
        "original_output_energy": ref_energy,
        "int4_weight_squared_error": float(np.sum(delta_weight**2)),
        "int4_output_squared_error": float(np.sum(delta_int4**2)),
        "svd_output_squared_error": float(np.sum(delta_svd**2)),
        "original_bf16_bytes": int(w.size * 2),
        "packed_int4_bytes_including_metadata": stored_bytes(packed, scales, COLS, BLOCK),
        "svd_rank14_fp16_factor_bytes": int(SVD_RANK * (ROWS + COLS + 1) * 2),
    }


def compare_windows(windows, seed):
    if type(seed) is not int or seed not in SEEDS:
        raise ValueError("unfrozen development seed")
    w = np.asarray(windows, dtype=np.float32)
    if (w.shape != (max(COUNTS), ROWS, COLS) or not np.isfinite(w).all()):
        raise ValueError("invalid frozen 16-expert window inventory")
    # Independent synthetic input fixture, used only as an algebraic probe.
    x = np.random.default_rng(26000 + seed).standard_normal(
        (BATCH, COLS)
    ).astype(np.float32)
    rows = [compare_expert(weight, x) for weight in w]
    grouped = []
    for count in COUNTS:
        subset = rows[:count]
        sums = lambda field: float(sum(v[field] for v in subset))
        weight_den = sums("original_weight_energy")
        output_den = sums("original_output_energy")
        grouped.append({
            "experts": count, "seed": seed, "sample_input_kind": "SYNTHETIC_GAUSSIAN_NOT_NATIVE_ACTIVATIONS",
            "original_bf16_bytes": int(sums("original_bf16_bytes")),
            "packed_int4_bytes_including_metadata": int(sums("packed_int4_bytes_including_metadata")),
            "svd_rank14_fp16_factor_bytes": int(sums("svd_rank14_fp16_factor_bytes")),
            "int4_relative_weight_squared_error": sums("int4_weight_squared_error") / weight_den,
            "int4_relative_output_squared_error": sums("int4_output_squared_error") / output_den,
            "svd_rank14_relative_output_squared_error": sums("svd_output_squared_error") / output_den,
            "per_expert": [
                {
                    "expert_id": i,
                    "int4_relative_weight_squared_error": row["int4_weight_squared_error"] / row["original_weight_energy"],
                    "int4_relative_output_squared_error": row["int4_output_squared_error"] / row["original_output_energy"],
                    "svd_rank14_relative_output_squared_error": row["svd_output_squared_error"] / row["original_output_energy"],
                }
                for i, row in enumerate(subset)
            ],
        })
    return grouped


def run(manifest_path, scratch_root, output, source_sha):
    if os.name != "posix" or type(source_sha) is not str or not SOURCE_SHA.fullmatch(source_sha):
        raise ValueError("G8 requires exact SHA and Linux custody resource probes")
    if not scratch_root.is_dir() or output.exists():
        raise ValueError("missing scratch root or existing output")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    if shutil.disk_usage(scratch_root).free < MAX_SHARD_BYTES + RESERVE_BYTES:
        raise RuntimeError("insufficient original shard staging capacity")
    upstream = build_manifest()
    check_upstream(manifest, upstream)
    index = pinned_metadata("model.safetensors.index.json", manifest)
    start = time.monotonic()
    baseline_rss = peak_rss_bytes()
    with tempfile.TemporaryDirectory(prefix="cnc-g8-original-", dir=scratch_root) as td:
        stage = Path(td) / SHARD
        request = Request(URL, headers={"Accept-Encoding": "identity",
                                        "User-Agent": "Wingless-CNC-G8-Development/1"})
        with urlopen(request, timeout=90) as response:
            custody = stage_verified(response, stage, manifest["files"][SHARD])
        header_receipt = inspect_header(stage, index)
        with stage.open("rb") as stream:
            width = int.from_bytes(stream.read(8), "little")
            if width != header_receipt["header_bytes"]:
                raise ValueError("post-custody header drift")
            header = parse_unique_json(stream.read(width))
            selected = select_experts(header, index["weight_map"])
            comparisons = []
            for seed in SEEDS:
                rng = np.random.default_rng(seed)
                row0 = int(rng.integers(0, selected["shape"][0] - ROWS + 1))
                col0 = int(rng.integers(0, selected["shape"][1] - COLS + 1))
                windows = np.stack([
                    read_bf16_window(stream, width, offset, selected["shape"], row0, col0)
                    for offset in selected["offsets"]
                ])
                for entry in compare_windows(windows, seed):
                    entry.update(
                        window_start=[row0, col0],
                        selected_layer=selected["layer"],
                        source_kind="ORIGINAL_SHA256_VERIFIED_BF16_WEIGHTS",
                        original_expert_names=selected["names"][:entry["experts"]],
                    )
                    comparisons.append(entry)
        if len(comparisons) != len(SEEDS) * len(COUNTS):
            raise ValueError("incomplete frozen comparison")
        if custody["observed_peak_staged_file_count"] != 1 or list(Path(td).iterdir()) != [stage]:
            raise ValueError("staging custody invariant violated")
        rss = peak_rss_bytes()
    receipt = {
        "schema": "wingless.cnc.g8.original-int4-development.v1",
        "status": RESULT, "model": MODEL, "revision": REVISION,
        "source_sha": source_sha, "source_shard": SHARD,
        "source_shard_sha256": manifest["files"][SHARD],
        "source_custody": custody, "original_header": header_receipt,
        "frozen_fit_seeds": list(SEEDS),
        "frozen_expert_counts": list(COUNTS),
        "frozen_window_shape": [ROWS, COLS],
        "int4_block_size": BLOCK, "int4_signed_range": [-7, 7],
        "svd_control_rank": SVD_RANK,
        "synthetic_activation_batch": BATCH,
        "comparison_count": len(comparisons),
        "comparisons": comparisons,
        "baseline_peak_rss_bytes": baseline_rss,
        "observed_peak_rss_bytes": rss,
        "elapsed_seconds": time.monotonic() - start,
        "staged_original_files_after_cleanup": 0,
        "authentic_model_activations": "NOT_CAPTURED",
        "full_model_inference": False,
        "gpu_measurement": False,
        "protected_heldout_used": False,
        "scientific_acceptance": "NOT_GRANTED",
        "production_activation": False,
    }
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8") as stream:
        json.dump(receipt, stream, indent=2, sort_keys=True)
        stream.write("\n")
    print("CNC_G8_VERIFIED_REAL_EXPERT_WEIGHT_WINDOWS", len(comparisons))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", required=True, type=Path)
    parser.add_argument("--scratch-root", required=True, type=Path)
    parser.add_argument("--source-sha", required=True)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    run(args.manifest, args.scratch_root, args.output, args.source_sha)
