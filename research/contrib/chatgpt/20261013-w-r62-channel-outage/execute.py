"""W-R62 frozen source-heldout ridge readout under single-channel acquisition outage."""
import importlib.util
import json
from pathlib import Path

PARENT = Path(__file__).resolve().parent.parent / "20261013-w-r61-entity-transfer" / "execute.py"
spec = importlib.util.spec_from_file_location("wingless_w61_pinned", PARENT)
assert spec and spec.loader
w61 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w61)

HELDOUT = tuple(range(9312, 9336))
AMPS = (1, 2, 3)
GENERATORS = w61.GENERATORS
GROUPS = w61.GROUPS
ARMS = ("intact_reference", "channel0_zeroed", "channel1_zeroed", "channel2_zeroed", "channel3_zeroed")
SITES = w61.m.TEST_SITES

def acquire(rows, arm):
    if arm == "intact_reference":
        return rows
    channel = ARMS.index(arm) - 1
    assert 0 <= channel < 4
    assert len(rows) == 16 and all(len(row) == 4 for row in rows)
    return tuple(tuple(0 if j == channel else value for j, value in enumerate(row))
                 for row in rows)

def case(seed, entity, generator, site, amplitude, model, arm):
    base, positive, negative = w61.m.w37.recoded_example(
        seed, entity, generator, "identity", site, amplitude, False)
    base = acquire(base, arm)
    positive = acquire(positive, arm)
    negative = acquire(negative, arm)
    a = w61.m.w42.feature(base, positive)
    b = w61.m.w42.feature(base, negative)
    weights = model["weights"]
    assert len(a) == len(b) == len(weights) == 16
    margin = sum(w * (p - n) for w, p, n in zip(weights, a, b))
    return w61.m.w45.credit(margin, True) + w61.m.w45.credit(-margin, False)

def run():
    trained = w61.fit()
    assert len(trained) == 9
    records = []
    summary = []
    for i in (1, 2, 3):
        model_name = f"{i}:withhold_xorshift_independent"
        model = trained[model_name]
        assert model["training_counts"]["xorshift_independent"] == 0
        for amp in AMPS:
            for generator in GENERATORS:
                for group, entities in GROUPS.items():
                    for arm in ARMS:
                        per_seed = []
                        all_credit = 0.0
                        all_count = 0
                        for seed in HELDOUT:
                            credit = 0.0
                            presentations = 0
                            for site in SITES:
                                for entity in entities:
                                    credit += case(seed, entity, generator, site, amp, model, arm)
                                    presentations += 2
                            assert presentations == len(SITES) * len(entities) * 2
                            per_seed.append({"seed": seed, "credit": credit,
                                             "presentations": presentations})
                            all_credit += credit
                            all_count += presentations
                        accuracy = all_credit / all_count
                        entry = {"model": model_name, "amplitude": amp,
                                 "generator": generator, "entity_group": group,
                                 "arm": arm, "accuracy": accuracy, "credit": all_credit,
                                 "presentations": all_count, "by_seed": per_seed}
                        records.append(entry)
                        summary.append({k: v for k, v in entry.items() if k != "by_seed"})
    lookup = {(x["model"], x["amplitude"], x["generator"], x["entity_group"], x["arm"]): x["accuracy"]
              for x in records}
    diffs = []
    for i in (1, 2, 3):
        model = f"{i}:withhold_xorshift_independent"
        for amp in AMPS:
            base = lookup[(model, amp, "xorshift_independent", "unseen_identity_transfer", "intact_reference")]
            for arm in ARMS[1:]:
                accuracy = lookup[(model, amp, "xorshift_independent", "unseen_identity_transfer", arm)]
                diffs.append({"model": model, "amplitude": amp, "arm": arm,
                              "intact_accuracy": base, "outage_accuracy": accuracy,
                              "outage_minus_intact": accuracy - base})
    assert len(diffs) == 36
    robust = all(x["outage_accuracy"] >= .55 and x["outage_minus_intact"] >= -.05 for x in diffs)
    partial = (sum(x["outage_accuracy"] >= .55 for x in diffs) >= 18 and
               sum(x["outage_minus_intact"] for x in diffs) / len(diffs) >= -.10)
    classification = "CHANNEL_OUTAGE_ROBUST" if robust else (
        "PARTIAL_CHANNEL_OUTAGE" if partial else "OUTAGE_TRANSFER_FAILURE")
    return {"schema": "wingless.w-r62-channel-outage-firstuse.v1",
            "parent_head": "4fd74701a7a4cb61a7e4156dd9651390c7d981d5",
            "classification": classification,
            "training_blocks": w61.DEV, "sealed_heldout_seeds": HELDOUT,
            "entity_groups": GROUPS, "source_generators": GENERATORS,
            "amplitudes": AMPS, "outage_arms": ARMS,
            "primary_xorshift_unseen_differences": diffs,
            "compact_results": summary, "all_sealed_trial_records": records,
            "training_provenance": {f"{i}:withhold_xorshift_independent":
                                     trained[f"{i}:withhold_xorshift_independent"]["development_sha256"]
                                     for i in (1, 2, 3)},
            "model_calls": 0, "accepted_ref_mutation": False,
            "claim_limit": "Synthetic unseen-source fixed-readout channel outage, not generalized reasoning."}

if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
