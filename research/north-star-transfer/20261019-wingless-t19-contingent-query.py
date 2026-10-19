"""T19 frozen TRAIN-only, contingent S3 role-query policy.
Internal diagnostic only: no sealed source, optimizer admission or promotion.
"""
from __future__ import annotations
import importlib.util, itertools, json, statistics, subprocess, sys, time
from pathlib import Path

HERE=Path(__file__).parent
ARMS=("trained-contingent","trained-T16B","untrained-same-cert","full-public-exact")
def load(name,file):
    spec=importlib.util.spec_from_file_location(name,HERE/file)
    module=importlib.util.module_from_spec(spec)
    sys.modules[name]=module
    spec.loader.exec_module(module)
    return module

def deps():
    kernel=load("t19_domain","20261018-wingless-t3-typed-composition.py")
    cert=load("t19_cert","20261019-wingless-t16b-persistent-bindings.py")
    source=load("t19_source","20261019-wingless-t19-s3-source.py")
    return kernel,cert,source

def rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Mem(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
             ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
             ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
             ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
             ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        kernel=ctypes.WinDLL("kernel32");psapi=ctypes.WinDLL("psapi")
        kernel.GetCurrentProcess.restype=wintypes.HANDLE
        psapi.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        psapi.GetProcessMemoryInfo.restype=wintypes.BOOL
        v=Mem();v.cb=ctypes.sizeof(v)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(v),v.cb):
            raise RuntimeError("T19_RSS_UNAVAILABLE")
        return int(v.PeakWorkingSetSize)
    import resource
    x=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(x)*1024 if sys.platform.startswith("linux") else int(x)

def reduced(cert,word):
    return cert.reduce_word(word)
def active(cert,word):
    return tuple(sorted({r for r,inv in reduced(cert,word)}))

def outputs(cert,rows,word,start):
    """Full set of answers under SAME public evidence; TRAIN-only diagnostic."""
    roles=reduced(cert,word)
    a=active(cert,word)
    if not roles:return {start}
    options=[cert.consistent_ops(1,rows[r]) for r in a]
    if any(not x for x in options):raise RuntimeError("T19_INCONSISTENT_PUBLIC")
    out=set()
    for joint in itertools.product(*options):
        mapped=dict(zip(a,joint));value=start
        for role,inv in roles:
            value=cert.apply_op(1,mapped[role],value,inv)
        out.add(value)
    return out

def features(cert,rows,word,role):
    rw=reduced(cert,word)
    n=sum(r==role for r,inv in rw)
    inv=sum(r==role and z for r,z in rw)
    acquired=sum(len(x)==2 for x in rows)
    x,y=rows[role][0]
    return (1.0,float(n),float(inv),float(x==y),
            float(acquired),float(n*acquired),float(role==0),float(role==2))

def fit(cert,train):
    """TRAIN labels are public answer-set contractions, never private truths.
    Include partial views contingent on any one already-revealed role.
    """
    vectors=[]
    for case in train:
        roles=case["roles"];used=active(cert,case["word"])
        initial=[r[:1] for r in roles]
        for prior in (-1,*used):
            view=list(initial)
            if prior>=0:view[prior]=roles[prior]
            before=len(outputs(cert,view,case["word"],case["start"]))
            for role in used:
                if role==prior:continue
                after=list(view);after[role]=roles[role]
                gain=before-len(outputs(cert,after,case["word"],case["start"]))
                vectors.append((features(cert,view,case["word"],role),float(gain)))
    weights=[0.0]*8
    # Immutable TRAIN-only ridge-like online updates and fixed six passes.
    for epoch in range(6):
        for x,target in vectors:
            pred=sum(a*b for a,b in zip(weights,x))
            error=max(-4.0,min(4.0,target-pred))
            denom=1.0+sum(t*t for t in x)
            for j in range(len(weights)):
                weights[j]+=0.25*error*x[j]/denom
    return tuple(round(x,10) for x in weights),len(vectors)

def solve_contingent(cert,view,request,word,start,weights):
    """Does not accept full fixture, private truth, or future public outcomes."""
    rows=[tuple(z) for z in view]
    used=active(cert,word);queries=0
    while True:
        value,proof=cert.shared_role_certificate(1,rows,word,start,budget=128)
        if proof["certified"]:
            return value,queries,True
        pending=[r for r in used if len(rows[r])==1]
        if not pending:break
        def rank(r):
            x=features(cert,rows,word,r)
            return (sum(a*b for a,b in zip(weights,x)),-r)
        nxt=max(pending,key=rank)
        pair=request(nxt)
        rows[nxt]=(rows[nxt][0],pair)
        queries+=1
    value=cert.public_exact_fallback(1,rows,word,start)
    return value,queries,True

