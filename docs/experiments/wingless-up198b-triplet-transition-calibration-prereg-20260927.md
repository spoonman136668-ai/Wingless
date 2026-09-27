# Wingless UP-198B — higher-order transition calibration

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-197B f5728ed889c5e6ef479952e18a2a81f01090c17a.

## Question

Does adding fixed directed adjacent-triplet history structure capture held-out calibration dependence that pairwise adjacent transitions miss?

## Frozen compositions

Reuse:
- mixed4: 0,5,1,6
- observe4: 5,6,7,8

All 24 permutations per composition.

## Frozen training and evaluation

Training phases:
- 55
- 56
- 57

Untouched evaluation phases:
- 61
- 62
- 63

## Frozen representations

Pairwise model:
- intercept + 12 directed adjacent-pair indicators, exactly the UP-197B representation.

Triplet model:
- intercept + 24 directed adjacent-triplet indicators corresponding to all ordered triples of distinct canonical ranks.

Both use ridge lambda 1e-6.
No feature selection, interaction search, phase identity, parity, or held-out fitting.

Constant baseline:
- training mean required factor per composition.

## Measurements

On untouched evaluation phases:
- actual required factor
- pairwise prediction
- triplet prediction
- constant-baseline prediction
- absolute errors

Per composition:
- mean and max absolute error for baseline, pairwise, and triplet models
- number of points where triplet beats pairwise
- number of points where triplet beats baseline

## Interpretation

A held-out triplet improvement would support higher-order local transition topology as part of Wingless's native history state. No improvement would imply that useful history information is not captured by short local n-grams alone.

## Bounds

Shadow calibration only. No held-out fitting, adaptive feature selection, architecture search, phase/parity input, correction activation, maintenance action, capacity change, live activation, or production authority.
