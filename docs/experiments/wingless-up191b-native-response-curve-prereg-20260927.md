# Wingless UP-191B — native calibration response curve

Status: preregistered scientific diagnostic experiment.

Scientific parent: sealed UP-190B bae5008f361cca459694cffc82fd2302162920fb.

## Question

After the binary correction failed out of regime, does the calibration factor required by Wingless vary systematically with its native correct-count state, or does perturbation history matter even at the same native count?

## Frozen pooled model

Reuse the pooled probability model trained on unshifted phases 26..30. No correction is fitted in this experiment.

## Frozen state profiles

Evaluation-only history profiles:
- control: no extra update;
- store1: 0;
- store2: 0,1;
- store4: 0,1,2,3;
- observe1: 5;
- observe2: 5,6;
- observe4: 5,6,7,8;
- mixed2: 0,5;
- mixed4: 0,5,1,6.

Profiles are fixed before the run. No profile is added or removed based on results.

## Evaluation phases

- 43;
- 44;
- 45.

## Measurements

Per profile × phase:
- native correct count;
- actual crossing mass;
- pooled predicted crossing mass;
- required multiplicative mass factor = actual / predicted;
- pooled absolute mass-ratio error.

Across all points:
- Pearson association between native correct count and required factor;
- for each repeated native count, minimum/maximum required factor and spread;
- maximum within-count factor spread.

## Interpretation

A strong smooth count/factor relationship with low within-count spread would justify testing a continuous native-state correction next. Large within-count spread would show that native correct count alone is insufficient and that history/state topology must also enter calibration.

## Bounds

Diagnostic only. No fitted correction, no held-out tuning, no adaptive profile search, no phase/parity predictor input, no maintenance action, no capacity change, no live activation, or production authority.
