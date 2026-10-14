"""W-R73: development-selected bounded temporal-state canonicalizers vs frozen ridge."""
import hashlib,importlib.util,json
from pathlib import Path
p=Path(__file__).resolve().parent.parent/"20261014-w-r72-three-phase"/"execute.py"
spec=importlib.util.spec_from_file_location("w72_science",p)
assert spec and spec.loader
w72=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w72)
m=w72.m
PRE=w72.w67
TRAIN=tuple(range(9008,9024))
VALID=tuple(range(9040,9048))+tuple(range(9072,9080))
SEALED=tuple(range(16000,16016))
MODES=("baseline","channel0_peak","channel0_gradient","four_channel_energy")
PHASES=("intact","plus1","minus1")
AMPS=(1,2,3)
PERMS=w72.PERMS
ENTITIES=w72.ENTITIES
SITES=w72.SITES
EXCLUDED="xorshift_independent"
assert EXCLUDED in m.GENERATORS
def normalize(rows,mode):
    if mode=="baseline":return rows
    assert len(rows)==16 and len(rows[0])==4
    if mode=="channel0_peak":
        values=[abs(rows[t][0]) for t in range(16)]
    elif mode=="channel0_gradient":
        values=[abs(rows[t][0]-rows[(t-1)%16][0]) for t in range(16)]
    elif mode=="four_channel_energy":
        values=[sum(abs(v) for v in rows[t]) for t in range(16)]
    else:raise ValueError("UNKNOWN_MODE")
    anchor=max(range(16),key=lambda i:(values[i],-i))
    return tuple(rows[(i+anchor)%16] for i in range(16))
def feature_pair(base,pos,neg,mode):
    a,b,c=(normalize(x,mode) for x in (base,pos,neg))
    p,n=m.w42.feature(a,b),m.w42.feature(a,c)
    assert len(p)==len(n)==16
    return tuple(x-y for x,y in zip(p,n))
def fit_mode(mode):
    ex=[]
    counts={g:0 for g in m.GENERATORS}
    for seed in TRAIN:
        for idx,g in enumerate(m.GENERATORS):
            if g==EXCLUDED:continue
            for entity in range(9):
                if entity%2!=(seed+idx)%2:continue
                for site in m.TRAIN_SITES:
                    for amplitude in (1,2):
                        a,b,c=m.w37.recoded_example(seed,entity,g,"identity",site,amplitude,True)
                        d=feature_pair(a,b,c,mode)
                        ex.extend(((d,1),(tuple(-v for v in d),-1)))
                        counts[g]+=2
    assert len(ex)==1728 and counts[EXCLUDED]==0
    assert all(counts[g]==864 for g in m.GENERATORS if g!=EXCLUDED)
    if mode=="baseline":
        original,c,sha=m.build_examples(TRAIN,EXCLUDED)
        assert ex==original and c==counts
    weights=m.w53.ridge_fit(ex)
    assert len(weights)==16
    return {"weights":weights,"examples":len(ex),"source_counts":counts,
            "weight_sha256":hashlib.sha256(json.dumps(weights,separators=(",",":")).encode()).hexdigest()}
def credit(seed,entity,site,amp,weights,mode,phase,order):
    a,b,c=m.w37.recoded_example(seed,entity,EXCLUDED,"identity",site,amp,False)
    a,b,c=(PRE.acquire(x,order,phase) for x in (a,b,c))
    d=feature_pair(a,b,c,mode)
    margin=sum(w*v for w,v in zip(weights,d))
    return m.w45.credit(margin,True)+m.w45.credit(-margin,False)
def validate_dev(all_models):
    by_mode={}
    included=[g for g in m.GENERATORS if g!=EXCLUDED]
    for mode in MODES:
        weights=all_models[mode]["weights"]
        results={phase:[0.,0] for phase in PHASES}
        # Validation exercises included generators only; no excluded-source preview.
        for seed in VALID:
            for generator in included:
                for entity in range(9):
                    for site in SITES:
                        amp=2
                        a,b,c=m.w37.recoded_example(seed,entity,generator,"identity",site,amp,False)
                        for phase in PHASES:
                            ar,br,cr=(PRE.acquire(x,(0,1,2,3),phase) for x in (a,b,c))
                            d=feature_pair(ar,br,cr,mode)
                            margin=sum(w*v for w,v in zip(weights,d))
                            results[phase][0]+=m.w45.credit(margin,True)+m.w45.credit(-margin,False)
                            results[phase][1]+=2
        by_mode[mode]={phase:float(v[0]/v[1]) for phase,v in results.items()}
    # Freeze a single candidate using development worst-phase score, then intact
    # score, then deterministic name. Baseline cannot be selected as candidate.
    selected=sorted(MODES[1:],key=lambda mode:
            (-min(by_mode[mode].values()),-by_mode[mode]["intact"],mode))[0]
    return selected,by_mode
def run():
    trained={mode:fit_mode(mode) for mode in MODES}
    chosen,dev_scores=validate_dev(trained)
    totals={mode:{p:{"credit":0.,"n":0} for p in PHASES}
            for mode in ("baseline",chosen)}
    rows=[]
    for mode in ("baseline",chosen):
        weights=trained[mode]["weights"]
        for amp in AMPS:
            for permutation,order in PERMS.items():
                for phase in PHASES:
                    for seed in SEALED:
                        total=sum(credit(seed,e,site,amp,weights,mode,phase,order)
                                  for e in ENTITIES for site in SITES)
                        n=2*len(ENTITIES)*len(SITES)
                        rows.append({"mode":mode,"amplitude":amp,"permutation":permutation,
                                     "phase":phase,"seed":seed,"credit":total,"presentations":n})
                        totals[mode][phase]["credit"]+=total
                        totals[mode][phase]["n"]+=n
    assert len(rows)==2*len(AMPS)*len(PERMS)*len(PHASES)*len(SEALED)==1152
    acc={mode:{p:totals[mode][p]["credit"]/totals[mode][p]["n"]
               for p in PHASES} for mode in totals}
    delta={p:acc[chosen][p]-acc["baseline"][p] for p in PHASES}
    robust=(min(acc[chosen].values())>=.60 and delta["intact"]>=-.05
            and delta["plus1"]>=.07 and delta["minus1"]>=-.05)
    partial=(min(acc[chosen].values())>=.55 and delta["intact"]>=-.08
             and delta["plus1"]>=.04)
    label="BOUNDED_CANONICAL_STATE_TRANSFER" if robust else (
          "PARTIAL_CANONICAL_STATE_TRANSFER" if partial else "NO_CANONICAL_STATE_GAIN")
    return {"schema":"wingless.w-r73-frozen-canonical-state.v1",
            "parent_head":"e127dde54ce1d4a9f992d0f201c88ab7ee263909",
            "classification":label,"development_seeds":TRAIN,
            "validation_seeds":VALID,"sealed_heldout_seeds":SEALED,
            "candidate_modes":MODES,"dev_phase_scores":dev_scores,
            "selected_mode":chosen,"frozen_model_metadata":{
                k:{z:v for z,v in d.items() if z!="weights"} for k,d in trained.items()},
            "readout_weights":16,"training_examples_per_candidate":1728,
            "extra_inference_normalizer_ops_bound_per_sequence":80,
            "accuracies":acc,"selected_minus_baseline":delta,
            "all_sealed_records":rows,"model_calls":0,"accepted_ref_mutation":False,
            "claim_limit":"Hand-designed circular raw-state canonicalizers, development-selected and ridge-fitted. Additional bounded normalization compute; not a learned state former or RSI."}
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
