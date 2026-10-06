#!/usr/bin/env python3
import argparse
import hashlib
import json
import math
import os
import pathlib
import re
import resource
import subprocess
import sys
import time

SCHEMA = "research.free-fanout.v1"
ENVELOPE_SCHEMA = "research.free-fanout-envelope.v1"
SUMMARY_SCHEMA = "research.free-fanout-summary.v1"
ID_RE = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$")
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
BOUND_PACKAGE_SHA = "BOUND_AT_DISPATCH"
ALLOWED_DATA = {"historical", "replay", "historical-replay"}
MAX_CANDIDATES = 32
MAX_TIMEOUT_SECONDS = 1800
RESOURCE_FIELDS = (
    "max_parameters",
    "max_context_bytes",
    "max_model_calls",
    "max_peak_rss_kib",
    "max_elapsed_seconds",
)

PROGRAM_CONTRACTS = {
    "Wingless": {
        "signature": "state_former(raw_input, history) -> fixed_size_state",
        "required_false": (
            "sealed_outcome_exposure",
            "post_result_tuning",
            "capacity_growth",
        ),
        "required_positive_int": (
            "fixed_state_dimension",
            "fixed_readout_capacity",
        ),
    },
    "Yggdrasil": {
        "signature": "cognition_consumer(retained_state, local_state) -> decision_state",
        "required_false": (
            "sealed_outcome_exposure",
            "post_result_tuning",
            "hidden_persistent_memory_growth",
            "addressing_change",
            "capacity_growth",
        ),
        "fixed_values": {
            "total_slots": 16,
            "active_slots": 7,
            "retained_slots": 9,
        },
    },
}

class FanoutError(RuntimeError):
    pass

def fail(code):
    raise FanoutError(code)

def load(path):
    try:
        with open(path, "r", encoding="utf-8") as f:
            return json.load(f)
    except (OSError, json.JSONDecodeError) as exc:
        raise FanoutError(f"FREE_FANOUT_JSON_LOAD_FAILED:{type(exc).__name__}") from exc

def canonical_json_sha256(value):
    raw = json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=True).encode("utf-8")
    return hashlib.sha256(raw).hexdigest()

def validate_repo_rel_path(value, *, suffix=None):
    if not isinstance(value, str) or not value or "\\" in value:
        fail("FREE_FANOUT_PATH_INVALID")
    raw_parts = value.split("/")
    if any(part in ("", ".", "..") for part in raw_parts):
        fail("FREE_FANOUT_PATH_INVALID")
    p = pathlib.PurePosixPath(value)
    if p.is_absolute() or len(p.parts) < 3 or p.parts[:2] != ("research", "fanout"):
        fail("FREE_FANOUT_PATH_INVALID")
    if suffix is not None and not value.endswith(suffix):
        fail("FREE_FANOUT_PATH_INVALID")
    return value

def validate_manifest_path(value):
    return validate_repo_rel_path(value, suffix=".json")

def _nonnegative_int(value):
    return isinstance(value, int) and not isinstance(value, bool) and value >= 0

def _positive_int(value):
    return isinstance(value, int) and not isinstance(value, bool) and value > 0

