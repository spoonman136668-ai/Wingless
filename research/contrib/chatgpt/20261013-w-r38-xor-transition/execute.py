"""W-R38: zero-shot XOR-invariant paired temporal observation feature."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r37-codebook-conditioned"/"execute.py"
spec=importlib.util.spec_from_file_location("w37_pinned",SOURCE)
assert spec and spec.loader
w37=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w37)

TRAIN=tuple(range(1760,1776))
HELDOUT=tuple(range(1776,1792))
ARMS=("arithmetic","xor_invariant")
CONDITIONS=("identity","xor55")
GENERATORS=w37.GENERATORS
TEST_SITES=w37.TEST_SITES
TRAIN_SITES=w37.TRAIN_SITES
BITS=4


def feature(base,candidate,arm):
    if arm=="arithmetic":
        v=w37.features(base,candidate)
    else:
        assert arm=="xor_invariant"
        assert len(base)==len(candidate)==16
        residual=[]
        for channel in range(4):
            for tick in range(16):
                nxt=(tick+1)%16
                mismatch=(base[tick][channel]^base[nxt][channel]^
                          candidate[tick][channel]^candidate[nxt][channel])
                residual.append(mismatch.bit_count())
        assert len(residual)==64
        top=sorted(residual,reverse=True)[:2]
        v=[sum(top)/2.0]*16
    assert len(v)==16
    return v


def train(arm):
    examples=[]
    identical=0
    for seed in TRAIN:
        for entity in range(9):
            for site in TRAIN_SITES:
                for amp in (1,2):
                    base,positive,negative=w37.recoded_example(
                        seed,entity,"affine_prior","identity",site,amp,True)
                    identical+=int(positive==negative)
                    examples.append((feature(base,positive,arm),1))
                    examples.append((feature(base,negative,arm),-1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for v,label in examples:
            if label*(bias+sum(w*x for w,x in zip(weights,v)))<=0:
                weights=[w+0.1*label*x for w,x in zip(weights,v)]
                bias+=0.1*label
                updates+=1
    return weights,bias,updates,identical


def eval_cell(arm,model,condition,generator,site):
    weights,bias,updates,identical=model
    by_seed=[]
    correct=total=indistinguishable=checks=0
    for seed in HELDOUT:
        tp=tn=fp=fn=same=verified=0
        for entity in range(9):
            base,positive,negative=w37.recoded_example(
                seed,entity,generator,condition,site,3,False)
            clean_base,clean_pos,clean_neg=w37.recoded_example(
                seed,entity,generator,"identity",site,3,False)
            same+=int(positive==negative)
            for truth,candidate,reference in ((True,positive,clean_pos),
                                               (False,negative,clean_neg)):
                fv=feature(base,candidate,arm)
                if arm=="xor_invariant":
                    assert fv==feature(clean_base,reference,arm)
                    verified+=1
                guess=bias+sum(w*x for w,x in zip(weights,fv))>=0
                if truth and guess:tp+=1
                elif not truth and not guess:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        indistinguishable+=same
        checks+=verified
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "observation_identical_pairs":same,
                        "feature_invariance_checks":verified})
    return {"generator":generator,"site":site,
            "correct":correct,"total":total,"accuracy":correct/total,
            "indistinguishable_pairs":indistinguishable,
            "invariance_checks":checks,"by_seed":by_seed}


def run():
    models={arm:train(arm) for arm in ARMS}
    results=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[eval_cell(arm,models[arm],condition,g,s)
                   for g in GENERATORS for s in TEST_SITES]
            results.append({"arm":arm,"condition":condition,
                            "train_examples":1728,"epochs":32,"feature_width":16,
                            "training_updates":models[arm][2],
                            "indistinguishable_training_pairs":models[arm][3],
                            "accuracy":sum(c["correct"] for c in cells)/
                                       sum(c["total"] for c in cells),
                            "invariance_checks":sum(c["invariance_checks"] for c in cells),
                            "cells":cells})
    by={(r["arm"],r["condition"]):r for r in results}
    invariant=(by["xor_invariant","identity"]["accuracy"]==
               by["xor_invariant","xor55"]["accuracy"] and
               by["xor_invariant","xor55"]["invariance_checks"]==6*16*9*2)
    loss=by["arithmetic","identity"]["accuracy"]-by["arithmetic","xor55"]["accuracy"]
    useful=(by["xor_invariant","identity"]["accuracy"]>=.65 and
            by["xor_invariant","xor55"]["accuracy"]>=.65)
    if invariant and useful and loss>=.10:
        label="ZERO_SHOT_XOR_TRANSFER"
    elif invariant and not useful:
        label="INVARIANT_BUT_NOT_USEFUL"
    else:
        label="OTHERWISE_NEGATIVE"
    return {"schema":"wingless.w-r38-xor-invariant-observed-feature.v1",
            "parent_head":"74fa95538c7df1092243d0637cee2bdd1b02a3d2",
            "classification":label,"development_seeds":TRAIN,
            "sealed_heldout_seeds":HELDOUT,"xor_feature_invariant":invariant,
            "arithmetic_identity_minus_xor_accuracy":loss,
            "rows":results,"model_calls":0,
            "claim_limit":"Exact synthetic XOR55 transition-feature invariance; utility and generalization are separate."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
