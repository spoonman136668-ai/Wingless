# Wingless UP-60B — confidence-carry higher-noise extension

Status: preregistered sandboxed robustness/self-evaluation experiment.

Scientific parent: UP-58B seal `850e010027508a9c463656b9f4ce5e73542048dd`.

This experiment is defined independently of UP-59B outcomes.

## Question

Does the confidence-carry policy family retain its observed advantage as memory noise rises beyond the UP-58B range?

## Frozen design

The confidence thresholds remain unchanged:

- 0.25;
- 0.5;
- 0.75.

They are compared against the unchanged closed-loop baseline on:

- untouched deterministic schedule bases 83M, 84M, and 85M;
- memory noise 0.08, 0.085, and 0.09;
- 64 writes per scenario.

The frozen pair-[0,5] geometry, decoder training, observer bank, propagation depths, relation head, and state-update semantics are unchanged.

No oracle state, retry, retraining, adaptive thresholding, threshold selection, production authority, or activation is permitted.

## Interpretation

This is a higher-noise robustness probe, not threshold selection. Any reversal or degradation is a valid scientific result and must be sealed unchanged.
