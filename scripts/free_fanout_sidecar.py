#!/usr/bin/env python3
import argparse
import hashlib
import json
import os
import pathlib
import re
import subprocess
import sys
import time

SCHEMA = "research.free-fanout.v1"
ID_RE = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$")
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
ALLOWED_DATA = {"historical", "replay", "historical-replay"}

def load(path):
    with open(path, "r", encoding="utf-8") as f:
        return json.load(f)

def validate(m, program, package_sha):
    errs = []
    if m.get("schema") != SCHEMA: errs.append("schema")
    if m.get("program") != program: errs.append("program")
    if m.get("package_sha") != package_sha or not SHA_RE.fullmatch(package_sha): errs.append("package_sha")
    if m.get("authority") != "external-evidence-only": errs.append("authority")
    if m.get("data_class") not in ALLOWED_DATA: errs.append("data_class")
    for k in ("sealed_inputs", "accepted_state_mutation", "successor_dispatch", "queue_write", "production_access"):
        if m.get(k) is not False: errs.append(k)
    ep = m.get("entrypoint")
    if not isinstance(ep, str) or not ep.startswith("research/fanout/") or ".." in pathlib.PurePosixPath(ep).parts:
        errs.append("entrypoint")
    timeout = m.get("candidate_timeout_seconds", 900)
    if not isinstance(timeout, int) or not (1 <= timeout <= 1800): errs.append("candidate_timeout_seconds")
    candidates = m.get("candidates")
    if not isinstance(candidates, list) or not (1 <= len(candidates) <= 32):
        errs.append("candidates")
    else:
        seen = set()
        for c in candidates:
            cid = c.get("id") if isinstance(c, dict) else None
            if not isinstance(cid, str) or not ID_RE.fullmatch(cid) or cid in seen:
                errs.append("candidate_id")
                break
            seen.add(cid)
    if errs:
        raise SystemExit("FREE_FANOUT_MANIFEST_INVALID " + ",".join(sorted(set(errs))))
    return m

def cmd_prepare(args):
    m = validate(load(args.manifest), args.program, args.package_sha)
    matrix = {"include": [{"id": c["id"]} for c in m["candidates"]]}
    print(json.dumps(matrix, separators=(",", ":")))

def cmd_run(args):
    m = validate(load(args.manifest), args.program, args.package_sha)
    candidate = next((c for c in m["candidates"] if c["id"] == args.candidate_id), None)
    if candidate is None:
        raise SystemExit("FREE_FANOUT_CANDIDATE_NOT_FOUND")
    root = pathlib.Path(args.output)
    root.mkdir(parents=True, exist_ok=True)
    result_path = root / "result.json"
    ep = pathlib.Path(m["entrypoint"])
    if not ep.is_file():
        raise SystemExit("FREE_FANOUT_ENTRYPOINT_MISSING")
    if ep.suffix == ".py":
        command = [sys.executable, str(ep)]
    elif ep.suffix == ".sh":
        command = ["bash", str(ep)]
    else:
        command = [str(ep)]
    command += ["--candidate", args.candidate_id, "--output", str(result_path)]
    env = os.environ.copy()
    for key in list(env):
        if key.endswith("_TOKEN") or key.endswith("_KEY") or key in {"GITHUB_TOKEN", "GH_TOKEN"}:
            env.pop(key, None)
    env["FREE_FANOUT_AUTHORITY"] = "external-evidence-only"
    env["FREE_FANOUT_DATA_CLASS"] = m["data_class"]
    start = time.monotonic()
    p = subprocess.run(
        command,
        env=env,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        timeout=m.get("candidate_timeout_seconds", 900),
        check=False,
    )
    elapsed = time.monotonic() - start
    (root / "stdout.txt").write_text(p.stdout, encoding="utf-8")
    (root / "stderr.txt").write_text(p.stderr, encoding="utf-8")
    if p.returncode != 0:
        raise SystemExit(f"FREE_FANOUT_CANDIDATE_FAILED exit={p.returncode}")
    if not result_path.is_file():
        raise SystemExit("FREE_FANOUT_RESULT_MISSING")
    result = load(result_path)
    if result.get("candidate_id") != args.candidate_id:
        raise SystemExit("FREE_FANOUT_RESULT_ID_MISMATCH")
    raw = result_path.read_bytes()
    envelope = {
        "schema": "research.free-fanout-envelope.v1",
        "program": args.program,
        "package_sha": args.package_sha,
        "candidate_id": args.candidate_id,
        "data_class": m["data_class"],
        "authority": "external-evidence-only",
        "elapsed_seconds": round(elapsed, 6),
        "result_sha256": hashlib.sha256(raw).hexdigest(),
    }
    (root / "candidate-envelope.json").write_text(
        json.dumps(envelope, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )

def cmd_aggregate(args):
    root = pathlib.Path(args.artifacts)
    envs = []
    for p in root.rglob("candidate-envelope.json"):
        envs.append(load(p))
    envs.sort(key=lambda x: x["candidate_id"])
    summary = {
        "schema": "research.free-fanout-summary.v1",
        "authority": "external-evidence-only",
        "candidate_count": len(envs),
        "candidates": envs,
        "scientific_classification": None,
        "promotion_authority": False,
        "successor_authority": False,
    }
    out = pathlib.Path(args.output)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n", encoding="utf-8")

def main():
    ap = argparse.ArgumentParser()
    sp = ap.add_subparsers(dest="cmd", required=True)
    p = sp.add_parser("prepare")
    p.add_argument("--manifest", required=True)
    p.add_argument("--program", required=True)
    p.add_argument("--package-sha", required=True)
    p.set_defaults(fn=cmd_prepare)
    p = sp.add_parser("run")
    p.add_argument("--manifest", required=True)
    p.add_argument("--program", required=True)
    p.add_argument("--package-sha", required=True)
    p.add_argument("--candidate-id", required=True)
    p.add_argument("--output", required=True)
    p.set_defaults(fn=cmd_run)
    p = sp.add_parser("aggregate")
    p.add_argument("--artifacts", required=True)
    p.add_argument("--output", required=True)
    p.set_defaults(fn=cmd_aggregate)
    a = ap.parse_args()
    a.fn(a)

if __name__ == "__main__":
    main()