def validate_manifest(m, program, package_sha):
    if not isinstance(m, dict):
        fail("FREE_FANOUT_MANIFEST_NOT_OBJECT")
    if m.get("schema") != SCHEMA:
        fail("FREE_FANOUT_MANIFEST_SCHEMA")
    if program not in PROGRAM_CONTRACTS or m.get("program") != program:
        fail("FREE_FANOUT_MANIFEST_PROGRAM")
    if not isinstance(package_sha, str) or not SHA_RE.fullmatch(package_sha):
        fail("FREE_FANOUT_PACKAGE_SHA_FORMAT")
    declared_package_sha = m.get("package_sha")
    if declared_package_sha not in (package_sha, BOUND_PACKAGE_SHA):
        fail("FREE_FANOUT_PACKAGE_SHA_MISMATCH")
    if declared_package_sha == BOUND_PACKAGE_SHA:
        m = dict(m)
        m["package_sha"] = package_sha
    if m.get("authority") != "external-evidence-only":
        fail("FREE_FANOUT_AUTHORITY")
    if m.get("data_class") not in ALLOWED_DATA:
        fail("FREE_FANOUT_DATA_CLASS")

    for key in (
        "sealed_inputs",
        "accepted_state_mutation",
        "successor_dispatch",
        "queue_write",
        "production_access",
        "promotion_authority",
    ):
        if m.get(key) is not False:
            fail(f"FREE_FANOUT_FORBIDDEN_{key.upper()}")

    entrypoint = m.get("entrypoint")
    validate_repo_rel_path(entrypoint)

    timeout = m.get("candidate_timeout_seconds")
    if not _positive_int(timeout) or timeout > MAX_TIMEOUT_SECONDS:
        fail("FREE_FANOUT_TIMEOUT")

    interface = m.get("candidate_interface")
    if not isinstance(interface, dict):
        fail("FREE_FANOUT_INTERFACE")
    contract = PROGRAM_CONTRACTS[program]
    if interface.get("signature") != contract["signature"]:
        fail("FREE_FANOUT_INTERFACE_SIGNATURE")
    if interface.get("training_data_class") != m["data_class"]:
        fail("FREE_FANOUT_INTERFACE_DATA_CLASS")
    for key in contract.get("required_false", ()):
        if interface.get(key) is not False:
            fail(f"FREE_FANOUT_INTERFACE_{key.upper()}")
    for key in contract.get("required_positive_int", ()):
        if not _positive_int(interface.get(key)):
            fail(f"FREE_FANOUT_INTERFACE_{key.upper()}")
    for key, expected in contract.get("fixed_values", {}).items():
        if interface.get(key) != expected:
            fail(f"FREE_FANOUT_INTERFACE_{key.upper()}")

    envelope = m.get("resource_envelope")
    if not isinstance(envelope, dict):
        fail("FREE_FANOUT_RESOURCE_ENVELOPE")
    if envelope.get("comparison") != "equal-or-lower-than-baseline":
        fail("FREE_FANOUT_RESOURCE_COMPARISON")
    if not isinstance(envelope.get("baseline_id"), str) or not envelope["baseline_id"].strip():
        fail("FREE_FANOUT_RESOURCE_BASELINE")
    for key in ("max_parameters", "max_context_bytes", "max_model_calls"):
        if not _nonnegative_int(envelope.get(key)):
            fail(f"FREE_FANOUT_RESOURCE_{key.upper()}")
    for key in ("max_peak_rss_kib", "max_elapsed_seconds"):
        if not _positive_int(envelope.get(key)):
            fail(f"FREE_FANOUT_RESOURCE_{key.upper()}")
    if envelope["max_elapsed_seconds"] > timeout:
        fail("FREE_FANOUT_RESOURCE_ELAPSED_EXCEEDS_TIMEOUT")

    candidates = m.get("candidates")
    if not isinstance(candidates, list) or not (1 <= len(candidates) <= MAX_CANDIDATES):
        fail("FREE_FANOUT_CANDIDATE_COUNT")
    seen = set()
    for candidate in candidates:
        if not isinstance(candidate, dict):
            fail("FREE_FANOUT_CANDIDATE")
        cid = candidate.get("id")
        if not isinstance(cid, str) or not ID_RE.fullmatch(cid):
            fail("FREE_FANOUT_CANDIDATE_ID")
        if cid in seen:
            fail("FREE_FANOUT_DUPLICATE_CANDIDATE_ID")
        seen.add(cid)
        config = candidate.get("config", {})
        if not isinstance(config, dict):
            fail("FREE_FANOUT_CANDIDATE_CONFIG")
        try:
            json.dumps(config, sort_keys=True, separators=(",", ":"))
        except (TypeError, ValueError) as exc:
            raise FanoutError("FREE_FANOUT_CANDIDATE_CONFIG_JSON") from exc
    return m

def load_validated_manifest(path, program, package_sha):
    return validate_manifest(load(path), program, package_sha)

