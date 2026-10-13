"""W-R46: prospective 48-seed, seed-cluster bootstrap replication of order-balanced paired ranking."""
import importlib.util
import json
import random
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r45-order-blind"/"execute.py"
spec=importlib.util.spec_from_file_location("w45_pinned",SOURCE)
assert spec and spec.loader
w45=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w45)
w43=w45.w43
w37=w45.w37
w42=w45.w42

TRAIN=tuple(range(2416,2432))
HELDOUT=tuple(range(2432,2480))
CONDITIONS=("identity","xor55")
GENERATORS=w45.GENERATORS
SITES=w45.SITES
BOOTSTRAP_SEED=46046
BOOTSTRAPS=4096

def fit():
    original=w43.TRAIN
    try:
        w43.TRAIN=TRAIN
        _,examples,identical=w43.training_data()
    finally:
        w43.TRAIN=original
    assert len(examples)==1728
    weights=[0.0]*16
    updates=0
    for _ in range(32):
        for x,label in examples:
            margin=sum(w*v for w,v in zip(weights,x))
            if label*margin<=0:
                weights=[w+.1*label*v for w,v in zip(weights,x)]
                updates+=1
    return weights,updates,identical

def run():
    weights,updates,train_identical=fit()
    by_condition={}
    for condition in CONDITIONS:
        by_seed={}
        by_family={}
        for generator in GENERATORS:
            cells=[]
            for site in SITES:
                by_seed_rows=[]
                for seed in HELDOUT:
                    true_first=true_second=0.0
                    tied=identical=0
                    for entity in range(9):
                        base,pos,neg=w37.recoded_example(
                            seed,entity,generator,condition,site,3,False)
                        barebase,barepos,bareneg=w37.recoded_example(
                            seed,entity,generator,"identity",site,3,False)
                        a=w42.feature(base,pos); b=w42.feature(base,neg)
                        assert a==w42.feature(barebase,barepos)
                        assert b==w42.feature(barebase,bareneg)
                        identical+=int(a==b)
                        direct=sum(w*(x-y) for w,x,y in zip(weights,a,b))
                        reverse=sum(w*(y-x) for w,x,y in zip(weights,a,b))
                        assert abs(direct+reverse)<=1e-7
                        true_first+=w45.credit(direct,True)
                        true_second+=w45.credit(reverse,False)
                        tied+=int(abs(direct)<=1e-10)+int(abs(reverse)<=1e-10)
                    assert 0<=true_first<=9 and 0<=true_second<=9
                    by_seed_rows.append({"seed":seed,"true_first_credit":true_first,
                                         "true_second_credit":true_second,
                                         "tied_presentations":tied,
                                         "identical_feature_pairs":identical})
                    target=by_seed.setdefault(seed,{"correct":0.0,"presentations":0,
                                                   "tied":0,"family_correct":{}})
                    target["correct"]+=true_first+true_second
                    target["presentations"]+=18
                    target["tied"]+=tied
                    target["family_correct"][generator]=(
                        target["family_correct"].get(generator,0)+true_first+true_second)
                cells.append({"site":site,"by_seed":by_seed_rows,
                              "correct_credit":sum(r["true_first_credit"]+r["true_second_credit"]
                                                   for r in by_seed_rows),
                              "presentations":48*18})
            by_family[generator]={"cells":cells,
                                  "accuracy":sum(c["correct_credit"] for c in cells)/
                                             sum(c["presentations"] for c in cells)}
        assert len(by_seed)==48
        assert all(row["presentations"]==3*2*9*2 for row in by_seed.values())
        clustered_scores=[by_seed[s]["correct"]/by_seed[s]["presentations"] for s in HELDOUT]
        point=sum(clustered_scores)/len(clustered_scores)
        assert .0<=point<=1.
        bootstrap_rng=random.Random(BOOTSTRAP_SEED)
        resamples=sorted(sum(clustered_scores[bootstrap_rng.randrange(48)]
                             for _ in range(48))/48 for _ in range(BOOTSTRAPS))
        bounds=(resamples[int(.025*BOOTSTRAPS)],resamples[int(.975*BOOTSTRAPS)-1])
        by_condition[condition]={"balanced_accuracy":point,"seed_cluster_interval":bounds,
                                 "bootstrap_seed":BOOTSTRAP_SEED,
                                 "bootstrap_resamples":BOOTSTRAPS,
                                 "family":by_family,
                                 "by_seed":[{"seed":s,**by_seed[s],
                                             "balanced_accuracy":clustered_scores[i]}
                                            for i,s in enumerate(HELDOUT)]}
    original=by_condition["identity"]
    other=by_condition["xor55"]
    assert original["balanced_accuracy"]==other["balanced_accuracy"]
    assert original["seed_cluster_interval"]==other["seed_cluster_interval"]
    assert all(original["family"][g]["accuracy"]==other["family"][g]["accuracy"]
               for g in GENERATORS)
    family_success=sum(original["family"][g]["accuracy"]>.52 for g in GENERATORS)
    if (original["balanced_accuracy"]>=.55 and
        original["seed_cluster_interval"][0]>.50 and family_success>=2):
        label="ROBUST_PARTIAL_REPLICATION"
    elif original["balanced_accuracy"]>=.53:
        label="BORDERLINE_SIGNAL"
    else:
        label="NOT_REPLICATED"
    return {"schema":"wingless.w-r46-cluster-replication.v1",
            "parent_head":"a0733122a2e35cfb05bbd04231847d68abc6dff5",
            "classification":label,"development_seeds":TRAIN,
            "sealed_heldout_seeds":HELDOUT,"training_updates":updates,
            "training_identical_feature_pairs":train_identical,
            "model_weights":weights,"balanced_order_constant_first_null":.5,
            "generator_families_above_52_percent":family_success,
            "results":by_condition,"model_calls":0,
            "claim_limit":"Fresh synthetic two-alternative balanced scoring with seed-cluster confidence interval; not standalone learned cognition."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
