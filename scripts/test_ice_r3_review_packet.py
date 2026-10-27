#!/usr/bin/env python3
"""R3 ICE source review packet: exact commit and hidden-worktree negative test.

Runs only inside a disposable branch checkout with authenticated historical
objects. All mutation occurs in a SECOND, disposable Git worktree.
"""
import argparse
import json
import subprocess
import tempfile
from pathlib import Path

from ice_r3_review_packet import build

TARGET = "research/native/itd02/r3/transitions.go"
REVIEWED = 352
NEW = 96


def run(root, *args, ok=True):
    p = subprocess.run(list(args), cwd=root, stdout=subprocess.PIPE,
                       stderr=subprocess.STDOUT, text=True, timeout=180)
    if ok and p.returncode != 0:
        raise RuntimeError("R3_ICE_TEST_COMMAND_FAILED: " + str(args) +
                           " " + p.stdout[-1000:])
    return p


def verify(root, scratch_parent):
    root = root.resolve()
    result = build(root)
    assert result["previous_reviewed_symbols"] == REVIEWED
    assert result["added_declaration_count"] == NEW
    assert result["prior_reviewed_identity_reuse"]["unchanged_declaration"] == 351
    assert result["prior_reviewed_identity_reuse"]["requires_fresh_review"] == 1
    assert result["historical_pin_requires_independent_reanchor"] is True
    assert result["semantic_atlas_publication_authorized"] is False
    assert len(result["unreviewed_declarations"]) == NEW
    assert all(x["review_status"] == "UNREVIEWED_NOT_AUTHORIZED"
               and x["authenticated_reviewer"] is None
               and x["semantic_behavior_claim"] is None
               and x["source_proof_approved"] is False
               for x in result["unreviewed_declarations"])
    print("R3_ICE_FORGED_OR_AUTOMATIC_REVIEW_APPROVAL=REJECTED")

    with tempfile.TemporaryDirectory(prefix="wingless-r3-ice-negative-",
                                     dir=scratch_parent) as directory:
        clone = Path(directory) / "source"
        run(root, "git", "worktree", "add", "--quiet", "--detach",
            str(clone), "HEAD")
        try:
            run(clone, "go", "run", "./cmd/ice", "build")
            assert build(clone)["added_declaration_count"] == NEW
            run(clone, "git", "update-index", "--assume-unchanged", TARGET)
            path = clone / TARGET
            with path.open("a", encoding="utf-8") as f:
                f.write("\n// R3 synthetic hidden source divergence; discard this worktree.\n")
            run(clone, "go", "run", "./cmd/ice", "build")
            dirty = run(clone, "git", "status", "--porcelain", "--untracked-files=all")
            if any(TARGET in line for line in dirty.stdout.splitlines()):
                raise AssertionError("R3_ICE_HIDDEN_SOURCE_FIXTURE_WAS_NOT_HIDDEN")
            try:
                build(clone)
            except SystemExit as e:
                if "HEAD_COMMITTED_SOURCE_MISMATCH:" + TARGET not in str(e):
                    raise
            else:
                raise AssertionError("R3_ICE_HIDDEN_SOURCE_TAMPER_ACCEPTED")
            print("R3_ICE_HIDDEN_WORKTREE_SOURCE_DRIFT=REJECTED")
        finally:
            run(root, "git", "worktree", "remove", "--force", str(clone))
    assert not run(root, "git", "status", "--porcelain",
                   "--untracked-files=all").stdout.splitlines() or all(
        line[3:].startswith(".ice/") for line in run(
            root, "git", "status", "--porcelain",
            "--untracked-files=all").stdout.splitlines())
    print("R3_ICE_UNREVIEWED_PACKET_NEGATIVE_ACCEPTANCE=PASS")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--scratch-parent", type=Path, required=True)
    a = parser.parse_args()
    verify(a.root, a.scratch_parent)
