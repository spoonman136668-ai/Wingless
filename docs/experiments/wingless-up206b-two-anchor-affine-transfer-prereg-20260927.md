# Wingless UP-206B — two-anchor affine history transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-205B 57554c8fe3a98d9ed21ffeb5f0c1a15d23bee3c7.

## Question

Does the relative history geometry across phases change mainly by a shared offset plus a shared scale?

## Frozen composition families

- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

All exact permutations per family.

## Frozen template phases

- 55
- 56
- 57

Untouched evaluation phases:
- 85
- 86
- 87

## Frozen anchors

Two anchors per composition:
1. canonical permutation;
2. reverse permutation.

No anchor search.

## Frozen predictors

Template residual:
- exact-permutation training mean minus composition-wide template mean.

One-anchor comparator:
- use canonical anchor only;
- offset = observed canonical factor - canonical template residual;
- prediction = offset + template residual.

Two-anchor affine:
- use canonical and reverse anchor observations;
- fit a shared scale and offset from those two anchors only;
- prediction = offset + scale × frozen template residual.

Both anchors are excluded from scoring.

No held-out feature fitting or model selection.

## Measurements

Per composition:
- one-anchor mean/max absolute error;
- two-anchor affine mean/max absolute error;
- affine-better-than-one-anchor count;
- fitted phase scale values.

## Interpretation

Improvement, especially on cross3, would support a compact phase transformation of stable history geometry: common-mode offset plus common scale.

## Bounds

Shadow calibration only. No anchor search, nonlinear transform, held-out fitting, maintenance action, live activation, or production authority.
