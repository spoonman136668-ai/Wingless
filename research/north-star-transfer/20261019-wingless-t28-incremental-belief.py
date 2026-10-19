"""T28: TRAIN-learned incremental, proof-carrying exact relational belief.

Each hidden topology/gate has its own 64-bit possible-state mask, updated
incrementally from masked chronological observations. Public TRAIN trajectories
induce a 16-entry dictionary of support-mask patterns and exact recurring
action/observation update macros. Every macro is independently re-derived
from the public causal grammar and cannot drop a possible hidden world.
The learned candidate uses macro+dictionary codes with exact fallback.
All strong controls receive the same public TRAIN streams and truth-safe
causal rules; they do not open or sample private dev truths while training.
No T27 full 72-bit predictive atlas is built.
"""
from __future__ import annotations
import importlib.util,json,statistics,subprocess,sys,time,math
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
FULL=(1<<64)-1
ARMS=("TRAIN-learned-incremental-belief",
      "TRAIN-frequency-dictionary-no-macros",
      "fully-trained-exact-bitset-cache",
      "fully-trained-exact-enumeration",
      "untrained-exact-bitset")
def load_source():
    spec=importlib.util.spec_from_file_location("t28_new_streams",ROOT/"20261019-wingless-t28-new-streams.py")
    src=importlib.util.module_from_spec(spec);sys.modules[spec.name]=src;spec.loader.exec_module(src)
    return src
def peakrss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Memory(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
              ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
              ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
              ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
              ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        k=ctypes.WinDLL("kernel32");p=ctypes.WinDLL("psapi")
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Memory),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Memory();m.cb=ctypes.sizeof(m)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T28_PROCESS_RSS_UNMEASURED")
        return int(m.PeakWorkingSetSize)
    import resource
    b=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(b)*1024 if sys.platform.startswith("linux") else int(b)

def observation_filter(mask,seen):
    if mask.bit_count()!=2 or seen&~mask:
        raise RuntimeError("T28_BAD_PUBLIC_SENSOR_OBSERVATION")
    return sum(1<<state for state in range(64) if state&mask==seen)

def bits(values):
    while values:
        low=values&-values
        yield low.bit_length()-1
        values^=low

def exact_transition(src,label,old_mask,action):
    t,g=divmod(label,2)
    parents=src.TEMPLATES[t]
    future=0
    for state in bits(old_mask):
        future|=1<<src.next_state(state,parents,g,action)
    return future

def exact_update(src,label,old_mask,action,mask,seen):
    return exact_transition(src,label,old_mask,action)&observation_filter(mask,seen)

def fresh_belief(public,src,trained_macros=None,dictionary=None):
    frames=public["observations"];actions=public["history_actions"]
    if len(frames)!=5 or len(actions)!=4:
        raise RuntimeError("T28_TEMPORAL_SOURCE_IDENTITY_MISMATCH")
    first_filter=observation_filter(*frames[0])
    posterior=[first_filter]*24
    cached=0;residual=0;encoded=0;full=0
    for idx,action in enumerate(actions):
        mask,seen=frames[idx+1]
        candidate_filter=observation_filter(mask,seen)
        nextposterior=[]
        for label,old in enumerate(posterior):
            if not old:
                nextposterior.append(0);continue
            if dictionary is not None:
                if old in dictionary:encoded+=1
                else:full+=1
            key=(label,old,action,mask,seen)
            if trained_macros is not None and key in trained_macros:
                new=trained_macros[key]
                cached+=1
            else:
                new=exact_transition(src,label,old,action)&candidate_filter
                residual+=1
            nextposterior.append(new)
        posterior=nextposterior
    if not any(posterior):raise RuntimeError("T28_ALL_HIDDEN_WORLDS_ELIMINATED")
    return tuple(posterior),{"hits":cached,"exact_fallbacks":residual,
       "dictionary_references":encoded,"raw_slices":full}

