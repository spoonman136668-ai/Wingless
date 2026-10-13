"""W-R49: three exact-example-order treatments on one new 48-seed paired cohort."""
import hashlib
import importlib.util
import json
import random
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r48-frozen-transfer"/"execute.py"
spec=importlib.util.spec_from_file_location("w48_pinned",SOURCE)
assert spec and spec.loader
w48=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w48)

w46=w48.w46
w43=w46.w43
w37=w48.w37
w42=w48.w42
w45=w48.w45

DEV=tuple(range(2864,2880))
HELDOUT=tuple(range(2880,2928))
GEN=w48.GENERATORS
SITES=w48.SITES
ARMS=("original_order","reversed_order","fixed_shuffled_order_49049")
CONDITIONS=("identity","xor55")
BOOTSTRAP_SEED=49050
BOOTSTRAPS=4096


def historical_development_examples():
    old=w43.TRAIN
    try:
        w43.TRAIN=DEV
        _,examples,identical=w43.training_data()
    finally:
        w43.TRAIN=old
    assert len(examples)==1728
    return examples,identical

def fit_all():
    examples,identical=historical_development_examples()
    shuffle=list(examples)
    random.Random(49049).shuffle(shuffle)
    ordered={"original_order":examples,
             "reversed_order":list(reversed(examples)),
             "fixed_shuffled_order_49049":shuffle}
    source_sig=hashlib.sha256(json.dumps(
        [[list(x),y] for x,y in examples],separators=(",",":")).encode()).hexdigest()
    out={}
    for arm in ARMS:
        weights=[0.0]*16
        updates=0
        for _ in range(32):
            for x,label in ordered[arm]:
                margin=sum(w*v for w,v in zip(weights,x))
                if label*margin<=0:
                    weights=[w+.1*label*v for w,v in zip(weights,x)]
                    updates+=1
        assert len(weights)==16
        out[arm]={"weights":weights,"updates":updates,
                  "weight_sha256":hashlib.sha256(
                      json.dumps(weights,separators=(",",":")).encode()).hexdigest()}
    return out,source_sig,identical

def run():
    models,source_sig,identical_train=fit_all()
    results={}
    for condition in CONDITIONS:
        per_arm={arm:{seed:{"correct":0.0,"presentations":0,"ties":0,
                            "family_correct":{g:0.0 for g in GEN},
                            "family_presentations":{g:0 for g in GEN}}
                      for seed in HELDOUT} for arm in ARMS}
        for generator in GEN:
            for site in SITES:
                for seed in HELDOUT:
                    for entity in range(9):
                        base,pos,neg=w37.recoded_example(
                            seed,entity,generator,condition,site,3,False)
                        origbase,origpos,origneg=w37.recoded_example(
                            seed,entity,generator,"identity",site,3,False)
                        a=w42.feature(base,pos);b=w42.feature(base,neg)
                        assert a==w42.feature(origbase,origpos)
                        assert b==w42.feature(origbase,origneg)
                        for arm in ARMS:
                            delta=sum(w*(x-y) for w,x,y in zip(models[arm]["weights"],a,b))
                            reverse=-delta
                            scored=w45.credit(delta,True)+w45.credit(reverse,False)
                            target=per_arm[arm][seed]
                            target["correct"]+=scored
                            target["presentations"]+=2
                            target["ties"]+=2*int(abs(delta)<=1e-10)
                            target["family_correct"][generator]+=scored
                            target["family_presentations"][generator]+=2
        results[condition]={}
        for arm in ARMS:
            by_seed=per_arm[arm]
            assert len(by_seed)==48
            assert all(row["presentations"]==3*2*9*2 for row in by_seed.values())
            scores=[by_seed[seed]["correct"]/by_seed[seed]["presentations"]
                    for seed in HELDOUT]
            mean=sum(scores)/48
            rng=random.Random(BOOTSTRAP_SEED)
            sampled=sorted(sum(scores[rng.randrange(48)] for _ in range(48))/48
                           for _ in range(BOOTSTRAPS))
            family={g:sum(by_seed[seed]["family_correct"][g] for seed in HELDOUT)/
                      sum(by_seed[seed]["family_presentations"][g] for seed in HELDOUT)
                    for g in GEN}
            results[condition][arm]={
                "balanced_accuracy":mean,"family_accuracies":family,
                "seed_cluster_interval":[sampled[int(.025*BOOTSTRAPS)],
                                         sampled[int(.975*BOOTSTRAPS)-1]],
                "by_seed":[{"seed":seed,**by_seed[seed]} for seed in HELDOUT],
                "training_update_count":models[arm]["updates"],
                "weights_sha256":models[arm]["weight_sha256"]}
    for arm in ARMS:
        assert results["identity"][arm]["balanced_accuracy"]==results["xor55"][arm]["balanced_accuracy"]
        assert results["identity"][arm]["family_accuracies"]==results["xor55"][arm]["family_accuracies"]
    primary=results["identity"]
    spread=max(v["balanced_accuracy"] for v in primary.values())-min(v["balanced_accuracy"] for v in primary.values())
    strict=all(row["balanced_accuracy"]>=.55 and
               sum(row["family_accuracies"][g]>.52 for g in GEN)>=2
               for row in primary.values()) and spread<=.05
    partial=(sum(row["balanced_accuracy"]>=.55 for row in primary.values())>=2
             and spread<=.10)
    classification="TRAIN_ORDER_STABLE" if strict else (
        "PARTIAL_ORDER_STABILITY" if partial else "OTHERWISE_UNSTABLE")
    return {"schema":"wingless.w-r49-training-order-sensitivity.v1",
            "parent_head":"28e4b397f6fa9586e9533f59eb84074174932de0",
            "classification":classification,
            "development_seeds":DEV,"sealed_heldout_seeds":HELDOUT,
            "same_development_examples_sha256":source_sig,
            "identical_development_feature_pairs":identical_train,
            "training_example_count_per_arm":1728,
            "accuracy_spread":spread,
            "results":results,"model_calls":0,
            "claim_limit":"Seed-cluster synthetic two-candidate training-order effect, not general standalone cognitive learning."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
