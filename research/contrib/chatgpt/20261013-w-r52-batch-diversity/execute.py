"""W-R52: fixed-budget three-generator full-batch training across 3 development cohorts."""
import hashlib
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r51-multicohort"/"execute.py"
spec=importlib.util.spec_from_file_location("w51_pinned",SOURCE)
assert spec and spec.loader
w51=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w51)
w47=w51.w50.w49.w48.w47
w37=w47.w37
w42=w47.w42
w45=w47.w46.w45

DEV_BLOCKS=(tuple(range(3424,3440)),tuple(range(3456,3472)),tuple(range(3488,3504)))
HELDOUT=tuple(range(3520,3568))
GEN=w47.GENERATORS
SITES=w47.SITES
TRAIN_SITES=w47.TRAIN_SITES
ARMS=("affine_only_full_batch","three_source_balanced_full_batch")
CONDITIONS=("identity","xor55")

def training_examples(block,arm):
    samples=[]
    counts={g:0 for g in GEN}
    for seed in block:
        for index,generator in enumerate(GEN):
            if arm=="affine_only_full_batch" and index!=0: continue
            for entity in range(9):
                if arm=="three_source_balanced_full_batch" and entity%3!=(index+seed)%3:continue
                for site in TRAIN_SITES:
                    for amp in (1,2):
                        base,pos,neg=w37.recoded_example(
                            seed,entity,generator,"identity",site,amp,True)
                        a=w42.feature(base,pos)
                        b=w42.feature(base,neg)
                        delta=tuple(x-y for x,y in zip(a,b))
                        samples.extend(((delta,1),(tuple(-v for v in delta),-1)))
                        counts[generator]+=2
    assert len(samples)==1728
    if arm=="three_source_balanced_full_batch":
        assert all(counts[g]==576 for g in GEN)
    else:
        assert counts[GEN[0]]==1728 and all(counts[g]==0 for g in GEN[1:])
    fingerprint=hashlib.sha256(json.dumps([[list(x),y] for x,y in samples],
                              separators=(",",":")).encode()).hexdigest()
    return samples,counts,fingerprint

def fit(block,arm):
    examples,counts,corpus_sha=training_examples(block,arm)
    weights=[0.0]*16
    updates=0
    for _ in range(32):
        gradient=[0]*16
        for feature,label in examples:
            if label*sum(w*x for w,x in zip(weights,feature))<=0:
                for j in range(16):gradient[j]+=label*feature[j]
                updates+=1
        weights=[w+(0.1/len(examples))*g for w,g in zip(weights,gradient)]
    return {"weights":weights,"updates":updates,
            "source_counts":counts,"corpus_sha256":corpus_sha,
            "weights_sha256":hashlib.sha256(json.dumps(weights,separators=(",",":")).encode()).hexdigest()}

def run():
    models={(i,arm):fit(block,arm) for i,block in enumerate(DEV_BLOCKS) for arm in ARMS}
    scores={}
    for condition in CONDITIONS:
        by_model={}
        for (idx,arm),model in models.items():
            family={g:{"correct":0.0,"total":0} for g in GEN}
            per_seed=[]
            for seed in HELDOUT:
                row={"seed":seed,"correct":0.0,"total":0,"ties":0,
                     "family_correct":{}}
                for generator in GEN:
                    subtotal=0.0
                    for site in SITES:
                        for entity in range(9):
                            base,pos,neg=w37.recoded_example(
                                seed,entity,generator,condition,site,3,False)
                            unbase,unpos,unneg=w37.recoded_example(
                                seed,entity,generator,"identity",site,3,False)
                            a=w42.feature(base,pos);b=w42.feature(base,neg)
                            assert a==w42.feature(unbase,unpos)
                            assert b==w42.feature(unbase,unneg)
                            margin=sum(w*(x-y) for w,x,y in zip(model["weights"],a,b))
                            credit=w45.credit(margin,True)+w45.credit(-margin,False)
                            subtotal+=credit
                            row["correct"]+=credit
                            row["total"]+=2
                            row["ties"]+=2*int(abs(margin)<=1e-10)
                            family[generator]["correct"]+=credit
                            family[generator]["total"]+=2
                    row["family_correct"][generator]=subtotal
                assert row["total"]==108
                per_seed.append(row)
            by_model[f"{idx+1}:{arm}"]={
                "accuracy":sum(row["correct"] for row in per_seed)/(48*108),
                "family_accuracy":{g:family[g]["correct"]/family[g]["total"] for g in GEN},
                "by_seed":per_seed,
                "weights_sha256":model["weights_sha256"],
                "training_corpus_sha256":model["corpus_sha256"],
                "source_training_example_counts":model["source_counts"],
                "batch_misclassification_updates":model["updates"]}
        scores[condition]=by_model
    assert scores["identity"].keys()==scores["xor55"].keys()
    for key in scores["identity"]:
        assert scores["identity"][key]["accuracy"]==scores["xor55"][key]["accuracy"]
        assert scores["identity"][key]["family_accuracy"]==scores["xor55"][key]["family_accuracy"]
    comparisons=[]
    for idx in range(1,4):
        single=scores["identity"][f"{idx}:affine_only_full_batch"]
        mixed=scores["identity"][f"{idx}:three_source_balanced_full_batch"]
        original_floor=min(single["family_accuracy"].values())
        new_floor=min(mixed["family_accuracy"].values())
        gain=new_floor-original_floor
        harms={g:mixed["family_accuracy"][g]-single["family_accuracy"][g] for g in GEN}
        strict=gain>=.05 and mixed["accuracy"]>=.55 and min(harms.values())>=-.05
        partial=gain>=.02
        comparisons.append({
            "development_block":idx,"min_source_gain":gain,
            "mixed_overall_accuracy":mixed["accuracy"],
            "affine_only_overall_accuracy":single["accuracy"],
            "source_deltas":harms,"strict_gate":strict,"partial_gate":partial})
    if all(c["strict_gate"] for c in comparisons):
        label="FAMILY_DIVERSITY_ROBUST"
    elif sum(c["partial_gate"] for c in comparisons)>=2:
        label="PARTIAL_FAMILY_GAIN"
    else:
        label="NO_CONFIRMED_DIVERSITY_GAIN"
    return {"schema":"wingless.w-r52-family-balanced-batch.v1",
            "parent_head":"300bf3d313a8bc98af120a6e32e046f3646385ad",
            "classification":label,
            "development_blocks":DEV_BLOCKS,"sealed_heldout_seeds":HELDOUT,
            "arm_comparisons":comparisons,
            "all_models_by_encoding":scores,
            "new_model_selection_from_heldout":False,
            "model_calls":0,
            "claim_limit":"Synthetic paired-choice full-batch equal-example generator-mix study; source-family/entity exposure changes by design; no real-world learning claims."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
