"""W-R25: frozen dual-format observed-only mask detector, unchanged learned weights."""
import importlib.util
import json
from pathlib import Path

PREVIOUS=Path(__file__).resolve().parent.parent/"20261013-w-r24-format"/"execute.py"
_spec=importlib.util.spec_from_file_location("w24_frozen",PREVIOUS)
assert _spec and _spec.loader
w24=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w24)
w22=w24.w22
w16=w22.w16

SEEDS=tuple(range(304,320))
FORMATS=w24.FORMATS
GENERATORS=w24.GENERATORS
ARMS=("old_zero_detector","observed_multi_marker_detector")
MISSING_ZERO=(0,0,0,0)
MISSING_255=(255,255,255,255)


def marker(row,previous):
    return row==MISSING_ZERO or row==MISSING_255 or row==previous


def new_feature(base,candidate):
    vals=[]
    for channel in range(4):
        aa=[row[channel] for row in base]
        bb=[row[channel] for row in candidate]
        for i in range(16):
            nxt=(i+1)%16
            missing_i=marker(candidate[i],candidate[i-1]) if i>0 else False
            missing_n=marker(candidate[nxt],candidate[nxt-1]) if nxt>0 else False
            if missing_i or missing_n:
                continue
            da=(aa[nxt]-aa[i])%256
            db=(bb[nxt]-bb[i])%256
            vals.append(w16.circular_distance(da,db))
    if len(vals)>=2:
        scalar=sum(sorted(vals,reverse=True)[:2])/2
    elif vals:
        scalar=vals[0]
    else:
        scalar=0.0
    return [scalar]*16


def check_frozen_development_feature_equivalence():
    count=0
    for seed in w22.TRAIN_SEEDS:
        for entity in range(9):
            for index in w22.TRAIN_SWAP_INDICES:
                for amp in (1,2):
                    n=(0,2,4)[(seed+entity+index+amp)%3]
                    for negative in (False,True):
                        a,b=w22.prepared(seed,entity,"affine_prior",amp,
                                         "development",n,negative,index)
                        assert new_feature(a,b)==w22.maskaware_feature(a,b)
                        count+=1
    assert count==1728
    return count


def predict(trained,base,observed,arm):
    weights,bias,updates=trained
    vec=(w22.maskaware_feature(base,observed)
         if arm=="old_zero_detector" else new_feature(base,observed))
    assert len(vec)==16
    return bias+sum(w*v for w,v in zip(weights,vec))>=0


def evaluate(arm,trained,generator,format_name,n):
    seeds=[]
    for seed in SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            for swap in w24.SWAPS:
                base,positive,negative=w24.example(seed,entity,generator,n,format_name,swap)
                assert positive[0]==negative[0] and sorted(positive)==sorted(negative)
                for truth,observed in ((True,positive),(False,negative)):
                    pred=predict(trained,base,observed,arm)
                    if truth and pred:tp+=1
                    elif not truth and not pred:tn+=1
                    elif not truth:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(s["tp"]+s["tn"] for s in seeds)
    return {"generator":generator,"format":format_name,"masked_ticks":n,
            "correct":correct,"total":576,"accuracy":correct/576,
            "by_seed":seeds}


def run():
    matching_development=check_frozen_development_feature_equivalence()
    trained=w22.train("maskaware_augmented")
    rows=[]
    for arm in ARMS:
        cells=[evaluate(arm,trained,g,f,n)
               for g in GENERATORS for f in FORMATS for n in (3,5)]
        assert len(cells)==18
        rows.append({"arm":arm,"training_examples":1728,
                     "feature_width":16,"epochs":32,"training_updates":trained[2],
                     "aggregate_accuracy":sum(x["correct"] for x in cells)/(18*576),
                     "cells":cells})
    old,new=rows
    gain=new["aggregate_accuracy"]-old["aggregate_accuracy"]
    if all(c["accuracy"]>=.8 for c in new["cells"]) and gain>=.1:
        classification="REPRESENTATION_BENEFIT"
    elif gain>=.05:
        classification="PARTIAL"
    else:
        classification="NO_GAIN"
    return {"schema":"wingless.w-r25-multi-observed-marker.v1",
            "parent_head":"2129317f1ad3a1247413991c68ae8c8e10da0dcd",
            "classification":classification,"aggregate_gain":gain,
            "matched_frozen_training_examples":matching_development,
            "training_seeds":w22.TRAIN_SEEDS,"sealed_heldout_seeds":SEEDS,
            "arms":rows,"model_calls":0,
            "claim_limit":"Hand-coded observed-data format detector and unchanged learned weights, not general reasoning"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
