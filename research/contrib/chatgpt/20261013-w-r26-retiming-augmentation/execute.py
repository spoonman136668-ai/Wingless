"""W-R26 frozen same-capacity retiming augmentation versus clean-only learner."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r24-monotone-warp"/"execute.py"
spec=importlib.util.spec_from_file_location("w24_frozen",SRC)
assert spec and spec.loader
w24=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w24)
w22=w24.w22
w16=w22.w16

TRAIN_SEEDS=tuple(range(336,352))
HELDOUT_SEEDS=tuple(range(352,368))
TRAIN_WARPS=w24.WARPS
TEST_WARPS={
    "dup_at_tick3":(0,1,2,2,4,5,6,7,8,9,10,11,12,13,14,15),
    "dup_at_tick9":(0,1,2,3,4,5,6,7,8,8,10,11,12,13,14,15),
    "dup_ticks3_4":(0,1,2,2,2,5,6,7,8,9,10,11,12,13,14,15)
}
ARMS=("clean_only","warp_augmented")


def create_pair(seed,entity,generator,amp,phase,warp,masked,swap,negative):
    base=w22.w20.base_sequence(seed,entity,generator)
    indices=(TRAIN_WARPS[warp] if phase=="development" else TEST_WARPS[warp])
    assert all(indices[i]==i for i in (0,11,12,13,14,15))
    assert tuple(sorted(indices))==indices
    timed=tuple(base[i] for i in indices)
    observed=w16.observed(timed,seed,entity,amp,
                         17 if phase=="development" else 79)
    mask=set(w22.MASK_ORDER[:masked])
    assert not mask.intersection({0,11,12,13,14})
    candidate=tuple((0,0,0,0) if i in mask else row
                    for i,row in enumerate(observed))
    if negative:
        other=(w22.swap_tail(candidate,swap) if phase=="development"
               else w22.w19.permute(candidate,swap))
    else:
        other=candidate
    assert other[0]==candidate[0] and sorted(other)==sorted(candidate)
    assert (other==candidate)==(not negative)
    return base,other


def train(arm):
    cases=[]
    names=tuple(TRAIN_WARPS)
    for seed in TRAIN_SEEDS:
        for entity in range(9):
            for swap in (2,5,8):
                for amp in (1,2):
                    which=(seed+entity+swap+amp)%3
                    warp=names[which] if arm=="warp_augmented" else "identity"
                    masked=(0,2,4)[which]
                    for negative in (False,True):
                        a,b=create_pair(seed,entity,"affine_prior",
                                        amp,"development",warp,masked,swap,negative)
                        f=w22.maskaware_feature(a,b)
                        assert len(f)==16
                        cases.append((f,-1 if negative else 1))
    assert len(cases)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for f,label in cases:
            signed=label*(bias+sum(w*v for w,v in zip(weights,f)))
            if signed<=0:
                weights=[w+.1*label*v for w,v in zip(weights,f)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def evaluate(arm,model,generator,warp,masked):
    weights,bias,updates=model
    by_seed=[]
    for seed in HELDOUT_SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            for swap in ("swap10_11","swap12_13"):
                ap,bp=create_pair(seed,entity,generator,3,"heldout",warp,masked,swap,False)
                an,bn=create_pair(seed,entity,generator,3,"heldout",warp,masked,swap,True)
                assert ap==an and bp[0]==bn[0] and sorted(bp)==sorted(bn)
                for positive,a,b in ((True,ap,bp),(False,an,bn)):
                    f=w22.maskaware_feature(a,b)
                    pred=bias+sum(w*v for w,v in zip(weights,f))>=0
                    if positive and pred:tp+=1
                    elif not positive and not pred:tn+=1
                    elif not positive:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(r["tp"]+r["tn"] for r in by_seed)
    return {"generator":generator,"warp":warp,"masked":masked,
            "accuracy":correct/(16*36),"correct":correct,"total":16*36,
            "heldout_by_seed":by_seed}


def run():
    arms=[]
    for arm in ARMS:
        model=train(arm)
        cells=[evaluate(arm,model,generator,warp,masked)
               for generator in w22.w20.GENERATORS
               for warp in TEST_WARPS
               for masked in (3,5)]
        assert len(cells)==18
        arms.append({"arm":arm,"width":16,"training_examples":1728,
                     "epochs":32,"training_updates":model[2],
                     "heldout_accuracy":sum(c["correct"] for c in cells)/sum(c["total"] for c in cells),
                     "cells":cells})
    comparison=arms[1]["heldout_accuracy"]-arms[0]["heldout_accuracy"]
    if all(c["accuracy"]>=.80 for c in arms[1]["cells"]) and comparison>=.10:
        label="AUGMENTATION_TRANSFER"
    elif comparison>=.05:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r26-equal-capacity-retiming-augmentation.v1",
            "parent_head":"58dc556daf8b71271981d18ff9c6f594a109d6b0",
            "classification":label,"augmented_minus_control":comparison,
            "development_seeds":TRAIN_SEEDS,"sealed_heldout_seeds":HELDOUT_SEEDS,
            "arms":arms,"model_calls":0,
            "claim_limit":"Synthetic fixed-width supervised monotone warp augmentation"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
