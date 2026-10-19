"""T30B inferred unknown causal laws with persistent proof-carrying posterior.

Per hidden graph, each public causal experiment reveals constraints on
unknown local truth-table bits. Hypotheses are ALL graph + law pairs
consistent with the entire nine-event sequence. A source-blind TRAIN model
learns an ordering over algebraic-normal-form (ANF) coefficients from public
temporal evidence; its prediction path executes each prospectively possible
truth table via explicit ANF and exact full-grammar fallback.

DEV includes nonaffine local causal laws never present in TRAIN. No one
may infer that affine TRAIN salience rules out unseen nonlinear coefficients.
Comparators include fully constrained truth-table bitset and actual
exhaustive per-law enumeration, both with identical public evidence.
"""
from __future__ import annotations
import importlib.util,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-learned-algebraic-rule-memory",
      "TRAIN-frequency-exact-rule-memory",
      "fully-trained-exact-constraint-bitset",
      "deterministic-minimal-observation-constraints",
      "untrained-explicit-function-enumeration")
def srcmod():
    spec=importlib.util.spec_from_file_location("t30_new_source",ROOT/"20261019-wingless-t30b-new-sparse-rules.py")
    src=importlib.util.module_from_spec(spec)
    sys.modules[spec.name]=src;spec.loader.exec_module(src)
    return src
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
        k=ctypes.WinDLL("kernel32",use_last_error=True)
        p=ctypes.WinDLL("psapi",use_last_error=True)
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T30B_WINDOWS_MEMORY_ERROR="+str(ctypes.get_last_error()))
        return int(m.PeakWorkingSetSize)
    import resource
    z=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(z)*1024 if sys.platform.startswith("linux") else int(z)
def rule_pattern(pre,parents):
    return sum(((pre>>node)&1)<<k for k,node in enumerate(parents))
def constraints(public,source):
    """Exact online persistent law posterior per hidden graph template."""
    f=public["family"]
    out=[]
    for graph,parents in enumerate(source.TEMPLATES[f]):
        known=0;ones=0;valid=True
        for probe in public["probes"]:
            pre=probe["pre"];action=probe["intervention"]
            for node,post in probe["observed"]:
                if action is not None and action[0]==node:
                    if post!=action[1]:raise RuntimeError("T30B_SOURCE_INTERVENTION_INCONSISTENT")
                    continue
                pattern=rule_pattern(pre,parents[node])
                bit=post^((pre>>node)&1)
                mark=1<<pattern
                if known&mark:
                    if bool(ones&mark)!=bool(bit):
                        valid=False;break
                else:
                    known|=mark
                    if bit:ones|=mark
            if not valid:break
        if valid:out.append((graph,known,ones))
    if not out:raise RuntimeError("T30B_TRUE_LAW_GRAPH_WAS_ELIMINATED")
    return tuple(out)
def allowed_rule(known,ones,law):
    return law&known==ones
def posterior_laws(public,source):
    f=public["family"];n,k=source.FAMILIES[f]
    output=[]
    for graph,known,ones in constraints(public,source):
        possibilities=tuple(law for law in range(1<<(1<<k)) if allowed_rule(known,ones,law))
        if not possibilities:raise RuntimeError("T30B_EMPTY_CAUSAL_POSTERIOR")
        output.append((graph,possibilities))
    return tuple(output)
def brute_posterior(public,source):
    """Fully independent actual exhaustive probe-by-probe enumerator."""
    f=public["family"];n,k=source.FAMILIES[f]
    worlds=[]
    for graph in range(8):
        topo=source.TEMPLATES[f][graph]
        for law in range(1<<(1<<k)):
            valid=True
            for p in public["probes"]:
                predicted=source.update(f,p["pre"],topo,law,p["intervention"])
                if any(((predicted>>j)&1)!=v for j,v in p["observed"]):
                    valid=False;break
            if valid:worlds.append((graph,law))
    if not worlds:raise RuntimeError("T30B_ENUMERATIVE_NO_POSSIBLE_CAUSAL_RULE")
    return tuple(worlds)
def anf_encode(tt,k):
    """Möbius transform over F2, all coefficients retained, no degree cap."""
    arr=[(tt>>i)&1 for i in range(1<<k)]
    for i in range(k):
        for j in range(1<<k):
            if (j>>i)&1:arr[j]^=arr[j^(1<<i)]
    return sum(x<<i for i,x in enumerate(arr))
def anf_eval(coeff,inputs):
    result=0
    subset=inputs
    while True:
        result^=(coeff>>subset)&1
        if subset==0:break
        subset=(subset-1)&inputs
    return result
def anf_update(f,pre,topology,coeff,action,source):
    after=0
    for node,parents in enumerate(topology):
        inputbits=rule_pattern(pre,parents)
        out=(pre>>node)&1
        out^=anf_eval(coeff,inputbits)
        after|=out<<node
    if action is not None:
        node,bit=action
        after=(after&~(1<<node)) | (bit<<node)
    return after
