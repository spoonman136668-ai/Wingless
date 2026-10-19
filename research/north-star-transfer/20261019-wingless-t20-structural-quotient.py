"""T20 frozen-development structural quotient: sound exact live-binding transducer.

The learned parameter is a TRAIN-only structural-shape compilation schedule.
It never certifies an answer without exact symbolic state-set agreement.
No T16 original, T19 DEV, independent first-use or optimizer admission.
"""
from __future__ import annotations
import importlib.util
import json
import statistics
import subprocess
import sys
import time
from pathlib import Path

HERE=Path(__file__).parent
MODES=("always","every-two","end-only")
ARMS=("trained-quotient","fixed-quotient","trained-T16B","full-public-exact")
STATE_CAP=128

def load(name,filename):
    spec=importlib.util.spec_from_file_location(name,HERE/filename)
    obj=importlib.util.module_from_spec(spec)
    sys.modules[name]=obj
    spec.loader.exec_module(obj)
    return obj

def deps():
    return (load("t20_kernel","20261018-wingless-t3-typed-composition.py"),
            load("t20_cert","20261019-wingless-t16b-persistent-bindings.py"),
            load("t20_source","20261019-wingless-t20-fresh-source.py"))

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
        x=Mem();x.cb=ctypes.sizeof(x)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(x),x.cb):
            raise RuntimeError("T20_RSS_UNAVAILABLE")
        return int(x.PeakWorkingSetSize)
    import resource
    value=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(value)*1024 if sys.platform.startswith("linux") else int(value)

def active(cert,word):
    return tuple(sorted({r for r,_ in cert.reduce_word(word)}))

