#!/usr/bin/env python3
"""G26 source-aware Wingless atlas regeneration: advisory, fail closed.

Requires a structurally rebuilt canonical .ice index at the exact target HEAD.
Previous labels transfer only if raw Go declaration bytes are identical.
Any changed declaration requires an explicit, source-bound new review receipt;
the receipt's reviewer identity is an assertion, NOT independent authentication.
No automatic human/behavioral review is conferred by this generator.
"""
import argparse
import copy
import hashlib
import json
import re
import subprocess
from pathlib import Path

REPO = "spoonman136668-ai/Wingless"
SCHEMA = "wingless.ice.semantic-code-atlas.v1"
ATLAS = ".ice/code-intelligence.json"
HTML = ".ice/code-intelligence.html"
ICE = ".ice/code-intelligence.ice"
START = '<script type="application/json" id="atlas">'
END = "</script>"


def require(ok, why):
    if not ok:
        raise SystemExit("G26_FAIL_CLOSED: " + why)


def git(root, *args):
    p = subprocess.run(["git", *args], cwd=root, capture_output=True)
    require(p.returncode == 0, "git " + " ".join(args) + ": " +
            p.stderr.decode("utf-8", errors="replace")[-280:])
    return p.stdout


def sha(data):
    return hashlib.sha256(data).hexdigest()


def node_key(n):
    return (n["file"], n["kind"], n["name"], n.get("receiver", ""))


def symbol_id(n):
    receiver = n.get("receiver", "")
    return ("wingless:" + n["file"] + "::" + n["kind"] + "::" +
            (receiver + "." if receiver else "") + n["name"] + "@" +
            str(n["start_line"]))


def source_bytes(root, source_sha, path):
    require(not Path(path).is_absolute() and ".." not in Path(path).parts,
            "invalid historical source path")
    return git(root, "show", source_sha + ":" + path)


def span(data, start, end):
    lines = data.splitlines(keepends=True)
    require(isinstance(start, int) and isinstance(end, int)
            and 1 <= start <= end <= len(lines), "invalid declaration span")
    return b"".join(lines[start-1:end])


def load(root, filename):
    return json.loads((root / filename).read_text(encoding="utf-8-sig"))


