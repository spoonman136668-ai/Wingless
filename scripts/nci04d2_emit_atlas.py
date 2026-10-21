#!/usr/bin/env python3
"""Read-only source-generated G28 PROVISIONAL receipts and ICE bundle emission.

Runs only in a disposable exact-head CI checkout. The reviewer assertion is
explicitly AUTOMATED AND UNVERIFIED, not human approval or scientific evidence.
Never mutates the GitHub repository or accepted branch. Emits proposal bytes
for independent source comparison and later git commit.
"""
import base64
import hashlib
import json
import pathlib
import re
import subprocess
import sys

CANONICAL=[
 ".ice/manifest.json",".ice/architecture-map.json",".ice/symbol-map.json",
 ".ice/dependency-map.json",".ice/test-map.json",".ice/implementations.json",
 ".ice/code-intelligence.json",".ice/code-intelligence.html",
 ".ice/code-intelligence.ice"
]

def run(args):
    p=subprocess.run(args,check=True,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
    return p.stdout

def main():
    root=pathlib.Path(__file__).resolve().parents[1]
    import os
    os.chdir(root)
    target=run(["git","rev-parse","HEAD"]).decode().strip()
    if not re.fullmatch("[0-9a-f]{40}",target):
        raise ValueError("target git commit invalid")
    # A source-only artifact commit legitimately follows its pinned source birth.
    # If the canonical index is already current, emit no replacement artifact.
    ready=subprocess.run(["go","run","./cmd/ice","validate"],stdout=subprocess.PIPE,stderr=subprocess.PIPE)
    if ready.returncode==0:
        print("NCI04D2_CANONICAL_ICE_CURRENT "+target)
        return
    run(["go","run","./cmd/ice","build"])
    old=json.loads((root/".ice/code-intelligence.json").read_text("utf-8"))
    current=json.loads((root/".ice/symbol-map.json").read_text("utf-8"))
    indexed=json.loads((root/".ice/manifest.json").read_text("utf-8"))
    oldkeys={(n["file"],n["kind"],n["name"],n.get("receiver","")) for n in old["nodes"]}
    old_reviewed=sum(n["evidence_basis"]=="source-reviewed" for n in old["nodes"])
    old_provisional=sum(n["evidence_basis"]=="lexical-provisional" for n in old["nodes"])
    assert len(oldkeys)==len(old["nodes"])>=352
    assert old_reviewed==352 and old_provisional==len(old["nodes"])-352
    assert len(current)>len(oldkeys)
    filehash={f["path"]:f["sha256"] for f in indexed["files"]}
    receipts=[]
    for n in current:
        key=(n["file"],n["kind"],n["name"],n.get("receiver",""))
        if key in oldkeys:
            continue
        path=n["file"]
        raw=(root/path).read_bytes()
        git_raw=run(["git","show",target+":"+path])
        if raw!=git_raw or hashlib.sha256(raw).hexdigest()!=filehash[path]:
            raise ValueError("new source not committed or digest mismatch "+path)
        lines=raw.splitlines(keepends=True)
        declaration=b"".join(lines[n["start"]-1:n["end"]])
        decoded=declaration.decode("utf-8")
        matches=re.findall(r"[A-Za-z_][A-Za-z_0-9]{3,}",decoded)
        if not matches:
            raise ValueError("symbol missing literal source witness "+path)
        receipts.append({
         "key":list(key),"target_source_commit":target,
         "file_sha256":filehash[path],
         "declaration_sha256":hashlib.sha256(declaration).hexdigest(),
         "start_line":n["start"],"end_line":n["end"],
         "reviewer":"automated-source-witness-NOT-INDEPENDENT",
         "label":"Provisional source symbol "+n["name"],
         "description":"Lexical Go declaration, without independently verified behavior",
         "literal_tokens":[matches[0]]
        })
    if not receipts:
        raise ValueError("expected source bound additions")
    target_file=root/".ice/nci04d2_provisional_receipts.json"
    target_file.write_text(json.dumps(receipts,sort_keys=True),encoding="utf-8")
    subprocess.run([sys.executable,"scripts/maintain_code_atlas.py","--root",".","--receipts",str(target_file)],
                   check=True)
    atlas=json.loads((root/".ice/code-intelligence.json").read_text("utf-8"))
    if atlas["source_commit"]!=target or atlas["coverage"]["source_reviewed_labels"]!=352:
        raise ValueError("source pin or historical reviewed baseline lost")
    total_provisional=old_provisional+len(receipts)
    if atlas["coverage"]["lexical_provisional_labels"]!=total_provisional:
        raise ValueError("incorrect inherited+new provisional evidence accounting")
    artifacts={}
    for path in CANONICAL:
        artifacts[path]=(root/path).read_text(encoding="utf-8")
    obj={"schema":"wingless.nci04d2.provisional-index-bundle.v1",
         "source_commit":target,"reviewed_count":352,
         "provisional_count":total_provisional,
         "independent_review":False,"artifacts":artifacts}
    b64=base64.b64encode(json.dumps(obj,ensure_ascii=True).encode()).decode("ascii")
    parts=[b64[i:i+14000] for i in range(0,len(b64),14000)]
    print("NCI04D2_BUNDLE_CHUNKS="+str(len(parts)))
    for i,part in enumerate(parts):
        print("NCI04D2_BUNDLE_CHUNK:"+str(i)+":"+part)
    print("NCI04D2_PROVISIONAL_SOURCE_ONLY_PASS "+target+" "+str(len(receipts)))

if __name__=="__main__":
    try:main()
    except Exception as exc:
        print("NCI04D2_PROVISIONAL_BUNDLE_FAILED "+str(exc),file=sys.stderr)
        sys.exit(1)
