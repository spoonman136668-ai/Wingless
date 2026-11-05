"""One-time Wingless original Qwen shard staging, without model execution.

Reuse G31's bounded HTTPS stream and byte-authentication contract. Never
download other shards, launch a model, execute GPU kernels, or publish weights.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import shutil
import tempfile
from urllib.request import Request, urlopen

from g31_original_shard_stage import URL, SHARD, RESERVE_BYTES, stage_verified
from g9u_verified_loader import PINNED_SHARD_SHA256
from runtime_feasibility import MODEL, REVISION, check_manifest

EXPECTED_BYTES = 3999417504
DIRNAME = "wingless-g9u-original"
SCHEMA = "wingless.gcen.g9u.persistent-shard-stage.v1"


def linked(path):
    return path.is_symlink() or (hasattr(path, "is_junction") and path.is_junction())


def verify(path):
    if linked(path) or not path.is_file() or path.stat().st_size != EXPECTED_BYTES:
        raise ValueError("original shard missing, linked or wrong byte count")
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(4 * 1024**2), b""):
            digest.update(chunk)
    if digest.hexdigest() != PINNED_SHARD_SHA256:
        raise ValueError("original shard SHA256 mismatch")


def stage(*, toolcache_root, manifest_path, output, source_sha):
    if not isinstance(source_sha, str) or len(source_sha) != 40 or any(
            x not in "0123456789abcdef" for x in source_sha):
        raise ValueError("invalid exact source SHA")
    root, manifest_path, output = Path(toolcache_root), Path(manifest_path), Path(output)
    if linked(root) or not root.is_dir() or linked(manifest_path):
        raise ValueError("untrusted root or manifest")
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    if (manifest["model"] != MODEL or manifest["revision"] != REVISION or
            manifest["files"][SHARD] != PINNED_SHARD_SHA256):
        raise ValueError("frozen checkpoint identity mismatch")
    if linked(output) or output.exists() or not output.parent.is_dir():
        raise ValueError("refusing to overwrite stage evidence")
    destination = root / DIRNAME
    if destination.exists():
        if linked(destination) or not destination.is_dir():
            raise ValueError("untrusted existing stage directory")
    else:
        destination.mkdir(mode=0o700)
    target = destination / SHARD
    if linked(target):
        raise ValueError("linked target forbidden")
    if target.exists():
        verify(target)
        status = "REUSED_VERIFIED_ORIGINAL_SHARD"
    else:
        if shutil.disk_usage(destination).free < EXPECTED_BYTES + RESERVE_BYTES:
            raise RuntimeError("insufficient original shard storage and reserve")
        with tempfile.TemporaryDirectory(prefix="g9u-once-", dir=destination) as temp:
            partial = Path(temp) / SHARD
            request = Request(URL, headers={
                "Accept-Encoding": "identity",
                "User-Agent": "Wingless-G9U-PinnedOriginalShard/1"})
            # Only one request; no resume, retry, or unverified source fallback.
            with urlopen(request, timeout=90) as response:
                custody = stage_verified(response, partial, PINNED_SHARD_SHA256,
                                         min_bytes=EXPECTED_BYTES,
                                         max_bytes=EXPECTED_BYTES)
            if custody["bytes_verified"] != EXPECTED_BYTES or target.exists():
                raise RuntimeError("shard changed during authenticated stage")
            os.rename(partial, target)
        verify(target)
        status = "DOWNLOADED_VERIFIED_ORIGINAL_SHARD"
    receipt = {
        "schema": SCHEMA, "status": status,
        "source_repo": "spoonman136668-ai/Wingless", "exact_source_sha": source_sha,
        "model": MODEL, "revision": REVISION, "shard": SHARD,
        "shard_sha256": PINNED_SHARD_SHA256, "shard_bytes": EXPECTED_BYTES,
        "stage_identity": DIRNAME + "/" + SHARD,
        "all_shards_downloaded": False, "original_expert_loaded": False,
        "gpu_kernel_executed": False, "full_model_forward": False,
        "scientific_acceptance": "NOT_GRANTED",
    }
    with output.open("x", encoding="utf-8") as sink:
        json.dump(receipt, sink, sort_keys=True, indent=2)
        sink.write("\n")
    print("G9U_SHARD_STAGE_PASS " + status, flush=True)
    return receipt


if __name__ == "__main__":
    p = argparse.ArgumentParser(description="Stage pinned original shard exactly once")
    p.add_argument("--toolcache-root", type=Path, required=True)
    p.add_argument("--manifest", type=Path, required=True)
    p.add_argument("--output", type=Path, required=True)
    p.add_argument("--source-sha", required=True)
    a = p.parse_args()
    stage(toolcache_root=a.toolcache_root, manifest_path=a.manifest,
          output=a.output, source_sha=a.source_sha)
