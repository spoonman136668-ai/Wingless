"""W-R39: fixed-budget XOR-invariant multi-statistic temporal representation."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r38-xor-transition"/"execute.py"
spec=importlib.util.spec_from_file_location("w38_pinned",SOURCE)
assert spec and spec.loader
w38=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w38)
w37=w38.w37

TRAIN=tuple(range(1824,1840))
HELDOUT=tuple(range(1840,1856))
ARMS=("scalar_repeated","top16_sorted","histogram")
CONDITIONS=("identity","xor55")
GENERATORS=w38.GENERATORS
TEST_SITES=w38.TEST_SITES
TRAIN_SITES=w38.TRAIN_SITES

def residuals(base,candidate):
    assert len(base)==len(candidate)==16
    return [((base[i][c]^base[(i+1)%16][c]^
              candidate[i][c]^candidate[(i+1)%16][c]).bit_count())
            for c in range(4) for i in range(16)]

def feature(base,candidate,arm):
    distances=residuals(base,candidate)
    assert len(distances)==64 and all(0<=x<=8 for x in distances)
    if arm=="scalar_repeated":
        vec=w38.feature(base,candidate,"xor_invariant")
    elif arm=="top16_sorted":
        vec=list(sorted(distances,reverse=True)[:16])
    else:
        assert arm=="histogram"
        vec=[sum(int(d==level) for d in distances) for level in range(9)]+[0]*7
    assert len(vec)==16
    return vec


def train(arm):
    examples=[]
    same=0
    for seed in TRAIN:
        for entity in range(9):
            for site in TRAIN_SITES:
                for amp in (1,2):
                    base,pos,neg=w37.recoded_example(
                        seed,entity,"affine_prior","identity",site,amp,True)
                    same+=int(pos==neg)
                    examples.append((feature(base,pos,arm),1))
                    examples.append((feature(base,neg,arm),-1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for v,label in examples:
            score=bias+sum(w*x for w,x in zip(weights,v))
            if label*score<=0:
                weights=[w+0.1*label*x for w,x in zip(weights,v)]
                bias+=0.1*label
                updates+=1
    return weights,bias,updates,same


def eval_cell(arm,model,condition,generator,site):
    weights,bias,_,_=model
    rows=[]
    correct=total=identical=checks=0
    for seed in HELDOUT:
        tp=tn=fp=fn=same=verified=0
        for entity in range(9):
            base,pos,neg=w37.recoded_example(
                seed,entity,generator,condition,site,3,False)
            rawbase,rawpos,rawneg=w37.recoded_example(
                seed,entity,generator,"identity",site,3,False)
            same+=int(pos==neg)
            for truth,candidate,raw in ((True,pos,rawpos),(False,neg,rawneg)):
                v=feature(base,candidate,arm)
                assert v==feature(rawbase,raw,arm)
                verified+=1
                pred=bias+sum(w*x for w,x in zip(weights,v))>=0
                if truth and pred:tp+=1
                elif not truth and not pred:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        identical+=same
        checks+=verified
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "identical_observation_pairs":same,
                     "feature_identity_checks":verified})
    return {"generator":generator,"site":site,
            "correct":correct,"total":total,"accuracy":correct/total,
            "identical_pairs":identical,"checks":checks,"by_seed":rows}


def run():
    trained={arm:train(arm) for arm in ARMS}
    results=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[eval_cell(arm,trained[arm],condition,g,s)
                   for g in GENERATORS for s in TEST_SITES]
            results.append({"arm":arm,"condition":condition,"width":16,
                            "train_examples":1728,"epochs":32,
                            "updates":trained[arm][2],
                            "identical_training_pairs":trained[arm][3],
                            "accuracy":sum(c["correct"] for c in cells)/
                                       sum(c["total"] for c in cells),
                            "identity_checks":sum(c["checks"] for c in cells),
                            "cells":cells})
    by={(r["arm"],r["condition"]):r["accuracy"] for r in results}
    baseline=by["scalar_repeated","identity"]
    candidate_gains={arm:by[arm,"identity"]-baseline for arm in ARMS if arm!="scalar_repeated"}
    invariant=all(by[arm,"identity"]==by[arm,"xor55"] for arm in ARMS)
    if (invariant and any(by[arm,"identity"]>=.65 and
                          by[arm,"xor55"]>=.65 and candidate_gains[arm]>=.10
                          for arm in candidate_gains)):
        classification="XOR_ROBUST_UTILITY"
    elif any(gain>=.05 for gain in candidate_gains.values()):
        classification="PARTIAL_GAIN"
    else:
        classification="NO_GAIN"
    return {"schema":"wingless.w-r39-multistat-xor-invariant.v1",
            "parent_head":"d9d51080086a39fcbf0f2fae0c08e7a78d6412de",
            "classification":classification,
            "development_seeds":TRAIN,"sealed_heldout_seeds":HELDOUT,
            "mathematical_xor_invariance_verified":invariant,
            "candidate_identity_accuracy_gains":candidate_gains,
            "results":results,"model_calls":0,
            "claim_limit":"Synthetic XOR-invariant bitwise temporal order statistics; added effective expressivity within unchanged numerical width is disclosed."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
