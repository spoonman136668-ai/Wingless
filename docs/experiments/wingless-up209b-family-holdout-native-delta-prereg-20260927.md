# Wingless UP-209B — family-holdout native trajectory-delta transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-208B 322ceca5b89ae730260360ef0f87ce927e2da837.

## Question

Are short native-state changes more family-invariant than absolute native-state geometry for predicting the common calibration shift?

## Frozen families

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

Each family is held out once. No heldout-family points are used to fit its model.

## Frozen phases

Template phases:
- 55,56,57

Training target phases:
- 59 through 72 inclusive

Each training point may use native geometry from the target phase and immediately preceding phase.

Untouched evaluation phases:
- 94,95,96

Each evaluation point may use only current and immediately preceding native geometry.

## Frozen predictors

Current native geometry:
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

One-step native deltas:
- delta_native_correct_count
- delta_mean_absolute_margin
- delta_near_zero_margin_count
- delta_min_absolute_margin

No phase number, parity, family identity, anchor measurement, future state, or heldout fitting.

## Frozen model

Per family-holdout fold:
- pooled linear ridge regression on the other three families;
- intercept plus 8 standardized predictors;
- ridge lambda 1e-6;
- target = canonical required-factor shift from frozen canonical template mean.

## Frozen comparators

- static exact-permutation template mean;
- oracle one-anchor shift measured from the heldout canonical history at evaluation time.

The canonical permutation is excluded from scoring.

## Measurements

Per heldout family:
- static mean/max absolute error;
- native-level+delta mean/max absolute error;
- oracle-anchor mean/max absolute error;
- native-better-than-static count.

## Interpretation

Improvement over UP-208B, especially on mixed4 and observe4, would indicate that native trajectory direction carries family-invariant calibration information missing from static geometry. Failure would reject this small trajectory descriptor and force a move away from generic pooled calibration.

## Bounds

Shadow calibration only. No heldout fitting, adaptive feature selection, family-ID input, phase/parity input, nonlinear model, future-state input, anchor measurement for native prediction, maintenance action, live activation, or production authority.
