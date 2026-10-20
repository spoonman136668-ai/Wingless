"""T41 source-blind learned causal composition selector, TRAIN-only.

The trained parameter is the prior probability of a compositional model.
Each newly observed Boolean-rule output changes posterior mixture weights
through an exact Bayes factor; all program and unrestricted law posteriors
remain exact. Source family (compositional vs unrestricted) and true rules
are private to source generation and evaluation, never to fitting/inference.
Fresh meta and validation fixtures; NO scientific original DEV is born.
"""
from __future__ import annotations
import json,random,statistics,subprocess,sys,time
from pathlib import Path
GATES=tuple(range(16));OP_COUNT=3
PROGRAMS=tuple((g,h,op) for g in GATES for h in GATES for op in range(OP_COUNT))
N_PROGRAMS=len(PROGRAMS)
FULL_RULES=tuple(range(256))
ARMS=("TRAIN-learned-sourceblind-Bayes-factor-selector",
      "fixed-50pct-sourceblind-mixture",
      "TRAIN-fitted-full256-categorical-EM",
      "TRAIN-fitted-compositional-factor-EM",
      "TRAIN-fitted-independent-Beta",
      "uniform-exact-unrestricted-256")
SEEDS={"C":{"META":41002026,"VALIDATION":41022026},
       "U":{"META":41012026,"VALIDATION":41032026}}
STRIDE=9817
META_N=32
VALID_N=16
OBSERVED=6
CHECKPOINTS=(2,4,6)
def peak_rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class M(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
                ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
                ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
                ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
                ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        k=ctypes.WinDLL("kernel32",use_last_error=True)
        p=ctypes.WinDLL("psapi",use_last_error=True)
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(M),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        mem=M();mem.cb=ctypes.sizeof(mem)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(mem),mem.cb):
            raise RuntimeError("T41_NATIVE_RSS_UNAVAILABLE")
        return int(mem.PeakWorkingSetSize)
    import resource
    r=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(r)*1024 if sys.platform.startswith("linux") else int(r)
def table(g,h,op):
    rule=0
    for q in range(8):
        a,b,c=q&1,(q>>1)&1,(q>>2)&1
        x=(g>>(a|(b<<1)))&1
        y=(h>>(b|(c<<1)))&1
        rule|=(x^y if op==0 else x&y if op==1 else x|y)<<q
    return rule
PROGRAM_LAWS=tuple(table(*program) for program in PROGRAMS)
COMPOSITION_PRIOR_RAW=[0.]*256
for law in PROGRAM_LAWS:COMPOSITION_PRIOR_RAW[law]+=1
COMPOSITION_PRIOR=tuple(v/N_PROGRAMS for v in COMPOSITION_PRIOR_RAW)
UNIFORM_PRIOR=(1/256,)*256
ELIGIBLE_PAIRS=tuple((g,h) for g in GATES for h in GATES if (g+3*h)%4!=0)
HOLDOUT_PAIRS=tuple((g,h) for g in GATES for h in GATES if (g+3*h)%4==0)
assert len(PROGRAMS)==768 and len(ELIGIBLE_PAIRS)==192 and len(HOLDOUT_PAIRS)==64
assert set(ELIGIBLE_PAIRS).isdisjoint(HOLDOUT_PAIRS)
def fixture(regime,split,index):
    if regime not in SEEDS or split not in ("META","VALIDATION"):
        raise RuntimeError("T41_SOURCE_FREE_PARTITION_ONLY")
    if not 0<=index<(META_N if split=="META" else VALID_N):
        raise RuntimeError("T41_INVALID_FROZEN_SOURCE_INDEX")
    r=random.Random(SEEDS[regime][split]+index*STRIDE)
    if regime=="U":
        law=r.randrange(256)
        program=None
    else:
        if split=="META" and index<16:
            g=index;h=(index+1)%16
            assert (g,h) in ELIGIBLE_PAIRS
        else:
            g,h=r.choice(ELIGIBLE_PAIRS if split=="META" else HOLDOUT_PAIRS)
        op=r.randrange(3)
        law=table(g,h,op)
        program=(g,h,op)
    order=tuple(r.sample(range(8),OBSERVED))
    public=tuple((q,(law>>q)&1) for q in order)
    return {"id":f"t41-sourcefree-{split}-{regime}-{index:03d}",
            "public":public,"_private_rule":law,
            "_private_regime":regime,"_private_program":program}
def make_meta():
    return tuple(fixture(regime,"META",i)
                 for regime in ("C","U") for i in range(META_N))
def make_validation():
    return tuple(fixture(regime,"VALIDATION",i)
                 for regime in ("C","U") for i in range(VALID_N))
