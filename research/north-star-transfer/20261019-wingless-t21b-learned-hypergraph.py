"""T21-B exact-certifying learned counterexample-hypergraph developmental trial.

TRAIN-only graph-message learning; no hidden TRAIN truth, no DEV training.
Candidate and baselines receive only public calibration and paid second pairs.
Scientific first-use and complete remote query cost are NOT qualified.
"""
from __future__ import annotations
import ctypes
import importlib.util
from fractions import Fraction
from functools import lru_cache
from itertools import permutations, product
import json
import statistics
import subprocess
import sys
import time
from pathlib import Path

HERE=Path(__file__).parent
WORLDS=tuple(product((0,1),repeat=3))
OBSERVATIONS=tuple(product((-1,0,1),repeat=3))
ORDERS=tuple(permutations(range(3)))
ARMS=("trained-hypergraph","trained-T16B","trained-best-static",
      "fixed-exact-adaptive","full-public-exact")
WEIGHTS=8
TRAIN_EPOCHS=18

def load(name,file):
    spec=importlib.util.spec_from_file_location(name,HERE/file)
    mod=importlib.util.module_from_spec(spec)
    sys.modules[name]=mod
    spec.loader.exec_module(mod)
    return mod

def cert_module():
    return load("t21b_cert","20261019-wingless-t16b-persistent-bindings.py")

def dependencies():
    return (load("t21b_kernel","20261018-wingless-t3-typed-composition.py"),
            cert_module(),load("t21b_source","20261019-wingless-t21b-fresh-source.py"))

def peak_rss():
    if sys.platform=="win32":
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
            raise RuntimeError("T21B_PEAK_RSS_UNAVAILABLE")
        return int(m.PeakWorkingSetSize)
    import resource
    rss=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(rss)*1024 if sys.platform.startswith("linux") else int(rss)

def public_worlds(cert,word,public,start):
    legal=tuple(cert.consistent_ops(1,public[r]) for r in range(3))
    if any(len(ops)!=2 for ops in legal):
        raise RuntimeError("T21B_INVALID_FIRST_PAIR")
    sequence=cert.reduce_word(word)
    outputs=[]
    for bits in WORLDS:
        current=start
        for role,inverse in sequence:
            current=cert.apply_op(1,legal[role][bits[role]],current,inverse)
        outputs.append(current)
    return tuple(outputs),legal

def remaining(obs):
    return tuple(i for i,b in enumerate(WORLDS)
      if all(x<0 or b[r]==x for r,x in enumerate(obs)))

def exact_answer(outputs,obs):
    values={outputs[i] for i in remaining(obs)}
    if len(values)==1:
        return next(iter(values))
    return None

def conflict_edges(outputs,indexes):
    """Distinct-output world-pair disagreements form typed hyperedges."""
    return tuple(tuple(j for j in range(3) if WORLDS[a][j]!=WORLDS[b][j])
      for pos,a in enumerate(indexes) for b in indexes[pos+1:]
      if outputs[a]!=outputs[b])

def graph_features(outputs,obs,role):
    if obs[role]>=0:
        raise RuntimeError("T21B_ALREADY_PAID_ROLE")
    indexes=remaining(obs)
    edges=conflict_edges(outputs,indexes)
    if not edges:
        return (0.0,)*WEIGHTS
    scale=float(len(edges))
    incidence=[edge for edge in edges if role in edge]
    degree=[sum(len(e)==length for e in incidence)/scale for length in (1,2,3)]
    before=len({outputs[i] for i in indexes})
    cert_gain=0.0
    edge_removal=0.0
    output_gain=0.0
    for bit in (0,1):
        following=tuple(bit if j==role else v for j,v in enumerate(obs))
        sub=remaining(following)
        probability=len(sub)/len(indexes)
        cert_gain+=probability*int(exact_answer(outputs,following) is not None)
        edge_removal+=probability*(1.0-len(conflict_edges(outputs,sub))/scale)
        output_gain+=probability*(before-len({outputs[i] for i in sub}))/3.0
    return (1.0,len(incidence)/scale,*degree,edge_removal,cert_gain,output_gain)