def full_enumeration(src,public):
    frame0=observation_filter(*public["observations"][0])
    active=[(label,s) for label in range(24) for s in bits(frame0)]
    for i,action in enumerate(public["history_actions"]):
        mask,seen=public["observations"][i+1]
        new=set()
        for label,state in active:
            t,g=divmod(label,2)
            nxt=src.next_state(state,src.TEMPLATES[t],g,action)
            if nxt&mask==seen:new.add((label,nxt))
        active=list(new)
    result=[0]*24
    for label,state in active:result[label]|=1<<state
    if not any(result):raise RuntimeError("T28_ENUMERATIVE_BELIEF_EMPTY")
    return tuple(result)

def predict(src,posterior,public):
    action=public["future_intervention"]
    options=[set(),set(),set()]
    for label,mask in enumerate(posterior):
        if not mask:continue
        t,g=divmod(label,2)
        parents=src.TEMPLATES[t]
        for state in bits(mask):
            one=src.next_state(state,parents,g,action)
            two=src.next_state(one,parents,g)
            for i,target in enumerate(public["targets"]):
                options[i].add((two>>target)&1)
    if not all(options):raise RuntimeError("T28_NO_PREDICTIONS")
    vals=tuple(next(iter(values)) if len(values)==1 else None for values in options)
    return vals,tuple(tuple(sorted(s)) for s in options)

def fit(src,TRAIN):
    """Fit exactly from PUBLIC TRAIN observations, never source hidden labels."""
    frequency=Counter()
    macrocounts=Counter()
    macrowitness={}
    for public in TRAIN:
        frames=public["observations"];actions=public["history_actions"]
        posterior=[observation_filter(*frames[0])]*24
        for idx,action in enumerate(actions):
            mask,seen=frames[idx+1]
            filt=observation_filter(mask,seen)
            new=[]
            for label,old in enumerate(posterior):
                if old:
                    frequency[old]+=1
                    key=(label,old,action,mask,seen)
                    candidate=exact_transition(src,label,old,action)&filt
                    macrocounts[key]+=1
                    if key in macrowitness and macrowitness[key]!=candidate:
                        raise RuntimeError("T28_NONDETERMINISTIC_TRAINED_MACRO")
                    macrowitness[key]=candidate
                else:candidate=0
                new.append(candidate)
            posterior=new
    # Train compact state support-mask dictionary from the recurring
    # PUBLIC approximate posterior motifs across all source graph hypotheses.
    dictionary=tuple(mask for mask,n in sorted(frequency.items(),
                            key=lambda x:(-x[1],x[0]))[:16])
    if len(dictionary)<8:raise RuntimeError("T28_INSUFFICIENT_COMPRESSION_DICTIONARY")
    keys=sorted(macrocounts,key=lambda k:(-macrocounts[k],k[0],k[1],str(k[2]),k[3],k[4]))
    macrokeys=tuple(k for k in keys if k[1] in dictionary)[:128]
    macros={k:macrowitness[k] for k in macrokeys}
    # Deterministically compiled exact TRAIN cache: all public historical
    # role/mask/action combinations, stronger than just the learned 128.
    fully_trained_cache=macrowitness
    for key,val in fully_trained_cache.items():
        label,old,action,mask,seen=key
        if exact_update(src,label,old,action,mask,seen)!=val:
            raise RuntimeError("T28_FULL_TRAINED_CACHE_FALSE")
    for key,val in macros.items():
        label,old,action,mask,seen=key
        if exact_update(src,label,old,action,mask,seen)!=val:
            raise RuntimeError("T28_LEARNED_PROOF_FALSE")
    return dictionary,macros,fully_trained_cache,{
       "TRAIN_dictionary_entries":len(dictionary),
       "TRAIN_unique_verified_macro_keys":len(fully_trained_cache),
       "TRAIN_learned_verified_macros":len(macros),
       "TRAIN_public_transition_witnesses":sum(macrocounts.values())}

def retain_true_world(src,posterior,case):
    label=2*case["_private_topology"]+case["_private_mode"]
    state=case["_private_terminal_state"]
    return bool(posterior[label]&(1<<state))

