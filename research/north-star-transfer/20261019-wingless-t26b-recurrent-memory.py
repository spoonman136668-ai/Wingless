"""T26B proof-carrying recurrent causal hypothesis memory in fresh U/V 5-frame worlds.

The latent state is a set of (persistent hidden graph, gate, initial state),
filtered over the entire chronological masked observation history. All
choices of persistent graph and initial state are retained without pruning.
Learn a source-induced table of local transition relations from public
TRAIN histories only, reprove every table entry against independent finite
Boolean algebra, and use that table on novel sequences with sound fallback.

Four equally TRAIN-informed controls:
  TRAIN-learned-recurrent-causal-memory (learned exact transition table)
  TRAIN-frequency-recurrent-exact-control (mode prior for traversal)
  fully-trained-symbolic-recurrent-Bayesian-control (exact symbolic filter)
  untrained-complete-recurrent-Bayesian-control (same exact completeness)
No guessed outputs, source-family oracle access, or protected old DEV usage.
"""
from __future__ import annotations
import importlib.util,itertools,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-learned-recurrent-causal-memory",
      "TRAIN-frequency-recurrent-exact-control",
      "fully-trained-symbolic-recurrent-Bayesian-control",
      "untrained-complete-recurrent-Bayesian-control")
TABLE_CAP=16
def load(name,file):
    spec=importlib.util.spec_from_file_location(name,ROOT/file)
    module=importlib.util.module_from_spec(spec);sys.modules[name]=module
    spec.loader.exec_module(module)
    return module
def source():
    return load("t26_source","20261019-wingless-t26b-fresh-sequences.py")
def process_rss():
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
        mem=Mem();mem.cb=ctypes.sizeof(mem)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(mem),mem.cb):
            raise RuntimeError("T26_RSS_UNAVAILABLE")
        return int(mem.PeakWorkingSetSize)
    import resource
    r=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(r)*1024 if sys.platform.startswith("linux") else int(r)

def next_exact(bits,parents,gate,intervention=None,table=None):
    future=0
    for node,(a,b) in enumerate(parents):
        own=(bits>>node)&1;left=(bits>>a)&1;right=(bits>>b)&1
        key=(gate,own,left,right)
        if table is not None and key in table:
            result=table[key]
        else:
            result=own^(left&right if gate==0 else left|right)
        future|=result<<node
    if intervention is not None:
        j,value=intervention
        future=(future&~(1<<j))|(value<<j)
    return future

def filter_chronology(public,src,table=None,mode_rank=(0,0)):
    """Exact unpruned filtering through 7 observed frames and 6 actions."""
    frames=public["observations"];actions=public["history_actions"]
    if len(frames)!=5 or len(actions)!=4:
        raise RuntimeError("T26_WRONG_TIME_HORIZON")
    if any(mask.bit_count()!=3 or observed & ~mask for mask,observed in frames):
        raise RuntimeError("T26_BAD_SENSOR_OBSERVATION")
    survived=[];visited=0
    modes=sorted(range(2),key=lambda g:(-mode_rank[g],g))
    for gate in modes:
        for t,parents in enumerate(src.TEMPLATES):
            for initial in range(64):
                visited+=1
                state=initial
                mask,obs=frames[0]
                if state&mask!=obs:continue
                valid=True
                for idx,action in enumerate(actions):
                    state=next_exact(state,parents,gate,action,table)
                    mask,obs=frames[idx+1]
                    if state&mask!=obs:
                        valid=False;break
                if valid:survived.append((t,gate,initial,state))
    if not survived:
        raise RuntimeError("T26_EMPTY_PUBLIC_CAUSAL_POSTERIOR")
    return tuple(survived),visited

