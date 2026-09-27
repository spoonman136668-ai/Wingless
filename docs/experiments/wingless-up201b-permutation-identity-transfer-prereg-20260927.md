# Wingless UP-201B — exact-history identity transfer

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-200B 56e9730fa7863308566b4c2489f8dd9baca8d8d0.

## Question

Is the history-dependent calibration mapping stable across phases even when it cannot be compressed by simple pairwise, triplet, or positional features?

## Frozen compositions

- mixed4: 0,5,1,6
- observe4: 5,6,7,8

All 24 permutations per composition.

## Frozen training and evaluation

Training phases:
- 55
- 56
- 57

Untouched evaluation phases:
- 70
- 71
- 72

## Frozen predictors

- constant baseline: composition-wide mean required factor on training phases.
- permutation identity: per-composition, per-exact-permutation mean required factor across training phases.

No interpolation, feature extraction, phase identity, parity, or held-out fitting.

## Measurements

Per composition:
- baseline mean/max absolute error;
- identity mean/max absolute error;
- points where identity beats baseline.

## Interpretation

Held-out improvement from exact identity would show a stable but non-compressible history mapping. No improvement would indicate that calibration dependence itself shifts across phases rather than being a persistent sequence signature.

## Bounds

Shadow calibration only. No held-out fitting, adaptive identity selection, feature search, correction activation, maintenance action, capacity change, live activation, or production authority.
