"""Wingless F61B source-only acceptance harness for EXISTING qualified workers.

This entrypoint never generates FIT, DIAG or original experiments and does not
dispatch controllers, scientific runs or independent acceptance.
Usage: python3 -B research/foundational/20261015-w-f61/qualify_source.py --expect-head EXACT_SHA
"""
from __future__ import annotations
import argparse, hashlib, json, pathlib, subprocess, sys
HERE=pathlib.Path(__file__).resolve().parent
ROOT=HERE.parents[2]
FROZEN={
    "research/foundational/20261015-w-f61/prereg.json":"d6fd61c0e8062ed27c19da63ff196ad08a5e6ed6",
    "research/foundational/20261015-w-f61/prereg-amendment.json":"08222165f73f4e9dfd1546b601dacdd719724bf1",
    "research/foundational/20261015-w-f61/source_feasibility.py":"f9372cf2cbba75f6b90add7ca388b607d9ee49bc",
    "research/foundational/20261015-w-f61/audit_fit.py":"f5fc5f0204280160ac8736f9777bc7a32f58c425",
    "research/foundational/20261015-w-f61/test_preflight.py":"72d01dce16f9c4d76f8dfce616fcf1587fd42933",
    "research/foundational/20261015-w-f60/source_feasibility.py":"f6863c7fdf31bd2cd20daaf6434640596b6ba5a4",
}
def git_blob(path):
    data=path.read_bytes()
    return hashlib.sha1(b"blob "+str(len(data)).encode()+bytes([0])+data).hexdigest()
def cmd(*args):
    return subprocess.check_output(args,cwd=ROOT,text=True).strip()
def run(stage,args):
    print("W_F61B_GATE_START="+stage,flush=True)
    subprocess.run(args,cwd=ROOT,check=True)
    print("W_F61B_GATE_PASS="+stage,flush=True)
def main():
    p=argparse.ArgumentParser()
    p.add_argument("--expect-head",required=True)
    args=p.parse_args()
    assert len(args.expect_head)==40 and all(x in "0123456789abcdef" for x in args.expect_head)
    assert cmd("git","rev-parse","HEAD")==args.expect_head,"HEAD_DRIFT"
    assert not cmd("git","status","--porcelain"),"DIRTY_WORKTREE"
    subprocess.run(["git","merge-base","--is-ancestor",
      "4f6fbfab77ab36a319742d0d5bc15dd920ae1c2e","HEAD"],cwd=ROOT,check=True)
    for name,blob in FROZEN.items():
        assert git_blob(ROOT/name)==blob,"BLOB_DRIFT: "+name
    old=json.loads((HERE/"prereg.json").read_text(encoding="utf-8"))
    amendment=json.loads((HERE/"prereg-amendment.json").read_text(encoding="utf-8"))
    assert old["stage"]=="PRECODE_FROZEN_NO_FIT_USE"
    assert old["fit"]["seeds"]==[751000,751016]
    assert amendment["original_prereg_blob_sha"]==FROZEN["research/foundational/20261015-w-f61/prereg.json"]
    assert amendment["fit"]["seeds"]==[771000,771016]
    assert amendment["diagnostic_UNAUTHORIZED"]["seeds"]==[771100,771116]
    assert amendment["original_UNRESERVED"]["seeds"]==[772000,772016]
    assert amendment["original_F61_cohort_status"].startswith("not pristine")
    assert amendment["source_gates"]["minimum_additional_correct_of_128"]==7
    assert amendment["independent_F3"] is False
    run("FOCUSED_F61B_SYNTHETIC_ONLY",[
       sys.executable,"-B","-m","unittest","discover","-s",str(HERE),"-p","test_preflight.py","-v"])
    run("ONE_FULL_GO_REGRESSION",["go","test","-count=1","./..."])
    run("GO_VET",["go","vet","./..."])
    run("GO_BUILD",["go","build","./..."])
    print("W_F61B_SOURCE_GATE=PASS",flush=True)
    print(json.dumps(dict(schema="wingless.f61b.source-qualification.v1",
        head=args.expect_head,stage="SOURCE_ONLY",pinned_blobs=len(FROZEN),
        full_regressions=1,FIT_executed=False,DIAG_executed=False,
        original_executed=False,independent_F3=False,
        science_eligible_without_owner_authorization=False),sort_keys=True))
if __name__=="__main__":
    try:main()
    except (AssertionError,subprocess.CalledProcessError,FileNotFoundError) as exc:
        print("W_F61B_SOURCE_GATE=FAIL",repr(exc),file=sys.stderr)
        sys.exit(1)
