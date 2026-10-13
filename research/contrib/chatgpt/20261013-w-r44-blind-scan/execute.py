"""W-R44: one observed candidate, fixed local alternatives, no true-site hint."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r43-paired-ranking"/"execute.py"
spec=importlib.util.spec_from_file_location("w43_pinned",SOURCE)
assert spec and spec.loader
w43=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w43)
w37=w43.w37
w42=w43.w42

TRAIN=tuple(range(2224,2240))
HELDOUT=tuple(range(2240,2256))
GENERATORS=w43.GENERATORS
TEST_SITES=w43.TEST_SITES
CONDITIONS=("identity","xor55")
ARMS=("known_site_oracle","blind_any_better","blind_majority")
LOCAL_SITES=tuple(range(14))

def switched(candidate,site):
    assert len(candidate)==16 and 0<=site<14
    tail=list(candidate[1:])
    tail[site],tail[site+1]=tail[site+1],tail[site]
    return candidate[:1]+tuple(tail)

def fit():
    previous=w43.TRAIN
    try:
        w43.TRAIN=TRAIN
        singles,pairs,identical=w43.training_data()
    finally:
        w43.TRAIN=previous
    assert len(singles)==len(pairs)==1728
    return w43.fit_pairs(pairs),identical

def decide(arm,base,candidate,model,true_site):
    feature=w42.feature(base,candidate)
    sites=(true_site,) if arm=="known_site_oracle" else LOCAL_SITES
    greater=equal=total=0
    for site in sites:
        alt=switched(candidate,site)
        other=w42.feature(base,alt)
        delta=tuple(a-b for a,b in zip(feature,other))
        margin=w43.score(model,delta)
        # Count genuinely stronger alternatives only, rather than smuggling
        # the true site, candidate class, or a known partner to blind methods.
        greater+=int(margin< -1e-10)
        equal+=int(abs(margin)<=1e-10)
        total+=1
    if arm=="blind_majority":
        prediction=greater<7
    else:
        prediction=greater==0
    return prediction,greater,equal,total

def evaluate(arm,model,condition,generator,site):
    by_seed=[]
    tp=tn=fp=fn=identical=equal_comparisons=compared=0
    for seed in HELDOUT:
        stp=stn=sfp=sfn=same=sties=scomparisons=0
        for entity in range(9):
            base,pos,neg=w37.recoded_example(seed,entity,generator,condition,site,3,False)
            plain_base,plain_pos,plain_neg=w37.recoded_example(
                seed,entity,generator,"identity",site,3,False)
            assert pos==switched(neg,site)
            assert w42.feature(base,pos)==w42.feature(plain_base,plain_pos)
            assert w42.feature(base,neg)==w42.feature(plain_base,plain_neg)
            same+=int(pos==neg)
            for truth,candidate in ((True,pos),(False,neg)):
                predicted,_,ties,calls=decide(arm,base,candidate,model,site)
                sties+=ties
                scomparisons+=calls
                if truth and predicted:stp+=1
                elif not truth and not predicted:stn+=1
                elif not truth:sfp+=1
                else:sfn+=1
        assert stp+stn+sfp+sfn==18
        tp+=stp;tn+=stn;fp+=sfp;fn+=sfn
        identical+=same;equal_comparisons+=sties;compared+=scomparisons
        by_seed.append({"seed":seed,"tp":stp,"tn":stn,"fp":sfp,"fn":sfn,
                        "identical_observation_pairs":same,
                        "comparisons":scomparisons,"equal_margin_comparisons":sties})
    return {"generator":generator,"site":site,
            "tp":tp,"tn":tn,"fp":fp,"fn":fn,
            "correct":tp+tn,"total":tp+tn+fp+fn,
            "balanced_accuracy":(tp+tn)/(tp+tn+fp+fn),
            "observation_identical_pairs":identical,
            "total_model_comparisons":compared,
            "equal_margin_comparisons":equal_comparisons,"by_seed":by_seed}

def run():
    model,train_identical=fit()
    results=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[evaluate(arm,model,condition,g,site)
                   for g in GENERATORS for site in TEST_SITES]
            results.append({"arm":arm,"condition":condition,
                            "train_examples":1728,"epochs":32,
                            "width":16,"training_updates":model["updates"],
                            "accuracy":sum(c["correct"] for c in cells)/
                                       sum(c["total"] for c in cells),
                            "model_comparisons":sum(c["total_model_comparisons"] for c in cells),
                            "cells":cells})
    by={(r["arm"],r["condition"]):r for r in results}
    assert all(by[arm,"identity"]["accuracy"]==by[arm,"xor55"]["accuracy"]
               for arm in ARMS)
    oracle=by["known_site_oracle","identity"]["accuracy"]
    better=max(by[arm,"identity"]["accuracy"] for arm in ARMS if arm!="known_site_oracle")
    if better>=.65 and better>=oracle-.10:
        label="BLIND_CONTRASTIVE_GAIN"
    elif better>=.55:
        label="PARTIAL_GAIN"
    else:
        label="NO_GAIN"
    return {"schema":"wingless.w-r44-blind-contrastive-counterfactuals.v1",
            "parent_head":"4721050ee908cb8461f26ec2ea1ac30588a59228",
            "classification":label,"development_seeds":TRAIN,
            "sealed_heldout_seeds":HELDOUT,"local_scan_sites":LOCAL_SITES,
            "training_identical_pairs":train_identical,
            "best_blind_accuracy":better,"known_site_oracle_accuracy":oracle,
            "results":results,"model_calls":0,
            "claim_limit":"Blind candidate-local swap enumeration is mechanism-specific, with 14x inference alternatives; true-site oracle is nondeployable."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
