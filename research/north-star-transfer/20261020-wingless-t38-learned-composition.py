"""T38: public-TRAIN-learned compositional causal-rule induction.

Fit P(G), P(H), P(OP) over public-consistent latent compositions using
weakly supervised EM. This transfers evidence about individual gate atoms
to previously unused pairs; no hidden gate labels or future truth seen.
The equally trained joint-composition control fits all 128 pair+op priors,
the full rule control fits 256 truth-table priors. Independent-pattern and
uniform exact models receive the same public evidence.

Probability of unknown future outcomes is NOT a logical certification.
Every hard certification is verified by enumerating all 256 possible
Boolean rules consistent with available public causal observations.
Original unseen DEV source creation is impossible through this module's
supported CLI. Only distinct TRAIN-meta and TRAIN-validation are opened.
"""
from __future__ import annotations
import importlib.util,json,statistics,subprocess,sys,time
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-learned-compositional-gate-factors",
      "TRAIN-fitted-joint-composition-categorical",
      "TRAIN-fitted-full-256-rule-empirical-Bayes",
      "TRAIN-fitted-independent-Beta-pattern",
      "uniform-exact-version-space")
def source():
    spec=importlib.util.spec_from_file_location("t38_new_causal_source",
                  ROOT/"20261020-wingless-t38-new-compositional-sources.py")
    mod=importlib.util.module_from_spec(spec)
    sys.modules[spec.name]=mod
    spec.loader.exec_module(mod)
    return mod
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
        kernel=ctypes.WinDLL("kernel32",use_last_error=True)
        psapi=ctypes.WinDLL("psapi",use_last_error=True)
        kernel.GetCurrentProcess.restype=wintypes.HANDLE
        psapi.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        psapi.GetProcessMemoryInfo.restype=wintypes.BOOL
        mem=Mem();mem.cb=ctypes.sizeof(mem)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(mem),mem.cb):
            raise RuntimeError("T38_PEAK_RSS_UNAVAILABLE:"+str(ctypes.get_last_error()))
        return int(mem.PeakWorkingSetSize)
    import resource
    x=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(x)*1024 if sys.platform.startswith("linux") else int(x)

def observed_rule_bits(public,src):
    known={}
    for trial in public["probes"]:
        for node,output in trial["observed"]:
            action=trial["do"]
            if action is not None and action[0]==node:
                if output!=action[1]:raise RuntimeError("T38_PUBLIC_INTERVENTION_INCONSISTENT")
                continue
            p,v=src.public_clause(trial["history"],public["topology"],node,output)
            if p in known and known[p]!=v:
                raise RuntimeError("T38_PUBLIC_CAUSAL_RULE_CONTRADICTION")
            known[p]=v
    return known
def consistent(rule,observed):
    return all(((rule>>p)&1)==bit for p,bit in observed.items())
def constraints(src,training):
    return [observed_rule_bits(public,src) for public in training]
def compositional_weights(g,h,op,model):
    return model["G"][g]*model["H"][h]*model["OP"][op]
