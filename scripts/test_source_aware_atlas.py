#!/usr/bin/env python3
"""G26 isolated source-aware atlas acceptance and fail-closed negatives.

All modifications happen in disposable exact-head clones, never in checkout.
"""
import argparse
import hashlib
import os
from concurrent.futures import ThreadPoolExecutor
import json
import math
import platform
import statistics
import subprocess
import sys
import tempfile
import time
from pathlib import Path

SOURCE = Path("inference/lifecycle.go")
ATLAS = Path(".ice/code-intelligence.json")


def require(ok, reason):
    if not ok:
        raise SystemExit("G26_TEST_FAIL: " + reason)


def cmd(root, *argv):
    p = subprocess.run([str(a) for a in argv], cwd=root,
                       stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                       text=True, timeout=240)
    return p.returncode, p.stdout


def good(root, *argv):
    code, out = cmd(root, *argv)
    require(code == 0, "command failed " + repr(argv) + ": " + out[-1600:])
    return out


def reject(root, expected, *argv):
    code, out = cmd(root, *argv)
    require(code != 0 and expected in out,
            "rejection mismatch " + repr(argv) + ": " + out[-1600:])
    print("G26_EXPECTED_REJECTION:", expected)


def clone(root, folder, commit):
    good(root, "git", "clone", "--quiet", "--no-hardlinks", "--no-checkout",
         root, folder)
    good(folder, "git", "checkout", "--quiet", "--detach", commit)
    require(not good(folder, "git", "status", "--porcelain").strip(),
            "scratch checkout dirty at start")


def make_commit(folder, old, new):
    require(old != new, "no fixture mutation")
    SOURCE_PATH = folder / SOURCE
    require(SOURCE_PATH.read_text("utf-8") == old, "fixture source drift")
    SOURCE_PATH.write_text(new, "utf-8")
    good(folder, "git", "add", str(SOURCE))
    good(folder, "git", "-c", "user.name=Wingless G26 Synthetic",
         "-c", "user.email=g26-fixture@example.invalid",
         "commit", "-qm", "G26 synthetic source first-use fixture")
    return good(folder, "git", "rev-parse", "HEAD").strip()


def fixture_unchanged(root, t, baseline):
    """Run unchanged, fail-closed provenance controls in an isolated clone."""
    last_mark = time.perf_counter()
    timings = {}

    def mark(stage):
        nonlocal last_mark
        now = time.perf_counter()
        timings[stage] = round(now - last_mark, 6)
        last_mark = now

    first = Path(t) / "unchanged-declarations"
    clone(root, first, baseline)
    original = (first / SOURCE).read_text("utf-8")
    comment = original + "\n// G26 synthetic file-only delta; declaration bytes unchanged.\n"
    comment_head = make_commit(first, original, comment)
    reject(first, "ICE_STALE", "go", "run", "./cmd/ice", "validate")
    good(first, "go", "run", "./cmd/ice", "build")
    good(first, "go", "run", "./cmd/ice", "validate")
    report = good(first, sys.executable, "-B",
                  "scripts/maintain_code_atlas.py", "--root", ".")
    q = json.loads(report.strip().splitlines()[-1])
    require(q["reused_declarations"] == 352 and
            q["fresh_review_receipts"] == 0 and
            q["target_source_commit"] == comment_head,
            "source-identical declarations did not preserve exact witnesses")
    derived = json.loads((first / ATLAS).read_text("utf-8"))
    require(derived["source_commit"] == comment_head, "new source identity missing")
    node = next(n for n in derived["nodes"]
                if n["file"] == str(SOURCE).replace("\\", "/"))
    require(node["provenance_lineage"]["mode"] ==
            "canonical-LF-declaration-byte-reuse" and
            node["review_evidence"]["source_commit"] == comment_head,
            "body-identical proof was not rebound")
    rels = Path(t) / "types-one.json"
    with rels.open("w", encoding="utf-8") as out:
        p = subprocess.run(["go", "run", "./cmd/ice", "relationships"],
                           cwd=first, stdout=out, stderr=subprocess.PIPE,
                           text=True, timeout=240)
        require(p.returncode == 0, "typed extraction failed: " + p.stderr)
    good(first, sys.executable, "-B", "scripts/validate_code_atlas.py",
         "--root", ".", "--relationships", rels)
    print("G26_EXACT_DECLARATION_REUSE=PASS", comment_head)
    mark("unchanged_declaration_reuse")

    return timings


