"""W-R30: first-use latent erasure vs legitimate marker indistinguishability."""
import importlib.util
import json
from pathlib import Path

SOURCE = Path(__file__).resolve().parent.parent / "20261013-w-r25-marker-invariance" / "execute.py"
spec = importlib.util.spec_from_file_location("w25_frozen", SOURCE)
assert spec and spec.loader
w25 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w25)

SEEDS = tuple(range(1000, 1016))
GENERATORS = ("affine_prior", "quadratic_irregular", "xorshift_independent")
ENCODINGS = ("all_zero", "all_255", "carry_forward")
TICKS = (5, 8)
DETECTORS = ("always_observed", "old_zero_detector", "observed_multi_marker_detector")


def observable_pair(seed, entity, generator, encoding, tick):
    base = w25.w24.w22.w20.base_sequence(seed, entity, generator)
    raw_observed = w25.w16.observed(base, seed, entity, 3, 79)
    changed = list(raw_observed)
    marker = ((0, 0, 0, 0) if encoding == "all_zero" else
              (255, 255, 255, 255) if encoding == "all_255" else changed[tick - 1])
    assert tick in (5, 8)
    assert raw_observed[tick] != marker  # no trivial naturally-marker original
    changed[tick] = marker
    acquired = tuple(changed)
    # Counterfactual latent cases:
    # A: measurement lost and the encoder filled in marker.
    # B: same bytes are a legitimate acquired value; no measurement lost.
    # Every model-observable input is IDENTICAL across worlds.
    erasure_world = (base, acquired, True)
    legitimate_world = (base, acquired, False)
    assert erasure_world[:2] == legitimate_world[:2]
    assert erasure_world[2] != legitimate_world[2]
    return erasure_world, legitimate_world


def predict_missing(row, previous, detector):
    if detector == "always_observed":
        return False
    if detector == "old_zero_detector":
        return row == (0, 0, 0, 0)
    assert detector == "observed_multi_marker_detector"
    return w25.marker(row, previous)


def cell(generator, encoding, tick):
    seed_rows = []
    for seed in SEEDS:
        counts = {name: {"true_erasure_detected": 0,
                         "legitimate_falsely_masked": 0,
                         "balanced_correct": 0}
                  for name in DETECTORS}
        matched = 0
        for entity in range(9):
            a, b = observable_pair(seed, entity, generator, encoding, tick)
            assert a[:2] == b[:2] and a[2] is True and b[2] is False
            for arm in DETECTORS:
                pa = predict_missing(a[1][tick], a[1][tick-1], arm)
                pb = predict_missing(b[1][tick], b[1][tick-1], arm)
                assert pa == pb
                counts[arm]["true_erasure_detected"] += int(pa)
                counts[arm]["legitimate_falsely_masked"] += int(pb)
                counts[arm]["balanced_correct"] += int(pa) + int(not pb)
            matched += 1
        assert matched == 9
        for arm in DETECTORS:
            assert counts[arm]["balanced_correct"] == 9  # max accuracy .50
        seed_rows.append({"seed": seed, "matched_observable_pairs": matched,
                          "arms": counts})
    assert len(seed_rows) == 16
    aggregated = {}
    for arm in DETECTORS:
        tp = sum(s["arms"][arm]["true_erasure_detected"] for s in seed_rows)
        fp = sum(s["arms"][arm]["legitimate_falsely_masked"] for s in seed_rows)
        correct = sum(s["arms"][arm]["balanced_correct"] for s in seed_rows)
        assert correct == 144
        aggregated[arm] = {"erasure_recall": tp/144,
                           "legitimate_false_mask_rate": fp/144,
                           "balanced_latent_mask_accuracy": correct/288,
                           "detected_erasures": tp,
                           "false_masked_valid_rows": fp}
    return {"generator": generator, "encoding": encoding, "tick": tick,
            "observationally_identical_pairs": 144,
            "arms": aggregated,
            "per_seed": seed_rows}


def run():
    cells = [cell(g, f, tick) for g in GENERATORS
             for f in ENCODINGS for tick in TICKS]
    assert len(cells) == 18
    pairs = sum(c["observationally_identical_pairs"] for c in cells)
    assert pairs == 2592
    assert all(x["balanced_latent_mask_accuracy"] == .5
               for c in cells for x in c["arms"].values())
    return {
        "schema": "wingless.w-r30-marker-alias-boundary.v1",
        "parent_head": "5617fd5f5a7ac4ec07b084ac08ede17e5c66f96d",
        "classification": "UNIDENTIFIABLE_WITHOUT_PROVENANCE",
        "sealed_heldout_seeds": SEEDS,
        "distinct_observed_samples": pairs,
        "paired_latent_worlds": 2*pairs,
        "cells": cells, "model_calls": 0,
        "claim_limit": "Constructed indistinguishable latent marker cases; not an estimate of natural collision prevalence or downstream model accuracy",
    }


if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
