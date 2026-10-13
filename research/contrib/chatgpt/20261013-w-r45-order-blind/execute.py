"""W-R45: genuine first/second balanced qualification of synthetic paired ranking."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r44-blind-scan"/"execute.py"
spec=importlib.util.spec_from_file_location("w44_pinned",SOURCE)
assert spec and spec.loader
w44=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w44)
w43=w44.w43
w42=w44.w42
w37=w44.w37

TRAIN=tuple(range(2336,2352))
HELDOUT=tuple(range(2352,2368))
GENERATORS=w43.GENERATORS
SITES=w43.TEST_SITES
CONDITIONS=("identity","xor55")
ARMS=("always_first_null","biased_original","antisymmetric_no_bias")

def fit_models():
    previous=w43.TRAIN
    try:
        w43.TRAIN=TRAIN
        singles,comparisons,identity=w43.training_data()
    finally:
        w43.TRAIN=previous
    assert len(singles)==len(comparisons)==1728
    biased=w43.fit_pairs(comparisons)
    weights=[0.0]*16
    updates=0
    for _ in range(32):
        for v,label in comparisons:
            score=sum(w*x for w,x in zip(weights,v))
            if label*score<=0:
                weights=[w+.1*label*x for w,x in zip(weights,v)]
                updates+=1
    symmetric={"weights":weights,"bias":0.0,"updates":updates}
    assert symmetric["bias"]==0.0
    return {"biased_original":biased,"antisymmetric_no_bias":symmetric},identity

def score(arm,model,a,b):
    if arm=="always_first_null":
        return 1.0
    delta=tuple(x-y for x,y in zip(a,b))
    return w43.score(model[arm],delta)

def credit(score,first_true):
    if abs(score)<=1e-10:return .5
    return float((score>0)==first_true)

def eval_cell(arm,model,condition,generator,site):
    rows=[]
    forward=reverse=total_pairs=tied_paired=0
    for seed in HELDOUT:
        first=second=ties=0.0
        same=0
        for entity in range(9):
            base,pos,neg=w37.recoded_example(seed,entity,generator,condition,site,3,False)
            barebase,barepos,bareneg=w37.recoded_example(seed,entity,generator,"identity",site,3,False)
            a=w42.feature(base,pos);b=w42.feature(base,neg)
            assert a==w42.feature(barebase,barepos)
            assert b==w42.feature(barebase,bareneg)
            same+=int(a==b)
            direct=score(arm,model,a,b)
            reversed_score=score(arm,model,b,a)
            if arm=="antisymmetric_no_bias":
                assert abs(direct+reversed_score)<=1e-7
            first+=credit(direct,True)
            second+=credit(reversed_score,False)
            ties+=int(abs(direct)<=1e-10)+int(abs(reversed_score)<=1e-10)
        assert first<=9 and second<=9
        forward+=first
        reverse+=second
        total_pairs+=9
        tied_paired+=ties
        rows.append({"seed":seed,"correct_credit_original_first":first,
                     "correct_credit_original_second":second,
                     "tied_presentations":ties,
                     "identical_feature_pairs":same})
    assert total_pairs==16*9
    return {"generator":generator,"site":site,
            "pairs":total_pairs,"correct_credit_original_first":forward,
            "correct_credit_original_second":reverse,
            "tied_presentations":tied_paired,
            "balanced_accuracy":(forward+reverse)/(2*total_pairs),
            "by_seed":rows}

def run():
    model,identity=fit_models()
    results=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[eval_cell(arm,model,condition,g,s)
                   for g in GENERATORS for s in SITES]
            first=sum(c["correct_credit_original_first"] for c in cells)
            second=sum(c["correct_credit_original_second"] for c in cells)
            count=sum(c["pairs"] for c in cells)
            results.append({"arm":arm,"condition":condition,
                            "matched_pairs":count,"presentations":count*2,
                            "correct_credit_original_first":first,
                            "correct_credit_original_second":second,
                            "forward_accuracy":first/count,
                            "reverse_accuracy":second/count,
                            "balanced_accuracy":(first+second)/(2*count),
                            "training_examples":0 if arm=="always_first_null" else 1728,
                            "training_updates":model[arm]["updates"] if arm!="always_first_null" else 0,
                            "cells":cells})
    by={(r["arm"],r["condition"]):r for r in results}
    assert all(by["always_first_null",cond]["balanced_accuracy"]==.5 for cond in CONDITIONS)
    assert all(by[arm,"identity"]["balanced_accuracy"]==
               by[arm,"xor55"]["balanced_accuracy"] for arm in ARMS)
    candidate=by["antisymmetric_no_bias","identity"]["balanced_accuracy"]
    if candidate>=.65 and candidate>=.5+.10:label="ORDER_BALANCED_RANK_SIGNAL"
    elif candidate>=.55:label="PARTIAL_SIGNAL"
    else:label="NO_CONFIRMED_SIGNAL"
    return {"schema":"wingless.w-r45-order-balanced-pair-rank.v1",
            "parent_head":"2f762cbc51f02cfe5cee94adc6d15cdfc7fc0721",
            "classification":label,
            "development_seeds":TRAIN,"sealed_heldout_seeds":HELDOUT,
            "training_identical_pairs":identity,
            "no_bias_balanced_accuracy":candidate,
            "prior_w43_forward_only_claim_status":"withheld_scientific_acceptance_orientation_confound",
            "results":results,"model_calls":0,
            "claim_limit":"Independent fresh first/second balanced synthetic two-candidate discrimination, not a standalone temporal classifier."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
