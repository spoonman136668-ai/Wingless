"""W-R21: frozen missing-tick falsification on W-R19's fixed trained models."""
import importlib.util
import json
from pathlib import Path

PRIOR = Path(__file__).resolve().parent.parent / "20261013-w-r20-generator-shift" / "execute.py"
_spec = importlib.util.spec_from_file_location("w20_frozen", PRIOR)
assert _spec and _spec.loader
w20 = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w20)

HELDOUT = tuple(range(208, 224))
GENERATORS = w20.GENERATORS
ARMS = w20.ARMS
SWAPS = w20.FAMILIES
MASK_ORDER = (2, 5, 8, 10, 3, 6)
OCCLUSIONS = (0, 2, 4, 6)


def example(seed, entity, generator, masked, swap, positive):
    base = w20.base_sequence(seed, entity, generator)
    observed = w20.w19.w16.observed(base, seed, entity, 3, 79)
    mask = set(MASK_ORDER[:masked])
    assert not mask.intersection({0, 11, 12, 13, 14})
    damaged = tuple((0,0,0,0) if tick in mask else row
                    for tick, row in enumerate(observed))
    transformed = damaged if positive else w20.w19.permute(damaged, swap)
    assert transformed[0] == damaged[0]
    assert sorted(transformed) == sorted(damaged)
    assert (transformed == damaged) == positive
    return base, transformed


def evaluate(kind, model, generator, mask_count, swap):
    weights, bias, updates = model
    rows = []
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            a_true,b_true = example(seed, entity, generator, mask_count, swap, True)
            a_false,b_false = example(seed, entity, generator, mask_count, swap, False)
            assert a_true == a_false
            assert b_true[0] == b_false[0] and sorted(b_true) == sorted(b_false)
            for positive,a,b in ((True,a_true,b_true),(False,a_false,b_false)):
                vec=w20.w19.feature(a,b,kind)
                assert len(vec)==16
                pred = bias+sum(w*x for w,x in zip(weights,vec)) >=0
                if positive and pred:tp+=1
                elif not positive and not pred:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn == 18
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(r["tp"]+r["tn"] for r in rows)
    return {"generator":generator,"occluded_ticks":mask_count,"swap":swap,
            "accuracy":correct/288,"correct":correct,"total":288,
            "heldout_by_seed":rows}


def run():
    arms=[]
    for kind in ARMS:
        model=w20.w19.train("single_adjacent",kind)
        cells=[evaluate(kind,model,generator,masked,swap)
               for generator in GENERATORS
               for masked in OCCLUSIONS
               for swap in SWAPS]
        assert len(cells)==24
        hard=[c for c in cells if c["occluded_ticks"]>=4]
        arms.append({
            "model":kind,"width":16,"training_examples":1728,"epochs":32,
            "training_updates":model[2],
            "aggregate_accuracy":sum(c["correct"] for c in cells)/sum(c["total"] for c in cells),
            "hard_occlusion_accuracy":sum(c["correct"] for c in hard)/sum(c["total"] for c in hard),
            "cells":cells
        })
    d={a["model"]:a for a in arms}
    top=d["mean_top_two_64"]
    first=d["first_channel_max"]
    hard=[c for c in top["cells"] if c["occluded_ticks"]>=4]
    gain=top["hard_occlusion_accuracy"]-first["hard_occlusion_accuracy"]
    if all(c["accuracy"]>=.8 for c in hard) and gain>=.05:
        classification="ROBUST_OCCLUSION"
    elif gain>=.03:
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {"schema":"wingless.w-r21-fixed-model-occlusion.v1",
            "parent_head":"52461b46db11ad7677ba058d20ffb52ccd8aa68c",
            "classification":classification,"hard_gain_top_vs_first":gain,
            "training_seeds":list(range(160,176)),
            "heldout_seeds":HELDOUT,"arms":arms,"model_calls":0,
            "claim_limit":"Synthetic fixed-model missing-tick transfer only"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
