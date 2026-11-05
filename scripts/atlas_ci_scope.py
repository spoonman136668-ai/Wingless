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
    "scripts/maintain_code_atlas.py",
    "scripts/test_source_aware_atlas.py",
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
    pr_trigger = header.split("  pull_request:\n", 1)[1].split("  workflow_dispatch:", 1)[0]
    require(re.search(r"(?m)^\s+paths(?:-ignore)?:", pr_trigger) is None,
            "PR-level path filtering can strand required status")
    require("    name: " + STATUS_JOB in workflow, "required job identity changed")
    require("        id: scope\n" in workflow,
            "atlas scope classification step absent")
    require("  push:\n    branches: [main]\n" in workflow,
            "main-only cache producer trigger missing")
    require(workflow.count("runs-on: windows-2022") == 2,
            "Atlas producer and original qualification must both run on Windows 2022")
    require("  qualify:\n    if: github.event_name != 'push'" in workflow,
            "required PR qualification separated from main cache seed")
    require("  seed_cache:\n    if: github.event_name == 'push'" in workflow,
            "cache producer may run on an untrusted PR")
    require(workflow.count("id: cache_identity") == 2,
            "compiler identity must be checked in producer and consumer")
    require(workflow.count("ATLAS_CACHE_IMAGE_IDENTITY_MISSING") == 2,
            "image identity verification missing")
    require(workflow.count("ATLAS_CACHE_TOOLCHAIN_MISMATCH") == 2,
            "toolchain identity verification missing")
    require(workflow.count("wingless-atlas-win2022-go1255-") == 2,
            "exact cache keys not bound in both jobs")
    require(workflow.count("actions/cache/restore@v4") == 1
            and workflow.count("actions/cache@v4") == 1,
            "PR must restore only; main must seed")
    require("restore-keys:" not in workflow,
            "inexact cache fallback prohibited")
    require("      - name: One complete integration regression\n" in workflow
            and "          go test ./..." in workflow,
            "full original Go regression changed")
    for step in (
        "Set up frozen Go toolchain",
        "Validate source-index freshness and typed relationships",
        "Bounded parallel tamper-negative and genuine Chromium acceptance",
        "Focused ICE regression",
        "One complete integration regression",
        "Vet, build, and verify clean checkout",
    ):
        marker = "      - name: " + step + "\n"
        require(workflow.count(marker) == 1, "missing/duplicate step: " + step)
        block = workflow.split(marker, 1)[1].split("\n      - name:", 1)[0]
        require("if: steps.scope.outputs.atlas_relevant == 'true'" in block,
                "full qualification can bypass scope: " + step)
    parallel_marker = "      - name: Bounded parallel tamper-negative and genuine Chromium acceptance\n"
    parallel_block = workflow.split(parallel_marker, 1)[1].split("\n      - name:", 1)[0]
    require(parallel_block.count("$jobs += Start-ThreadJob") == 2,
            "both mandatory acceptance jobs must start exactly once")
    for required in (
        "python scripts/test_code_atlas_tamper.py",
        "python -B scripts/test_atlas_browser.py",
        "PINNED_SELENIUM_INSTALL_FAILED",
        "ATLAS_TAMPER_NEGATIVE_ACCEPTANCE_FAILED",
        "ATLAS_CHROMIUM_BROWSER_ACCEPTANCE_FAILED",
        "ATLAS_TAMPER_PASS_RECEIPT_MISSING",
        "ATLAS_CHROMIUM_PASS_RECEIPT_MISSING",
        "WINGLESS_ATLAS_NEGATIVE_QUALIFICATION_PASS",
        "WINGLESS_ATLAS_CHROMIUM_BROWSER_ACCEPTANCE_PASS",
        "$jobs | Wait-Job",
        "Receive-Job -Job $job -ErrorAction Stop",
        "if ($job.State -ne 'Completed')",
        "ATLAS_PARALLEL_ACCEPTANCE_FAILED",
        "ATLAS_BOUNDED_PARALLEL_ACCEPTANCE_PASS",
    ):
        require(required in parallel_block, "parallel acceptance contract changed: " + required)
    require(parallel_block.count("ConvertFrom-Json -ErrorAction Stop") == 2,
            "both workers must parse structured original receipts")
    for required in (
        "ATLAS_TAMPER_PASS_RECEIPT_IDENTITY_MISMATCH",
        "ATLAS_CHROMIUM_PASS_RECEIPT_IDENTITY_MISMATCH",
        "$receipt.head -ne $expectedHead",
        "$receipt.expected_rejections -ne 9",
        "$receipt.positive_control -ne 'PASS'",
        "$receipt.source_checkout_changed -ne $false",
        "$receipt.browser_test_commit -ne $expectedHead",
        "$receipt.original_checkout_dirty -ne $false",
        "$receipt.steps -notcontains 'actual_click_cross_file_relationship'",
        "$receipt.steps -notcontains 'actual_click_source_link_new_tab'",
        "$receipt.browser -notlike '*real Selenium UI*'",
        "$receipts.Count -ne 1",
    ):
        require(required in parallel_block, "parallel receipt identity contract changed: " + required)
    entries = json.loads((root / ".ice/manifest.json").read_text("utf-8"))["files"]
    representative = [
        "new_source.go", "src/new_source.go", "README.md", "src/new.md",
        "src/new.json", "src/new.ps1", ".github/workflows/fresh.yml",
        "src/new.yaml", "go.mod", ".gitattributes", ".gitignore",
        "scripts/validate_code_atlas.py", "scripts/test_code_atlas_tamper.py",
        "scripts/test_atlas_browser.py", "scripts/atlas_ci_scope.py",
        "scripts/maintain_code_atlas.py", "scripts/test_source_aware_atlas.py",
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
