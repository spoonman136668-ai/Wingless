#!/usr/bin/env python3
"""Fail-closed, source-bound validation of the Wingless read-only semantic atlas."""
import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

REPOSITORY = "spoonman136668-ai/Wingless"
SCHEMA = "wingless.ice.semantic-code-atlas.v1"
ALLOWED_DIFF = {
    ".ice/code-intelligence.json",
    ".ice/code-intelligence.ice",
    ".ice/code-intelligence.html",
    ".ice/manifest.json",
    ".ice/architecture-map.json",
    "scripts/validate_code_atlas.py",
    ".github/workflows/wingless-code-atlas.yml",
}

def require(condition, message):
    if not condition:
        raise SystemExit("ATLAS_FAIL_CLOSED: " + message)

def read_json(p):
    return json.loads(p.read_text(encoding="utf-8-sig"))

def git(root, *args):
    p = subprocess.run(["git", *args], cwd=root, capture_output=True, text=True, check=True)
    return p.stdout.strip()

def identity(symbol):
    receiver = (symbol.get("receiver", "") + ".") if symbol.get("receiver") else ""
    return (
        f"wingless:{symbol['file']}::{symbol['kind']}::"
        f"{receiver}{symbol['name']}@{symbol['start']}"
    )

def key(symbol):
    receiver = (symbol.get("receiver", "") + ".") if symbol.get("receiver") else ""
    return f"{symbol['file']}|{receiver}{symbol['name']}"

def matches_type_target(target, name, receiver, local_package):
    if not target.endswith("." + name):
        return False
    if "/" + local_package + "." not in target and "/" + local_package + "." not in target.replace("(", ""):
        # go/types format for methods can be "(*github.../pkg.Type).Method".
        if "/" + local_package + "." not in target.replace(")", ""):
            return False
    if receiver and receiver not in target:
        return False
    return True

