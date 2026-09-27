# Wingless UP-192B — native geometry residual

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-191B native-state response e6cf06e3bc44ca603df9cf9dc8f54b87dc2d36ba.

## Question

When native correct-count is insufficient across history families, do frozen endogenous margin-geometry features explain the remaining variation in required failure-mass calibration?

## Frozen predictor

Reuse the pooled probability model trained on phases 26..30. No correction is fit or applied.

## Frozen history profiles

Union of the two prior fixed diagnostic families:
- control;
- store1, store2, store_forward, store_reverse;
- observe1, observe2, observe_forward, observe_reverse;
- mixed2, mixed_forward, mixed_reverse.

Evaluation phases: 43, 44, 45.

## Frozen native features

Measured from the shifted Wingless gate state before evaluation:
- native correct count;
- mean signed margin;
- mean absolute margin;
- minimum margin;
- count with absolute margin below 0.01;
- total negative-margin mass.

For each profile × phase, also measure required aggregate calibration factor = actual crossings / pooled predicted crossing mass.

Report Pearson association of each frozen native feature with required factor. No feature selection controls the harness.

## Interpretation

If a geometry feature tracks required factor where correct-count does not, it becomes a candidate additional native risk coordinate. If none does, the relevant state descriptor is likely higher-dimensional or trajectory-dependent.

## Bounds

Diagnostic only. No fitted correction, no held-out tuning, no adaptive profile or feature search, no phase/parity predictor input, no maintenance action, no capacity change, no live activation, or production authority.
