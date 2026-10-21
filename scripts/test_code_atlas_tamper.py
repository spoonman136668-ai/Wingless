#!/usr/bin/env python3
"""Fail-closed negative qualification of the read-only Wingless ICE atlas.

Mutations are performed exclusively in a disposable local Git clone.
The checked-out PR tree, committed Wingless source, and research state remain untouched.
"""
import argparse
import fnmatch
import re
import json
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

ATLAS = Path(".ice/code-intelligence.json")
HTML = Path(".ice/code-intelligence.html")
SOURCE = Path("inference/backend.go")
PROOF_MARKER = "ATLAS_FAIL_CLOSED:"
MAX_TEST_SECONDS = 180


def execute(argv, *, cwd, timeout=MAX_TEST_SECONDS):
    return subprocess.run(
        [str(arg) for arg in argv], cwd=cwd, text=True,
        stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
        timeout=timeout, check=False,
    )


def require(ok, message):
    if not ok:
        raise SystemExit("TAMPER_QUALIFICATION_FAILED: " + message)


def jbytes(data):
    return (json.dumps(data, ensure_ascii=False, indent=2) + "\n").encode("utf-8")


def replace_json(data, change):
    doc = json.loads(data.decode("utf-8-sig"))
    change(doc)
    return jbytes(doc)


def run_negative(label, work, relative, transform, validator_command, expected):
    target = work / relative
    original = target.read_bytes()
    try:
        corrupted = transform(original)
        require(corrupted != original, label + ": mutation had no effect")
        target.write_bytes(corrupted)
        result = execute(validator_command, cwd=work)
        require(result.returncode != 0, label + ": forged input was incorrectly ACCEPTED")
        require(expected in result.stdout,
                label + ": rejected for unexpected reason; output=" + result.stdout[-1600:])
        print("EXPECTED_REJECTION: " + label + " / " + expected)
    finally:
        target.write_bytes(original)
    status = execute(["git", "status", "--porcelain"], cwd=work)
    require(status.returncode == 0 and not status.stdout.strip(),
            label + ": scratch checkout not fully restored")


def mutate_proof(doc):
    node = next(n for n in doc["nodes"] if n.get("review_evidence", {}).get("literal_tokens"))
    node["review_evidence"]["literal_tokens"][0] = "UNSUPPORTED_ATLAS_PROOF_MARKER_926374"


def mutate_dependency(doc):
    node = next(n for n in doc["nodes"] if n.get("file_level_local_package_dependencies"))
    node["file_level_local_package_dependencies"] = ["forged-internal-package"]


def mutate_contract(doc):
    node = next(n for n in doc["nodes"] if n.get("test_contract", {}).get("classification") == "test_case")
    node["test_contract"]["classification"] = "test_fixture_or_helper"


def mutate_identity(doc):
    doc["nodes"][0]["start_line"] += 1


def mutate_viewer(raw):
    src = raw.decode("utf-8")
    begin = '<script type="application/json" id="atlas">'
    start = src.find(begin)
    require(start >= 0, "viewer embedding not found")
    start += len(begin)
    end = src.find("</script>", start)
    require(end >= 0, "viewer JSON terminator not found")
    embedded = json.loads(src[start:end])
    embedded["nodes"][0]["label"] += " — FORGED_HTML_LABEL"
    encoded = json.dumps(embedded, ensure_ascii=True, separators=(",", ":")).replace("<", "\\u003c")
    return (src[:start] + encoded + src[end:]).encode("utf-8")



def mutate_provisional_promotion(doc):
    node = next(n for n in doc["nodes"]
                if n.get("evidence_basis") == "lexical-provisional")
    node["evidence_basis"] = "source-reviewed"

def mutate_provisional_source_token(doc):
    node = next(n for n in doc["nodes"]
                if n.get("evidence_basis") == "lexical-provisional" and
                n.get("provenance_lineage", {}).get("mode") == "new-symbol-source-bound-provisional")
    node["provenance_lineage"]["literal_tokens"] = [
        "FORGED_PROVISIONAL_TOKEN_MUST_NOT_EXIST_782149"]