def exact_adaptive_actions(outputs,active):
    @lru_cache(None)
    def value(obs):
        if exact_answer(outputs,obs) is not None:
            return Fraction(0)
        present=remaining(obs)
        costs=[]
        for role in active:
            if obs[role]>=0: continue
            cost=Fraction(1)
            for bit in (0,1):
                nxt=tuple(bit if r==role else x for r,x in enumerate(obs))
                n=sum(WORLDS[i][role]==bit for i in present)
                if not n: raise AssertionError("T21B_ZERO_MASS")
                cost+=Fraction(n,len(present))*value(nxt)
            costs.append((cost,role))
        if not costs: raise RuntimeError("T21B_NO_EXACT_ADAPTIVE_ACTION")
        return min(costs)[0]
    def options(obs):
        if exact_answer(outputs,obs) is not None: return ()
        present=remaining(obs)
        scores=[]
        for role in active:
            if obs[role]>=0: continue
            cost=Fraction(1)
            for bit in (0,1):
                nxt=tuple(bit if r==role else x for r,x in enumerate(obs))
                n=sum(WORLDS[i][role]==bit for i in present)
                cost+=Fraction(n,len(present))*value(nxt)
            scores.append((cost,role))
        best=min(cost for cost,_ in scores)
        return tuple(role for cost,role in scores if cost==best)
    return value,options

def static_cost(outputs,active,ordering):
    total=0
    for bits in WORLDS:
        obs=[-1,-1,-1]
        while exact_answer(outputs,tuple(obs)) is None:
            role=next(r for r in ordering if r in active and obs[r]<0)
            obs[role]=bits[role]
            total+=1
    return Fraction(total,len(WORLDS))

def fit_static(cert,train):
    totals={order:Fraction(0) for order in ORDERS}
    for case in train:
        public=tuple((row[0],) for row in case["roles"])
        output,_=public_worlds(cert,case["word"],public,case["start"])
        active=tuple(sorted({r for r,_ in cert.reduce_word(case["word"])}))
        for order in ORDERS:
            totals[order]+=static_cost(output,active,order)
    chosen=min(ORDERS,key=lambda order:(totals[order],order))
    return chosen,len(train)*len(ORDERS)

def fit_hypergraph(cert,train):
    examples=[]
    for case in train:
        public=tuple((row[0],) for row in case["roles"])
        outputs,_=public_worlds(cert,case["word"],public,case["start"])
        active=tuple(sorted({r for r,_ in cert.reduce_word(case["word"])}))
        _,teacher=exact_adaptive_actions(outputs,active)
        for obs in OBSERVATIONS:
            if exact_answer(outputs,obs) is not None: continue
            legal=tuple(r for r in active if obs[r]<0)
            if len(legal)<2: continue
            features={r:graph_features(outputs,obs,r) for r in legal}
            optimal=teacher(obs)
            if not optimal or not set(optimal)<=set(legal):
                raise RuntimeError("T21B_INVALID_PUBLIC_TEACHER")
            examples.append((legal,features,optimal))
    weight=[0.0]*WEIGHTS;updates=0
    for epoch in range(TRAIN_EPOCHS):
        rate=0.16/(1+epoch/6.0)
        for legal,feat,optimal in examples:
            predicted=max(legal,key=lambda r:(
                sum(weight[k]*feat[r][k] for k in range(WEIGHTS)),-r))
            if predicted in optimal: continue
            target=min(optimal)
            for k in range(WEIGHTS):
                weight[k]+=rate*(feat[target][k]-feat[predicted][k])
            updates+=1
    return tuple(weight),{"training_public_states":len(examples),
                          "training_updates":updates,"epochs":TRAIN_EPOCHS,
                          "message_weights":WEIGHTS}

def solve_graph(cert,word,start,public,request,ranker):
    outputs,operators=public_worlds(cert,word,public,start)
    active=tuple(sorted({r for r,_ in cert.reduce_word(word)}))
    obs=[-1,-1,-1];bought=0
    while True:
        answer=exact_answer(outputs,tuple(obs))
        if answer is not None:
            return answer,bought
        eligible=tuple(r for r in active if obs[r]<0)
        if not eligible: raise RuntimeError("T21B_UNCERTIFIED_COMPLETE")
        role=ranker(outputs,tuple(obs),eligible)
        if role not in eligible: raise RuntimeError("T21B_UNAUTHORIZED_ROLE")
        second=request(role)
        consistent=cert.consistent_ops(1,(public[role][0],second))
        if len(consistent)!=1:
            raise RuntimeError("T21B_SECOND_PAIR_NOT_DISCRIMINATING")
        obs[role]=operators[role].index(consistent[0])
        bought+=1

def solve_T16B(cert,word,start,public,request,rank):
    rows=[tuple(p) for p in public]
    used=tuple(sorted({r for r,_ in cert.reduce_word(word)}))
    bought=0
    while True:
        value,proof=cert.shared_role_certificate(1,rows,word,start,budget=128)
        if proof["certified"]: return value,bought
        pending=[r for r in used if len(rows[r])==1]
        if not pending: raise RuntimeError("T21B_T16B_UNCERTIFIED")
        next_role=max(pending,key=lambda r:(rank[r],-r))
        rows[next_role]=(rows[next_role][0],request(next_role))
        bought+=1

