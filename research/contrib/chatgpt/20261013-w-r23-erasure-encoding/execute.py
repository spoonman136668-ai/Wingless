"""W-R23: unseen missing-row sentinel encoding transfer, without model retraining."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r22-maskaware"/"execute.py"
_spec=importlib.util.spec_from_file_location("w22_frozen",PARENT)
assert _spec and _spec.loader
w22=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w22)
w20=w22.w20
w16=w22.w16
w19=w22.w19

HELDOUT=tuple(range(272,288))
SENTINELS=((0,0,0,0),(255,255,255,255),(127,127,127,127))
MASK_COUNTS=(3,5)
MASK_ORDER=(2,5,8,10,3,6)
ARMS=("naive_augmented","maskaware_zero_only","maskaware_declared_sentinel_vocabulary")
SWAPS=("swap10_11","swap12_13")


def example(seed,entity,generator,masked,sentinel,swap,negative):
    base=w20.base_sequence(seed,entity,generator)
    observed=w16.observed(base,seed,entity,3,79)
    mask=set(MASK_ORDER[:masked])
    assert not mask.intersection({0,11,12,13,14})
    erased=tuple(sentinel if tick in mask else row
                 for tick,row in enumerate(observed))
    candidate=w19.permute(erased,swap) if negative else erased
    assert candidate[0]==erased[0] and sorted(candidate)==sorted(erased)
    assert (candidate==erased)==(not negative)
    return base,candidate


def visible_transition_feature(a,b,skip_values):
    differences=[]
    for channel in range(4):
        va=[r[channel] for r in a]
        vb=[r[channel] for r in b]
        for i in range(16):
            nxt=(i+1)%16
            if b[i] in skip_values or b[nxt] in skip_values:
                continue
            da=(va[nxt]-va[i])%256
            db=(vb[nxt]-vb[i])%256
            differences.append(w16.circular_distance(da,db))
    value=(sum(sorted(differences,reverse=True)[:2])/2
           if len(differences)>=2 else
           differences[0] if differences else 0.0)
    return [value]*16


def features(a,b,arm):
    if arm=="naive_augmented":
        return w19.feature(a,b,"mean_top_two_64")
    if arm=="maskaware_zero_only":
        return w22.maskaware_feature(a,b)
    assert arm=="maskaware_declared_sentinel_vocabulary"
    return visible_transition_feature(a,b,set(SENTINELS))


def evaluate(arm,model,generator,masked,sentinel):
    weights,bias,updates=model
    by_seed=[]
    aliases=0
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            base=w20.base_sequence(seed,entity,generator)
            observed=w16.observed(base,seed,entity,3,79)
            aliases+=sum(row in SENTINELS and tick not in MASK_ORDER[:masked]
                         for tick,row in enumerate(observed))
            for swap in SWAPS:
                ap,bp=example(seed,entity,generator,masked,sentinel,swap,False)
                an,bn=example(seed,entity,generator,masked,sentinel,swap,True)
                assert ap==an and bp[0]==bn[0] and sorted(bp)==sorted(bn)
                for positive,a,b in ((True,ap,bp),(False,an,bn)):
                    f=features(a,b,arm)
                    assert len(f)==16
                    pred=bias+sum(x*y for x,y in zip(weights,f))>=0
                    if positive and pred:tp+=1
                    elif not positive and not pred:tn+=1
                    elif not positive:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(s["tp"]+s["tn"] for s in by_seed)
    return {"generator":generator,"mask_count":masked,"sentinel":list(sentinel),
            "accuracy":correct/(16*36),"correct":correct,"total":16*36,
            "unmasked_natural_sentinel_alias_count":aliases,"heldout_by_seed":by_seed}


def run():
    models={
        "naive_augmented":w22.train("naive_augmented"),
        "maskaware_zero_only":w22.train("maskaware_augmented"),
        "maskaware_declared_sentinel_vocabulary":w22.train("maskaware_augmented"),
    }
    assert models["maskaware_zero_only"]==models["maskaware_declared_sentinel_vocabulary"]
    outcomes=[]
    for arm in ARMS:
        cells=[evaluate(arm,models[arm],generator,masked,sentinel)
               for generator in w20.GENERATORS
               for masked in MASK_COUNTS
               for sentinel in SENTINELS]
        assert len(cells)==18
        nonzero=[c for c in cells if c["sentinel"]!=[0,0,0,0]]
        outcomes.append({"arm":arm,"width":16,"training_pairs":1728,"epochs":32,
                         "train_updates":models[arm][2],
                         "total_accuracy":sum(c["correct"] for c in cells)/sum(c["total"] for c in cells),
                         "nonzero_accuracy":sum(c["correct"] for c in nonzero)/sum(c["total"] for c in nonzero),
                         "cells":cells})
    mapping={o["arm"]:o for o in outcomes}
    new=mapping["maskaware_declared_sentinel_vocabulary"]
    old=mapping["maskaware_zero_only"]
    diff=new["nonzero_accuracy"]-old["nonzero_accuracy"]
    nonzero_cells=[c for c in new["cells"] if c["sentinel"]!=[0,0,0,0]]
    if all(c["accuracy"]>=0.85 for c in nonzero_cells) and diff>=0.10:
        classification="SENTINEL_VOCABULARY_TRANSFER"
    elif diff>=0.05:
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {"schema":"wingless.w-r23-erasure-encoding-invariance.v1",
            "parent_head":"c551d41a7b8d4f77b9349033fa788d621a800519",
            "classification":classification,
            "delta_vocabulary_vs_zero_only_nonzero":diff,
            "development_seeds":list(range(240,256)),"sealed_heldout_seeds":HELDOUT,
            "arms":outcomes,"model_calls":0,
            "claim_limit":"Finite synthetic missing-data encodings, not general perception"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
