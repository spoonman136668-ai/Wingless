"""T29 scalable exact causal separator learning vs fully trained optimized controls.

Actual new F/G transition grammars and different 7/9 entity counts. Hidden
graph-mode assignments and state support persist through 5 masked frames.
Learn which intervention-target conditional separators are useful from
PUBLIC TRAIN ambiguity statistics; prove the projected state sufficient
for EVERY possible hidden state and EVERY public structural graph/mode.
No source-private labels enter fit or inference. Candidate is exact or
abstains. Full-state bitset and deterministic ancestor baselines share
identical TRAIN data and use no weaker controls.
"""
from __future__ import annotations
import importlib.util,json,statistics,subprocess,sys,time
from collections import Counter,defaultdict
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-learned-certified-causal-separator",
      "TRAIN-informed-frequency-projection",
      "fully-trained-exact-transition-bitset",
      "deterministic-minimal-ancestor-projection",
      "untrained-exact-enumeration")
def source():
    spec=importlib.util.spec_from_file_location("t29_src",ROOT/"20261019-wingless-t29-new-structural-worlds.py")
    m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m
    spec.loader.exec_module(m);return m
def rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Mem(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
              ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
              ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
              ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
              ("PagefileUsage",ctypes.c_size_t),
              ("PeakPagefileUsage",ctypes.c_size_t)]
        k=ctypes.WinDLL("kernel32",use_last_error=True)
        p=ctypes.WinDLL("psapi",use_last_error=True)
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T29_RSS_UNAVAILABLE_WIN32_ERROR="+str(ctypes.get_last_error()))
        return int(m.PeakWorkingSetSize)
    import resource
    a=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(a)*1024 if sys.platform.startswith("linux") else int(a)
def each_bit(value):
    while value:
        l=value&-value
        yield l.bit_length()-1
        value^=l
def filter_bits(mask,seen,n):
    if mask.bit_count() not in (2,3) or seen&~mask:
        raise RuntimeError("T29_SENSOR_CONTRACT_FAILED")
    return sum(1<<state for state in range(1<<n) if state&mask==seen)
def transition_map(src,family,t,g,action,cache):
    key=(family,t,g,action)
    if key not in cache:
        par=src.TEMPLATES[family][t]
        cache[key]=tuple(src.next_state(family,state,par,g,action)
                         for state in range(1<<src.NODES[family]))
    return cache[key]
def apply_map(support,mapping):
    result=0
    for state in each_bit(support):
        result|=1<<mapping[state]
    return result
def posterior(src,pub,cache,enumerative=False):
    f=pub["family"];n=src.NODES[f];frames=pub["observations"]
    if len(frames)!=5 or len(pub["history_actions"])!=4:
        raise RuntimeError("T29_WRONG_HISTORY_HORIZON")
    k=filter_bits(*frames[0],n)
    belief=[k]*16
    for i,action in enumerate(pub["history_actions"]):
        observed=filter_bits(*frames[i+1],n)
        updated=[]
        for label,current in enumerate(belief):
            if not current:updated.append(0);continue
            t,g=divmod(label,2)
            mapping=transition_map(src,f,t,g,action,cache)
            nxt=apply_map(current,mapping)&observed
            updated.append(nxt)
        belief=updated
    if not any(belief):raise RuntimeError("T29_TRUE_WORLD_FAMILY_ELIMINATED")
    return tuple(belief)
def future_outputs(src,f,t,g,action,target,cache):
    key=("future",f,t,g,action,target)
    if key not in cache:
        n=src.NODES[f]
        par=src.TEMPLATES[f][t]
        first=transition_map(src,f,t,g,action,cache)
        second=transition_map(src,f,t,g,None,cache)
        cache[key]=tuple((second[first[state]]>>target)&1 for state in range(1<<n))
    return cache[key]
def has_any_disagreement(belief,outmap):
    a=False;b=False
    for state in each_bit(belief):
        if outmap[state]:b=True
        else:a=True
        if a and b:return True
    return False
