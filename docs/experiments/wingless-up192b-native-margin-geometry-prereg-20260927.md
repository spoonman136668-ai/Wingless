# Wingless UP-192B — native margin-geometry diagnostic

Status: preregistered scientific diagnostic experiment.

Scientific parent: sealed UP-191B ae9469a8672adc110f6ec1c9382432fef9a9230b.

## Question

Do low-dimensional native margin-geometry features explain the required calibration factor better than native correct count alone?

## Frozen states

Reuse the exact nine UP-191B state profiles and phases 43, 44, 45. No profile search or new perturbation is allowed.

## Frozen probability model

Reuse the pooled phase-26..30 probability model unchanged. No correction is fitted.

## Native features

Before each evaluation path, aggregate across the six native subject states:
- native correct count;
- mean absolute probability margin;
- count of margins with absolute value < 0.01;
- minimum absolute margin.

## Measurements

For each state:
- required mass factor = actual crossing mass / pooled predicted crossing mass;
- the four native geometry features above.

Across all 27 states:
- Pearson association of required factor with each feature.

## Interpretation

If a margin-geometry feature materially outperforms correct count, it becomes a candidate coordinate for the next frozen calibration model. If all remain weak, the calibration state likely requires richer topology/history rather than a low-dimensional scalar summary.

## Bounds

Diagnostic only. No fitted correction, held-out tuning, adaptive feature selection, adaptive profile search, maintenance action, capacity change, live activation, or production authority.
