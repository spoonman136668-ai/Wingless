"""W-R40: position-aware XOR-invariant fixed-width temporal state pilots."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r39-multistat"/"execute.py"
spec=importlib.util.spec_from_file_location("w39_pinned",SOURCE)
assert spec and spec.loader
w39=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w39)
w37=w39.w37

TRAIN=tuple(range(1856,1872))
HELDOUT=tuple(range(1872,1888))
ARMS=("scalar_repeated","ordered_lag1","ordered_lag2")
CONDITIONS=("identity","xor55")
GENERATORS=w39.GENERATORS
TEST_SITES=w39.TEST_SITES
TRAIN_SITES=w39.TRAIN_SITES


def feature(base,candidate,arm):
    assert len(base)==len(candidate)==16
    if arm=="scalar_repeated":
        vec=w39.feature(base,candidate,"scalar_repeated")
    else:
        step=1 if arm=="ordered_lag1" else 2
        assert arm in ("ordered_lag1","ordered_lag2")
        vec=[]
        for tick in range(16):
            nxt=(tick+step)%16
            v=sum((base[tick][c]^base[nxt][c]^
                   candidate[tick][c]^candidate[nxt][c]).bit_count()
                  for c in range(4))
            vec.append(v)
    assert len(vec)==16
    return vec


def train(arm):
    examples=[]
    identical=0
    for seed in TRAIN:
        for entity in range(9):
            for site in TRAIN_SITES:
                for amplitude in (1,2):
                    base,pos,neg=w37.recoded_example(
                        seed,entity,"affine_prior","identity",site,amplitude,True)
                    identical+=int(pos==neg)
                    examples.append((feature(base,pos,arm),1))
                    examples.append((feature(base,neg,arm),-1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in examples:
            if label*(bias+sum(w*x for w,x in zip(weights,vec)))<=0:
                weights=[w+0.1*label*x for w,x in zip(weights,vec)]
                bias+=0.1*label
                updates+=1
    return weights,bias,updates,identical


def evaluate(arm,model,condition,generator,site):
    weights,bias,updates,identical=model
    by_seed=[]
    correct=total=identical_pairs=verified=0
    for seed in HELDOUT:
        tp=tn=fp=fn=same=checks=0
        for entity in range(9):
            base,pos,neg=w37.recoded_example(
                seed,entity,generator,condition,site,3,False)
            origbase,origpos,origneg=w37.recoded_example(
                seed,entity,generator,"identity",site,3,False)
            same+=int(pos==neg)
            for truth,candidate,original in ((True,pos,origpos),(False,neg,origneg)):
                vec=feature(base,candidate,arm)
                assert vec==feature(origbase,original,arm)
                checks+=1
                prediction=bias+sum(w*x for w,x in zip(weights,vec))>=0
                if truth and prediction:tp+=1
                elif not truth and not prediction:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        identical_pairs+=same
        verified+=checks
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "identical_pairs":same,"feature_equivalence_checks":checks})
    return {"generator":generator,"site":site,"correct":correct,"total":total,
            "accuracy":correct/total,"identical_pairs":identical_pairs,
            "feature_equivalence_checks":verified,"by_seed":by_seed}


def run():
    models={arm:train(arm) for arm in ARMS}
    results=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[evaluate(arm,models[arm],condition,g,site)
                   for g in GENERATORS for site in TEST_SITES]
            results.append({"arm":arm,"condition":condition,
                            "training_examples":1728,"epochs":32,"width":16,
                            "training_updates":models[arm][2],
                            "identical_training_pairs":models[arm][3],
                            "accuracy":sum(c["correct"] for c in cells)/
                                       sum(c["total"] for c in cells),
                            "feature_equivalence_checks":sum(c["feature_equivalence_checks"] for c in cells),
                            "cells":cells})
    by={(r["arm"],r["condition"]):r for r in results}
    invariant=all(by[arm,"identity"]["accuracy"]==by[arm,"xor55"]["accuracy"] and
                  by[arm,"xor55"]["feature_equivalence_checks"]==6*16*9*2
                  for arm in ARMS)
    assert invariant
    baseline=by["scalar_repeated","identity"]["accuracy"]
    gains={arm:by[arm,"identity"]["accuracy"]-baseline for arm in ARMS if arm!="scalar_repeated"}
    if any(by[arm,"identity"]["accuracy"]>=.65 and
           by[arm,"xor55"]["accuracy"]>=.65 and gains[arm]>=.10
           for arm in gains):
        classification="POSITIONAL_ROBUST_UTILITY"
    elif any(x>=.05 for x in gains.values()):
        classification="PARTIAL_GAIN"
    else:
        classification="NO_GAIN"
    return {"schema":"wingless.w-r40-positional-xor-state.v1",
            "parent_head":"9bd4d26f3fc6c47a15f73a4e222973e1c0add30f",
            "classification":classification,
            "development_seeds":TRAIN,"sealed_heldout_seeds":HELDOUT,
            "xor_feature_invariance":invariant,"accuracy_gains_over_scalar":gains,
            "results":results,"model_calls":0,
            "claim_limit":"Synthetic position-preserving XOR55-invariant temporal state with fixed 16 numerical weights; task utility is independently tested."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