def fixture_guardrails(root, t, baseline):
    """Run unchanged, fail-closed provenance controls in an isolated clone."""
    last_mark = time.perf_counter()
    timings = {}

    def mark(stage):
        nonlocal last_mark
        now = time.perf_counter()
        timings[stage] = round(now - last_mark, 6)
        last_mark = now

    # An uncommitted Go change plus a rebuilt canonical ICE index must
    # never masquerade as the frozen target commit.
    dirty = Path(t) / "dirty-source-checkout"
    clone(root, dirty, baseline)
    unchanged = (dirty / SOURCE).read_text("utf-8")
    (dirty / SOURCE).write_text(
        unchanged + "\n// G27 uncommitted source delta.\n", "utf-8")
    good(dirty, "go", "run", "./cmd/ice", "build")
    reject(dirty, "G27_DIRTY_NON_ICE_WORKTREE",
           sys.executable, "-B", "scripts/maintain_code_atlas.py",
           "--root", ".")
    mark("dirty_worktree_rejection")
    # Even if local Git status is deliberately blinded, independent
    # HEAD-blob verification must reject a modified worktree file.
    hidden = Path(t) / "hidden-working-tree-change"
    clone(root, hidden, baseline)
    hidden_old = (hidden / SOURCE).read_text("utf-8")
    good(hidden, "git", "update-index", "--assume-unchanged",
         SOURCE.as_posix())
    (hidden / SOURCE).write_text(
        hidden_old + "\n// G27 hidden source blob divergence.\n",
        "utf-8")
    good(hidden, "go", "run", "./cmd/ice", "build")
    reject(hidden, "G27_TARGET_SOURCE_CONTENT_MISMATCH",
           sys.executable, "-B", "scripts/maintain_code_atlas.py",
           "--root", ".")
    mark("hidden_source_blob_rejection")
    # A committed module context change must not inherit 352 behavioral
    # descriptions solely because Go declaration bodies match.
    module = Path(t) / "module-context-change"
    clone(root, module, baseline)
    modfile = Path("go.mod")
    mod_before = (module / modfile).read_text("utf-8")
    (module / modfile).write_text(
        mod_before + "\n// G27 synthetic module context drift\n", "utf-8")
    good(module, "git", "add", "go.mod")
    good(module, "git", "-c", "user.name=Wingless G27 Synthetic",
         "-c", "user.email=g27-fixture@example.invalid",
         "commit", "-qm", "G27 synthetic module-context fixture")
    good(module, "go", "run", "./cmd/ice", "build")
    reject(module, "G27_DEPENDENCY_CONTEXT_CHANGED",
           sys.executable, "-B", "scripts/maintain_code_atlas.py",
           "--root", ".")
    print("G27_GIT_SOURCE_IDENTITY_AND_MODULE_CONTEXT_NEGATIVES_PASS")
    mark("module_context_rejection")
    return timings


