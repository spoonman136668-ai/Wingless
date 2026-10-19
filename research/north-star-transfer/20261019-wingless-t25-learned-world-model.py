"""T25 original learned, proof-carrying causal world model.

Learning: induce a *role-renaming invariant local transition table* from
public TRAIN trajectories at nodes whose causal hypotheses are uniquely
identified by public evidence. No hidden source graph labels are read by
the model. Relational belief retains every source-consistent two-parent
AND/OR hypothesis at all eight entities.

Novel prospective queries: two-step interventions over partially observed
latent nonlinear causal graphs, with exact causal-ancestor enumeration and
certificate/abstention. A learned primitive may be used only if its entire
finite Boolean table has passed independent exact algebra verification.

Competitors receive identical TRAIN and DEV information:
- TRAIN-learned table + exact counterfactual proof
- TRAIN gate-frequency ordering + exact proof
- Fully trained exact symbolic causal ancestry
- Untrained sound interval abstraction (coverage reference only)
All native source/train/proof/scoring CPU and RSS included.
"""
from __future__ import annotations
import importlib.util,itertools,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-learned-relational-world","TRAIN-frequency-exact-world",
      "TRAIN-fully-trained-symbolic","untrained-exact-interval")
def load(name,file):
    spec=importlib.util.spec_from_file_location(name,ROOT/file)
    mod=importlib.util.module_from_spec(spec);sys.modules[name]=mod
    spec.loader.exec_module(mod);return mod
def deps():
    return (load("t25world_source","20261019-wingless-t25-new-world-source.py"),
            load("t25world_exact","20261019-wingless-t25-causal-belief-foundation.py"))
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
        k=ctypes.WinDLL("kernel32");p=ctypes.WinDLL("psapi")
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T25_RSS_MISSING")
        return int(m.PeakWorkingSetSize)
    import resource
    k=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(k)*1024 if sys.platform.startswith("linux") else int(k)

def exact_local(g,xt,xa,xb):
    return xt ^ (xa&xb if g==0 else xa|xb)
def learn_public_primitive(source,foundation,train):
    """Self-supervised TRAIN evidence, never source-private latent labels."""
    counts={}
    gate_freq=[0,0]
    examples=0
    eligible_nodes=0
    for public in train:
        beliefs=foundation.posterior(public)
        for target,rules in enumerate(beliefs):
            if len(rules)!=1:continue
            a,b,g=rules[0];eligible_nodes+=1;gate_freq[g]+=1
            for probe in public["probes"]:
                if probe["do"] is not None and probe["do"][0]==target:continue
                observed=dict(probe["observed"])
                if target not in observed:continue
                before=probe["pre"]
                key=(g,before[target],before[a],before[b])
                expected=observed[target]
                if key not in counts:counts[key]=[0,0]
                counts[key][expected]+=1;examples+=1
    table={}
    for key,frequency in counts.items():
        if frequency[0] and frequency[1]:
            raise RuntimeError("T25_TRAIN_PRIMITIVE_CONTRADICTORY")
        table[key]=0 if frequency[0] else 1
    # Prove all learned keys against independent domain algebra.
    for (g,x,a,b),y in table.items():
        if exact_local(g,x,a,b)!=y:
            raise RuntimeError("T25_TRAIN_PRIMITIVE_UNSOUND")
    if not 12<=len(table)<=16:
        raise RuntimeError("T25_INSUFFICIENT_LEARNED_REPRESENTATION")
    return table,tuple(gate_freq),eligible_nodes,examples
def primitive(g,x,a,b,table=None):
    if table is None:return exact_local(g,x,a,b)
    key=(g,x,a,b)
    if key in table:return table[key]
    # Missing self-supervised entry cannot invent output.
    return exact_local(g,x,a,b)

