"""W-R20: frozen out-of-generator transfer, no retraining on heldout generators."""
import importlib.util
import json
from pathlib import Path

PREVIOUS = Path(__file__).resolve().parent.parent / "20261013-w-r19-multichannel" / "execute.py"
spec = importlib.util.spec_from_file_location("w19_frozen", PREVIOUS)
assert spec and spec.loader
w19 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w19)

HELDOUT = tuple(range(192, 208))
GENERATORS = ("affine_prior", "quadratic_irregular", "xorshift_independent")
ARMS = ("first_channel_max", "mean_four_max", "median_four_max", "mean_top_two_64")
FAMILIES = ("swap10_11", "swap12_13")
MASK = (1 << 32) - 1


def base_sequence(seed, entity, generator):
    if generator == "affine_prior":
        return w19.w16.original(seed, entity)
    if generator == "quadratic_irregular":
        return tuple(
            tuple((43*seed + 29*entity + 17*c + (5+2*c)*t
                   + (11+3*c)*t*t + 7*(((seed+c)*t) % 19)) % 256
                  for c in range(4))
            for t in range(16)
        )
    assert generator == "xorshift_independent"
    seq = []
    for t in range(16):
        channels = []
        for c in range(4):
            state = ((seed*0x9e3779b1) ^ (entity*0x85ebca77)
                     ^ (c*0xc2b2ae3d) ^ ((t+1)*0x27d4eb2f)) & MASK
            state ^= (state << 13) & MASK
            state &= MASK
            state ^= state >> 17
            state &= MASK
            state ^= (state << 5) & MASK
            state &= MASK
            channels.append(state & 255)
        seq.append(tuple(channels))
    return tuple(seq)


def pair(seed, entity, amp, generator, negative_transform, positive):
    original = base_sequence(seed, entity, generator)
    candidate = w19.w16.observed(original, seed, entity, amp, 79)
    compared = candidate if positive else w19.permute(candidate, negative_transform)
    assert original and len(original) == 16
    assert candidate[0] == compared[0]
    assert sorted(candidate) == sorted(compared)
    assert (candidate == compared) == positive
    return original, compared


def evaluate(kind, trained, generator, amp, transform):
    weights, bias, updates = trained
    rows = []
    for seed in HELDOUT:
        tp = tn = fp = fn = 0
        for entity in range(9):
            a_true, b_true = pair(seed, entity, amp, generator, transform, True)
            a_false, b_false = pair(seed, entity, amp, generator, transform, False)
            assert a_true == a_false and b_true[0] == b_false[0]
            assert sorted(b_true) == sorted(b_false)
            for positive, a, b in ((True, a_true, b_true), (False, a_false, b_false)):
                vector = w19.feature(a, b, kind)
                assert len(vector) == 16
                prediction = bias + sum(w*v for w, v in zip(weights, vector)) >= 0
                if positive and prediction: tp += 1
                elif not positive and not prediction: tn += 1
                elif not positive: fp += 1
                else: fn += 1
        assert tp+tn+fp+fn == 18
        rows.append({"seed": seed, "tp": tp, "tn": tn, "fp": fp, "fn": fn})
    correct = sum(x["tp"]+x["tn"] for x in rows)
    return {"generator": generator, "amplitude": amp, "swap": transform,
            "correct": correct, "total": 16*18, "accuracy": correct/(16*18),
            "heldout_by_seed": rows}


def run():
    arms = []
    for kind in ARMS:
        trained = w19.train("single_adjacent", kind)
        cells = [evaluate(kind, trained, generator, amp, swap)
                 for generator in GENERATORS
                 for amp in (1,2,3)
                 for swap in FAMILIES]
        assert len(cells) == 18
        nonaffine = [c for c in cells if c["generator"] != "affine_prior"]
        arms.append({
            "model": kind, "training_examples": 1728,
            "width": len(trained[0]), "epochs": 32, "training_updates": trained[2],
            "nonaffine_accuracy": sum(c["correct"] for c in nonaffine)/
                                  sum(c["total"] for c in nonaffine),
            "aggregate_accuracy": sum(c["correct"] for c in cells)/
                                  sum(c["total"] for c in cells),
            "cells": cells,
        })
    arm_map = {a["model"]: a for a in arms}
    top = arm_map["mean_top_two_64"]
    baseline = arm_map["first_channel_max"]
    delta = top["nonaffine_accuracy"] - baseline["nonaffine_accuracy"]
    nonaffine_cells = [x for x in top["cells"] if x["generator"] != "affine_prior"]
    if all(c["accuracy"] >= .80 for c in nonaffine_cells) and delta >= .05:
        label = "CROSS_GENERATOR_TRANSFER"
    elif delta >= .03:
        label = "PARTIAL"
    else:
        label = "NEGATIVE"
    return {
        "schema": "wingless.w-r20-out-of-generator.v1",
        "parent_head": "42909d78f08e2d64007969f99224ad2d6c4915f8",
        "classification": label, "nonaffine_top_two_vs_first_delta": delta,
        "training_development_seeds": list(range(160,176)),
        "sealed_heldout_seeds": HELDOUT, "generator_families": GENERATORS,
        "arms": arms, "model_calls": 0,
        "claim_limit": "Synthetic deterministic generator transfer only, no generalized reasoning"
    }


if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
