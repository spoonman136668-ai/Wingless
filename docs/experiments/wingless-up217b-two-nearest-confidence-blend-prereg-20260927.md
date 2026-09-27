# Wingless UP-217B — training-confidence two-nearest calibration blend

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-216B 1de95869087acfed0f23858a0108c456332c8c67.

## Question

When the frozen native router flags a route as uncertain using only the UP-216B training-derived confidence threshold, does blending the two nearest frozen family calibrators outperform both ordinary routing and static fallback?

## Frozen router/calibrators

Unchanged:
- regimes: mixed4, observe4, store4, cross3
- template phases 55,56,57
- training phases 58 through 72
- native features: native_correct_count, mean_absolute_margin, near_zero_margin_count, min_absolute_margin
- pooled standardization
- nearest-centroid router
- one family-specific ridge offset model per regime
- ridge lambda 1e-6

## Frozen uncertainty threshold

Exactly UP-216B:
- route all 60 training regime-phase states;
- retain correctly routed training states;
- threshold = minimum confidence margin among those correct training states.

No evaluation data contributes.

## Frozen fallback rule

For each evaluation state:
- if confidence margin >= threshold: use ordinary nearest-family routed calibrator;
- if confidence margin < threshold: compute the offset predicted by each of the two nearest frozen family calibrators and use their equal arithmetic mean.

No weighting parameter is fit.

## Untouched evaluation phases

- 196 through 207 inclusive

48 route decisions total.

## Comparators

- static exact-permutation template mean
- ordinary routed calibration
- two-nearest confidence blend
- oracle-family calibrator
- oracle-anchor upper bound

Canonical points excluded from permutation scoring.

## Measurements

Overall:
- training-derived threshold
- route accuracy
- blend decision count

Per true regime:
- route correct / total
- blend decision count
- static/routed/blend/oracle-family/oracle-anchor MAE
- blend-better-than-routed count
- blend-better-than-static count

## Interpretation

Improved error on uncertain states would support native uncertainty as a bounded model-routing signal without evaluation tuning. Worsening error would reject two-nearest blending and indicate that uncertainty is detectable but not safely correctable with the current calibrator family.

## Bounds

Shadow calibration only. No evaluation-derived threshold, adaptive weighting, evaluation-label fitting, retraining, phase/parity input, nonlinear router/model, anchor measurement for routed/blend/oracle-family prediction, maintenance action, live activation, or production authority.