def train(src,publics,arm):
    cases=constraints(src,publics)
    if arm==ARMS[0]:
        G=[1/8]*8;H=[1/8]*8;OP=[0.5,0.5]
        for sweep in range(4):
            cg=[0.05]*8;ch=[0.05]*8;co=[0.05]*2
            for evidence in cases:
                possible=[(g,h,op,G[g]*H[h]*OP[op])
                    for g,h,op,rule in src.HYPOTHESES if consistent(rule,evidence)]
                denom=sum(w for _,_,_,w in possible)
                if denom<=0:raise RuntimeError("T38_FACTORIZED_TRAIN_ZERO_EVIDENCE")
                for g,h,op,w in possible:
                    q=w/denom;cg[g]+=q;ch[h]+=q;co[op]+=q
            G=[x/sum(cg) for x in cg]
            H=[x/sum(ch) for x in ch]
            OP=[x/sum(co) for x in co]
        return {"G":tuple(G),"H":tuple(H),"OP":tuple(OP),"cases":len(cases)}
    if arm==ARMS[1]:
        joint=[1/128]*128
        for sweep in range(4):
            counts=[0.05/128]*128
            for evidence in cases:
                available=[i for i,(_,_,_,rule) in enumerate(src.HYPOTHESES)
                           if consistent(rule,evidence)]
                denom=sum(joint[i] for i in available)
                if denom<=0:raise RuntimeError("T38_JOINT_TRAIN_ZERO_EVIDENCE")
                for i in available:counts[i]+=joint[i]/denom
            total=sum(counts)
            joint=[v/total for v in counts]
        return {"joint":tuple(joint),"cases":len(cases)}
    if arm==ARMS[2]:
        full=[1/256]*256
        for sweep in range(4):
            counts=[0.05/256]*256
            for evidence in cases:
                available=[rule for rule in range(256) if consistent(rule,evidence)]
                denom=sum(full[r] for r in available)
                if denom<=0:raise RuntimeError("T38_FULL_RULE_TRAIN_ZERO_EVIDENCE")
                for rule in available:counts[rule]+=full[rule]/denom
            total=sum(counts);full=[p/total for p in counts]
        return {"full":tuple(full),"cases":len(cases)}
    if arm==ARMS[3]:
        bits=[[2.,2.] for _ in range(8)]
        for evidence in cases:
            for pattern,value in evidence.items():bits[pattern][value]+=1
        return {"independent":tuple(p[1]/sum(p) for p in bits),"cases":len(cases)}
    if arm==ARMS[4]:
        return {"uniform":True,"cases":len(cases)}
    raise RuntimeError("T38_UNKNOWN_COMPARATOR")
def future(src,public,known,model,arm):
    all_possible=tuple(rule for rule in range(256) if consistent(rule,known))
    if not all_possible:raise RuntimeError("T38_EXACT_CERTIFY_NO_COMPATIBLE_LAW")
    if arm in ARMS[:2]:
        support=[(i,g,h,op,rule) for i,(g,h,op,rule) in enumerate(src.HYPOTHESES)
                  if consistent(rule,known)]
        if not support:raise RuntimeError("T38_CAUSAL_GATE_GRAMMAR_NO_POSTERIOR")
        if arm==ARMS[0]:
            weights=[compositional_weights(g,h,op,model) for i,g,h,op,rule in support]
        else:
            weights=[model["joint"][i] for i,g,h,op,rule in support]
    elif arm==ARMS[2]:
        support=all_possible
        weights=[model["full"][rule] for rule in support]
    else:
        support=()
        weights=()
    norm=sum(weights)
    if support and norm<=0:raise RuntimeError("T38_TRAINED_CAUSAL_PRIOR_SUPPORT_ZERO")
    out=[]
    for q in public["queries"]:
        node=q["node"]
        action=q["do"]
        if action is not None and node==action[0]:
            allowed={action[1]}
            estimated=float(action[1])
        else:
            pattern=src.pattern(q["history"],public["topology"][node])
            residual=(q["history"][0]>>node)&1
            allowed={residual^((rule>>pattern)&1) for rule in all_possible}
            if arm==ARMS[3]:
                p=model["independent"][pattern]
                estimated=p if residual==0 else 1-p
            elif arm==ARMS[4]:
                estimated=sum((residual^((rule>>pattern)&1)) for rule in all_possible)/len(all_possible)
            else:
                if arm in ARMS[:2]:
                    weighted=sum(w*(residual^((item[4]>>pattern)&1))
                             for item,w in zip(support,weights))
                else:
                    weighted=sum(w*(residual^((rule>>pattern)&1))
                             for rule,w in zip(support,weights))
                estimated=weighted/norm
        if not 0<=estimated<=1:raise RuntimeError("T38_INVALID_PROBABILITY")
        exact=next(iter(allowed)) if len(allowed)==1 else None
        out.append({"P1":estimated,"certified":exact,"unresolved":exact is None,
                    "full_256_version_space_size":len(all_possible)})
    return out
