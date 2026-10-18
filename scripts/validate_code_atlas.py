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
    "scripts/test_code_atlas_tamper.py",
    "scripts/test_atlas_browser.py",
    "scripts/atlas_ci_scope.py",
    "scripts/maintain_code_atlas.py",
    "scripts/test_source_aware_atlas.py",
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
    historical_source_cache = {}
    reviewed = 0
    evidenced = 0
    dependency_bound = 0
    test_contracts = 0
    test_cases = 0
    test_helpers = 0
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
        proof = node.get("review_evidence")
        if basis == "source-reviewed":
            reviewed += 1
            require(proof is not None, "source-reviewed symbol has no explicit proof: " + symbol_id)
            if proof is not None:
                evidenced += 1
                require(isinstance(proof, dict), "review proof not an object")
                require(proof.get("source_commit") == source_commit, "review source commit drift: " + symbol_id)
                require(proof.get("file") == s["file"], "review file identity drift: " + symbol_id)
                require(proof.get("file_sha256") == indexed_files[s["file"]]["sha256"], "review source hash drift: " + symbol_id)
                require(proof.get("start_line") == s["start"] and proof.get("end_line") == s["end"], "reviewed span drift: " + symbol_id)
                tokens = proof.get("literal_tokens")
                require(isinstance(tokens, list) and 1 <= len(tokens) <= 12, "inadequate review tokens: " + symbol_id)
                lines = (root / s["file"]).read_text(encoding="utf-8").splitlines()
                require(s["end"] <= len(lines), "source span outside file: " + symbol_id)
                span = "\\n".join(lines[s["start"] - 1:s["end"]])
                for token in tokens:
                    require(isinstance(token, str) and len(token) >= 4 and token.strip() == token,
                            "invalid source review token: " + symbol_id)
                    require(token in span, "source evidence token mismatch: " + symbol_id + " " + repr(token))
        # G26 reuse is checked against immutable source bytes; a new self-asserted
        # review cannot independently authenticate a human or admit changed code.
        lineage = node.get("provenance_lineage")
        if lineage is not None:
            require(isinstance(lineage, dict), "G26 invalid provenance lineage")
            prior = lineage.get("prior_source_commit", "")
            require(isinstance(prior, str) and
                    re.fullmatch(r"[0-9a-f]{40}", prior) is not None,
                    "G26 previous source identity invalid")
            old_start = lineage.get("prior_start_line")
            old_end = lineage.get("prior_end_line")
            require(isinstance(old_start, int) and isinstance(old_end, int) and
                    1 <= old_start <= old_end, "G26 historical declaration lines invalid")
            cache_id = (prior, s["file"])
            if cache_id not in historical_source_cache:
                read = subprocess.run(
                    ["git", "show", prior + ":" + s["file"]],
                    cwd=root, capture_output=True)
                require(read.returncode == 0, "G26 old source unavailable")
                historical_source_cache[cache_id] = read.stdout
            old_source = historical_source_cache[cache_id]
            require(hashlib.sha256(old_source).hexdigest() ==
                    lineage.get("prior_file_sha256"),
                    "G26 historical file hash mismatch")
            old_lines = old_source.splitlines(keepends=True)
            now_lines = (root / s["file"]).read_bytes().splitlines(keepends=True)
            require(old_end <= len(old_lines) and s["end"] <= len(now_lines),
                    "G26 historical or current span exceeds source")
            old_decl = b"".join(old_lines[old_start-1:old_end]).replace(b"\r\n", b"\n")
            now_decl = b"".join(now_lines[s["start"]-1:s["end"]]).replace(b"\r\n", b"\n")
            prev_digest = hashlib.sha256(old_decl).hexdigest()
            now_digest = hashlib.sha256(now_decl).hexdigest()
            require(prev_digest == lineage.get("prior_declaration_sha256") and
                    now_digest == lineage.get("current_declaration_sha256"),
                    "G26 declaration provenance digest mismatch")
            if lineage.get("mode") == "canonical-LF-declaration-byte-reuse":
                require(prev_digest == now_digest,
                        "G26 claimed reused declaration body changed")
            else:
                require(lineage.get("mode") ==
                        "fresh-review-receipt-unverified-identity",
                        "G26 unknown review lineage class")
                require(False,
                        "G26_UNAPPROVED_FRESH_REVIEW: source witnesses do not independently authenticate a reviewer")
        dependency_proof = node.get("file_level_local_package_dependencies")
        require(dependency_proof is not None, "missing file-level dependency evidence: " + symbol_id)
        if dependency_proof is not None:
            require(basis == "source-reviewed" and proof is not None,
                    "file-level package dependencies require a reviewed symbol: " + symbol_id)
            require(isinstance(dependency_proof, list) and all(isinstance(x, str) for x in dependency_proof),
                    "invalid local package import evidence: " + symbol_id)
            expected_imports = sorted(path.removeprefix("github.com/spoonman136668-ai/Wingless/")
                    for path in (dependencies.get(s["file"]) or [])
                    if path.startswith("github.com/spoonman136668-ai/Wingless/"))
            require(sorted(dependency_proof) == expected_imports,
                    "local package dependency identity drift: " + symbol_id)
            dependency_bound += 1
        contract = node.get("test_contract")
        if s["file"].endswith("_test.go"):
            require(basis == "source-reviewed" and proof is not None and isinstance(contract, dict),
                    "unreviewed test contract or missing source proof: " + symbol_id)
            classification = ("test_case" if s["kind"] == "function" and s["name"].startswith("Test")
                              else "test_fixture_or_helper")
            require(contract.get("classification") == classification,
                    "incorrect test-case versus helper classification: " + symbol_id)
            require(contract.get("test_file") == s["file"], "test file contract identity drift: " + symbol_id)
            require(contract.get("source_anchor") in proof["literal_tokens"],
                    "test contract witness not in reviewed source tokens: " + symbol_id)
            require(isinstance(contract.get("interpretation"), str)
                    and "do not prove untested environments" in contract["interpretation"],
                    "test contract overstates authority: " + symbol_id)
            test_contracts += 1
            if classification == "test_case":
                test_cases += 1
            else:
                test_helpers += 1
        else:
            require(contract is None, "non-test symbol has test assertion contract: " + symbol_id)
        if basis != "source-reviewed":
            require(proof is None, "unreviewed symbol must not carry behavioral proof: " + symbol_id)
            require("not a verified behavioral description" in node["description"], "unsupported behavioral inference")
    require(atlas["coverage"].get("source_dependency_bound_labels") == dependency_bound,
            "source dependency evidence count drift")
    require(atlas["coverage"].get("source_evidenced_test_contracts") == test_contracts,
            "test contract coverage drift")
    require(atlas["coverage"].get("test_cases") == test_cases and
            atlas["coverage"].get("test_fixtures_and_helpers") == test_helpers,
            "test fixture/test case classification drift")
    require(test_contracts == 116, "frozen test contract cohort not fully reviewed")
    require(atlas["coverage"]["source_reviewed_labels"] == reviewed, "reviewed count drift")
    require(atlas["coverage"].get("source_evidenced_new_labels") == evidenced,
            "new source-evidence count drift")
    require(evidenced == reviewed == len(nodes), "source evidence not uniform across all symbols")
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
    const_values = dict(re.findall(r'([A-Za-z]+)\s+State\s*=\s*"([A-Z_]+)"', source))
    require(const_values, "Go lifecycle enum values absent")
    transitions = {}
    for state, destinations in re.findall(r"([A-Za-z]+):\s*\{([^}]+)\}", source):
        require(state in const_values, "unknown lifecycle source state: " + state)
        options = [v.strip() for v in destinations.split(",") if v.strip()]
        require(all(v in const_values for v in options), "unknown lifecycle destination state")
        transitions[const_values[state]] = [const_values[v] for v in options]
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
        "new_source_evidenced_labels": evidenced,
        "fully_evidenced_symbol_count": evidenced,
        "dependency_index_bound_labels": dependency_bound,
        "source_verified_test_contracts": test_contracts,
        "test_cases": test_cases,
        "test_fixtures_and_helpers": test_helpers,
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