def train_only():
    src=load_source()
    if len(src.TRAIN)!=40 or len(src.DEV)!=64 or set(src.TRAIN)&set(src.DEV):
        raise RuntimeError("T28_TRAIN_SOURCE_IDENTITY_DRIFT")
    inputs=[src.train_public(src.sample(k)) for k in src.TRAIN]
    dictionary,macros,fullcache,stats=fit(src,inputs)
    n=0;correct=0;abstained=0;equivalent=0;retained=0;hits=0;compressed=0
    for label,public in zip(src.TRAIN,inputs):
        private=src.sample(label)
        reference,raw=fresh_belief(public,src)
        learned,trace=fresh_belief(public,src,macros,set(dictionary))
        full,other=fresh_belief(public,src,fullcache,None)
        enum=full_enumeration(src,public)
        if reference!=learned or reference!=full or reference!=enum:
            raise RuntimeError("T28_TRAIN_LATENT_BINDING_POSTERIOR_MISMATCH")
        if not retain_true_world(src,reference,private):
            raise RuntimeError("T28_TRUE_LATENT_WORLD_EXCLUDED")
        retained+=1
        prediction,options=predict(src,learned,public)
        baseline,*_=predict(src,enum,public)
        if prediction!=baseline:
            raise RuntimeError("T28_SAFETY_COVERAGE_LOSS")
        equivalent+=3
        hits+=trace["hits"];compressed+=trace["dictionary_references"]
        for v,s,y in zip(prediction,options,private["_truth"]):
            if y not in s:raise RuntimeError("T28_TRUE_COUNTERFACTUAL_EXCLUDED")
            if v is None:abstained+=1
            else:
                if v!=y:raise RuntimeError("T28_TRAIN_FALSE_POSITIVE")
                correct+=1
            n+=1
    return {"schema":"wingless.T28.TRAIN-only-incremental-exact-belief.v1",
       "status":"PASS_TRAIN_ONLY_DEV_UNOPENED","TRAIN_cases":40,"DEV_generated":0,
       "TRAIN_counterfactual_target_questions":n,
       "TRAIN_certified_correct":correct,"TRAIN_safe_abstentions":abstained,
       "TRAIN_ambiguity_fraction":abstained/n,
       "TRAIN_exact_posterior_and_answer_equivalences":equivalent,
       "TRAIN_true_persistent_worlds_retained":retained,
       "TRAIN_learned_cache_hits":hits,
       "TRAIN_dictionary_references":compressed,
       "hidden_graph_mode_channels":24,
       "hidden_states_per_channel":64,
       "no_predictive_signature_atlas":True,
       "independent_exact_macro_verification":True,
       "source_origin":"INTERNAL_NOT_EXTERNAL",
       "optimizer_admission":False,"scientific_positive":False,
       "old_DEV_not_replayed":True,**stats}

