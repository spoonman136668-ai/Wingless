"""W-R23: fixed-model hidden missingness encoding falsification."""
import importlib.util
import json
from pathlib import Path

PREVIOUS = Path(__file__).resolve().parent.parent / "20261013-w-r22-maskaware" / "execute.py"
_spec=importlib.util.spec_from_file_location("w22_frozen",PREVIOUS)
assert _spec and _spec.loader
w22=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w22)
w20=w22.w20
w19=w22.w19
w16=w22.w16

HELDOUT=tuple(range(272,288))
GENERATORS=w20.GENERATORS
ARMS=w22.ARMS
FILL_MODES=("explicit_all_zero","forward_fill","linear_interpolate")
OCCLUSIONS=(3,5)


def fill_missing(observed, count, mode):
    positions=set(w22.MASK_ORDER[:count])
    assert not positions.intersection({0,11,12,13,14})
    rows=list(observed)
    for index in sorted(positions):
        if mode=="explicit_all_zero":
            rows[index]=(0,0,0,0)
            continue
        prev=max(i for i in range(index) if i not in positions)
        if mode=="forward_fill":
            rows[index]=observed[prev]
        else:
            assert mode=="linear_interpolate"
            nxt=min(i for i in range(index+1,len(observed)) if i not in positions)
            rows[index]=tuple((observed[prev][c]+observed[nxt][c])//2
                              for c in range(4))
    return tuple(rows)


def pair(seed,entity,generator,count,mode,swap,positive):
    base=w20.base_sequence(seed,entity,generator)
    observed=w16.observed(base,seed,entity,3,79)
    candidate=fill_missing(observed,count,mode)
    compared=candidate if positive else w19.permute(candidate,swap)
    assert compared[0]==candidate[0] and sorted(compared)==sorted(candidate)
    assert (compared==candidate)==positive
    return base,compared


def evaluate(arm,model,generator,masked,mode):
    weights,bias,updates=model
    rows=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            for swap in ("swap10_11","swap12_13"):
                a_pos,b_pos=pair(seed,entity,generator,masked,mode,swap,True)
                a_neg,b_neg=pair(seed,entity,generator,masked,mode,swap,False)
                assert a_pos==a_neg and b_pos[0]==b_neg[0]
                assert sorted(b_pos)==sorted(b_neg)
                for positive,a,b in ((True,a_pos,b_pos),(False,a_neg,b_neg)):
                    feature=w22.featurize(a,b,arm)
                    assert len(feature)==16
                    predicted=bias+sum(w*x for w,x in zip(weights,feature))>=0
                    if positive and predicted:tp+=1
                    elif not positive and not predicted:tn+=1
                    elif not positive:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(s["tp"]+s["tn"] for s in rows)
    return {"generator":generator,"masked_ticks":masked,"encoding":mode,
            "accuracy":correct/(16*36),"correct":correct,"total":16*36,
            "heldout_by_seed":rows}


def run():
    arms=[]
    for arm in ARMS:
        trained=w22.train(arm)
        cells=[evaluate(arm,trained,gen,masked,encoding)
               for gen in GENERATORS for masked in OCCLUSIONS
               for encoding in FILL_MODES]
        assert len(cells)==18
        arms.append({
            "arm":arm,"width":16,"training_examples":1728,"epochs":32,
            "training_updates":trained[2],
            "accuracy":sum(x["correct"] for x in cells)/sum(x["total"] for x in cells),
            "cells":cells
        })
    by={x["arm"]:x for x in arms}
    target=by["maskaware_augmented"]
    zeros=[r for r in target["cells"] if r["encoding"]=="explicit_all_zero"]
    nonzero=[r for r in target["cells"] if r["encoding"]!="explicit_all_zero"]
    zero_acc=sum(x["correct"] for x in zeros)/sum(x["total"] for x in zeros)
    hidden_acc=sum(x["correct"] for x in nonzero)/sum(x["total"] for x in nonzero)
    if zero_acc>=.85 and hidden_acc<=.55:
        classification="SENTINEL_DEPENDENCE"
    elif all(x["accuracy"]>=.8 for x in nonzero):
        classification="ENCODING_ROBUST"
    else:
        classification="MIXED"
    return {
        "schema":"wingless.w-r23-hidden-erasure-encoding.v1",
        "parent_head":"c551d41a7b8d4f77b9349033fa788d621a800519",
        "classification":classification,
        "maskaware_aug_zero_accuracy":zero_acc,
        "maskaware_aug_hidden_accuracy":hidden_acc,
        "heldout_seeds":HELDOUT,
        "arms":arms,"model_calls":0,
        "claim_limit":"Finite synthetic dropout substitution transfer only"
    }


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
