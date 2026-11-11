#!/usr/bin/env python3
"""NR-M1S: strictly read-only frozen-source reviewed-declaration delta.

This does not authorize fresh review, infer behavior, or modify reviewed
semantics. Produces a machine-readable list of exact changed declarations.
"""
import argparse
import hashlib
import json
import subprocess
from pathlib import Path

def git(*args):
    r=subprocess.run(["git",*args],capture_output=True)
    if r.returncode:raise SystemExit("NRM1S_MISSING_HISTORICAL_SOURCE: "+repr(args)+" "+r.stderr.decode("utf8","replace")[-300:])
    return r.stdout

def sha(x):return hashlib.sha256(x).hexdigest()
def key(x):return x["file"],x["kind"],x["name"],x.get("receiver","")
def span(raw,lo,hi):
    lines=raw.splitlines(keepends=True)
    if not isinstance(lo,int) or not isinstance(hi,int) or not 1<=lo<=hi<=len(lines):
        raise ValueError("NRM1S_SYMBOL_SPAN_OUT_OF_RANGE")
    return b"".join(lines[lo-1:hi]).replace(b"\r\n",b"\n")

def main():
    ap=argparse.ArgumentParser()
    ap.add_argument("--output",type=Path,required=True)
    args=ap.parse_args()
    prior=json.loads(Path(".ice/code-intelligence.json").read_text("utf8"))
    current=json.loads(Path(".ice/symbol-map.json").read_text("utf8"))
    original={key(x):x for x in prior["nodes"]}
    current_map={key(x):x for x in current}
    if len(original)!=len(prior["nodes"]) or len(current_map)!=len(current):
        raise ValueError("NRM1S_DUPLICATE_DECLARATIONS")
    if len(original)<352 or not set(original).issubset(current_map):
        raise ValueError("NRM1S_HISTORICAL_REVIEW_DISAPPEARED")
    source=prior["source_commit"]
    head=git("rev-parse","HEAD").decode().strip()
    # All previous-reviewed files are checked against their original Git bytes.
    result=[]
    historical={}
    for k,old in sorted(original.items()):
        path=old["file"]
        if path not in historical:
            historical[path]=git("show",source+":"+path)
        oldbytes=historical[path]
        if sha(oldbytes)!=old["file_sha256"]:
            raise ValueError("NRM1S_PREVIOUS_REVIEW_FILE_SHA_MISMATCH "+path)
        new=current_map[k]
        present=Path(path).read_bytes()
        target=git("show",head+":"+path)
        if present.replace(b"\r\n",b"\n")!=target.replace(b"\r\n",b"\n"):
            raise ValueError("NRM1S_NEW_SOURCE_NOT_COMMITTED "+path)
        before=span(oldbytes,old["start_line"],old["end_line"])
        after=span(present,new["start"],new["end"])
        if before!=after:
            result.append({
              "file":path,"kind":k[1],"name":k[2],"receiver":k[3],
              "previous_source_commit":source,"current_source_commit":head,
              "previous_start_line":old["start_line"],"previous_end_line":old["end_line"],
              "current_start_line":new["start"],"current_end_line":new["end"],
              "prior_declaration_sha256":sha(before),"current_declaration_sha256":sha(after),
              "independent_review":False,"approved":False,
              "disposition":"BLOCKS_G28_UNTIL_INDEPENDENT_REVIEW_OR_AUTHORIZED_NATIVE_SOURCE_RECONCILIATION"
            })
    out={"schema":"nr-m1s-exact-declaration-review-delta-v1",
         "prior_reviewed_source_commit":source,"candidate_head":head,
         "historical_reviewed_symbols":len(original),
         "current_symbols":len(current_map),
         "new_unreviewed_symbols":len(current_map)-len(original),
         "modified_previously_reviewed_symbols":len(result),
         "independent_review_completed":False,
         "gate":"BLOCKED" if result else "REVIEWED_DECLARATIONS_BYTE_IDENTICAL",
         "modified":result}
    args.output.parent.mkdir(parents=True,exist_ok=True)
    args.output.write_text(json.dumps(out,sort_keys=True,indent=2)+"\n","utf8")
    print(json.dumps({k:v for k,v in out.items() if k!="modified"},sort_keys=True))
    for x in result[:100]:
        print("NRM1S_INDEPENDENT_REVIEW_REQUIRED",x["file"],x["name"],x["prior_declaration_sha256"],x["current_declaration_sha256"])
if __name__=="__main__":main()
