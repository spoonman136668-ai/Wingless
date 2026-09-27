# Wingless UP-216B — training-derived confidence fallback

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-215B 1b2da82bd2321a12bbcd98607b00aa9235980c1e.

## Question

Can a confidence threshold derived strictly from training data identify uncertain native routes well enough for a bounded static-calibration fallback to improve heldout performance?

## Frozen regimes
- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

## Frozen router/calibrators
Unchanged from UP-213B / UP-214B:
- template phases 55,56,57
- training phases 58 through 72
- native features: correct count, mean absolute margin, near-zero margin count, minimum absolute margin
- pooled standardization
- nearest-centroid router
- one family-specific ridge offset model per regime, lambda 1e-6

## Frozen confidence threshold

Using only the 60 training regime-phase states:
1. route every training state with the frozen centroids;
2. retain correctly routed training states;
3. threshold = minimum confidence margin (second-nearest distance - nearest distance) among those correct training states.

No evaluation data contributes to the threshold.

## Frozen evaluation
Fresh phases:
- 184 through 195 inclusive

48 route decisions total.

## Frozen fallback rule
- if evaluation confidence margin >= training-derived threshold: use normal native-routed family calibrator;
- if margin < threshold: use static exact-permutation template mean, i.e. no phase offset correction.

Comparator predictions:
- static
- ordinary routed
- confidence-fallback
- oracle-family calibrator
- oracle-anchor upper bound

Canonical points excluded from permutation scoring.

## Measurements
Overall:
- training-derived threshold
- route accuracy
- fallback decision count

Per true regime:
- route correct/total
- fallback count
- static/routed/fallback/oracle-family/oracle-anchor MAE
- fallback-better-than-routed count
- fallback-better-than-static count

## Interpretation

Improved heldout error would show Wingless can use its own routing uncertainty to avoid unreliable calibration without evaluation tuning. Worsening error would falsify static fallback even if uncertainty detection itself is valid.

## Bounds
Shadow calibration only. No evaluation-derived threshold, adaptive fallback, evaluation-label fitting, retraining, phase/parity input, nonlinear router/model, anchor measurement for routed/fallback/oracle-family prediction, maintenance action, live activation, or production authority.