def resolve_package_entrypoint(package_root, rel_path):
    validate_repo_rel_path(rel_path)
    root = pathlib.Path(package_root).resolve(strict=True)
    allowed = (root / "research" / "fanout").resolve(strict=True)
    target = (root / pathlib.PurePosixPath(rel_path)).resolve(strict=True)
    if not target.is_file() or not target.is_relative_to(allowed):
        fail("FREE_FANOUT_ENTRYPOINT_MISSING_OR_ESCAPED")
    return target

def validate_result(result, candidate_id, manifest):
    if not isinstance(result, dict):
        fail("FREE_FANOUT_RESULT_NOT_OBJECT")
    if result.get("candidate_id") != candidate_id:
        fail("FREE_FANOUT_RESULT_ID_MISMATCH")
    interface = manifest["candidate_interface"]
    if result.get("interface") != interface["signature"]:
        fail("FREE_FANOUT_RESULT_INTERFACE_MISMATCH")

    metric = result.get("metric")
    if (
        not isinstance(metric, dict)
        or not isinstance(metric.get("name"), str)
        or not metric["name"]
        or not isinstance(metric.get("value"), (int, float))
        or isinstance(metric.get("value"), bool)
        or not math.isfinite(float(metric["value"]))
    ):
        fail("FREE_FANOUT_RESULT_METRIC")

    usage = result.get("resource_usage")
    if not isinstance(usage, dict):
        fail("FREE_FANOUT_RESULT_RESOURCE_USAGE")
    limits = manifest["resource_envelope"]
    for result_key, limit_key in (
        ("parameters", "max_parameters"),
        ("context_bytes", "max_context_bytes"),
        ("model_calls", "max_model_calls"),
    ):
        value = usage.get(result_key)
        if not _nonnegative_int(value):
            fail(f"FREE_FANOUT_RESULT_RESOURCE_{result_key.upper()}")
        if value > limits[limit_key]:
            fail(f"FREE_FANOUT_RESOURCE_LIMIT_{result_key.upper()}")

    if manifest["program"] == "Wingless":
        if result.get("fixed_state_dimension") != interface["fixed_state_dimension"]:
            fail("FREE_FANOUT_RESULT_STATE_DIMENSION")
        if result.get("fixed_readout_capacity") != interface["fixed_readout_capacity"]:
            fail("FREE_FANOUT_RESULT_READOUT_CAPACITY")
    elif manifest["program"] == "Yggdrasil":
        for key in ("total_slots", "active_slots", "retained_slots"):
            if result.get(key) != interface[key]:
                fail(f"FREE_FANOUT_RESULT_{key.upper()}")
        if result.get("hidden_persistent_memory_growth") is not False:
            fail("FREE_FANOUT_RESULT_HIDDEN_MEMORY_GROWTH")
    return result

def candidate_by_id(manifest, candidate_id):
    for candidate in manifest["candidates"]:
        if candidate["id"] == candidate_id:
            return candidate
    fail("FREE_FANOUT_CANDIDATE_NOT_FOUND")

def build_candidate_env(home_dir, tmp_dir, data_class):
    home_dir.mkdir(parents=True, exist_ok=True)
    tmp_dir.mkdir(parents=True, exist_ok=True)
    return {
        "PATH": os.environ.get("PATH", "/usr/local/bin:/usr/bin:/bin"),
        "HOME": str(home_dir),
        "TMPDIR": str(tmp_dir),
        "LANG": "C.UTF-8",
        "LC_ALL": "C.UTF-8",
        "PYTHONHASHSEED": "0",
        "FREE_FANOUT_AUTHORITY": "external-evidence-only",
        "FREE_FANOUT_DATA_CLASS": data_class,
    }

