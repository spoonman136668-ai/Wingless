# Wingless UP-222B — absolute-distance OOD rejection

Status: preregistered scientific shadow-calibration diagnostic.

Scientific parent: sealed UP-221B 4d9a0a4c40b2a65c6ada6a511b9e24c986e6aee9.

## Question

Can a training-frozen absolute distance-to-known-regime gate catch harmful unseen-family routing that the relative confidence-margin gate misses?

## Frozen known regimes and training

Known regimes:
- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

Template phases:
- 55,56,57

Training phases:
- 58 through 72

Native features:
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

Router:
- pooled standardization
- nearest centroid among four known regimes

Calibrators:
- one known-regime linear ridge offset model
- ridge lambda 1e-6

## Frozen training-derived gates

Margin gate:
- reject if nearest-vs-second-nearest distance margin is below the minimum margin among correctly routed training states.

Absolute-distance gate:
- reject if nearest-centroid distance is above the maximum nearest-centroid distance among correctly routed training states.

Union gate:
- reject if either frozen gate rejects.

No unseen-family information enters either threshold.

## Unseen families

Same preregistered unseen compositions as UP-221B:
- mixed4b: 2,7,3,8
- observe3: 5,6,8
- store3: 0,2,3
- cross4: 0,5,8,13

## Fresh untouched evaluation phases

- 320 through 351 inclusive

32 phases per unseen family; 128 route decisions.

## Evaluation label

A forced route is harmful only for evaluation when:
- routed calibration MAE > static template MAE.

The harmful label is never used by any gate.

## Measurements

Per gate × unseen family, plus overall:
- decisions
- rejected / accepted
- harmful forced routes
- harmful rejected
- harmful accepted
- benign rejected
- rejection rate
- harmful-route rejection rate
- benign-rejection rate

15 summaries total.

## Interpretation

If absolute distance or the union gate recovers harmful high-margin OOD cases with bounded benign rejection, Wingless gains a safer abstention signal outside known calibration regimes. If not, centroid geometry is insufficient for OOD risk.

## Bounds

Shadow diagnostic only. No unseen-family fitting, new centroid, new calibrator, evaluation-derived threshold, adaptive family selection, phase/parity input, nonlinear model, maintenance action, live activation, or production authority.
