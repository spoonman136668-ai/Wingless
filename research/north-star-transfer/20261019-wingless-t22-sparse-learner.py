"""T22 frozen sparse exact group identity rewrite learner and five controls.

Source authoring module imports neither candidate nor frozen baseline.
The learned object is a bounded certified symbolic identity dictionary.
Never uses a 216-bit compositional motif relation. The 216 universal checks
are TRAIN-only proof steps; inference only rewrites with authenticated
identities then runs the unchanged sound T16B certifier.
"""
from __future__ import annotations
import importlib.util,itertools,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
HERE=Path(__file__).resolve().parent
ARMS=("learned-certified-sparse-rewrite",
      "trained-frequency-certified-rewrite",
      "fixed-rewrite-equal-cap",
      "equally-trained-T16B","full-public-exact")
MAX_CANDIDATE_TESTS=384
MAX_RULES=16
MAX_RULE_BYTES=4096
def load(name,filename):
    spec=importlib.util.spec_from_file_location(name,HERE/filename)
    module=importlib.util.module_from_spec(spec)
    sys.modules[name]=module
    spec.loader.exec_module(module)
    return module

def deps():
    return (load("t22_cert","20261019-wingless-t16b-persistent-bindings.py"),
            load("t22_source","20261019-wingless-t22-new-sources.py"))

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
            raise RuntimeError("T22_RSS_UNAVAILABLE")
        return int(x.PeakWorkingSetSize)
    import resource
    x=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(x)*1024 if sys.platform.startswith("linux") else int(x)

def universal_identity(cert,motif):
    """Global-role universal identity, independently checked 216 times."""
    for assignment in itertools.product(cert.PERMS,repeat=3):
        for initial in range(3):
            state=initial
            for role,inv in motif:
                state=cert.apply_op(1,assignment[role],state,inv)
            if state!=initial:
                return False
    return True

def train_candidates(cert,train):
    counts=Counter()
    for case in train:
        word=tuple(case["word"])
        for n in range(4,min(16,len(word))+1):
            for i in range(len(word)-n+1):
                motif=word[i:i+n]
                if cert.reduce_word(motif): # freely trivial motifs have no new representation.
                    counts[motif]+=1
    return counts

def fixed_candidates():
    # Mathematically motivated, untrained equal-cap standard dictionary.
    unary=tuple((role,inv) for role in range(3) for inv in (False,True))
    out=[((role,inv),)*6 for role,inv in unary]
    for i in range(3):
        for j in range(3):
            if i==j:continue
            comm=((i,False),(j,False),(i,True),(j,True))
            out.append(comm*3)
    return out

def fit(cert,train,arm):
    assert arm in ARMS[:3]
    hist=train_candidates(cert,train) # equal TRAIN available to each.
    if arm=="fixed-rewrite-equal-cap":
        ranked=fixed_candidates()
    elif arm=="learned-certified-sparse-rewrite":
        ranked=sorted(hist,key=lambda w:(
            -(hist[w]-1)*(len(w)-1)/(32+len(w)*2),-hist[w],-len(w),w))
    else:
        ranked=sorted(hist,key=lambda w:(-hist[w],-len(w),w))
    rules=[];updates=0;payload=0;checked=set()
    for word in ranked:
        if len(rules)>=MAX_RULES or updates>=MAX_CANDIDATE_TESTS:break
        if word in checked:continue
        checked.add(word)
        cost=32+2*len(word)
        if payload+cost>MAX_RULE_BYTES:continue
        updates+=1
        if universal_identity(cert,word):
            rules.append(word)
            payload+=cost
    return tuple(rules),updates,payload,len(hist)

def reduce_with_rules(cert,word,rules):
    """Delete certified exact universal identities; preserve other roles."""
    current=tuple(cert.reduce_word(word))
    cuts=0;passes=0
    while rules and passes<8:
        passes+=1
        ordered=sorted(rules,key=lambda w:(-len(w),w))
        out=[];i=0;new=0
        while i<len(current):
            match=next((w for w in ordered if current[i:i+len(w)]==w),None)
            if match:
                i+=len(match);cuts+=1;new+=1
            else:
                out.append(current[i]);i+=1
        updated=tuple(cert.reduce_word(out))
        if not new or updated==current:
            break
        current=updated
    return current,cuts

def exact_public_solve(cert,visible,request,word,start,rank):
    rows=[tuple(x) for x in visible]
    active=tuple(sorted({role for role,inv in cert.reduce_word(word)}))
    queries=0
    while True:
        value,proof=cert.shared_role_certificate(1,rows,word,start,budget=128)
        if proof["certified"]:
            return value,queries
        remaining=[j for j in active if len(rows[j])==1]
        if not remaining:break
        priority=max(remaining,key=lambda j:(rank[j],-j))
        rows[priority]=(rows[priority][0],request(priority))
        queries+=1
    return cert.public_exact_fallback(1,rows,word,start),queries

