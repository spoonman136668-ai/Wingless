#!/usr/bin/env python3
"""Read-only, exact-HEAD source review queue for Wingless R3 ICE expansion.

This produces *review material*, never provenance approval, source-backed
semantic labels, or an ICE promotion. Canonical ICE is rebuilt separately in a
disposable worktree. The reviewed 352-symbol atlas remains untouched.
"""
import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

OLD_FILES = 138
OLD_DECLARATIONS = 352
EXPECTED_ADDED_FILES = 16
EXPECTED_ADDED_DECLARATIONS = 96
EXPECTED_R3_DECLARATIONS = 33
REPO = "spoonman136668-ai/Wingless"
R3_PREFIX = "research/native/itd02/r3/"


def require(ok, reason):
    if not ok:
        raise SystemExit("R3_ICE_REVIEW_PACKET_BLOCKED:" + reason)


def command(root, *argv):
    p = subprocess.run(["git", *argv], cwd=root, capture_output=True)
    require(p.returncode == 0, "GIT_" + "_".join(argv[:2]) +
            ":" + p.stderr.decode("utf-8", errors="replace")[-240:])
    return p.stdout


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


def identity(node):
    return (node["file"], node["kind"], node["name"], node.get("receiver", ""))


def declaration(raw, start, end):
    lines = raw.splitlines(keepends=True)
    require(type(start) is int and type(end) is int and
            1 <= start <= end <= len(lines), "INVALID_DECLARATION_LINE_RANGE")
    return b"".join(lines[start-1:end]).replace(b"\r\n", b"\n")


def read_json(root, filename):
    return json.loads((root / filename).read_text(encoding="utf-8"))


