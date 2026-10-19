"""T31: TRAIN-learned conditional causal law proposal, exact posterior.

The learned conditional per-parent-bit law priors ACTUALLY order unseen
graph/mode/law world proposals in each counterfactual proof. They never
exclude a world or certify an answer. As soon as two outputs are witnessed,
UNKNOWN is sound; an exact certification requires scanning every world.
All costs include source, model fit, inference, and native memory.
"""
from __future__ import annotations
import importlib.util,json,math,statistics,subprocess,sys,time
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("learned-conditional-causal-program",
      "train-fitted-exact-constraint-policy",
      "train-frequency-conditional-ablation",
      "independent-explicit-world-enumerator")
ORDERS=("natural","reverse","popcount","reversed-popcount")
def source_module():
    spec=importlib.util.spec_from_file_location("t31_source",ROOT/"20261019-wingless-t31-new-causal-program-worlds.py")
    obj=importlib.util.module_from_spec(spec);sys.modules[spec.name]=obj
    spec.loader.exec_module(obj);return obj
def peak_rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Mem(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
             ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
             ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
             ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
             ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        k=ctypes.WinDLL("kernel32",use_last_error=True)
        p=ctypes.WinDLL("psapi",use_last_error=True)
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T31_WINDOWS_NATIVE_RSS_ERROR="+str(ctypes.get_last_error()))
        return int(m.PeakWorkingSetSize)
    import resource
    val=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(val)*1024 if sys.platform.startswith("linux") else int(val)
def public_posterior(public,source):
    """Reusable exact public constraints for graph+mode+unknown-rule worlds."""
    f=public["family"];out=[]
    for gid,graph in enumerate(source.TEMPLATES[f]):
        for mode in (0,1):
            known=ones=0;valid=True
            for p in public["probes"]:
                for node,output in p["observed"]:
                    if p["do"] is not None and p["do"][0]==node:
                        if output!=p["do"][1]:raise RuntimeError("T31_BAD_INTERVENTION_SOURCE")
                        continue
                    pat=source.pattern(p["pre"],graph[node])
                    bit=output^((p["pre"]>>node)&1)
                    if node%2:bit^=mode*((p["pre"]>>((node+1)%source.FAMILIES[f][0]))&1)
                    mask=1<<pat
                    if known&mask:
                        if bool(ones&mask)!=bool(bit):valid=False;break
                    else:
                        known|=mask
                        if bit:ones|=mask
                if not valid:break
            if valid:out.append((gid,mode,known,ones))
    if not out:raise RuntimeError("T31_ZERO_PUBLIC_WORLDS")
    return out
def possible_worlds(public,source):
    f=public["family"];bits=1<<source.FAMILIES[f][1]
    return [(gid,mode,law) for gid,mode,known,ones in public_posterior(public,source)
      for law in range(1<<bits) if law&known==ones]
def brute_worlds(public,source):
    """Independent brute acceptance: re-simulate every full causal probe."""
    f=public["family"];bits=1<<source.FAMILIES[f][1]
    result=[]
    for gid,graph in enumerate(source.TEMPLATES[f]):
        for mode in (0,1):
            for law in range(1<<bits):
                match=True
                for p in public["probes"]:
                    after=source.update(f,p["pre"],graph,law,mode,p["do"])
                    if any(((after>>node)&1)!=v for node,v in p["observed"]):
                        match=False;break
                if match:result.append((gid,mode,law))
    if not result:raise RuntimeError("T31_BRUTE_NO_WORLDS")
    return result
def future(source,public,world):
    f=public["family"];gid,mode,law=world;p=public["query"]
    state=p["start"]
    graph=source.TEMPLATES[f][gid]
    for step in range(3):
        state=source.update(f,state,graph,law,mode,p["do"] if step==0 else None)
    return tuple((state>>node)&1 for node in p["targets"])
def train_conditional(source,train):
    """Self-supervised PUBLIC-only posterior bit conditional priors."""
    model={}
    for f,(_,k) in source.FAMILIES.items():
        model[f]={"pattern_ones":[1]*(1<<k),
                  "pattern_total":[2]*(1<<k),
                  "graph_counts":[1]*16}
    for public in train:
        f=public["family"];m=model[f]
        for gid,mode,law in possible_worlds(public,source):
            m["graph_counts"][gid*2+mode]+=1
            for i in range(len(m["pattern_ones"])):
                m["pattern_total"][i]+=1
                m["pattern_ones"][i]+=(law>>i)&1
    for f,m in model.items():
        m["log_odds"]=tuple(math.log((m["pattern_ones"][i])/
              (m["pattern_total"][i]-m["pattern_ones"][i]))
              for i in range(len(m["pattern_ones"])))
        m["graph_prior"]=tuple(math.log(x) for x in m["graph_counts"])
    return model
def train_frequency(source,train):
    """Strong public TRAIN-informed marginal law frequency prior."""
    counts={f:{} for f in source.FAMILIES}
    for public in train:
        d=counts[public["family"]]
        for _,_,law in possible_worlds(public,source):d[law]=d.get(law,0)+1
    return counts
def ordering(world,arm,fitted,f,policy):
    gid,mode,law=world
    if arm==ARMS[0]:
        coeff=fitted[f]["log_odds"]
        return (fitted[f]["graph_prior"][gid*2+mode]+
             sum(((law>>i)&1)*v for i,v in enumerate(coeff)), -gid,-mode,-law)
    if arm==ARMS[2]:
        return (fitted[f].get(law,0),-gid,-mode,-law)
    if arm==ARMS[1]:
        if policy=="reverse":return (-law,-gid,-mode)
        if policy=="popcount":return (law.bit_count(),-law,-gid,-mode)
        if policy=="reversed-popcount":return (-law.bit_count(),-law,-gid,-mode)
    return (law,gid,mode)
def run_decisions(public,source,arm,fitted=None,policy="natural",ablate=False):
    worlds=brute_worlds(public,source) if arm==ARMS[3] else possible_worlds(public,source)
    if arm in (ARMS[0],ARMS[2]):
        if ablate:
            ordered=worlds[:]
        else:
            ordered=sorted(worlds,key=lambda w:ordering(w,arm,fitted,public["family"],policy),reverse=True)
    elif arm==ARMS[1]:
        ordered=sorted(worlds,key=lambda w:ordering(w,arm,fitted,public["family"],policy),reverse=True)
    else:ordered=worlds
    outputs=[set(),set(),set()];visits=0
    # EACH target is an independently sound decision; early UNKNOWN never
    # claims a certified bit. The source's unknown law never reaches solver.
    for target in range(3):
        for world in ordered:
            value=future(source,public,world)[target]
            outputs[target].add(value);visits+=1
            if len(outputs[target])==2:break
    return (tuple(next(iter(s)) if len(s)==1 else None for s in outputs),
            tuple(tuple(sorted(s)) for s in outputs),
            visits,len(worlds),
            tuple(ordered[:3]))
def fit_exact_policy(source,train):
    """TRAIN-fitted optimized exact competitor, same PUBLIC information."""
    fitted={}
    for f in source.FAMILIES:
        cohort=[p for p in train if p["family"]==f]
        costs={}
        for policy in ORDERS:
            costs[policy]=sum(run_decisions(p,source,ARMS[1],None,policy)[2] for p in cohort)
        fitted[f]=min(ORDERS,key=lambda p:(costs[p],ORDERS.index(p)))
    return fitted
def train_only():
    source=source_module()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    train=[source.public_train(source.build(name)) for name in source.TRAIN]
    learned=train_conditional(source,train)
    freq=train_frequency(source,train)
    exact_fit=fit_exact_policy(source,train)
    true_retained=certified=unknown=checks=0
    influenced_orders=influenced_work=0
    for identity,public in zip(source.TRAIN,train):
        item=source.build(identity)
        truth=(item["_private_graph"],item["_private_mode"],item["_private_law"])
        fast=possible_worlds(public,source);brute=brute_worlds(public,source)
        if set(fast)!=set(brute):raise RuntimeError("T31_EXACT_POSTERIOR_DISAGREES_WITH_BRUTE")
        if truth not in fast:raise RuntimeError("T31_HIDDEN_LAW_OR_MODE_EXCLUDED")
        true_retained+=1
        learner=run_decisions(public,source,ARMS[0],learned)
        control=run_decisions(public,source,ARMS[1],None,exact_fit[public["family"]])
        ablated=run_decisions(public,source,ARMS[0],learned,ablate=True)
        other=run_decisions(public,source,ARMS[2],freq)
        if learner[0]!=control[0] or learner[0]!=other[0]:
            raise RuntimeError("T31_LEARNED_CERTIFICATE_NOT_EXACT")
        if learner[4]!=ablated[4]:influenced_orders+=1
        if learner[2]!=ablated[2]:influenced_work+=1
        for result,truthbit in zip(learner[0],item["_truth"]):
            if result is None:unknown+=1
            elif result==truthbit:certified+=1
            else:raise RuntimeError("T31_TRAIN_HARD_FALSE_POSITIVE")
            checks+=1
    # Source-free direct 272-function algebraic law space coverage;
    # generated without opening any DEV identity.
    for f,(_,k) in source.FAMILIES.items():
        graph=source.TEMPLATES[f][0]
        for law in range(1<<(1<<k)):
            for mode in (0,1):
                for pre in (0,1,(1<<source.FAMILIES[f][0])-1):
                    expected=0
                    for j,parents in enumerate(graph):
                        bit=((pre>>j)&1)^((law>>source.pattern(pre,parents))&1)
                        if j%2:bit^=mode*((pre>>((j+1)%source.FAMILIES[f][0]))&1)
                        expected|=bit<<j
                    if source.update(f,pre,graph,law,mode)!=expected:
                        raise RuntimeError("T31_UNKNOWN_CAUSAL_ALGEBRA_MISMATCH")
    return {"schema":"wingless.T31.TRAIN-only.v1",
      "TRAIN":40,"DEV_opened":0,"retained_true_graph_mode_law":true_retained,
      "exact_posterior_agreed_with_independent_enumerator":40,
      "counterfactual_decision_equivalences":checks,
      "certified":certified,"safe_abstentions":unknown,
      "uncertainty_fraction":unknown/checks,
      "learned_parameter_changed_executed_proposal_order_cases":influenced_orders,
      "learned_parameter_changed_inference_work_cases":influenced_work,
      "trained_exact_competitor_policies":exact_fit,
      "full_272_law_space_cross_mode_validated":True,
      "valid_for_first_use":true_retained==40 and checks==120 and
        influenced_orders>0 and influenced_work>0 and .15<=unknown/checks<=.90,
      "DEV_replayed":False,"optimizer_admission":False,"scientific_positive":False}
def one(arm,repeat):
    assert arm in ARMS and 0<=repeat<5
    start=time.process_time_ns()
    source=source_module()
    train=[source.public_train(source.build(name)) for name in source.TRAIN]
    learned=train_conditional(source,train) if arm==ARMS[0] else None
    freq=train_frequency(source,train) if arm==ARMS[2] else None
    exact_fit=fit_exact_policy(source,train) if arm==ARMS[1] else {}
    rows=[];ok=unknown=visits=0
    for identity in source.DEV:
        item=source.build(identity)
        public=source.public_dev(item)
        model=learned if arm==ARMS[0] else freq if arm==ARMS[2] else None
        result,sets,work,count,_=run_decisions(public,source,arm,model,
                                  exact_fit.get(public["family"],"natural"))
        num_ok=num_no=0
        for answer,possibilities,truth in zip(result,sets,item["_truth"]):
            if truth not in possibilities:raise RuntimeError("T31_PRIVATE_TRUTH_EXCLUDED:"+identity)
            if answer is None:unknown+=1;num_no+=1
            else:
                if answer!=truth:raise RuntimeError("T31_HARD_FALSE_POSITIVE:"+identity)
                ok+=1;num_ok+=1
        visits+=work
        rows.append({"case":identity,"certified":num_ok,"unknown":num_no,
                    "exact_worlds":count,"world_target_visits":work})
    if ok+unknown!=192:raise RuntimeError("T31_INCOMPLETE_PROSPECTIVE")
    return {"schema":"wingless.T31.original-internal-process.v1","arm":arm,
      "repeat":repeat,"cases":rows,"correctly_certified":ok,
      "safe_abstentions":unknown,"hard_false_positives":0,
      "world_target_visits":visits,"complete_CPU_ns":time.process_time_ns()-start,
      "peak_RSS_bytes":peak_rss(),"new_internal_not_external":True,
      "no_old_source_reuse":True,"optimizer_admission":False}
def run():
    result={}
    for arm in ARMS:
        batch=[]
        for i in range(5):
            x=subprocess.run([sys.executable,str(Path(__file__).resolve()),"--arm",arm,str(i)],
                            text=True,capture_output=True)
            if x.returncode:raise RuntimeError("T31_ORIGINAL_CHILD_FAILED:"+arm+" "+str(i)+"\n"+x.stderr)
            batch.append(json.loads(x.stdout))
        item=batch[0]
        if any(x["cases"]!=item["cases"] for x in batch):
            raise RuntimeError("T31_NONDETERMINISTIC_SOURCE_OR_RESULT")
        result[arm]={"certified":item["correctly_certified"],
          "safe_abstentions":item["safe_abstentions"],
          "hard_false_positives":item["hard_false_positives"],
          "world_target_visits":item["world_target_visits"],
          "per_case":item["cases"],
          "five_CPU_ns":[x["complete_CPU_ns"] for x in batch],
          "five_peak_RSS_bytes":[x["peak_RSS_bytes"] for x in batch],
          "median_CPU_ns":int(statistics.median(x["complete_CPU_ns"] for x in batch)),
          "max_peak_RSS_bytes":max(x["peak_RSS_bytes"] for x in batch)}
    if not all(x["certified"]+x["safe_abstentions"]==192 and x["hard_false_positives"]==0 for x in result.values()):
        raise RuntimeError("T31_PROSPECTIVE_CERTIFICATION_FAILURE")
    return {"schema":"wingless.T31.original-new-causal-conditional-program.v1",
      "new_DEV_worlds":64,"counterfactuals":192,"results":result,
      "scientific_positive":False,"optimizer_admission":False,
      "external_independent_source":False,"no_old_DEV_replayed":True}
if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--train-only":print(json.dumps(train_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":print(json.dumps(one(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:print(json.dumps(run(),sort_keys=True))
    else:raise SystemExit("T31_INVALID_ARGS")
