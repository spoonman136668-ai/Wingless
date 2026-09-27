# Wingless UP-218B — late confidence-blend transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-217B d2af9082b10cf0362861bd8ce402cffb9b0bfe73.

## Question

Does the training-confidence two-nearest calibration blend remain useful on a longer, later heldout phase window with no retraining or threshold changes?

## Frozen regimes
- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

## Frozen training machinery

Exactly UP-217B:
- template phases 55,56,57
- training phases 58 through 72
- native features: native_correct_count, mean_absolute_margin, near_zero_margin_count, min_absolute_margin
- pooled standardization
- nearest-centroid router
- one family-specific ridge offset model per regime
- ridge lambda 1e-6
- uncertainty threshold = minimum confidence margin among correctly routed training states
- uncertain route fallback = equal mean of the two nearest frozen calibrators

No evaluation-derived fitting.

## Untouched evaluation phases
- 208 through 231 inclusive

24 phases per regime; 96 route decisions total.

## Comparators
- static exact-permutation template
- ordinary nearest-family routed calibration
- confidence-triggered two-nearest blend
- oracle-family calibrator
- oracle-anchor upper bound

Canonical permutation excluded from error scoring.

## Measurements
Overall:
- training-derived threshold
- route accuracy
- blend decision count

Per regime:
- route correct / total
- blend decision count
- static/routed/blend/oracle-family/oracle-anchor MAE
- blend-better-than-routed count
- blend-better-than-static count

## Interpretation

Persistent improvement over ordinary routing, particularly in mixed4, supports native uncertainty-aware calibration as a transferable state-routing mechanism. Loss of benefit bounds the blend to the earlier heldout window.

## Bounds

Shadow calibration only. No evaluation-derived threshold, adaptive weighting, evaluation-label fitting, retraining, phase/parity input, nonlinear router/model, anchor measurement for routed/blend/oracle-family prediction, maintenance action, live activation, or production authority.
