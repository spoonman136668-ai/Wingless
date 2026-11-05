"""CNC-G9M: fail-closed external original-forward capture custody preflight.

Checks local evidence bytes and declared identities only. It cannot verify that
an original Qwen model actually emitted them. Does not parse or upload traces,
download weights, call an inference service, or approve scientific use.
"""
import argparse
import hashlib
import json
from pathlib import Path
import re

MODEL = "Qwen/Qwen3-30B-A3B"
REVISION = "ad44e777bcd18fa416d9da3bd8f70d33ebb85d39"
SOURCE = re.compile(r"[0-9a-f]{40}\Z")
SHA256 = re.compile(r"[0-9a-f]{64}\Z")
TENSOR = re.compile(r"model\.layers\.[0-9]+\.mlp\.experts\.[0-9]+\.(?:gate_proj|up_proj|down_proj)\.weight\Z")
SCHEMA = "wingless.cnc.g9m.native-forward-capture-custody.v1"
MAX_TRACE_BYTES = 1 << 30


def _records(value, name):
    if (not isinstance(value, list) or not value or len(value) > 10000
        or any(type(x) is not str or not 1 <= len(x) <= 128
               or not re.fullmatch(r"[A-Za-z0-9_.:-]+", x) for x in value)
        or len(set(value)) != len(value)):
        raise ValueError("invalid " + name)
    return set(value)


def _safe_trace(root, relpath):
    if type(relpath) is not str or not relpath or Path(relpath).is_absolute():
        raise ValueError("trace must be a relative local file")
    bits = Path(relpath).parts
    if len(bits) > 8 or any(v in ("..", ".", "") for v in bits):
        raise ValueError("trace path escapes capture directory")
    root = root.resolve(strict=True)
    path = root.joinpath(*bits)
    if path.is_symlink() or not path.is_file() or not path.resolve(strict=True).is_relative_to(root):
        raise ValueError("trace path is missing, symlinked or outside root")
    size = path.stat().st_size
    if not 0 < size <= MAX_TRACE_BYTES:
        raise ValueError("trace size exceeds frozen custody limit")
    return path, size


def verify_capture(manifest, capture_dir):
    if not isinstance(manifest, dict) or manifest.get("schema") != SCHEMA:
        raise ValueError("missing frozen capture schema")
    if manifest.get("model") != MODEL or manifest.get("revision") != REVISION:
        raise ValueError("original checkpoint identity mismatch")
    if manifest.get("cohort") != "development" or manifest.get("protected_heldout_used") is not False:
        raise ValueError("protected or nondevelopment records forbidden")
    if manifest.get("capture_kind") != "native_model_forward_hook":
        raise ValueError("native forward-hook capture claim required")
    if type(manifest.get("source_commit")) is not str or not SOURCE.fullmatch(manifest["source_commit"]):
        raise ValueError("source commit identity absent")
    if type(manifest.get("workflow_run_id")) is not int or manifest["workflow_run_id"] <= 0:
        raise ValueError("source workflow identity absent")
    if type(manifest.get("tensor_name")) is not str or not TENSOR.fullmatch(manifest["tensor_name"]):
        raise ValueError("exact original expert source tensor missing")
    calibration = _records(manifest.get("calibration_record_ids"), "calibration cohort")
    evaluation = _records(manifest.get("evaluation_record_ids"), "evaluation cohort")
    if calibration & evaluation:
        raise ValueError("calibration and evaluation record overlap")
    traces = manifest.get("trace_files")
    if not isinstance(traces, list) or len(traces) != 2:
        raise ValueError("exactly two separate trace files required")
    filenames = set()
    digests = {}
    for item in traces:
        if not isinstance(item, dict) or set(item) != {"cohort", "path", "sha256"}:
            raise ValueError("trace entry malformed")
        cohort = item["cohort"]
        if cohort not in ("calibration", "evaluation") or cohort in digests:
            raise ValueError("duplicate/unknown trace cohort")
        name = item["path"]
        if name in filenames:
            raise ValueError("same trace path reused")
        filenames.add(name)
        digest = item["sha256"]
        if type(digest) is not str or not SHA256.fullmatch(digest):
            raise ValueError("trace checksum absent")
        path, size = _safe_trace(capture_dir, name)
        hashed = hashlib.sha256()
        with path.open("rb") as handle:
            while chunk := handle.read(1024*1024):
                hashed.update(chunk)
        if hashed.hexdigest() != digest:
            raise ValueError("trace byte checksum mismatch")
        digests[cohort] = {"sha256": digest, "bytes": size}
    if set(digests) != {"calibration", "evaluation"}:
        raise ValueError("disjoint trace cohorts missing")
    return {
        "schema": "wingless.cnc.g9m.custody-report.v1",
        "source_commit": manifest["source_commit"],
        "source_workflow_run_id": manifest["workflow_run_id"],
        "model": MODEL, "revision": REVISION,
        "tensor_name": manifest["tensor_name"],
        "record_counts": {"calibration": len(calibration), "evaluation": len(evaluation)},
        "trace_file_integrity": digests,
        "local_sha256_verified": True,
        "original_forward_authenticity": "UNVERIFIED_REQUIRES_INDEPENDENT_ATTESTATION",
        "gpu_resource_qualification": "UNVERIFIED",
        "actual_full_model_inference": "UNVERIFIED",
        "source_model_activation_evaluation_admitted": False,
        "scientific_acceptance": "NOT_GRANTED",
        "published_raw_activations": False,
    }


def main():
    parser = argparse.ArgumentParser(description="Offline verify claimed original-forward capture bytes; never admit authenticity")
    parser.add_argument("--manifest", type=Path, required=True)
    parser.add_argument("--capture-dir", type=Path, required=True)
    args = parser.parse_args()
    data = json.loads(args.manifest.read_text(encoding="utf-8"))
    print(json.dumps(verify_capture(data, args.capture_dir), sort_keys=True))


if __name__ == "__main__":
    main()
