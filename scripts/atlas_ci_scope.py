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
    require(re.search(r"(?m)^\s+paths(?:-ignore)?:", header) is None,
            "PR-level path filtering can strand required status")
    # The always-on protected status is an aggregate: no independent lane may
    # grant acceptance, and skipped/cancelled/failed lanes must fail closed.
    def contract(w):
        require(w.count("    name: " + STATUS_JOB + "\n") == 1,
                "missing/duplicate required aggregate status")
        require("        id: scope\n" in w, "atlas scope classifier absent")
        job_parts = re.split(r"(?m)^  (scope|browser|tamper|integration|qualification):\n", w)
        jobs = {job_parts[i]: job_parts[i + 1]
                for i in range(1, len(job_parts) - 1, 2)}
        require(set(jobs) == {"scope", "browser", "tamper", "integration", "qualification"},
                "atlas scope or parallel lanes missing")
        require("      atlas_relevant: ${{ steps.scope.outputs.atlas_relevant }}" in jobs["scope"],
                "scope result not forwarded")
        require("scripts/atlas_ci_scope.py --root . --self-test" in jobs["scope"]
                and "--expected-head" in jobs["scope"]
                and "--output" in jobs["scope"],
                "unconditional fail-closed source/scope classification missing")
        per_lane = {
            "browser": ["Genuine Chromium browser navigation acceptance"],
            "tamper": ["Set up frozen Go toolchain",
                       "Validate source-index freshness and typed relationships",
                       "Tamper-resistance negative acceptance in disposable exact-head clone"],
            "integration": ["Set up frozen Go toolchain", "Focused ICE regression",
                            "One complete integration regression",
                            "Vet, build, and verify clean checkout"],
        }
        for lane, steps in per_lane.items():
            block = jobs[lane]
            require("    needs: [scope]\n" in block
                    and "    if: needs.scope.outputs.atlas_relevant == 'true'\n" in block,
                    "lane bypasses positive source classification: " + lane)
            require("ref: ${{ github.event.pull_request.head.sha || github.sha }}" in block
                    and "persist-credentials: false" in block
                    and "fetch-depth: 0" in block
                    and "ATLAS_EXACT_HEAD_DRIFT" in block
                    and "ATLAS_SOURCE_CHECKOUT_DIRTY" in block,
                    "lane source checkout custody weakened: " + lane)
            for step in steps:
                require(block.count("      - name: " + step + "\n") == 1,
                        "missing/duplicate native qualification step: " + step)
        summary = jobs["qualification"]
        require("    needs: [scope, browser, tamper, integration]\n" in summary
                and "    if: always()\n" in summary,
                "aggregate not fail-closed across every lane")
        for field in ("SCOPE_RESULT", "ATLAS_RELEVANT", "BROWSER_RESULT",
                      "TAMPER_RESULT", "INTEGRATION_RESULT"):
            require(field + ":" in summary, "aggregate missing result: " + field)
        for expected in ('test "$SCOPE_RESULT" = success',
                         'test "$ATLAS_RELEVANT" = false',
                         'test "$ATLAS_RELEVANT" = true',
                         'test "$BROWSER_RESULT" = skipped',
                         'test "$TAMPER_RESULT" = skipped',
                         'test "$INTEGRATION_RESULT" = skipped',
                         'test "$BROWSER_RESULT" = success',
                         'test "$TAMPER_RESULT" = success',
                         'test "$INTEGRATION_RESULT" = success',
                         "ATLAS_SCOPE_INVALID"):
            require(expected in summary, "aggregate condition removed: " + expected)

    contract(workflow)
    # Deterministic negative controls: never accept deletion of a required lane,
    # conditional omission of the aggregate, or silent weakening of status gates.
    negatives = (
        ("    if: always()\n", "    if: false\n"),
        ('test "$BROWSER_RESULT" = success', 'test "$BROWSER_RESULT" != failure'),
        ("      - name: One complete integration regression\n",
         "      - name: Optional integration regression\n"),
        ("    if: needs.scope.outputs.atlas_relevant == 'true'\n",
         "    if: always()\n"),
        ("    needs: [scope, browser, tamper, integration]\n",
         "    needs: [scope, browser, integration]\n"),
        ("      atlas_relevant: ${{ steps.scope.outputs.atlas_relevant }}",
         "      atlas_relevant: 'false'"),
    )
    rejected = 0
    for before, after in negatives:
        require(before in workflow, "negative mutation anchor absent")
        mutated = workflow.replace(before, after, 1)
        try:
            contract(mutated)
        except SystemExit:
            rejected += 1
        else:
            require(False, "negative contract mutation accepted")
    require(rejected == len(negatives), "not all negatives rejected")
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
        "aggregate_negative_controls": rejected,
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
