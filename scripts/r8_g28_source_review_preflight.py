#!/usr/bin/env python3
"""Read-only G28 migration inventory for the pinned Wingless R8 source lineage.

Machine source witnesses are NEVER review receipts. This tool cannot write ICE,
promote a label, grant first-use, create a PR, or change any accepted ref.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

SCHEMA = "wingless.r8.g28.source-review-preflight.v1"
ACCEPTED = "9a92fcae40c811c724b4214da03fff34126f6b4f"
R8_INDEXED = "21cde395d77520318079bf45098f8441d6c20b0b"
HISTORICAL_REVIEW = "f87fa798050124d5275874c7b623ab897f19324e"
ALLOWED_SOURCE_DRIFT = {
    ".github/workflows/native.yml",
    ".github/workflows/wingless-code-atlas.yml",
    "cmd/wingless/main.go",
}
SHA40 = re.compile(r"[0-9a-f]{40}\Z")
# This is the exact frozen, additive test step in the EXISTING Atlas workflow.
# Any other workflow modification fails closed rather than becoming a broad grant.
WORKFLOW = ".github/workflows/wingless-code-atlas.yml"
WORKFLOW_GIT_BLOB = "39b7e20f867bc091a0d4c0a3b8bd32a8327d7efc"
INDEX_BLOBS = {".ice/manifest.json":"9fdf7828b6bcd8856af5b862079415b5fefca496",".ice/architecture-map.json":"f12e143111d69a9a618935c0010742844e63ed34"}
LOCAL_ONLY = frozenset({
    "scripts/r8_g28_source_review_preflight.py",
    "scripts/test_r8_g28_source_review_preflight.py",
    WORKFLOW,
    *INDEX_BLOBS,
})


class R8G28Denied(ValueError):
    pass


def require(ok: bool, code: str) -> None:
    if not ok:
        raise R8G28Denied("R8_G28_" + code)


def digest(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def key(s: dict) -> tuple[str, str, str, str]:
    return (s["file"], s.get("receiver", ""), s["name"], s["kind"])


def verified_source(read, commit: str, path: str, expected: str) -> bytes:
    require(SHA40.fullmatch(commit) is not None, "COMMIT_INVALID")
    require(path and not path.startswith("/") and "\\" not in path
            and all(p not in ("", ".", "..") for p in path.split("/")),
            "PATH_INVALID")
    raw = read(commit, path)
    require(type(raw) is bytes and digest(raw) == expected, "SOURCE_HASH_DRIFT")
    return raw


def declaration(raw: bytes, s: dict) -> bytes:
    lines = raw.replace(b"\r\n", b"\n").splitlines(keepends=True)
    start, end = s.get("start", s.get("start_line")), s.get("end", s.get("end_line"))
    require(type(start) is int and type(end) is int
            and 1 <= start <= end <= len(lines), "DECLARATION_SPAN_INVALID")
    return b"".join(lines[start - 1:end])


def inventory(old_manifest: dict, target_manifest: dict, old_atlas: dict,
              read, *, target_sha: str) -> dict:
    """Compare immutable sources; NEVER create authenticated review evidence."""
    require(target_sha == R8_INDEXED, "TARGET_SHA_NOT_PINNED")
    require(old_atlas.get("source_commit") == HISTORICAL_REVIEW
            and old_atlas.get("authority") == "advisory-navigation-only",
            "HISTORICAL_ATLAS_UNTRUSTED")
    old_files = {x["path"]: x for x in old_manifest["files"]}
    new_files = {x["path"]: x for x in target_manifest["files"]}
    old_s = {key(s): s for s in old_manifest["symbols"]}
    new_s = {key(s): s for s in target_manifest["symbols"]}
    reviewed_nodes = {key(n): n for n in old_atlas["nodes"]}
    require(len(old_files) == len(old_manifest["files"])
            and len(new_files) == len(target_manifest["files"])
            and len(old_s) == len(old_manifest["symbols"])
            and len(new_s) == len(target_manifest["symbols"])
            and len(reviewed_nodes) == len(old_atlas["nodes"]),
            "DUPLICATE_INDEX_IDENTITY")
    require(len(old_files) == 138 and len(new_files) == 195
            and len(old_s) == 352 and len(new_s) == 655
            and set(old_files) <= set(new_files)
            and set(old_s) == set(reviewed_nodes)
            and set(old_s) <= set(new_s),
            "FROZEN_COHORT_DRIFT")
    changed_files = {
        path for path, old in old_files.items()
        if new_files[path]["sha256"] != old["sha256"]
    }
    require(changed_files == ALLOWED_SOURCE_DRIFT, "CHANGED_OLD_FILE_SCOPE")
    require(len(set(new_files) - set(old_files)) == 57, "NEW_FILE_SCOPE")
    cache: dict[tuple[str, str], bytes] = {}

    def get(sha: str, path: str, expected: str) -> bytes:
        k = (sha, path)
        if k not in cache:
            cache[k] = verified_source(read, sha, path, expected)
        else:
            require(digest(cache[k]) == expected, "CACHED_SOURCE_IDENTITY_DRIFT")
        return cache[k]

    retained, changed, provisional = 0, [], []
    for k, old in sorted(old_s.items()):
        current, node = new_s[k], reviewed_nodes[k]
        require(node.get("evidence_basis") == "source-reviewed"
                and isinstance(node.get("review_evidence"), dict),
                "HISTORICAL_REVIEW_MISSING")
        proof = node["review_evidence"]
        require(proof.get("source_commit") == HISTORICAL_REVIEW
                and proof.get("file") == k[0]
                and proof.get("file_sha256") == old_files[k[0]]["sha256"]
                and proof.get("start_line") == old["start"]
                and proof.get("end_line") == old["end"],
                "HISTORICAL_PROOF_DRIFT")
        old_src = get(ACCEPTED, k[0], old_files[k[0]]["sha256"])
        cur_src = get(target_sha, k[0], new_files[k[0]]["sha256"])
        prior, now = declaration(old_src, old), declaration(cur_src, current)
        tokens = proof.get("literal_tokens")
        require(isinstance(tokens, list) and len(tokens) >= 1
                and all(isinstance(t, str) and t.encode("utf-8") in prior
                        for t in tokens), "HISTORICAL_LITERAL_WITNESS_DRIFT")
        if prior == now:
            require(all(t.encode("utf-8") in now for t in tokens),
                    "REUSED_REVIEW_WITNESS_DRIFT")
            retained += 1
        else:
            changed.append({
                "key": list(k), "old_span": [old["start"], old["end"]],
                "new_span": [current["start"], current["end"]],
                "old_declaration_sha256": digest(prior),
                "new_declaration_sha256": digest(now),
                "review_required": True,
                "reviewed": False,
            })
    for k in sorted(set(new_s) - set(old_s)):
        cur = new_s[k]
        raw = get(target_sha, k[0], new_files[k[0]]["sha256"])
        fragment = declaration(raw, cur)
        require(len(fragment) > 0, "NEW_DECLARATION_EMPTY")
        provisional.append({
            "key": list(k), "target_source_commit": target_sha,
            "file_sha256": new_files[k[0]]["sha256"],
            "declaration_sha256": digest(fragment),
            "start_line": cur["start"], "end_line": cur["end"],
            "evidence_basis": "lexical-provisional",
            "review_evidence": None, "independent_approval": False,
            "source_url": "https://github.com/spoonman136668-ai/Wingless/blob/"
                          + target_sha + "/" + k[0] + "#L"
                          + str(cur["start"]) + "-L" + str(cur["end"]),
        })
    require(retained == 351 and len(changed) == 1 and
            changed[0]["key"] == ["cmd/wingless/main.go", "", "main", "function"]
            and len(provisional) == 303, "REVIEW_ACCOUNTING_DRIFT")
    return {
        "schema": SCHEMA, "accepted_source_sha": ACCEPTED,
        "structural_source_sha": target_sha,
        "historical_review_commit": HISTORICAL_REVIEW,
        "historical_reviewed_identities": 352,
        "byte_identical_reviewed_declarations": retained,
        "changed_historical_declarations": changed,
        "new_unreviewed_provisional_declarations": provisional,
        "new_symbol_count": len(provisional),
        "source_review_completed": False,
        "source_admission_authorized": False,
        "scientific_first_use": False,
        "status": "BLOCKED_INDEPENDENT_SOURCE_REVIEW",
    }


def git(root: Path, *args: str) -> bytes:
    return subprocess.check_output(["git", *args], cwd=root, stderr=subprocess.PIPE)


def run(root: Path) -> dict:
    root = root.resolve()
    head = git(root, "rev-parse", "HEAD").decode().strip()
    require(SHA40.fullmatch(head) is not None, "HEAD_INVALID")
    require(subprocess.run(["git", "merge-base", "--is-ancestor", R8_INDEXED, head],
                           cwd=root, stdout=subprocess.DEVNULL,
                           stderr=subprocess.DEVNULL).returncode == 0,
            "R8_SOURCE_ANCESTRY_DRIFT")
    require(not git(root, "status", "--porcelain").strip(), "DIRTY_CHECKOUT")
    changed = set(git(root, "diff", "--name-only", R8_INDEXED, head).decode().splitlines())
    require(changed <= LOCAL_ONLY, "SOURCE_OR_REVIEW_POLICY_MUTATION")
    require(git(root, "rev-parse", "HEAD:" + WORKFLOW).decode().strip()
            == WORKFLOW_GIT_BLOB, "WORKFLOW_CONTENT_DRIFT")
    for indexed_path, indexed_sha in INDEX_BLOBS.items():
        require(git(root, "rev-parse", "HEAD:" + indexed_path).decode().strip()
                == indexed_sha, "GENERATED_ICE_BLOB_DRIFT")
    for ref in (ACCEPTED, HISTORICAL_REVIEW, R8_INDEXED):
        require(git(root, "cat-file", "-t", ref).strip() == b"commit",
                "MISSING_IMMUTABLE_SOURCE")
    def source(ref: str, path: str) -> bytes:
        return git(root, "show", ref + ":" + path)
    def load(ref: str, path: str) -> dict:
        return json.loads(source(ref, path))
    return inventory(load(ACCEPTED, ".ice/manifest.json"),
                     load(R8_INDEXED, ".ice/manifest.json"),
                     load(ACCEPTED, ".ice/code-intelligence.json"),
                     source, target_sha=R8_INDEXED)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, default=Path("."))
    args = parser.parse_args()
    print(json.dumps(run(args.root), sort_keys=True, separators=(",", ":")))


if __name__ == "__main__":
    main()
