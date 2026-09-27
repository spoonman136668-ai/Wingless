# Wingless UP-220B — native confidence-risk strata

Status: preregistered scientific shadow-calibration diagnostic.

Scientific parent: sealed UP-219B f02efe61c0a43ad1f1c4fc2abc8f34e2d559bcd8.

## Question

Does heldout routing/calibration risk increase monotonically across confidence strata defined entirely from training-state native route margins?

## Frozen training machinery

Exactly the accepted router/calibrators:
- regimes: mixed4, observe4, store4, cross3
- template phases 55,56,57
- training phases 58 through 72
- native features: native_correct_count, mean_absolute_margin, near_zero_margin_count, min_absolute_margin
- pooled standardization
- nearest-centroid router
- one family-specific linear ridge offset model per regime
- ridge lambda 1e-6

## Frozen confidence boundaries

Using only correctly routed training states:
- minimum confidence margin
- 25th percentile
- 50th percentile
- 75th percentile

Percentiles use the sorted training margins and floor index:
- floor((n-1) × percentile)

Five frozen strata:
1. below_training_min
2. training_low: min <= margin < q25
3. training_mid: q25 <= margin < q50
4. training_high: q50 <= margin < q75
5. training_top: margin >= q75

No evaluation-derived boundaries.

## Untouched evaluation phases

- 256 through 287 inclusive

32 phases per regime; 128 route decisions.

## Per-decision measurements

- confidence margin
- frozen confidence stratum
- route correct / incorrect
- routed calibration MAE across noncanonical permutations
- static calibration MAE
- oracle-family calibration MAE

## Aggregate measurements

Per stratum:
- decisions
- misroutes
- empirical misroute rate
- mean routed MAE
- mean static MAE
- mean oracle-family MAE

## Interpretation

A monotonic increase in misroute/error risk as confidence falls would convert the binary warning from UP-219B into a coarse calibrated native risk curve. Non-monotonicity would bound confidence to binary gating rather than graded risk estimation.

## Bounds

Shadow diagnostic only. No evaluation-derived boundaries, adaptive binning, retraining, phase/parity input, nonlinear router/model, anchor measurement, maintenance action, live activation, or production authority.