def mutate_inherited_provisional_token(doc):
    node = next(n for n in doc["nodes"]
                if n.get("evidence_basis") == "lexical-provisional" and
                n.get("provenance_lineage", {}).get("mode") == "provisional-source-byte-reuse")
    node["provenance_lineage"]["literal_tokens"] = [
        "FORGED_PRIOR_PROVISIONAL_SOURCE_WITNESS_461782"]

def mutate_typed(raw):
    doc = json.loads(raw.decode("utf-8-sig"))
    require(doc.get("relationships"), "no relationships to remove")
    doc["relationships"] = []
    return jbytes(doc)


def verify_indexed_path_triggers(work):
    """Check all indexed paths, future sensitive paths and unconditional PR status."""
    result = execute([sys.executable, "-B",
                      str(work / "scripts/atlas_ci_scope.py"),
                      "--root", str(work), "--self-test"], cwd=work)
    require(result.returncode == 0 and
            "ATLAS_CI_SCOPE_CONTRACT_PASS" in result.stdout,
            "always-on PR scope contract failed: " + result.stdout[-1500:])
    print(result.stdout.strip())

def main():
    p = argparse.ArgumentParser()
    p.add_argument("--root", type=Path, required=True)
    p.add_argument("--relationships", type=Path, required=True)
    p.add_argument("--scratch-parent", type=Path, default=None)
    args = p.parse_args()
    root = args.root.resolve()
    report = args.relationships.resolve()
    require((root / ATLAS).is_file() and report.is_file(), "required inputs unavailable")
    require(not execute(["git", "status", "--porcelain"], cwd=root).stdout.strip(),
            "source checkout must be clean before destructive negative tests")
    commit = execute(["git", "rev-parse", "HEAD"], cwd=root)
    require(commit.returncode == 0 and len(commit.stdout.strip()) == 40, "unverified checkout HEAD")
    current_commit = commit.stdout.strip()
    with tempfile.TemporaryDirectory(prefix="wingless-atlas-negative-", dir=args.scratch_parent) as tmp:
        work = Path(tmp) / "scratch"
        clone = execute(["git", "clone", "--quiet", "--no-hardlinks", "--no-checkout",
                         str(root), str(work)], cwd=root, timeout=180)
        require(clone.returncode == 0, "failed scratch clone: " + clone.stdout[-1200:])
        checkout = execute(["git", "checkout", "--quiet", "--detach", current_commit], cwd=work)
        require(checkout.returncode == 0, "failed exact-head scratch checkout: " + checkout.stdout[-1200:])
        require(not execute(["git", "status", "--porcelain"], cwd=work).stdout.strip(), "scratch not clean")
        verify_indexed_path_triggers(work)

        validation = [sys.executable, str(work / "scripts/validate_code_atlas.py"),
                      "--root", str(work), "--relationships", str(report)]
        positive = execute(validation, cwd=work)
        require(positive.returncode == 0 and "ATLAS_SEMANTIC_QUALIFICATION_PASS" in positive.stdout,
                "positive control failed on exact head: " + positive.stdout[-1600:])
        print("POSITIVE_CONTROL_PASS: " + current_commit)
        run_negative("source-hash-mutation", work, SOURCE,
                     lambda b: b + b"\n// negative qualification source drift only\n",
                     validation, PROOF_MARKER + " source hash mismatch")
        run_negative("reviewed-proof-forgery", work, ATLAS,
                     lambda b: replace_json(b, mutate_proof),
                     validation, PROOF_MARKER + " source evidence token mismatch")
        run_negative("dependency-forgery", work, ATLAS,
                     lambda b: replace_json(b, mutate_dependency),
                     validation, PROOF_MARKER + " local package dependency identity drift")
        run_negative("test-contract-forgery", work, ATLAS,
                     lambda b: replace_json(b, mutate_contract),
                     validation, PROOF_MARKER + " incorrect test-case versus helper classification")
        run_negative("line-identity-forgery", work, ATLAS,
                     lambda b: replace_json(b, mutate_identity),
                     validation, PROOF_MARKER + " symbol lines mismatch")
        run_negative("embedded-viewer-divergence", work, HTML,
                     mutate_viewer, validation, PROOF_MARKER + " offline browser differs from machine index")

        # The independent type report is stored outside the scratch checkout.
        # A forged report must not satisfy any of the source-reviewed Go call edges.
        relation_copy = Path(tmp) / "typed-forgery.json"
        relation_copy.write_bytes(report.read_bytes())
        typed_validation = [sys.executable, str(work / "scripts/validate_code_atlas.py"),
                            "--root", str(work), "--relationships", str(relation_copy)]
        run_negative("type-resolved-edge-forgery", work, Path("../typed-forgery.json"),
                     mutate_typed, typed_validation, PROOF_MARKER + " type-unverified reviewed relationship")

        # G28 candidate-head additions are always unreviewed; actively prove
        # they cannot impersonate independent review or cite fabricated tokens.
        run_negative("provisional-promotion-forgery", work, ATLAS,
                     lambda b: replace_json(b, mutate_provisional_promotion),
                     validation, PROOF_MARKER + " source-reviewed symbol has no explicit proof")
        run_negative("provisional-source-token-forgery", work, ATLAS,
                     lambda b: replace_json(b, mutate_provisional_source_token),
                     validation, PROOF_MARKER + " G28 provisional receipt source tokens invalid")
        run_negative("inherited-provisional-source-token-forgery", work, ATLAS,
                     lambda b: replace_json(b, mutate_inherited_provisional_token),
                     validation, PROOF_MARKER + " G28 inherited provisional tokens do not bind source")

        # The original source-derived ICE gate must independently detect file drift.
        go_validation = ["go", "run", "./cmd/ice", "validate"]
        run_negative("canonical-ice-source-staleness", work, SOURCE,
                     lambda b: b + b"\n// negative qualification index staleness only\n",
                     go_validation, "ICE_STALE:")
        # G26/G28 maintenance of the original 352 human-reviewed declarations
        # must be tested on its *qualified historical source*; trying to relabel
        # the 51 new provisional declarations as historical reviewed entries is
        # explicitly forbidden. This second checkout is disposable and source-pinned.
        reviewed_base = "48ba343b403bc9b4389ce1777130e546b3c446bd"
        ancestry = execute(["git", "merge-base", "--is-ancestor",
                            reviewed_base, current_commit], cwd=work)
        require(ancestry.returncode == 0,
                "historical reviewed base does not precede candidate")
        old_work = Path(tmp) / "original-reviewed-352"
        hist_clone = execute(["git", "clone", "--quiet", "--no-hardlinks",
                              "--no-checkout", str(root), str(old_work)], cwd=root)
        require(hist_clone.returncode == 0, "historical scratch clone failed")
        hist_checkout = execute(["git", "checkout", "--quiet", "--detach",
                                 reviewed_base], cwd=old_work)
        require(hist_checkout.returncode == 0, "historical reviewed head unavailable")
        g26 = execute([sys.executable, "-B", "scripts/test_source_aware_atlas.py",
                       "--root", str(old_work), "--scratch-parent", str(args.scratch_parent or Path(tmp))],
                      cwd=old_work, timeout=900)
        require(g26.returncode == 0 and
                "G26_ISOLATED_MAINTENANCE_QUALIFICATION_PASS" in g26.stdout,
                "G26 source-aware regeneration qualification failed: " + g26.stdout[-2400:])
        print(g26.stdout)
        require(not execute(["git", "status", "--porcelain"], cwd=root).stdout.strip(),
                "original checkout was modified")
        print(json.dumps({
            "result": "WINGLESS_ATLAS_NEGATIVE_QUALIFICATION_PASS",
            "head": current_commit,
            "expected_rejections": 11,
            "positive_control": "PASS",
            "source_checkout_changed": False,
            "scratch_cleanup": "automatic",
        }, sort_keys=True))



if __name__ == "__main__":
    main()