def standard_answer(src,pub,belief,cache):
    family=pub["family"];answers=[]
    for target in pub["targets"]:
        possible=set()
        for label,mask in enumerate(belief):
            if not mask:continue
            t,g=divmod(label,2)
            seq=future_outputs(src,family,t,g,pub["future_intervention"],target,cache)
            if has_any_disagreement(mask,seq):
                possible.update((0,1))
            else:
                first=(mask&-mask).bit_length()-1
                possible.add(seq[first])
            if len(possible)>1:break
        if not possible:raise RuntimeError("T29_NO_CAUSAL_ANSWER")
        answers.append(next(iter(possible)) if len(possible)==1 else None)
    return tuple(answers)

def project(value,coordinates):
    out=0
    for i,bit in enumerate(coordinates):
        if(value>>bit)&1:out|=1<<i
    return out
def exact_projection_map(output,coordinates):
    """Independent exhaustive finite state sufficiency proof."""
    values={}
    for state,val in enumerate(output):
        key=project(state,coordinates)
        if key in values and values[key]!=val:return None
        values[key]=val
    return values
def direct_ancestors(src,f,t,action,target):
    parents=src.TEMPLATES[f][t]
    one={target,*parents[target]}
    need=set()
    for node in one:
        if action[0]==node:continue
        need.add(node);need.update(parents[node])
    return tuple(sorted(need))
def certify_projection(src,f,t,g,action,target,order,cache):
    output=future_outputs(src,f,t,g,action,target,cache)
    # Dependency cone is sound by graph induction; greedily remove parts
    # that an independent exhaustive Boolean proof shows irrelevant.
    active=list(direct_ancestors(src,f,t,action,target))
    if exact_projection_map(output,active) is None:
        raise RuntimeError("T29_GRAPH_INFLUENCE_CONE_UNSOUND")
    for bit in order:
        if bit not in active:continue
        trial=[b for b in active if b!=bit]
        if exact_projection_map(output,trial) is not None:
            active=trial
    proof=exact_projection_map(output,active)
    if proof is None:raise RuntimeError("T29_CERTIFICATE_FAILED")
    return tuple(active),proof
def compile_macro(src,f,action,target,order,cache):
    configs=[]
    for t in range(8):
        for g in range(2):
            coords,table=certify_projection(src,f,t,g,action,target,order,cache)
            configs.append((coords,table))
    return tuple(configs)
def macro_answer(mask,macro):
    coords,table=macro
    options=set()
    for state in each_bit(mask):
        code=project(state,coords)
        if code not in table:raise RuntimeError("T29_MISSING_CERTIFIED_CAUSAL_CODE")
        options.add(table[code])
        if len(options)>1:break
    return options
def prediction(src,pub,belief,macro_store,cache):
    if macro_store is None:
        return standard_answer(src,pub,belief,cache),0,0
    f=pub["family"];action=pub["future_intervention"]
    outputs=[];hits=0;fallback=0
    for target in pub["targets"]:
        key=(f,action,target)
        macro=macro_store.get(key)
        if macro is None:
            # Full exact map, trained candidate has no right to guess.
            fallback+=1
            value=standard_single(src,pub,belief,target,cache)
        else:
            possible=set();hits+=1
            for label,support in enumerate(belief):
                if not support:continue
                possible.update(macro_answer(support,macro[label]))
                if len(possible)>1:break
            value=next(iter(possible)) if len(possible)==1 else None
        outputs.append(value)
    return tuple(outputs),hits,fallback
def standard_single(src,pub,belief,target,cache):
    fake=dict(pub);fake["targets"]=(target,)
    return standard_answer(src,fake,belief,cache)[0]