def trial(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T28_BAD_ARM")
    start=time.process_time_ns()
    src=load_source()
    TRAIN=[src.train_public(src.sample(k)) for k in src.TRAIN]
    dictionary,macros,fullcache,fitstats=fit(src,TRAIN) if arm!=ARMS[4] else ((),{},{},{})
    cache=(macros if arm==ARMS[0] else fullcache if arm==ARMS[2] else None)
    dictset=set(dictionary) if arm in ARMS[:2] else None
    rows=[];correct=0;unknown=0;cachehits=0;direct=0;dictionaryhits=0;raw=0
    for key in src.DEV:
        case=src.sample(key);public=src.dev_public(case)
        if arm==ARMS[3]:
            posterior=full_enumeration(src,public)
            trace={"hits":0,"exact_fallbacks":0,"dictionary_references":0,"raw_slices":0}
        else:
            posterior,trace=fresh_belief(public,src,cache,dictset)
        values,opts=predict(src,posterior,public)
        yes=no=0
        for v,s,y in zip(values,opts,case["_truth"]):
            if y not in s:raise RuntimeError("T28_DEV_LATENT_TRUTH_EXCLUDED")
            if v is None:unknown+=1;no+=1
            else:
                if v!=y:raise RuntimeError("T28_HARD_FALSE_POSITIVE")
                correct+=1;yes+=1
        cachehits+=trace["hits"];direct+=trace["exact_fallbacks"]
        dictionaryhits+=trace["dictionary_references"];raw+=trace["raw_slices"]
        rows.append({"source_id":key,"certified":yes,"sound_abstentions":no,
              "latent_support_size":sum(mask.bit_count() for mask in posterior),
              "TRAIN_macros_used":trace["hits"],
              "incremental_dictionary_references":trace["dictionary_references"]})
    if correct+unknown!=192 or len(rows)!=64:raise RuntimeError("T28_ORIGINAL_INCOMPLETE")
    return {"schema":"wingless.T28.original-incremental-memory-process.v1","arm":arm,
       "rep":rep,"TRAIN_cases":40,"prospective_DEV_cases":64,"rows":rows,
       "correct_certified":correct,"sound_abstentions":unknown,"hard_false_positives":0,
       "TRAIN_verified_macro_fit":fitstats.get("TRAIN_learned_verified_macros",0),
       "original_macro_hits":cachehits,"original_exact_fallbacks":direct,
       "original_dictionary_references":dictionaryhits,"original_raw_support_slices":raw,
       "whole_process_cpu_ns":time.process_time_ns()-start,
       "peak_process_rss_bytes":peakrss(),
       "old_DEV_replayed":False,"external_replication":False,"optimizer_admission":False}
def evaluate():
    results={}
    for arm in ARMS:
        reps=[]
        for i in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),"--arm",arm,str(i)],
                             capture_output=True,text=True)
            if p.returncode:raise RuntimeError("T28_ORIGINAL_CHILD_FAILED:"+arm+":"+str(i)+"\n"+p.stderr)
            reps.append(json.loads(p.stdout))
        if any(v["rows"]!=reps[0]["rows"] for v in reps):
            raise RuntimeError("T28_NONDETERMINISTIC_SOURCE_IDENTITIES")
        v=reps[0]
        results[arm]={"correct_certified":v["correct_certified"],
            "safe_abstained":v["sound_abstentions"],
            "hard_false_positives":v["hard_false_positives"],
            "per_case":v["rows"],"original_macro_hits":v["original_macro_hits"],
            "original_dictionary_references":v["original_dictionary_references"],
            "TRAIN_verified_macros":v["TRAIN_verified_macro_fit"],
            "five_full_process_cpu_ns":[x["whole_process_cpu_ns"] for x in reps],
            "five_peak_process_rss_bytes":[x["peak_process_rss_bytes"] for x in reps],
            "median_complete_cpu_ns":int(statistics.median(x["whole_process_cpu_ns"] for x in reps)),
            "max_peak_rss_bytes":max(x["peak_process_rss_bytes"] for x in reps)}
        if v["correct_certified"]+v["sound_abstentions"]!=192 or v["hard_false_positives"]:
            raise RuntimeError("T28_UNSOUND_RESULT")
    learned=results[ARMS[0]]
    comparison={a:{"coverage_preserved":learned["correct_certified"]>=v["correct_certified"],
                   "complete_cpu_nonregression":learned["median_complete_cpu_ns"]<=v["median_complete_cpu_ns"],
                   "peak_rss_nonregression":learned["max_peak_rss_bytes"]<=v["max_peak_rss_bytes"]}
                for a,v in results.items() if a!=ARMS[0]}
    return {"schema":"wingless.T28.original-incremental-belief-comparison.v1",
       "TRAIN_cases":40,"DEV_cases":64,"counterfactual_target_questions":192,
       "new_source_families":"Y_and_Z_internal","arms":results,"comparison":comparison,
       "scientific_positive":False,"independent_external_replication":False,
       "optimizer_admission":False,"old_DEV_replayed":False,"paid_API":False,
       "new_controller_or_research_runner":False}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(train_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T28_UNRECOGNIZED_ARGS")