def fixture_changed(root, t, baseline):
    """Run unchanged, fail-closed provenance controls in an isolated clone."""
    last_mark = time.perf_counter()
    timings = {}

    def mark(stage):
        nonlocal last_mark
        now = time.perf_counter()
        timings[stage] = round(now - last_mark, 6)
        last_mark = now

    second = Path(t) / "changed-declaration"
    clone(root, second, baseline)
    old = (second / SOURCE).read_text("utf-8")
    before = "invalid lifecycle transition %s -> %s"
    after = "invalid G26 lifecycle transition %s -> %s"
    require(old.count(before) == 1, "fixture declaration anchor changed")
    new_head = make_commit(second, old, old.replace(before, after))
    good(second, "go", "run", "./cmd/ice", "build")
    reject(second, "G26_REVIEW_REQUIRED", sys.executable, "-B",
           "scripts/maintain_code_atlas.py", "--root", ".")
    require(json.loads((second / ATLAS).read_text("utf-8"))["source_commit"]
            != new_head, "failed generation unexpectedly published atlas")
    manifest = json.loads((second / ".ice/manifest.json").read_text("utf-8"))
    old_atlas = json.loads((second / ATLAS).read_text("utf-8"))
    entry = next(f for f in manifest["files"] if f["path"] == SOURCE.as_posix())
    item = next(n for n in old_atlas["nodes"]
                if n["file"] == SOURCE.as_posix() and n["name"] == "Move")
    new_source = (second / SOURCE).read_bytes()
    lines = new_source.splitlines(keepends=True)
    text = b"".join(lines[item["start_line"]-1:item["end_line"]])
    # Comment/body fixture preserves original declaration positions.
    fixture = [{
        "key": [item["file"], item["kind"], item["name"],
                item.get("receiver", "")],
        "target_source_commit": new_head,
        "file_sha256": entry["sha256"],
        "declaration_sha256": hashlib.sha256(text).hexdigest(),
        "start_line": item["start_line"],
        "end_line": item["end_line"],
        "reviewer": "G26 SYNTHETIC FIXTURE ONLY",
        "label": "Review new lifecycle transition rejection path",
        "description": "Rejects transitions missing from the existing explicit lifecycle allowlist; synthetic review fixture only.",
        "literal_tokens": [after],
    }]
    receipt = Path(t) / "review-fixture.json"
    receipt.write_text(json.dumps(fixture), "utf-8")
    fixture[0]["declaration_sha256"] = "0" * 64
    invalid = Path(t) / "invalid-review.json"
    invalid.write_text(json.dumps(fixture), "utf-8")
    reject(second, "fresh-review receipt source identity drift",
           sys.executable, "-B", "scripts/maintain_code_atlas.py",
           "--root", ".", "--receipts", invalid)
    positive = good(second, sys.executable, "-B",
                    "scripts/maintain_code_atlas.py",
                    "--root", ".", "--receipts", receipt)
    out = json.loads(positive.strip().splitlines()[-1])
    require(out["fresh_review_receipts"] == 1 and
            out["reused_declarations"] == 351 and
            out["independent_approval"] is False,
            "changed declaration incorrectly reused or approved")
    report = Path(t) / "types-two.json"
    with report.open("w", encoding="utf-8") as out:
        p = subprocess.run(["go", "run", "./cmd/ice", "relationships"],
                           cwd=second, stdout=out, stderr=subprocess.PIPE,
                           text=True, timeout=240)
        require(p.returncode == 0, "typed extraction failed: " + p.stderr)
    reject(second, "G26_UNAPPROVED_FRESH_REVIEW",
           sys.executable, "-B", "scripts/validate_code_atlas.py",
           "--root", ".", "--relationships", report)
    print("G26_CHANGED_DECLARATION_FRESH_RECEIPT_REQUIRES_INDEPENDENT_APPROVAL=PASS", new_head)
    mark("changed_declaration_fresh_review")

    return timings


def fixture_new_symbol(root, t, baseline):
    """Run unchanged, fail-closed provenance controls in an isolated clone."""
    last_mark = time.perf_counter()
    timings = {}

    def mark(stage):
        nonlocal last_mark
        now = time.perf_counter()
        timings[stage] = round(now - last_mark, 6)
        last_mark = now

    # G28: an additional Go declaration must be source-witnessed,
    # remain explicitly unreviewed, and never inherit the 352 old reviews.
    added = Path(t) / "provisional-new-symbol"
    clone(root, added, baseline)
    existing = (added / SOURCE).read_text("utf-8")
    with_new = (existing +
                "\nfunc G28AtlasSourceWitnessFixture() int { return 1 }\n")
    added_head = make_commit(added, existing, with_new)
    good(added, "go", "run", "./cmd/ice", "build")
    reject(added, "G28_NEW_SYMBOL_SOURCE_RECEIPTS_REQUIRED",
           sys.executable, "-B", "scripts/maintain_code_atlas.py",
           "--root", ".")
    idx = json.loads((added / ".ice/symbol-map.json").read_text("utf-8"))
    fresh_node = next(n for n in idx
                      if n["file"] == SOURCE.as_posix()
                      and n["name"] == "G28AtlasSourceWitnessFixture")
    raw = (added / SOURCE).read_bytes()
    decl = b"".join(raw.splitlines(keepends=True)[
        fresh_node["start"]-1:fresh_node["end"]])
    review = [{
        "key": [SOURCE.as_posix(), fresh_node["kind"], fresh_node["name"], ""],
        "target_source_commit": added_head,
        "file_sha256": hashlib.sha256(raw).hexdigest(),
        "declaration_sha256": hashlib.sha256(decl).hexdigest(),
        "start_line": fresh_node["start"], "end_line": fresh_node["end"],
        "reviewer": "G28 SYNTHETIC UNVERIFIED",
        "label": "Synthetic source-witnessed new symbol",
        "description": "Synthetic test-only lexical definition with no independently verified behavior",
        "literal_tokens": ["G28AtlasSourceWitnessFixture"],
    }]
    review_file = Path(t) / "g28-added.json"
    review_file.write_text(json.dumps(review), "utf-8")
    review[0]["file_sha256"] = "0" * 64
    forged_file = Path(t) / "g28-forged.json"
    forged_file.write_text(json.dumps(review), "utf-8")
    reject(added, "G28_NEW_SYMBOL_RECEIPT_IDENTITY_DRIFT",
           sys.executable, "-B", "scripts/maintain_code_atlas.py",
           "--root", ".", "--receipts", forged_file)
    accepted = good(added, sys.executable, "-B",
                    "scripts/maintain_code_atlas.py",
                    "--root", ".", "--receipts", review_file)
    result = json.loads(accepted.strip().splitlines()[-1])
    require(result["new_source_bound_provisional"] == 1
            and result["reused_declarations"] == 352,
            "G28 new source-bound symbol result incorrect")
    expanded = json.loads((added / ATLAS).read_text("utf-8"))
    marker = next(n for n in expanded["nodes"]
                  if n["name"] == "G28AtlasSourceWitnessFixture")
    require(len(expanded["nodes"]) == 353
            and marker["evidence_basis"] == "lexical-provisional"
            and marker["review_evidence"] is None
            and marker["provenance_lineage"]["independent_approval"] is False,
            "G28 provisional additions were silently promoted to reviewed")
    g28_types = Path(t) / "types-g28.json"
    with g28_types.open("w", encoding="utf-8") as out:
        p = subprocess.run(["go", "run", "./cmd/ice", "relationships"],
                           cwd=added, stdout=out, stderr=subprocess.PIPE,
                           text=True, timeout=240)
        require(p.returncode == 0, "G28 typed extraction failed: " + p.stderr)
    good(added, sys.executable, "-B", "scripts/validate_code_atlas.py",
         "--root", ".", "--relationships", g28_types)
    print("G28_NEW_SOURCE_BOUND_PROVISIONAL=PASS", added_head)
    marker["evidence_basis"] = "source-reviewed"
    (added / ATLAS).write_text(json.dumps(expanded), "utf-8")
    reject(added, "source-reviewed symbol has no explicit proof",
           sys.executable, "-B", "scripts/validate_code_atlas.py",
           "--root", ".", "--relationships", g28_types)
    print("G28_UNAUTHENTICATED_REVIEW_PROMOTION_REJECTED=PASS")
    mark("new_symbol_provisional_and_promotion_negative")

    return timings