def learn_temporal_program(src,public_train):
    """TRAIN-only belief-conditioned finite temporal table, no hidden labels.

    Every surviving public-compatible temporal world contributes an observed
    transition instance only if the post-state bit is directly revealed by
    a public sensor and not overwritten by intervention. Candidate truth
    table keys are mode+own+two parent bits. Table values MUST be separately
    proved by the finite truth evaluator before being trusted.
    """
    counts=Counter();mode_freq=[0,0];updates=0;posterior_sizes=[]
    for public in public_train:
        posterior,_=filter_chronology(public,src)
        posterior_sizes.append(len(posterior))
        for template_idx,gate,initial,terminal in posterior:
            mode_freq[gate]+=1
            state=initial
            par=src.TEMPLATES[template_idx]
            for step,action in enumerate(public["history_actions"]):
                nxt=next_exact(state,par,gate,action)
                next_mask,next_obs=public["observations"][step+1]
                for node,(a,b) in enumerate(par):
                    if action is not None and action[0]==node:continue
                    if not (next_mask>>node)&1:continue
                    key=(gate,(state>>node)&1,(state>>a)&1,(state>>b)&1)
                    observed=(next_obs>>node)&1
                    counts[(key,observed)]+=1
                    updates+=1
                state=nxt
    table={}
    # Select at most 16 possible Boolean rows by TRAIN observed support.
    keys={key for key,value in counts}
    ranked=sorted(keys,key=lambda x:(-sum(counts[x,y] for y in (0,1)),x))
    for key in ranked[:TABLE_CAP]:
        observed_vals=tuple(y for y in (0,1) if counts[(key,y)]>0)
        if len(observed_vals)!=1:
            raise RuntimeError("T26_INCONSISTENT_PUBLIC_TEMPORAL_LABEL")
        table[key]=observed_vals[0]
    for (gate,own,left,right),answer in table.items():
        reference=own^(left&right if gate==0 else left|right)
        if answer!=reference:
            raise RuntimeError("T26_LEARNED_TEMPORAL_TRANSITION_FALSE_CERTIFICATE")
    return table,tuple(mode_freq),updates,tuple(posterior_sizes)

def counterfactual(posterior,public,src,table=None):
    answers=[set(),set(),set()]
    do=public["future_intervention"];targets=public["targets"]
    for idx,gate,initial,last_state in posterior:
        parents=src.TEMPLATES[idx]
        one=next_exact(last_state,parents,gate,do,table)
        two=next_exact(one,parents,gate,None,table)
        for j,target in enumerate(targets):
            answers[j].add((two>>target)&1)
    result=tuple(next(iter(s)) if len(s)==1 else None for s in answers)
    return result,tuple(tuple(sorted(s)) for s in answers)

def train_only():
    src=source()
    if len(src.TRAIN)!=40 or len(src.DEV)!=64 or set(src.TRAIN)&set(src.DEV):
        raise RuntimeError("T26_SOURCE_COHORT_VIOLATION")
    public_train=[src.train_public(src.sample(i)) for i in src.TRAIN]
    table,mode_freq,updates,sizes=learn_temporal_program(src,public_train)
    confirmed=0;ambiguous=0;world_contained=0;learned_equal=0;tested=0
    for label,public in zip(src.TRAIN,public_train):
        item=src.sample(label)
        posterior,_=filter_chronology(public,src)
        if not any(t==item["_private_topology"] and g==item["_private_mode"]
                   and initial==item["_private_initial"]
                   for t,g,initial,last in posterior):
            raise RuntimeError("T26_TRUE_TEMPORAL_WORLD_ELIMINATED")
        world_contained+=1
        exact,possible=counterfactual(posterior,public,src)
        trained_posterior,_=filter_chronology(public,src,table,mode_freq)
        trained,_=counterfactual(trained_posterior,public,src,table)
        if trained!=exact:
            raise RuntimeError("T26_TRAIN_LEARNED_RECURRENT_STATE_NOT_EXACT")
        learned_equal+=3
        for got,options,truth in zip(exact,possible,item["_truth"]):
            if truth not in options:raise RuntimeError("T26_GROUND_TRUTH_EXCLUDED")
            if got is None:ambiguous+=1
            else:
                confirmed+=1
                if got!=truth:raise RuntimeError("T26_FALSE_POSITIVE")
            tested+=1
    return {"schema":"wingless.T26B.TRAIN-only-prospective-readiness.v1",
        "status":"PASS_TRAIN_ONLY_NO_DEV","TRAIN_cases":40,"DEV_cases_generated":0,
        "history_frames":5,"masked_sensor_count_per_frame":3,"entities":6,
        "world_hypotheses_at_birth":12*2*64,
        "TRAIN_worlds_truth_retained":world_contained,
        "TRAIN_future_target_queries":tested,
        "TRAIN_soundly_certified":confirmed,"TRAIN_safely_abstained":ambiguous,
        "TRAIN_ambiguous_target_fraction":ambiguous/tested,
        "TRAIN_recurrent_exact_equivalences":learned_equal,
        "learned_temporal_entries":len(table),"TRAIN_public_temporal_examples":updates,
        "TRAIN_mode_rank":mode_freq,
        "mean_surviving_hypotheses":statistics.mean(sizes),
        "max_surviving_hypotheses":max(sizes),
        "proof_valid":True,"source_origin":"INTERNAL_NOT_EXTERNALLY_REPLICATED",
        "scientific_acceptance":False,"optimizer_admission":False,
        "previous_consumed_DEV_used":False}

