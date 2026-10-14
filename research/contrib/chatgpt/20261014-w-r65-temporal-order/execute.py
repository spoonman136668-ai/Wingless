"""W-R65: prospective original-source temporal-order causal perturbation."""
import importlib.util
import json
from pathlib import Path

PARENT = Path(__file__).resolve().parent.parent / "20261014-w-r64-source-specificity" / "execute.py"
spec = importlib.util.spec_from_file_location("w64_frozen", PARENT)
assert spec and spec.loader
w64 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w64)

HELDOUT = tuple(range(9792, 9816))
CONDITIONS = ("intact", "erase_source0", "reverse_source0", "reverse_source1",
              "shift_source0", "shift_source1")
AMPS = (1, 2, 3)

def acquire(rows, order, condition):
    assert len(rows) == 16 and sorted(order) == [0, 1, 2, 3]
    assert condition in CONDITIONS
    altered = [list(row) for row in rows]
    if condition != "intact":
        source = int(condition[-1])
        if condition.startswith("erase_"):
            values = [0] * 16
        elif condition.startswith("reverse_"):
            values = [rows[15 - t][source] for t in range(16)]
        elif condition.startswith("shift_"):
            values = [rows[(t + 1) % 16][source] for t in range(16)]
        else:
            raise AssertionError(condition)
        for t, v in enumerate(values):
            altered[t][source] = v
    return tuple(tuple(row[channel] for channel in order) for row in altered)

def rank(seed, entity, site, amp, weights, order, condition):
    m = w64.w63.w61.m.m
    base, pos, neg = m.w37.recoded_example(
        seed, entity, "xorshift_independent", "identity", site, amp, False)
    base, pos, neg = (acquire(a, order, condition) for a in (base, pos, neg))
    p, n = m.w42.feature(base, pos), m.w42.feature(base, neg)
    assert len(p) == len(n) == len(weights) == 16
    margin = sum(w * (a - b) for w, a, b in zip(weights, p, n))
    return m.w45.credit(margin, True) + m.w45.credit(-margin, False)

def run():
    trained = w64.w63.w61.fit()
    records = []
    for i in (1, 2, 3):
        model = f"{i}:withhold_xorshift_independent"
        fit = trained[model]
        assert fit["training_counts"]["xorshift_independent"] == 0
        weights = fit["weights"]
        assert len(weights) == 16
        for amp in AMPS:
            for permutation, order in w64.PERMS.items():
                for condition in CONDITIONS:
                    by_seed = []
                    for seed in HELDOUT:
                        score = sum(rank(seed, entity, site, amp, weights, order, condition)
                                    for entity in w64.ENTITIES for site in w64.SITES)
                        by_seed.append({"seed": seed, "correct_credit": score,
                                        "presentations": len(w64.ENTITIES)*len(w64.SITES)*2})
                    credit = sum(x["correct_credit"] for x in by_seed)
                    presentations = sum(x["presentations"] for x in by_seed)
                    records.append({"model": model, "amplitude": amp, "permutation": permutation,
                                    "condition": condition, "credit": credit,
                                    "presentations": presentations, "accuracy": credit / presentations,
                                    "by_seed": by_seed})
    assert len(records) == 3*3*4*6
    scores = {(r["model"], r["amplitude"], r["permutation"], r["condition"]): r["accuracy"]
              for r in records}
    comparisons = []
    for i in (1, 2, 3):
        model = f"{i}:withhold_xorshift_independent"
        for amp in AMPS:
            for perm in w64.PERMS:
                k = (model, amp, perm)
                intact = scores[k + ("intact",)]
                delta = {name: scores[k + (name,)] - intact
                         for name in CONDITIONS if name != "intact"}
                comparisons.append({"model": model, "amplitude": amp,
                                    "permutation": perm, "intact_accuracy": intact,
                                    "delta_from_intact": delta})
    assert len(comparisons) == 36
    means = {name: sum(x["delta_from_intact"][name] for x in comparisons)/36
             for name in CONDITIONS if name != "intact"}
    temporal = (means["erase_source0"] <= -.06
                and means["reverse_source0"] <= -.025
                and means["shift_source0"] <= -.025
                and means["reverse_source0"] <= means["reverse_source1"] - .015
                and means["shift_source0"] <= means["shift_source1"] - .015)
    static = (means["erase_source0"] <= -.06
              and means["reverse_source0"] > -.025
              and means["shift_source0"] > -.025)
    classification = ("TEMPORAL_ORDER_DEPENDENCE" if temporal
                      else "NONORDER_SOURCE0_DEPENDENCE" if static
                      else "MIXED_OR_NONREPLICATION")
    return {"schema": "wingless.w-r65-temporal-order-causal-perturbation.v1",
            "parent_head": "d5fe061ddedab79893b2eb3c5ba5020c9ee2d8a2",
            "classification": classification, "sealed_heldout_seeds": HELDOUT,
            "training_blocks": w64.w63.w61.DEV, "entities": w64.ENTITIES,
            "generator": "xorshift_independent", "conditions": CONDITIONS,
            "mean_accuracy_change": means, "comparisons": comparisons,
            "all_sealed_records": records, "model_calls": 0,
            "accepted_ref_mutation": False,
            "claim_limit": "Paired synthetic temporal perturbation of fixed source-withheld ridge; not RSI."}

if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