def verify_canonical_derived_diff_contract():
    """Allow only the derived Go-index paths; forbid unrelated mutations.

    The independent go run ./cmd/ice validate compares each derived file
    byte-for-byte with canonical source. This test protects the narrow
    filename allowlist, not a claim of independent semantic approval.
    """
    from validate_code_atlas import ALLOWED_DIFF
    derived = {
        ".ice/dependency-map.json",
        ".ice/symbol-map.json",
        ".ice/test-map.json",
    }
    require(derived.issubset(ALLOWED_DIFF),
            "G26_CANONICAL_DERIVED_SOURCE_INDEX_NOT_ADMITTED")
    forbidden = {
        "research/jepa/stream.go",
        ".github/workflows/native.yml",
        ".ice/unknown-generated.json",
        "go.mod",
        "research/status/accepted-science.json",
    }
    require(not (forbidden & ALLOWED_DIFF),
            "G26_UNRELATED_SOURCE_OR_RESEARCH_SCOPE_ADMITTED")


def run(root, scratch_parent, *, workers_override=None):
    started = time.perf_counter()
    timings = {}
    root = root.resolve()
    verify_canonical_derived_diff_contract()
    baseline = good(root, "git", "rev-parse", "HEAD").strip()
    require(not good(root, "git", "status", "--porcelain").strip(),
            "real checkout dirty before testing")
    # Only independent disposable exact-head clones run concurrently.
    # Cap Go subprocess contention; all futures must succeed before acceptance.
    max_workers = min(2, os.cpu_count() or 1) if workers_override is None else workers_override
    require(max_workers in (1, 2), "fixture worker limit must be 1 or 2")
    with tempfile.TemporaryDirectory(prefix="g26-atlas-", dir=scratch_parent) as t:
        groups = (
            ("new_symbol", fixture_new_symbol),
            ("unchanged", fixture_unchanged),
            ("changed", fixture_changed),
            ("guardrails", fixture_guardrails),
        )
        required_stages = {
            "unchanged": {"unchanged_declaration_reuse"},
            "guardrails": {"dirty_worktree_rejection",
                           "hidden_source_blob_rejection", "module_context_rejection"},
            "changed": {"changed_declaration_fresh_review"},
            "new_symbol": {"new_symbol_provisional_and_promotion_negative"},
        }
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = {
                name: executor.submit(task, root, Path(t), baseline)
                for name, task in groups
            }
            # Join every lane; one missing result or exception fails qualification.
            for name in ("unchanged", "guardrails", "changed", "new_symbol"):
                stage_timings = futures[name].result()
                require(set(stage_timings) == required_stages[name],
                        "missing/duplicate fixture evidence: " + name)
                require(not set(timings).intersection(stage_timings),
                        "duplicate fixture timing evidence: " + name)
                timings.update(stage_timings)
        clean_started = time.perf_counter()
        require(not good(root, "git", "status", "--porcelain").strip(),
                "original checkout modified")
        timings["final_source_checkout_cleanliness"] = round(
            time.perf_counter() - clean_started, 6)
    receipt = {
        "timing_schema": "wingless.atlas.g26-source-aware-diagnostic.v1",
        "timings_seconds": timings,
        "parallel_execution": max_workers > 1,
        "max_workers": max_workers,
        "total_elapsed_seconds": round(time.perf_counter() - started, 6),
        "result": "G26_ISOLATED_MAINTENANCE_QUALIFICATION_PASS",
        "reused_unchanged": 352,
        "changed_declaration_without_receipt": "REJECTED",
        "forged_fresh_receipt": "REJECTED",
        "synthetic_fresh_receipt": "ACCEPTED_AS_UNVERIFIED_ASSERTION",
        "dirty_source_identity": "REJECTED",
        "hidden_source_blob_divergence": "REJECTED",
        "module_context_drift": "REJECTED",
        "original_checkout_changed": False,
    }
    print(json.dumps(receipt, sort_keys=True))
    return receipt