def mask_signature(observations):
    mask=sig=0
    for pattern,bit in observations:
        if mask&(1<<pattern):
            if (sig>>pattern)&1!=bit:raise RuntimeError("T41_INCONSISTENT_PUBLIC_OBSERVATION")
        mask|=1<<pattern
        if bit:sig|=1<<pattern
    return mask,sig
def restricted(prior,mask,sig):
    values=tuple(w if (rule&mask)==sig else 0. for rule,w in enumerate(prior))
    mass=sum(values)
    if mass<=0:return mass,None
    return mass,values
def p_rule(prior,mask,sig,target):
    mass,values=restricted(prior,mask,sig)
    if mass==0:return None
    return sum(w for law,w in enumerate(values) if (law>>target)&1)/mass
def fit_mixture(public_train):
    # Unsupervised public-only mixture prior; no case category labels.
    pi=0.5
    for _ in range(4):
        acc=0.
        for observed in public_train:
            mask,sig=mask_signature(observed)
            s,_=restricted(COMPOSITION_PRIOR,mask,sig)
            u,_=restricted(UNIFORM_PRIOR,mask,sig)
            z=pi*s+(1-pi)*u
            if z<=0:raise RuntimeError("T41_ZERO_MIXTURE_TRAIN_PUBLIC_LIKELIHOOD")
            acc+=pi*s/z
        pi=(.05+acc)/(.1+len(public_train))
    if not .05<=pi<=.95:raise RuntimeError("T41_FITTED_PI_OUTSIDE_FROZEN_BOUNDS")
    return pi
def fit_categorical(public_train,law_space,iterations=4):
    size=len(law_space)
    prior=[1/size]*size
    for _ in range(iterations):
        accum=[.05/size]*size
        for observed in public_train:
            mask,sig=mask_signature(observed)
            allowed=[i for i,law in enumerate(law_space) if (law&mask)==sig]
            denom=sum(prior[i] for i in allowed)
            if denom<=0:raise RuntimeError("T41_CONTROL_EM_ZERO_SUPPORT")
            for i in allowed:accum[i]+=prior[i]/denom
        den=sum(accum);prior=[v/den for v in accum]
    return tuple(prior)
def fit_composition_factors(public_train):
    g=[1/16]*16;h=[1/16]*16;o=[1/3]*3
    for _ in range(4):
        a=[.05]*16;b=[.05]*16;c=[.05]*3
        for observed in public_train:
            mask,sig=mask_signature(observed)
            allowed=[i for i,law in enumerate(PROGRAM_LAWS) if (law&mask)==sig]
            if not allowed:continue # if public rules are not compositional, source-blind safe fallback
            total=sum(g[PROGRAMS[i][0]]*h[PROGRAMS[i][1]]*o[PROGRAMS[i][2]]
                      for i in allowed)
            if total<=0:raise RuntimeError("T41_TRAIN_COMPOSITION_PRIOR_ZERO")
            for i in allowed:
                gg,hh,oo=PROGRAMS[i]
                w=g[gg]*h[hh]*o[oo]/total
                a[gg]+=w;b[hh]+=w;c[oo]+=w
        g=[x/sum(a) for x in a]
        h=[x/sum(b) for x in b]
        o=[x/sum(c) for x in c]
    full=[0.]*256
    for i,law in enumerate(PROGRAM_LAWS):
        gg,hh,oo=PROGRAMS[i]
        full[law]+=g[gg]*h[hh]*o[oo]
    return tuple(full)
def fit_beta(public_train):
    c=[[2.,2.] for _ in range(8)]
    for rows in public_train:
        for q,v in rows:c[q][v]+=1
    return tuple(row[1]/sum(row) for row in c)
def fit(arm,public_train):
    if arm==ARMS[0]:return fit_mixture(public_train)
    if arm==ARMS[1]:return 0.5
    if arm==ARMS[2]:return fit_categorical(public_train,FULL_RULES)
    if arm==ARMS[3]:return fit_composition_factors(public_train)
    if arm==ARMS[4]:return fit_beta(public_train)
    if arm==ARMS[5]:return UNIFORM_PRIOR
    raise RuntimeError("T41_UNKNOWN_FITTED_ARM")
def posterior_structure_weight(pi,mask,sig):
    a,_=restricted(COMPOSITION_PRIOR,mask,sig)
    b,_=restricted(UNIFORM_PRIOR,mask,sig)
    z=pi*a+(1-pi)*b
    if z<=0:raise RuntimeError("T41_MIXTURE_EMPTY_POSTERIOR")
    return pi*a/z