def trial(arm,rep):
    assert arm in ARMS and 0<=rep<5
    t0=time.process_time_ns()
    cert,source=deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert not set(source.TRAIN)&set(source.DEV)
    train=[source.public_train_case(source.build(k)) for k in source.TRAIN]
    rank=cert.score_role_requests([(1,c["roles"],c["word"],c["start"]) for c in train])
    rules=();checks=0;bytes_used=0;candidate_count=0
    if arm in ARMS[:3]:
        rules,checks,bytes_used,candidate_count=fit(cert,train,arm)
    cases=[];total_cuts=0;reduced_tokens=0
    for label in source.DEV:
        case=source.build(label)
        visible=source.public_dev_interface(case)
        original=visible["word"]
        if arm in ARMS[:3]:
            transformed,cuts=reduce_with_rules(cert,original,rules)
            reduced_tokens+=len(original)-len(transformed)
            total_cuts+=cuts
            value,q=exact_public_solve(cert,visible["first"],
                  visible["request"],transformed,visible["start"],rank)
        elif arm=="equally-trained-T16B":
            value,q=exact_public_solve(cert,visible["first"],visible["request"],
                original,visible["start"],rank)
        else:
            used=tuple(sorted({role for role,inv in cert.reduce_word(original)}))
            for j in used:
                visible["request"](j)
            value=cert.public_exact_fallback(1,case["roles"],original,visible["start"])
            q=visible["charged"]()
        if value!=case["private_truth"] or q!=visible["charged"]():
            raise RuntimeError("T22_FALSE_POSITIVE_OR_QUERY_ACCOUNTING:"+label)
        cases.append({"id":label,"correct":True,"queries":q})
    return {"schema":"wingless.t22.new_internal_per_process.v1","arm":arm,
       "repetition":rep,"rows":cases,
       "candidate_rule_checks":checks,"dictionary_rules":len(rules),
       "dictionary_bytes":bytes_used,"train_candidate_motifs":candidate_count,
       "DEV_exact_rewrite_hits":total_cuts,"DEV_reduced_tokens":reduced_tokens,
       "complete_cpu_ns":time.process_time_ns()-t0,
       "process_peak_rss_bytes":rss(),
       "source_kind":"NEW_INTERNAL_NOT_EXTERNALLY_AUTHORED",
       "old_dev_replayed":False,"scientific_acceptance":False,"optimizer_admission":False}

def evaluate():
    outcomes=[]
    for arm in ARMS:
        for rep in range(5):
            run=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                  "--arm",arm,str(rep)],capture_output=True,text=True)
            if run.returncode:
                raise RuntimeError("T22_CHILD_FAIL "+arm+" "+str(rep)+"\n"+run.stderr)
            outcomes.append(json.loads(run.stdout))
    ans={}
    for arm in ARMS:
        r=[x for x in outcomes if x["arm"]==arm]
        assert len(r)==5 and all(x["rows"]==r[0]["rows"] for x in r)
        ans[arm]={"correct":sum(x["correct"] for x in r[0]["rows"]),
         "queries":sum(x["queries"] for x in r[0]["rows"]),
         "rows":r[0]["rows"],"rule_checks":r[0]["candidate_rule_checks"],
         "rule_count":r[0]["dictionary_rules"],"dictionary_bytes":r[0]["dictionary_bytes"],
         "train_candidates":r[0]["train_candidate_motifs"],
         "DEV_rewrite_hits":r[0]["DEV_exact_rewrite_hits"],
         "DEV_reduced_tokens":r[0]["DEV_reduced_tokens"],
         "cpu_five_process_ns":[x["complete_cpu_ns"] for x in r],
         "rss_five_process_bytes":[x["process_peak_rss_bytes"] for x in r],
         "median_complete_cpu_ns":int(statistics.median(x["complete_cpu_ns"] for x in r)),
         "maximum_peak_rss_bytes":max(x["process_peak_rss_bytes"] for x in r)}
    assert all(x["correct"]==64 for x in ans.values())
    learned=ans[ARMS[0]]
    cmp={arm:{
      "correctness_nonregression":learned["correct"]>=ref["correct"],
      "public_query_nonregression":learned["queries"]<=ref["queries"],
      "CPU_nonregression":learned["median_complete_cpu_ns"]<=ref["median_complete_cpu_ns"],
      "RSS_nonregression":learned["maximum_peak_rss_bytes"]<=ref["maximum_peak_rss_bytes"],
    } for arm,ref in ans.items() if arm!=ARMS[0]}
    return {"schema":"wingless.t22.prospective-internal-sparse-identity.v1",
      "new_source_origin":"INTERNAL_SOURCE_BLIND_TWO_FAMILIES_G_H",
      "results":ans,"comparisons":cmp,
      "strict_query_advantage_vs_trained_T16B":learned["queries"]<ans["equally-trained-T16B"]["queries"],
      "strict_cpu_advantage_vs_trained_T16B":learned["median_complete_cpu_ns"]<ans["equally-trained-T16B"]["median_complete_cpu_ns"],
      "joint_scientific_acceptance":False,
      "remote_public_service_cost":"UNMEASURED",
      "independent_external_replication":False,"optimizer_admission":False,
      "no_old_DEV_replay":True,"no_controller_changes":True,"no_merge":True}

if __name__=="__main__":
    if len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(evaluate(),sort_keys=True))
    else:
        raise SystemExit("T22_INVALID_ARGS")
