# Wingless UP-219B — native route-confidence error calibration

Status: preregistered scientific shadow-calibration diagnostic.

Scientific parent: sealed UP-218B 94aaf811a0e2ab0bdb62a197691e028360bac08d.

## Question

Does Wingless's own frozen native route-confidence margin identify states where routed calibration error is elevated?

## Frozen machinery

Exactly the UP-218B training machinery:
- regimes: mixed4, observe4, store4, cross3
- template phases 55,56,57
- training phases 58 through 72
- native features: native_correct_count, mean_absolute_margin, near_zero_margin_count, min_absolute_margin
- pooled standardization
- one centroid per regime
- nearest-centroid routing
- one family-specific linear ridge offset model per regime
- ridge lambda 1e-6
- confidence threshold = minimum confidence margin among correctly routed training states

No retraining or threshold change.

## Untouched evaluation phases

- 232 through 255 inclusive

24 phases per regime; 96 route decisions.

## Frozen confidence strata

- uncertain: confidence margin < frozen training threshold
- confident: confidence margin >= frozen training threshold

## Measurement unit

One regime × phase route decision.

For each route decision:
- route correct / incorrect
- confidence margin
- routed calibration MAE across that regime's noncanonical permutations
- static calibration MAE across the same permutations
- oracle-family calibration MAE

## Aggregate measurements

Per confidence stratum:
- decisions
- misroutes
- misroute rate
- mean routed MAE
- mean static MAE
- mean oracle-family MAE

Per regime:
- uncertain / confident decision counts
- misroutes
- mean routed MAE in each stratum

## Interpretation

Higher misroute rate and routed MAE in the frozen uncertain stratum would demonstrate that native internal confidence is informative about calibration failure risk before an external error label is observed.

## Bounds

Shadow diagnostic only. No evaluation-derived threshold, adaptive weighting, retraining, phase/parity input, nonlinear router/model, anchor measurement, maintenance action, live activation, or production authority.