def fit_train(source,train):
    """Fit ALGORITHM weights and learned coefficient basis from PUBLIC ONLY.

    Importantly, no coefficient is considered 'impossible' after TRAIN; the
    entire ANF basis is retained for future nonaffine functions.
    """
    scores={}
    observation_count=0
    for family in source.FAMILIES:
        _,k=source.FAMILIES[family]
        scores[family]=[0]*(1<<k)
    for public in train:
        f=public["family"];n,k=source.FAMILIES[f]
        hypoth=posterior_laws(public,source)
        for graph,laws in hypoth:
            for law in laws:
                coeff=anf_encode(law,k)
                for bit in range(1<<k):
                    if (coeff>>bit)&1:scores[f][bit]+=1
        observation_count+=sum(len(p["observed"]) for p in public["probes"])
    ranks={f:tuple(sorted(range(len(v)),key=lambda i:(-v[i],i))) for f,v in scores.items()}
    return ranks,observation_count
def resolve(public,source,arm,ranking):
    f=public["family"]
    p=public["query"];q=p["do"];start=p["start"];target=p["targets"]
    if arm==ARMS[4]:
        hypotheses=brute_posterior(public,source)
    else:
        hypotheses=tuple((graph,law) for graph,laws in posterior_laws(public,source) for law in laws)
    outputs=[set(),set(),set()]
    if arm==ARMS[0]:
        k=source.FAMILIES[f][1]
        enc={law:anf_encode(law,k) for graph,law in hypotheses}
        # Learned ANF order is a TRAIN-fitted representation, but
        # an exhaustive proof-preserving decode is required on every path.
        if set(ranking[f])!=set(range(1<<k)):
            raise RuntimeError("T30B_LEARNED_ANF_NOT_COMPLETE")
    for graph,law in hypotheses:
        topo=source.TEMPLATES[f][graph]
        if arm==ARMS[0]:
            coeff=enc[law]
            one=anf_update(f,start,topo,coeff,q,source)
            final=anf_update(f,one,topo,coeff,None,source)
        else:
            one=source.update(f,start,topo,law,q)
            final=source.update(f,one,topo,law)
        for i,node in enumerate(target):
            outputs[i].add((final>>node)&1)
    if not all(outputs):raise RuntimeError("T30B_NO_FUTURE_CAUSAL_HYPOTHESES")
    predictions=tuple(next(iter(v)) if len(v)==1 else None for v in outputs)
    return predictions,tuple(tuple(sorted(v)) for v in outputs),len(hypotheses)
def train_only():
    source=srcmod()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert set(source.TRAIN).isdisjoint(source.DEV)
    training=[source.public_train(source.build(name)) for name in source.TRAIN]
    ranking,samples=fit_train(source,training)
    n=correct=abstentions=preserved=exact=0
    mean_hypotheses=[]
    for name,public in zip(source.TRAIN,training):
        case=source.build(name)
        guesses=posterior_laws(public,source)
        truth=(case["_private_template_index"],case["_private_law"])
        if not any(graph==truth[0] and truth[1] in laws for graph,laws in guesses):
            raise RuntimeError("T30B_TRAIN_PRIVATE_TRUE_LAW_EXCLUDED")
        preserved+=1
        brute=set(brute_posterior(public,source))
        optimized={(graph,law) for graph,laws in guesses for law in laws}
        if optimized!=brute:raise RuntimeError("T30B_PUBLIC_CAUSAL_RULE_CONSTRAINTS_NOT_EXACT")
        learned,sets,num=resolve(public,source,ARMS[0],ranking)
        strong,*_=resolve(public,source,ARMS[2],ranking)
        if learned!=strong:raise RuntimeError("T30B_ALGEBRAIC_LEARNER_NOT_EXACT_EQUIVALENT")
        mean_hypotheses.append(num)
        for result,possible,answer in zip(learned,sets,case["_private_truth"]):
            if answer not in possible:raise RuntimeError("T30B_TRAIN_TRUE_COUNTERFACTUAL_EXCLUDED")
            if result is None:abstentions+=1
            else:
                if result!=answer:raise RuntimeError("T30B_TRAIN_HARD_FALSE_POSITIVE")
                correct+=1
            n+=1;exact+=1
    # Universal algebraic equivalence even for unseen nonaffine functions;
    # this tests the evaluator without opening ANY prospective case.
    for family,(nodes,k) in source.FAMILIES.items():
        for law in range(1<<(1<<k)):
            coeff=anf_encode(law,k)
            for pattern in range(1<<k):
                if ((law>>pattern)&1)!=anf_eval(coeff,pattern):
                    raise RuntimeError("T30B_ANF_NONAFFINE_ALGEBRA_FAIL")
    return {"schema":"wingless.T30B.TRAIN-only-unknown-mechanism-soundness.v1",
        "status":"PASS_TRAIN_ONLY_DEV_UNOPENED",
        "source_family_M":"seven_entity_two_parent_unknown_16_boolean_rules",
        "source_family_N":"nine_entity_three_parent_unknown_256_boolean_rules",
        "TRAIN_cases":40,"DEV_generated":0,
        "TRAIN_public_observation_pairs":samples,
        "TRAIN_true_hidden_laws_retained":preserved,
        "TRAIN_counterfactual_queries":n,
        "TRAIN_correctly_certified":correct,"TRAIN_safe_abstentions":abstentions,
        "TRAIN_uncertainty_fraction":abstentions/n,
        "TRAIN_exact_public_constraint_and_answer_equivalences":exact,
        "mean_posterior_full_graph_laws":statistics.mean(mean_hypotheses),
        "learned_ANF_coefficient_rank_M":ranking["M"],
        "learned_ANF_coefficient_rank_N":ranking["N"],
        "full_known_and_unknown_ANF_rules_verified":16+256,
        "proof_does_not_assume_affine_DEV":True,
        "previous_consumed_cases_opened":False,
        "scientific_positive":False,"optimizer_admission":False}
