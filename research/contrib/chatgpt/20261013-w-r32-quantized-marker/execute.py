"""W-R32: untouched clean observations quantized to test natural marker alias incidence."""
import importlib.util
import json
from pathlib import Path

SOURCE = Path(__file__).resolve().parent.parent / "20261013-w-r31-legitimate-markers" / "execute.py"
spec = importlib.util.spec_from_file_location("w31_pinned", SOURCE)
assert spec and spec.loader
w31 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w31)

HELDOUT = tuple(range(1080, 1096))
GENERATORS = w31.GENERATORS
BITS = (8, 4, 2, 1)
SITES = (10, 12)
ARMS = w31.ARMS


def quantize(row, bits):
    levels = (1 << bits) - 1
    return tuple(((value >> (8 - bits)) * 255) // levels for value in row)


def observed_pair(seed, entity, generator, bits, site):
    base = w31.w25.w24.w22.w20.base_sequence(seed, entity, generator)
    raw = w31.w25.w16.observed(base, seed, entity, 3, 79)
    positive = tuple(quantize(row, bits) for row in raw)
    tail = list(positive[1:])
    tail[site], tail[site+1] = tail[site+1], tail[site]
    negative = positive[:1] + tuple(tail)
    assert len(positive) == len(negative) == 16
    assert positive[0] == negative[0]
    assert sorted(positive) == sorted(negative)
    return base, positive, negative


def cell(arm, model, generator, bits, site):
    by_seed = []
    correct = total = alias_count = alias_slots = identical_pairs = 0
    for seed in HELDOUT:
        tp = tn = fp = fn = aliases = ident = 0
        for entity in range(9):
            base, positive, negative = observed_pair(seed, entity, generator, bits, site)
            aliases += sum(w31.w25.marker(positive[t], positive[t-1])
                           for t in range(1, len(positive)))
            ident += int(positive == negative)
            for truth, candidate in ((True, positive), (False, negative)):
                prediction = w31.classify(model, base, candidate, arm)
                if truth and prediction: tp += 1
                elif not truth and not prediction: tn += 1
                elif not truth: fp += 1
                else: fn += 1
        assert tp + tn + fp + fn == 18
        total += 18
        correct += tp + tn
        alias_count += aliases
        alias_slots += 9*15
        identical_pairs += ident
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "genuine_marker_aliases": aliases,
                        "legitimate_marker_opportunities":9*15,
                        "observation_identical_pairs":ident})
    return {"generator":generator,"bits":bits,"tail_swap_index":site,
            "correct":correct,"total":total,"accuracy":correct/total,
            "genuine_marker_alias_count":alias_count,
            "genuine_marker_slot_count":alias_slots,
            "genuine_marker_incidence":alias_count/alias_slots,
            "identical_pairs":identical_pairs,"by_seed":by_seed}


def run():
    model = w31.w25.w22.train("maskaware_augmented")
    assert len(model[0]) == 16
    results = []
    for arm in ARMS:
        cells = [cell(arm, model, generator, bits, site)
                 for generator in GENERATORS for bits in BITS for site in SITES]
        assert len(cells) == 24
        results.append({"arm":arm,"width":16,"training_examples":1728,
                        "epochs":32,"training_updates":model[2],
                        "aggregate_accuracy":sum(c["correct"] for c in cells)/
                                             sum(c["total"] for c in cells),
                        "cells":cells})
    by_arm = {r["arm"]:r for r in results}
    # Marker prevalence is an observation property and must be identical
    # for all arms on exactly the same frozen cohort.
    incidence={}
    accuracy={}
    for bits in BITS:
        reference = [c for c in by_arm["multi_marker"]["cells"] if c["bits"]==bits]
        incidence[str(bits)] = (sum(c["genuine_marker_alias_count"] for c in reference)/
                                sum(c["genuine_marker_slot_count"] for c in reference))
        accuracy[str(bits)]={a["arm"]:sum(c["correct"] for c in a["cells"] if c["bits"]==bits)/
                                      sum(c["total"] for c in a["cells"] if c["bits"]==bits)
                             for a in results}
        for r in results:
            comparable=[c for c in r["cells"] if c["bits"]==bits]
            assert [c["genuine_marker_alias_count"] for c in comparable] == [
                   c["genuine_marker_alias_count"] for c in reference]
    rise=incidence["1"]-incidence["8"]
    penalty=accuracy["1"]["all_observed"]-accuracy["1"]["multi_marker"]
    if rise>=0.05 and penalty>=0.05:
        classification="QUANTIZED_MARKER_COST"
    elif rise>=0.05:
        classification="ALIAS_WITHOUT_CONFIRMED_COST"
    else:
        classification="NO_ALIAS_GROWTH"
    return {"schema":"wingless.w-r32-natural-quantized-markers.v1",
            "parent_head":"0c8679500be0322dcaeb17637eb56867744e73b6",
            "classification":classification,"sealed_heldout_seeds":HELDOUT,
            "generators":GENERATORS,"quantizer_bits":BITS,
            "genuine_marker_incidence_by_bits":incidence,
            "downstream_accuracy_by_bits":accuracy,
            "one_bit_alias_rise":rise,"one_bit_all_observed_minus_marker":penalty,
            "arms":results,"model_calls":0,
            "claim_limit":"Artificially quantized but genuinely observed deterministic synthetic sequences, not field sensor incidence or semantic learning"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
