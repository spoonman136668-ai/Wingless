"""W-R43: falsify pair-relative readout gains under explicit two-candidate input."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r42-nonlinear"/"execute.py"
spec=importlib.util.spec_from_file_location("w42_pinned",SOURCE)
assert spec and spec.loader
w42=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w42)
w37=w42.w37
TRAIN=tuple(range(2112,2128))
HELDOUT=tuple(range(2128,2144))
GENERATORS=w42.GENERATORS
TEST_SITES=w42.TEST_SITES
TRAIN_SITES=w42.TRAIN_SITES
CONDITIONS=("identity","xor55")
ARMS=("independent_perceptron","paired_difference_perceptron")

def training_data():
    singles=[]
    comparisons=[]
    identical_pairs=0
    for seed in TRAIN:
        for entity in range(9):
            for site in TRAIN_SITES:
                for amplitude in (1,2):
                    base,pos,neg=w37.recoded_example(
                        seed,entity,"affine_prior","identity",site,amplitude,True)
                    x=w42.feature(base,pos)
                    y=w42.feature(base,neg)
                    identical_pairs+=int(x==y)
                    singles.extend(((x,1),(y,-1)))
                    delta=tuple(a-b for a,b in zip(x,y))
                    comparisons.extend(((delta,1),(tuple(-a for a in delta),-1)))
    assert len(singles)==len(comparisons)==1728
    return singles,comparisons,identical_pairs

def fit_pairs(examples):
    weights=[0.0]*16
    bias=0.0
    updates=0
    for epoch in range(32):
        for vec,label in examples:
            if label*(bias+sum(w*x for w,x in zip(weights,vec)))<=0:
                weights=[w+.1*label*x for w,x in zip(weights,vec)]
                bias+=.1*label
                updates+=1
    return {"weights":weights,"bias":bias,"updates":updates}

def score(model,vector):
    return model["bias"]+sum(w*x for w,x in zip(model["weights"],vector))

def evaluate(arm,model,condition,generator,site):
    by_seed=[]
    prefer=wrong=ties=0
    feature_ties=0
    for seed in HELDOUT:
        p=w=t=identical=0
        for entity in range(9):
            base,pos,neg=w37.recoded_example(
                seed,entity,generator,condition,site,3,False)
            unbase,unpos,unneg=w37.recoded_example(
                seed,entity,generator,"identity",site,3,False)
            a=w42.feature(base,pos)
            b=w42.feature(base,neg)
            assert a==w42.feature(unbase,unpos) and b==w42.feature(unbase,unneg)
            identical+=int(a==b)
            if arm=="independent_perceptron":
                # Both options are presented to this evaluator; the standalone
                # classifier is not being scored on its original interface.
                diff=score(model,a)-score(model,b)
            else:
                diff=score(model,tuple(x-y for x,y in zip(a,b)))
            if diff>1e-10:p+=1
            elif diff< -1e-10:w+=1
            else:t+=1
        assert p+w+t==9
        prefer+=p
        wrong+=w
        ties+=t
        feature_ties+=identical
        by_seed.append({"seed":seed,"preferred_true":p,"preferred_false":w,
                        "score_ties":t,"identical_feature_pairs":identical})
    total=prefer+wrong+ties
    assert total==16*9
    return {"generator":generator,"site":site,"true_preferred":prefer,
            "false_preferred":wrong,"tie_count":ties,
            "identical_feature_pairs":feature_ties,
            "balanced_tie_credit_accuracy":(prefer+0.5*ties)/total,
            "by_seed":by_seed}

def run():
    singles,comparisons,train_identical=training_data()
    independent=w42.fit_linear(singles)
    paired=fit_pairs(comparisons)
    models={ARMS[0]:independent,ARMS[1]:paired}
    results=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[evaluate(arm,models[arm],condition,g,site)
                   for g in GENERATORS for site in TEST_SITES]
            pos=sum(c["true_preferred"] for c in cells)
            bad=sum(c["false_preferred"] for c in cells)
            tie=sum(c["tie_count"] for c in cells)
            results.append({"arm":arm,"condition":condition,
                            "training_examples":1728,"epochs":32,
                            "weights":16,"training_updates":models[arm]["updates"],
                            "paired_preference_accuracy":(pos+0.5*tie)/(pos+bad+tie),
                            "matched_pairs":pos+bad+tie,"tie_count":tie,
                            "cells":cells})
    by={(r["arm"],r["condition"]):r["paired_preference_accuracy"] for r in results}
    assert all(by[a,"identity"]==by[a,"xor55"] for a in ARMS)
    gain=by[ARMS[1],"identity"]-by[ARMS[0],"identity"]
    if by[ARMS[1],"identity"]>=.65 and gain>=.10:
        label="PAIRED_RANK_GAIN"
    elif gain>=.05:
        label="PARTIAL_GAIN"
    else:
        label="NO_GAIN"
    return {"schema":"wingless.w-r43-paired-relative-scoring.v1",
            "parent_head":"d75b6bb35a037abc9613282f6375d66b0c0eff1b",
            "classification":label,"development_seeds":TRAIN,
            "sealed_heldout_seeds":HELDOUT,
            "identical_development_feature_pairs":train_identical,
            "paired_minus_independent_ranking_accuracy_gain":gain,
            "results":results,"model_calls":0,
            "claim_limit":"Diagnostic two-alternative paired ranking task only; comparison interface is NOT single-candidate cognition."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