def learn(src,training,cache,arm):
    if arm not in ARMS[:2] and arm!=ARMS[3]:return {},{},0
    uses=Counter();ambiguities=Counter();features=Counter()
    for public in training:
        f=public["family"];q=public["future_intervention"]
        belief=posterior(src,public,cache)
        for t in public["targets"]:
            key=(f,q,t);uses[key]+=1
            res=standard_single(src,public,belief,t,cache)
            if res is None:ambiguities[key]+=1
            for mask,seen in public["observations"]:
                for b in range(src.NODES[f]):
                    if mask>>b&1:features[f,b]+=1
    chosen={}
    for family in ("F","G"):
        local=[key for key in uses if key[0]==family]
        if arm==ARMS[0]:
            # Public TRAIN disagreement teaches where proof representation
            # can amortize the most nontrivial counterfactual work.
            ordering=sorted(local,key=lambda k:(-(uses[k]+2*ambiguities[k]),k))
        elif arm==ARMS[1]:
            ordering=sorted(local,key=lambda k:(-uses[k],k))
        else:
            ordering=sorted(local)
        for key in ordering[:16]:chosen[key]=None
    macros={}
    proof_count=0
    for f,action,target in chosen:
        n=src.NODES[f]
        # TRAIN-learned salience affects elimination trajectory, and is
        # never based on held-out source/latent hidden truth.
        if arm==ARMS[0]:
            order=tuple(sorted(range(n),key=lambda i:(features[f,i],i)))
        elif arm==ARMS[1]:
            order=tuple(sorted(range(n),reverse=True))
        else:
            order=tuple(range(n))
        compiled=compile_macro(src,f,action,target,order,cache)
        macros[f,action,target]=compiled
        proof_count+=sum(len(v[0])<n for v in compiled)
    return macros,{"compiled_query_types":len(macros),
        "universally_certified_nontrivial_projection_count":proof_count,
        "TRAIN_seen_future_query_types":len(uses)},proof_count

def TRAIN_only():
    src=source();cache={}
    if len(src.TRAIN)!=40 or len(src.DEV)!=64:
        raise RuntimeError("T29_ORIGINAL_COHORT_IDENTITY_MISMATCH")
    training=[src.train_public(src.sample(s)) for s in src.TRAIN]
    learner,stats,nontrivial=learn(src,training,cache,ARMS[0])
    frequency,_,_=learn(src,training,cache,ARMS[1])
    exact=0;abstained=0;retained=0;comparisons=0;hits=0
    for name,pub in zip(src.TRAIN,training):
        item=src.sample(name)
        belief=posterior(src,pub,cache)
        label=item["_hidden_template"]*2+item["_hidden_mode"]
        if not(belief[label]&(1<<item["_hidden_current_state"])):
            raise RuntimeError("T29_TRUE_NEW_GRAMMAR_WORLD_EXCLUDED")
        retained+=1
        ref=standard_answer(src,pub,belief,cache)
        cand,h,_=prediction(src,pub,belief,learner,cache)
        freq,*_=prediction(src,pub,belief,frequency,cache)
        if cand!=ref or freq!=ref:
            raise RuntimeError("T29_LEARNED_CAUSAL_SEPARATOR_NOT_LOSSLESS")
        hits+=h
        for v,truth in zip(cand,item["_private_future_answers"]):
            comparisons+=1
            if v is None:abstained+=1
            else:
                if v!=truth:raise RuntimeError("T29_TRAIN_HARD_FALSE_POSITIVE")
                exact+=1
    return {"schema":"wingless.T29.original-TRAIN-only-structural-sufficiency.v1",
        "status":"PASS_TRAIN_ONLY_DEV_UNOPENED",
        "source_families":"F_seven_entities_and_G_nine_entities",
        "TRAIN_cases":40,"DEV_generated":0,"TRUE_hidden_worlds_retained":retained,
        "counterfactual_queries":comparisons,"certified":exact,"abstained":abstained,
        "ambiguity_fraction":abstained/comparisons,
        "exact_learned_equals_strong_control_queries":comparisons,
        "learned_training_macro_hits":hits,"universally_proven_nontrivial_separators":nontrivial,
        "source_structural_novelty":"DIFFERENT_ENTITY_COUNT_PARENT_ARITY_AND_LOGIC",
        "training_from_public_only":True,
        "scientific_positive":False,"optimizer_admission":False,
        "old_DEV_replayed":False,**stats}
