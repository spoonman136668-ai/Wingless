# Wingless UP-205B — anchor-offset family transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-204B d64d7b33463536887bf759a15215367eabd26800.

## Question

Does the single-anchor phase-offset + frozen relative-history template model generalize beyond mixed4 and observe4 to distinct store4 and cross3 history families?

## Frozen compositions

- store4: 0,1,2,3
- cross3: 0,5,13

All permutations per composition:
- store4: 24
- cross3: 6

## Frozen template phases

- 55
- 56
- 57

Untouched evaluation phases:
- 82
- 83
- 84

## Frozen anchor and predictor

The canonical permutation is the single anchor for each composition.

From template phases:
- compute each exact permutation's mean required factor;
- subtract the composition-wide template mean to form a frozen residual template.

For each untouched phase:
- measure only the canonical anchor's required factor;
- phase offset = measured anchor factor - frozen anchor residual;
- predict every non-anchor permutation as phase offset + its frozen residual.

Comparator:
- exact-permutation training mean with no phase-offset correction.

The anchor is excluded from error scoring.

No held-out anchor search or feature fitting.

## Measurements

Per composition:
- comparator mean/max absolute error;
- anchor-offset mean/max absolute error;
- count where anchor-offset beats comparator.

## Interpretation

Improvement in both families would support a reusable state decomposition: stable relative history geometry plus a low-dimensional phase offset.

## Bounds

Shadow calibration only. No anchor search, held-out fitting, nonlinear correction, maintenance action, live activation, or production authority.
