# Wingless UP-200B — combined history calibration

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-199B 0b107d58dded41b4c239a52e61160fa38aac1489.

## Question

Does combining fixed directed pairwise transitions with fixed global item-position features improve held-out calibration beyond either representation alone?

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
- 67
- 68
- 69

## Frozen representations

- constant baseline: training mean factor
- pairwise: 12 directed adjacent-pair features
- position: 16 item-by-position features
- combined: concatenated 28 pairwise + position features

All learned models use ridge lambda 1e-6.

No feature selection, architecture search, phase identity, parity, or held-out fitting.

## Measurements

Per composition on untouched phases:
- mean and max absolute error for baseline, pairwise, position, and combined models
- combined-better-than-pairwise count
- combined-better-than-position count
- combined-better-than-baseline count

## Interpretation

A held-out improvement from the frozen combined representation would support distributed history information across local transition and global position cues. No improvement would indicate these simple order coordinates remain insufficient.

## Bounds

Shadow calibration only. No held-out fitting, adaptive feature selection, architecture search, correction activation, maintenance action, capacity change, live activation, or production authority.
