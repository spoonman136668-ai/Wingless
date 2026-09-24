# Wingless UP-58B — confidence-carry generalization

Status: preregistered sandboxed robustness/self-evaluation experiment.

Scientific parent: UP-57B seal `35b4000c1b25372419ad1f740e19146694c84952`.

## Question

UP-57B showed that the endogenous decoder margin can be used causally to reduce closed-loop error amplification on three untouched schedules.

UP-58B asks whether that effect generalizes rather than depending on those particular schedules or the 64-write horizon.

## Frozen design

No threshold is selected from UP-57B.

All three preregistered confidence-carry cutoffs are carried forward unchanged: 0.25, 0.5, and 0.75.

They are compared against the unchanged closed-loop baseline on:

- new deterministic schedule bases 70M, 71M, and 72M;
- memory noise 0.06, 0.07, and 0.075;
- write horizons 32 and 128.

The frozen pair-[0,5] geometry, decoder training, observer bank, propagation depths, relation head, and state-update semantics are unchanged.

No oracle state, retry, retraining, adaptive thresholding, or production authority is permitted.

## Interpretation

This experiment measures whether self-evaluation-based carry-forward is a robust policy family across new conditions. It does not choose or promote a threshold. Any later policy selection must use a separate untouched confirmation design.
