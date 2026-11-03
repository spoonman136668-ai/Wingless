"""G31 single original Qwen3-MoE shard: pinned byte custody and physical staging.

One-shot development diagnostic. Never constructs or executes the original model,
never decodes tensors into resident weights and never accesses held-out data.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import resource
import shutil
import struct
import tempfile
import time
from urllib.request import Request, urlopen

from runtime_feasibility import MODEL, REVISION, check_manifest
from out_of_core_preflight import pinned_metadata, SHARDS
from collect_upstream_manifest import build_manifest

SHARD = SHARDS[0]
MAX_SHARD_BYTES = 6 * 1024**3
MIN_SHARD_BYTES = 1 * 1024**3
RESERVE_BYTES = 8 * 1024**3
CHUNK_BYTES = 4 * 1024**2
HEADER_LIMIT = 32 * 1024**2
MAX_RSS_INCREASE = 1024**3
DTYPE_BYTES = {"BF16": 2, "F16": 2, "F32": 4, "I64": 8, "I32": 4, "I16": 2, "I8": 1, "U8": 1}
RESULT = "G31_ORIGINAL_SINGLE_SHARD_BYTES_VERIFIED"
URL = f"https://huggingface.co/{MODEL}/resolve/{REVISION}/{SHARD}"


def parse_unique_json(raw):
    def pairs(items):
        found = {}
        for key, value in items:
            if key in found:
                raise ValueError("duplicate JSON key")
            found[key] = value
        return found
    return json.loads(raw, object_pairs_hook=pairs)


def check_upstream(manifest, upstream):
    check_manifest(manifest)
    check_manifest(upstream)
    if manifest != upstream:
        raise ValueError("upstream LFS and metadata custody mismatch")
    if upstream["files"][SHARD] != manifest["files"][SHARD]:
        raise ValueError("original shard digest differs from pinned source")


def inspect_header(shard_path, index, *, header_limit=HEADER_LIMIT):
    """Read ONLY the Safetensors header; never map or materialize real tensors."""
    weight_map = index.get("weight_map")
    if not isinstance(weight_map, dict) or not weight_map:
        raise ValueError("pinned checkpoint index absent")
    expected = {key for key, part in weight_map.items() if part == SHARD}
    if not expected or set(weight_map.values()) != set(SHARDS):
        raise ValueError("checkpoint index not complete")
    size = shard_path.stat().st_size
    with shard_path.open("rb") as stream:
        first = stream.read(8)
        if len(first) != 8:
            raise ValueError("short Safetensors header length")
        width = struct.unpack("<Q", first)[0]
        if not 0 < width <= header_limit or 8 + width >= size:
            raise ValueError("header size exceeds allowed physical shard")
        header_raw = stream.read(width)
        if len(header_raw) != width:
            raise ValueError("truncated Safetensors header")
    header = parse_unique_json(header_raw)
    if not isinstance(header, dict):
        raise ValueError("invalid Safetensors header")
    keys = set(header) - {"__metadata__"}
    if keys != expected:
        raise ValueError("source/index tensor key mapping mismatch")
    spans = []
    dtypes = {}
    for key in sorted(keys):
        item = header[key]
        if not isinstance(item, dict) or set(item) != {"dtype", "shape", "data_offsets"}:
            raise ValueError("unrecognized tensor header")
        dtype, dims, offsets = item["dtype"], item["shape"], item["data_offsets"]
        if dtype not in DTYPE_BYTES or not isinstance(dims, list) or not isinstance(offsets, list):
            raise ValueError("invalid tensor type or dimensions")
        if not dims or len(offsets) != 2 or any(type(n) is not int or n <= 0 for n in dims):
            raise ValueError("invalid tensor shape")
        if any(type(n) is not int for n in offsets):
            raise ValueError("invalid tensor offset")
        count = 1
        for n in dims:
            count *= n
        lo, hi = offsets
        if not 0 <= lo < hi <= size - 8 - width or hi - lo != count * DTYPE_BYTES[dtype]:
            raise ValueError("tensor payload boundary mismatch")
        spans.append((lo, hi))
        dtypes[dtype] = dtypes.get(dtype, 0) + 1
    cursor = 0
    for lo, hi in sorted(spans):
        if lo != cursor:
            raise ValueError("noncontiguous or overlapping tensor payload")
        cursor = hi
    if cursor != size - 8 - width:
        raise ValueError("Safetensors file length did not reconcile")
    return {"header_bytes": width, "tensor_count": len(keys), "dtype_inventory": dtypes,
            "indexed_payload_bytes": cursor, "header_and_payload_reconciled": True}


def stage_verified(response, dest, expected_sha, *, max_bytes=MAX_SHARD_BYTES,
                   min_bytes=MIN_SHARD_BYTES):
    """Consume a single verified response into one isolated stage file. No retries."""
    if dest.exists() or not isinstance(expected_sha, str) or len(expected_sha) != 64:
        raise ValueError("invalid stage destination or digest")
    if response.status != 200 or response.headers.get("Content-Range"):
        raise ValueError("partial or unexpected upstream response")
    if not str(response.geturl()).startswith("https://"):
        raise ValueError("untrusted non-HTTPS final response")
    encoding = response.headers.get("Content-Encoding")
    if encoding not in (None, "identity"):
        raise ValueError("content encoding not byte-identical")
    declared = response.headers.get("Content-Length")
    if declared is not None:
        if not declared.isdigit() or not min_bytes <= int(declared) <= max_bytes:
            raise ValueError("upstream content length outside original-shard bound")
    hasher = hashlib.sha256()
    total = 0
    staged_count_peak = 0
    rss_peak_kb = resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    allocated_peak = 0
    try:
        with dest.open("xb") as sink:
            staged_count_peak = max(staged_count_peak, len(tuple(dest.parent.iterdir())))
            if staged_count_peak != 1:
                raise RuntimeError("more than one shard was staged")
            while True:
                chunk = response.read(CHUNK_BYTES)
                if not chunk:
                    break
                total += len(chunk)
                if total > max_bytes:
                    raise ValueError("upstream shard exceeded byte cap")
                hasher.update(chunk)
                sink.write(chunk)
                rss_peak_kb = max(rss_peak_kb, resource.getrusage(resource.RUSAGE_SELF).ru_maxrss)
                allocated_peak = max(allocated_peak, dest.stat().st_blocks * 512)
        if not min_bytes <= total <= max_bytes:
            raise ValueError("source shard size invalid")
        if declared is not None and total != int(declared):
            raise ValueError("truncated upstream shard")
        digest = hasher.hexdigest()
        if digest != expected_sha:
            raise ValueError("pinned original shard SHA256 mismatch")
        return {"bytes_verified": total, "observed_sha256": digest,
                "physical_allocated_bytes_peak": max(allocated_peak, dest.stat().st_blocks * 512),
                "observed_peak_staged_file_count": staged_count_peak,
                "transfer_peak_process_rss_bytes": rss_peak_kb * 1024}
    except BaseException:
        dest.unlink(missing_ok=True)
        raise


def run(manifest_path, scratch_root, output, source_sha):
    if os.name != "posix" or len(source_sha) != 40 or any(x not in "0123456789abcdef" for x in source_sha):
        raise ValueError("Linux-only exact-SHA diagnostic")
    if not scratch_root.is_dir() or output.exists():
        raise ValueError("invalid scratch root or existing receipt")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    if shutil.disk_usage(scratch_root).free < MAX_SHARD_BYTES + RESERVE_BYTES:
        raise RuntimeError("insufficient scratch space for bounded original shard")
    upstream = build_manifest()  # independently retrieve pinned official LFS hashes
    check_upstream(manifest, upstream)
    index = pinned_metadata("model.safetensors.index.json", manifest)
    before_rss = resource.getrusage(resource.RUSAGE_SELF).ru_maxrss * 1024
    free_before = shutil.disk_usage(scratch_root).free
    started = time.monotonic()
    with tempfile.TemporaryDirectory(prefix="gcen-g31-one-original-shard-", dir=scratch_root) as temporary:
        root = Path(temporary)
        stage = root / SHARD
        req = Request(URL, headers={"Accept-Encoding": "identity",
                                    "User-Agent": "Wingless-GCEN-G31-Development/1"})
        with urlopen(req, timeout=90) as response:
            metrics = stage_verified(response, stage, manifest["files"][SHARD])
        if list(root.iterdir()) != [stage]:
            raise RuntimeError("unexpected staged original files")
        details = inspect_header(stage, index)
        if metrics["observed_peak_staged_file_count"] != 1 or details["tensor_count"] <= 0:
            raise RuntimeError("original shard staging invariant failed")
        peak_rss = resource.getrusage(resource.RUSAGE_SELF).ru_maxrss * 1024
        if peak_rss - before_rss > MAX_RSS_INCREASE:
            raise RuntimeError("unbounded memory growth during one-shard stage")
        receipt = {
            "schema": "wingless.gcen.g31.original-one-shard-custody.v1",
            "result": RESULT, "source_sha": source_sha,
            "model": MODEL, "revision": REVISION,
            "shard": SHARD, "expected_sha256": manifest["files"][SHARD],
            "upstream_lfs_digest_matched_frozen_manifest": True,
            "physical_scratch_free_before_bytes": free_before,
            "single_staged_shard": metrics, "header": details,
            "baseline_process_rss_bytes": before_rss,
            "peak_process_rss_bytes": peak_rss,
            "elapsed_seconds": time.monotonic() - started,
            "original_weight_shards_verified_this_gate": 1,
            "full_checkpoint_custody": "UNVERIFIED",
            "entire_model_loaded": False, "original_model_forward_executed": False,
            "authentic_activation_trace": False,
            "original_model_compression_accuracy": "UNVERIFIED",
            "scientific_acceptance": "NOT_GRANTED",
            "production_activation": False
        }
    # The entire temporary directory (including original weights) is now gone.
    receipt["original_staging_files_after_cleanup"] = 0
    receipt["physical_scratch_free_after_bytes"] = shutil.disk_usage(scratch_root).free
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8") as writer:
        json.dump(receipt, writer, indent=2, sort_keys=True)
        writer.write("\n")
    print("GCEN_G31_ORIGINAL_SINGLE_SHARD_CUSTODY_PASS")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", required=True, type=Path)
    parser.add_argument("--scratch-root", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--source-sha", required=True)
    args = parser.parse_args()
    run(args.manifest, args.scratch_root, args.output, args.source_sha)
