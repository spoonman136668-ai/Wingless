"""T41 frozen source-blind model selection and fair source-free controls.

A learned forecast is never a hard truth certificate. Exact consistency uses
all 256 possible three-input Boolean laws independent of the learned prior.
No original science/DEV opened; native scientific admission is not authorized.
"""
from __future__ import annotations
import importlib.util,json,math,statistics,subprocess,sys,time
from pathlib import Path

HERE=Path(__file__).resolve().parent
ARMS=("sourceblind-mixture","trained-compositional","trained-joint768",
      "trained-full256","trained-beta","uniform-exact","half-mixture")
MIX_GRID=(0.0,0.25,0.5,0.75,1.0)
REGIMES=("unrestricted256","uniform_composition768","nonuniform_composition768")
BUDGETS=(2,4,6)

def import_source(name,filename):
    spec=importlib.util.spec_from_file_location(name,HERE/filename)
    mod=importlib.util.module_from_spec(spec)
    sys.modules[name]=mod
    spec.loader.exec_module(mod)
    return mod

def deps():
    return (import_source("t41_t40_math","20261020-wingless-t40-exact-learnability-audit.py"),
            import_source("t41_source","20261020-wingless-t41-fresh-sourcefree.py"))

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
            raise RuntimeError("T41_RSS_UNAVAILABLE")
        return int(m.PeakWorkingSetSize)
    import resource
    val=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(val)*1024 if sys.platform.startswith("linux") else int(val)

def public_likelihood(prior,mask,signature):
    return sum(w for law,w in enumerate(prior) if law&mask==signature)

def selector_weight(comp_prior,validation):
    scores=[]
    for fraction in MIX_GRID:
        score=0.0
        for mask,sig in validation:
            c=public_likelihood(comp_prior,mask,sig)
            u=1/(1<<mask.bit_count())
            m=fraction*c+(1-fraction)*u
            if m<=0:
                score=float("-inf")
                break
            score+=math.log(m)
        scores.append((score,-fraction,fraction))
    return max(scores)[2]

def fit_model(core,arm,public_train):
    if arm in ("sourceblind-mixture","half-mixture"):
        if arm=="sourceblind-mixture":
            base=core.project_programs(core.train_factors(public_train[:48]))
            mixture=selector_weight(base,public_train[48:])
        else:
            mixture=0.5
        fitted=core.project_programs(core.train_factors(public_train))
        uniform=1/256
        weights=tuple(mixture*w+(1-mixture)*uniform for w in fitted)
        return weights,mixture
    if arm=="trained-compositional":
        return core.project_programs(core.train_factors(public_train)),1.0
    if arm=="trained-joint768":
        return core.project_programs(core.train_joint(public_train,768,core.LAW)),None
    if arm=="trained-full256":
        return core.train_joint(public_train,256,tuple(range(256))),None
    if arm=="trained-beta":
        return core.independent_beta(public_train),None
    if arm=="uniform-exact":
        return (1/256,)*256,None
    raise ValueError("T41_UNKNOWN_ARM")

def inferred_probability(model,arm,mask,signature,q):
    assert not (mask&(1<<q))
    if arm=="trained-beta":
        return model[q],False
    if arm=="uniform-exact":
        return 0.5,False
    mass=0.;positive=0.
    for law,w in enumerate(model):
        if law&mask==signature:
            mass+=w
            if (law>>q)&1:
                positive+=w
    if mass<=1e-14:
        # A restricted prior may assign zero probability to a valid
        # out-of-grammar public history. Treat as a real support failure,
        # use explicitly noncertifying neutral prediction, and report it.
        return 0.5,True
    return min(1.0,max(0.0,positive/mass)),False

def exact_hard_certificate(mask,signature,q):
    matching=[law for law in range(256) if law&mask==signature]
    assert len(matching)==1<<(8-mask.bit_count())
    outcomes={(law>>q)&1 for law in matching}
    return next(iter(outcomes)) if len(outcomes)==1 else None

def preeval_math_properties():
    # No source gen, no TRAIN, no previously used DEV.
    checks=0
    for k in BUDGETS:
        for mask in range(256):
            if mask.bit_count()!=k:
                continue
            bits=[j for j in range(8) if not mask&(1<<j)]
            for sig in range(256):
                if sig&~mask:continue
                for q in bits:
                    assert exact_hard_certificate(mask,sig,q) is None
                    checks+=1
    assert checks==sum(math.comb(8,k)*(1<<k)*(8-k) for k in BUDGETS)
    assert len(MIX_GRID)==5
    # Trained-prior selector must never have access to the generating label.
    return {"T41_PUBLIC_ONLY_EXACT_256_HARD_PROOF_CHECKS":checks,
      "sourcefree_eval_consumed":False,"scientific_DEV_opened":False,
      "T40_previous_public_cases_replayed":False}

