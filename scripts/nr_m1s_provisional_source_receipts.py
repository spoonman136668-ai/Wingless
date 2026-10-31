#!/usr/bin/env python3
"""Build unreviewed, source-bound G28 provisional receipts for NEW Go declarations.

Must run after the existing canonical Go ICE builder against exact HEAD. The
approved, unmodified scripts/maintain_code_atlas.py independently verifies all
source bytes. This script NEVER declares independent human review.
"""
import argparse
import hashlib
import json
import subprocess
from pathlib import Path

def require(ok,message):
    if not ok:raise SystemExit("NRM1S_UNQUALIFIED_SOURCE_RECEIPT: "+message)

def git(*args):
    p=subprocess.run(["git",*args],capture_output=True)
    require(p.returncode==0,"git "+repr(args)+" "+p.stderr.decode("utf8","replace")[-200:])
    return p.stdout

def key(s):
    return (s["file"],s["kind"],s["name"],s.get("receiver",""))

def main():
    parser=argparse.ArgumentParser()
    parser.add_argument("--output",type=Path,required=True)
    x=parser.parse_args()
    target=git("rev-parse","HEAD").decode("ascii").strip()
    require(len(target)==40 and all(c in "0123456789abcdef" for c in target),"invalid source head")
    old=json.loads(Path(".ice/code-intelligence.json").read_text("utf8"))
    # The reviewed source must be an ACTUAL immutable ancestor, not a
    # fabricated review continuity or unrelated historical hash.
    reviewed_source=old.get("source_commit","")
    require(len(reviewed_source)==40 and
            all(c in "0123456789abcdef" for c in reviewed_source),
            "invalid previous reviewed source identity")
    relation=subprocess.run(["git","merge-base","--is-ancestor",reviewed_source,target])
    require(relation.returncode==0,
            "previous reviewed source is NOT an ancestor of this exact head")
    original={key(s) for s in old["nodes"]}
    require(len(original)==len(old["nodes"]) and len(original)>=352,"lost historical reviewed nodes")
    symbols=json.loads(Path(".ice/symbol-map.json").read_text("utf8"))
    current={key(s):s for s in symbols}
    require(len(current)==len(symbols),"duplicate new Go declarations")
    require(original.issubset(current),"a historically reviewed symbol disappeared/was renamed")
    manifest=json.loads(Path(".ice/manifest.json").read_text("utf8"))
    hashes={f["path"]:f["sha256"] for f in manifest["files"]}
    receipts=[]
    for k in sorted(set(current)-original):
        s=current[k]
        path=s["file"]
        require(path in hashes and path.endswith(".go"),"new symbol outside indexed Go source: "+path)
        data=Path(path).read_bytes()
        require(hashlib.sha256(data).hexdigest()==hashes[path],"ICE file hash mismatch: "+path)
        require(git("show",target+":"+path).replace(b"\r\n",b"\n")==
                data.replace(b"\r\n",b"\n"),"PR head source mismatch: "+path)
        lines=data.splitlines(keepends=True)
        start,end=s["start"],s["end"]
        require(isinstance(start,int) and isinstance(end,int) and
                1<=start<=end<=len(lines),"invalid symbol range: "+repr(k))
        decl=b"".join(lines[start-1:end])
        token=next((w for w in ("func","type","const","import")
                    if w.encode() in decl),None)
        if token is None and len(s["name"])>=4 and s["name"].encode() in decl:
            token=s["name"]
        require(token is not None,"cannot derive lexical source witness: "+repr(k))
        receipts.append({
            "key":list(k),
            "target_source_commit":target,
            "file_sha256":hashlib.sha256(data).hexdigest(),
            "declaration_sha256":hashlib.sha256(decl).hexdigest(),
            "start_line":start,
            "end_line":end,
            "reviewer":"AUTOMATED_SOURCE_ONLY_NO_HUMAN_REVIEW",
            "label":str(s["kind"])+" "+str(s["name"])+" unreviewed native source declaration",
            "description":"Unreviewed source-indexed declaration; this is only lexical provenance, with no independently authenticated behavioral review.",
            "literal_tokens":[token]
        })
    x.output.parent.mkdir(parents=True,exist_ok=True)
    x.output.write_text(json.dumps(receipts,indent=2,sort_keys=True)+"\n","utf8")
    print(json.dumps({"status":"NRM1S_PROVISIONAL_RECEIPTS_SOURCE_BOUND_NOT_REVIEWED",
      "head":target,"new_provisional_symbols":len(receipts),
      "historical_reviewed_symbols":len(original),"independent_approval":False},sort_keys=True))
if __name__=="__main__":
    main()