def trial(arm,rep):
    assert arm in ARMS and 0<=rep<5
    t0=time.process_time_ns()
    kernel,cert,source=deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert not set(source.TRAIN)&set(source.DEV)
    train=[source.build(label,kernel) for label in source.TRAIN]
    weights=(0.0,)*8;updates=0
    if arm=="trained-contingent":
        weights,updates=fit(cert,train)
    rank=cert.score_role_requests([(x["domain"],x["roles"],x["word"],x["start"]) for x in train]) if arm=="trained-T16B" else (0,0,0,0)
    rows=[]
    for label in source.DEV:
        case=source.build(label,kernel)
        if arm=="trained-contingent":
            visible,request,charged=source.observed_interface(case)
            answer,queries,sound=solve_contingent(cert,visible,request,case["word"],case["start"],weights)
            assert queries==charged()
        elif arm=="full-public-exact":
            answer=cert.public_exact_fallback(1,case["roles"],case["word"],case["start"])
            queries=len(active(cert,case["word"]));sound=True
        else:
            answer,meta=cert.bound_query_solve(1,case["roles"],case["word"],case["start"],rank,cap=128)
            queries=meta["queries"];sound=meta["sound"]
        good=(answer==case["private_truth"])
        if not good or not sound:raise RuntimeError("T19_HARD_FALSE_POSITIVE_OR_BASELINE_LOSS:"+label)
        rows.append({"id":label,"correct":good,"extra_public_requests":queries})
    assert len(rows)==64
    return {"schema":"wingless.t19.fresh-internal-s3-per-process.v1","arm":arm,
            "repetition":rep,"trained_cases":40,"dev_cases":64,"rows":rows,
            "updates":updates,"coefficients":weights if arm=="trained-contingent" else None,
            "complete_cpu_ns":time.process_time_ns()-t0,"process_peak_rss_bytes":rss(),
            "source":"SELF_AUTHORED_NEW_INTERNAL_DEV_ONLY","T16_originals_opened":0,
            "real_remote_query_service_cost":"UNQUALIFIED","scientific_acceptance":False,
            "optimizer_admission":False}

def evaluate():
    allrows=[]
    for arm in ARMS:
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),"--arm",arm,str(rep)],
                             capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T19_CHILD_FAILURE "+arm+" rep="+str(rep)+"\n"+p.stderr)
            allrows.append(json.loads(p.stdout))
    out={}
    for arm in ARMS:
        xs=[x for x in allrows if x["arm"]==arm]
        assert len(xs)==5 and all(z["rows"]==xs[0]["rows"] for z in xs)
        out[arm]={"correct":sum(x["correct"] for x in xs[0]["rows"]),
                  "queries":sum(x["extra_public_requests"] for x in xs[0]["rows"]),
                  "rows":xs[0]["rows"],
                  "cpu_five_process_ns":[x["complete_cpu_ns"] for x in xs],
                  "rss_five_process_bytes":[x["process_peak_rss_bytes"] for x in xs],
                  "median_complete_cpu_ns":int(statistics.median(x["complete_cpu_ns"] for x in xs)),
                  "maximum_peak_rss_bytes":max(x["process_peak_rss_bytes"] for x in xs),
                  "updates":xs[0]["updates"],"coefficients":xs[0]["coefficients"]}
    learned=out["trained-contingent"];trained=out["trained-T16B"]
    comparisons={k:{"strict_query_savings":learned["queries"]<v["queries"],
                    "complete_cpu_nonregression":learned["median_complete_cpu_ns"]<=v["median_complete_cpu_ns"],
                    "peak_rss_nonregression":learned["maximum_peak_rss_bytes"]<=v["maximum_peak_rss_bytes"],
                    "correctness_nonregression":learned["correct"]>=v["correct"]}
                 for k,v in out.items() if k!="trained-contingent"}
    assert all(v["correct"]==64 for v in out.values())
    return {"schema":"wingless.t19.contingent-s3-train-only-development.v1",
      "source":"NEW S3 SELF_AUTHORED INTERNAL ONLY 40 TRAIN 64 DEV",
      "results":out,"comparison":comparisons,
      "trained_query_advantage_observed":learned["queries"]<trained["queries"],
      "all_measured_resource_gates_pass":all(all(x.values()) for x in comparisons.values()),
      "remote_service_cost_qualified":False,
      "T16_originals_consumed":0,"independent_science":False,
      "scientific_acceptance":False,"optimizer_admission":False,
      "no_new_controller":True,"no_merges":True}

def main():
    if len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(evaluate(),sort_keys=True))
    else:raise SystemExit("T19_INVALID_ARGS")
if __name__=="__main__":main()
