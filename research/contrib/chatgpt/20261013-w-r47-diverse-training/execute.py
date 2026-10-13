"""W-R47: prospective equal-example single-family vs mixed-family two-order ranking."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r46-cluster-replication"/"execute.py"
spec=importlib.util.spec_from_file_location("w46_pinned",SOURCE)
assert spec and spec.loader
w46=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w46)
w37=w46.w37
w42=w46.w42
w45=w46.w45
TRAIN=tuple(range(2560,2576))
HELDOUT=tuple(range(2576,2608))
GENERATORS=w46.GENERATORS
SITES=w46.SITES
TRAIN_SITES=w46.w43.TRAIN_SITES
ARMS=("affine_only","diverse_mix")
CONDITIONS=("identity","xor55")

def fit(arm):
    examples=[]
    source_counts={g:0 for g in GENERATORS}
    identical=0
    for seed in TRAIN:
        for index,generator in enumerate(GENERATORS):
            if arm=="affine_only" and index!=0:continue
            for entity in range(9):
                if arm=="diverse_mix" and entity%3!=(index+seed)%3:continue
                for site in TRAIN_SITES:
                    for amplitude in (1,2):
                        base,pos,neg=w37.recoded_example(
                            seed,entity,generator,"identity",site,amplitude,True)
                        a=w42.feature(base,pos);b=w42.feature(base,neg)
                        identical+=int(a==b)
                        delta=tuple(x-y for x,y in zip(a,b))
                        examples.append((delta,1))
                        examples.append((tuple(-v for v in delta),-1))
                        source_counts[generator]+=2
    assert len(examples)==1728
    if arm=="affine_only":
        assert source_counts=={"affine_prior":1728,"quadratic_irregular":0,"xorshift_independent":0}
    else:
        assert all(n==576 for n in source_counts.values())
    weights=[0.0]*16
    updates=0
    for _ in range(32):
        for x,label in examples:
            if label*sum(w*v for w,v in zip(weights,x))<=0:
                weights=[w+.1*label*v for w,v in zip(weights,x)]
                updates+=1
    return {"weights":weights,"updates":updates,"source_counts":source_counts,
            "identical_train_pairs":identical}

def evaluate(arm,model,condition,generator,site):
    records=[]
    correct=total=tied=0.0
    for seed in HELDOUT:
        scored=0.0
        ties=identical=0
        for entity in range(9):
            base,pos,neg=w37.recoded_example(seed,entity,generator,condition,site,3,False)
            rawbase,rawpos,rawneg=w37.recoded_example(seed,entity,generator,"identity",site,3,False)
            a=w42.feature(base,pos);b=w42.feature(base,neg)
            assert a==w42.feature(rawbase,rawpos)
            assert b==w42.feature(rawbase,rawneg)
            identical+=int(a==b)
            score=sum(w*(x-y) for w,x,y in zip(model["weights"],a,b))
            negscore=-score
            assert abs(score+negscore)<=1e-10
            scored+=w45.credit(score,True)+w45.credit(negscore,False)
            ties+=2*int(abs(score)<=1e-10)
        correct+=scored
        total+=18
        tied+=ties
        records.append({"seed":seed,"correct_credit":scored,"presentations":18,
                        "tied_presentations":ties,"identical_feature_pairs":identical})
    assert total==32*9*2
    return {"arm":arm,"generator":generator,"site":site,
            "correct_credit":correct,"presentations":total,
            "accuracy":correct/total,"tied_presentations":tied,"by_seed":records}

def run():
    trained={arm:fit(arm) for arm in ARMS}
    results=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[evaluate(arm,trained[arm],condition,g,s)
                   for g in GENERATORS for s in SITES]
            families={g:{"correct_credit":sum(c["correct_credit"] for c in cells if c["generator"]==g),
                         "presentations":sum(c["presentations"] for c in cells if c["generator"]==g)}
                      for g in GENERATORS}
            for val in families.values():
                val["accuracy"]=val["correct_credit"]/val["presentations"]
            results.append({"arm":arm,"condition":condition,
                            "development_example_count":1728,"epochs":32,
                            "weights":16,"training_updates":trained[arm]["updates"],
                            "training_source_counts":trained[arm]["source_counts"],
                            "training_feature_ties":trained[arm]["identical_train_pairs"],
                            "balanced_accuracy":sum(c["correct_credit"] for c in cells)/
                                                sum(c["presentations"] for c in cells),
                            "family":families,"cells":cells})
    by={(r["arm"],r["condition"]):r for r in results}
    assert all(by[arm,"identity"]["balanced_accuracy"]==by[arm,"xor55"]["balanced_accuracy"]
               for arm in ARMS)
    src=by["affine_only","identity"];mixed=by["diverse_mix","identity"]
    floors=[min(r["family"][g]["accuracy"] for g in GENERATORS) for r in (src,mixed)]
    worst_gain=floors[1]-floors[0]
    harms={g:mixed["family"][g]["accuracy"]-src["family"][g]["accuracy"] for g in GENERATORS}
    if worst_gain>=.05 and mixed["balanced_accuracy"]>=.55 and min(harms.values())>=-.05:
        label="DIVERSITY_MINIMAX_GAIN"
    elif worst_gain>=.02:
        label="PARTIAL_GAIN"
    else:
        label="NO_GAIN"
    return {"schema":"wingless.w-r47-equal-budget-family-transfer.v1",
            "parent_head":"e2629063bb0057c2ea9062e3dc7ab69ec5636585",
            "classification":label,"development_seeds":TRAIN,
            "sealed_heldout_seeds":HELDOUT,"family_worst_accuracy_gain":worst_gain,
            "family_accuracy_differences":harms,"results":results,
            "model_calls":0,
            "claim_limit":"Synthetic two-candidate generator-mixture supervision only; exposure and entity distribution differ even when example counts match."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