def generate(root, receipts_path=None):
    root = root.resolve()
    target = git(root, "rev-parse", "HEAD").decode("ascii").strip()
    require(re.fullmatch(r"[0-9a-f]{40}", target) is not None,
            "target must be exact HEAD commit")
    old = load(root, ATLAS)
    idx = load(root, ".ice/manifest.json")
    symbols = load(root, ".ice/symbol-map.json")
    deps = load(root, ".ice/dependency-map.json")
    require(old.get("schema") == SCHEMA and old.get("repository") == REPO,
            "unrecognized prior atlas")
    previous = old.get("source_commit", "")
    require(re.fullmatch(r"[0-9a-f]{40}", previous) is not None
            and previous != target, "source lineage is not a new commit")
    require(git(root, "cat-file", "-t", previous).strip() == b"commit",
            "historical source unavailable")
    require(git(root, "merge-base", previous, target).strip() ==
            previous.encode("ascii"), "prior source is not an ancestor")
    require(idx.get("authority") == "advisory syntax only; source wins",
            "canonical ICE authority changed")
    # G27 review hardening: a structurally consistent working tree is not
    # sufficient proof that target HEAD actually contains those source bytes.
    # Only canonical ICE outputs may be dirty after 'go run ./cmd/ice build'.
    status = git(root, "status", "--porcelain", "--untracked-files=all")
    for entry in status.decode("utf-8", errors="replace").splitlines():
        require(entry[3:].startswith(".ice/"),
                "G27_DIRTY_NON_ICE_WORKTREE: " + entry[3:])
    # Context-sensitive descriptions cannot inherit earlier review if module
    # resolution or text-conversion rules changed even when declarations did not.
    for context_path in ("go.mod", ".gitattributes"):
        require(source_bytes(root, previous, context_path) ==
                source_bytes(root, target, context_path),
                "G27_DEPENDENCY_CONTEXT_CHANGED: " + context_path)
    # All source-bearing files must exactly match the advertised target commit.
    # The canonical LF transform is safe for this repository's explicit
    # '.gitattributes: * text=auto eol=lf' contract.
    for item in idx["files"]:
        path = item["path"]
        raw = (root / path).read_bytes()
        committed = source_bytes(root, target, path)
        require(sha(raw) == item["sha256"],
                "G27_STRUCTURAL_INDEX_MISMATCH: " + path)
        require(raw.replace(b"\r\n", b"\n") ==
                committed.replace(b"\r\n", b"\n"),
                "G27_TARGET_SOURCE_CONTENT_MISMATCH: " + path)
    files = {f["path"]: f for f in idx["files"]}
    require(len(files) == len(idx["files"]), "duplicate current indexed file")
    old_by_key = {node_key(n): n for n in old["nodes"]}
    current_by_key = {
        node_key(s): s for s in symbols
    }
    require(len(old_by_key) == len(old["nodes"]) and
            len(current_by_key) == len(symbols), "duplicate symbol identity")
    require(set(old_by_key).issubset(current_by_key),
            "G28_REMOVED_OR_RENAMED_REVIEWED_SYMBOL")
    require(len(old_by_key) >= 352 and len(idx["files"]) >= 138
            and {n["file"] for n in old["nodes"]}.issubset(files),
            "G28_FROZEN_REVIEWED_BASELINE_MISSING")
    added_keys = set(current_by_key) - set(old_by_key)
    require(not added_keys or receipts_path is not None,
            "G28_NEW_SYMBOL_SOURCE_RECEIPTS_REQUIRED")
    require(git(root, "rev-parse", "HEAD").decode().strip() == target,
            "target identity moved")
    receipts = []
    if receipts_path is not None:
        receipts = json.loads(receipts_path.read_text(encoding="utf-8"))
    require(isinstance(receipts, list), "receipt input must be a list")
    receipt_by_key = {}
    for item in receipts:
        require(isinstance(item, dict) and isinstance(item.get("key"), list)
                and len(item["key"]) == 4, "invalid review receipt key")
        key = tuple(item["key"])
        require(key not in receipt_by_key, "duplicate fresh-review receipt")
        receipt_by_key[key] = item

    historical_sources = {}
    current_sources = {}
    new_nodes = []
    id_mapping = {}
    used_receipts = set()
    reused = 0
    fresh = 0
    provisional_added = 0
    for key, current in current_by_key.items():
        oldnode = old_by_key.get(key)
        path = key[0]
        require(path in files, "current node outside canonical source files")
        if oldnode is None:
            receipt = receipt_by_key.get(key)
            require(receipt is not None, "G28_NEW_SYMBOL_REVIEW_RECEIPT_REQUIRED " + repr(key))
            if path not in current_sources:
                raw = (root / path).read_bytes()
                require(sha(raw) == files[path]["sha256"],
                        "G28_NEW_SYMBOL_SOURCE_HASH_DRIFT: " + path)
                current_sources[path] = raw
            new_file = current_sources[path]
            new_decl = span(new_file, current["start"], current["end"])
            require(receipt.get("target_source_commit") == target
                    and receipt.get("file_sha256") == sha(new_file)
                    and receipt.get("declaration_sha256") == sha(new_decl)
                    and receipt.get("start_line") == current["start"]
                    and receipt.get("end_line") == current["end"],
                    "G28_NEW_SYMBOL_RECEIPT_IDENTITY_DRIFT: " + repr(key))
            require(isinstance(receipt.get("reviewer"), str)
                    and len(receipt["reviewer"].strip()) >= 3,
                    "G28_NEW_SYMBOL_REVIEWER_ASSERTION_REQUIRED")
            for field in ("description", "label"):
                require(isinstance(receipt.get(field), str)
                        and len(receipt[field].strip()) >= 12,
                        "G28_NEW_SYMBOL_DESCRIPTION_REQUIRED: " + field)
            tokens = receipt.get("literal_tokens")
            require(isinstance(tokens, list) and 1 <= len(tokens) <= 12
                    and all(isinstance(t, str) and len(t) >= 4
                            and t.strip() == t and t.encode("utf-8") in new_decl
                            for t in tokens),
                    "G28_NEW_SYMBOL_LITERAL_WITNESS_INVALID")
            receiver = current.get("receiver", "")
            node = {
                "id": symbol_id({
                    "file": path, "kind": current["kind"], "name": current["name"],
                    "receiver": receiver, "start_line": current["start"]}),
                "file": path, "file_sha256": sha(new_file),
                "kind": current["kind"], "name": current["name"],
                "start_line": current["start"], "end_line": current["end"],
                "label": "PROVISIONAL: " + receipt["label"],
                "description": (receipt["description"]
                                + " — not a verified behavioral description"),
                "evidence_basis": "lexical-provisional",
                "source_url": ("https://github.com/" + REPO + "/blob/"
                               + target + "/" + path + "#L"
                               + str(current["start"]) + "-L"
                               + str(current["end"])),
                "syntax_calls": sorted(set(current.get("calls") or [])),
                "review_evidence": None,
                "file_level_local_package_dependencies": None,
                "provenance_lineage": {
                    "mode": "new-symbol-source-bound-provisional",
                    "target_source_commit": target,
                    "file_sha256": sha(new_file),
                    "current_declaration_sha256": sha(new_decl.replace(b"\r\n", b"\n")),
                    "literal_tokens": tokens,
                    "reviewer_assertion": receipt["reviewer"],
                    "independent_approval": False,
                },
            }
            new_nodes.append(node)
            used_receipts.add(key)
            provisional_added += 1
            continue
        if path not in historical_sources:
            raw = source_bytes(root, previous, path)
            require(sha(raw) == oldnode["file_sha256"],
                    "historical source file hash mismatch: " + path)
            historical_sources[path] = raw
        old_file = historical_sources[path]
        require(sha(old_file) == oldnode["file_sha256"],
                "historical node file hash mismatch")
        if path not in current_sources:
            raw = (root / path).read_bytes()
            require(sha(raw) == files[path]["sha256"],
                    "current canonical file hash mismatch: " + path)
            current_sources[path] = raw
        new_file = current_sources[path]
        old_decl = span(old_file, oldnode["start_line"], oldnode["end_line"])
        new_decl = span(new_file, current["start"], current["end"])
        oldproof = oldnode.get("review_evidence")
        require(oldnode.get("evidence_basis") == "source-reviewed"
                and isinstance(oldproof, dict)
                and oldproof.get("source_commit") == previous
                and oldproof.get("file_sha256") == sha(old_file)
                and oldproof.get("start_line") == oldnode["start_line"]
                and oldproof.get("end_line") == oldnode["end_line"],
                "historical reviewed witness invalid: " + repr(key))
        for token in oldproof.get("literal_tokens", []):
            require(isinstance(token, str) and token.encode("utf-8") in old_decl,
                    "historical token outside declaration: " + repr(key))
        n = copy.deepcopy(oldnode)
        # LF/CRLF checkout normalization is not a source-body edit.
        old_canonical = old_decl.replace(b"\r\n", b"\n")
        new_canonical = new_decl.replace(b"\r\n", b"\n")
        if old_canonical == new_canonical:
            require(key not in receipt_by_key,
                    "review receipt for unchanged declaration is forbidden")
            reused += 1
            n["provenance_lineage"] = {
                "mode": "canonical-LF-declaration-byte-reuse",
                "prior_source_commit": previous,
                "prior_start_line": oldnode["start_line"],
                "prior_end_line": oldnode["end_line"],
                "prior_file_sha256": sha(old_file),
                "prior_declaration_sha256": sha(old_canonical),
                "current_declaration_sha256": sha(new_canonical),
            }
        else:
            receipt = receipt_by_key.get(key)
            require(receipt is not None,
                    "G26_REVIEW_REQUIRED " + repr(key))
            require(receipt.get("target_source_commit") == target and
                    receipt.get("file_sha256") == sha(new_file) and
                    receipt.get("declaration_sha256") == sha(new_decl) and
                    receipt.get("start_line") == current["start"] and
                    receipt.get("end_line") == current["end"],
                    "fresh-review receipt source identity drift: " + repr(key))
            require(isinstance(receipt.get("reviewer"), str) and
                    len(receipt["reviewer"].strip()) >= 3,
                    "fresh reviewer assertion missing")
            for field in ("description", "label"):
                require(isinstance(receipt.get(field), str) and
                        len(receipt[field].strip()) >= 12,
                        "fresh descriptive review missing: " + field)
            tokens = receipt.get("literal_tokens")
            require(isinstance(tokens, list) and 1 <= len(tokens) <= 12 and
                    all(isinstance(t, str) and len(t) >= 4 and
                        t.strip() == t and t.encode("utf-8") in new_decl
                        for t in tokens), "fresh witness tokens not source-bound")
            if path.endswith("_test.go"):
                require("test_contract" in receipt,
                        "changed test requires refreshed test contract")
                n["test_contract"] = receipt["test_contract"]
            n["label"] = receipt["label"]
            n["description"] = receipt["description"]
            n["review_evidence"]["literal_tokens"] = tokens
            n["provenance_lineage"] = {
                "mode": "fresh-review-receipt-unverified-identity",
                "prior_source_commit": previous,
                "prior_start_line": oldnode["start_line"],
                "prior_end_line": oldnode["end_line"],
                "prior_file_sha256": sha(old_file),
                "prior_declaration_sha256": sha(old_decl),
                "current_declaration_sha256": sha(new_decl),
                "reviewer_assertion": receipt["reviewer"],
                "independent_approval": False,
            }
            fresh += 1
            used_receipts.add(key)
        n["start_line"] = current["start"]
        n["end_line"] = current["end"]
        n["id"] = symbol_id(n)
        n["file_sha256"] = sha(new_file)
        n["source_url"] = ("https://github.com/" + REPO + "/blob/" +
                           target + "/" + path + "#L" +
                           str(n["start_line"]) + "-L" + str(n["end_line"]))
        n["syntax_calls"] = sorted(set(current.get("calls") or []))
        historical_deps = sorted(oldnode.get(
            "file_level_local_package_dependencies") or [])
        current_deps = sorted(
            x.removeprefix("github.com/" + REPO + "/")
            for x in (deps.get(path) or [])
            if x.startswith("github.com/" + REPO + "/"))
        require(historical_deps == current_deps,
                "G27_IMPORT_CONTEXT_CHANGED: " + path)
        n["file_level_local_package_dependencies"] = sorted(
            x.removeprefix("github.com/" + REPO + "/")
            for x in (deps.get(path) or [])
            if x.startswith("github.com/" + REPO + "/"))
        proof = n["review_evidence"]
        proof.update(source_commit=target, file=path,
                     file_sha256=sha(new_file),
                     start_line=n["start_line"], end_line=n["end_line"])
        for token in proof["literal_tokens"]:
            require(token.encode("utf-8") in new_decl,
                    "rebound token outside new source: " + repr(key))
        id_mapping[oldnode["id"]] = n["id"]
        new_nodes.append(n)
    require(used_receipts == set(receipt_by_key),
            "unused/foreign reviewer receipt")
    revised = copy.deepcopy(old)
    revised["source_commit"] = target
    revised["nodes"] = new_nodes
    revised["coverage"]["symbols"] = len(new_nodes)
    revised["coverage"]["source_files"] = len(files)
    reviewed_count = len(new_nodes) - provisional_added
    require(reviewed_count >= 352, "G28_HISTORICAL_REVIEW_LOST")
    revised["coverage"]["source_reviewed_labels"] = reviewed_count
    revised["coverage"]["source_evidenced_new_labels"] = reviewed_count
    revised["coverage"]["source_dependency_bound_labels"] = reviewed_count
    revised["coverage"]["all_reviewed_source_evidenced_labels"] = reviewed_count
    revised["coverage"]["lexical_provisional_labels"] = provisional_added
    for edge in revised.get("relationships", []):
        if edge.get("kind") != "imports_local_package":
            require(edge["from"] in id_mapping and edge["to"] in id_mapping,
                    "relationship refers to missing symbol")
            edge["from"] = id_mapping[edge["from"]]
            edge["to"] = id_mapping[edge["to"]]
    for machine in revised.get("state_machines", []):
        require(machine["owner_symbol"] in id_mapping,
                "state machine owner disappeared")
        machine["owner_symbol"] = id_mapping[machine["owner_symbol"]]
    new_html = (root / HTML).read_text(encoding="utf-8")
    require(new_html.count(START) == 1, "HTML JSON start ambiguous")
    before, rest = new_html.split(START, 1)
    require(rest.count(END) >= 1, "HTML JSON end missing")
    embedded, tail = rest.split(END, 1)
    require(json.loads(embedded) == old, "HTML no longer matches old atlas")
    serial = json.dumps(revised, ensure_ascii=True, separators=(",", ":"))
    serial = serial.replace("<", "\\u003c")
    viewer = before + START + serial + END + tail
    require(json.loads(viewer.split(START, 1)[1].split(END, 1)[0]) == revised,
            "HTML regenerated machine data mismatch")
    narrative = (root / ICE).read_text(encoding="utf-8")
    orig_marker = "SOURCE_COMMIT: " + previous
    require(narrative.splitlines().count(orig_marker) == 1, "durable header source ambiguous")
    require("STATUS: ADVISORY" in narrative, "durable authority is not advisory")
    narrative = narrative.replace(orig_marker, "SOURCE_COMMIT: " + target, 1)
    if provisional_added:
        narrative = re.sub(r"(?m)^EVIDENCE: .*?$",
                           "EVIDENCE: .ice/code-intelligence.json ("
                           + str(len(new_nodes)) + " source-indexed symbols; "
                           + str(provisional_added) + " provisional)", narrative)
        narrative = re.sub(r"(?m)^COVERAGE: .*?$",
                           "COVERAGE: " + str(len(new_nodes)) + " symbols in "
                           + str(len(files)) + " indexed files; "
                           + str(reviewed_count) + " source-reviewed labels; "
                           + str(provisional_added) + " unreviewed provisional labels.",
                           narrative)
    narrative += (
        "\nG26 SOURCE-AWARE ATLAS MAINTENANCE\n"
        "TARGET_SOURCE_COMMIT: " + target + "\n"
        "PREVIOUS_SOURCE_COMMIT: " + previous + "\n"
        "EXACT_DECLARATION_REUSE_COUNT: " + str(reused) + "\n"
        "FRESH_RECEIPT_COUNT: " + str(fresh) + "\n"
        "FRESH_RECEIPTS_ARE_SELF_ASSERTED: TRUE\n"
        "INDEPENDENT_HUMAN_REVIEW_REQUIRED_FOR_PROMOTION: TRUE\n"
        "Historical evidence is retained as lineage, never silently renewed "
        "for changed declarations. Syntax-only edges remain provisional.\n"
    )
    # Produce all bytes before writing; refuse to overwrite unrelated snapshots.
    encoded = json.dumps(revised, ensure_ascii=False, indent=2) + "\n"
    (root / ATLAS).write_text(encoded, encoding="utf-8")
    (root / HTML).write_text(viewer, encoding="utf-8")
    (root / ICE).write_text(narrative, encoding="utf-8")
    print(json.dumps({
        "result": "G26_SOURCE_AWARE_GENERATION_PASS",
        "target_source_commit": target,
        "previous_source_commit": previous,
        "reused_declarations": reused,
        "fresh_review_receipts": fresh,
        "new_source_bound_provisional": provisional_added,
        "total_declarations": len(new_nodes),
        "independent_approval": False,
        "outputs": [ATLAS, HTML, ICE],
    }, sort_keys=True))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--receipts", type=Path)
    args = parser.parse_args()
    generate(args.root, args.receipts)


if __name__ == "__main__":
    main()
