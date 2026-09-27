# Wingless UP-202B — exact-history trend transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-201B 6529a78c26a1ce8a86830918dd9a4673b94b3fe1.

## Question

For history families whose exact permutation identity does not transfer as a fixed level, is the calibration dependence drifting systematically across phase?

## Frozen compositions

- mixed4: 0,5,1,6
- observe4: 5,6,7,8

All 24 permutations per composition.

## Frozen training and evaluation

Training phases:
- 55
- 56
- 57

Untouched evaluation phases:
- 73
- 74
- 75

## Frozen predictors

Per composition:
- constant baseline: composition-wide training mean.
- identity mean: exact-permutation training mean.
- identity trend: exact-permutation least-squares linear trend over training phase number, extrapolated unchanged.

No feature extraction, phase/parity class, nonlinear trend, breakpoint, held-out fitting, or model selection.

## Measurements

Per composition on untouched phases:
- baseline mean/max absolute error;
- identity-mean mean/max absolute error;
- identity-trend mean/max absolute error;
- points where identity trend beats identity mean;
- points where identity trend beats baseline.

## Interpretation

Held-out improvement from the frozen per-history trend would support systematic phase drift. Failure to improve would indicate that the phase dependence is not captured by a simple stable linear trajectory.

## Bounds

Shadow calibration only. No held-out fitting, adaptive trend selection, nonlinear fitting, correction activation, maintenance action, capacity change, live activation, or production authority.
