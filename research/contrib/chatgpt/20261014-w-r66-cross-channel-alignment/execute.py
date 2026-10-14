"""W-R66 frozen cross-channel alignment vs marginal temporal-order perturbation."""
import importlib.util
import json
from pathlib import Path

PARENT = Path(__file__).resolve().parent.parent / "20261014-w-r65-temporal-order" / "execute.py"
spec = importlib.util.spec_from_file_location("w65_pinned", PARENT)
assert spec and spec.loader
w65 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w65)

HELDOUT = tuple(range(10016, 10040))
CONDITIONS = (
    "intact", "shift_source0", "shift_source1", "shift_sources01",
    "shift_sources23", "shift_all", "reverse_source0", "reverse_all"
)
AMPS = (1, 2, 3)

def acquire(rows, order, condition):
    assert len(rows) == 16 and sorted(order) == [0, 1, 2, 3]
    assert condition in CONDITIONS
    if condition == "intact":
        channels, mode = (), "none"
    elif condition == "shift_source0":
        channels, mode = (0,), "shift"
    elif condition == "shift_source1":
        channels, mode = (1,), "shift"
    elif condition == "shift_sources01":
        channels, mode = (0, 1), "shift"
    elif condition == "shift_sources23":
        channels, mode = (2, 3), "shift"
    elif condition == "shift_all":
        channels, mode = (0, 1, 2, 3), "shift"
    elif condition == "reverse_source0":
        channels, mode = (0,), "reverse"
    else:
        channels, mode = (0, 1, 2, 3), "reverse"
    out = [list(row) for row in rows]
    for t in range(16):
        for ch in channels:
            source_t = (t + 1) % 16 if mode == "shift" else 15 - t
            out[t][ch] = rows[source_t][ch]
    return tuple(tuple(row[ch] for ch in order) for row in out)

def rank(seed, entity, site, amp, weights, order, condition):
    m = w65.w64.w63.w61.m.m
    baseline, pos, neg = m.w37.recoded_example(
        seed, entity, "xorshift_independent", "identity", site, amp, False)
    baseline, pos, neg = (
        acquire(rows, order, condition) for rows in (baseline, pos, neg))
    p, n = m.w42.feature(baseline, pos), m.w42.feature(baseline, neg)
    assert len(p) == len(n) == len(weights) == 16
    margin = sum(w * (a - b) for w, a, b in zip(weights, p, n))
    return m.w45.credit(margin, True) + m.w45.credit(-margin, False)

def run():
    trained = w65.w64.w63.w61.fit()
    records = []
    for i in (1, 2, 3):
        model = f"{i}:withhold_xorshift_independent"
        fit = trained[model]
        assert fit["training_counts"]["xorshift_independent"] == 0
        weights = fit["weights"]
        assert len(weights) == 16
        for amp in AMPS:
            for permutation, order in w65.w64.PERMS.items():
                for condition in CONDITIONS:
                    seed_rows = []
                    for seed in HELDOUT:
                        credit = sum(rank(seed, entity, site, amp, weights, order, condition)
                                     for entity in w65.w64.ENTITIES for site in w65.w64.SITES)
                        seed_rows.append({
                            "seed": seed, "correct_credit": credit,
                            "presentations": len(w65.w64.ENTITIES)*len(w65.w64.SITES)*2
                        })
                    numerator = sum(x["correct_credit"] for x in seed_rows)
                    denominator = sum(x["presentations"] for x in seed_rows)
                    records.append({
                        "model": model, "amplitude": amp,
                        "permutation": permutation, "condition": condition,
                        "credit": numerator, "presentations": denominator,
                        "accuracy": numerator / denominator, "by_seed": seed_rows
                    })
    assert len(records) == 3*3*4*len(CONDITIONS)
    scores = {(r["model"], r["amplitude"], r["permutation"], r["condition"]):
              r["accuracy"] for r in records}
    comparisons = []
    for i in (1, 2, 3):
        model = f"{i}:withhold_xorshift_independent"
        for amp in AMPS:
            for perm in w65.w64.PERMS:
                k = (model, amp, perm)
                intact = scores[k + ("intact",)]
                changes = {name: scores[k + (name,)] - intact
                           for name in CONDITIONS if name != "intact"}
                comparisons.append({
                    "model": model, "amplitude": amp, "permutation": perm,
                    "intact_accuracy": intact, "delta": changes
                })
    assert len(comparisons) == 36
    mean = {name: sum(c["delta"][name] for c in comparisons)/36
            for name in CONDITIONS if name != "intact"}
    individuals = (
        mean["shift_source0"] <= -0.06 and mean["shift_source1"] <= -0.06)
    rescue = (
        individuals and mean["shift_all"] >= -0.04
        and mean["shift_all"] >= max(mean["shift_source0"], mean["shift_source1"])+0.04)
    absent = individuals and mean["shift_all"] <= -0.06
    classification = (
        "SYNCHRONOUS_ALIGNMENT_RESCUE" if rescue else
        "NO_GLOBAL_ALIGNMENT_RESCUE" if absent else
        "MIXED_OR_INCONCLUSIVE")
    return {
        "schema": "wingless.w-r66-frozen-alignment-ablation.v1",
        "parent_head": "be4517b888d0315e1da392eb066a750e3ea96ab2",
        "classification": classification,
        "sealed_heldout_seeds": HELDOUT,
        "training_blocks": w65.w64.w63.w61.DEV,
        "generator": "xorshift_independent",
        "conditions": CONDITIONS, "mean_accuracy_change": mean,
        "comparisons": comparisons, "all_sealed_records": records,
        "model_calls": 0, "accepted_ref_mutation": False,
        "claim_limit": "Synthetic fixed-readout shift/reversal; global time-origin confounding remains possible; not RSI."
    }

if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