def run_candidate(manifest, package_root, candidate_id, output_root):
    candidate = candidate_by_id(manifest, candidate_id)
    output_root = pathlib.Path(output_root).resolve()
    output_root.mkdir(parents=True, exist_ok=True)
    result_path = output_root / "result.json"
    config_path = output_root / "candidate-config.json"
    config_path.write_text(json.dumps(candidate.get("config", {}), indent=2, sort_keys=True) + "\n", encoding="utf-8")

    entrypoint = resolve_package_entrypoint(package_root, manifest["entrypoint"])
    if entrypoint.suffix == ".py":
        command = [sys.executable, str(entrypoint)]
    elif entrypoint.suffix == ".sh":
        command = ["bash", str(entrypoint)]
    elif os.access(entrypoint, os.X_OK):
        command = [str(entrypoint)]
    else:
        fail("FREE_FANOUT_ENTRYPOINT_NOT_EXECUTABLE")
    command += [
        "--candidate", candidate_id,
        "--candidate-config", str(config_path),
        "--output", str(result_path),
    ]

    before = resource.getrusage(resource.RUSAGE_CHILDREN)
    started = time.monotonic()
    try:
        proc = subprocess.run(
            command,
            cwd=str(pathlib.Path(package_root).resolve(strict=True)),
            env=build_candidate_env(output_root / "home", output_root / "tmp", manifest["data_class"]),
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            timeout=manifest["candidate_timeout_seconds"],
            check=False,
        )
    except subprocess.TimeoutExpired as exc:
        (output_root / "stdout.txt").write_text(exc.stdout or "", encoding="utf-8")
        (output_root / "stderr.txt").write_text(exc.stderr or "", encoding="utf-8")
        fail("FREE_FANOUT_CANDIDATE_TIMEOUT")
    elapsed = time.monotonic() - started
    after = resource.getrusage(resource.RUSAGE_CHILDREN)

    (output_root / "stdout.txt").write_text(proc.stdout, encoding="utf-8")
    (output_root / "stderr.txt").write_text(proc.stderr, encoding="utf-8")
    if proc.returncode != 0:
        fail(f"FREE_FANOUT_CANDIDATE_FAILED_EXIT_{proc.returncode}")
    if not result_path.is_file():
        fail("FREE_FANOUT_RESULT_MISSING")

    result = validate_result(load(result_path), candidate_id, manifest)
    peak_rss_kib = int(after.ru_maxrss)
    measurements = {
        "elapsed_seconds": round(elapsed, 6),
        "peak_rss_kib": peak_rss_kib,
        "user_cpu_seconds": round(max(0.0, after.ru_utime - before.ru_utime), 6),
        "system_cpu_seconds": round(max(0.0, after.ru_stime - before.ru_stime), 6),
        "parameters": result["resource_usage"]["parameters"],
        "context_bytes": result["resource_usage"]["context_bytes"],
        "model_calls": result["resource_usage"]["model_calls"],
    }
    limits = manifest["resource_envelope"]
    if elapsed > limits["max_elapsed_seconds"]:
        fail("FREE_FANOUT_RESOURCE_LIMIT_ELAPSED")
    if peak_rss_kib > limits["max_peak_rss_kib"]:
        fail("FREE_FANOUT_RESOURCE_LIMIT_PEAK_RSS")

    raw = result_path.read_bytes()
    envelope = {
        "schema": ENVELOPE_SCHEMA,
        "program": manifest["program"],
        "package_sha": manifest["package_sha"],
        "candidate_id": candidate_id,
        "candidate_config_sha256": canonical_json_sha256(candidate.get("config", {})),
        "candidate_interface_sha256": canonical_json_sha256(manifest["candidate_interface"]),
        "data_class": manifest["data_class"],
        "authority": "external-evidence-only",
        "resource_baseline_id": limits["baseline_id"],
        "resource_measurements": measurements,
        "metric": result["metric"],
        "result_sha256": hashlib.sha256(raw).hexdigest(),
        "scientific_classification": None,
        "promotion_authority": False,
        "successor_authority": False,
    }
    (output_root / "candidate-envelope.json").write_text(
        json.dumps(envelope, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    return envelope

def aggregate(manifest, artifacts_root):
    artifacts_root = pathlib.Path(artifacts_root)
    envelopes = []
    for path in artifacts_root.rglob("candidate-envelope.json"):
        envelope = load(path)
        if envelope.get("schema") != ENVELOPE_SCHEMA:
            fail("FREE_FANOUT_AGGREGATE_ENVELOPE_SCHEMA")
        if envelope.get("program") != manifest["program"]:
            fail("FREE_FANOUT_AGGREGATE_PROGRAM")
        if envelope.get("package_sha") != manifest["package_sha"]:
            fail("FREE_FANOUT_AGGREGATE_PACKAGE_SHA")
        if envelope.get("data_class") != manifest["data_class"]:
            fail("FREE_FANOUT_AGGREGATE_DATA_CLASS")
        if envelope.get("authority") != "external-evidence-only":
            fail("FREE_FANOUT_AGGREGATE_AUTHORITY")
        if envelope.get("scientific_classification") is not None:
            fail("FREE_FANOUT_AGGREGATE_CLASSIFICATION")
        if envelope.get("promotion_authority") is not False:
            fail("FREE_FANOUT_AGGREGATE_PROMOTION_AUTHORITY")
        if envelope.get("successor_authority") is not False:
            fail("FREE_FANOUT_AGGREGATE_SUCCESSOR_AUTHORITY")
        envelopes.append(envelope)

    expected = sorted(candidate["id"] for candidate in manifest["candidates"])
    found = [env.get("candidate_id") for env in envelopes]
    if len(found) != len(set(found)):
        fail("FREE_FANOUT_AGGREGATE_DUPLICATE_CANDIDATE")
    if sorted(found) != expected:
        fail("FREE_FANOUT_AGGREGATE_CANDIDATE_SET_MISMATCH")

    envelopes.sort(key=lambda item: item["candidate_id"])
    return {
        "schema": SUMMARY_SCHEMA,
        "program": manifest["program"],
        "package_sha": manifest["package_sha"],
        "authority": "external-evidence-only",
        "candidate_count": len(envelopes),
        "candidate_interface_sha256": canonical_json_sha256(manifest["candidate_interface"]),
        "resource_baseline_id": manifest["resource_envelope"]["baseline_id"],
        "candidates": envelopes,
        "scientific_classification": None,
        "promotion_authority": False,
        "successor_authority": False,
    }

def _manifest_from_cli(args):
    validate_manifest_path(args.manifest)
    return load_validated_manifest(args.manifest, args.program, args.package_sha)

def cmd_prepare(args):
    manifest = _manifest_from_cli(args)
    matrix = {"include": [{"id": candidate["id"]} for candidate in manifest["candidates"]]}
    print(json.dumps(matrix, separators=(",", ":")))

def cmd_run(args):
    manifest = _manifest_from_cli(args)
    run_candidate(manifest, args.package_root, args.candidate_id, args.output)

def cmd_aggregate(args):
    manifest = _manifest_from_cli(args)
    summary = aggregate(manifest, args.artifacts)
    out = pathlib.Path(args.output)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n", encoding="utf-8")

def main():
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(dest="cmd", required=True)

    prepare = subparsers.add_parser("prepare")
    prepare.add_argument("--manifest", required=True)
    prepare.add_argument("--program", required=True)
    prepare.add_argument("--package-sha", required=True)
    prepare.set_defaults(fn=cmd_prepare)

    run = subparsers.add_parser("run")
    run.add_argument("--manifest", required=True)
    run.add_argument("--program", required=True)
    run.add_argument("--package-sha", required=True)
    run.add_argument("--package-root", required=True)
    run.add_argument("--candidate-id", required=True)
    run.add_argument("--output", required=True)
    run.set_defaults(fn=cmd_run)

    aggregate_cmd = subparsers.add_parser("aggregate")
    aggregate_cmd.add_argument("--manifest", required=True)
    aggregate_cmd.add_argument("--program", required=True)
    aggregate_cmd.add_argument("--package-sha", required=True)
    aggregate_cmd.add_argument("--artifacts", required=True)
    aggregate_cmd.add_argument("--output", required=True)
    aggregate_cmd.set_defaults(fn=cmd_aggregate)

    args = parser.parse_args()
    try:
        args.fn(args)
    except FanoutError as exc:
        raise SystemExit(str(exc)) from exc

if __name__ == "__main__":
    main()