def build(root):
    root = root.resolve()
    head = command(root, "rev-parse", "HEAD").decode("ascii").strip()
    require(re.fullmatch("[0-9a-f]{40}", head) is not None, "INVALID_TARGET_COMMIT")
    # Only a deterministic canonical ICE scratch build may be uncommitted.
    dirty = command(root, "status", "--porcelain", "--untracked-files=all").decode("utf-8")
    require(all(line[3:].startswith(".ice/") for line in dirty.splitlines()),
            "NON_ICE_WORKTREE_DIRTY")
    old = json.loads(command(root, "show", "HEAD:.ice/manifest.json"))
    reviewed = json.loads(command(root, "show", "HEAD:.ice/code-intelligence.json"))
    now = read_json(root, ".ice/manifest.json")
    require(old["authority"] == now["authority"] ==
            "advisory syntax only; source wins", "CANONICAL_ICE_AUTHORITY_DRIFT")
    require(reviewed["authority"] == "advisory-navigation-only" and
            reviewed["repository"] == REPO, "REVIEWED_ATLAS_AUTHORITY_DRIFT")
    require(reviewed["coverage"]["symbols"] == OLD_DECLARATIONS and
            reviewed["coverage"]["source_files"] == OLD_FILES and
            len(reviewed["nodes"]) == OLD_DECLARATIONS,
            "REVIEWED_COVERAGE_DRIFT")

    old_files = {f["path"]: f for f in old["files"]}
    new_files = {f["path"]: f for f in now["files"]}
    old_symbols = {identity(s): s for s in old["symbols"]}
    new_symbols = {identity(s): s for s in now["symbols"]}
    old_reviewed = {
        (n["file"], n["kind"], n["name"], n.get("receiver", "")): n
        for n in reviewed["nodes"]
    }
    require(len(old_files) == len(old["files"]) == OLD_FILES and
            len(new_files) == len(now["files"]) == OLD_FILES+EXPECTED_ADDED_FILES and
            len(old_symbols) == len(old["symbols"]) == OLD_DECLARATIONS and
            len(new_symbols) == len(now["symbols"]) ==
            OLD_DECLARATIONS+EXPECTED_ADDED_DECLARATIONS,
            "STRUCTURAL_COHORT_COUNT_DRIFT")
    require(set(old_symbols) == set(old_reviewed),
            "REVIEWED_SYMBOL_LINEAGE_IDENTITY_DRIFT")
    new_keys = sorted(set(new_symbols)-set(old_symbols))
    removed_keys = set(old_symbols)-set(new_symbols)
    added_files = sorted(set(new_files)-set(old_files))
    require(len(new_keys) == EXPECTED_ADDED_DECLARATIONS and not removed_keys and
            len(added_files) == EXPECTED_ADDED_FILES,
            "UNEXPECTED_GRAPH_DELTA")
    require(sum(k[0].startswith(R3_PREFIX) for k in new_keys) ==
            EXPECTED_R3_DECLARATIONS, "R3_DECLARATION_SCOPE_DRIFT")

    for filename, file in new_files.items():
        path = root / filename
        require(path.is_file() and not path.is_symlink() and
                digest(path.read_bytes()) == file["sha256"],
                "REBUILT_FILE_HASH_DRIFT:" + filename)
        # Canonical ICE rebuild only proves agreement with WORKTREE content.
        # Bind all reviewed source candidates independently to immutable HEAD,
        # even when a malicious or stale Git assume-unchanged bit hides edits.
        actual = path.read_bytes().replace(b"\r\n", b"\n")
        committed = command(root, "show", "HEAD:" + filename).replace(b"\r\n", b"\n")
        require(actual == committed,
                "HEAD_COMMITTED_SOURCE_MISMATCH:" + filename)
    source_ref = reviewed["source_commit"]
    require(re.fullmatch("[0-9a-f]{40}", source_ref) is not None,
            "INVALID_PRIOR_SOURCE_ID")
    require(command(root, "cat-file", "-t", source_ref).strip() == b"commit",
            "AUTHENTICATED_HISTORICAL_SOURCE_UNAVAILABLE")
    # Protected main used a squash integration; the authentic historical
    # review source may be a non-ancestor even when every source blob is
    # available. Report the custody boundary; never silently re-anchor it.
    ancestry = subprocess.run(["git","merge-base","--is-ancestor",source_ref,head],
                              cwd=root, capture_output=True)
    require(ancestry.returncode in (0,1), "HISTORICAL_SOURCE_ANCESTRY_CHECK_ERROR")
    historical_ancestor = ancestry.returncode == 0

    historic_files = {}
    old_review_status = {"unchanged_declaration": 0,
                         "requires_fresh_review": 0}
    old_review_deltas = []
    for key in sorted(old_symbols):
        before = old_reviewed[key]
        after = new_symbols[key]
        proof = before.get("review_evidence")
        require(before.get("evidence_basis") == "source-reviewed" and
                isinstance(proof, dict) and
                proof.get("source_commit") == source_ref and
                proof.get("file") == key[0],
                "UNAUTHENTICATED_HISTORICAL_LABEL:" + str(key))
        path = key[0]
        if path not in historic_files:
            historic_files[path] = command(root, "show", source_ref + ":" + path)
        old_raw = historic_files[path]
        require(digest(old_raw) == before["file_sha256"] ==
                proof["file_sha256"], "HISTORIC_SOURCE_SHA_MISMATCH:" + path)
        left = declaration(old_raw, before["start_line"], before["end_line"])
        right = declaration((root / path).read_bytes(), after["start"], after["end"])
        require(proof["start_line"] == before["start_line"] and
                proof["end_line"] == before["end_line"] and
                all(t.encode("utf-8") in left for t in proof.get("literal_tokens", [])),
                "HISTORIC_SOURCE_REVIEW_ANCHOR_DRIFT:" + str(key))
        if left == right:
            old_review_status["unchanged_declaration"] += 1
        else:
            old_review_status["requires_fresh_review"] += 1
            old_review_deltas.append({"identity": list(key),
                                      "historical_decl_sha256": digest(left),
                                      "current_decl_sha256": digest(right),
                                      "status": "REVIEW_REQUIRED_UNAUTHENTICATED"})

    requests = []
    for key in new_keys:
        entry = new_symbols[key]
        filename = entry["file"]
        file = new_files[filename]
        proof = declaration((root / filename).read_bytes(),
                            entry["start"], entry["end"])
        require(len(proof) > 0, "EMPTY_GO_DECLARATION:" + str(key))
        requests.append({
            "identity": list(key),
            "file_sha256": file["sha256"],
            "declaration_sha256": digest(proof),
            "source_line_start": entry["start"],
            "source_line_end": entry["end"],
            "source_url": ("https://github.com/" + REPO + "/blob/" +
                           head + "/" + filename +
                           "#L" + str(entry["start"]) +
                           "-L" + str(entry["end"])),
            "syntax_calls_provisional": sorted(set(entry.get("calls", []))),
            "file_imports_unreviewed": sorted(file.get("imports", [])),
            "test_role_unreviewed": (
                "test_case" if filename.endswith("_test.go") and
                entry["kind"] == "function" and entry["name"].startswith("Test")
                else "test_helper" if filename.endswith("_test.go")
                else None),
            "review_status": "UNREVIEWED_NOT_AUTHORIZED",
            "authenticated_reviewer": None,
            "semantic_behavior_claim": None,
            "source_proof_approved": False,
        })
    result = {
        "schema": "wingless.ice.r3.exact-source-review-queue.v1",
        "authority": "advisory SOURCE-ONLY; not semantic review or acceptance",
        "repository": REPO,
        "target_source_commit": head,
        "historical_review_source_commit": source_ref,
        "historical_source_is_current_ancestor": historical_ancestor,
        "historical_pin_requires_independent_reanchor": not historical_ancestor,
        "previous_reviewed_symbols": OLD_DECLARATIONS,
        "previous_reviewed_files": OLD_FILES,
        "current_syntax_symbols": len(new_symbols),
        "current_syntax_files": len(new_files),
        "new_source_files": added_files,
        "added_declaration_count": len(new_keys),
        "removed_declaration_count": len(removed_keys),
        "r3_added_declaration_count": EXPECTED_R3_DECLARATIONS,
        "prior_reviewed_identity_reuse": old_review_status,
        "prior_reviewed_changed_declarations": old_review_deltas,
        "needs_independent_authentication": True,
        "semantic_atlas_publication_authorized": False,
        "scientific_first_use": False,
        "source_only": True,
        "unreviewed_declarations": requests,
    }
    return result


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    require(not args.output.exists(), "EXISTING_REVIEW_PACKET_MUST_NOT_OVERWRITE")
    require(args.output.parent.is_dir() and
            not args.output.resolve().is_relative_to(args.root.resolve()),
            "REVIEW_PACKET_MUST_BE_OUTSIDE_SOURCE")
    content = json.dumps(build(args.root), indent=2, sort_keys=True) + "\n"
    args.output.write_text(content, encoding="utf-8")
    result = json.loads(content)
    print("R3_ICE_REVIEW_PACKET_SHA256="+digest(content.encode("utf-8")))
    print("R3_ICE_REVIEW_REQUESTS="+str(result["added_declaration_count"]))
    print("R3_ICE_R3_DECLARATIONS="+str(result["r3_added_declaration_count"]))
    print("R3_ICE_PRIOR_EXACT_BODY_REUSE="+
          str(result["prior_reviewed_identity_reuse"]["unchanged_declaration"]))
    print("R3_ICE_PRIOR_REVIEW_REQUIRED="+
          str(result["prior_reviewed_identity_reuse"]["requires_fresh_review"]))
    print("R3_ICE_AUTHENTICATED_SEMANTIC_REVIEW=NOT_PERFORMED")


if __name__ == "__main__":
    main()