def trial(arm,repetition):
    if arm not in ARMS or repetition not in range(5):
        raise ValueError("T26_INVALID_ARM_OR_REPETITION")
    t0=time.process_time_ns()
    src=source()
    train=[src.train_public(src.sample(label)) for label in src.TRAIN]
    table,mode_freq,updates,sizes=learn_temporal_program(src,train)
    active_table=table if arm==ARMS[0] else None
    active_rank=mode_freq if arm in ARMS[:2] else (0,0)
    rows=[];correct=0;abstained=0;visited=0
    for label in src.DEV:
        item=src.sample(label)
        visible=src.dev_public(item)
        posterior,count=filter_chronology(visible,src,active_table,active_rank)
        result,sets=counterfactual(posterior,visible,src,active_table)
        yes=0;unknown=0
        for answer,possible,truth in zip(result,sets,item["_truth"]):
            if truth not in possible:raise RuntimeError("T26_HIDDEN_TRUTH_NOT_IN_PROOF")
            if answer is None:unknown+=1
            else:
                if answer!=truth:raise RuntimeError("T26_HARD_FALSE_POSITIVE:"+label)
                yes+=1
        correct+=yes;abstained+=unknown;visited+=count
        rows.append({"id":label,"certified":yes,"abstained":unknown,
                     "hypotheses_remaining":len(posterior),"candidate_states_checked":count})
    if len(rows)!=64 or correct+abstained!=192:
        raise RuntimeError("T26_PROSPECTIVE_INCOMPLETE")
    return {"schema":"wingless.T26B.original-internal-sequence-process.v1",
        "arm":arm,"rep":repetition,"rows":rows,
        "certified":correct,"abstained":abstained,"hard_false_positives":0,
        "public_TRAIN_sequences":40,"prospective_sequences":64,
        "TRAIN_learned_temporal_entries":len(table),"TRAIN_examples":updates,
        "initial_global_hypotheses":1536,"candidate_states_visited":visited,
        "full_process_cpu_ns":time.process_time_ns()-t0,
        "peak_process_rss_bytes":process_rss(),
        "source_origin":"INTERNAL_NOT_EXTERNAL",
        "not_independently_replicated":True,
        "old_DEV_not_reopened":True,"optimizer_admission":False}

def evaluate():
    result={}
    for arm in ARMS:
        children=[]
        for rep in range(5):
            c=subprocess.run([sys.executable,str(Path(__file__).resolve()),"--arm",arm,str(rep)],
                             capture_output=True,text=True)
            if c.returncode:
                raise RuntimeError("T26_ORIGINAL_CHILD_FAILED:"+arm+":"+str(rep)+"\n"+c.stderr)
            children.append(json.loads(c.stdout))
        if any(v["rows"]!=children[0]["rows"] for v in children):
            raise RuntimeError("T26_NONDETERMINISTIC_CASE_OUTCOMES")
        base=children[0]
        result[arm]={"certified":base["certified"],"abstained":base["abstained"],
            "hard_false_positives":base["hard_false_positives"],
            "per_case":base["rows"],
            "learned_temporal_entries":base["TRAIN_learned_temporal_entries"],
            "public_TRAIN_examples":base["TRAIN_examples"],
            "initial_hypothesis_count":base["initial_global_hypotheses"],
            "five_full_process_cpu_ns":[x["full_process_cpu_ns"] for x in children],
            "five_peak_process_rss_bytes":[x["peak_process_rss_bytes"] for x in children],
            "median_full_process_cpu_ns":int(statistics.median(x["full_process_cpu_ns"] for x in children)),
            "max_peak_process_rss_bytes":max(x["peak_process_rss_bytes"] for x in children)}
    if not all(x["certified"]+x["abstained"]==192 and x["hard_false_positives"]==0
               for x in result.values()):
        raise RuntimeError("T26_COMPARISON_SCIENTIFIC_INVALID")
    learned=result[ARMS[0]]
    comparison={k:{"correct_certification_coverage_nondecreasing":learned["certified"]>=v["certified"],
       "complete_CPU_nondecreasing":learned["median_full_process_cpu_ns"]<=v["median_full_process_cpu_ns"],
       "peak_RSS_nondecreasing":learned["max_peak_process_rss_bytes"]<=v["max_peak_process_rss_bytes"]}
       for k,v in result.items() if k!=ARMS[0]}
    return {"schema":"wingless.T26B.original-sequence-world-model-comparison.v1",
        "arms":result,"comparison":comparison,"TRAIN_worlds":40,"prospective_DEV_worlds":64,
        "prospective_target_questions":192,"new_source_families":"U_and_V",
        "source_origin":"INTERNAL_NOT_EXTERNAL",
        "learned_cognition_positive":False,"joint_scientific_acceptance":False,
        "independent_external_replication":False,
        "optimizer_admission":False,"old_DEV_replayed":False,
        "no_paid_models":True,"no_controller_or_research_runner_modification":True}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(train_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T26_INVALID_ARGS")