def property_checks():
    """Public mathematical fixtures only; does NOT even import T21-B source."""
    cert=cert_module()
    count=0;renames=0;oracle_checks=0
    for order in ORDERS:
      for flags in product((False,True),repeat=3):
        word=tuple((order[j],flags[j]) for j in range(3))
        for images in ((0,0,0),(0,1,2),(1,0,2),(2,1,0)):
          first=tuple(((role,images[role]),) for role in range(3))
          for start in range(3):
            outputs,operators=public_worlds(cert,word,first,start)
            active=tuple(sorted({r for r,_ in cert.reduce_word(word)}))
            dp,options=exact_adaptive_actions(outputs,active)
            if dp((-1,-1,-1))>min(static_cost(outputs,active,o) for o in ORDERS):
                raise RuntimeError("T21B_ADAPTIVE_EXCEEDS_STATIC_BOUND")
            oracle_checks+=1
            for obs in OBSERVATIONS:
                candidate=exact_answer(outputs,obs)
                pairs=[tuple(row) for row in first]
                for role,b in enumerate(obs):
                    if b>=0:
                        op=operators[role][b]
                        other=(first[role][0][0]+1)%3
                        pairs[role]=(pairs[role][0],(other,op[other]))
                exact,proof=cert.shared_role_certificate(1,pairs,word,start,budget=4096)
                if (candidate is not None)!=proof["certified"] or candidate!=exact:
                    raise RuntimeError("T21B_PROPERTY_CERTIFICATION_DISAGREEMENT")
                count+=1
            if oracle_checks<=18:
                for rename in ORDERS:
                    shifted_first=[None]*3
                    for role in range(3):
                        shifted_first[rename[role]]=first[role]
                    rewritten=tuple((rename[r],inv) for r,inv in word)
                    new_outputs,_=public_worlds(cert,rewritten,tuple(shifted_first),start)
                    old_index={b:i for i,b in enumerate(WORLDS)}
                    for i,newbits in enumerate(WORLDS):
                        oldbits=tuple(newbits[rename[role]] for role in range(3))
                        if new_outputs[i]!=outputs[old_index[oldbits]]:
                            raise RuntimeError("T21B_ROLE_RENAME_OUTPUT")
                    obs=(-1,-1,-1)
                    if exact_answer(outputs,obs) is None:
                        for role in range(3):
                            f=graph_features(outputs,obs,role)
                            g=graph_features(new_outputs,obs,rename[role])
                            if any(abs(x-y)>1e-10 for x,y in zip(f,g)):
                                raise RuntimeError("T21B_NON_EQUIVARIANT_MESSAGE")
                    renames+=1
    if count!=48*4*3*27 or renames!=18*6:
        raise RuntimeError("T21B_INCOMPLETE_PUBLIC_PROPERTY")
    return {"exact_source_free_equivalences":count,"public_oracle_bounds":oracle_checks,
            "role_rename_equivariance":renames,"no_train_or_dev_import":True}

