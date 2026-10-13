"""W-R25: frozen same-width learned monotone alignment under missing ticks."""
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

DEV=tuple(range(304,320))
HELDOUT=tuple(range(320,336))
ARMS=("maskaware_transition","monotone_dp")
TRAIN_SWAP_INDICES=(2,5,8)
GAP_PENALTY=.10
MAX_INDEX_DEVIATION=2


def dp_feature(a,b):
    # Entirely observed zero-coded erasures, no privileged mask variable.
    visible=[(i,row) for i,row in enumerate(b) if row!=(0,0,0,0)]
    assert visible and visible[0][0]==0 and visible[-1][0]==15
    origin_a=a[0]
    origin_b=b[0]
    aa=[tuple((row[c]-origin_a[c])%256 for c in range(4)) for row in a]
    bb=[(tick,tuple((row[c]-origin_b[c])%256 for c in range(4)))
        for tick,row in visible]
    n=len(aa)
    m=len(bb)
    infinity=float("inf")
    table=[[infinity]*m for _ in range(n)]
    for i in range(n):
        for j,(tick,other) in enumerate(bb):
            if abs(i-tick)>MAX_INDEX_DEVIATION:
                continue
            cost=sum(w16.circular_distance(aa[i][c],other[c])
                     for c in range(4))/4
            if i==0 and j==0:
                table[i][j]=cost
                continue
            predecessors=[]
            if i>0 and j>0:
                predecessors.append(table[i-1][j-1])
            if i>0:
                predecessors.append(table[i-1][j]+GAP_PENALTY)
            if j>0:
                predecessors.append(table[i][j-1]+GAP_PENALTY)
            if predecessors:
                table[i][j]=cost+min(predecessors)
    total=table[n-1][m-1]
    assert 0<=total<float("inf")
    return [total/16]*16


def features(a,b,arm):
    if arm=="maskaware_transition":
        return w22.maskaware_feature(a,b)
    assert arm=="monotone_dp"
    return dp_feature(a,b)


def train(arm):
    cases=[]
    for seed in DEV:
        for entity in range(9):
            for swap in TRAIN_SWAP_INDICES:
                for amp in (1,2):
                    masked=(0,2,4)[(seed+entity+swap+amp)%3]
                    for negative in (False,True):
                        a,b=w22.prepared(seed,entity,"affine_prior",
                                         amp,"development",masked,negative,swap)
                        vector=features(a,b,arm)
                        assert len(vector)==16
                        cases.append((vector,-1 if negative else 1))
    assert len(cases)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vector,label in cases:
            margin=label*(bias+sum(w*v for w,v in zip(weights,vector)))
            if margin<=0:
                weights=[w+.1*label*v for w,v in zip(weights,vector)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def heldout_pair(seed,entity,generator,warp,mask_count,swap,negative):
    base=w22.w20.base_sequence(seed,entity,generator)
    retimed=tuple(base[i] for i in w24.WARPS[warp])
    observed=w16.observed(retimed,seed,entity,3,79)
    mask=set(w22.MASK_ORDER[:mask_count])
    assert not mask.intersection({0,11,12,13,14})
    clean=tuple((0,0,0,0) if i in mask else row
                for i,row in enumerate(observed))
    other=w22.w19.permute(clean,swap) if negative else clean
    assert other[0]==clean[0] and sorted(other)==sorted(clean)
    assert (other==clean)==(not negative)
    return base,other


def evaluate(arm,model,generator,warp,masked):
    weights,bias,updates=model
    by_seed=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            for swap in ("swap10_11","swap12_13"):
                positive_a,positive_b=heldout_pair(seed,entity,generator,warp,masked,swap,False)
                negative_a,negative_b=heldout_pair(seed,entity,generator,warp,masked,swap,True)
                assert positive_a==negative_a and positive_b[0]==negative_b[0]
                assert sorted(positive_b)==sorted(negative_b)
                for label,a,b in ((1,positive_a,positive_b),(-1,negative_a,negative_b)):
                    f=features(a,b,arm)
                    pred=1 if bias+sum(w*v for w,v in zip(weights,f))>=0 else -1
                    if label==1 and pred==1:tp+=1
                    elif label==-1 and pred==-1:tn+=1
                    elif label==-1:fp+=1
                    else:fn+=1
        assert tp+tn+fp+fn==36
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(s["tp"]+s["tn"] for s in by_seed)
    return {"generator":generator,"warp":warp,"mask_count":masked,
            "correct":correct,"total":16*36,"accuracy":correct/(16*36),
            "by_seed":by_seed}


def run():
    rows=[]
    for arm in ARMS:
        model=train(arm)
        cells=[evaluate(arm,model,generator,warp,masked)
               for generator in w22.w20.GENERATORS
               for warp in w24.WARPS
               for masked in (3,5)]
        hard=[c for c in cells if c["warp"]!="identity"]
        rows.append({"arm":arm,"width":16,"train_pairs":1728,"epochs":32,
                     "updates":model[2],
                     "overall_accuracy":sum(c["correct"] for c in cells)/sum(c["total"] for c in cells),
                     "warped_accuracy":sum(c["correct"] for c in hard)/sum(c["total"] for c in hard),
                     "cells":cells})
    baseline=rows[0]
    candidate=rows[1]
    difficult=[c for c in candidate["cells"] if c["warp"]!="identity"]
    gain=candidate["warped_accuracy"]-baseline["warped_accuracy"]
    if all(c["accuracy"]>=.8 for c in difficult) and gain>=.1:
        classification="MONOTONE_ALIGNMENT_TRANSFER"
    elif gain>=.05:
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {"schema":"wingless.w-r25-same-width-monotone-dp.v1",
            "parent_head":"de6234a46160970e4c23f74d11e9580683047319",
            "classification":classification,"warped_dp_vs_baseline_gain":gain,
            "development_seeds":DEV,"sealed_heldout_seeds":HELDOUT,
            "arms":rows,"model_calls":0,
            "claim_limit":"Finite synthetic monotone local alignment; DP computation larger than baseline"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
