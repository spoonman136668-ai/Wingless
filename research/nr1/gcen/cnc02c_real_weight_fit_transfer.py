"""CNC-02C: authentic original BF16 Qwen expert-window compressed learning.

Uses the pre-existing G31/G32 SHA256 original-shard custody and G32's frozen
first 4 experts, five seeds and 128x256 windows. Training and independent
transfer activations are synthetic; not Qwen model activations or inference.
"""
import argparse
import json
import os
from pathlib import Path
import shutil
import tempfile
import time
from urllib.request import Request, urlopen

import numpy as np
from cnc02b_fit_transfer import compare_window, ROWS, COLS, STEPS
from cnc02d_matched_controls import compare_controls
from cnc02f_consensus_learning import compare_consensus
from collect_upstream_manifest import build_manifest
from expert_scaling_comparison import SEEDS
from g31_original_shard_stage import (
    MAX_SHARD_BYTES, RESERVE_BYTES, SHARD, URL, check_upstream,
    inspect_header, parse_unique_json, peak_rss_bytes, stage_verified,
)
from g32_real_expert_scaling import read_bf16_window, select_experts
from out_of_core_preflight import pinned_metadata
from runtime_feasibility import MODEL, REVISION, check_manifest

EXPERT_COUNT = 4
SCHEMA = "wingless.cnc02c.real-original-compressed-fit-transfer.v1"


def evaluate_verified_windows(windows, *, seed, comparator=None):
    if type(seed) is not int or seed not in SEEDS:
        raise ValueError("unfrozen G32 seed")
    w = np.asarray(windows, dtype=np.float32)
    if w.shape != (EXPERT_COUNT, ROWS, COLS) or not np.isfinite(w).all():
        raise ValueError("invalid original BF16 windows")
    evaluator = compare_window if comparator is None else comparator
    if evaluator not in (compare_window, compare_controls, compare_consensus):
        raise ValueError("unfrozen experiment comparator")
    results = []
    for idx, weight in enumerate(w):
        outcome = dict(evaluator(weight, steps=STEPS))
        outcome["original_expert_id"] = idx
        outcome["source_window_seed"] = seed
        results.append(outcome)
    return results


def run(manifest_path, scratch_root, output, source_sha, *, comparison_mode="direct"):
    if comparison_mode not in ("direct", "three-arm", "consensus"):
        raise ValueError("unfrozen comparison mode")
    if (os.name != "posix" or type(source_sha) is not str or
            len(source_sha) != 40 or any(c not in "0123456789abcdef" for c in source_sha)):
        raise ValueError("exact source SHA and Linux custody required")
    if not scratch_root.is_dir() or output.exists():
        raise ValueError("missing scratch or existing immutable output")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    if shutil.disk_usage(scratch_root).free < MAX_SHARD_BYTES + RESERVE_BYTES:
        raise RuntimeError("insufficient isolated scratch space")
    check_upstream(manifest, build_manifest())
    index = pinned_metadata("model.safetensors.index.json", manifest)
    baseline_rss = peak_rss_bytes()
    start = time.monotonic()
    with tempfile.TemporaryDirectory(prefix="cnc02c-original-", dir=scratch_root) as td:
        stage = Path(td) / SHARD
        with urlopen(Request(URL, headers={"Accept-Encoding": "identity",
                     "User-Agent": "Wingless-CNC02C-Development/1"}), timeout=90) as response:
            custody = stage_verified(response, stage, manifest["files"][SHARD])
        header_receipt = inspect_header(stage, index)
        with stage.open("rb") as stream:
            header_bytes = int.from_bytes(stream.read(8), "little")
            if header_bytes != header_receipt["header_bytes"]:
                raise ValueError("original header drift")
            header = parse_unique_json(stream.read(header_bytes))
            selected = select_experts(header, index["weight_map"])
            if len(selected["names"]) < EXPERT_COUNT:
                raise ValueError("missing original experts")
            rows = []
            for seed in SEEDS:
                rng = np.random.default_rng(seed)
                row0 = int(rng.integers(0, selected["shape"][0] - ROWS + 1))
                col0 = int(rng.integers(0, selected["shape"][1] - COLS + 1))
                windows = np.stack([
                    read_bf16_window(stream, header_bytes, offset,
                                     selected["shape"], row0, col0)
                    for offset in selected["offsets"][:EXPERT_COUNT]
                ])
                for result in evaluate_verified_windows(
                        windows, seed=seed,
                        comparator=(compare_controls if comparison_mode == "three-arm" else
                                    compare_consensus if comparison_mode == "consensus" else None)):
                    result["original_expert_name"] = selected["names"][result["original_expert_id"]]
                    result["source_window_start"] = [row0, col0]
                    result["original_selected_layer"] = selected["layer"]
                    rows.append(result)
        if len(rows) != EXPERT_COUNT * len(SEEDS):
            raise ValueError("incomplete original-window trial inventory")
        if (custody["observed_peak_staged_file_count"] != 1 or
                list(Path(td).iterdir()) != [stage]):
            raise ValueError("original-shard staging violated")
        peak_rss = peak_rss_bytes()
    receipt = {
        "schema": (SCHEMA if comparison_mode == "direct" else
                   "wingless.cnc02e.real-original-three-arm.v1" if comparison_mode == "three-arm" else
                   "wingless.cnc02g.original-consensus.v1"),
        "result": ("REAL_ORIGINAL_WEIGHTS_COMPRESSED_FIT_TRANSFER_COMPLETE" if comparison_mode == "direct" else
                   "REAL_ORIGINAL_WEIGHTS_THREE_ARM_COMPARISON_COMPLETE" if comparison_mode == "three-arm" else
                   "REAL_ORIGINAL_WEIGHTS_FROZEN_CONSENSUS_COMPLETE"),
        "source_sha": source_sha, "model": MODEL, "revision": REVISION,
        "source_shard": SHARD, "source_shard_sha256": manifest["files"][SHARD],
        "source_custody": custody, "original_header": header_receipt,
        "frozen_source_window_seeds": list(SEEDS),
        "source_window_shape": [ROWS, COLS], "original_experts_per_seed": EXPERT_COUNT,
        "direct_compressed_training_steps_per_expert": STEPS,
        "result_count": len(rows), "comparisons": rows,
        "baseline_rss_bytes": baseline_rss, "peak_rss_bytes": peak_rss,
        "elapsed_seconds": time.monotonic() - start,
        "staged_original_files_after_cleanup": 0,
        "activation_kind": "SYNTHETIC_INDEPENDENT_GAUSSIAN_FIT_TRANSFER",
        "original_model_activations": "NOT_CAPTURED",
        "fp32_baseline_matched": False,
        "equal_update_step_counts_only": comparison_mode in ("three-arm", "consensus"),
        "comparison_mode": comparison_mode,
        "full_model_forward": False, "gpu_execution": False,
        "scientific_acceptance": "NOT_GRANTED", "production_activation": False,
    }
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8") as fp:
        json.dump(receipt, fp, indent=2, sort_keys=True)
        fp.write("\n")
    print("CNC02C_REAL_QWEN_WEIGHT_WINDOWS_ADAPTATION", len(rows))


if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--manifest", required=True, type=Path)
    p.add_argument("--scratch-root", required=True, type=Path)
    p.add_argument("--output", required=True, type=Path)
    p.add_argument("--source-sha", required=True)
    p.add_argument("--comparison-mode", choices=("direct", "three-arm", "consensus"), default="direct")
    a = p.parse_args()
    run(a.manifest, a.scratch_root, a.output, a.source_sha,
        comparison_mode=a.comparison_mode)
