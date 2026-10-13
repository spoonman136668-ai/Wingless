"""W-R27: synthetic counterfactual observability test of carry-forward markers."""
import hashlib
import importlib.util
import json
from pathlib import Path

SOURCE = Path(__file__).resolve().parent.parent / "20261013-w-r25-marker-invariance" / "execute.py"
spec = importlib.util.spec_from_file_location("w25_pinned", SOURCE)
assert spec and spec.loader
w25 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w25)

SEEDS = tuple(range(400, 416))
GENERATORS = w25.GENERATORS
POSITIONS = (3, 6, 9)
ZERO = (0, 0, 0, 0)
FULL = (255, 255, 255, 255)


def run():
    by_generator = []
    digester = hashlib.sha256()
    total_pairs = excluded = marked_pairs = zero_only_pairs = 0
    for generator in GENERATORS:
        cells = []
        for seed in SEEDS:
            n_pairs = skipped = 0
            for entity in range(9):
                base = w25.w24.w20.base_sequence(seed, entity, generator)
                observed = w25.w16.observed(base, seed, entity, 3, 79)
                assert len(observed) == 16
                for pos in POSITIONS:
                    previous, actual = observed[pos - 1], observed[pos]
                    # Exclusions are frozen structural degeneracies, never outcomes.
                    if previous == actual or previous in (ZERO, FULL):
                        skipped += 1
                        continue
                    # World C: the clean latent sensor value naturally repeats.
                    # World M: a different true value is erased and carried forward.
                    clean_latent = tuple(observed[:pos]) + (previous,) + tuple(observed[pos + 1:])
                    missing_latent = tuple(observed)
                    clean_visible = clean_latent
                    missing_visible = tuple(missing_latent[:pos]) + (previous,) + tuple(missing_latent[pos + 1:])
                    assert clean_visible == missing_visible
                    assert clean_latent[pos] == previous
                    assert missing_latent[pos] != previous
                    marked = w25.marker(clean_visible[pos], clean_visible[pos - 1])
                    assert marked == w25.marker(missing_visible[pos], missing_visible[pos - 1])
                    marked_pairs += int(marked)
                    zero_only_pairs += int(clean_visible[pos] == ZERO)
                    n_pairs += 1
                    # Bind paired fixture to its exact visible bytes and latent labels.
                    digester.update(json.dumps([generator, seed, entity, pos,
                                clean_visible, 0, missing_visible, 1],
                                separators=(",", ":")).encode())
            total_pairs += n_pairs
            excluded += skipped
            cells.append({"seed": seed, "balanced_opposite_label_pairs": n_pairs,
                          "structural_exclusions": skipped})
        by_generator.append({"generator": generator, "cells": cells})
    assert total_pairs > 0
    # In every pair both algorithms receive exactly the same observation, so
    # predicting whether a carry-forward event occurred is chance-limited.
    observed_marker_balanced_accuracy = 0.5
    classification = ("DETECTOR_IDENTIFIABILITY_LIMIT"
                      if marked_pairs == total_pairs and zero_only_pairs == 0
                      else "PARTIAL_COUNTEREXAMPLE")
    return {"schema": "wingless.w-r27-carry-forward-identifiability.v1",
            "parent_head": "5617fd5f5a7ac4ec07b084ac08ede17e5c66f96d",
            "classification": classification,
            "development_or_training_updates": 0,
            "sealed_heldout_seeds": SEEDS,
            "original_generator_families": GENERATORS,
            "sampled_positions": POSITIONS,
            "balanced_opposite_label_pairs": total_pairs,
            "structural_exclusions": excluded,
            "marked_pairs": marked_pairs,
            "allzero_marker_pairs": zero_only_pairs,
            "w25_marker_heldout_balanced_accuracy": observed_marker_balanced_accuracy,
            "fixture_sha256": digester.hexdigest(),
            "by_generator": by_generator,
            "model_calls": 0,
            "claim_limit": "Synthetic event identifiability for observed-only carry-forward marker; not swap-label model accuracy or language cognition"}


if __name__ == "__main__":
    print(json.dumps(run(), sort_keys=True, indent=2))