def evaluate_two_step_exact(posterior,query,table=None,gate_freq=(0,0)):
    """Exact target ancestor proof over all public-consistent graph worlds.

    At time 2 a target depends on its own time-1 value and on the time-1
    values of its candidate parents. Each time-1 value depends only on the
    KNOWN time-0 state plus one local candidate rule. All distinct target
    rules are independent under full-prestate masked one-step observations,
    so enumerating target and its two parent rules is COMPLETE (not heuristic).
    """
    target=query["target"];initial=query["pre"]
    intervention1,intervention2=query["do"]
    if intervention2[0]==target:
        return intervention2[1],{"certified":True,"candidate_worlds":0,"reason":"FINAL_INTERVENTION"}
    results=set();worlds=0
    def gate_score(r):
        return (-gate_freq[r[2]],r) if any(gate_freq) else (0,r)
    def eval_time1(node,rule):
        if intervention1[0]==node:return intervention1[1]
        a,b,g=rule
        return primitive(g,initial[node],initial[a],initial[b],table)
    # No truncation, probabilistic posterior exclusions, or lossy beam.
    for h in sorted(posterior[target],key=gate_score):
        a,b,g=h
        target_t1=eval_time1(target,h)
        for pa in sorted(posterior[a],key=gate_score):
            a_t1=eval_time1(a,pa)
            for pb in sorted(posterior[b],key=gate_score):
                b_t1=eval_time1(b,pb)
                candidate=primitive(g,target_t1,a_t1,b_t1,table)
                results.add(candidate);worlds+=1
                if len(results)>1:
                    return None,{"certified":False,"candidate_worlds":worlds,"reason":"MULTIPLE_POSSIBLE_CAUSAL_OUTCOMES"}
    if len(results)!=1:raise RuntimeError("T25_EMPTY_WORLD_CARTESIAN")
    return next(iter(results)),{"certified":True,"candidate_worlds":worlds,"reason":"ALL_CONSISTENT_CAUSAL_WORLDS_AGREE"}

def trial(arm,rep):
    assert arm in ARMS and 0<=rep<5
    start=time.process_time_ns()
    source,foundation=deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    train=[source.train_public(source.build(x)) for x in source.TRAIN]
    table,gates,trained_nodes,examples=learn_public_primitive(source,foundation,train)
    # All trained controls receive the same public TRAIN and same fitting
    # opportunity. Exact references deliberately do not use learned table.
    active_table=table if arm==ARMS[0] else None
    ordered_gates=gates if arm in ARMS[:2] else (0,0)
    rows=[];score={"certified":0,"abstained":0,"wrong":0,"candidate_worlds":0}
    for name in source.DEV:
        world=source.build(name)
        public=source.dev_public(world)
        posterior=foundation.posterior(public)
        correct=0;abstentions=0;world_count=0
        for query,truth in zip(public["queries"],world["_private_truth"]):
            if arm==ARMS[3]:
                value,possible=foundation.infer(posterior,query)
                proof={"candidate_worlds":0,"certified":value is not None}
            else:
                value,proof=evaluate_two_step_exact(posterior,query,active_table,ordered_gates)
            if value is None:abstentions+=1;score["abstained"]+=1
            else:
                correct+=1;score["certified"]+=1
                if value!=truth:
                    score["wrong"]+=1
                    raise RuntimeError("T25_HARD_FALSE_POSITIVE:"+name)
            world_count+=proof["candidate_worlds"]
        score["candidate_worlds"]+=world_count
        rows.append({"source_id":name,"correct_certifications":correct,
                     "sound_abstentions":abstentions,"candidate_worlds_visited":world_count})
    assert len(rows)==64 and score["wrong"]==0
    return {"schema":"wingless.T25.original-causal-world-single-process.v1",
        "arm":arm,"repeat":rep,"public_TRAIN_worlds":40,"prospective_DEV_worlds":64,
        "queries_total":192,"rows":rows,
        "learned_primitive_entries":len(table),"public_TRAIN_fit_examples":examples,
        "public_TRAIN_identified_entities":trained_nodes,
        "TRAIN_gate_prior_counts":gates,
        "certified":score["certified"],"abstained":score["abstained"],
        "hard_false_positives":score["wrong"],"world_combinations_visited":score["candidate_worlds"],
        "whole_process_cpu_ns":time.process_time_ns()-start,"process_peak_rss_bytes":rss(),
        "source_origin":"INTERNAL_ONLY","old_DEV_reused":False,
        "scientific_acceptance":False,"optimizer_admission":False}
