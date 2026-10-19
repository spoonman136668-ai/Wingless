"""T27 proof-carrying TRAIN-learned predictive-state sufficient statistic.

Learns (TRAIN public only) a <=16-feature set of future intervention-target
coordinates (12 interventions x 6 entity outputs). Exact independently
verified global atlas proves when a feature code determines a full 72-bit
counterfactual-response signature. Ambiguous feature codes fall back to
full exact signature computation; no hidden label or guessed answer.

The proof lives over ALL 12*2*64 template/mode/terminal-state worlds,
not just sampled positive cases. The posterior retains true stable world
bindings through chronological masked observations. At query time, a set
of unique full predictive signatures is a lossless sufficient statistic
for every prospective two-step intervention/counterfactual in the
preregistered query family.
"""
from __future__ import annotations
import importlib.util,itertools,json,statistics,subprocess,sys,time
from collections import defaultdict,Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
FEATURE_CAP=16
ARMS=("TRAIN-learned-predictive-state-quotient",
      "TRAIN-frequency-predictive-quotient-control",
      "deterministic-full-signature-quotient",
      "fully-TRAIN-exact-recurrent-Bayesian",
      "untrained-unpruned-exact-recurrent")
def source():
    name="t27_frozen_source"
    spec=importlib.util.spec_from_file_location(name,ROOT/"20261019-wingless-t27-new-predictive-worlds.py")
    m=importlib.util.module_from_spec(spec);sys.modules[name]=m;spec.loader.exec_module(m)
    return m
def rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class M(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
             ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
             ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
             ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
             ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        k=ctypes.WinDLL("kernel32");p=ctypes.WinDLL("psapi")
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(M),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        z=M();z.cb=ctypes.sizeof(z)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(z),z.cb):
            raise RuntimeError("T27_RSS_UNAVAILABLE")
        return int(z.PeakWorkingSetSize)
    import resource
    u=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(u)*1024 if sys.platform.startswith("linux") else int(u)
def worlds(public,src,prior=(0,0)):
    frames=public["observations"];actions=public["history_actions"]
    if len(frames)!=5 or len(actions)!=4 or any(mask.bit_count()!=2 or obs&~mask for mask,obs in frames):
        raise RuntimeError("T27_NOT_SPARSE_FIVE_FRAME_SOURCE")
    surviving=[]
    gate_order=sorted((0,1),key=lambda g:(-prior[g],g))
    for gate in gate_order:
        for t,parents in enumerate(src.TEMPLATES):
            for start in range(64):
                mask,obs=frames[0]
                if start&mask!=obs:continue
                state=start;ok=True
                for step,action in enumerate(actions):
                    state=src.next_state(state,parents,gate,action)
                    mask,obs=frames[step+1]
                    if state&mask!=obs:
                        ok=False;break
                if ok:surviving.append((t,gate,state))
    if not surviving:raise RuntimeError("T27_EMPTY_POSTERIOR")
    return tuple(surviving)
