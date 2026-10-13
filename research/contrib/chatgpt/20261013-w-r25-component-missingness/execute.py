"""W-R25: component-zero masking vs row-zero masking, same 16-wide supervised budget."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r24-channel-dropout"/"execute.py"
_spec=importlib.util.spec_from_file_location("w24_frozen",SOURCE)
assert _spec and _spec.loader
w24=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w24)
w22=w24.w22
w16=w22.w16
w20=w22.w20

DEV=tuple(range(304,320))
HELDOUT=tuple(range(320,336))
TRAIN_SWAPS=(2,5,8)
TEST_SWAPS=("swap10_11","swap12_13")
TEST_CASES=("none","irrelevant_full_rows_2_5","single_channel_zero_2_5",
            "four_channels_separate_ticks_2_5_8_10")
TRAIN_CASES=("none","single_channel_zero_2_5","four_channels_separate_ticks_2_5_8_10")
ARMS=("rowaware_clean","rowaware_augmented","componentzero_clean","componentzero_augmented")


def swapped(rows,transform):
    if isinstance(transform,int):
        index=transform
    else:
        index={"swap10_11":10,"swap12_13":12}[transform]
    tail=list(rows[1:])
    tail[index],tail[index+1]=tail[index+1],tail[index]
    result=rows[:1]+tuple(tail)
    assert result!=rows and result[0]==rows[0]
    assert sorted(result)==sorted(rows)
    return result


def pair(seed,entity,amp,generator,transform,case,positive,train):
    original=w20.base_sequence(seed,entity,generator)
    measured=w16.observed(original,seed,entity,amp,17 if train else 79)
    altered=measured if positive else swapped(measured,transform)
    compared=w24.drop(altered,case,transform if not train else "swap10_11")
    # Drop positions are absolute timeline positions, not privileged label inputs.
    assert len(compared)==len(original)==16
    return original,compared


def features(original,candidate,arm):
    if arm.startswith("rowaware"):
        return w22.maskaware_feature(original,candidate)
    assert arm.startswith("componentzero")
    diffs=[]
    for channel in range(4):
        for tick in range(16):
            next_tick=(tick+1)%16
            # Candidate value 0 is only an observed symbol, not an externally supplied mask.
            if candidate[tick][channel]==0 or candidate[next_tick][channel]==0:
                continue
            ref_delta=(original[next_tick][channel]-original[tick][channel])%256
            seen_delta=(candidate[next_tick][channel]-candidate[tick][channel])%256
            diffs.append(w16.circular_distance(ref_delta,seen_delta))
    assert diffs
    top=sorted(diffs,reverse=True)[:2]
    mean=sum(top)/len(top)
    return [mean]*16


def train(arm):
    samples=[]
    for seed in DEV:
        for entity in range(9):
            for ix in TRAIN_SWAPS:
                for amp in (1,2):
                    case=TRAIN_CASES[(seed+entity+ix+amp)%3] if arm.endswith("augmented") else "none"
                    for pos in (True,False):
                        original,candidate=pair(seed,entity,amp,"affine_prior",ix,case,pos,True)
                        samples.append((features(original,candidate,arm),1 if pos else -1))
    assert len(samples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vector,label in samples:
            score=bias+sum(w*x for w,x in zip(weights,vector))
            if label*score<=0:
                weights=[w+.1*label*x for w,x in zip(weights,vector)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def evaluate(arm,trained,generator,transform,case):
    weights,bias,updates=trained
    seeds=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            for pos in (True,False):
                original,candidate=pair(seed,entity,3,generator,transform,case,pos,False)
                vector=features(original,candidate,arm)
                predicted=bias+sum(w*x for w,x in zip(weights,vector))>=0
                if pos and predicted:tp+=1
                elif not pos and not predicted:tn+=1
                elif not pos:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    good=sum(s["tp"]+s["tn"] for s in seeds)
    return {"generator":generator,"swap":transform,"case":case,
            "accuracy":good/(16*18),"correct":good,"total":16*18,"seed_rows":seeds}


def run():
    models=[]
    for arm in ARMS:
        trained=train(arm)
        cells=[evaluate(arm,trained,g,swap,case)
               for g in w20.GENERATORS for swap in TEST_SWAPS for case in TEST_CASES]
        assert len(cells)==24
        partial=[c for c in cells if c["case"] in TRAIN_CASES[1:]]
        models.append({"arm":arm,"feature_width":16,"training_pairs":1728,"epochs":32,
                       "updates":trained[2],
                       "overall_accuracy":sum(c["correct"] for c in cells)/sum(c["total"] for c in cells),
                       "partial_accuracy":sum(c["correct"] for c in partial)/sum(c["total"] for c in partial),
                       "cells":cells})
    by={m["arm"]:m for m in models}
    a=by["componentzero_augmented"]
    b=by["rowaware_augmented"]
    gain=a["partial_accuracy"]-b["partial_accuracy"]
    cells=[c for c in a["cells"] if c["case"] in TRAIN_CASES[1:]]
    if all(c["accuracy"]>=.80 for c in cells) and gain>=.10:
        label="CHANNEL_ZERO_TRANSFER"
    elif gain>=.05:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r25-component-observation.v1",
            "parent_head":"800b2fc679e584d50c51a6f0bf5a28a9118dcd2d",
            "classification":label,"partial_gain_component_vs_row":gain,
            "development_seeds":DEV,"sealed_heldout_seeds":HELDOUT,
            "models":models,"model_calls":0,
            "claim_limit":"Finite synthetic zero-symbol missingness mechanism"}
    
    
if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
