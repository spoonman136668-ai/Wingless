# Wingless UP-212B — native history-regime identification

Status: preregistered scientific shadow-state diagnostic.

Scientific parent: sealed UP-211B 6a64c22843a6cd663e0e56f849fc5dfd4af9423e.

## Question

Can Wingless identify which history regime it is currently in from native state geometry alone, without being handed the family label?

## Frozen regimes

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

## Frozen native features

Canonical history only:
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

No phase number, parity, permutation identity, explicit family ID, calibration target, future state, or outcome label is available to the classifier at evaluation time.

## Frozen training phases

- 58 through 72 inclusive

For each regime:
- compute standardized native features;
- compute one centroid in standardized feature space.

Standardization statistics are pooled over the 60 training points.

## Frozen classifier

Nearest Euclidean centroid.

No learned nonlinear boundary and no heldout fitting.

## Untouched evaluation phases

- 121 through 132 inclusive

48 evaluation points total.

## Measurements

- total accuracy;
- per-regime correct / total;
- 4x4 confusion counts;
- mean nearest-centroid distance;
- mean confidence margin = second-nearest distance - nearest distance.

Comparator:
- uniform four-regime chance level = 0.25.

## Interpretation

High later-phase accuracy would show that family-specific calibration need not require an external regime label: the regime is recoverable from native state itself. Low accuracy would mean the current native geometry is insufficient to route safely between family-specific calibrators.

## Bounds

Shadow diagnostic only. No evaluation-label fitting, adaptive feature selection, phase/parity input, calibration-target input, nonlinear classifier, live activation, or production authority.
