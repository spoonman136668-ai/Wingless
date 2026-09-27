# Wingless UP-213B — native-routed family-specific offset calibration

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-212B 8eb0000d98f4fd10d524f0b7f53b660f39226561.

## Question

Can Wingless use only native state to route itself to a family-specific calibration model and recover the phase-wide offset without an external family label?

## Frozen regimes

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

## Frozen template phases

- 55,56,57

Exact-permutation template means are frozen per regime.

## Frozen training phases

- 58 through 72 inclusive

Native features:
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

Pooled training statistics standardize all four features.

## Frozen router

- one centroid per regime in standardized native feature space;
- nearest Euclidean centroid;
- no evaluation-label fitting.

## Frozen family-specific calibrators

For each regime separately:
- linear ridge model with intercept plus the four standardized native features;
- ridge lambda 1e-6;
- target = canonical required-factor phase shift from that regime's frozen canonical template mean.

## Untouched evaluation phases

- 133,134,135

## Frozen predictions

Static:
- exact-permutation template mean, no phase correction.

Native-routed:
- classify regime from current native geometry;
- apply that predicted regime's frozen offset model;
- add predicted shift to the true history's frozen permutation template.

Oracle-family:
- use the true regime's frozen offset model, but no anchor measurement.

Oracle-anchor upper bound:
- measure the true canonical required factor at the evaluation phase and apply its exact common shift.

Canonical points are excluded from permutation-error scoring.

## Measurements

Per true regime:
- routing accuracy by phase;
- static MAE/max error;
- native-routed MAE/max error;
- oracle-family MAE/max error;
- oracle-anchor MAE/max error;
- routed-better-than-static count.

Overall:
- route accuracy across 12 regime-phase decisions.

## Interpretation

If native-routed error approaches oracle-family error, native state is sufficient to choose the needed regime-specific calibration autonomously. A gap concentrated in mixed4 would localize the remaining representation problem to routing ambiguity.

## Bounds

Shadow calibration only. No evaluation-label fitting, adaptive model selection, phase/parity input, anchor measurement for native-routed/oracle-family prediction, nonlinear router/model, maintenance action, live activation, or production authority.
