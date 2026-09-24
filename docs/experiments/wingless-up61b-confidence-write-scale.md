# Wingless UP-61B — confidence-carry high-noise write-scale cross

Status: preregistered scientific robustness/scale experiment.

Scientific parent: UP-60B seal `f639951748145749a75403a6cb0e59bd375aea15`.

## Question

UP-60B showed that the frozen confidence-carry family retained a commit-accuracy advantage at memory noise 0.08 through 0.09 for 64-write scenarios.

UP-61B asks whether that advantage survives when the same higher-noise regime is crossed with much longer write horizons.

## Frozen design

Confidence thresholds remain unchanged:

- 0.25;
- 0.5;
- 0.75.

They are compared against the unchanged closed-loop baseline on:

- untouched deterministic schedule bases 89M, 90M, and 91M;
- memory noise 0.08, 0.085, and 0.09;
- write horizons 128, 256, and 512.

The frozen pair-[0,5] geometry, decoder training, observer bank, propagation depths, relation head, and state-update semantics are unchanged.

No oracle state, retry, retraining, adaptive thresholding, threshold selection, result-informed stopping, production authority, or activation is permitted.

## Interpretation

This is a fixed high-noise/long-horizon scale cross. Improvement, reversal, degradation, or interaction effects are all valid scientific results and must be sealed exactly as observed.
