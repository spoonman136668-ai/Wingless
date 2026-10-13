"""W-R23: masks applied after ordering, finite observational identifiability test."""
import importlib.util
import json
from pathlib import Path

SOURCE = Path(__file__).resolve().parent.parent / "20261013-w-r22-maskaware" / "execute.py"
_spec = importlib.util.spec_from_file_location("w22_frozen", SOURCE)
assert _spec and _spec.loader
w22 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w22)

HELDOUT = tuple(range(272,288))
GENERATORS = w22.w20.GENERATORS
SWAPS = ("swap10_11","swap12_13")
CASES = ("none","irrelevant_2_5","one_swap_row","both_swap_rows")
ARMS = w22.ARMS
ZERO=(0,0,0,0)


def swap_rows(swap):
    assert swap in SWAPS
    # W19's tail indices are zero-based after the preserved tick-0 anchor.
    return (11,12) if swap=="swap10_11" else (13,14)


def observed_pair(seed,entity,generator,swap,case):
    original=w22.w20.base_sequence(seed,entity,generator)
    measured=w22.w16.observed(original,seed,entity,3,79)
    altered=w22.w19.permute(measured,swap)
    assert measured!=altered
    if case=="none":
        missing=set()
    elif case=="irrelevant_2_5":
        missing={2,5}
    elif case=="one_swap_row":
        missing={swap_rows(swap)[0]}
    else:
        assert case=="both_swap_rows"
        missing=set(swap_rows(swap))
    assert 0 not in missing
    positive=tuple(ZERO if tick in missing else row
                   for tick,row in enumerate(measured))
    negative=tuple(ZERO if tick in missing else row
                   for tick,row in enumerate(altered))
    equivalent=(positive==negative)
    assert equivalent==(case=="both_swap_rows")
    return original,positive,negative,equivalent


def test(arm,trained,generator,swap,case):
    weights,bias,_=trained
    by_seed=[]
    ambiguous=0
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        eq_count=0
        for entity in range(9):
            original,positive,negative,equivalent=observed_pair(seed,entity,generator,swap,case)
            eq_count+=int(equivalent)
            for truth,candidate in ((True,positive),(False,negative)):
                features=w22.featurize(original,candidate,arm)
                assert len(features)==16
                pred=(bias+sum(w*x for w,x in zip(weights,features)))>=0
                if truth and pred:tp+=1
                elif not truth and not pred:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "equivalent_positive_negative_pairs":eq_count})
        ambiguous+=eq_count
    correct=sum(r["tp"]+r["tn"] for r in by_seed)
    return {"generator":generator,"swap":swap,"case":case,
            "accuracy":correct/(16*18),"correct":correct,"total":16*18,
            "equivalent_pairs":ambiguous,"total_pairs":16*9,
            "theoretical_balanced_accuracy_ceiling":
                (1.0-0.5*ambiguous/(16*9)),
            "heldout_by_seed":by_seed}


def run():
    models=[]
    for arm in ARMS:
        trained=w22.train(arm)
        cells=[test(arm,trained,gen,swap,case)
               for gen in GENERATORS for swap in SWAPS for case in CASES]
        assert len(cells)==24
        models.append({"arm":arm,"trained_width":16,
                       "training_examples":1728,"epochs":32,
                       "training_updates":trained[2],
                       "accuracy":sum(c["correct"] for c in cells)/
                                  sum(c["total"] for c in cells),
                       "cells":cells})
    for model in models:
        for cell in model["cells"]:
            assert (cell["equivalent_pairs"]==cell["total_pairs"]) == (
                cell["case"]=="both_swap_rows")
            if cell["case"]=="both_swap_rows":
                assert cell["accuracy"]==0.5
    by={m["arm"]:m for m in models}
    identified=[c for c in by["maskaware_clean"]["cells"]
                if c["case"]!="both_swap_rows"]
    if all(c["accuracy"]>=.80 for c in identified):
        label="IDENTIFIABILITY_BOUNDARY"
    else:
        label="PARTIAL_BOUNDARY"
    return {"schema":"wingless.w-r23-overlap-occlusion-observability.v1",
            "parent_head":"c551d41a7b8d4f77b9349033fa788d621a800519",
            "classification":label,"development_seeds":list(w22.TRAIN_SEEDS),
            "sealed_heldout_seeds":HELDOUT,"models":models,"model_calls":0,
            "claim_limit":"Synthetic information-theoretic paired observation boundary only"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
