"""W-R24: frozen partial-channel dropout of original W-R22 models."""
import importlib.util
import json
from pathlib import Path

SRC=Path(__file__).resolve().parent.parent/"20261013-w-r22-maskaware"/"execute.py"
spec=importlib.util.spec_from_file_location("w22_frozen",SRC)
assert spec and spec.loader
w22=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w22)

SEEDS=tuple(range(288,304))
GENERATORS=w22.w20.GENERATORS
SWAPS=("swap10_11","swap12_13")
ARMS=w22.ARMS
CASES=("none","irrelevant_full_rows_2_5","single_channel_zero_2_5",
       "four_channels_separate_ticks_2_5_8_10","both_swapped_full_rows")
ZERO=(0,0,0,0)


def drop(rows,case,swap):
    result=[list(row) for row in rows]
    if case=="irrelevant_full_rows_2_5":
        for t in (2,5):result[t]=list(ZERO)
    elif case=="single_channel_zero_2_5":
        for t in (2,5):result[t][0]=0
    elif case=="four_channels_separate_ticks_2_5_8_10":
        for channel,t in enumerate((2,5,8,10)):result[t][channel]=0
    elif case=="both_swapped_full_rows":
        for t in ((11,12) if swap=="swap10_11" else (13,14)):
            result[t]=list(ZERO)
    else:
        assert case=="none"
    return tuple(tuple(row) for row in result)


def pair(seed,entity,generator,swap,case):
    original=w22.w20.base_sequence(seed,entity,generator)
    measured=w22.w16.observed(original,seed,entity,3,79)
    reordered=w22.w19.permute(measured,swap)
    assert reordered!=measured
    positive=drop(measured,case,swap)
    negative=drop(reordered,case,swap)
    same=(positive==negative)
    assert same==(case=="both_swapped_full_rows")
    assert positive[0]==negative[0]
    return original,positive,negative,same


def test(arm,model,generator,swap,case):
    weights,bias,updates=model
    by_seed=[]
    identical=0
    for seed in SEEDS:
        tp=tn=fp=fn=0
        for entity in range(9):
            a,p,n,same=pair(seed,entity,generator,swap,case)
            identical+=int(same)
            for positive,candidate in ((True,p),(False,n)):
                f=w22.featurize(a,candidate,arm)
                predicted=bias+sum(w*x for w,x in zip(weights,f))>=0
                if positive and predicted:tp+=1
                elif not positive and not predicted:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(s["tp"]+s["tn"] for s in by_seed)
    if case=="both_swapped_full_rows":
        assert identical==16*9
        assert correct==16*9
    else:
        assert identical==0
    return {"generator":generator,"swap":swap,"case":case,
            "correct":correct,"total":16*18,"accuracy":correct/(16*18),
            "identical_positive_negative_pairs":identical,
            "heldout_by_seed":by_seed}


def run():
    models=[]
    for arm in ARMS:
        model=w22.train(arm)
        cells=[test(arm,model,gen,swap,case)
               for gen in GENERATORS for swap in SWAPS for case in CASES]
        assert len(cells)==30
        partial=[c for c in cells if c["case"] in (
            "single_channel_zero_2_5","four_channels_separate_ticks_2_5_8_10")]
        models.append({"arm":arm,"width":16,"epochs":32,"train_pairs":1728,
                       "train_updates":model[2],
                       "overall_accuracy":sum(c["correct"] for c in cells)/
                                          sum(c["total"] for c in cells),
                       "partial_dropout_accuracy":sum(c["correct"] for c in partial)/
                                                   sum(c["total"] for c in partial),
                       "cells":cells})
    by={m["arm"]:m for m in models}
    target=by["maskaware_augmented"]
    baseline=by["naive_augmented"]
    gain=target["partial_dropout_accuracy"]-baseline["partial_dropout_accuracy"]
    partial_cells=[c for c in target["cells"] if c["case"] in (
        "single_channel_zero_2_5","four_channels_separate_ticks_2_5_8_10")]
    if all(c["accuracy"]>=.80 for c in partial_cells) and gain>=.05:
        classification="PARTIAL_CHANNEL_TRANSFER"
    elif gain>=.05:
        classification="MIXED"
    else:
        classification="NEGATIVE"
    return {"schema":"wingless.w-r24-channel-specific-dropout.v1",
            "parent_head":"8b4afcc73937bd7f3b1f5621e1aaa00ba6d7094e",
            "classification":classification,"gain_partial_dropout":gain,
            "sealed_heldout_seeds":SEEDS,"models":models,"model_calls":0,
            "claim_limit":"Synthetic channel erasure falsification, not general intelligence"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
