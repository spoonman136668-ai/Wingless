"""W-R33: bounded equal-budget quantized temporal prediction retraining."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r32-quantized-marker"/"execute.py"
spec=importlib.util.spec_from_file_location("w32_pinned", SOURCE)
assert spec and spec.loader
w32=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w32)
w31=w32.w31
w25=w31.w25
w22=w25.w22
w20=w22.w20
w16=w25.w16

TRAIN=tuple(range(1136,1152))
HELDOUT=tuple(range(1152,1168))
BITS=(8,4,2,1)
ARMS=("frozen_maskaware","refit_multimarker","refit_allobserved")
GENERATORS=w32.GENERATORS
SITES=w32.SITES
INDICES=w22.TRAIN_SWAP_INDICES


def feature(base, candidate, arm):
    if arm=="frozen_maskaware":
        result=w22.maskaware_feature(base,candidate)
    elif arm=="refit_multimarker":
        result=w25.new_feature(base,candidate)
    else:
        assert arm=="refit_allobserved"
        result=w31.representation(base,candidate,"all_observed")
    assert len(result)==16
    return result


def train(arm,bits):
    assert arm!="frozen_maskaware"
    examples=[]
    identical_pairs=0
    for seed in TRAIN:
        for entity in range(9):
            base=w20.base_sequence(seed,entity,"affine_prior")
            for site in INDICES:
                for amplitude in (1,2):
                    original=w16.observed(base,seed,entity,amplitude,17)
                    positive=tuple(w32.quantize(row,bits) for row in original)
                    tail=list(positive[1:])
                    tail[site],tail[site+1]=tail[site+1],tail[site]
                    negative=positive[:1]+tuple(tail)
                    assert positive[0]==negative[0]
                    assert sorted(positive)==sorted(negative)
                    identical_pairs+=int(positive==negative)
                    examples.append((feature(base,positive,arm),1))
                    examples.append((feature(base,negative,arm),-1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for epoch in range(32):
        for v,label in examples:
            score=bias+sum(w*x for w,x in zip(weights,v))
            if label*score<=0:
                weights=[w+0.1*label*x for w,x in zip(weights,v)]
                bias+=0.1*label
                updates+=1
    return (weights,bias,updates,identical_pairs)


def evaluate(arm,model,generator,bits,site):
    weights,bias,updates,identical_training=model
    by_seed=[]
    correct=total=identical=0
    for seed in HELDOUT:
        tp=tn=fp=fn=same=0
        for entity in range(9):
            base,positive,negative=w32.observed_pair(seed,entity,generator,bits,site)
            same+=int(positive==negative)
            for is_positive,candidate in ((True,positive),(False,negative)):
                v=feature(base,candidate,arm)
                guess=bias+sum(w*x for w,x in zip(weights,v))>=0
                if is_positive and guess: tp+=1
                elif not is_positive and not guess: tn+=1
                elif not is_positive: fp+=1
                else: fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        identical+=same
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "identical_counterfactual_pairs":same})
    return {"generator":generator,"bits":bits,"swap_site":site,
            "correct":correct,"total":total,"accuracy":correct/total,
            "identical_heldout_pairs":identical,"by_seed":by_seed}


def run():
    frozen=w22.train("maskaware_augmented")
    assert len(frozen[0])==16
    results=[]
    trained_models={}
    for bits in BITS:
        for arm in ARMS:
            model=((frozen[0],frozen[1],frozen[2],None)
                   if arm=="frozen_maskaware" else train(arm,bits))
            trained_models[(arm,bits)]=model
            cells=[evaluate(arm,model,generator,bits,site)
                   for generator in GENERATORS for site in SITES]
            assert len(cells)==6
            results.append({"bits":bits,"arm":arm,"width":16,
                            "training_examples":1728,
                            "training_population":("W-R22 historical original 240..255"
                                                   if arm=="frozen_maskaware" else "new development 1136..1151"),
                            "epochs":32,"training_updates":model[2],
                            "training_identical_pairs":model[3],
                            "aggregate_accuracy":sum(c["correct"] for c in cells)/
                                                 sum(c["total"] for c in cells),
                            "cells":cells})
    by={(r["bits"],r["arm"]):r for r in results}
    gains={}
    for bits in (4,2):
        prior=by[(bits,"frozen_maskaware")]["aggregate_accuracy"]
        gains[str(bits)]={arm:by[(bits,arm)]["aggregate_accuracy"]-prior
                          for arm in ("refit_multimarker","refit_allobserved")}
    peak=max(v for scores in gains.values() for v in scores.values())
    if peak>=0.10: classification="RECOVERY"
    elif peak>=0.05: classification="PARTIAL_RECOVERY"
    else: classification="NO_RECOVERY"
    return {"schema":"wingless.w-r33-equal-budget-quantized-refit.v1",
            "parent_head":"78613735739aa99fa30b5093b9ff55389e177922",
            "classification":classification,
            "development_seeds":TRAIN,"sealed_heldout_seeds":HELDOUT,
            "quantizer_bits":BITS,"heldout_gains_relative_to_historical_frozen":gains,
            "results":results,"model_calls":0,
            "claim_limit":"Synthetic clean quantized pairs under equal-budget new refits; fixed original weights are a historical benchmark, not new-training matched."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