def forecast(model,arm,mask,sig,target):
    # No private source category, no truth table, no future answer.
    if arm in ARMS[:2]:
        w=posterior_structure_weight(model,mask,sig)
        su=p_rule(COMPOSITION_PRIOR,mask,sig,target)
        uv=p_rule(UNIFORM_PRIOR,mask,sig,target)
        if su is None:
            if abs(w)>1e-12:raise RuntimeError("T41_CONTRADICTORY_ZERO_STRUCTURED_SUPPORT")
            return uv,w
        return w*su+(1-w)*uv,w
    if arm==ARMS[4]:
        if mask&(1<<target):return float((sig>>target)&1),None
        return model[target],None
    val=p_rule(model,mask,sig,target)
    if val is None:
        # Robust source-blind compositional fallback if the public history
        # is incompatible with every structured rule.
        if arm==ARMS[3]:return p_rule(UNIFORM_PRIOR,mask,sig,target),None
        raise RuntimeError("T41_FULL_RULE_ZERO_POSTERIOR")
    return val,None
def independent_certification(observed,target,law):
    mask,sig=mask_signature(observed)
    allowed=[r for r in range(256) if (r&mask)==sig]
    if not allowed:raise RuntimeError("T41_ALL256_CAUSAL_POSTERIOR_EMPTY")
    possible={((rule>>target)&1) for rule in allowed}
    if ((law>>target)&1) not in possible:
        raise RuntimeError("T41_HIDDEN_SOURCE_RULE_EXCLUDED")
    return next(iter(possible)) if len(possible)==1 else None
def score(arm,model,public,private_law):
    result={"correct":0,"unresolved":0,"brier_unknown":0.,"n":0,
            "posteriors":[],"forecasts":[]}
    for k in CHECKPOINTS:
        visible=public[:k]
        mask,sig=mask_signature(visible)
        for q in range(8):
            p,weight=forecast(model,arm,mask,sig,q)
            if not 0<=p<=1:raise RuntimeError("T41_UNBOUNDED_BELIEF")
            cert=independent_certification(visible,q,private_law)
            truth=(private_law>>q)&1
            if cert is None:
                result["unresolved"]+=1
                result["brier_unknown"]+=(p-truth)**2
            else:
                if cert!=truth or abs(p-truth)>1e-8:
                    raise RuntimeError("T41_INCORRECT_HARD_CERTIFICATION_OR_PROBABILITY")
                result["correct"]+=1
            result["n"]+=1
            result["forecasts"].append(round(p,11))
            if weight is not None:result["posteriors"].append(round(weight,11))
    return result
def validate():
    meta=make_meta()
    train_public=[x["public"] for x in meta] # explicit removal of private source labels
    trained={a:fit(a,train_public) for a in ARMS}
    learned_pi=trained[ARMS[0]]
    if abs(learned_pi-.5)<1e-12:raise RuntimeError("T41_LEARNED_PI_NOT_DIFFERENT_FROM_FIXED_CONTROL")
    gseen={x["_private_program"][0] for x in meta if x["_private_program"] is not None}
    hseen={x["_private_program"][1] for x in meta if x["_private_program"] is not None}
    if gseen!=set(range(16)) or hseen!=set(range(16)):
        raise RuntimeError("T41_PUBLIC_META_ATOM_ROLE_COVERAGE_NOT_COMPLETE")
    samples=make_validation()
    scores={rg:{a:{"correct":0,"unresolved":0,"Brier_sum":0.} for a in ARMS}
            for rg in ("C","U")}
    learned_changed=0;posterior_changes=0
    for item in samples:
        rg=item["_private_regime"]
        if rg=="C" and (item["_private_program"][0],item["_private_program"][1]) not in HOLDOUT_PAIRS:
            raise RuntimeError("T41_NEW_COMPOSITION_NOT_IN_HELDOUT_PARTITION")
        results={a:score(a,trained[a],item["public"],item["_private_rule"]) for a in ARMS}
        learned=results[ARMS[0]];fixed=results[ARMS[1]]
        if any(abs(p-q)>1e-10 for p,q in zip(learned["forecasts"],fixed["forecasts"])):
            learned_changed+=1
        if len(set(learned["posteriors"]))>1:posterior_changes+=1
        for a,result in results.items():
            dest=scores[rg][a]
            dest["correct"]+=result["correct"]
            dest["unresolved"]+=result["unresolved"]
            dest["Brier_sum"]+=result["brier_unknown"]
    for group in scores.values():
        for entry in group.values():
            entry["Brier"]=entry["Brier_sum"]/entry["unresolved"]
    structured=scores["C"]
    best=min((structured[a]["Brier"],a) for a in ARMS[1:])
    model=structured[ARMS[0]]["Brier"]
    gain=(best[0]-model)/best[0] if best[0] else 0.
    unrestricted=scores["U"]
    u=unrestricted[ARMS[0]]["Brier"]
    floor=unrestricted[ARMS[5]]["Brier"]
    regression=(u-floor)/floor if floor else 0.
    return {"schema":"wingless.T41.sourceblind-structure-selection-TRAIN-validation.v1",
      "status":"PASS_TECHNICAL_EXACT_TRAIN_ONLY",
      "public_META_cases":len(meta),"META_structured":32,"META_unrestricted":32,
      "TRAIN_validation_structured_cases":16,"TRAIN_validation_unrestricted_cases":16,
      "scientific_DEV_created":0,"old_research_DEV_reused":False,
      "G_role_complete_atoms":len(gseen),"H_role_complete_atoms":len(hseen),
      "learned_pi":learned_pi,
      "structural_parameter_ablation_changed_worlds":learned_changed,
      "posterior_structural_weight_changed_with_observation_cases":posterior_changes,
      "private_regime_labels_not_used_in_fit_or_inference":True,
      "hard_certification_oracle_independent_all_256_rules":True,
      "results":scores,"structured_gain_vs_best_strong_control":gain,
      "strongest_structured_control":best[1],
      "unrestricted_regression_relative_uniform":regression,
      "structured_5pct_gain_pass":gain>=.05,
      "unrestricted_max_1pct_regression_pass":regression<=.01,
      "optimizer_admission":False}
