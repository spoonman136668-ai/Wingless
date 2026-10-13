"""W-R24: zero-retraining missing-value format falsification."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r23-observability"/"execute.py"
spec=importlib.util.spec_from_file_location("w23_frozen",SRC)
assert spec and spec.loader
w23=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w23)
w22=w23.w22
w20=w22.w20
w16=w22.w16
w19=w22.w19

SEEDS=tuple(range(288,304))
GENERATORS=w20.GENERATORS
ARMS=w22.ARMS
FORMATS=("all_zero","all_255","carry_forward")
MASK_ORDER=(2,5,8,10,3,6)
SWAPS=("swap10_11","swap12_13")


def example(seed,entity,generator,mask_count,format_name,swap):
    base=w20.base_sequence(seed,entity,generator)
    original=w16.observed(base,seed,entity,3,79)
    mask=set(MASK_ORDER[:mask_count])
    assert mask_count in (3,5) and not mask.intersection({0,11,12,13,14})
    if format_name=="all_zero":
        filler=lambda t:(0,0,0,0)
    elif format_name=="all_255":
        filler=lambda t:(255,255,255,255)
    else:
        assert format_name=="carry_forward"
        filler=lambda t:original[t-1]
    compared=tuple(filler(t) if t in mask else row
                   for t,row in enumerate(original))
    changed=w19.permute(compared,swap)
    assert compared[0]==changed[0] and sorted(compared)==sorted(changed)
    assert compared!=changed
    return base,compared,changed


def evaluate(arm,trained,generator,format_name,mask_count):
    weights,bias,updates=trained
    by_seed=[]
    for seed in SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            for swap in SWAPS:
                base,positive,negative=example(seed,entity,generator,mask_count,
                                               format_name,swap)
                for truth,candidate in ((True,positive),(False,negative)):
                    vec=w22.featurize(base,candidate,arm)
                    assert len(vec)==16
                    prediction=bias+sum(w*v for w,v in zip(weights,vec))>=0
                    if truth and prediction:tp+=1
                    elif not truth and not prediction:tn+=1
                    elif not truth:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(x["tp"]+x["tn"] for x in by_seed)
    return {"generator":generator,"format":format_name,"masked_ticks":mask_count,
            "correct":correct,"total":16*36,"accuracy":correct/(16*36),
            "by_seed":by_seed}


def run():
    arms=[]
    for arm in ARMS:
        trained=w22.train(arm)
        cells=[evaluate(arm,trained,g,f,n)
               for g in GENERATORS for f in FORMATS for n in (3,5)]
        assert len(cells)==18
        arms.append({"arm":arm,"width":16,"training_examples":1728,
                     "epochs":32,"training_updates":trained[2],
                     "aggregate_accuracy":sum(c["correct"] for c in cells)/(18*576),
                     "cells":cells})
    by={arm["arm"]:arm for arm in arms}
    informed=by["maskaware_augmented"]
    naive=by["naive_augmented"]
    zeros=[c for c in informed["cells"] if c["format"]=="all_zero"]
    alternatives=[c for c in informed["cells"] if c["format"]!="all_zero"]
    allgood=all(c["accuracy"]>=.85 for c in informed["cells"])
    gain=informed["aggregate_accuracy"]-naive["aggregate_accuracy"]
    if allgood and gain>=.05:
        label="SENTINEL_INVARIANCE"
    elif all(c["accuracy"]>=.85 for c in zeros) and any(c["accuracy"]<.75 for c in alternatives):
        label="FORMAT_SPECIFIC"
    elif all(c["accuracy"]>=.85 for c in zeros):
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r24-missing-format.v1",
            "parent_head":"c572851852e4ae6c2eb5ba8d028a47cdc5fcb903",
            "classification":label,"gain_over_naive":gain,
            "training_seeds":w22.TRAIN_SEEDS,"sealed_heldout_seeds":SEEDS,
            "arms":arms,"model_calls":0,
            "claim_limit":"Synthetic input sentinel transfer, no semantic cognition"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
