"""W-R48 frozen-weight new-cohort two-order ranking transfer; no new model fitting."""
import hashlib
import importlib.util
import json
import random
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r47-diverse-training"/"execute.py"
spec=importlib.util.spec_from_file_location("w47_pinned",SOURCE)
assert spec and spec.loader
w47=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w47)
w46=w47.w46
w42=w47.w42
w37=w47.w37
w45=w47.w45

HELDOUT=tuple(range(2720,2768))
CONDITIONS=("identity","xor55")
GENERATORS=w46.GENERATORS
SITES=w46.SITES
MODELS=("frozen_W46","frozen_W47_affine","frozen_W47_diverse")
BOOTSTRAP_SEED=46048
BOOTSTRAPS=4096

def reconstruct_before_first_use():
    # Source-versioned historical DEV fits. The study never invokes train()
    # or fit() on the prospective heldout seeds or their labels.
    w46weights,updates,identical=w46.fit()
    a=w47.fit("affine_only")
    b=w47.fit("diverse_mix")
    result={
        "frozen_W46":{"weights":w46weights,"historical_updates":updates,
                      "historical_identical_pairs":identical,"historical_dev":"2416..2431"},
        "frozen_W47_affine":{**a,"historical_dev":"2560..2575"},
        "frozen_W47_diverse":{**b,"historical_dev":"2560..2575"}
    }
    for name,row in result.items():
        assert len(row["weights"])==16
        assert all(isinstance(v,(int,float)) for v in row["weights"])
        row["frozen_weights_sha256"]=hashlib.sha256(
            json.dumps(row["weights"],separators=(",",":"),allow_nan=False).encode()).hexdigest()
    return result

def compare(weights,a,b):
    forward=sum(w*(x-y) for w,x,y in zip(weights,a,b))
    reverse=sum(w*(y-x) for w,x,y in zip(weights,a,b))
    assert abs(forward+reverse)<=1e-7
    return w45.credit(forward,True)+w45.credit(reverse,False),int(abs(forward)<=1e-10)*2

def run():
    models=reconstruct_before_first_use()
    results={}
    for condition in CONDITIONS:
        scores={name:{} for name in MODELS}
        for generator in GENERATORS:
            for site in SITES:
                for seed in HELDOUT:
                    for model_name in MODELS:
                        scores[model_name].setdefault(seed,{
                            "correct_credit":0.0,"presentations":0,
                            "tied_presentations":0,
                            "family_correct":{g:0.0 for g in GENERATORS},
                            "family_presentations":{g:0 for g in GENERATORS}})
                    for entity in range(9):
                        base,pos,neg=w37.recoded_example(
                            seed,entity,generator,condition,site,3,False)
                        cleanbase,cleanpos,cleanneg=w37.recoded_example(
                            seed,entity,generator,"identity",site,3,False)
                        a=w42.feature(base,pos);b=w42.feature(base,neg)
                        assert a==w42.feature(cleanbase,cleanpos)
                        assert b==w42.feature(cleanbase,cleanneg)
                        for model_name in MODELS:
                            correct,ties=compare(models[model_name]["weights"],a,b)
                            target=scores[model_name][seed]
                            target["correct_credit"]+=correct
                            target["presentations"]+=2
                            target["tied_presentations"]+=ties
                            target["family_correct"][generator]+=correct
                            target["family_presentations"][generator]+=2
        summary={}
        for model_name in MODELS:
            by_seed=scores[model_name]
            assert len(by_seed)==48
            assert all(row["presentations"]==3*2*9*2 for row in by_seed.values())
            grouped=[by_seed[s]["correct_credit"]/by_seed[s]["presentations"] for s in HELDOUT]
            score=sum(grouped)/48
            rng=random.Random(BOOTSTRAP_SEED)
            simulated=sorted(sum(grouped[rng.randrange(48)] for _ in range(48))/48
                             for i in range(BOOTSTRAPS))
            bounds=[simulated[int(.025*BOOTSTRAPS)],
                    simulated[int(.975*BOOTSTRAPS)-1]]
            families={g:{"accuracy":sum(by_seed[s]["family_correct"][g] for s in HELDOUT)/
                                   sum(by_seed[s]["family_presentations"][g] for s in HELDOUT),
                         "correct_credit":sum(by_seed[s]["family_correct"][g] for s in HELDOUT),
                         "presentations":sum(by_seed[s]["family_presentations"][g] for s in HELDOUT)}
                      for g in GENERATORS}
            summary[model_name]={"balanced_accuracy":score,
                                 "seed_cluster_bootstrap_interval":bounds,
                                 "families":families,
                                 "by_seed":[{"seed":s,**by_seed[s]} for s in HELDOUT],
                                 "frozen_weights_sha256":models[model_name]["frozen_weights_sha256"]}
        results[condition]=summary
    for name in MODELS:
        a=results["identity"][name]
        b=results["xor55"][name]
        assert a["balanced_accuracy"]==b["balanced_accuracy"]
        assert a["seed_cluster_bootstrap_interval"]==b["seed_cluster_bootstrap_interval"]
        assert a["families"]==b["families"]
    primary=results["identity"]["frozen_W46"]
    eligible=sum(primary["families"][g]["accuracy"]>.52 for g in GENERATORS)
    if (primary["balanced_accuracy"]>=.55 and
        primary["seed_cluster_bootstrap_interval"][0]>.5 and eligible>=2):
        classification="FROZEN_COHORT_TRANSFER"
    elif primary["balanced_accuracy"]>=.53:
        classification="PARTIAL_TRANSFER"
    else:
        classification="NO_TRANSFER"
    return {
        "schema":"wingless.w-r48-frozen-cross-cohort.v1",
        "parent_head":"1dc0a393aee9b23bc879cfe7d87f7cb78864d066",
        "classification":classification,
        "sealed_heldout_seeds":HELDOUT,
        "historical_dev_models":models,
        "new_training_updates":0,
        "bootstrap_seed":BOOTSTRAP_SEED,"bootstrap_resamples":BOOTSTRAPS,
        "primary_generators_above_52_percent":eligible,
        "results":results,"model_calls":0,
        "claim_limit":"Fresh synthetic fixed-weight two-alternative ranking; not source-born real-world transfer, independent semantic cognition or autonomous model improvement."
    }

if __name__=="__main__":
    print(json.dumps(run(),indent=2,sort_keys=True))
