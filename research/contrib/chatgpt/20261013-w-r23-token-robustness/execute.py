"""W-R23: frozen missing-data encoding transfer; unchanged W-R22 learned models."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r22-maskaware"/"execute.py"
spec=importlib.util.spec_from_file_location("w22_frozen",SOURCE)
assert spec and spec.loader
w22=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w22)
w20=w22.w20
w19=w22.w19

SEEDS=tuple(range(272,288))
GENERATORS=("affine_prior","quadratic_irregular","xorshift_independent")
STYLES=("all_zero_seen","all_255_unseen","previous_row_fill_unseen")
MASK_COUNTS=(2,4,6)
SWAPS=("swap10_11","swap12_13")
ARMS=("naive_augmented","maskaware_clean","maskaware_augmented")
MASK_ORDER=(2,5,8,10,3,6)

def transform(observed, count, style):
    absent=set(MASK_ORDER[:count])
    assert not absent.intersection({0,11,12,13,14})
    filled=[]
    previous_unmasked=None
    for tick,row in enumerate(observed):
        if tick not in absent:
            previous_unmasked=row
            filled.append(row)
        elif style=="all_zero_seen":
            filled.append((0,0,0,0))
        elif style=="all_255_unseen":
            filled.append((255,255,255,255))
        else:
            assert style=="previous_row_fill_unseen" and previous_unmasked is not None
            filled.append(previous_unmasked)
    result=tuple(filled)
    assert len(result)==16
    return result

def paired(seed,entity,generator,mask_count,style,swap):
    source=w20.base_sequence(seed,entity,generator)
    observed=w22.w16.observed(source,seed,entity,3,79)
    positive=transform(observed,mask_count,style)
    negative=w19.permute(positive,swap)
    assert positive[0]==negative[0]
    assert sorted(positive)==sorted(negative)
    assert positive!=negative
    return source,positive,negative

def evaluate(arm,model,generator,count,style,swap):
    weights,bias,updates=model
    by_seed=[]
    for seed in SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            base,p,n=paired(seed,entity,generator,count,style,swap)
            for positive,candidate in ((True,p),(False,n)):
                vec=w22.featurize(base,candidate,arm)
                pred=bias+sum(w*x for w,x in zip(weights,vec))>=0
                if positive and pred:tp+=1
                elif not positive and not pred:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    good=sum(r["tp"]+r["tn"] for r in by_seed)
    return {"generator":generator,"mask_count":count,"fill_style":style,"swap":swap,
            "accuracy":good/(16*18),"correct":good,"total":16*18,
            "heldout_by_seed":by_seed}

def run():
    arms=[]
    for arm in ARMS:
        model=w22.train(arm)
        cells=[evaluate(arm,model,g,c,s,swap)
               for g in GENERATORS for c in MASK_COUNTS for s in STYLES for swap in SWAPS]
        assert len(cells)==54
        summaries=cells
        novel=[c for c in cells if c["fill_style"]!="all_zero_seen"]
        arms.append({"arm":arm,"width":len(model[0]),"training_examples":1728,
                     "epochs":32,"training_updates":model[2],
                     "aggregate_accuracy":sum(x["correct"] for x in cells)/sum(x["total"] for x in cells),
                     "novel_fill_accuracy":sum(x["correct"] for x in novel)/sum(x["total"] for x in novel),
                     "cells":summaries})
    by={a["arm"]:a for a in arms}
    models=("maskaware_clean","maskaware_augmented")
    if all(all(c["accuracy"]>=.8 for c in by[arm]["cells"] if
               c["fill_style"]!="all_zero_seen") for arm in models):
        classification="ENCODING_INVARIANT"
    elif by["maskaware_augmented"]["novel_fill_accuracy"]>=.7:
        classification="PARTIAL"
    else:
        classification="ENCODING_DEPENDENT"
    return {"schema":"wingless.w-r23-missingness-encoding.v1",
            "parent_head":"c551d41a7b8d4f77b9349033fa788d621a800519",
            "classification":classification,
            "training_seeds":w22.TRAIN_SEEDS,"sealed_heldout_seeds":SEEDS,
            "arms":arms,"model_calls":0,
            "claim_limit":"Synthetic dropout encoding transfer, not general intelligence"}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
