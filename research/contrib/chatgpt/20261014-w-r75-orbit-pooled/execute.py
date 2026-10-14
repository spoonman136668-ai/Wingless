"""W-R75: frozen learned ridge readouts on circular-orbit pooled temporal features.

Preserve relative channel alignment, eliminate unknown shared circular origin.
The orbit encoder is fixed; only the 16-dimensional ridge readout is trained.
Two finite compute budgets (4 vs 16 orbit evaluations) are compared with the
same example count and an unchanged native-feature baseline.
"""
import hashlib,importlib.util,json
from pathlib import Path
parent=Path(__file__).resolve().parent.parent/"20261014-w-r73-canonical-temporal-state"/"execute.py"
spec=importlib.util.spec_from_file_location("w73_pinned_orbit",parent)
assert spec and spec.loader
w73=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w73)
m=w73.m
PRE=w73.PRE
TRAIN=w73.TRAIN
VALID=w73.VALID
SEALED=tuple(range(18240,18256))
MODES=("baseline","orbit4","orbit16")
PHASES=("intact","plus1","minus1")
AMPS=w73.AMPS
PERMS=w73.PERMS
ENTITIES=w73.ENTITIES
SITES=w73.SITES
EXCLUDED=w73.EXCLUDED
def phase_roll(rows,offset):
    assert len(rows)==16
    return tuple(rows[(i+offset)%16] for i in range(16))
def feature_pair(a,b,c,mode):
    if mode=="baseline":
        p,n=m.w42.feature(a,b),m.w42.feature(a,c)
        assert len(p)==len(n)==16
        return tuple(x-y for x,y in zip(p,n))
    offsets=(0,4,8,12) if mode=="orbit4" else tuple(range(16))
    assert mode in MODES
    total=[0.]*16
    for offset in offsets:
        x,y,z=(phase_roll(arr,offset) for arr in (a,b,c))
        positive,negative=m.w42.feature(x,y),m.w42.feature(x,z)
        for i in range(16):total[i]+=positive[i]-negative[i]
    return tuple(value/len(offsets) for value in total)
def fit_mode(mode):
    items=[]
    counts={g:0 for g in m.GENERATORS}
    for seed in TRAIN:
        for gen_index,g in enumerate(m.GENERATORS):
            if g==EXCLUDED:continue
            for entity in range(9):
                if entity%2!=(seed+gen_index)%2:continue
                for site in m.TRAIN_SITES:
                    for amp in (1,2):
                        a,b,c=m.w37.recoded_example(seed,entity,g,"identity",site,amp,True)
                        d=feature_pair(a,b,c,mode)
                        items.extend(((d,1),(tuple(-v for v in d),-1)))
                        counts[g]+=2
    assert len(items)==1728 and counts[EXCLUDED]==0
    assert all(counts[g]==864 for g in m.GENERATORS if g!=EXCLUDED)
    if mode=="baseline":
        existing,baseline_counts,_=m.build_examples(TRAIN,EXCLUDED)
        assert items==existing and counts==baseline_counts
    weights=m.w53.ridge_fit(items)
    assert len(weights)==16
    return {"weights":weights,"examples":len(items),
            "weight_sha256":hashlib.sha256(json.dumps(weights,separators=(",",":")).encode()).hexdigest()}
def score(seed,entity,site,amp,weights,mode,phase,order,source):
    a,b,c=m.w37.recoded_example(seed,entity,source,"identity",site,amp,False)
    a,b,c=(PRE.acquire(arr,order,phase) for arr in (a,b,c))
    d=feature_pair(a,b,c,mode)
    margin=sum(x*y for x,y in zip(weights,d))
    return m.w45.credit(margin,True)+m.w45.credit(-margin,False)
def fit_and_select():
    trained={mode:fit_mode(mode) for mode in MODES}
    by_mode={}
    generators=tuple(g for g in m.GENERATORS if g!=EXCLUDED)
    for mode in MODES:
        c={phase:{"credit":0.0,"n":0} for phase in PHASES}
        weights=trained[mode]["weights"]
        for seed in VALID:
            for generator in generators:
                for entity in range(9):
                    for site in SITES:
                        for phase in PHASES:
                            value=score(seed,entity,site,2,weights,mode,phase,(0,1,2,3),generator)
                            c[phase]["credit"]+=value
                            c[phase]["n"]+=2
        by_mode[mode]={phase:v["credit"]/v["n"] for phase,v in c.items()}
    selected=sorted(MODES[1:],key=lambda mode:
                    (-min(by_mode[mode].values()),-by_mode[mode]["intact"],mode))[0]
    return trained,selected,by_mode
def run():
    trained,chosen,dev=fit_and_select()
    totals={mode:{phase:{"credit":0.,"n":0} for phase in PHASES}
            for mode in ("baseline",chosen)}
    records=[]
    for mode in ("baseline",chosen):
        weights=trained[mode]["weights"]
        for amp in AMPS:
            for pname,permutation in PERMS.items():
                for phase in PHASES:
                    for seed in SEALED:
                        total=sum(score(seed,entity,site,amp,weights,mode,phase,
                                        permutation,EXCLUDED)
                                  for entity in ENTITIES for site in SITES)
                        n=2*len(ENTITIES)*len(SITES)
                        records.append({"mode":mode,"amplitude":amp,"permutation":pname,
                                        "phase":phase,"seed":seed,"credit":total,"presentations":n})
                        totals[mode][phase]["credit"]+=total
                        totals[mode][phase]["n"]+=n
    assert len(records)==2*len(AMPS)*len(PERMS)*len(PHASES)*len(SEALED)==1152
    accuracy={mode:{phase:totals[mode][phase]["credit"]/totals[mode][phase]["n"]
                    for phase in PHASES} for mode in totals}
    delta={phase:accuracy[chosen][phase]-accuracy["baseline"][phase] for phase in PHASES}
    improved=(min(accuracy[chosen].values())>=.60 and delta["plus1"]>=.07
              and delta["intact"]>=-.05 and delta["minus1"]>=-.05)
    partial=(min(accuracy[chosen].values())>=.55 and delta["plus1"]>=.04
             and delta["intact"]>=-.08 and delta["minus1"]>=-.08)
    verdict="ORBIT_POOLED_READOUT_TRANSFER" if improved else (
             "PARTIAL_ORBIT_POOLED_READOUT_GAIN" if partial else "NO_ORBIT_POOLED_GAIN")
    return {"schema":"wingless.w-r75-orbit-pooled-state.v1",
            "parent_head":"33ad17eeafed3a7941f0b1258bc67913be247aea",
            "classification":verdict,"development_seeds":TRAIN,
            "validation_seeds":VALID,"sealed_heldout_seeds":SEALED,
            "selected_mode":chosen,"development_scores":dev,
            "candidate_models":{mode:{"training_examples":obj["examples"],"weight_sha256":obj["weight_sha256"],
                                      "orbit_evaluations":1 if mode=="baseline" else (4 if mode=="orbit4" else 16)}
                                for mode,obj in trained.items()},
            "readout_dimension":16,"source_excluded":EXCLUDED,
            "accuracy":accuracy,"selected_minus_baseline":delta,
            "complete_heldout_records":records,"model_calls":0,
            "accepted_ref_mutation":False,
            "claim_limit":"Deterministic circular orbit pooling plus supervised fixed-width ridge; 4× or 16× feature extraction compute. Not learned state-former, RSI or fair equal-FLOP comparison."}
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
