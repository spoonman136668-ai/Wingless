"""W-R31: controlled legitimate marker collision cost in fixed temporal learner."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r30-marker-alias"/"execute.py"
spec=importlib.util.spec_from_file_location("w30_frozen",SOURCE)
assert spec and spec.loader
w30=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w30)
w25=w30.w25

HELDOUT=tuple(range(1024,1040))
GENERATORS=w30.GENERATORS
FORMATS=("ordinary_control","valid_all_zero","valid_all_255","valid_copy_previous")
SITES=(10,12)
ARMS=("old_zero","multi_marker","all_observed")


def observed_pair(seed,entity,generator,format_name,site):
    base=w25.w24.w22.w20.base_sequence(seed,entity,generator)
    raw=w25.w16.observed(base,seed,entity,3,79)
    positive=list(raw)
    tick=site+1
    assert tick in (11,13)
    if format_name=="valid_all_zero":
        positive[tick]=(0,0,0,0)
    elif format_name=="valid_all_255":
        positive[tick]=(255,255,255,255)
    elif format_name=="valid_copy_previous":
        positive[tick]=positive[tick-1]
    else:
        assert format_name=="ordinary_control"
    positive=tuple(positive)
    tail=list(positive[1:])
    tail[site],tail[site+1]=tail[site+1],tail[site]
    negative=positive[:1]+tuple(tail)
    assert positive[0]==negative[0]
    assert sorted(positive)==sorted(negative)
    assert positive!=negative
    return base,positive,negative


def representation(base,candidate,arm):
    if arm=="old_zero":
        features=w25.w22.maskaware_feature(base,candidate)
    elif arm=="multi_marker":
        features=w25.new_feature(base,candidate)
    else:
        assert arm=="all_observed"
        features=w25.w22.w20.w19.feature(base,candidate,"mean_top_two_64")
    assert len(features)==16
    return features


def classify(model,base,candidate,arm):
    weights,bias,updates=model
    v=representation(base,candidate,arm)
    return bias+sum(w*f for w,f in zip(weights,v))>=0


def evaluation(arm,model,generator,fmt,site):
    rows=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        flagged_legit=0
        for entity in range(9):
            base,positive,negative=observed_pair(seed,entity,generator,fmt,site)
            tick=site+1
            # This observation is genuinely acquired, NOT dropped. Any
            # detector mask is a false exclusion regardless of pair label.
            flagged_legit+=int(w25.marker(positive[tick],positive[tick-1]))
            for truth,candidate in ((True,positive),(False,negative)):
                predicted=classify(model,base,candidate,arm)
                if truth and predicted:tp+=1
                elif not truth and not predicted:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "positive_legitimate_marker_exclusions":flagged_legit})
    correct=sum(r["tp"]+r["tn"] for r in rows)
    return {"generator":generator,"format":fmt,"tail_swap_index":site,
            "accuracy":correct/288,"correct":correct,"total":288,
            "legitimate_markers_falsely_excluded_on_positives":
                sum(r["positive_legitimate_marker_exclusions"] for r in rows),
            "by_seed":rows}


def run():
    trained=w25.w22.train("maskaware_augmented")
    assert len(trained[0])==16
    arms=[]
    for arm in ARMS:
        cells=[evaluation(arm,trained,g,f,s)
               for g in GENERATORS for f in FORMATS for s in SITES]
        assert len(cells)==24
        arms.append({"arm":arm,"width":16,"training_examples":1728,
                     "epochs":32,"training_updates":trained[2],
                     "aggregate_accuracy":sum(c["correct"] for c in cells)/
                                          sum(c["total"] for c in cells),
                     "cells":cells})
    d={arm["arm"]:arm for arm in arms}
    multi={(c["generator"],c["format"],c["tail_swap_index"]):c
           for c in d["multi_marker"]["cells"]}
    all_rows={(c["generator"],c["format"],c["tail_swap_index"]):c
              for c in d["all_observed"]["cells"]}
    cost_rows=[multi[k]["accuracy"]-all_rows[k]["accuracy"]
               for k in multi if k[1]!="ordinary_control"]
    noncontrol_gain=sum(cost_rows)/len(cost_rows)
    if noncontrol_gain<=-.10:
        label="LEGITIMATE_MARKER_COST"
    elif any(x<=-.10 for x in cost_rows) and any(x>=.10 for x in cost_rows):
        label="TRADEOFF"
    else:
        label="NO_CONFIRMED_COST"
    return {"schema":"wingless.w-r31-legitimate-marker-downstream.v1",
            "parent_head":"8e648949b708666d153320ef5fc3f8b77cae0b37",
            "classification":label,
            "multi_marker_minus_all_observed_noncontrol":noncontrol_gain,
            "heldout_seeds":HELDOUT,"training_seeds":w25.w22.TRAIN_SEEDS,
            "arms":arms,"model_calls":0,
            "claim_limit":"Synthetic intentionally valid marker values overlapping swaps; not natural marker incidence or generalized reasoning"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
