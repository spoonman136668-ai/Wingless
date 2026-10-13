"""W-R35: frozen W-R34 models under four independent fresh heldout cohorts."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r34-codec-alignment"/"execute.py"
spec=importlib.util.spec_from_file_location("w34_pinned",SOURCE)
assert spec and spec.loader
w34=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w34)

COHORTS=tuple(tuple(range(start,start+16)) for start in (1440,1456,1472,1488))
BITS=(4,2)
ARMS=w34.ARMS
GENERATORS=w34.GENERATORS
SITES=w34.SITES


def cell(arm,bits,model,generator,site,seeds):
    weights,bias,updates,identical_training=model
    by_seed=[]
    correct=total=identical_pairs=0
    for seed in seeds:
        tp=tn=fp=fn=same=0
        for entity in range(9):
            base,positive,negative=w34.w32.observed_pair(seed,entity,generator,bits,site)
            same+=int(positive==negative)
            for positive_truth,candidate in ((True,positive),(False,negative)):
                feat=w34.representation(base,candidate,arm,bits)
                predicted=bias+sum(w*x for w,x in zip(weights,feat))>=0
                if positive_truth and predicted: tp+=1
                elif not positive_truth and not predicted: tn+=1
                elif not positive_truth: fp+=1
                else: fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        identical_pairs+=same
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "observation_identical_pairs":same})
    return {"generator":generator,"swap_site":site,"correct":correct,"total":total,
            "accuracy":correct/total,"identical_pairs":identical_pairs,"by_seed":by_seed}


def run():
    models={(bits,arm):w34.train(arm,bits) for bits in BITS for arm in ARMS}
    rows=[]
    deltas={}
    for bits in BITS:
        deltas[str(bits)]=[]
        for index,seeds in enumerate(COHORTS):
            comparisons=[]
            for arm in ARMS:
                model=models[(bits,arm)]
                cells=[cell(arm,bits,model,g,site,seeds)
                       for g in GENERATORS for site in SITES]
                assert len(cells)==6
                comparisons.append({"arm":arm,"bits":bits,"width":16,
                                    "training_examples":1728,"epochs":32,
                                    "training_updates":model[2],
                                    "training_identical_pairs":model[3],
                                    "correct":sum(x["correct"] for x in cells),
                                    "total":sum(x["total"] for x in cells),
                                    "accuracy":sum(x["correct"] for x in cells)/
                                               sum(x["total"] for x in cells),
                                    "cells":cells})
            by={x["arm"]:x for x in comparisons}
            delta=by["symmetric_refit"]["accuracy"]-by["asymmetric_refit"]["accuracy"]
            deltas[str(bits)].append(delta)
            rows.append({"cohort_index":index,"heldout_seeds":seeds,
                         "precision_bits":bits,"delta_symmetric_minus_asymmetric":delta,
                         "arms":comparisons})
    pooled_4=sum(deltas["4"])/len(deltas["4"])
    if sum(x<=-.05 for x in deltas["4"])>=3 and pooled_4<=-.05:
        label="ROBUST_ASYMMETRIC_ADVANTAGE"
    elif any(x>0 for x in deltas["4"]) and any(x<0 for x in deltas["4"]):
        label="MIXED_TRANSFER"
    else:
        label="NO_REPLICATION"
    return {"schema":"wingless.w-r35-independent-cohort-replication.v1",
            "parent_head":"99cfa8607008a2616660d05d45e5022e757f02c2",
            "classification":label,
            "original_frozen_training_seeds":w34.TRAIN,
            "new_heldout_cohorts":COHORTS,
            "cohort_deltas":deltas,"pooled_four_bit_delta":pooled_4,
            "rows":rows,"model_calls":0,
            "claim_limit":"Only synthetic independent seed transfer for frozen code/weights; no real-language transfer."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