def evaluation():
    src=source()
    if len(src.TRAIN_META)!=40 or len(src.TRAIN_VALID)!=16 or len(src.DEV)!=64:
        raise RuntimeError("T38_FROZEN_SOURCE_COHORT_DRIFT")
    if set(src.TRAIN_META)&set(src.TRAIN_VALID) or (set(src.TRAIN_META)|set(src.TRAIN_VALID))&set(src.DEV):
        raise RuntimeError("T38_HOLDOUT_IDENTITY_OVERLAP")
    training=[src.META(src.build(name)) for name in src.TRAIN_META]
    models={arm:train(src,training,arm) for arm in ARMS}
    scores={arm:{"all_squared":0.0,"unresolved_squared":0.0,"n":0,"uncertain":0,
        "certified":0,"false_certified":0,"rows":[]} for arm in ARMS}
    changed=0;source_true_retained=0;seen_train_pairs=set()
    for name in src.TRAIN_META:
        world=src.build(name)
        seen_train_pairs.add((world["_private_g"],world["_private_h"]))
    for name in src.TRAIN_VALID:
        world=src.build(name);public=src.VALID(world)
        known=observed_rule_bits(public,src)
        if not consistent(world["_private_rule"],known):
            raise RuntimeError("T38_TRUE_HIDDEN_CAUSAL_PROGRAM_ELIMINATED")
        if (world["_private_g"],world["_private_h"]) not in src.TRAIN_PAIRS:
            raise RuntimeError("T38_META_VALID_SOURCE_COMPOSITION_LEAK")
        source_true_retained+=1
        predictions={a:future(src,public,known,models[a],a) for a in ARMS}
        # Source-free ablation: turn OFF all learned factor probabilities,
        # but retain identical compositional grammar and evidence.
        uniform_comp={"G":(1/8,)*8,"H":(1/8,)*8,"OP":(0.5,0.5)}
        ablated=future(src,public,known,uniform_comp,ARMS[0])
        if any(abs(z["P1"]-a["P1"])>1e-9 for z,a in zip(predictions[ARMS[0]],ablated)):
            changed+=1
        for a,rows in predictions.items():
            correct=unknown=0
            for q,p in enumerate(rows):
                y=world["_answers"][q]
                if p["certified"] is not None:
                    if p["certified"]!=y:raise RuntimeError("T38_WRONG_HARD_CAUSAL_CERTIFICATION")
                    correct+=1;scores[a]["certified"]+=1
                else:
                    unknown+=1;scores[a]["uncertain"]+=1
                    scores[a]["unresolved_squared"]+=(p["P1"]-y)**2
                scores[a]["all_squared"]+=(p["P1"]-y)**2;scores[a]["n"]+=1
            scores[a]["rows"].append({"source_id":name,
              "certified":correct,"unresolved":unknown,
              "P1":[round(p["P1"],10) for p in rows],
              "private_truth_scoring_only":world["_answers"]})
    for a,s in scores.items():
        s["Brier_unresolved"]=s["unresolved_squared"]/max(s["uncertain"],1)
        s["Brier_all"]=s["all_squared"]/max(s["n"],1)
    if any((s["certified"],s["uncertain"])!=(
       scores[ARMS[4]]["certified"],scores[ARMS[4]]["uncertain"]) for s in scores.values()):
        raise RuntimeError("T38_EXACT_CERTIFICATION_MISMATCH")
    best=min(scores[a]["Brier_unresolved"] for a in ARMS[1:])
    learner=scores[ARMS[0]]["Brier_unresolved"]
    return {"schema":"wingless.T38.TRAIN-only-compositional-transfer-validation.v1",
       "source":src,"models":models,
       "TRAIN_meta":40,"TRAIN_validation":16,"DEV_opened":0,
       "score":scores,"true_causal_rules_retained":source_true_retained,
       "new_source_train_pairs_seen":len(seen_train_pairs),
       "source_partition_frozen":True,
       "learned_factor_ablated_changes_forecasts_in_cases":changed,
       "best_control_unresolved_Brier":best,
       "learned_Brier_relative_gain":(best-learner)/best if best else 0.0,
       "full_256_exact_certifier":True,
       "old_consumed_DEV_reused":False,"optimizer_admission":False}
def proof():
    e=evaluation()
    scores=e["score"];s=scores[ARMS[0]]
    return {"schema":"wingless.T38.original-TRAIN-only-compositional-causal-rule-proof.v1",
       "status":"PASS_TECHNICAL_TRAIN_ONLY",
       "TRAIN_meta_cases":e["TRAIN_meta"],"TRAIN_validation_cases":e["TRAIN_validation"],
       "DEV_generated":0,"TRAIN_worlds_true_rule_retained":e["true_causal_rules_retained"],
       "TRAIN_heldout_forecasts":s["n"],"correct_certified":s["certified"],
       "sound_abstentions":s["uncertain"],"hard_false_certified":0,
       "unresolved_fraction":s["uncertain"]/s["n"],
       "learned_parameters_changed_actual_belief_worlds":e["learned_factor_ablated_changes_forecasts_in_cases"],
       "learned_Brier_unresolved":s["Brier_unresolved"],
       "competitor_Brier_unresolved":{a:z["Brier_unresolved"] for a,z in scores.items() if a!=ARMS[0]},
       "learned_Brier_gain_vs_best_equally_trained":e["learned_Brier_relative_gain"],
       "full_256_truth_table_certification_verified":True,
       "source_new_pair_partition_sealed":True,
       "old_DEV_unopened":True,"optimizer_admission":False}
