# Wingless UP-197B — directed-transition calibration

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-196B 6a4ed7ef52e69cf835d143402d670e13800793c6.

## Question

Can a fixed representation of *which directed local transitions occurred* predict the history-dependent calibration factor on untouched future phases better than a history-agnostic constant baseline?

## Frozen compositions

Reuse the two full-permutation composition families from UP-196B:
- mixed4: 0,5,1,6
- observe4: 5,6,7,8

All 24 permutations are used.

## Frozen feature representation

Map each item to its canonical rank 0..3.

For a permutation, emit 12 binary features indicating the presence of each possible directed adjacent transition a>b where a != b.

No feature selection, interaction search, permutation distance, phase identity, or parity is included.

## Frozen training and evaluation

Training phases:
- 55
- 56
- 57

These are the already-observed UP-196B phases.

Untouched evaluation phases:
- 58
- 59
- 60

Fit one ridge model per composition:
- intercept + 12 directed-transition features
- lambda = 1e-6

Baseline predictor per composition is the training mean required factor.

## Measurements

On untouched evaluation phases:
- actual required factor
- transition-model predicted factor
- baseline predicted factor
- absolute errors

Per composition:
- baseline mean and max absolute error
- transition-model mean and max absolute error
- fraction of evaluation points where the transition model improves on baseline

## Interpretation

A material held-out improvement supports directed local transition topology as a useful history coordinate. Little or no improvement means the relevant history representation remains higher-order than pairwise adjacency.

## Bounds

Shadow calibration only. No held-out fitting, adaptive feature selection, phase/parity input, correction activation, maintenance action, capacity change, live activation, or production authority.
