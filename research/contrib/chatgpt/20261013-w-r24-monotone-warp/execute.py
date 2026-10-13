"""W-R24: frozen monotone retiming out-of-distribution test of W-R22 weights."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r22-maskaware"/"execute.py"
_spec=importlib.util.spec_from_file_location("w22_frozen",SRC)
assert _spec and _spec.loader
w22=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w22)

HELDOUT=tuple(range(288,304))
WARPS={
"identity":tuple(range(16)),
"one_local_duplicate":(0,1,2,3,4,5,5,7,8,9,10,11,12,13,14,15),
"two_local_duplicates":(0,1,2,3,4,4,4,7,8,9,10,11,12,13,14,15),
}
ARMS=("naive_augmented","maskaware_augmented")
SWAPS=("swap10_11","swap12_13")


def example(seed,entity,generator,warp,mask_count,swap,negative):
    source=w22.w20.base_sequence(seed,entity,generator)
    indices=WARPS[warp]
    assert tuple(sorted(indices))==indices and indices[0]==0 and indices[-1]==15
    assert all(indices[i]==i for i in (11,12,13,14))
    retimed=tuple(source[i] for i in indices)
    seen=w22.w16.observed(retimed,seed,entity,3,79)
    masked=set(w22.MASK_ORDER[:mask_count])
    assert not masked.intersection({0,11,12,13,14})
    candidate=tuple((0,0,0,0) if i in masked else row
                    for i,row in enumerate(seen))
    other=w22.w19.permute(candidate,swap) if negative else candidate
    assert other[0]==candidate[0] and sorted(other)==sorted(candidate)
    assert (other==candidate)==(not negative)
    return source,other


def eval_cell(arm,model,generator,warp,mask_count):
    weights,bias,updates=model
    by_seed=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            for swap in SWAPS:
                ap,bp=example(seed,entity,generator,warp,mask_count,swap,False)
                an,bn=example(seed,entity,generator,warp,mask_count,swap,True)
                assert ap==an and bp[0]==bn[0] and sorted(bp)==sorted(bn)
                for positive,a,b in ((True,ap,bp),(False,an,bn)):
                    features=w22.featurize(a,b,arm)
                    assert len(features)==16
                    predicted=bias+sum(w*v for w,v in zip(weights,features))>=0
                    if positive and predicted:tp+=1
                    elif not positive and not predicted:tn+=1
                    elif not positive:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(row["tp"]+row["tn"] for row in by_seed)
    return {"generator":generator,"warp":warp,"mask_count":mask_count,
            "accuracy":correct/(16*36),"correct":correct,"total":16*36,
            "heldout_by_seed":by_seed}


def run():
    outcomes=[]
    for arm in ARMS:
        model=w22.train(arm)
        cells=[eval_cell(arm,model,generator,warp,mask_count)
               for generator in w22.w20.GENERATORS
               for warp in WARPS
               for mask_count in (3,5)]
        assert len(cells)==18
        warped=[c for c in cells if c["warp"]!="identity"]
        outcomes.append({"arm":arm,"width":16,"training_examples":1728,
                         "epochs":32,"train_updates":model[2],
                         "overall_accuracy":sum(c["correct"] for c in cells)/sum(c["total"] for c in cells),
                         "warped_accuracy":sum(c["correct"] for c in warped)/sum(c["total"] for c in warped),
                         "cells":cells})
    m={a["arm"]:a for a in outcomes}
    sensitive=m["maskaware_augmented"]
    naive=m["naive_augmented"]
    warped=[c for c in sensitive["cells"] if c["warp"]!="identity"]
    if all(c["accuracy"]>=.80 for c in warped) and naive["warped_accuracy"]<=.60:
        label="ORDER_PRESERVING_WARP_TRANSFER"
    elif sensitive["warped_accuracy"]>.60 and (sensitive["warped_accuracy"]-
          naive["warped_accuracy"])>=.05:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r24-monotone-warp.v1",
            "parent_head":"0e208895e192febd8873975dd2ed86eb771afc54",
            "classification":label,"arms":outcomes,
            "training_seeds":tuple(range(240,256)),"sealed_heldout_seeds":HELDOUT,
            "model_calls":0,
            "claim_limit":"Finite deterministic order-preserving time warp falsification"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