def validate(root, relationship_file):
    root = root.resolve()
    atlas = read_json(root / ".ice/code-intelligence.json")
    indexed = read_json(root / ".ice/symbol-map.json")
    architecture = read_json(root / ".ice/architecture-map.json")
    manifest = read_json(root / ".ice/manifest.json")
    dependencies = read_json(root / ".ice/dependency-map.json")
    typed = read_json(relationship_file)
    require(atlas.get("schema") == SCHEMA, "invalid semantic schema")
    require(atlas.get("repository") == REPOSITORY, "wrong repository boundary")
    require(atlas.get("authority") == "advisory-navigation-only", "authority changed")
    source_commit = atlas.get("source_commit")
    require(re.fullmatch(r"[0-9a-f]{40}", source_commit or "") is not None, "invalid source SHA")
    require(git(root, "cat-file", "-t", source_commit) == "commit", "source commit unavailable")
    changed = set(git(root, "diff", "--name-only", source_commit, "HEAD").splitlines())
    require(changed.issubset(ALLOWED_DIFF), "unexpected changed paths: " + repr(sorted(changed - ALLOWED_DIFF)))
    require(architecture == manifest.get("files"), "original ICE file coverage diverged")
    indexed_files = {f["path"]: f for f in architecture}
    require(len(indexed_files) == len(architecture), "duplicated file path")
    for file, entry in indexed_files.items():
        data = (root / file).read_bytes()
        require(hashlib.sha256(data).hexdigest() == entry["sha256"], "source hash mismatch: " + file)
    nodes = atlas.get("nodes", [])
    require(len(nodes) == len(indexed), "symbol coverage count differs from source ICE")
    require(atlas["coverage"]["source_files"] == len(architecture), "source file count drift")
    require(atlas["coverage"]["symbols"] == len(indexed), "symbol count drift")
    expected = {identity(s): s for s in indexed}
    require(len(expected) == len(indexed), "duplicate raw symbol identity")
    actual = {n["id"]: n for n in nodes}
    require(len(actual) == len(nodes) and actual.keys() == expected.keys(), "symbol identity drift")
    reviewed = 0
    for symbol_id, node in actual.items():
        s = expected[symbol_id]
        for field in ("file", "kind", "name"):
            require(node[field] == s[field], "symbol field mismatch: " + symbol_id)
        require(node["start_line"] == s["start"] and node["end_line"] == s["end"], "symbol lines mismatch: " + symbol_id)
        require(node.get("receiver", "") == s.get("receiver", ""), "receiver mismatch: " + symbol_id)
        require(node["file_sha256"] == indexed_files[s["file"]]["sha256"], "symbol file hash mismatch: " + symbol_id)
        require(node["syntax_calls"] == sorted(set(s.get("calls", []))), "syntax calls mismatch: " + symbol_id)
        require(node["source_url"].endswith(f"/{s['file']}#L{s['start']}-L{s['end']}"), "pinned span mismatch: " + symbol_id)
        require("/blob/" + source_commit + "/" in node["source_url"], "source hyperlink not immutable")
        basis = node["evidence_basis"]
        require(basis in ("source-reviewed", "lexical-provisional"), "unknown label evidence")
        if basis == "source-reviewed":
            reviewed += 1
        else:
            require("not a verified behavioral description" in node["description"], "unsupported behavioral inference")
    require(atlas["coverage"]["source_reviewed_labels"] == reviewed, "reviewed count drift")
    require(atlas["coverage"]["lexical_provisional_labels"] == len(nodes)-reviewed, "provisional count drift")

    ids = set(actual)
    edges = atlas.get("relationships", [])
    reviewed_edges = []
    for e in edges:
        require(e["evidence_basis"] in ("source-reviewed", "source-derived-import", "syntax-only-provisional"), "unrecognized edge provenance")
        if e["kind"] == "imports_local_package":
            require(e["from"] in dependencies, "unknown importer")
            require(
                "github.com/spoonman136668-ai/Wingless/" + e["to"] in (dependencies[e["from"]] or []),
                "import not supported by original dependency index"
            )
        else:
            require(e["from"] in ids and e["to"] in ids, "dangling symbol relationship")
            if e["evidence_basis"] == "source-reviewed":
                reviewed_edges.append(e)
                source, dest = actual[e["from"]], actual[e["to"]]
                # The legacy syntax extractor omits nested Go selectors such as
                # s.runner.Run. Never mislabel that parser limitation as a source
                # disproof; the type-resolved call-site check below is mandatory
                # for EVERY source-reviewed relationship, including nested calls.
                require(source["file"] != dest["file"] or source["id"] != dest["id"],
                        "self-edge cannot prove an independently reviewed relationship")
    require(len(reviewed_edges) == atlas["coverage"]["source_reviewed_relationships"], "reviewed edge count drift")

    require(typed.get("authority", "").startswith("derived; active platform"), "type report authority mismatch")
    resolved = typed.get("relationships", [])
    # Go encodes an uninitialized []string as JSON null; null means no diagnostics.
    diagnostics = typed.get("diagnostics") or []
    require(isinstance(diagnostics, list), "missing typecheck diagnostics")
    confirmed = 0
    for e in reviewed_edges:
        src = actual[e["from"]]
        dst = actual[e["to"]]
        pkg = indexed_files[dst["file"]].get("package", "")
        candidates = [r for r in resolved if
            r.get("kind") == "resolved_call_target"
            and r.get("file") == src["file"]
            and src["start_line"] <= r.get("line", -1) <= src["end_line"]
            and matches_type_target(r.get("to", ""), dst["name"], dst.get("receiver", ""), pkg)
        ]
        require(candidates, "type-unverified reviewed relationship: " + e["label"]
                + "; typecheck diagnostics=" + repr(diagnostics[:3]))
        confirmed += 1

    machine = atlas.get("state_machines", [])
    require(len(machine) == 1 and machine[0]["name"] == "Inference lifecycle", "state machine identity drift")
    source = (root / "inference/lifecycle.go").read_text(encoding="utf-8")
    fragment = source.split("allowed := map[State][]State{", 1)
    require(len(fragment) == 2, "lifecycle transition table missing")
    fragment = fragment[1].split("}", 1)[0]  # Read below from full source instead for nested braces.
    transitions = {a: [v.strip() for v in b.split(",") if v.strip()]
                   for a, b in re.findall(r"([A-Za-z]+):\s*\{([^}]+)\}", source)}
    require(transitions == machine[0]["transitions"], "lifecycle graph does not match source")

    html = (root / ".ice/code-intelligence.html").read_text(encoding="utf-8")
    match = re.search(r'<script type="application/json" id="atlas">(.*?)</script>', html, re.S)
    require(match is not None, "HTML lacks embedded machine data")
    require(json.loads(match.group(1)) == atlas, "offline browser differs from machine index")
    narrative = (root / ".ice/code-intelligence.ice").read_text(encoding="utf-8")
    require("SOURCE_COMMIT: " + source_commit in narrative, "durable ICE missing source provenance")
    require("ADVISORY" in narrative and "provisional" in narrative.lower(), "durable ICE overclaims authority")
    summary = {
        "result": "ATLAS_SEMANTIC_QUALIFICATION_PASS",
        "source_commit": source_commit,
        "repository": REPOSITORY,
        "source_files": len(architecture),
        "symbols": len(nodes),
        "source_reviewed_labels": reviewed,
        "typed_confirmed_reviewed_edges": confirmed,
        "typecheck_diagnostics": len(diagnostics),
    }
    print(json.dumps(summary, sort_keys=True))

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--relationships", type=Path, required=True)
    args = parser.parse_args()
    validate(args.root, args.relationships)
