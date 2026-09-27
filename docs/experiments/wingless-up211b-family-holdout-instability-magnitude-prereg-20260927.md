# Wingless UP-211B — family-holdout instability magnitude

Status: preregistered scientific shadow-risk experiment.

Scientific parent: sealed UP-210B 7842b7338be06a648f326fad3b21a5ffd828304e.

## Question

Even when drift direction is family-dependent, can absolute native-state change predict the magnitude of calibration instability in an unseen history family?

## Frozen families

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

Each family is held out once.

## Frozen training phases

- 59 through 72 inclusive

Target:
- absolute value of canonical required-factor drift from phase-1 to phase.

## Frozen evaluation phases

- 109 through 120 inclusive

## Frozen predictors

Absolute one-step native changes only:
- abs(delta_native_correct_count)
- abs(delta_mean_absolute_margin)
- abs(delta_near_zero_margin_count)
- abs(delta_min_absolute_margin)

No signed direction, absolute state level, family identity, phase number, parity, future state, or anchor measurement.

## Frozen model

Per holdout fold:
- linear ridge regression on the other three families;
- intercept plus four standardized predictors;
- ridge lambda 1e-6;
- negative predicted magnitudes clipped to zero.

## Frozen comparators

1. zero-instability predictor: 0.
2. training-mean predictor: mean absolute drift magnitude across the training families and phases.

## Measurements

Per heldout family:
- zero baseline MAE;
- training-mean baseline MAE;
- native magnitude MAE;
- max errors;
- native-better-than-training-mean count;
- Pearson correlation between predicted and actual magnitudes where variance permits.

## Interpretation

Beating the training-mean baseline on an unseen family would support a family-invariant native instability signal even when signed calibration dynamics remain family-specific.

## Bounds

Shadow diagnostic only. No heldout fitting, adaptive feature selection, family-ID input, signed-drift input, absolute-state input, phase/parity input, nonlinear model, future-state input, live activation, or production authority.
