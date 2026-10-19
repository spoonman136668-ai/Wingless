"""T18 internal query-identifiability experiment, separate from T16 originals.

Three-role frozen source; genuine strongest trained T16B control and matching
untrained exact certifier. Clairvoyant minimum-query subset is a NONELIGIBLE
hindsight bound computed only after per-arm runs, never fed to policies.
"""
import importlib.util
import itertools
import json
import statistics
import subprocess
import sys
import time
from pathlib import Path

D=Path(__file__).parent
ARMS=("trained-t16b","untrained-same-cert","full-exact")

def load(name,filename):
    spec=importlib.util.spec_from_file_location(name,D/filename)
    obj=importlib.util.module_from_spec(spec)
    sys.modules[name]=obj
    spec.loader.exec_module(obj)
    return obj

def deps():
    k=load("t18_k","20261018-wingless-t3-typed-composition.py")
    c=load("t18_cert","20261019-wingless-t16b-persistent-bindings.py")
    s=load("t18_src","20261019-wingless-t18-headroom-source.py")
    return k,c,s

def peak_rss_bytes():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Pmc(ctypes.Structure):
            _fields_=[
             ("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
             ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
             ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),
             ("QuotaPagedPoolUsage",ctypes.c_size_t),
             ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),
             ("QuotaNonPagedPoolUsage",ctypes.c_size_t),
             ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        kernel=ctypes.WinDLL("kernel32")
        psapi=ctypes.WinDLL("psapi")
        kernel.GetCurrentProcess.restype=wintypes.HANDLE
        psapi.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Pmc),wintypes.DWORD]
        psapi.GetProcessMemoryInfo.restype=wintypes.BOOL
        c=Pmc();c.cb=ctypes.sizeof(c)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(c),c.cb):
            raise OSError("T18_WIN32_PEAK_WORKING_SET_UNAVAILABLE")
        return int(c.PeakWorkingSetSize)
    import resource
    x=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(x)*1024 if sys.platform.startswith("linux") else int(x)

def admissible_public(f):
    return f["domain"],f["roles"],f["word"],f["start"]

def cases(s,k,ids):
    return [s.build(label,k) for label in ids]

def arm_trial(name,trial):
    assert name in ARMS and trial in range(5)
    tick=time.process_time_ns()
    k,c,s=deps()
    train=cases(s,k,s.TRAIN)
    # Controls use exactly the same TRAIN public role calibration data.
    weights=c.score_role_requests([admissible_public(f) for f in train]) if name=="trained-t16b" else (0,0,0,0)
    dev=cases(s,k,s.DEV)
    rows=[]
    for f in dev:
        domain,roles,word,start=admissible_public(f)
        if name=="full-exact":
            answer=c.public_exact_fallback(domain,roles,word,start)
            queries=len(set(j for j,_ in c.reduce_word(word)))
        else:
            answer,diagnostic=c.bound_query_solve(domain,roles,word,start,weights,cap=128)
            queries=diagnostic["queries"]
        rows.append({"case":f["source_id"],"domain":domain,
                     "correct":int(answer==f["private_truth"]),
                     "queries":queries})
    assert len(rows)==48 and all(x["correct"] for x in rows), "T18_BASELINE_POSITIVE_LOSS"
    return {"schema":"wingless.t18.internal-per-process.v1",
      "arm":name,"trial":trial,"train_cases":len(train),"dev_cases":len(dev),
      "rows":rows,"cpu_ns":time.process_time_ns()-tick,
      "whole_process_peak_rss_bytes":peak_rss_bytes(),
      "source_only_internal":True,"external_firstuse":False}

def hindsight_min_queries(c,f):
    domain,roles,word,start=admissible_public(f)
    active=tuple(sorted({r for r,_ in c.reduce_word(word)}))
    for n in range(len(active)+1):
        for selected in itertools.combinations(active,n):
            evidence=[role[:1] for role in roles]
            for r in selected:evidence[r]=roles[r]
            value,proof=c.shared_role_certificate(domain,evidence,word,start,budget=4096)
            if proof["certified"] and value==f["private_truth"]:
                return n
    raise AssertionError("T18_NO_FULL_EVIDENCE_CERTIFICATE")

def evaluate():
    trials=[]
    for name in ARMS:
        for i in range(5):
            cmd=[sys.executable,str(Path(__file__).resolve()),"--arm",name,str(i)]
            p=subprocess.run(cmd,capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T18_FOCUSED_CHILD_FAILED "+name+" trial "+str(i)+"\n"+p.stderr)
            trials.append(json.loads(p.stdout))
    metrics={}
    for name in ARMS:
        xs=[x for x in trials if x["arm"]==name]
        one=xs[0]["rows"]
        assert len(xs)==5 and all(x["rows"]==one for x in xs)
        metrics[name]={
          "correct":sum(x["correct"] for x in one),
          "queries":sum(x["queries"] for x in one),
          "per_case":one,
          "five_process_cpu_ns":[x["cpu_ns"] for x in xs],
          "five_process_rss_bytes":[x["whole_process_peak_rss_bytes"] for x in xs],
          "median_full_cpu_ns":int(statistics.median(x["cpu_ns"] for x in xs)),
          "maximum_whole_process_rss_bytes":max(x["whole_process_peak_rss_bytes"] for x in xs)
        }
    k,c,s=deps()
    dev=cases(s,k,s.DEV)
    lower={f["source_id"]:hindsight_min_queries(c,f) for f in dev}
    incumbent={r["case"]:r for r in metrics["trained-t16b"]["per_case"]}
    oracle_better=sum(lower[x]<r["queries"] for x,r in incumbent.items())
    bydomain={}
    for f in dev:
        group=bydomain.setdefault("F11" if f["domain"]==0 else "S3",
                                  {"cases":0,"hindsight_strict_headroom":0,
                                   "control_queries":0,"hindsight_queries":0})
        group["cases"]+=1
        group["hindsight_strict_headroom"]+=lower[f["source_id"]]<incumbent[f["source_id"]]["queries"]
        group["control_queries"]+=incumbent[f["source_id"]]["queries"]
        group["hindsight_queries"]+=lower[f["source_id"]]
    return {
      "schema":"wingless.t18.query-identifiability-new-synthetic-diagnostic.v1",
      "source":"48 NEW internal DEV, 24 NEW internal TRAIN, 3 roles",
      "frozen_controls":ARMS,"results":metrics,
      "per_case_noneligible_hindsight_lower_bound":lower,
      "hindsight_oracle_not_deployable":True,
      "hindsight_strict_headroom_cases":oracle_better,
      "hindsight_strict_headroom_fraction":oracle_better/len(dev),
      "headroom_at_least_10pct":oracle_better>=5,
      "per_domain":bydomain,
      "real_public_query_service_cpu_not_metered":True,
      "source_independent_external_birth":False,
      "T16_protected_originals_opened":0,
      "scientific_acceptance":False,
      "optimizer_admission":False
    }

def main():
    if len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(arm_trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(evaluate(),sort_keys=True))
    else:raise SystemExit("T18_INVALID_ARGUMENTS")

if __name__=="__main__":main()
