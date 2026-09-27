# Wingless UP-190B — endogenous context-shift calibration

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-189B 62429a84c67c974773efe5ddb51d7324be880507.

## Question

Does the frozen UP-188B/UP-189B native-context mass correction remain calibrated when evaluation begins from Wingless-native states perturbed away from the repeating unshifted regime?

## Frozen predictor

Reuse unchanged:
- pooled probability model trained on unshifted phases 26..30;
- context_A if native correct count <= 621, context_B otherwise;
- context_A factor = 1.0103904235372925;
- context_B factor = 0.8983526543771945.

No factors, thresholds, bins, or model cells are recomputed.

## Frozen evaluation shifts

Before generating evaluation rows, apply one fixed sequence of ordinary historical updates to the native gate state:
- store_shift: indices 0,1,2,3;
- observe_shift: indices 5,6,7,8;
- mixed_shift: indices 0,5,1,6.

These profiles are fixed before seeing their resulting native correct counts. There is no adaptive search for a convenient state.

## Untouched evaluation phases

- 43;
- 44;
- 45.

## Measurements

Per shift profile × phase × model:
- native correct count;
- whether that count lies outside the previously observed 620/623 pair;
- actual crossing count;
- predicted crossing mass and absolute mass-ratio error;
- Brier score;
- expected calibration error;
- AUROC;
- average precision.

Aggregate mean and maximum absolute mass-ratio error for pooled and corrected models.

## Interpretation

Persistence would support calibration beyond simple phase repetition. Degradation would locate the boundary of the two-factor correction. If the frozen perturbations fail to leave the 620/623 native-count pair, the experiment is informative as a failed OOD construction and must not be retuned post hoc.

## Bounds

Shadow only. Evaluation-state perturbation only. No held-out fitting, phase/parity predictor input, adaptive shift search, maintenance action, capacity change, extra model calls, live activation, or production authority.
