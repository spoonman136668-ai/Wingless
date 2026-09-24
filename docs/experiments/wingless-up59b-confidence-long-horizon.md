# Wingless UP-59B — confidence-carry long-horizon scale ladder

Status: preregistered sandboxed robustness/self-evaluation experiment.

Scientific parent: UP-58B seal `850e010027508a9c463656b9f4ce5e73542048dd`.

## Question

UP-58B showed that all three pre-existing confidence-carry thresholds improved commit accuracy over the closed-loop baseline across 18 untouched conditions while preserving relational accuracy.

UP-59B asks whether that benefit survives substantially longer write horizons without threshold selection or retraining.

## Frozen design

The confidence thresholds remain unchanged:

- 0.25;
- 0.5;
- 0.75.

They are compared against the unchanged closed-loop baseline on:

- untouched deterministic schedule bases 80M, 81M, and 82M;
- memory noise 0.07 and 0.075;
- write horizons 256 and 512.

The frozen pair-[0,5] geometry, decoder training, observer bank, propagation depths, relation head, and state-update semantics are unchanged from UP-58B.

No oracle state, retry, retraining, adaptive thresholding, threshold selection, production authority, or activation is permitted.

## Interpretation

This is a pure scale/robustness test. A degradation or reversal at long horizons is a valid scientific result and must be sealed as observed. Passing results do not select or promote any threshold; threshold choice, if ever studied, requires a separate untouched preregistered confirmation experiment.
