# Wingless UP-207B — native-state phase-offset proxy

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-206B dc08d20b1bde815dbb69921705519521dbca58bb.

## Question

Can the phase-wide calibration shift be predicted from Wingless-native state geometry without measuring an anchor history at evaluation time?

## Frozen composition families

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

## Frozen residual template

Template phases:
- 55
- 56
- 57

For every family and exact permutation:
- exact-permutation template mean required factor is frozen.
- canonical template mean is frozen.

The target phase shift for a family is:
- canonical required factor at that phase minus canonical template mean.

## Frozen native predictors

Computed only from the canonical history at that phase:
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

No phase number, parity, history identity feature, or held-out anchor measurement.

## Frozen fitting

Training phases:
- 58 through 72 inclusive

One pooled ridge model across all four families.
Ridge lambda:
- 1e-6

Untouched evaluation phases:
- 88
- 89
- 90

## Frozen comparators

1. static exact-permutation template mean
2. oracle one-anchor offset using the actual canonical required factor at evaluation time

Native-offset prediction:
- predicted phase shift from the ridge model;
- predicted permutation factor = frozen exact-permutation template mean + predicted shift.

The canonical point is excluded from permutation scoring.

## Measurements

Per family:
- static mean/max absolute error
- native-offset mean/max absolute error
- oracle-anchor mean/max absolute error
- native better than static count
- native versus oracle gap

## Interpretation

If native-offset prediction materially beats static calibration and approaches the oracle anchor, a low-dimensional native-state signal explains part of the common phase shift.

## Bounds

Shadow calibration only. No held-out fitting, adaptive feature selection, phase/parity input, nonlinear model, correction activation, maintenance action, live activation, or production authority.
