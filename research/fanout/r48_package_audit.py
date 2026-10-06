import argparse
import json
import os
import pathlib
import subprocess
import sys

INTERFACE = "state_former(raw_input, history) -> fixed_size_state"
PREREG_PATH = pathlib.Path("docs/experiments/wlm-lm-external-future-data-learned-state-former-feature-attribution-r48.json")
MANIFEST_PATH = pathlib.Path("docs/experiments/wlm-lm-external-future-data-r48-manifest.json")
EXPECTED_PREREG_BLOB = "2d63f093e9b7d06b7f5ad1ef6d427bbb12496301"
EXPECTED_MANIFEST_BLOB = "83b988f17ad4c51351dce390ba0b8eeafb1e3595"
EXPECTED_PREREG_SHA = "d39cd28bcf2543f331ac4557ac5e6a3cc37500c5"
EXPECTED_PARENT_SHA = "465609d9845adfa464af363d00313b9d1d10698e"
EXPECTED_MANIFEST_SHA = "73c4bcdb292a27ccd403ef691db09c9ec436ad137d907423615fe69cced8b1fb"
EXPECTED_EXPERIMENT = "WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-FEATURE-ATTRIBUTION-R48"

MODES = {
    "identity": None,
    "unitary-focused": ["go", "test", "./unitary", "-count=1"],
    "cmd-focused": ["go", "test", "./cmd/wlm-lm-external-future-data-learned-state-former-feature-attribution-r48", "-count=1"],
    "full-regression": ["go", "test", "./...", "-count=1"],
}

def fail(message):
    raise SystemExit(message)

def git_blob(path):
    proc = subprocess.run(["git", "hash-object", str(path)], text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False)
    if proc.returncode != 0:
        fail("R48_FANOUT_GIT_HASH_FAILED:" + proc.stderr.strip())
    return proc.stdout.strip()

def verify_identity():
    for p in (PREREG_PATH, MANIFEST_PATH):
        if not p.is_file() or p.is_symlink():
            fail("R48_FANOUT_IDENTITY_FILE_MISSING:" + str(p))
    if git_blob(PREREG_PATH) != EXPECTED_PREREG_BLOB:
        fail("R48_FANOUT_PREREG_BLOB_MISMATCH")
    if git_blob(MANIFEST_PATH) != EXPECTED_MANIFEST_BLOB:
        fail("R48_FANOUT_MANIFEST_BLOB_MISMATCH")
    prereg = json.loads(PREREG_PATH.read_text(encoding="utf-8"))
    manifest = json.loads(MANIFEST_PATH.read_text(encoding="utf-8"))
    if prereg.get("experiment_id") != EXPECTED_EXPERIMENT:
        fail("R48_FANOUT_EXPERIMENT_ID_MISMATCH")
    if prereg.get("parent_sha") != EXPECTED_PARENT_SHA:
        fail("R48_FANOUT_PARENT_SHA_MISMATCH")
    if manifest.get("manifest_sha256") != EXPECTED_MANIFEST_SHA:
        fail("R48_FANOUT_MANIFEST_SHA_MISMATCH")
    if not pathlib.Path("unitary/wlm_lm_external_future_data_learned_state_former_feature_attribution_r48.go").is_file():
        fail("R48_FANOUT_UNITARY_SURFACE_MISSING")
    if not pathlib.Path("cmd/wlm-lm-external-future-data-learned-state-former-feature-attribution-r48/main.go").is_file():
        fail("R48_FANOUT_CMD_SURFACE_MISSING")
    if not pathlib.Path(".github/workflows/external-future-data-learned-state-former-feature-attribution-r48.yml").is_file():
        fail("R48_FANOUT_WORKFLOW_SURFACE_MISSING")

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--candidate", required=True)
    parser.add_argument("--candidate-config", required=True)
    parser.add_argument("--output", required=True)
    args = parser.parse_args()

    if os.environ.get("FREE_FANOUT_AUTHORITY") != "external-evidence-only":
        fail("R48_FANOUT_AUTHORITY_ENV_INVALID")
    if os.environ.get("FREE_FANOUT_DATA_CLASS") != "historical-replay":
        fail("R48_FANOUT_DATA_CLASS_ENV_INVALID")

    cfg = json.loads(pathlib.Path(args.candidate_config).read_text(encoding="utf-8"))
    if set(cfg) != {"mode"}:
        fail("R48_FANOUT_CONFIG_KEYS_INVALID")
    mode = cfg["mode"]
    if mode not in MODES:
        fail("R48_FANOUT_MODE_INVALID")
    if args.candidate != mode:
        fail("R48_FANOUT_CANDIDATE_MODE_MISMATCH")

    verify_identity()
    command = MODES[mode]
    if command is not None:
        proc = subprocess.run(command, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False)
        sys.stdout.write(proc.stdout)
        sys.stderr.write(proc.stderr)
        if proc.returncode != 0:
            fail("R48_FANOUT_AUDIT_COMMAND_FAILED:" + mode)

    result = {
        "candidate_id": args.candidate,
        "interface": INTERFACE,
        "metric": {"name": "package_audit_pass", "value": 1},
        "resource_usage": {"parameters": 0, "context_bytes": 0, "model_calls": 0},
        "fixed_state_dimension": 6,
        "fixed_readout_capacity": 7,
    }
    out = pathlib.Path(args.output)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(result, sort_keys=True, separators=(",", ":")) + "\n", encoding="utf-8")

if __name__ == "__main__":
    main()
