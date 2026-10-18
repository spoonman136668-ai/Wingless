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


def mutate_typed(raw):
    doc = json.loads(raw.decode("utf-8-sig"))
    require(doc.get("relationships"), "no relationships to remove")
    doc["relationships"] = []
    return jbytes(doc)


def verify_indexed_path_triggers(work):
    """Assert that source-only changes cannot bypass the GitHub PR freshness gate."""
    workflow = (work / ".github/workflows/wingless-code-atlas.yml").read_text(encoding="utf-8")
    match = re.search(r"(?m)^  pull_request:\n    paths:\n((?:      - '.+'\n)+)", workflow)
    require(match is not None, "pull_request indexed-source path gate missing")
    patterns = re.findall(r"(?m)^      - '([^']+)'$", match.group(1))
    require(patterns, "no indexed-source patterns")
    indexed = json.loads((work / ".ice/manifest.json").read_text(encoding="utf-8"))["files"]
    # Synthetic paths cover newly introduced files not yet in the existing index.
    sources = [record["path"] for record in indexed] + [
        "new_source.go", "sub/new_source.go", "README.md", "sub/new.md",
        "src/new.json", "src/new.ps1", ".github/workflows/fresh.yml",
        "src/new.yaml", "go.mod", ".gitattributes", ".gitignore",
        "scripts/validate_code_atlas.py", "scripts/test_code_atlas_tamper.py",
    ]
    def match_path(path, pattern):
        return fnmatch.fnmatchcase(path, pattern) or (
            pattern.startswith("**/") and fnmatch.fnmatchcase(path, pattern[3:])
        )
    uncovered = sorted(set(path for path in sources
                           if not any(match_path(path, rule) for rule in patterns)))
    require(not uncovered, "indexed source can bypass PR gate: " + repr(uncovered))
    print(json.dumps({"result": "ATLAS_PR_SOURCE_PATH_COVERAGE_PASS",
                      "indexed_files": len(indexed), "globs": len(patterns),
                      "future_path_probes": len(sources) - len(indexed)}, sort_keys=True))


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

        # The original source-derived ICE gate must independently detect file drift.
        go_validation = ["go", "run", "./cmd/ice", "validate"]
        run_negative("canonical-ice-source-staleness", work, SOURCE,
                     lambda b: b + b"\n// negative qualification index staleness only\n",
                     go_validation, "ICE_STALE:")
        require(not execute(["git", "status", "--porcelain"], cwd=root).stdout.strip(),
                "original checkout was modified")
        print(json.dumps({
            "result": "WINGLESS_ATLAS_NEGATIVE_QUALIFICATION_PASS",
            "head": current_commit,
            "expected_rejections": 8,
            "positive_control": "PASS",
            "source_checkout_changed": False,
            "scratch_cleanup": "automatic",
        }, sort_keys=True))
        # A real Chrome browser acceptance occurs after existing fail-closed
        # negative verification. Changes are confined to a disposable clone.
        # Pin Selenium to make the driver contract reproducible; no model/API.
        install = execute([sys.executable, "-m", "pip", "install",
                           "--disable-pip-version-check", "--quiet",
                           "selenium==4.35.0"], cwd=work, timeout=240)
        require(install.returncode == 0,
                "browser driver dependency setup failed: " + install.stdout[-1600:])
        browser = execute([sys.executable, str(work / "scripts/test_atlas_browser.py"),
                           "--root", str(work)], cwd=work, timeout=240)
        print(browser.stdout)
        require(browser.returncode == 0 and
                "WINGLESS_ATLAS_CHROMIUM_BROWSER_ACCEPTANCE_PASS" in browser.stdout,
                "real Chromium UI acceptance failed: " + browser.stdout[-5000:])
        require(not execute(["git", "status", "--porcelain"], cwd=work).stdout.strip(),
                "browser qualification modified the scratch worktree")



if __name__ == "__main__":
    main()
