#!/usr/bin/env python3
import argparse
import hashlib
import json
import pathlib

EXPORT_SCHEMA = "research.portable-checkpoint-export.v1"
CHECKPOINT_SCHEMA = "research.portable-checkpoint.v1"
SHA1_RE = __import__("re").compile(r"^[0-9a-f]{40}$")
SHA256_RE = __import__("re").compile(r"^[0-9a-f]{64}$")
ALLOWED_ROLES = {"learned_parameters", "resumable_state", "replay_input", "replay_output", "metadata"}
REQUIRED_ROLES = {"learned_parameters", "resumable_state", "replay_input", "replay_output"}
PROJECTS = {"Wingless", "Yggdrasil"}

class PortabilityError(RuntimeError):
    pass

def fail(code):
    raise PortabilityError(code)

def canonical_bytes(value):
    return json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=True).encode("utf-8")

def sha256_bytes(raw):
    return hashlib.sha256(raw).hexdigest()

def load_json(path):
    try:
        return json.loads(pathlib.Path(path).read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        raise PortabilityError(f"PORTABLE_CHECKPOINT_JSON_LOAD_FAILED:{type(exc).__name__}") from exc

def validate_rel_path(value):
    if not isinstance(value, str) or not value or "\\" in value:
        fail("PORTABLE_CHECKPOINT_PATH_INVALID")
    parts = value.split("/")
    if any(part in ("", ".", "..") for part in parts):
        fail("PORTABLE_CHECKPOINT_PATH_INVALID")
    p = pathlib.PurePosixPath(value)
    if p.is_absolute():
        fail("PORTABLE_CHECKPOINT_PATH_INVALID")
    return value

def source_file(root, rel):
    validate_rel_path(rel)
    root = pathlib.Path(root).resolve(strict=True)
    candidate = root / pathlib.PurePosixPath(rel)
    probe = root
    for part in pathlib.PurePosixPath(rel).parts:
        probe = probe / part
        if probe.is_symlink():
            fail("PORTABLE_CHECKPOINT_SYMLINK_REJECTED")
    resolved = candidate.resolve(strict=True)
    if not resolved.is_file() or not resolved.is_relative_to(root):
        fail("PORTABLE_CHECKPOINT_PATH_ESCAPED")
    return resolved

def require_sha(value, kind):
    regex = SHA1_RE if kind == "sha1" else SHA256_RE
    if not isinstance(value, str) or not regex.fullmatch(value):
        fail(f"PORTABLE_CHECKPOINT_{kind.upper()}_INVALID")

def validate_resource_envelope(value):
    if not isinstance(value, dict):
        fail("PORTABLE_CHECKPOINT_RESOURCE_ENVELOPE")
    if value.get("comparison") != "equal-or-lower-than-baseline":
        fail("PORTABLE_CHECKPOINT_RESOURCE_COMPARISON")
    if not isinstance(value.get("baseline_id"), str) or not value["baseline_id"].strip():
        fail("PORTABLE_CHECKPOINT_RESOURCE_BASELINE")
    for key in ("max_parameters", "max_context_bytes", "max_model_calls"):
        v = value.get(key)
        if not isinstance(v, int) or isinstance(v, bool) or v < 0:
            fail(f"PORTABLE_CHECKPOINT_RESOURCE_{key.upper()}")

def validate_interface(project, interface):
    if not isinstance(interface, dict):
        fail("PORTABLE_CHECKPOINT_INTERFACE")
    if project == "Wingless":
        if interface.get("signature") != "state_former(raw_input, history) -> fixed_size_state":
            fail("PORTABLE_CHECKPOINT_WINGLESS_SIGNATURE")
        for key in ("fixed_state_dimension", "fixed_readout_capacity"):
            v = interface.get(key)
            if not isinstance(v, int) or isinstance(v, bool) or v <= 0:
                fail(f"PORTABLE_CHECKPOINT_WINGLESS_{key.upper()}")
        if interface.get("capacity_growth") is not False:
            fail("PORTABLE_CHECKPOINT_WINGLESS_CAPACITY_GROWTH")
    elif project == "Yggdrasil":
        if interface.get("signature") != "cognition_consumer(retained_state, local_state) -> decision_state":
            fail("PORTABLE_CHECKPOINT_YGGDRASIL_SIGNATURE")
        expected = {"total_slots": 16, "active_slots": 7, "retained_slots": 9}
        for key, val in expected.items():
            if interface.get(key) != val:
                fail(f"PORTABLE_CHECKPOINT_YGGDRASIL_{key.upper()}")
        for key in ("hidden_persistent_memory_growth", "addressing_change", "capacity_growth"):
            if interface.get(key) is not False:
                fail(f"PORTABLE_CHECKPOINT_YGGDRASIL_{key.upper()}")

def validate_spec(spec):
    if not isinstance(spec, dict) or spec.get("schema") != EXPORT_SCHEMA:
        fail("PORTABLE_CHECKPOINT_SPEC_SCHEMA")
    project = spec.get("project")
    if project not in PROJECTS:
        fail("PORTABLE_CHECKPOINT_PROJECT")
    for key in ("source_commit_sha", "source_tree_sha"):
        require_sha(spec.get(key), "sha1")
    for key in ("source_schema_sha256", "qualification_contract_sha256"):
        require_sha(spec.get(key), "sha256")
    if not isinstance(spec.get("toolchain_id"), str) or not spec["toolchain_id"].strip():
        fail("PORTABLE_CHECKPOINT_TOOLCHAIN")
    if not isinstance(spec.get("mechanism_id"), str) or not spec["mechanism_id"].strip():
        fail("PORTABLE_CHECKPOINT_MECHANISM_ID")
    lineage = spec.get("lineage")
    if not isinstance(lineage, dict) or not isinstance(lineage.get("experiment_id"), str) or not lineage["experiment_id"].strip():
        fail("PORTABLE_CHECKPOINT_LINEAGE")
    parent = lineage.get("parent_checkpoint_sha256")
    if parent is not None:
        require_sha(parent, "sha256")
    validate_resource_envelope(spec.get("resource_envelope"))
    validate_interface(project, spec.get("interface"))

    files = spec.get("files")
    if not isinstance(files, list) or not files:
        fail("PORTABLE_CHECKPOINT_FILES")
    seen, roles = set(), set()
    for entry in files:
        if not isinstance(entry, dict):
            fail("PORTABLE_CHECKPOINT_FILE_ENTRY")
        rel = validate_rel_path(entry.get("path"))
        role = entry.get("role")
        if role not in ALLOWED_ROLES:
            fail("PORTABLE_CHECKPOINT_FILE_ROLE")
        if rel in seen:
            fail("PORTABLE_CHECKPOINT_DUPLICATE_PATH")
        seen.add(rel)
        roles.add(role)
    if not REQUIRED_ROLES.issubset(roles):
        fail("PORTABLE_CHECKPOINT_REQUIRED_ROLE_MISSING")

    replay = spec.get("replay")
    if not isinstance(replay, dict) or not isinstance(replay.get("vectors"), list) or not replay["vectors"]:
        fail("PORTABLE_CHECKPOINT_REPLAY")
    file_map = {e["path"]: e["role"] for e in files}
    ids = set()
    for vector in replay["vectors"]:
        if not isinstance(vector, dict) or not isinstance(vector.get("id"), str) or not vector["id"]:
            fail("PORTABLE_CHECKPOINT_REPLAY_VECTOR")
        if vector["id"] in ids:
            fail("PORTABLE_CHECKPOINT_REPLAY_DUPLICATE_ID")
        ids.add(vector["id"])
        ip, op = vector.get("input_path"), vector.get("output_path")
        if file_map.get(ip) != "replay_input" or file_map.get(op) != "replay_output":
            fail("PORTABLE_CHECKPOINT_REPLAY_PATH_ROLE")

    authority = spec.get("authority")
    if authority != {
        "accepted_state_mutation": False,
        "execution_authority": False,
        "promotion_authority": False,
        "production_authority": False,
        "successor_authority": False,
        "retraining_required": False,
    }:
        fail("PORTABLE_CHECKPOINT_AUTHORITY")
    return spec

def export_checkpoint(spec_path, root, output):
    spec = validate_spec(load_json(spec_path))
    root = pathlib.Path(root).resolve(strict=True)
    out = pathlib.Path(output)
    if out.exists() and any(out.iterdir()):
        fail("PORTABLE_CHECKPOINT_OUTPUT_NOT_EMPTY")
    out.mkdir(parents=True, exist_ok=True)
    files_root = out / "files"
    files_root.mkdir()

    records = []
    for entry in spec["files"]:
        src = source_file(root, entry["path"])
        raw = src.read_bytes()
        dst = files_root / pathlib.PurePosixPath(entry["path"])
        dst.parent.mkdir(parents=True, exist_ok=True)
        dst.write_bytes(raw)
        records.append({
            "path": entry["path"],
            "role": entry["role"],
            "size_bytes": len(raw),
            "sha256": sha256_bytes(raw),
        })
    records.sort(key=lambda x: x["path"])

    replay_vectors = []
    by_path = {r["path"]: r for r in records}
    for vector in spec["replay"]["vectors"]:
        replay_vectors.append({
            "id": vector["id"],
            "input_path": vector["input_path"],
            "input_sha256": by_path[vector["input_path"]]["sha256"],
            "output_path": vector["output_path"],
            "output_sha256": by_path[vector["output_path"]]["sha256"],
        })
    replay_vectors.sort(key=lambda x: x["id"])

    payload = {
        "schema": CHECKPOINT_SCHEMA,
        "project": spec["project"],
        "mechanism_id": spec["mechanism_id"],
        "qualification_identity": {
            "source_commit_sha": spec["source_commit_sha"],
            "source_tree_sha": spec["source_tree_sha"],
            "source_schema_sha256": spec["source_schema_sha256"],
            "qualification_contract_sha256": spec["qualification_contract_sha256"],
            "toolchain_id": spec["toolchain_id"],
        },
        "lineage": spec["lineage"],
        "interface": spec["interface"],
        "resource_envelope": spec["resource_envelope"],
        "files": records,
        "replay": {
            "vectors": replay_vectors,
            "runtime_replay_required_before_promotion": True,
        },
        "portability": {
            "artifact_boundary": "language-neutral",
            "retraining_required": False,
            "restore_supported": True,
            "replay_vectors_bundled": True,
            "requalification_required_before_promotion": True,
        },
        "authority": spec["authority"],
        "scientific_classification": None,
    }
    payload["checkpoint_content_sha256"] = sha256_bytes(canonical_bytes(payload))
    (out / "manifest.json").write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return payload

def verify_checkpoint(checkpoint):
    root = pathlib.Path(checkpoint).resolve(strict=True)
    manifest = load_json(root / "manifest.json")
    if manifest.get("schema") != CHECKPOINT_SCHEMA:
        fail("PORTABLE_CHECKPOINT_MANIFEST_SCHEMA")
    expected = manifest.get("checkpoint_content_sha256")
    require_sha(expected, "sha256")
    copy = dict(manifest)
    copy.pop("checkpoint_content_sha256", None)
    if sha256_bytes(canonical_bytes(copy)) != expected:
        fail("PORTABLE_CHECKPOINT_MANIFEST_TAMPER")
    if manifest.get("scientific_classification") is not None:
        fail("PORTABLE_CHECKPOINT_CLASSIFICATION_AUTHORITY")
    auth = manifest.get("authority")
    if not isinstance(auth, dict) or any(auth.get(k) is not False for k in (
        "accepted_state_mutation", "execution_authority", "promotion_authority",
        "production_authority", "successor_authority", "retraining_required",
    )):
        fail("PORTABLE_CHECKPOINT_AUTHORITY")
    validate_interface(manifest.get("project"), manifest.get("interface"))
    validate_resource_envelope(manifest.get("resource_envelope"))
    for key in ("source_commit_sha", "source_tree_sha"):
        require_sha(manifest.get("qualification_identity", {}).get(key), "sha1")
    for key in ("source_schema_sha256", "qualification_contract_sha256"):
        require_sha(manifest.get("qualification_identity", {}).get(key), "sha256")
    if not manifest.get("qualification_identity", {}).get("toolchain_id"):
        fail("PORTABLE_CHECKPOINT_TOOLCHAIN")

    seen = set()
    for record in manifest.get("files", []):
        rel = validate_rel_path(record.get("path"))
        if rel in seen:
            fail("PORTABLE_CHECKPOINT_DUPLICATE_PATH")
        seen.add(rel)
        p = source_file(root / "files", rel)
        raw = p.read_bytes()
        if len(raw) != record.get("size_bytes") or sha256_bytes(raw) != record.get("sha256"):
            fail("PORTABLE_CHECKPOINT_FILE_TAMPER")
    file_map = {r["path"]: r for r in manifest["files"]}
    for vector in manifest.get("replay", {}).get("vectors", []):
        if file_map.get(vector.get("input_path"), {}).get("sha256") != vector.get("input_sha256"):
            fail("PORTABLE_CHECKPOINT_REPLAY_INPUT_TAMPER")
        if file_map.get(vector.get("output_path"), {}).get("sha256") != vector.get("output_sha256"):
            fail("PORTABLE_CHECKPOINT_REPLAY_OUTPUT_TAMPER")
    return manifest

def restore_checkpoint(checkpoint, destination):
    manifest = verify_checkpoint(checkpoint)
    dst = pathlib.Path(destination)
    if dst.exists() and any(dst.iterdir()):
        fail("PORTABLE_CHECKPOINT_RESTORE_NOT_EMPTY")
    dst.mkdir(parents=True, exist_ok=True)
    src_root = pathlib.Path(checkpoint).resolve(strict=True) / "files"
    for record in manifest["files"]:
        src = source_file(src_root, record["path"])
        out = dst / pathlib.PurePosixPath(record["path"])
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_bytes(src.read_bytes())
    receipt = {
        "schema": "research.portable-checkpoint-restore.v1",
        "project": manifest["project"],
        "mechanism_id": manifest["mechanism_id"],
        "checkpoint_content_sha256": manifest["checkpoint_content_sha256"],
        "retraining_performed": False,
        "runtime_replay_required": True,
        "requalification_required": True,
        "promotion_authority": False,
    }
    (dst / "restore-receipt.json").write_text(json.dumps(receipt, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return receipt

def main():
    ap = argparse.ArgumentParser()
    sp = ap.add_subparsers(dest="cmd", required=True)
    p = sp.add_parser("export")
    p.add_argument("--spec", required=True)
    p.add_argument("--root", required=True)
    p.add_argument("--output", required=True)
    p = sp.add_parser("verify")
    p.add_argument("--checkpoint", required=True)
    p = sp.add_parser("restore")
    p.add_argument("--checkpoint", required=True)
    p.add_argument("--destination", required=True)
    args = ap.parse_args()
    try:
        if args.cmd == "export":
            export_checkpoint(args.spec, args.root, args.output)
        elif args.cmd == "verify":
            verify_checkpoint(args.checkpoint)
        else:
            restore_checkpoint(args.checkpoint, args.destination)
    except PortabilityError as exc:
        raise SystemExit(str(exc)) from exc

if __name__ == "__main__":
    main()