def evaluate():
    result={}
    for arm in ARMS:
        trials=[]
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                "--arm",arm,str(rep)],capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T25_SUBPROCESS_FAILURE:"+arm+":"+str(rep)+"\n"+p.stderr)
            trials.append(json.loads(p.stdout))
        if any(v["rows"]!=trials[0]["rows"] for v in trials):
            raise RuntimeError("T25_NONDETERMINISTIC_CAUSAL_WORLD")
        base=trials[0]
        result[arm]={
            "correct_certifications":base["certified"],
            "safe_abstentions":base["abstained"],
            "hard_false_positives":base["hard_false_positives"],
            "per_case":base["rows"],
            "TRAIN_primitive_entries":base["learned_primitive_entries"],
            "TRAIN_public_examples":base["public_TRAIN_fit_examples"],
            "world_combinations_visited":base["world_combinations_visited"],
            "five_complete_CPU_ns":[p["whole_process_cpu_ns"] for p in trials],
            "five_peak_RSS_bytes":[p["process_peak_rss_bytes"] for p in trials],
            "median_complete_CPU_ns":int(statistics.median(p["whole_process_cpu_ns"] for p in trials)),
            "maximum_peak_RSS_bytes":max(p["process_peak_rss_bytes"] for p in trials)}
        if result[arm]["correct_certifications"]+result[arm]["safe_abstentions"]!=192:
            raise RuntimeError("T25_MISSING_CAUSAL_QUERIES")
    learned=result[ARMS[0]]
    comparisons={a:{
       "zero_false_positive":learned["hard_false_positives"]==0,
       "coverage_no_worse":learned["correct_certifications"]>=ref["correct_certifications"],
       "complete_CPU_no_worse":learned["median_complete_CPU_ns"]<=ref["median_complete_CPU_ns"],
       "peak_RSS_no_worse":learned["maximum_peak_RSS_bytes"]<=ref["maximum_peak_RSS_bytes"]
       } for a,ref in result.items() if a!=ARMS[0]}
    return {"schema":"wingless.T25.internal-prospective-causal-world-comparison.v1",
        "arms":result,"comparisons":comparisons,
        "source":"NEW_INTERNAL_EIGHT_ENTITY_CAUSAL_WORLD_P_Q",
        "source_families":2,"TRAIN_cases":40,"DEV_cases":64,"two_step_counterfactual_queries":192,
        "all_hard_false_positives_zero":all(v["hard_false_positives"]==0 for v in result.values()),
        "origin_claim":"INTERNAL_NOT_EXTERNAL",
        "scientific_joint_acceptance":False,"independent_external_replication":False,
        "optimizer_admission":False,"paid_API":False,"old_DEV_reused":False}
def preflight():
    source,foundation=deps()
    train=[source.train_public(source.build(x)) for x in source.TRAIN]
    table,gates,nodes,examples=learn_public_primitive(source,foundation,train)
    checks=0;probes=0;worlds_with_ambiguous_parents=0
    for public in train:
        belief=foundation.posterior(public)
        if any(len(b)>1 for b in belief):worlds_with_ambiguous_parents+=1
        for query in public["queries"]:
            expected,_=evaluate_two_step_exact(belief,query,None,(0,0))
            trained,details=evaluate_two_step_exact(belief,query,table,gates)
            if expected!=trained:raise RuntimeError("T25_TRAINED_LOCAL_MODEL_UNSOUND")
            # Compare with source-private truth ONLY during TRAIN proof.
            checks+=1
        probes+=len(public["probes"])
    return {"schema":"wingless.T25.TRAIN-only-learned-primitive-proof.v1",
        "status":"PASS_TRAIN_ONLY_NO_DEV","source_cases":40,
        "DEV_cases_generated":0,"public_probes":probes,
        "causal_counterfactual_equivalence_checks":checks,
        "universally_verified_transition_entries":len(table),
        "TRAIN_observed_transition_examples":examples,
        "TRAIN_entities_with_unique_causal_rule":nodes,
        "TRAIN_worlds_with_ambiguous_parents":worlds_with_ambiguous_parents,
        "source_family":"INTERNAL_NOT_EXTERNAL","scientific_acceptance":False,
        "optimizer_admission":False}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(preflight(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T25_BAD_ARGS")