def shape(cert,word):
    rw=cert.reduce_word(word)
    counts=tuple(sum(role==j for role,_ in rw) for j in range(3))
    return (len(rw)//3,tuple(sorted(counts)),sum(bool(i) for _,i in rw)//3)

def quotient_certificate(cert,rows,word,start,mode="always",cap=STATE_CAP):
    """Enumerate exact public-consistent role bindings through a quotient
    on output state and only *still-live* registers. No unsafe merging.
    """
    if mode not in MODES or cap<1:
        raise ValueError("T20_INVALID_COMPILER")
    sequence=cert.reduce_word(word)
    if not sequence:
        return start,{"certified":True,"reason":"IDENTITY","work":0,
                      "merge_work":0,"max_frontier":1}
    used={role for role,_ in sequence}
    options={j:cert.consistent_ops(1,rows[j]) for j in used}
    if any(not options[j] for j in used):
        return None,{"certified":False,"reason":"INCONSISTENT_PUBLIC",
                     "work":0,"merge_work":0,"max_frontier":0}
    frontier=[(start,(None,None,None))]
    work=0;merge_work=0;max_frontier=1
    for i,(role,inv) in enumerate(sequence):
        future={r for r,_ in sequence[i+1:]}
        nxt=[]
        for value,bindings in frontier:
            possible=(bindings[role],) if bindings[role] is not None else options[role]
            for op in possible:
                b=list(bindings)
                b[role]=op
                result=cert.apply_op(1,op,value,inv)
                # Forget only identities absent from every remaining suffix.
                for j in range(3):
                    if j not in future:
                        b[j]=None
                nxt.append((result,tuple(b)))
        work+=len(nxt)
        collapse=(mode=="always" or (mode=="every-two" and i%2==1)
                  or i==len(sequence)-1)
        if collapse:
            merge_work+=len(nxt)
            frontier=list(set(nxt))
        else:
            frontier=nxt
        max_frontier=max(max_frontier,len(frontier))
        if len(frontier)>cap:
            return None,{"certified":False,"reason":"BOUNDED_UNKNOWN",
                         "work":work,"merge_work":merge_work,
                         "max_frontier":max_frontier}
    outcomes={value for value,_ in frontier}
    if len(outcomes)==1:
        return next(iter(outcomes)),{"certified":True,"reason":"EXACT_QUOTIENT",
              "work":work,"merge_work":merge_work,"max_frontier":max_frontier}
    return None,{"certified":False,"reason":"MULTIPLE_VALID_OUTCOMES",
             "work":work,"merge_work":merge_work,"max_frontier":max_frontier}

def prove_sound_on_public_algebra(cert):
    """Non-source mathematical adversarial coverage; never DEV data."""
    from itertools import product
    operators=cert.PERMS
    checks=0
    # Distinct nontrivial word grammars, repeat bindings and inverse edges.
    words=[]
    for n in (1,2,3,4):
        for w in product(tuple(product(range(3),(False,True))),repeat=n):
            if n==4 and len(words)>=300:
                break
            words.append(w)
    for w in words:
        for r in range(3):
            # Each first-pair row has two legal S3 operators; no oracle value.
            rows=[((r,operators[(j*2+r)%6][r]),) for j in range(3)]
            if any(not cert.consistent_ops(1,rows[j]) for j in range(3)):
                raise AssertionError("T20_ILLEGAL_PUBLIC_FIXTURE")
            exact,reference=cert.shared_role_certificate(1,rows,w,r,budget=4096)
            for mode in MODES:
                candidate,proof=quotient_certificate(cert,rows,w,r,mode,cap=STATE_CAP)
                if proof["reason"]=="BOUNDED_UNKNOWN":
                    continue
                if proof["certified"]!=reference["certified"] or candidate!=exact:
                    raise AssertionError("T20_QUOTIENT_UNSOUND:"+str((w,r,mode,proof)))
                checks+=1
    if checks<100:
        raise AssertionError("T20_INSUFFICIENT_ALGEBRA_COVERAGE")
    return checks

def fit_compiler(cert,train):
    """Learn a bounded structural-shape-to-safe-schedule table from public
    TRAIN calibration pairs and words, never hidden truth or DEV labels.
    """
    totals={};updates=0
    for case in train:
        sig=shape(cert,case["word"])
        values=totals.setdefault(sig,{m:0 for m in MODES})
        views=[tuple((r[0],) for r in case["roles"]),case["roles"]]
        for view in views:
            for mode in MODES:
                _val,proof=quotient_certificate(cert,view,case["word"],
                                                case["start"],mode)
                # Deterministic public surrogate, includes compression work.
                values[mode]+=proof["work"]+proof["merge_work"]
                updates+=1
    policy={}
    for sig,costs in totals.items():
        policy[sig]=min(MODES,key=lambda m:(costs[m],MODES.index(m)))
    return policy,updates

def solve_quotient(cert,visible,request,word,start,rank,mode):
    rows=[tuple(p) for p in visible];queries=0;proofs=0
    while True:
        value,proof=quotient_certificate(cert,rows,word,start,mode)
        proofs+=1
        if proof["certified"]:
            return value,queries,proofs
        remaining=[r for r in active(cert,word) if len(rows[r])==1]
        if not remaining:
            break
        selected=max(remaining,key=lambda j:(rank[j],-j))
        rows[selected]=(rows[selected][0],request(selected))
        queries+=1
    # Exhaustive verified fallback, only after all active rows are public.
    return cert.public_exact_fallback(1,rows,word,start),queries,proofs

def child(arm,repetition):
    assert arm in ARMS and 0<=repetition<5
    start_ns=time.process_time_ns()
    kernel,cert,source=deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert not set(source.TRAIN)&set(source.DEV)
    train=[source.build(label,kernel) for label in source.TRAIN]
    rank=cert.score_role_requests([(x["domain"],x["roles"],x["word"],x["start"])
                                    for x in train])
    compiler={};updates=0
    if arm=="trained-quotient":
        compiler,updates=fit_compiler(cert,train)
    out=[];total_proofs=0
    for label in source.DEV:
        case=source.build(label,kernel)
        public,request,charged=source.observed_interface(case)
        if arm in ("trained-quotient","fixed-quotient"):
            mode=compiler.get(shape(cert,case["word"]),"always") if arm=="trained-quotient" else "always"
            answer,queries,proofs=solve_quotient(cert,public,request,case["word"],
                                                  case["start"],rank,mode)
        elif arm=="trained-T16B":
            # Existing trained exact control uses the same observed pair rows.
            answer,meta=cert.bound_query_solve(1,case["roles"],case["word"],
                                               case["start"],rank,cap=128)
            queries=meta["queries"];proofs=0
        else:
            for j in active(cert,case["word"]):
                request(j)
            answer=cert.public_exact_fallback(1,case["roles"],case["word"],
                                               case["start"])
            queries=charged();proofs=0
        if arm in ("trained-quotient","fixed-quotient") and queries!=charged():
            raise RuntimeError("T20_PUBLIC_REQUEST_ACCOUNTING_MISMATCH")
        if answer!=case["private_truth"]:
            raise RuntimeError("T20_HARD_FALSE_POSITIVE:"+label)
        total_proofs+=proofs
        out.append({"id":label,"correct":True,"extra_public_queries":queries})
    return {"schema":"wingless.t20.internal-s3-process.v1","arm":arm,
      "repetition":repetition,"trained_cases":40,"dev_cases":64,
      "rows":out,"compiled_signatures":len(compiler),"train_updates":updates,
      "certifier_invocations":total_proofs,
      "complete_cpu_ns":time.process_time_ns()-start_ns,
      "process_peak_rss_bytes":rss(),"independent_science":False,
      "scientific_acceptance":False,"optimizer_admission":False,
      "T19_DEV_replayed":False,"T16_originals_opened":0,
      "remote_query_service_cost":"UNMEASURED"}

def run():
    data=[]
    for arm in ARMS:
        for repetition in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                              "--child",arm,str(repetition)],text=True,capture_output=True)
            if p.returncode:
                raise RuntimeError("T20_CHILD_FAILURE "+arm+" "+str(repetition)+"\n"+p.stderr)
            data.append(json.loads(p.stdout))
    result={}
    for arm in ARMS:
        cases=[d for d in data if d["arm"]==arm]
        assert len(cases)==5 and all(c["rows"]==cases[0]["rows"] for c in cases)
        result[arm]={"correct":sum(int(z["correct"]) for z in cases[0]["rows"]),
          "queries":sum(z["extra_public_queries"] for z in cases[0]["rows"]),
          "median_cpu_ns":int(statistics.median(c["complete_cpu_ns"] for c in cases)),
          "max_rss_bytes":max(c["process_peak_rss_bytes"] for c in cases),
          "cpu_five_process_ns":[c["complete_cpu_ns"] for c in cases],
          "rss_five_process_bytes":[c["process_peak_rss_bytes"] for c in cases],
          "compiled_signatures":cases[0]["compiled_signatures"],
          "train_updates":cases[0]["train_updates"]}
    learned=result["trained-quotient"]
    comparison={k:{"correctness_parity":learned["correct"]==v["correct"],
      "query_nonregression":learned["queries"]<=v["queries"],
      "cpu_nonregression":learned["median_cpu_ns"]<=v["median_cpu_ns"],
      "rss_nonregression":learned["max_rss_bytes"]<=v["max_rss_bytes"]}
      for k,v in result.items() if k!="trained-quotient"}
    assert all(v["correct"]==64 for v in result.values())
    return {"schema":"wingless.t20.internal-only-structural-quotient.v1",
      "source":"FRESH SELF-AUTHORED S3 TRAIN40 DEV64",
      "results":result,"comparison":comparison,
      "strict_query_gain_vs_trained":learned["queries"]<result["trained-T16B"]["queries"],
      "strict_cpu_gain_vs_trained":learned["median_cpu_ns"]<result["trained-T16B"]["median_cpu_ns"],
      "complete_joint_gate_pass":False,
      "complete_joint_gate_reason":"REMOTE_PUBLIC_QUERY_SERVICE_COST_UNMEASURED",
      "original_source_first_use":False,"independent_science":False,
      "scientific_acceptance":False,"optimizer_admission":False,
      "T19_DEV_replayed":False,"T16_originals_opened":0,
      "no_new_controller":True,"no_workstation":True,"no_merge":True}

if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--property":
        _k,c,_s=deps()
        print(json.dumps({"T20_EXACT_QUOTIENT_PROPERTY_PASS":prove_sound_on_public_algebra(c),
                          "no_DEV_source_consumed":True}))
    elif len(sys.argv)==4 and sys.argv[1]=="--child":
        print(json.dumps(child(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(run(),sort_keys=True))
    else:
        raise SystemExit("T20_INVALID_ARGUMENTS")