def native(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T38_NATIVE_ARM_INVALID")
    begin=time.process_time_ns()
    src=source()
    train_public=[src.META(src.build(s)) for s in src.TRAIN_META]
    fit=train(src,train_public,arm)
    entries=[];correct=unknown=0;unresolved_sum=0.
    for name in src.TRAIN_VALID:
        world=src.build(name);public=src.VALID(world)
        predictions=future(src,public,observed_rule_bits(public,src),fit,arm)
        yes=no=0
        for p,y in zip(predictions,world["_answers"]):
            if p["certified"] is None:
                unknown+=1;no+=1
                unresolved_sum+=(p["P1"]-y)**2
            else:
                if p["certified"]!=y:raise RuntimeError("T38_NATIVE_FALSE_HARD_CERTIFICATION")
                correct+=1;yes+=1
        entries.append({"id":name,"correct":yes,"unresolved":no,
            "probabilities":[round(x["P1"],10) for x in predictions]})
    return {"schema":"wingless.T38.full-native-TRAIN-process.v1","arm":arm,
      "repeat":rep,"meta_TRAIN_cases":40,"validation_cases":16,
      "DEV_generated":0,"correct_certified":correct,"unresolved":unknown,
      "hard_false_certifications":0,"Brier_unresolved":unresolved_sum/max(unknown,1),
      "per_case":entries,"complete_CPU_ns":time.process_time_ns()-begin,
      "peak_RSS_bytes":rss()}
def cost():
    out={}
    for arm in ARMS:
        trials=[]
        for rep in range(5):
            proc=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                    "--arm",arm,str(rep)],capture_output=True,text=True)
            if proc.returncode:raise RuntimeError("T38_NATIVE_CHILD_FAILED:"+arm+"\n"+proc.stderr[-3000:])
            trials.append(json.loads(proc.stdout))
        first=trials[0]
        if any(r["per_case"]!=first["per_case"] for r in trials):
            raise RuntimeError("T38_TRAIN_NONDETERMINISTIC")
        out[arm]={"TRAIN_validation_cases":16,"correct_certified":first["correct_certified"],
         "sound_abstentions":first["unresolved"],"hard_false_certifications":0,
         "Brier_unresolved":first["Brier_unresolved"],
         "five_full_CPU_ns":[x["complete_CPU_ns"] for x in trials],
         "five_peak_RSS_bytes":[x["peak_RSS_bytes"] for x in trials],
         "median_full_CPU_ns":int(statistics.median(x["complete_CPU_ns"] for x in trials)),
         "max_peak_RSS_bytes":max(x["peak_RSS_bytes"] for x in trials)}
    candidate=out[ARMS[0]]
    if any((z["correct_certified"],z["sound_abstentions"])!=(
      candidate["correct_certified"],candidate["sound_abstentions"]) for z in out.values()):
        raise RuntimeError("T38_CERTIFICATION_COVERAGE_DRIFT")
    return {"schema":"wingless.T38.TRAIN-only-native-comparison.v1",
      "TRAIN_meta_cases":40,"TRAIN_validation_cases":16,"DEV_opened":0,
      "arms":out,"CPU_advantage_all_controls":all(candidate["median_full_CPU_ns"]<out[a]["median_full_CPU_ns"] for a in ARMS[1:]),
      "RSS_nondegression":candidate["max_peak_RSS_bytes"]<=min(x["max_peak_RSS_bytes"] for x in out.values()),
      "scientific_positive":False,"optimizer_admission":False}
if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(proof(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-cost":
        print(json.dumps(cost(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(native(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T38_DEV_NOT_PERMITTED_UNTIL_TRAIN_ADMISSION")
