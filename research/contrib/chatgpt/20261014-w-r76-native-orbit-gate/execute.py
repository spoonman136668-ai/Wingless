"""W-R76: development-selected learnable mixture of native and orbit pooled state.

The single mixing coefficient is trained by discrete bounded development
optimization; every model has exactly sixteen ridge weights and the same
1728 signed training examples. Intermediate candidates expend up to
17 feature extractions per example. No heldout feedback enters the choice.
"""
import hashlib,importlib.util,json
from pathlib import Path
parent=Path(__file__).resolve().parent.parent/"20261014-w-r75-orbit-pooled"/"execute.py"
spec=importlib.util.spec_from_file_location("w75_frozen",parent)
assert spec and spec.loader
w75=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w75)
m=w75.m
PRE=w75.PRE
TRAIN=w75.TRAIN
VALID=w75.VALID
SEALED=tuple(range(19040,19056))
EXCLUDED=w75.EXCLUDED
PHASES=w75.PHASES
AMPS=w75.AMPS
PERMS=w75.PERMS
SITES=w75.SITES
ENTITIES=w75.ENTITIES
ALPHAS=(0.0,0.25,0.50,0.75,1.0)
def mixed_features(a,b,c,alpha):
    baseline=w75.feature_pair(a,b,c,"baseline")
    if alpha==0.0:return baseline
    invariant=w75.feature_pair(a,b,c,"orbit16")
    if alpha==1.0:return invariant
    return tuple((1-alpha)*u+alpha*v for u,v in zip(baseline,invariant))
def fit_alpha(alpha):
    samples=[]
    counts={g:0 for g in m.GENERATORS}
    for seed in TRAIN:
        for idx,g in enumerate(m.GENERATORS):
            if g==EXCLUDED:continue
            for entity in range(9):
                if entity%2!=(seed+idx)%2:continue
                for site in m.TRAIN_SITES:
                    for amplitude in (1,2):
                        a,b,c=m.w37.recoded_example(seed,entity,g,"identity",site,amplitude,True)
                        d=mixed_features(a,b,c,alpha)
                        samples.extend(((d,1),(tuple(-v for v in d),-1)))
                        counts[g]+=2
    assert len(samples)==1728 and counts[EXCLUDED]==0
    assert all(counts[g]==864 for g in m.GENERATORS if g!=EXCLUDED)
    if alpha==0.0:
        original,origcounts,_=m.build_examples(TRAIN,EXCLUDED)
        assert original==samples and counts==origcounts
    weights=m.w53.ridge_fit(samples)
    assert len(weights)==16
    return {"weights":weights,"examples":len(samples),
            "sha256":hashlib.sha256(json.dumps(weights).encode()).hexdigest()}
def credit(seed,entity,site,amp,weights,alpha,phase,order,source):
    a,b,c=m.w37.recoded_example(seed,entity,source,"identity",site,amp,False)
    a,b,c=(PRE.acquire(rows,order,phase) for rows in (a,b,c))
    delta=mixed_features(a,b,c,alpha)
    margin=sum(w*v for w,v in zip(weights,delta))
    return m.w45.credit(margin,True)+m.w45.credit(-margin,False)
def fit_and_select():
    fitted={alpha:fit_alpha(alpha) for alpha in ALPHAS}
    included=tuple(g for g in m.GENERATORS if g!=EXCLUDED)
    dev={}
    for alpha in ALPHAS:
        aggregate={phase:{"score":0.0,"n":0} for phase in PHASES}
        weights=fitted[alpha]["weights"]
        for seed in VALID:
            for generator in included:
                for entity in range(9):
                    for site in SITES:
                        for phase in PHASES:
                            val=credit(seed,entity,site,2,weights,alpha,phase,
                                       (0,1,2,3),generator)
                            aggregate[phase]["score"]+=val
                            aggregate[phase]["n"]+=2
        dev[str(alpha)]={p:v["score"]/v["n"] for p,v in aggregate.items()}
    chosen=min(ALPHAS[1:],key=lambda alpha:
          (-min(dev[str(alpha)].values()),-dev[str(alpha)]["intact"],alpha))
    return fitted,chosen,dev
def run():
    fitted,chosen,development=fit_and_select()
    scores={str(alpha):{phase:{"score":0.0,"n":0} for phase in PHASES}
            for alpha in (0.0,chosen)}
    records=[]
    for alpha in (0.0,chosen):
        weights=fitted[alpha]["weights"]
        for amp in AMPS:
            for name,order in PERMS.items():
                for phase in PHASES:
                    for seed in SEALED:
                        total=sum(credit(seed,entity,site,amp,weights,alpha,phase,
                                         order,EXCLUDED)
                                  for entity in ENTITIES for site in SITES)
                        n=2*len(ENTITIES)*len(SITES)
                        records.append({"alpha":str(alpha),"amplitude":amp,
                                        "permutation":name,"phase":phase,"seed":seed,
                                        "credit":total,"presentations":n})
                        scores[str(alpha)][phase]["score"]+=total
                        scores[str(alpha)][phase]["n"]+=n
    assert len(records)==2*3*len(PERMS)*len(PHASES)*len(SEALED)==1152
    accuracy={key:{phase:x["score"]/x["n"] for phase,x in phase_group.items()}
              for key,phase_group in scores.items()}
    delta={phase:accuracy[str(chosen)][phase]-accuracy["0.0"][phase]
           for phase in PHASES}
    robust=(min(accuracy[str(chosen)].values())>=.62 and
            delta["intact"]>=-.08 and delta["minus1"]>=-.08 and delta["plus1"]>=.08)
    partial=(min(accuracy[str(chosen)].values())>=.57 and
             delta["intact"]>=-.12 and delta["minus1"]>=-.12 and delta["plus1"]>=.04)
    verdict="LEARNED_NATIVE_ORBIT_MIX_TRANSFER" if robust else (
             "PARTIAL_NATIVE_ORBIT_COMPROMISE" if partial else "NO_NATIVE_ORBIT_MIX_GAIN")
    return {"schema":"wingless.w-r76-native-orbit-mixture.v1",
            "parent_head":"dd8ef01f85d422bf95fe267f5b73e1f076026961",
            "classification":verdict,"development_seeds":TRAIN,
            "validation_seeds":VALID,"sealed_heldout_seeds":SEALED,
            "alpha_candidates":ALPHAS,"selected_alpha":chosen,
            "development_scores":development,
            "model_hashes":{str(k):v["sha256"] for k,v in fitted.items()},
            "per_model_training_examples":1728,"ridge_weight_dimension":16,
            "feature_compute_multiplier_max":17,"accuracy":accuracy,
            "chosen_minus_baseline":delta,"complete_heldout_records":records,
            "model_calls":0,"accepted_ref_mutation":False,
            "claim_limit":"Development-trained bounded scalar mixing gate + supervised fixed ridge; feature work up to17x, not learned recurrent state former, equal-FLOP or RSI."}
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
