# Wingless UP-223B — training-manifold OOD rejection

Status: preregistered scientific shadow-calibration diagnostic.

Scientific parent: sealed UP-222B c1c98b5fa5ffbc5760139ecf234406b2e874d53c.

## Question

Does distance to the actual cloud of correctly routed training states reject harmful unseen-family routes that remain inside the frozen centroid/margin envelope?

## Frozen known regimes and model

Exactly UP-222B:
- mixed4, observe4, store4, cross3
- template phases 55,56,57
- training phases 58 through 72
- four native features
- pooled standardization
- nearest-centroid routing
- family-specific ridge calibrators, lambda 1e-6.

## Frozen baseline gate

Reject when either:
- route margin < minimum correctly-routed training margin; or
- nearest centroid distance > maximum correctly-routed training centroid distance.

## Frozen manifold gate

Training exemplars:
- only correctly routed training states;
- same standardized four-feature vectors.

Training-manifold threshold:
- for each correctly routed training state, compute Euclidean distance to its nearest *other* correctly routed training state;
- threshold = maximum of those leave-one-out nearest-neighbor distances.

Augmented reject:
- baseline reject OR nearest-training-exemplar distance > frozen manifold threshold.

No unseen-family outcome is used to define any threshold.

## Fresh heldout evaluation

Unseen families:
- mixed4b: 2,7,3,8
- observe3: 5,6,8
- store3: 0,2,3
- cross4: 0,5,8,13

Evaluation phases:
- 352 through 383 inclusive

32 decisions per family; 128 total.

A forced route is evaluation-labeled harmful only when routed MAE > static MAE.

## Measurements

Per family and overall, for baseline and augmented gates:
- decisions
- harmful / benign routes
- rejected / accepted
- harmful rejected / accepted
- benign rejected
- harmful rejection rate
- benign rejection rate.

Also report frozen margin, centroid-distance, and manifold-distance thresholds.

## Interpretation

Improved harmful-route rejection, especially on store3, would show that harmful aliases occupy holes in the training manifold even when centroid confidence looks normal. No improvement would indicate the current four-feature native state is genuinely aliased and requires a new state dimension rather than a new gate.

## Bounds

Shadow diagnostic only. No unseen-family fitting, no new native feature, no new centroid/calibrator, no evaluation-derived threshold, no adaptive gate selection, no nonlinear model, no phase/parity input, no maintenance action, no live activation, or production authority.
