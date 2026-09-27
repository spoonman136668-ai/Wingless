# Wingless UP-214B — late-phase native-routed offset transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-213B e39f9e0c9cc3c453d205d0fac4f13237fa530d28.

## Question

Does the frozen native router plus frozen family-specific calibrators remain useful across a much longer later-phase holdout?

## Frozen regimes
- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

## Frozen template phases
- 55,56,57

## Frozen training phases
- 58 through 72 inclusive

## Frozen native features
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

## Frozen router
- pooled standardization from training data;
- one centroid per regime;
- nearest Euclidean centroid.

## Frozen family-specific calibrators
- one linear ridge model per regime;
- intercept plus four standardized native features;
- ridge lambda 1e-6;
- target = canonical required-factor phase shift from frozen canonical template mean.

No retraining or recalibration from UP-213B design.

## Untouched evaluation phases
- 136 through 159 inclusive

24 phases per regime; 96 route decisions total.

## Predictions
Static:
- frozen exact-permutation template mean.

Native-routed:
- infer regime from current native state;
- apply routed regime's frozen offset model;
- add predicted shift to true history's frozen permutation template.

Oracle-family:
- true regime's frozen offset model, no anchor measurement.

Oracle-anchor upper bound:
- measured true canonical shift at that evaluation phase.

Canonical permutation excluded from error scoring.

## Measurements
Per true regime:
- route correct / total;
- static MAE/max error;
- routed MAE/max error;
- oracle-family MAE/max error;
- oracle-anchor MAE/max error;
- routed-better-than-static count.

Overall:
- route accuracy across 96 decisions.

## Interpretation

Stable routing/calibration across this longer holdout would support persistent native regime identification rather than a short-window coincidence. Concentrated mixed4 errors would confirm the remaining ambiguity is localized rather than system-wide.

## Bounds

Shadow calibration only. No evaluation-label fitting, adaptive model selection, retraining, phase/parity input, anchor measurement for routed/oracle-family prediction, nonlinear router/model, maintenance action, live activation, or production authority.
