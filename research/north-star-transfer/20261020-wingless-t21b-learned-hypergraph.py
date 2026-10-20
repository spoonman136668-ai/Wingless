"""T21-B learned conflict-hypergraph decision policy; new INTERNAL DEV only.

Safety is enforced by the unchanged exact shared-binding certificate.
The learned component only selects which second public role row to buy.
No protected originals, no previous DEV, no independent science or promotion.
"""
from __future__ import annotations
import importlib.util
import json
import statistics
import subprocess
import sys
import time
from functools import lru_cache
from fractions import Fraction
from itertools import permutations,product
from pathlib import Path

HERE=Path(__file__).resolve().parent
WORLDS=tuple(product((0,1),repeat=3))
ORDERS=tuple(permutations(range(3)))
OBSERVATIONS=tuple(product((-1,0,1),repeat=3))
ARMS=("learned-hypergraph","optimized-trained-static","trained-T16B",
      "fixed-exact-adaptive","untrained-hypergraph","full-public-exact")
EPOCHS=3
ETA=0.25

def module(name,file):
    spec=importlib.util.spec_from_file_location(name,HERE/file)
    m=importlib.util.module_from_spec(spec)
    sys.modules[name]=m
    spec.loader.exec_module(m)
    return m

def dependencies():
    return (module("t21b_cert","20261019-wingless-t16b-persistent-bindings.py"),
            module("t21b_src","20261020-wingless-t21b-fresh-source.py"))

