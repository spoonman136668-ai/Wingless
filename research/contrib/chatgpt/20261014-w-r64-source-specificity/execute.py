"""W-R64: frozen original-source channel erasure specificity across wiring frames."""
import importlib.util
import json
from pathlib import Path

PARENT = Path(__file__).resolve().parent.parent / "20261013-w-r63-channel-role" / "execute.py"
spec = importlib.util.spec_from_file_location("w63_pinned", PARENT)
assert spec and spec.loader
w63 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w63)

HELDOUT = tuple(range(9664, 9688))
PERMS = w63.PERMS
AMPS = (1, 2, 3)
CONDITIONS = ("intact", "erase_source0", "erase_source1", "erase_source2", "erase_source3")
GENERATOR = "xorshift_independent"
ENTITIES = w63.GROUPS["unseen_identity_transfer"]
SITES = w63.SITES

def acquire(rows, order, condition):
    assert sorted(order) == [0, 1, 2, 3]
    source = None if condition == "intact" else int(condition[-1])
    assert condition in CONDITIONS
    return tuple(tuple(0 if source == channel else row[channel]
                       for channel in order) for row in rows)

def rank(seed, entity, site, amplitude, weights, order, condition):
    base, pos, neg = w63.w61.m.m.w37.recoded_example(
        seed, entity, GENERATOR, "identity", site, amplitude, False)
    base = acquire(base, order, condition)
    pos = acquire(pos, order, condition)
    neg = acquire(neg, order, condition)
    p = w63.w61.m.m.w42.feature(base, pos)
    n = w63.w61.m.m.w42.feature(base, neg)
    assert len(p) == len(n) == len(weights) == 16
    margin = sum(w*(a-b) for w,a,b in zip(weights,p,n))
    return w63.w61.m.m.w45.credit(margin, True) + w63.w61.m.m.w45.credit(-margin, False)

def run():
    trained = w63.w61.fit()
    rows = []
    for i in (1,2,3):
        model = f"{i}:withhold_xorshift_independent"
        fit = trained[model]
        assert fit["training_counts"][GENERATOR] == 0
        weights = fit["weights"]
        assert len(weights) == 16
        for amp in AMPS:
            for name,order in PERMS.items():
                for condition in CONDITIONS:
                    by_seed = []
                    for seed in HELDOUT:
                        credit = sum(rank(seed,entity,site,amp,weights,order,condition)
                                     for entity in ENTITIES for site in SITES)
                        presentations = len(ENTITIES)*len(SITES)*2
                        by_seed.append({"seed":seed,"correct_credit":credit,"presentations":presentations})
                    numerator = sum(r["correct_credit"] for r in by_seed)
                    denominator = sum(r["presentations"] for r in by_seed)
                    rows.append({"model":model,"amplitude":amp,"permutation":name,
                                 "condition":condition,"credit":numerator,
                                 "presentations":denominator,"accuracy":numerator/denominator,
                                 "by_seed":by_seed})
    assert len(rows) == 3*3*4*5
    scores = {(x["model"],x["amplitude"],x["permutation"],x["condition"]):x["accuracy"] for x in rows}
    comparisons = []
    for i in (1,2,3):
        model=f"{i}:withhold_xorshift_independent"
        for amp in AMPS:
            for name in PERMS:
                prefix = (model,amp,name)
                intact=scores[prefix+("intact",)]
                delta={str(source):scores[prefix+(f"erase_source{source}",)]-intact
                       for source in range(4)}
                comparisons.append({"model":model,"amplitude":amp,"permutation":name,
                                    "intact_accuracy":intact,"loss_vs_intact":delta})
    assert len(comparisons)==36
    mean = {str(c):sum(x["loss_vs_intact"][str(c)] for x in comparisons)/36
            for c in range(4)}
    if mean["0"]<=-.07 and mean["0"]<=min(mean[str(c)] for c in (1,2,3))-.035:
        classification="SOURCE0_DOMINANT_REPLICATED"
    elif sum(v<=-.05 for v in mean.values())>=2:
        classification="MULTICHANNEL_SENSITIVITY"
    else:
        classification="NO_CLEAR_SOURCE_RANK"
    return {"schema":"wingless.w-r64-source-channel-specificity.v1",
            "parent_head":"635d8409a175c483dde48b624149ba1f6505d1a7",
            "classification":classification,
            "sealed_heldout_seeds":HELDOUT,"development_blocks":w63.w61.DEV,
            "generator":GENERATOR,"entities":ENTITIES,
            "source_erase_mean_loss":mean,"comparisons":comparisons,"all_sealed_records":rows,
            "model_calls":0,"accepted_ref_mutation":False,
            "claim_limit":"Synthetic individual signal ablation with fixed readout and zero-fill intervention, not RSI."}

if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
