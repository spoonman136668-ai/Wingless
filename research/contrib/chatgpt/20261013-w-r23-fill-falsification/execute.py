"""W-R23: frozen nonzero-imputation falsification of W-R22's zero-row awareness."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r22-maskaware"/"execute.py"
spec=importlib.util.spec_from_file_location("w22_frozen",SRC)
assert spec and spec.loader
w22=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w22)
w20=w22.w20
w19=w22.w19
w16=w22.w16

HELDOUT=tuple(range(272,288))
MASK_ORDER=(1,4,7,9,15,6)
FILLS=("all_four_zero","last_observed_four_channel_row","all_four_127")
GENERATORS=("affine_prior","quadratic_irregular","xorshift_independent")
SWAPS=("swap10_11","swap12_13")


def fill_candidate(original,missing,policy):
    if policy=="all_four_zero":
        return tuple((0,0,0,0) if index in missing else row
                     for index,row in enumerate(original))
    if policy=="all_four_127":
        return tuple((127,127,127,127) if index in missing else row
                     for index,row in enumerate(original))
    assert policy=="last_observed_four_channel_row"
    updated=[]
    for index,row in enumerate(original):
        if index in missing:
            assert updated
            updated.append(updated[-1])
        else:
            updated.append(row)
    return tuple(updated)


def pair(seed,entity,generator,count,fill,swap,positive):
    base=w20.base_sequence(seed,entity,generator)
    observed=w16.observed(base,seed,entity,3,79)
    missing=set(MASK_ORDER[:count])
    assert count in (2,4,6)
    assert not missing.intersection({0,11,12,13,14})
    candidate=fill_candidate(observed,missing,fill)
    counterpart=candidate if positive else w19.permute(candidate,swap)
    assert counterpart[0]==candidate[0]
    assert sorted(counterpart)==sorted(candidate)
    assert (counterpart==candidate)==positive
    return base,counterpart


def evaluate(kind,trained,fill,generator,count,swap):
    weights,bias,_=trained
    seeds=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            pbase,positive=pair(seed,entity,generator,count,fill,swap,True)
            nbase,negative=pair(seed,entity,generator,count,fill,swap,False)
            assert pbase==nbase and positive[0]==negative[0]
            assert sorted(positive)==sorted(negative)
            for groundtruth,base,compared in ((True,pbase,positive),
                                              (False,nbase,negative)):
                vec=w22.featurize(base,compared,kind)
                prediction=bias+sum(w*x for w,x in zip(weights,vec))>=0
                if groundtruth and prediction:tp+=1
                elif not groundtruth and not prediction:tn+=1
                elif not groundtruth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(x["tp"]+x["tn"] for x in seeds)
    return {"fill":fill,"generator":generator,"masked_ticks":count,
            "swap":swap,"correct":correct,"total":288,"accuracy":correct/288,
            "by_seed":seeds}


def run():
    arms=[]
    for kind in w22.ARMS:
        model=w22.train(kind)
        cells=[evaluate(kind,model,fill,gen,count,swap)
               for fill in FILLS for gen in GENERATORS
               for count in (2,4,6) for swap in SWAPS]
        assert len(cells)==54
        scores={fill:sum(c["correct"] for c in cells if c["fill"]==fill)/
                     sum(c["total"] for c in cells if c["fill"]==fill)
                for fill in FILLS}
        arms.append({"arm":kind,"width":len(model[0]),"epochs":32,
                     "training_examples":1728,"training_updates":model[2],
                     "accuracy_by_fill":scores,
                     "aggregate_accuracy":sum(c["correct"] for c in cells)/
                                          sum(c["total"] for c in cells),
                     "cells":cells})
    selected=next(a for a in arms if a["arm"]=="maskaware_augmented")
    s=selected["accuracy_by_fill"]
    if all(s[fill]>=.8 for fill in FILLS):
        label="ROBUST_IMPUTATION"
    elif s["all_four_zero"]>=.9 and all(s[fill]<=.7 for fill in FILLS[1:]):
        label="SENTINEL_DEPENDENT"
    else:
        label="MIXED"
    return {"schema":"wingless.w-r23-imputation-falsification.v1",
            "parent_head":"c551d41a7b8d4f77b9349033fa788d621a800519",
            "classification":label,"training_seeds":w22.TRAIN_SEEDS,
            "heldout_seeds":HELDOUT,"arms":arms,"model_calls":0,
            "claim_limit":"No new sensor or semantic reasoning, deterministic missing-row study"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
