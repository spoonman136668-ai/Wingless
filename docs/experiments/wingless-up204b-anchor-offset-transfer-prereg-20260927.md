# Wingless UP-204B — single-anchor offset transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-203B 8a379456d273d79b9410c47e1009b3ca05a2cafa.

## Question

Can one measured reference history per phase remove the phase-wide offset and recover the frozen relative order template for the remaining histories?

## Frozen compositions

- mixed4: 0,5,1,6
- observe4: 5,6,7,8

All 24 permutations per composition.

## Frozen template phases

- 55
- 56
- 57

Untouched evaluation phases:
- 79
- 80
- 81

## Frozen anchor

The canonical permutation for each composition is the single anchor.

## Frozen predictor

From template phases:
- compute each permutation's mean required factor;
- subtract the composition-wide template mean to form a residual template.

For each untouched phase:
- measure the anchor permutation's required factor;
- phase offset = measured anchor factor - frozen anchor residual;
- predict each non-anchor permutation as phase offset + its frozen residual.

The anchor point is excluded from error scoring.

Comparator:
- exact-permutation training mean with no phase offset correction.

No held-out fitting beyond the single preregistered anchor measurement.

## Measurements

Per composition:
- comparator mean/max absolute error over 23 non-anchor permutations × 3 phases;
- anchor-offset mean/max absolute error;
- count where anchor-offset beats comparator.

## Interpretation

Improvement would show that a low-dimensional phase-wide offset plus stable relative history geometry is sufficient to recover calibration.

## Bounds

Shadow calibration only. No anchor search, held-out feature fitting, nonlinear correction, maintenance action, live activation, or production authority.
