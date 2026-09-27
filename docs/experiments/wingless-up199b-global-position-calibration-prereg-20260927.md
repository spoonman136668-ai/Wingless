# Wingless UP-199B — global positional history calibration

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-198B cb538e82e05c087d2939ba130339b04870aeb953.

## Question

Does global item position capture held-out history-dependent calibration better than local pairwise transitions?

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
- 64
- 65
- 66

## Frozen representations

Pairwise model:
- intercept + 12 directed adjacent-pair indicators, unchanged.

Global positional model:
- intercept + 16 binary item-by-position indicators.
- canonical item ranks 0..3 crossed with sequence positions 0..3.

Both use ridge lambda 1e-6.

Constant baseline:
- training mean required factor per composition.

No feature selection, interaction search, phase identity, parity, or held-out fitting.

## Measurements

Per composition on untouched phases:
- baseline mean/max absolute error;
- pairwise mean/max absolute error;
- positional mean/max absolute error;
- positional-better-than-pairwise count;
- positional-better-than-baseline count.

## Interpretation

Held-out positional improvement would support global order placement as part of the missing history state. No improvement would push the representation beyond simple local or absolute-position encodings.

## Bounds

Shadow calibration only. No held-out fitting, adaptive feature selection, architecture search, phase/parity input, correction activation, maintenance action, capacity change, live activation, or production authority.
