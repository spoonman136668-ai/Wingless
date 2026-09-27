# Wingless UP-191B — native-state response curve

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-190B endogenous-context shift bae5008f361cca459694cffc82fd2302162920fb.

## Question

The frozen two-factor correction failed once Wingless-native correct count moved from the 620/623 regime into 635–647. Does the calibration factor required by those shifted states vary systematically with native correct count, or does perturbation history still matter at the same or similar native count?

## Frozen predictor

Reuse the unchanged pooled probability model trained on unshifted phases 26..30.

Do not apply or fit a new correction.

For each evaluation state, measure only the descriptive factor that would make aggregate predicted crossing mass equal observed crossing mass:

required_factor = actual_crossings / pooled_predicted_crossings_sum.

This factor is an observation, not a fitted predictor.

## Frozen perturbation histories

Evaluate six histories fixed before seeing their results:

- store_forward: 0,1,2,3;
- store_reverse: 3,2,1,0;
- observe_forward: 5,6,7,8;
- observe_reverse: 8,7,6,5;
- mixed_forward: 0,5,1,6;
- mixed_reverse: 6,1,5,0.

Evaluation phases:
- 43;
- 44;
- 45.

No profile is added, removed, or searched after observing native state.

## Measurements

Per profile × phase:
- native correct count;
- actual crossings;
- pooled predicted crossing mass;
- pooled predicted/actual ratio;
- absolute pooled mass-ratio error;
- required calibration factor.

Aggregate by native correct count:
- number of observations;
- number of distinct perturbation histories;
- minimum required factor;
- maximum required factor;
- factor spread.

Also report Pearson correlation between native correct count and required factor across all fixed evaluation points.

## Interpretation

A small within-count factor spread plus a strong systematic count/factor relationship would support native correct count as a useful continuous calibration coordinate.

Large factor spread at the same count, or weak/non-monotonic structure, would show that count alone is insufficient and developmental/perturbation history must remain part of the state description.

## Bounds

Descriptive shadow analysis only. No new correction is fit or applied. No held-out fitting, adaptive profile search, phase/parity predictor input, maintenance action, capacity change, extra model calls, live activation, or production authority.
