# Wingless UP-208B — family-holdout native offset transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-207B 8ce15c4234206ebb1bab899a88ea7d822df5a87e.

## Question

Does the native-state phase-offset relationship generalize to an entire history family excluded from fitting?

## Frozen families

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

Each family is held out once. The ridge model is fit on the other three families only.

## Frozen phases

Template phases:
- 55,56,57

Training phases:
- 58 through 72 inclusive

Untouched evaluation phases:
- 91,92,93

## Frozen native predictors

Canonical history only:
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

No phase number, parity, family identity, held-out anchor measurement, or held-out fitting.

## Frozen model

Per holdout fold:
- pooled ridge regression on the three non-heldout families;
- ridge lambda 1e-6;
- target = canonical required-factor phase shift from its frozen template mean.

## Frozen comparators

- static exact-permutation template mean;
- oracle one-anchor shift using the heldout family's canonical required factor at evaluation time.

The canonical point is excluded from scoring.

## Measurements

Per heldout family:
- static mean/max absolute error;
- family-holdout native-offset mean/max absolute error;
- oracle-anchor mean/max absolute error;
- native-better-than-static count.

## Interpretation

Material improvement on a family never used to fit the native-offset mapping supports cross-family native-state generalization. Failure in a specific family identifies family-dependent state geometry rather than a universal offset law.

## Bounds

Shadow calibration only. No heldout fitting, adaptive feature selection, family-ID input, phase/parity input, nonlinear model, anchor measurement for native prediction, maintenance action, live activation, or production authority.
