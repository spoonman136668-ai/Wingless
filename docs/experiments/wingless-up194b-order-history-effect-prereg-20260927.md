# Wingless UP-194B — order-history calibration effect

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-193B 50aaa41f36c10283ca677124f5f07cb2b1be294f.

## Question

When perturbation composition is held fixed, does changing only update order change Wingless's required aggregate failure-risk calibration?

## Frozen predictor

Reuse the pooled phase-26..30 probability model unchanged. No correction is fit or applied.

## Frozen matched order pairs

Each pair contains exactly the same update indices:
- store4: 0,1,2,3 versus 3,2,1,0;
- observe4: 5,6,7,8 versus 8,7,6,5;
- mixed4: 0,5,1,6 versus 6,1,5,0;
- cross3: 0,5,13 versus 13,5,0.

## Untouched evaluation phases

- 49;
- 50;
- 51.

## Measurements

Per pair × phase:
- forward and reverse native correct count;
- forward and reverse required mass factor;
- absolute factor difference;
- absolute native-count difference;
- whether the two histories end at the same native correct count.

Aggregate:
- mean and maximum factor difference;
- mean factor difference among same-count pairs;
- number of same-count pairs with nonzero factor difference.

## Interpretation

A nonzero calibration difference under identical update composition is direct evidence of order/history dependence. A nonzero difference even when native correct count matches shows count cannot summarize that history effect.

## Bounds

Diagnostic only. No fitted correction, no adaptive pair search, no held-out fitting, no phase/parity predictor input, no maintenance action, no capacity change, no live activation, or production authority.