def peak_rss_bytes():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Mem(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
              ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
              ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
              ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
              ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        k=ctypes.WinDLL("kernel32");p=ctypes.WinDLL("psapi")
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T21B_RSS_UNAVAILABLE")
        return int(m.PeakWorkingSetSize)
    import resource
    val=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(val)*1024 if sys.platform.startswith("linux") else int(val)

def public_worlds(cert,visible,word,start):
    """Only first input/output calibration pairs, not hidden operators."""
    ops=tuple(cert.consistent_ops(1,visible[j]) for j in range(3))
    if any(len(x)!=2 for x in ops):
        raise AssertionError("T21B_S3_TWO_OP_PUBLIC_INVARIANT")
    outputs=[]
    for bits in WORLDS:
        state=start
        for role,inverse in word:
            state=cert.apply_op(1,ops[role][bits[role]],state,inverse)
        outputs.append(state)
    return tuple(outputs),ops

def eligible(obs):
    return tuple(i for i,w in enumerate(WORLDS)
                 if all(obs[j]==-1 or obs[j]==w[j] for j in range(3)))

def deterministic(outputs,obs):
    indices=eligible(obs)
    assert indices
    values={outputs[i] for i in indices}
    return len(values)==1

@lru_cache(maxsize=20000)
def best_exact_action(outputs,obs):
    """Uniform public-consistent-world exact DP; no observed DEV truth."""
    if deterministic(outputs,obs):
        return (Fraction(0),-1)
    poss=eligible(obs)
    actions=[]
    for r in range(3):
        if obs[r]!=-1:
            continue
        q=Fraction(1)
        for bit in (0,1):
            branch=tuple(bit if j==r else obs[j] for j in range(3))
            count=sum(WORLDS[i][r]==bit for i in poss)
            q+=Fraction(count,len(poss))*best_exact_action(outputs,branch)[0]
        actions.append((q,r))
    return min(actions)

def conflict_graph_features(outputs,obs,word,role):
    """Structural public hyperedges between output-disagreeing legal worlds.
    No explicit role ID enters features, so scored geometry is equivariant.
    """
    valid=eligible(obs)
    if obs[role]!=-1 or deterministic(outputs,obs):
        raise AssertionError("T21B_FEATURE_ON_CERTIFIED_OR_OBSERVED")
    conflicts=[]
    for i,k in enumerate(valid):
        for z in valid[i+1:]:
            if outputs[k]!=outputs[z]:
                conflicts.append(tuple(j for j in range(3)
                      if WORLDS[k][j]!=WORLDS[z][j] and obs[j]==-1))
    assert conflicts
    m=len(conflicts)
    covering=sum(role in edge for edge in conflicts)/m
    single=sum(edge==(role,) for edge in conflicts)/m
    double=sum(role in edge and len(edge)==2 for edge in conflicts)/m
    triple=sum(role in edge and len(edge)==3 for edge in conflicts)/m
    splits=[]
    for bit in (0,1):
        ids=[k for k in valid if WORLDS[k][role]==bit]
        assert ids
        splits.append(len({outputs[k] for k in ids}))
    # Exact certificate remains separate from this fallible learned score.
    guaranteed_certificates=sum(x==1 for x in splits)/2
    diversity_reduction=(len({outputs[k] for k in valid})-sum(splits)/2)/3
    appearances=sum(r==role for r,_ in word)/len(word)
    inverse_frac=sum(r==role and inv for r,inv in word)/len(word)
    return (covering,single,double,triple,guaranteed_certificates,
            diversity_reduction,appearances,inverse_frac)

def choose_learned(outputs,obs,word,weights):
    ranks=[]
    for r in range(3):
        if obs[r]>=0:
            continue
        features=conflict_graph_features(outputs,obs,word,r)
        ranks.append((sum(w*x for w,x in zip(weights,features)),-r,r))
    return max(ranks)[2]

def fit_hypergraph(cert,train):
    """Public-only imitation of exact adaptive teacher; no hidden truth,
    no DEV labels, fixed three epochs and frozen hyperparameters.
    """
    weights=[0.0]*8
    public=[]
    for case in train:
        first=tuple((r[0],) for r in case["full_rows"])
        outputs,_ops=public_worlds(cert,first,case["word"],case["start"])
        public.append((outputs,case["word"]))
    updates=0;teacher_states=0
    for epoch in range(EPOCHS):
        for outputs,word in public:
            for obs in OBSERVATIONS:
                if deterministic(outputs,obs):
                    continue
                target=best_exact_action(outputs,obs)[1]
                pred=choose_learned(outputs,obs,word,weights)
                teacher_states+=1
                if target!=pred:
                    good=conflict_graph_features(outputs,obs,word,target)
                    bad=conflict_graph_features(outputs,obs,word,pred)
                    for j in range(8):
                        weights[j]+=ETA*(good[j]-bad[j])
                    updates+=1
    return tuple(weights),updates,teacher_states

def trained_static_order(cert,train):
    """Optimize ALL six static orders using exactly the same TRAIN public
    ambiguity as the adaptive learner, not a weaker fixed baseline.
    """
    costs={order:0 for order in ORDERS}
    for case in train:
        first=tuple((r[0],) for r in case["full_rows"])
        outputs,_=public_worlds(cert,first,case["word"],case["start"])
        for ordering in ORDERS:
            for world in WORLDS:
                obs=[-1]*3;n=0
                while not deterministic(outputs,tuple(obs)):
                    role=next(r for r in ordering if obs[r]<0)
                    obs[role]=world[role];n+=1
                costs[ordering]+=n
    return min(ORDERS,key=lambda r:(costs[r],r)),costs

def solve(cert,word,start,initial,request,arm,weights,order,rank):
    visible=[tuple(p) for p in initial]
    outputs,options=public_worlds(cert,visible,word,start)
    obs=[-1]*3
    purchases=0;checks=0
    while True:
        value,proof=cert.shared_role_certificate(1,visible,word,start,budget=64)
        checks+=1
        same=deterministic(outputs,tuple(obs))
        if proof["certified"]!=same:
            raise AssertionError("T21B_LEARNED_AND_INDEPENDENT_CERT_MISMATCH")
        if same:
            assert value is not None
            return value,purchases,checks
        if arm=="full-public-exact":
            role=next(r for r in range(3) if obs[r]<0)
        elif arm=="fixed-exact-adaptive":
            role=best_exact_action(outputs,tuple(obs))[1]
        elif arm=="learned-hypergraph":
            role=choose_learned(outputs,tuple(obs),word,weights)
        elif arm=="untrained-hypergraph":
            role=choose_learned(outputs,tuple(obs),word,(0.0,)*8)
        elif arm=="optimized-trained-static":
            role=next(r for r in order if obs[r]<0)
        elif arm=="trained-T16B":
            role=max((r for r in range(3) if obs[r]<0),key=lambda r:(rank[r],-r))
        else:
            raise ValueError("T21B_UNREGISTERED_ARM")
        second=request(role)
        matches=[bit for bit,perm in enumerate(options[role])
                 if perm[second[0]]==second[1]]
        if len(matches)!=1 or obs[role]>=0:
            raise AssertionError("T21B_PURCHASE_NOT_UNIQUELY_IDENTIFYING")
        obs[role]=matches[0]
        visible[role]=(visible[role][0],second)
        purchases+=1
        if purchases>3:
            raise AssertionError("T21B_QUERY_OVERFLOW")

def exact_public_property(cert):
    """Source-free adversarial algebra only; before first DEV read."""
    words=tuple(tuple((perm[j],j%2==0) for j in range(3))
                for perm in ORDERS)
    sample=0
    for w in words:
        for vals in product(range(3),repeat=3):
            initial=tuple(((j,vals[j]),) for j in range(3))
            for start in range(3):
                outputs,opts=public_worlds(cert,initial,w,start)
                for obs in OBSERVATIONS:
                    rows=[tuple(p) for p in initial]
                    for role,bit in enumerate(obs):
                        if bit>=0:
                            x=(role+1)%3
                            rows[role]=(rows[role][0],(x,opts[role][bit][x]))
                    answer,proof=cert.shared_role_certificate(1,rows,w,start,budget=64)
                    assert proof["certified"]==deterministic(outputs,obs)
                    if proof["certified"]:
                        assert answer==outputs[eligible(obs)[0]]
                    sample+=1
    if sample!=13122: # 6 permutations x 27 first pairs x 3 starts x 27 observation states
        raise AssertionError("T21B_PUBLIC_PROPERTY_INCOMPLETE:"+str(sample))
    return sample

def child(arm,rep):
    assert arm in ARMS and 0<=rep<5
    began=time.process_time_ns()
    cert,source=dependencies()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert not set(source.TRAIN)&set(source.DEV)
    train=[source.build(label,cert) for label in source.TRAIN]
    rank=cert.score_role_requests([(1,c["full_rows"],c["word"],c["start"]) for c in train])
    order,static_cost=trained_static_order(cert,train)
    weights=(0.0,)*8;updates=0;teacher_states=0
    if arm=="learned-hypergraph":
        weights,updates,teacher_states=fit_hypergraph(cert,train)
    queries=0;checks=0;rows=[]
    for label in source.DEV:
        case=source.build(label,cert)
        public,request,charged=source.public_interface(case)
        answer,n,p=solve(cert,case["word"],case["start"],public,request,
                         arm,weights,order,rank)
        if n!=charged():
            raise AssertionError("T21B_UNCHARGED_QUERY")
        if answer!=case["private_truth"]:
            raise RuntimeError("T21B_HARD_FALSE_POSITIVE:"+label)
        rows.append({"id":label,"correct":True,"purchases":n})
        queries+=n;checks+=p
    return {"schema":"wingless.t21b.whole-process.v1","arm":arm,
      "rep":rep,"correct":len(rows),"rows":rows,"public_queries":queries,
      "train":len(train),"dev":len(rows),"train_updates":updates,
      "teacher_states":teacher_states,"train_optimized_static_order":order,
      "train_static_expected_world_costs":static_cost if rep==0 else None,
      "trained_hypergraph_coefficients":weights if rep==0 and arm=="learned-hypergraph" else None,
      "verified_certificates":checks,
      "process_cpu_ns":time.process_time_ns()-began,"peak_rss_bytes":peak_rss_bytes(),
      "independent_science":False,"optimizer_admission":False,
      "T20_DEV_replayed":False,"T16_originals_opened":0,
      "remote_query_service_cost":"UNMEASURED"}

def run():
    results={}
    for arm in ARMS:
        samples=[]
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                "--child",arm,str(rep)],capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T21B_PROCESS_FAILED:"+arm+":"+str(rep)+"\n"+p.stderr)
            samples.append(json.loads(p.stdout))
        assert len(samples)==5 and all(x["rows"]==samples[0]["rows"] for x in samples)
        results[arm]={
          "correct":samples[0]["correct"],
          "queries":samples[0]["public_queries"],
          "cpu_samples_ns":[x["process_cpu_ns"] for x in samples],
          "rss_samples_bytes":[x["peak_rss_bytes"] for x in samples],
          "median_cpu_ns":int(statistics.median(x["process_cpu_ns"] for x in samples)),
          "max_rss_bytes":max(x["peak_rss_bytes"] for x in samples),
          "train_updates":samples[0]["train_updates"],
          "teacher_states":samples[0]["teacher_states"],
          "certificate_invocations":samples[0]["verified_certificates"],
          "trained_static_order":samples[0]["train_optimized_static_order"],
          "learned_coefficients":samples[0]["trained_hypergraph_coefficients"]}
    learned=results["learned-hypergraph"]
    comparisons={name:{"query_nonregression":learned["queries"]<=r["queries"],
      "cpu_nonregression":learned["median_cpu_ns"]<=r["median_cpu_ns"],
      "rss_nonregression":learned["max_rss_bytes"]<=r["max_rss_bytes"],
      "strict_query_gain":learned["queries"]<r["queries"]}
      for name,r in results.items() if name!="learned-hypergraph"}
    assert all(r["correct"]==64 and len(r["cpu_samples_ns"])==5
               and len(r["rss_samples_bytes"])==5 for r in results.values())
    return {"schema":"wingless.t21b.internal-only-learned-hypergraph.v1",
      "status":"DEVELOPMENT_ONLY_NOT_INDEPENDENT",
      "source":"fresh self-authored T21B TRAIN40 DEV64",
      "results":results,"comparisons":comparisons,
      "remote_query_service_cost":"UNMEASURED",
      "joint_science_and_optimizer_acceptance":False,
      "scientific_acceptance":False,"optimizer_admission":False,
      "T20_DEV_replayed":False,"T16_originals_opened":0,
      "repetitions_per_arm":5,"no_merge":True,"no_local_workstation":True}

if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--property":
        c,_=dependencies()
        print(json.dumps({"T21B_PUBLIC_PROPERTY_PASS":exact_public_property(c),
           "DEV_consumed":False,"T20_DEV_replayed":False}))
    elif len(sys.argv)==4 and sys.argv[1]=="--child":
        print(json.dumps(child(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(run(),sort_keys=True))
    else:
        raise SystemExit("T21B_BAD_ARGUMENTS")
