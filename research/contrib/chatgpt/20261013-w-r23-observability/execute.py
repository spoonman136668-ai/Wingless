"""W-R23: first-use observational equivalence at hidden swap sites."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r22-maskaware"/"execute.py"
spec=importlib.util.spec_from_file_location("w22_frozen",SRC)
assert spec and spec.loader
w22=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w22)
w20=w22.w20
w19=w22.w19
w16=w22.w16

HELDOUT=tuple(range(272,288))
ARMS=w22.ARMS
GENERATORS=w20.GENERATORS
SWAPS=("swap10_11","swap12_13")
VISIBILITIES=("visible","one_hidden","both_hidden")
BACKGROUND=(2,5,8)


def pair(seed,entity,generator,swap,visibility):
    assert swap in SWAPS and visibility in VISIBILITIES
    base=w20.base_sequence(seed,entity,generator)
    observed=w16.observed(base,seed,entity,3,79)
    tail_index={"swap10_11":10,"swap12_13":12}[swap]
    target_ticks=(tail_index+1,tail_index+2)
    masked=set(BACKGROUND)
    if visibility in ("one_hidden","both_hidden"):
        masked.add(target_ticks[0])
    if visibility=="both_hidden":
        masked.add(target_ticks[1])
    assert 0 not in masked
    candidate=tuple((0,0,0,0) if tick in masked else row
                    for tick,row in enumerate(observed))
    tail=list(candidate[1:])
    tail[tail_index],tail[tail_index+1]=tail[tail_index+1],tail[tail_index]
    negative=candidate[:1]+tuple(tail)
    assert candidate[0]==negative[0] and sorted(candidate)==sorted(negative)
    assert (candidate==negative)==(visibility=="both_hidden")
    assert all(candidate[t]==(0,0,0,0) for t in masked)
    return base,candidate,negative


def evaluate(arm,trained,generator,swap,visibility):
    weights,bias,updates=trained
    rows=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        equivalences=0
        for entity in range(9):
            base,positive,negative=pair(seed,entity,generator,swap,visibility)
            if positive==negative:
                equivalences+=1
                assert visibility=="both_hidden"
            for truth,candidate in ((True,positive),(False,negative)):
                vec=w22.featurize(base,candidate,arm)
                assert len(vec)==16
                predicted_positive=bias+sum(w*v for w,v in zip(weights,vec))>=0
                if truth and predicted_positive:tp+=1
                elif not truth and not predicted_positive:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        assert equivalences==(9 if visibility=="both_hidden" else 0)
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "paired_observation_equivalences":equivalences})
    correct=sum(x["tp"]+x["tn"] for x in rows)
    return {"generator":generator,"swap":swap,"visibility":visibility,
            "correct":correct,"total":16*18,"accuracy":correct/(16*18),
            "by_seed":rows}


def run():
    arms=[]
    for arm in ARMS:
        trained=w22.train(arm)
        cells=[evaluate(arm,trained,g,s,v)
               for g in GENERATORS for s in SWAPS for v in VISIBILITIES]
        assert len(cells)==18
        arms.append({"arm":arm,"feature_width":16,"training_examples":1728,
                     "epochs":32,"training_updates":trained[2],
                     "accuracy":sum(c["correct"] for c in cells)/(18*288),
                     "cells":cells})
    invis_ok=all(
        c["accuracy"]==.5 and
        all(s["paired_observation_equivalences"]==9 for s in c["by_seed"])
        for arm in arms for c in arm["cells"] if c["visibility"]=="both_hidden")
    visible_ok=all(
        c["accuracy"]>=.85
        for arm in arms if arm["arm"].startswith("maskaware")
        for c in arm["cells"] if c["visibility"]=="visible")
    classification=("OBSERVABILITY_BOUNDARY" if invis_ok and visible_ok else
                    "PARTIAL" if invis_ok else "INVALID")
    return {"schema":"wingless.w-r23-observability-boundary.v1",
            "parent_head":"c551d41a7b8d4f77b9349033fa788d621a800519",
            "classification":classification,"equivalence_control_pass":invis_ok,
            "visible_control_pass":visible_ok,"development_seeds":w22.TRAIN_SEEDS,
            "sealed_heldout_seeds":HELDOUT,"arms":arms,
            "model_calls":0,
            "claim_limit":"Synthetic observation equivalence; invisible label changes cannot be inferred"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
