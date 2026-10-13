"""W-R50: order-invariant full-batch versus online updates on new sealed synthetic cohorts."""
import hashlib
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r49-order-sensitivity"/"execute.py"
spec=importlib.util.spec_from_file_location("w49_source",SOURCE)
assert spec and spec.loader
w49=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w49)
w37=w49.w37
w42=w49.w42
w45=w49.w45

DEV=tuple(range(3040,3056))
HELDOUT=tuple(range(3056,3104))
ARMS=("online_original","online_reverse","order_invariant_full_batch")
CONDITIONS=("identity","xor55")
GEN=w49.GEN
SITES=w49.SITES

def train_models():
    before=w49.DEV
    try:
        w49.DEV=DEV
        examples,same=w49.historical_development_examples()
    finally:
        w49.DEV=before
    assert len(examples)==1728
    treatments={"online_original":examples,"online_reverse":list(reversed(examples))}
    records={}
    for arm in ARMS:
        weights=[0.0]*16
        updates=0
        batch_nonzero=0
        for epoch in range(32):
            if arm!="order_invariant_full_batch":
                for vec,label in treatments[arm]:
                    margin=label*sum(w*x for w,x in zip(weights,vec))
                    if margin<=0:
                        weights=[w+0.1*label*x for w,x in zip(weights,vec)]
                        updates+=1
            else:
                # Read all candidate margins from the SAME pre-epoch weights.
                # Fixed-width integer gradient has exact commutative sum;
                # never adapt within the epoch.
                gradient=[0]*16
                mistaken=0
                for vec,label in examples:
                    if label*sum(w*x for w,x in zip(weights,vec))<=0:
                        for j in range(16): gradient[j]+=label*vec[j]
                        mistaken+=1
                weights=[w+(0.1/len(examples))*g for w,g in zip(weights,gradient)]
                updates+=mistaken
                batch_nonzero+=int(any(gradient))
        records[arm]={
            "weights":weights,
            "misclassification_updates":updates,
            "nonzero_batch_epochs":batch_nonzero,
            "weights_sha256":hashlib.sha256(json.dumps(weights,separators=(",",":")).encode()).hexdigest()}
    return records,hashlib.sha256(json.dumps([[list(v),label] for v,label in examples],separators=(",",":")).encode()).hexdigest(),same

def run():
    models,corpus_sha,identical=train_models()
    results={}
    for condition in CONDITIONS:
        summary={}
        for arm in ARMS:
            by_seed=[]
            family={g:{"credit":0.0,"presentations":0} for g in GEN}
            for seed in HELDOUT:
                per={"seed":seed,"credit":0.0,"presentations":0,"ties":0,"family_correct":{}}
                for g in GEN:
                    subtotal=0.0
                    for site in SITES:
                        for entity in range(9):
                            base,pos,neg=w37.recoded_example(seed,entity,g,condition,site,3,False)
                            origin,op,on=w37.recoded_example(seed,entity,g,"identity",site,3,False)
                            va=w42.feature(base,pos);vb=w42.feature(base,neg)
                            assert va==w42.feature(origin,op) and vb==w42.feature(origin,on)
                            margin=sum(w*(a-b) for w,a,b in zip(models[arm]["weights"],va,vb))
                            credited=w45.credit(margin,True)+w45.credit(-margin,False)
                            per["credit"]+=credited
                            per["presentations"]+=2
                            per["ties"]+=2*int(abs(margin)<=1e-10)
                            family[g]["credit"]+=credited
                            family[g]["presentations"]+=2
                            subtotal+=credited
                    per["family_correct"][g]=subtotal
                assert per["presentations"]==108
                by_seed.append(per)
            family_acc={g:family[g]["credit"]/family[g]["presentations"] for g in GEN}
            summary[arm]={"accuracy":sum(s["credit"] for s in by_seed)/(48*108),
                          "family_accuracy":family_acc,
                          "by_seed":by_seed,
                          "training_updates":models[arm]["misclassification_updates"],
                          "nonzero_batch_epochs":models[arm]["nonzero_batch_epochs"],
                          "weights_sha256":models[arm]["weights_sha256"]}
        results[condition]=summary
    for arm in ARMS:
        assert results["identity"][arm]["accuracy"]==results["xor55"][arm]["accuracy"]
        assert results["identity"][arm]["family_accuracy"]==results["xor55"][arm]["family_accuracy"]
    scores=results["identity"]
    batch=scores["order_invariant_full_batch"]["accuracy"]
    mean_online=(scores["online_original"]["accuracy"]+scores["online_reverse"]["accuracy"])/2
    gain=batch-mean_online
    good_families=sum(scores["order_invariant_full_batch"]["family_accuracy"][g]>.52 for g in GEN)
    if batch>=.55 and gain>=.05 and good_families>=2:
        classification="BATCH_STABLE_USEFUL"
    elif gain>=.02:
        classification="PARTIAL_BATCH_GAIN"
    else:
        classification="NO_GAIN"
    return {"schema":"wingless.w-r50-order-invariant-batch.v1",
            "parent_head":"1ea9d0b4cc218c2fedb1ab0fcc47ed93ad5876a1",
            "classification":classification,"development_seeds":DEV,
            "sealed_heldout_seeds":HELDOUT,
            "original_training_corpus_sha256":corpus_sha,
            "identical_development_feature_pairs":identical,
            "training_examples_per_arm":1728,"training_epochs":32,
            "batch_minus_online_mean_accuracy":gain,
            "batch_source_families_above_52pct":good_families,
            "results":results,"model_calls":0,
            "claim_limit":"Exact-order-invariant synthetic fixed-width full-batch optimization compared with online training; equal epoch/data/parameter budgets, different update semantics, two-candidate comparison only."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