def action_at(idx):
    return (idx//2,idx%2)
def response6(src,t,g,state,action):
    nxt=src.next_state(state,src.TEMPLATES[t],g,action)
    return src.next_state(nxt,src.TEMPLATES[t],g)
def full_signature(src,world):
    t,g,state=world;bits=0
    for act in range(12):
        final=response6(src,t,g,state,action_at(act))
        bits|=final<<(6*act)
    return bits
def selected_code(src,world,coordinates):
    """Infer at most one 6-bit response per selected public action."""
    t,g,state=world
    outputs={}
    code=0
    for j,coord in enumerate(coordinates):
        action=coord//6
        if action not in outputs:
            outputs[action]=response6(src,t,g,state,action_at(action))
        if (outputs[action]>>(coord%6))&1:code|=1<<j
    return code
def from_signature(signature,coordinates):
    out=0
    for j,coord in enumerate(coordinates):
        if (signature>>coord)&1:out|=1<<j
    return out
def all_signature_atlas(src):
    states={}
    for t in range(12):
        for g in range(2):
            for state in range(64):
                w=(t,g,state)
                states[w]=full_signature(src,w)
    if len(states)!=1536:raise RuntimeError("T27_ATLAS_INCOMPLETE")
    return states
def choose_features(TRAIN_signatures,mode):
    """Induce learned predictive coordinates from PUBLIC TRAIN posteriors."""
    assert mode in ARMS[:2]
    # Each TRAIN world contributes a distribution over distinct available
    # predictive signatures; no ground-truth hidden graph enters selection.
    counts=Counter(x for case in TRAIN_signatures for x in set(case))
    chosen=[]
    for n in range(FEATURE_CAP):
        best=None
        for coord in range(72):
            if coord in chosen:continue
            if mode==ARMS[1]:
                ones=sum(weight for sig,weight in counts.items() if sig>>coord&1)
                zeros=sum(counts.values())-ones
                score=min(ones,zeros)
            else:
                # Partition reduction weights actual TRAIN candidate sets,
                # never the hidden true member; learned role transfer.
                score=0
                for sigs in TRAIN_signatures:
                    groups=defaultdict(lambda:[0,0])
                    for sig in set(sigs):
                        key=from_signature(sig,chosen)
                        groups[key][(sig>>coord)&1]+=1
                    score+=sum(a*b for a,b in groups.values())
            candidate=(score,-coord)
            if best is None or candidate>best[0]:
                best=(candidate,coord)
        chosen.append(best[1])
    if len(set(chosen))!=16:raise RuntimeError("T27_LEARNED_BASIS_NOT_FROZEN")
    return tuple(chosen)
def certified_atlas_map(atlas,coordinates):
    """Independent exhaustive full-future verifier, not TRAIN sampled."""
    grouped={}
    ambiguous=set()
    for sig in atlas.values():
        code=from_signature(sig,coordinates)
        if code in grouped and grouped[code]!=sig:
            ambiguous.add(code)
        else:grouped[code]=sig
    for code in ambiguous:grouped.pop(code,None)
    # A surviving code may NEVER give an incorrect full-future response.
    for sig in atlas.values():
        code=from_signature(sig,coordinates)
        if code in grouped and grouped[code]!=sig:
            raise RuntimeError("T27_FEATURE_CERTIFICATE_FALSE")
    return grouped,ambiguous

def prepare_model(src,TRAIN,arm):
    PRIOR=[0,0]
    posteriors=[]
    for public in TRAIN:
        posterior=worlds(public,src)
        posteriors.append(posterior)
        for _,gate,_ in posterior:PRIOR[gate]+=1
    if arm in ARMS[:2]:
        atlas=all_signature_atlas(src)
        sigs=[tuple(atlas[w] for w in posterior) for posterior in posteriors]
        coords=choose_features(sigs,arm)
        map_,collisions=certified_atlas_map(atlas,coords)
        return tuple(PRIOR),coords,map_,len(collisions),len(set(atlas.values()))
    return tuple(PRIOR),(),{},0,0

def prediction(src,public,arm,prior,coords,certified_map):
    posterior=worlds(public,src,prior if arm in ARMS[:2] else (0,0))
    if arm in ARMS[:3]:
        signatures=set()
        good=0;fallback=0
        for w in posterior:
            if coords:
                code=selected_code(src,w,coords)
                if code in certified_map:
                    sig=certified_map[code];good+=1
                else:
                    sig=full_signature(src,w);fallback+=1
            else:
                sig=full_signature(src,w)
            signatures.add(sig)
        intervention=public["future_intervention"]
        action_index=2*intervention[0]+intervention[1]
        outputs=tuple(
            frozenset((sig>>(action_index*6+node))&1 for sig in signatures)
            for node in public["targets"])
        result=tuple(next(iter(s)) if len(s)==1 else None for s in outputs)
        return result,outputs,len(posterior),len(signatures),good,fallback
    # Strong original exact full-support recurrent baseline, no quotient.
    answers=[]
    for w in posterior:
        t,g,state=w
        final=response6(src,t,g,state,public["future_intervention"])
        answers.append(tuple((final>>node)&1 for node in public["targets"]))
    sets=tuple(frozenset(row[j] for row in answers) for j in range(3))
    values=tuple(next(iter(s)) if len(s)==1 else None for s in sets)
    return values,sets,len(posterior),len(set(answers)),0,0

def TRAIN_only():
    src=source()
    if len(src.TRAIN)!=40 or len(src.DEV)!=64 or set(src.TRAIN)&set(src.DEV):
        raise RuntimeError("T27_COHORT_DRIFT")
    TRAIN=[src.train_public(src.sample(k)) for k in src.TRAIN]
    prior,coords,cert,collision,unique=prepare_model(src,TRAIN,ARMS[0])
    refprior,refcoords,refcert,refcoll,refunique=prepare_model(src,TRAIN,ARMS[1])
    if len(coords)!=16 or len(refcoords)!=16:raise RuntimeError("T27_BASIS_NOT_FITTED")
    queried=0;ambiguity=0;survived=0;verified=0;exact=0;merged=0;hits=0;fallback=0
    for name,public in zip(src.TRAIN,TRAIN):
        case=src.sample(name)
        poster=worlds(public,src)
        # Exact graph/mode/last-state is preserved although initial state may
        # be forgotten as a Markov sufficient statistic.
        if (case["_private_topology"],case["_private_mode"],case["_private_terminal_state"]) not in poster:
            raise RuntimeError("T27_TRUE_LATENT_WORLD_EXCLUDED")
        survived+=1
        result,sets,states,classes,h,f=prediction(src,public,ARMS[0],prior,coords,cert)
        baseline,*_=prediction(src,public,ARMS[3],(0,0),(),{})
        freq,*_=prediction(src,public,ARMS[1],refprior,refcoords,refcert)
        if baseline!=result or baseline!=freq:raise RuntimeError("T27_TRAIN_LOSSLESS_EQUIVALENCE_FAILED")
        merged+=states-classes;hits+=h;fallback+=f
        for a,s,t in zip(result,sets,case["_truth"]):
            verified+=1
            if t not in s:raise RuntimeError("T27_FALSE_EXCLUSION")
            if a is None:ambiguity+=1
            else:
                if a!=t:raise RuntimeError("T27_FALSE_POSITIVE")
                exact+=1
            queried+=1
    return {"schema":"wingless.T27.TRAIN-only-predictive-sufficiency.v1",
        "status":"PASS_TRAIN_ONLY_DEV_UNOPENED","TRAIN_cases":40,"DEV_generated":0,
        "TRUE_hidden_worlds_retained":survived,
        "TRAIN_question_count":queried,"ambiguity_count":ambiguity,
        "ambiguity_fraction":ambiguity/queried,"correctly_certified":exact,
        "TRAIN_exact_answer_equivalence_checks":verified,
        "learned_feature_coordinates":coords,"frequency_feature_coordinates":refcoords,
        "learned_certified_global_code_count":len(cert),
        "learned_ambiguous_codes_requiring_exact_residual":collision,
        "full_global_predictive_signature_classes":unique,
        "TRAIN_latent_membership_merge_count":merged,
        "TRAIN_learned_verified_code_hits":hits,
        "TRAIN_full_signature_fallbacks":fallback,
        "global_hidden_state_certification_worlds":1536,
        "full_signature_bits_per_state":72,
        "no_source_private_labels_in_learning":True,
        "no_old_DEV_replay":True,"scientific_positive":False,
        "optimizer_admission":False}
def trial(arm,rep):
    assert arm in ARMS and rep in range(5)
    begin=time.process_time_ns()
    src=source()
    train=[src.train_public(src.sample(x)) for x in src.TRAIN]
    prior,coords,cert,collisions,unique=prepare_model(src,train,arm)
    rows=[];certified=0;abstained=0;hits=0;fallback=0;merged=0
    for name in src.DEV:
        case=src.sample(name);public=src.dev_public(case)
        answer,poss,states,classes,h,f=prediction(src,public,arm,prior,coords,cert)
        yes=0;unknown=0
        for a,s,true in zip(answer,poss,case["_truth"]):
            if true not in s:raise RuntimeError("T27_ORIGINAL_FALSE_EXCLUSION:"+name)
            if a is None:unknown+=1
            else:
                if a!=true:raise RuntimeError("T27_ORIGINAL_FALSE_POSITIVE:"+name)
                yes+=1
        certified+=yes;abstained+=unknown;hits+=h;fallback+=f;merged+=states-classes
        rows.append({"source_id":name,"certified":yes,"abstained":unknown,
           "latent_hypotheses":states,"predictive_classes":classes,
           "verified_learned_code_hits":h,"full_signature_fallbacks":f})
    if len(rows)!=64 or certified+abstained!=192:raise RuntimeError("T27_MISSING_ORIGINAL_CASE")
    return {"schema":"wingless.T27.original-finite-causal-belief-single-process.v1",
        "arm":arm,"rep":rep,"rows":rows,"certified":certified,"abstained":abstained,
        "hard_false_positives":0,"TRAIN_cases":40,"DEV_cases":64,
        "learned_feature_count":len(coords),"globally_certified_codes":len(cert),
        "TRAIN_code_collisions":collisions,"global_predictive_classes":unique,
        "original_posterior_merge_count":merged,
        "original_verified_feature_hits":hits,"original_exact_residual_fallbacks":fallback,
        "whole_process_CPU_ns":time.process_time_ns()-begin,
        "peak_process_RSS_bytes":rss(),
        "internal_origin":True,"independent_replication":False,
        "optimizer_admission":False,"old_DEV_replayed":False}
def evaluate():
    result={}
    for arm in ARMS:
        records=[]
        for rep in range(5):
            z=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                  "--arm",arm,str(rep)],capture_output=True,text=True)
            if z.returncode:
                raise RuntimeError("T27_ORIGINAL_CHILD_FAILED:"+arm+":"+str(rep)+"\n"+z.stderr)
            records.append(json.loads(z.stdout))
        if any(r["rows"]!=records[0]["rows"] for r in records):
            raise RuntimeError("T27_NONDETERMINISM")
        a=records[0]
        result[arm]={"correct_certified":a["certified"],"safe_abstained":a["abstained"],
            "hard_false_positives":a["hard_false_positives"],
            "per_case":a["rows"],
            "learned_feature_count":a["learned_feature_count"],
            "global_verified_codes":a["globally_certified_codes"],
            "ambiguous_code_fallback_groups":a["TRAIN_code_collisions"],
            "world_to_predictive_class_merges":a["original_posterior_merge_count"],
            "full_signature_fallbacks":a["original_exact_residual_fallbacks"],
            "verified_learned_feature_hits":a["original_verified_feature_hits"],
            "five_complete_cpu_ns":[x["whole_process_CPU_ns"] for x in records],
            "five_peak_rss_bytes":[x["peak_process_RSS_bytes"] for x in records],
            "median_complete_cpu_ns":int(statistics.median(x["whole_process_CPU_ns"] for x in records)),
            "max_peak_rss_bytes":max(x["peak_process_RSS_bytes"] for x in records)}
    if not all(x["correct_certified"]+x["safe_abstained"]==192 and x["hard_false_positives"]==0 for x in result.values()):
        raise RuntimeError("T27_UNSOUND_OR_INCOMPLETE")
    cand=result[ARMS[0]]
    comparison={arm:{
      "no_baseline_correct_loss":cand["correct_certified"]>=obj["correct_certified"],
      "CPU_nonregression":cand["median_complete_cpu_ns"]<=obj["median_complete_cpu_ns"],
      "RSS_nonregression":cand["max_peak_rss_bytes"]<=obj["max_peak_rss_bytes"]
      } for arm,obj in result.items() if arm!=ARMS[0]}
    return {"schema":"wingless.T27.internal-original-predictive-state-comparison.v1",
      "TRAIN_cases":40,"prospective_cases":64,"counterfactual_queries":192,
      "source_families":"W_and_X_INTERNAL_ONLY",
      "arms":result,"comparison":comparison,
      "scientific_joint_acceptance":False,
      "independent_external_replication":False,"optimizer_admission":False,
      "old_DEV_replayed":False,"no_new_controllers":True,"no_paid_API":True}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(TRAIN_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T27_UNKNOWN_ARGS")