def child(arm,rep):
    assert arm in ARMS and 0<=rep<5
    began=time.process_time_ns()
    kernel,cert,source=dependencies()
    assert len(source.TRAIN)==72 and len(source.DEV)==96
    assert not set(source.TRAIN)&set(source.DEV)
    train=[source.build(label,kernel) for label in source.TRAIN]
    rank=cert.score_role_requests([(c["domain"],c["roles"],c["word"],c["start"]) for c in train])
    order,static_train=fit_static(cert,train) if arm=="trained-best-static" else ((0,1,2),0)
    weights,learning=fit_hypergraph(cert,train) if arm=="trained-hypergraph" else ((0.0,)*WEIGHTS,{})
    results=[]
    for label in source.DEV:
        case=source.build(label,kernel)
        public,request,charged=source.observed_interface(case)
        word=case["word"];start=case["start"]
        if arm=="trained-hypergraph":
            def choose(outputs,obs,legal):
                return max(legal,key=lambda r:(
                    sum(weights[k]*graph_features(outputs,obs,r)[k] for k in range(WEIGHTS)),-r))
            answer,paid=solve_graph(cert,word,start,public,request,choose)
        elif arm=="trained-T16B":
            answer,paid=solve_T16B(cert,word,start,public,request,rank)
        elif arm=="trained-best-static":
            def choose(outputs,obs,legal):
                return next(r for r in order if r in legal)
            answer,paid=solve_graph(cert,word,start,public,request,choose)
        elif arm=="fixed-exact-adaptive":
            outputs,_=public_worlds(cert,word,public,start)
            active=tuple(sorted({r for r,_ in cert.reduce_word(word)}))
            _cost,optimal=exact_adaptive_actions(outputs,active)
            def choose(outputs,obs,legal):
                return min(optimal(obs))
            answer,paid=solve_graph(cert,word,start,public,request,choose)
        else:
            active=tuple(sorted({r for r,_ in cert.reduce_word(word)}))
            rows=[tuple(row) for row in public]
            for role in active:
                rows[role]=(rows[role][0],request(role))
            answer=cert.public_exact_fallback(1,rows,word,start)
            paid=len(active)
        if paid!=charged():
            raise RuntimeError("T21B_MISMATCHED_PUBLIC_QUERY_ACCOUNTING:"+label)
        if answer!=case["private_truth"]:
            raise RuntimeError("T21B_HARD_FALSE_POSITIVE:"+label)
        results.append((label,paid))
    return {"arm":arm,"rep":rep,"correct":len(results),
        "queries":sum(v for _,v in results),"per_case_queries":results,
        "whole_process_cpu_ns":time.process_time_ns()-began,
        "whole_process_peak_rss_bytes":peak_rss(),
        "training_public_states":learning.get("training_public_states",0),
        "training_updates":learning.get("training_updates",0),
        "trained_coefficients":learning.get("message_weights",0),
        "static_train_objective_evaluations":static_train,
        "T20_DEV_replayed":False,"T16_originals_opened":0,
        "remote_query_service_cost":"UNMEASURED","independent_science":False,
        "scientific_acceptance":False,"optimizer_admission":False}

def aggregate():
    items=[]
    for arm in ARMS:
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
               "--child",arm,str(rep)],text=True,capture_output=True)
            if p.returncode:
                raise RuntimeError("T21B_FAILED_PROCESS "+arm+" "+str(rep)+"\n"+p.stderr)
            items.append(json.loads(p.stdout))
    result={}
    for arm in ARMS:
        cases=[v for v in items if v["arm"]==arm]
        assert len(cases)==5 and all(x["per_case_queries"]==cases[0]["per_case_queries"] for x in cases)
        result[arm]={
          "correct":cases[0]["correct"],"queries":cases[0]["queries"],
          "cpu_five_process_ns":[x["whole_process_cpu_ns"] for x in cases],
          "rss_five_process_bytes":[x["whole_process_peak_rss_bytes"] for x in cases],
          "median_cpu_ns":int(statistics.median(x["whole_process_cpu_ns"] for x in cases)),
          "max_peak_rss_bytes":max(x["whole_process_peak_rss_bytes"] for x in cases),
          "training_public_states":cases[0]["training_public_states"],
          "training_updates":cases[0]["training_updates"],
          "trained_coefficients":cases[0]["trained_coefficients"],
          "static_train_objective_evaluations":cases[0]["static_train_objective_evaluations"]}
    assert len(result)==5 and all(v["correct"]==96 for v in result.values())
    learned=result["trained-hypergraph"]
    comparison={arm:{
        "queries_nonregression":learned["queries"]<=values["queries"],
        "cpu_nonregression":learned["median_cpu_ns"]<=values["median_cpu_ns"],
        "rss_nonregression":learned["max_peak_rss_bytes"]<=values["max_peak_rss_bytes"]}
        for arm,values in result.items() if arm!="trained-hypergraph"}
    return {"schema":"wingless.t21b.internal-learned-hypergraph.v1",
      "origin":"FRESH_SELF_AUTHORED_INTERNAL_TRAIN72_DEV96",
      "results":result,"comparison":comparison,
      "learned_strict_query_gain_vs_best_trained_static":
        learned["queries"]<result["trained-best-static"]["queries"],
      "learned_strict_query_gain_vs_trained_T16B":
        learned["queries"]<result["trained-T16B"]["queries"],
      "joint_acceptance":False,
      "joint_acceptance_reason":"REMOTE_QUERY_SERVICE_COST_UNMEASURED_AND_NO_INDEPENDENT_SOURCE",
      "T20_DEV_replayed":False,"T16_originals_opened":0,
      "scientific_acceptance":False,"independent_science":False,
      "optimizer_admission":False,"merge_authorized":False,
      "no_workstation":True,"no_new_controller":True,"no_CKB_plane_link":True}

if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--property":
        print(json.dumps(property_checks(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--child":
        print(json.dumps(child(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(aggregate(),sort_keys=True))
    else:
        raise SystemExit("T21B_INVALID_ARGUMENTS")
