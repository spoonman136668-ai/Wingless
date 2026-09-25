# Wingless UP-62B — decoder-margin calibration at longer write horizons

Status: preregistered scientific calibration/scale experiment.

Scientific parent: sealed UP-61B margin-calibration Windows evidence `3ec4c5975388b109892a572e850c8d8729b2d1c8`.

## Question

UP-61B showed that the highest frozen decoder-margin bin had higher commit correctness than the lowest bin in all 12 tested 64-write conditions, while the four bins were not strictly monotonic in every individual condition. UP-62B asks whether that descriptive confidence signal persists as closed-loop error accumulation grows.

## Frozen design

Unchanged substrate and diagnostic semantics:

- frozen pair [0,5];
- existing observer bank and classifier/regressor training;
- training depth 0;
- held-out depths 32, 128, 512, 1024;
- 48 deterministic scenarios;
- baseline closed-loop update only;
- no confidence-carry intervention;
- margin bins [0.00,0.25), [0.25,0.50), [0.50,0.75), [0.75,1.00].

Untouched deterministic schedule bases:

- 92M;
- 93M.

Memory-noise levels:

- 0.08;
- 0.085;
- 0.09.

Write horizons:

- 128;
- 256.

For every schedule/noise/write condition, record event count and empirical commit correctness in each frozen bin.

No oracle state, carry policy, retry, retraining, adaptive binning, threshold selection, result-informed stopping, production authority, or activation is permitted.

## Interpretation

Persistence, flattening, inversion, sparsity, or failure of the margin/correctness relationship are all valid results. The bins remain descriptive and are not candidate policy thresholds.
