"""G32: frozen development-only real BF16 expert-window sharing comparison.

Reuses G31 authenticated single-shard custody, G31 header-only validation and
the existing 4/8/16 expert independent-SVD versus shared-basis comparator.
No full checkpoint, original model forward, task-quality, or protected cohort.
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
from collect_upstream_manifest import build_manifest
from expert_scaling_comparison import COUNTS, COLS, ROWS, SEEDS, evaluate_group
from g31_original_shard_stage import (
    MAX_SHARD_BYTES, RESERVE_BYTES, SHARD, URL, check_upstream,
    inspect_header, parse_unique_json, peak_rss_bytes, stage_verified,
)
from out_of_core_preflight import pinned_metadata
from runtime_feasibility import MODEL, REVISION, check_manifest

EXPERT_IDS = tuple(range(max(COUNTS)))
PROJECTED_WEIGHT = "down_proj.weight"
RESULT = "G32_REAL_BF16_EXPERT_WINDOW_SCALING_SCREEN_COMPLETE"


def select_experts(header, weight_map):
    """Metadata-only frozen first-complete-layer selection; no error optimization."""
    if not isinstance(header, dict) or not isinstance(weight_map, dict):
        raise ValueError("missing original weight metadata")
    for layer in range(48):
        names = [f"model.layers.{layer}.mlp.experts.{expert}.{PROJECTED_WEIGHT}"
                 for expert in EXPERT_IDS]
        if not all(weight_map.get(name) == SHARD and name in header for name in names):
            continue
        descriptions = [header[name] for name in names]
        first = descriptions[0]
        if not isinstance(first, dict):
            raise ValueError("invalid original expert header")
        shape = first.get("shape")
        if first.get("dtype") != "BF16" or not isinstance(shape, list) or len(shape) != 2:
            raise ValueError("original expert dtype/shape incompatible with frozen test")
        if any(type(x) is not int or x <= 0 for x in shape) or shape[0] < ROWS or shape[1] < COLS:
            raise ValueError("original projection smaller than frozen window")
        for desc in descriptions:
            if not isinstance(desc, dict) or desc.get("dtype") != "BF16" or desc.get("shape") != shape:
                raise ValueError("distinct experts must share original BF16 shape")
            offsets = desc.get("data_offsets")
            if not isinstance(offsets, list) or len(offsets) != 2 or any(type(n) is not int for n in offsets):
                raise ValueError("original expert tensor offsets missing")
            if offsets[0] < 0 or offsets[1] - offsets[0] != 2 * shape[0] * shape[1]:
                raise ValueError("original expert physical extent mismatch")
        return {"layer": layer, "names": names, "shape": shape,
                "offsets": [desc["data_offsets"][0] for desc in descriptions]}
    raise ValueError("no complete original first-16-expert layer in staged shard")


def read_bf16_window(stream, header_bytes, start, shape, row0, col0,
                     *, rows=ROWS, cols=COLS):
    """Read only the requested 2D BF16 window from a SHA-verified file."""
    height, width = shape
    if any(type(n) is not int for n in (header_bytes, start, row0, col0, rows, cols)):
        raise ValueError("noninteger expert window bounds")
    if header_bytes <= 0 or start < 0 or rows <= 0 or cols <= 0:
        raise ValueError("invalid window size")
    if row0 < 0 or col0 < 0 or row0 + rows > height or col0 + cols > width:
        raise ValueError("expert window outside verified tensor")
    data = np.empty((rows, cols), dtype=np.float32)
    for i in range(rows):
        stream.seek(8 + header_bytes + start + 2 * ((row0 + i) * width + col0))
        raw = stream.read(cols * 2)
        if len(raw) != cols * 2:
            raise ValueError("short authentic BF16 read")
        data[i] = (np.frombuffer(raw, dtype="<u2").astype(np.uint32) << 16).view(np.float32)
    if not np.isfinite(data).all():
        raise ValueError("nonfinite original BF16 expert weights")
    return data


def run(manifest_path, scratch_root, output, source_sha):
    if os.name != "posix" or len(source_sha) != 40 or any(c not in "0123456789abcdef" for c in source_sha):
        raise ValueError("G32 requires exact SHA and Linux resource probes")
    if not scratch_root.is_dir() or output.exists():
        raise ValueError("missing scratch root or receipt already exists")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    if shutil.disk_usage(scratch_root).free < MAX_SHARD_BYTES + RESERVE_BYTES:
        raise RuntimeError("insufficient isolated scratch for one original shard")
    upstream = build_manifest()
    check_upstream(manifest, upstream)
    index = pinned_metadata("model.safetensors.index.json", manifest)
    baseline_rss = peak_rss_bytes()
    started = time.monotonic()
    with tempfile.TemporaryDirectory(prefix="gcen-g32-verified-original-", dir=scratch_root) as td:
        stage = Path(td) / SHARD
        request = Request(URL, headers={"Accept-Encoding": "identity",
                                        "User-Agent": "Wingless-GCEN-G32-Development/1"})
        with urlopen(request, timeout=90) as response:
            custody = stage_verified(response, stage, manifest["files"][SHARD])
        verified_header = inspect_header(stage, index)
        with stage.open("rb") as stream:
            width = int.from_bytes(stream.read(8), "little")
            if width != verified_header["header_bytes"]:
                raise ValueError("post-custody header changed")
            header = parse_unique_json(stream.read(width))
            selected = select_experts(header, index["weight_map"])
            all_results = []
            for seed in SEEDS:
                rng = np.random.default_rng(seed)
                row0 = int(rng.integers(0, selected["shape"][0] - ROWS + 1))
                col0 = int(rng.integers(0, selected["shape"][1] - COLS + 1))
                windows = np.stack([
                    read_bf16_window(stream, width, start, selected["shape"], row0, col0)
                    for start in selected["offsets"]
                ])
                if windows.shape != (max(COUNTS), ROWS, COLS) or not np.isfinite(windows).all():
                    raise ValueError("failed frozen original expert sampling")
                for experts in COUNTS:
                    entry = evaluate_group(windows[:experts])
                    entry.update(seed=seed, window_start=[row0, col0],
                                 distinct_expert_names=selected["names"][:experts],
                                 source_type="AUTHENTIC_PINNED_BF16_ORIGINAL_WEIGHTS")
                    all_results.append(entry)
        if len(all_results) != len(SEEDS) * len(COUNTS):
            raise RuntimeError("frozen development result inventory incomplete")
        if custody["observed_peak_staged_file_count"] != 1 or list(Path(td).iterdir()) != [stage]:
            raise RuntimeError("one original shard staging violated")
        original_peak_rss = peak_rss_bytes()
        receipt = {
            "schema": "wingless.gcen.g32.real-expert-sharing-development.v1",
            "result": RESULT, "source_sha": source_sha,
            "model": MODEL, "revision": REVISION, "original_shard": SHARD,
            "frozen_original_shard_sha256": manifest["files"][SHARD],
            "source_custody": custody, "header": verified_header,
            "selection_policy": "FIRST_LAYER_WITH_16_DISTINCT_DOWN_PROJECTIONS_ONLY",
            "selected_layer": selected["layer"],
            "source_projection_shape": selected["shape"],
            "frozen_expert_counts": list(COUNTS), "frozen_window_shape": [ROWS, COLS],
            "frozen_fit_window_seeds": list(SEEDS),
            "development_comparisons": all_results,
            "expected_comparison_count": len(SEEDS) * len(COUNTS),
            "baseline_process_rss_bytes": baseline_rss,
            "observed_peak_process_rss_bytes": original_peak_rss,
            "elapsed_seconds": time.monotonic() - started,
            "original_model_forward_executed": False,
            "inference_output_quality": "UNVERIFIED",
            "full_model_compression": "UNVERIFIED",
            "protected_heldout_used": False,
            "scientific_acceptance": "NOT_GRANTED",
            "production_activation": False,
        }
    receipt["staged_original_files_after_cleanup"] = 0
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8") as f:
        json.dump(receipt, f, indent=2, sort_keys=True)
        f.write("\n")
    print("GCEN_G32_VERIFIED_REAL_EXPERT_WINDOWS", len(all_results))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", required=True, type=Path)
    parser.add_argument("--scratch-root", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--source-sha", required=True)
    args = parser.parse_args()
    run(args.manifest, args.scratch_root, args.output, args.source_sha)
