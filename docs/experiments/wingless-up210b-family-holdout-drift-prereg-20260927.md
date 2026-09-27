# Wingless UP-210B — family-holdout native drift prediction

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-209B 3d286792ae090cbf748314d2a678532b73a77424.

## Question

Even when absolute calibration is family-dependent, can native state change predict the direction and magnitude of phase-to-phase calibration drift in an unseen history family?

## Frozen families

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

Each family is held out once. No heldout-family point is used for fitting.

## Frozen training phases

Target phases:
- 59 through 72 inclusive

Each training target is:
- canonical_required_factor(phase) - canonical_required_factor(phase-1)

## Frozen evaluation phases

- 97 through 108 inclusive

## Frozen native predictors

Only one-step changes in canonical native geometry:
- delta_native_correct_count
- delta_mean_absolute_margin
- delta_near_zero_margin_count
- delta_min_absolute_margin

No absolute native levels, family identity, phase number, parity, future state, or anchor measurement.

## Frozen model

Per family-holdout fold:
- pooled linear ridge regression on the other three families;
- intercept plus four standardized native-delta predictors;
- ridge lambda 1e-6.

## Frozen comparator

Zero-drift predictor:
- predicted calibration drift = 0.

## Measurements

Per heldout family:
- zero-drift mean/max absolute error;
- native-drift mean/max absolute error;
- native-better-than-zero count;
- nonzero actual-drift points;
- correct drift-direction count on nonzero points.

## Interpretation

Improvement over zero drift on an unseen family would support a family-invariant dynamic signal even though absolute calibration remains family-specific. Failure would reject this compact pooled drift model.

## Bounds

Shadow diagnostic only. No heldout fitting, adaptive feature selection, family-ID input, phase/parity input, absolute-state input, nonlinear model, future-state input, live activation, or production authority.