def one_process(arm,rep):
    assert arm in ARMS and 0<=rep<5
    began=time.process_time_ns()
    core,source=deps()
    assert source.REGIMES==REGIMES and source.N_TRAIN==64 and source.N_EVAL==96
    datasets={}
    # Every arm bears identical source-free fixture construction.
    for regime in REGIMES:
        train,ct=source.cases(regime,"train",core)
        evaluation,ce=source.cases(regime,"eval",core)
        datasets[regime]=(train,evaluation,ct,ce)
    scores={};mix_weights={};source_custody={}
    support_failures=0
    examined=0
    wrong_certifications=0
    for regime,(train,evaluation,ct,ce) in datasets.items():
        public=[x["public"] for x in train]
        predictor,mix=fit_model(core,arm,public)
        mix_weights[regime]=mix
        source_custody[regime]={"train":ct,"evaluation":ce}
        row={}
        for k in BUDGETS:
            loss=0.0;count=0
            for fixture in evaluation:
                mask,signature=fixture["public_observations"][k]
                rule=fixture["private_truth_rule"]
                assert rule&mask==signature
                for q in range(8):
                    if mask&(1<<q):continue
                    hard=exact_hard_certificate(mask,signature,q)
                    if hard is not None:
                        wrong_certifications+=int(hard!=((rule>>q)&1))
                    forecast,unsupported=inferred_probability(predictor,arm,mask,signature,q)
                    assert 0<=forecast<=1
                    support_failures+=int(unsupported)
                    truth=(rule>>q)&1
                    loss+=(forecast-truth)**2
                    count+=1
            assert count==source.N_EVAL*(8-k)
            row[str(k)]={"Brier":loss/count,"targets":count}
            examined+=count
        scores[regime]=row
    assert wrong_certifications==0
    assert examined==source.N_EVAL*(6+4+2)*3
    return {"schema":"wingless.T41.sourcefree-process.v1","arm":arm,
      "rep":rep,"sourcefree_scores":scores,"selector_weights":mix_weights,
      "meta_custody":source_custody if rep==0 else None,
      "out_of_support_forecasts":support_failures,
      "unseen_target_exact_false_certifications":wrong_certifications,
      "complete_process_cpu_ns":time.process_time_ns()-began,
      "peak_process_rss_bytes":peak_rss_bytes(),"scored_targets":examined,
      "scientific_DEV_opened":0,"original_previous_cases_replayed":False,
      "source_born_independent":False,"remote_query_service_cost":"UNMEASURED",
      "scientific_acceptance":False,"optimizer_admission":False}

def run():
    all_results={}
    for arm in ARMS:
        proc=[]
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
             "--child",arm,str(rep)],capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T41_ORIGINAL_CHILD_FAILED:"+arm+":"+str(rep)+"\n"+p.stderr)
            proc.append(json.loads(p.stdout))
        assert all(x["sourcefree_scores"]==proc[0]["sourcefree_scores"] for x in proc)
        all_results[arm]={
          "scores":proc[0]["sourcefree_scores"],
          "selector_weights":proc[0]["selector_weights"],
          "source_generation_audit":proc[0]["meta_custody"],
          "scored_targets":proc[0]["scored_targets"],
          "out_of_support_forecasts":proc[0]["out_of_support_forecasts"],
          "hard_false_certifications":proc[0]["unseen_target_exact_false_certifications"],
          "five_cpu_ns":[x["complete_process_cpu_ns"] for x in proc],
          "five_peak_rss_bytes":[x["peak_process_rss_bytes"] for x in proc],
          "median_complete_cpu_ns":int(statistics.median(x["complete_process_cpu_ns"] for x in proc)),
          "max_process_rss_bytes":max(x["peak_process_rss_bytes"] for x in proc)}
    learned=all_results["sourceblind-mixture"]
    dominant={}
    for regime in REGIMES:
        dominant[regime]={}
        for k in BUDGETS:
            best=min(((all_results[arm]["scores"][regime][str(k)]["Brier"],arm)
                 for arm in ARMS if arm!="sourceblind-mixture"))
            actual=learned["scores"][regime][str(k)]["Brier"]
            dominant[regime][str(k)]={"learned":actual,"best_control":best[1],
              "best_control_Brier":best[0],"Brier_regret":actual-best[0],
              "learner_beats_strongest_control":actual<best[0]}
    uniform=learned["scores"]["unrestricted256"]
    safety=all(uniform[str(k)]["Brier"]<=.26 for k in BUDGETS)
    structured=all(dominant[regime][str(k)]["learner_beats_strongest_control"]
               for regime in REGIMES[1:] for k in BUDGETS)
    influence=any(learned["scores"][regime]!=all_results["half-mixture"]["scores"][regime]
                  for regime in REGIMES)
    # Cost condition is intentionally strict against every eligible strong
    # control; sourcefree internal success is not scientific first-use.
    native=all(learned["median_complete_cpu_ns"]<=r["median_complete_cpu_ns"]
        and learned["max_process_rss_bytes"]<=r["max_process_rss_bytes"]
        for a,r in all_results.items() if a!="sourceblind-mixture")
    return {"schema":"wingless.T41.sourceblind-regime-model-selection.v1",
     "status":"INTERNAL_SOURCEFREE_COMPARISON_NOT_SCIENCE","arms":all_results,
     "comparisons":dominant,"unrestricted_Brier_noncollapse":safety,
     "structured_all_strong_controls_beaten":structured,
     "learned_selector_changes_forecasts_vs_fixed_half":influence,
     "two_OS_full_native_required_native_regression_pass_this_OS":native,
     "complete_joint_admission":False,
     "joint_admission_reason":"UNMEASURED_QUERY_SERVICE_COST_AND_NO_INDEPENDENT_SCIENTIFIC_FIRST_USE",
     "scientific_source_generated":False,"scientific_DEV_opened":0,
     "T40_previous_cases_replayed":False,
     "no_hard_false_certification":all(r["hard_false_certifications"]==0 for r in all_results.values()),
     "scientific_acceptance":False,"optimizer_admission":False,
     "no_merge":True,"no_paid_inference":True,"no_new_controller":True}

if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--property":
        print(json.dumps(preeval_math_properties(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--child":
        print(json.dumps(one_process(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(run(),sort_keys=True))
    else:
        raise SystemExit("T41_INVALID_INVOCATION")