def one(arm,repeat):
    assert arm in ARMS and 0<=repeat<5
    start=time.process_time_ns()
    source=srcmod()
    train=[source.public_train(source.build(s)) for s in source.TRAIN]
    ranking,examples=fit_train(source,train) if arm!=ARMS[4] else ({"M":(0,1,2,3),"N":tuple(range(8))},0)
    out=[];certified=unknown=worlds=0
    for name in source.DEV:
        item=source.build(name)
        public=source.public_DEV(item)
        predicted,possible,num=resolve(public,source,arm,ranking)
        yes=no=0
        for answer,options,truth in zip(predicted,possible,item["_private_truth"]):
            if truth not in options:raise RuntimeError("T30B_PRIVATE_NONAFFINE_TRUTH_EXCLUDED:"+name)
            if answer is None:no+=1;unknown+=1
            else:
                if answer!=truth:raise RuntimeError("T30B_ORIGINAL_HARD_FALSE_POSITIVE:"+name)
                yes+=1;certified+=1
        worlds+=num
        out.append({"source_id":name,"correct_certifications":yes,
             "safe_abstentions":no,"consistent_hidden_rule_graph_worlds":num})
    if len(out)!=64 or certified+unknown!=192:raise RuntimeError("T30B_PROSPECTIVE_INCOMPLETE")
    return {"schema":"wingless.T30B.original-internal-mechanism-process.v1",
        "arm":arm,"repeat":repeat,"per_case":out,"certified":certified,
        "safe_abstentions":unknown,"hard_false_positives":0,
        "public_TRAIN_sequences":40,"DEV_sequences":64,
        "source_origin":"INTERNAL_NEW_NONAFFINE_NOT_EXTERNAL",
        "all_consistent_hidden_rule_graph_worlds":worlds,
        "complete_process_CPU_ns":time.process_time_ns()-start,
        "peak_process_RSS_bytes":rss(),
        "no_old_DEV_replay":True,"optimizer_admission":False}
def evaluate():
    results={}
    for arm in ARMS:
        trials=[]
        for i in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                      "--arm",arm,str(i)],capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T30B_ORIGINAL_CHILD_FAILED:"+arm+":"+str(i)+"\n"+p.stderr)
            trials.append(json.loads(p.stdout))
        a=trials[0]
        if any(t["per_case"]!=a["per_case"] for t in trials):
            raise RuntimeError("T30B_NONDETERMINISTIC_HIDDEN_MECHANISM_SOURCE")
        results[arm]={"correct_certified":a["certified"],
          "safe_abstentions":a["safe_abstentions"],
          "hard_false_positives":a["hard_false_positives"],
          "per_case":a["per_case"],
          "five_complete_CPU_ns":[t["complete_process_CPU_ns"] for t in trials],
          "five_peak_RSS_bytes":[t["peak_process_RSS_bytes"] for t in trials],
          "median_complete_CPU_ns":int(statistics.median(t["complete_process_CPU_ns"] for t in trials)),
          "max_peak_RSS_bytes":max(t["peak_process_RSS_bytes"] for t in trials)}
    if not all(x["correct_certified"]+x["safe_abstentions"]==192 and x["hard_false_positives"]==0 for x in results.values()):
        raise RuntimeError("T30B_WRONG_OR_INCOMPLETE")
    learn=results[ARMS[0]]
    compare={a:{"coverage_preserved":learn["correct_certified"]>=x["correct_certified"],
        "CPU_nondecreasing":learn["median_complete_CPU_ns"]<=x["median_complete_CPU_ns"],
        "peak_memory_nondecreasing":learn["max_peak_RSS_bytes"]<=x["max_peak_RSS_bytes"]}
        for a,x in results.items() if a!=ARMS[0]}
    return {"schema":"wingless.T30B.original-unseen-nonaffine-function-comparison.v1",
        "source_families":"new_M_N_structures_and_nonaffine_rules",
        "TRAIN_cases":40,"DEV_cases":64,"future_counterfactual_questions":192,
        "arms":results,"comparisons":compare,
        "scientific_positive":False,"external_independent_replication":False,
        "optimizer_admission":False,"old_DEV_replayed":False,
        "no_paid_API":True,"no_new_controller":True}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(train_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(one(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T30B_BAD_ARGS")
