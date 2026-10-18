#!/usr/bin/env python3
"""G26 isolated source-aware atlas acceptance and fail-closed negatives.

All modifications happen in disposable exact-head clones, never in checkout.
"""
import argparse
import hashlib
import json
import subprocess
import sys
import tempfile
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


def run(root, scratch_parent):
    root = root.resolve()
    baseline = good(root, "git", "rev-parse", "HEAD").strip()
    require(not good(root, "git", "status", "--porcelain").strip(),
            "real checkout dirty before testing")
    with tempfile.TemporaryDirectory(prefix="g26-atlas-", dir=scratch_parent) as t:
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
        good(second, sys.executable, "-B", "scripts/validate_code_atlas.py",
             "--root", ".", "--relationships", report)
        print("G26_CHANGED_DECLARATION_FRESH_RECEIPT=PASS", new_head)

        require(not good(root, "git", "status", "--porcelain").strip(),
                "original checkout modified")
    print(json.dumps({
        "result": "G26_ISOLATED_MAINTENANCE_QUALIFICATION_PASS",
        "reused_unchanged": 352,
        "changed_declaration_without_receipt": "REJECTED",
        "forged_fresh_receipt": "REJECTED",
        "synthetic_fresh_receipt": "ACCEPTED_AS_UNVERIFIED_ASSERTION",
        "original_checkout_changed": False,
    }, sort_keys=True))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--scratch-parent", type=Path, default=None)
    args = parser.parse_args()
    run(args.root, args.scratch_parent)
