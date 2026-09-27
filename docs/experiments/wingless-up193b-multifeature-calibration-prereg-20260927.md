# Wingless UP-193B — compact multi-feature native calibration

Status: preregistered scientific calibration experiment.

Scientific parent: sealed UP-192B f60d940777f04fdd5aee86b2f8483a0d7faa8f6c.

## Question

Can a small frozen combination of native state features calibrate failure mass on disjoint states where every single tested scalar was weak?

## Training data

Use only the sealed UP-191B/192B state family:
- phases 43, 44, 45;
- profiles control, store1, store2, store4, observe1, observe2, observe4, mixed2, mixed4.

Target is required pooled mass factor = actual crossing mass / pooled predicted crossing mass.

## Frozen features

Exactly four:
1. native correct count;
2. mean absolute native margin;
3. near-zero margin count at |margin| < 0.01;
4. minimum absolute native margin.

Standardize each feature using training data only.

## Frozen model

Linear ridge regression with intercept.
- ridge lambda = 1e-6;
- ridge applies to feature coefficients, not intercept;
- no feature selection;
- no clipping of the fitted factor itself;
- per-event corrected probabilities remain bounded to [0,1].

## Disjoint evaluation

Phases:
- 46;
- 47;
- 48.

Profiles never used in training:
- store3: 0,1,2;
- observe3: 5,6,7;
- mixed3: 0,5,1;
- report2: 13,14;
- cross3: 0,5,13.

## Measurements

Per evaluation state:
- actual required factor;
- fitted factor;
- pooled absolute mass-ratio error;
- corrected absolute mass-ratio error.

Aggregate mean and maximum errors.

## Bounds

No held-out fitting, no adaptive feature selection, no profile search, no phase/parity predictor input, no maintenance action, no capacity change, no live activation, or production authority.
