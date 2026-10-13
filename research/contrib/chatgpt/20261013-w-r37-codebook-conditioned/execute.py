"""W-R37 equal-budget supervised transfer across identity and XOR byte codebooks."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r36-codebook-transfer"/"execute.py"
spec=importlib.util.spec_from_file_location("w36_pinned",SOURCE)
assert spec and spec.loader
w36=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w36)
w34=w36.w34
w32=w34.w32
w20=w34.w20
w16=w34.w16

TRAIN=tuple(range(1664,1680))
HELDOUT=tuple(range(1680,1696))
CONDITIONS=("identity","xor55")
BITS=4
GENERATORS=w34.GENERATORS
TEST_SITES=w34.SITES
TRAIN_SITES=w34.INDICES


def features(base,candidate):
    result=w34.representation(base,candidate,"asymmetric_refit",BITS)
    assert len(result)==16
    return result


def recoded_example(seed,entity,generator,condition,site,amplitude,training):
    base=w20.base_sequence(seed,entity,generator)
    original=w16.observed(base,seed,entity,amplitude,17 if training else 79)
    candidate=tuple(w32.quantize(row,BITS) for row in original)
    tail=list(candidate[1:])
    tail[site],tail[site+1]=tail[site+1],tail[site]
    negative=candidate[:1]+tuple(tail)
    assert candidate[0]==negative[0]
    assert sorted(candidate)==sorted(negative)
    raw=w36.recode(base,condition)
    positive=w36.recode(candidate,condition)
    negative=w36.recode(negative,condition)
    assert sorted(positive)==sorted(negative)
    return raw,positive,negative


def train(condition):
    examples=[]
    identical=0
    for seed in TRAIN:
        for entity in range(9):
            for site in TRAIN_SITES:
                for amplitude in (1,2):
                    base,pos,neg=recoded_example(
                        seed,entity,"affine_prior",condition,site,amplitude,True)
                    identical+=int(pos==neg)
                    examples.append((features(base,pos),1))
                    examples.append((features(base,neg),-1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in examples:
            score=bias+sum(w*x for w,x in zip(weights,vec))
            if label*score<=0:
                weights=[w+0.1*label*x for w,x in zip(weights,vec)]
                bias+=0.1*label
                updates+=1
    return weights,bias,updates,identical


def evaluate(model,condition,generator,site):
    weights,bias,_,_=model
    by_seed=[]
    total=correct=identical=0
    for seed in HELDOUT:
        tp=tn=fp=fn=same=0
        for entity in range(9):
            base,pos,neg=recoded_example(seed,entity,generator,condition,site,3,False)
            same+=int(pos==neg)
            for truth,candidate in ((True,pos),(False,neg)):
                v=features(base,candidate)
                predicted=bias+sum(w*x for w,x in zip(weights,v))>=0
                if truth and predicted:tp+=1
                elif not truth and not predicted:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        total+=18
        correct+=tp+tn
        identical+=same
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "observation_identical_pairs":same})
    return {"generator":generator,"site":site,"correct":correct,"total":total,
            "accuracy":correct/total,"identical_pairs":identical,"by_seed":by_seed}


def run():
    rows=[]
    trained={c:train(c) for c in CONDITIONS}
    for training_condition in CONDITIONS:
        model=trained[training_condition]
        for test_condition in CONDITIONS:
            cells=[evaluate(model,test_condition,g,s)
                   for g in GENERATORS for s in TEST_SITES]
            assert len(cells)==6
            rows.append({"training_condition":training_condition,
                         "test_condition":test_condition,
                         "accuracy":sum(x["correct"] for x in cells)/
                                    sum(x["total"] for x in cells),
                         "width":16,"train_examples":1728,"epochs":32,
                         "training_updates":model[2],
                         "identical_training_pairs":model[3],
                         "cells":cells})
    by={(x["training_condition"],x["test_condition"]):x["accuracy"] for x in rows}
    gain=by["xor55","xor55"]-by["identity","xor55"]
    if gain>=.10 and by["xor55","xor55"]>=.65 and by["identity","identity"]>=.65:
        label="CODEBOOK_CONDITIONED_RECOVERY"
    elif gain>=.05:
        label="PARTIAL_GAIN"
    else:
        label="NO_GAIN"
    return {"schema":"wingless.w-r37-codebook-conditioned-refit.v1",
            "parent_head":"595a78cf0759e205c4e86510e3fc7a10cfdd53e3",
            "classification":label,
            "development_seeds":TRAIN,"sealed_heldout_seeds":HELDOUT,
            "xor_conditioned_accuracy_gain":gain,
            "rows":rows,"model_calls":0,
            "claim_limit":"Synthetic equal-budget codebook-conditioned supervised refit; not zero-shot representation invariance."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