def native(arm,repeat):
    if arm not in ARMS or repeat not in range(5):
        raise RuntimeError("T41_INVALID_ORIGINAL_ARM")
    start=time.process_time_ns()
    meta=make_meta()
    model=fit(arm,[x["public"] for x in meta])
    values={"C":{"correct":0,"unresolved":0,"Brier_sum":0.},
            "U":{"correct":0,"unresolved":0,"Brier_sum":0.}}
    public_fingerprints=[]
    for item in make_validation():
        result=score(arm,model,item["public"],item["_private_rule"])
        target=values[item["_private_regime"]]
        for src,dst in (("correct","correct"),("unresolved","unresolved"),("brier_unknown","Brier_sum")):
            target[dst]+=result[src]
        public_fingerprints.append(item["id"]+":"+str(len(item["public"])))
    return {"schema":"wingless.T41.unchanged-native-TRAIN-process.v1",
      "arm":arm,"repeat":repeat,"META_cases":64,"TRAIN_validation_cases":32,
      "scientific_DEV_created":0,"results":values,
      "source_fingerprints":public_fingerprints,
      "whole_CPU_ns":time.process_time_ns()-start,"peak_RSS_bytes":peak_rss()}
def cost():
    costs={}
    for arm in ARMS:
        trials=[]
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                  "--arm",arm,str(rep)],capture_output=True,text=True)
            if p.returncode:raise RuntimeError("T41_NATIVE_CHILD_FAILURE:"+arm+"\n"+p.stderr[-3000:])
            trials.append(json.loads(p.stdout))
        first=trials[0]
        if any(x["results"]!=first["results"] or x["source_fingerprints"]!=first["source_fingerprints"] for x in trials):
            raise RuntimeError("T41_UNSTABLE_FROZEN_NATIVE_SOURCE")
        costs[arm]={"original_five_whole_CPU_ns":[x["whole_CPU_ns"] for x in trials],
          "original_five_peak_RSS_bytes":[x["peak_RSS_bytes"] for x in trials],
          "median_whole_CPU_ns":int(statistics.median(x["whole_CPU_ns"] for x in trials)),
          "max_peak_RSS_bytes":max(x["peak_RSS_bytes"] for x in trials),
          "correct_certified":sum(v["correct"] for v in first["results"].values()),
          "safely_unresolved":sum(v["unresolved"] for v in first["results"].values()),
          "Brier_by_source_regime":{rg:v["Brier_sum"]/v["unresolved"] for rg,v in first["results"].items()}}
    reference=costs[ARMS[-1]]
    if any((v["correct_certified"],v["safely_unresolved"])!=(
       reference["correct_certified"],reference["safely_unresolved"]) for v in costs.values()):
        raise RuntimeError("T41_NATIVE_CERTIFICATION_COVERAGE_DIFFERENCE")
    learned=costs[ARMS[0]]
    return {"schema":"wingless.T41.original-TRAIN-native-cost.v1",
        "sourcefree_META":64,"sourcefree_VALIDATION":32,"scientific_DEV_created":0,
        "arms":costs,
        "strict_CPU_better_than_all_controls":all(learned["median_whole_CPU_ns"]<costs[a]["median_whole_CPU_ns"] for a in ARMS[1:]),
        "peak_RSS_nondegression_vs_best_control":learned["max_peak_RSS_bytes"]<=min(x["max_peak_RSS_bytes"] for x in costs.values()),
        "scientific_positive":False,"optimizer_admission":False}
if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(validate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--native":
        print(json.dumps(cost(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(native(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T41_NO_PROSPECTIVE_DEV_MODE")
