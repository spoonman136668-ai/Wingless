# Wingless UP-221B — unseen-family confidence rejection

Status: preregistered scientific shadow-calibration diagnostic.

Scientific parent: sealed UP-220B 9cd5ec9eb8bd3b6da3790887b603c4427f0edb8b.

## Question

Does the frozen below-training-min native confidence gate reject harmful forced routing when Wingless encounters history families that were never represented in router/calibrator training?

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
- nearest centroid among the four known regimes

Calibrators:
- one known-regime linear ridge offset model
- ridge lambda 1e-6

Confidence gate:
- reject if route margin < the frozen minimum margin among correctly routed training states.

No retraining.

## Preregistered unseen families

These exact compositions are not members of the four known regimes:
- mixed4b: 2,7,3,8
- observe3: 5,6,8
- store3: 0,2,3
- cross4: 0,5,8,13

The router must still choose one of the four known regimes if forced.

## Untouched evaluation phases

- 288 through 319 inclusive

32 phases per unseen family; 128 route decisions.

## Per-decision outcomes

For each unseen family × phase:
- native route-confidence margin
- rejected by frozen gate or accepted
- forced routed calibration MAE across noncanonical permutations
- static template MAE
- oracle-anchor MAE

A forced route is labeled harmful only for evaluation when:
- routed MAE > static MAE.

The harmful label is never used by the gate.

## Aggregate measurements

Per unseen family and overall:
- decisions
- rejected
- accepted
- harmful forced routes
- harmful rejected
- harmful accepted
- benign rejected
- rejection rate
- harmful-route rejection rate
- mean routed/static/oracle-anchor MAE for rejected and accepted decisions.

## Interpretation

If the frozen gate preferentially rejects harmful forced routing on unseen families, native confidence supports safe abstention outside known calibration regimes. If harmful unseen routes remain high-confidence, the binary gate is in-distribution only.

## Bounds

Shadow diagnostic only. No unseen-family fitting, new centroid, new calibrator, evaluation-derived threshold, adaptive family selection, phase/parity input, nonlinear model, maintenance action, live activation, or production authority.
