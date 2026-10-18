#!/usr/bin/env python3
"""Fail-closed scope for the always-present Wingless atlas CI status check.

The status check runs on every PR. Exact source/atlas/browser-interface changes
perform the full Windows qualification; unrelated files do not consume that budget.
"""
import argparse
import json
import re
import subprocess
from pathlib import Path

SCANNED_EXTENSIONS = (".go", ".md", ".json", ".ps1", ".yml", ".yaml")
ATLAS_EXTENSIONS = (".html", ".ice")
SPECIAL_BASENAMES = {"go.mod", ".gitattributes", ".gitignore"}
ATLAS_SCRIPTS = {
    "scripts/validate_code_atlas.py",
    "scripts/test_code_atlas_tamper.py",
    "scripts/test_atlas_browser.py",
    "scripts/atlas_ci_scope.py",
}
WORKFLOW = ".github/workflows/wingless-code-atlas.yml"
STATUS_JOB = "Read-only Windows Atlas Qualification"


def require(ok, message):
    if not ok:
        raise SystemExit("ATLAS_CI_SCOPE_FAIL_CLOSED: " + message)


def relevant(path):
    path = path.replace("\\", "/")
    name = path.rsplit("/", 1)[-1].lower()
    return (path in ATLAS_SCRIPTS
            or name in SPECIAL_BASENAMES
            or name.endswith(SCANNED_EXTENSIONS + ATLAS_EXTENSIONS))


def git(root, *args):
    proc = subprocess.run(["git", *args], cwd=root, capture_output=True)
    require(proc.returncode == 0,
            "git failed: " + " ".join(args) + " / " +
            proc.stderr.decode("utf-8", errors="replace")[-300:])
    return proc.stdout


def self_test(root):
    workflow = (root / WORKFLOW).read_text(encoding="utf-8")
    header = workflow.split("\npermissions:", 1)[0]
    require(re.search(r"(?m)^on:\n  pull_request:\s*\n  workflow_dispatch:", header) is not None,
            "pull_request not unconditional")
    require(re.search(r"(?m)^\s+paths(?:-ignore)?:", header) is None,
            "PR-level path filtering can strand required status")
    require("    name: " + STATUS_JOB in workflow, "required job identity changed")
    require("        id: scope\n" in workflow,
            "atlas scope classification step absent")
    for step in (
        "Set up frozen Go toolchain",
        "Validate source-index freshness and typed relationships",
        "Tamper-resistance negative acceptance in disposable exact-head clone",
        "Genuine Chromium browser navigation acceptance",
        "Focused ICE regression",
        "One complete integration regression",
        "Vet, build, and verify clean checkout",
    ):
        marker = "      - name: " + step + "\n"
        require(workflow.count(marker) == 1, "missing/duplicate step: " + step)
        block = workflow.split(marker, 1)[1].split("\n      - name:", 1)[0]
        require("if: steps.scope.outputs.atlas_relevant == 'true'" in block,
                "full qualification can bypass scope: " + step)
    entries = json.loads((root / ".ice/manifest.json").read_text("utf-8"))["files"]
    representative = [
        "new_source.go", "src/new_source.go", "README.md", "src/new.md",
        "src/new.json", "src/new.ps1", ".github/workflows/fresh.yml",
        "src/new.yaml", "go.mod", ".gitattributes", ".gitignore",
        "scripts/validate_code_atlas.py", "scripts/test_code_atlas_tamper.py",
        "scripts/test_atlas_browser.py", "scripts/atlas_ci_scope.py",
        ".ice/code-intelligence.html", ".ice/code-intelligence.ice",
        "viewer.html", "deep/new_viewer.html", "notes.ice", "deep/new.ice",
        ".github/workflows/wingless-code-atlas.yml",
    ]
    uncovered = [item for item in [e["path"] for e in entries] + representative
                 if not relevant(item)]
    require(not uncovered, "indexed or future sensitive paths skipped: " + repr(uncovered))
    for name in ("notes.txt", "other/LICENSE", "scripts/unrelated.py"):
        require(not relevant(name), "unrelated path unnecessarily triggers full suite: " + name)
    print(json.dumps({
        "result": "ATLAS_CI_SCOPE_CONTRACT_PASS",
        "indexed_files": len(entries),
        "future_path_probes": len(representative),
        "always_on_status": STATUS_JOB,
        "unrelated_path_skips": 3,
    }, sort_keys=True))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, default=Path("."))
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--event")
    parser.add_argument("--base", default="")
    parser.add_argument("--expected-head", default="")
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()
    root = args.root.resolve()
    if args.self_test:
        self_test(root)
        return
    require(args.event in ("pull_request", "workflow_dispatch"),
            "unknown trigger: " + str(args.event))
    head = git(root, "rev-parse", "HEAD").decode("ascii").strip()
    require(re.fullmatch(r"[0-9a-f]{40}", args.expected_head) is not None
            and head == args.expected_head, "checkout is not expected PR head")
    if args.event == "workflow_dispatch":
        changed = []
        affected = True
    else:
        require(re.fullmatch(r"[0-9a-f]{40}", args.base) is not None,
                "PR base SHA is missing or malformed")
        merge_base = git(root, "merge-base", args.base, head).decode("ascii").strip()
        require(re.fullmatch(r"[0-9a-f]{40}", merge_base) is not None,
                "merge base missing")
        changed_bytes = git(root, "diff", "--no-renames", "--name-only",
                            "-z", merge_base, head)
        changed = [p.decode("utf-8", errors="surrogateescape")
                   for p in changed_bytes.split(b"\0") if p]
        # Go's canonical ICE builder rejects symlinks anywhere in traversed source.
        # Any symlink mode change must therefore run the full validation even if
        # its filename would not otherwise be indexed.
        raw = git(root, "diff", "--no-renames", "--raw", "-z",
                  merge_base, head)
        symlink_change = b"120000" in raw
        affected = any(relevant(path) for path in changed) or symlink_change
    require(args.output is not None, "GITHUB_OUTPUT path required")
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("a", encoding="utf-8") as handle:
        handle.write("atlas_relevant=" + str(affected).lower() + "\n")
    print(json.dumps({
        "result": "ATLAS_CI_SCOPE_PASS",
        "head": head,
        "event": args.event,
        "changed_count": len(changed),
        "full_qualification_required": affected,
    }, sort_keys=True))


if __name__ == "__main__":
    main()