def trial(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T29_BAD_ARM")
    begin=time.process_time_ns()
    src=source();cache={}
    TRAIN=[src.train_public(src.sample(s)) for s in src.TRAIN]
    macros,fitstats,proofs=learn(src,TRAIN,cache,arm)
    results=[];hits=0;fallback=0;certified=0;abstained=0
    for name in src.DEV:
        item=src.sample(name)
        pub=src.dev_public(item)
        belief=posterior(src,pub,cache)
        val,h,f=prediction(src,pub,belief,macros if arm in ARMS[:2] or arm==ARMS[3] else None,cache)
        yes=0;no=0
        for answer,truth in zip(val,item["_private_future_answers"]):
            if answer is None:no+=1;abstained+=1
            else:
                if answer!=truth:raise RuntimeError("T29_ORIGINAL_HARD_FALSE_POSITIVE")
                yes+=1;certified+=1
        hits+=h;fallback+=f
        results.append({"id":name,"correct_certified":yes,"sound_abstentions":no,
          "source_family":pub["family"],
          "posterior_hidden_world_count":sum(b.bit_count() for b in belief),
          "learned_macro_hits":h,"exact_fallbacks":f})
    if certified+abstained!=192 or len(results)!=64:
        raise RuntimeError("T29_ORIGINAL_MISSING_CASE")
    return {"schema":"wingless.T29.original-finite-structural-learner-process.v1",
        "arm":arm,"rep":rep,"TRAIN_cases":40,"DEV_cases":64,
        "certified":certified,"sound_abstentions":abstained,"hard_false_positives":0,
        "macro_hits":hits,"full_exact_fallbacks":fallback,
        "source_bound_case_results":results,"TRAIN_proved_separators":proofs,
        "whole_process_CPU_ns":time.process_time_ns()-begin,
        "peak_process_RSS_bytes":rss(),
        "external_independent_replication":False,
        "old_DEV_replayed":False,"optimizer_admission":False}
def evaluate():
    results={}
    for arm in ARMS:
        trials=[]
        for i in range(5):
            q=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                   "--arm",arm,str(i)],capture_output=True,text=True)
            if q.returncode:raise RuntimeError("T29_ORIGINAL_SUBPROCESS_FAILED "+arm+" "+str(i)+"\n"+q.stderr)
            trials.append(json.loads(q.stdout))
        if any(t["source_bound_case_results"]!=trials[0]["source_bound_case_results"] for t in trials):
            raise RuntimeError("T29_NONDETERMINISM")
        v=trials[0]
        results[arm]={"correct_certified":v["certified"],
          "sound_abstentions":v["sound_abstentions"],"hard_false_positives":v["hard_false_positives"],
          "per_case":v["source_bound_case_results"],"TRAIN_proved_separators":v["TRAIN_proved_separators"],
          "original_learned_macro_hits":v["macro_hits"],"full_exact_fallbacks":v["full_exact_fallbacks"],
          "five_complete_CPU_ns":[t["whole_process_CPU_ns"] for t in trials],
          "five_peak_RSS_bytes":[t["peak_process_RSS_bytes"] for t in trials],
          "median_complete_CPU_ns":int(statistics.median(t["whole_process_CPU_ns"] for t in trials)),
          "max_peak_RSS_bytes":max(t["peak_process_RSS_bytes"] for t in trials)}
    if any(v["correct_certified"]+v["sound_abstentions"]!=192 or v["hard_false_positives"] for v in results.values()):
        raise RuntimeError("T29_ORIGINAL_SAFETY_CRITERIA_FAILED")
    cand=results[ARMS[0]]
    comp={k:{"no_baseline_answer_loss":cand["correct_certified"]>=v["correct_certified"],
       "complete_CPU_nonregression":cand["median_complete_CPU_ns"]<=v["median_complete_CPU_ns"],
       "peak_RSS_nonregression":cand["max_peak_RSS_bytes"]<=v["max_peak_RSS_bytes"]}
       for k,v in results.items() if k!=ARMS[0]}
    return {"schema":"wingless.T29.original-structural-generalization-comparison.v1",
      "source_families":"F_and_G_NEW_INTERNAL_GRAMMARS",
      "TRAIN_cases":40,"DEV_cases":64,"counterfactual_questions":192,
      "arms":results,"comparison":comp,
      "independent_external_replication":False,"scientific_positive":False,
      "optimizer_admission":False,"old_DEV_replayed":False,
      "no_paid_models":True,"no_new_controller":True}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(TRAIN_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T29_BAD_ARGUMENTS")