def compare_workers(root, scratch_parent):
    """Optional development-only counterbalanced, same-root 1-vs-2 worker probe.

    Run *all* original provenance fixtures on every arm. This is not an
    acceptance-speedup claim: OS page cache and Go build caches can still vary.
    """
    root = root.resolve()
    baseline = good(root, "git", "rev-parse", "HEAD").strip()
    require(not good(root, "git", "status", "--porcelain").strip(),
            "original checkout dirty before worker comparison")
    require((os.cpu_count() or 1) >= 2,
            "comparison requires two logical CPUs; do not oversubscribe")
    expected_stages = {
        "unchanged_declaration_reuse", "dirty_worktree_rejection",
        "hidden_source_blob_rejection", "module_context_rejection",
        "changed_declaration_fresh_review",
        "new_symbol_provisional_and_promotion_negative",
        "final_source_checkout_cleanliness",
    }
    order = (1, 2, 2, 1)
    samples = {1: [], 2: []}
    for workers in order:
        require(good(root, "git", "rev-parse", "HEAD").strip() == baseline,
                "source HEAD changed during worker comparison")
        receipt = run(root, scratch_parent, workers_override=workers)
        require(receipt["result"] == "G26_ISOLATED_MAINTENANCE_QUALIFICATION_PASS"
                and receipt["max_workers"] == workers
                and receipt["parallel_execution"] == (workers == 2)
                and set(receipt["timings_seconds"]) == expected_stages,
                "incomplete or inconsistent worker comparison fixture receipt")
        seconds = receipt["total_elapsed_seconds"]
        require(math.isfinite(seconds) and seconds > 0,
                "invalid worker elapsed time")
        samples[workers].append(seconds)
        require(not good(root, "git", "status", "--porcelain").strip(),
                "worker comparison modified original source checkout")
    serial_median = statistics.median(samples[1])
    parallel_median = statistics.median(samples[2])
    print(json.dumps({
        "result": "G26_WORKER_COMPARISON_DIAGNOSTIC_PASS",
        "schema": "wingless.atlas.worker-comparison.v1",
        "source_head": baseline,
        "platform": platform.platform(),
        "python_version": platform.python_version(),
        "go_version": good(root, "go", "version").strip(),
        "logical_cpu_count": os.cpu_count(),
        "worker_order": order,
        "serial_seconds": samples[1],
        "parallel_seconds": samples[2],
        "serial_median_seconds": serial_median,
        "parallel_median_seconds": parallel_median,
        "observed_median_speedup_ratio": round(serial_median / parallel_median, 6),
        "qualification": "DEVELOPMENT_DIAGNOSTIC_ONLY",
        "limitations": "OS/page caches, Go compile caches and runner contention may vary",
        "original_checkout_changed": False,
    }, sort_keys=True))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--scratch-parent", type=Path, default=None)
    parser.add_argument("--workers", type=int, choices=(1, 2), default=None,
                        help="Optional 1/2-worker development probe; default retains qualified 2-worker cap")
    parser.add_argument("--compare-workers", action="store_true",
                        help="Opt-in ABBA same-source original G26/G27/G28 full-suite comparison")
    args = parser.parse_args()
    require(not (args.compare_workers and args.workers is not None),
            "comparison defines its own 1/2-worker schedule")
    if args.compare_workers:
        compare_workers(args.root, args.scratch_parent)
    else:
        run(args.root, args.scratch_parent, workers_override=args.workers)
