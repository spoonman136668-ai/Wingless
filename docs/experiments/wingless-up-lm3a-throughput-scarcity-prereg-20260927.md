# Wingless UP-LM3A — throughput scarcity

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM2Z 8f94f6a8a8a387115ae2c9da7f29182989457e09.

## Question

UP-LM2Z found that earliest-deadline ordering first gains value at total budget 5 and becomes universally better by budget 6 under one action per round. Is that advantage caused by total action scarcity, by per-round timing/throughput scarcity, or both?

## Frozen scenarios

Reuse UP-LM2Z:
- identity rotations 3 and 11;
- arm permutations identity, reverse, rotate2;
- six independent arms per scenario;
- deferred levels 4,5,6;
- suffix_reported and alternating_reported layouts;
- exact recall cap 16;
- 12 global pressure rounds.

## Frozen total budgets

- 5;
- 6.

## Frozen per-round throughput

Allow at most:
- 1 action per round;
- 2 actions per round;
- 3 actions per round.

The total action budget remains fixed. Unused capacity in one round does not increase the total budget.

## Policies

- baseline: no actions;
- earliest_deadline: repeatedly choose the currently shortest native countdown until round throughput or total budget is exhausted;
- fixed_order: repeatedly choose the first pending arm in the current permuted order under the same limits.

Each action marks one pending dependency completed before that round's pressure writes.

## Measurements

Per budget × throughput aggregated across both rotations and all three permutations:
- baseline failed;
- earliest failed;
- fixed failed;
- actions used;
- failures prevented;
- prevention per action;
- earliest-deadline advantage.

## Interpretation

If more per-round throughput reduces or removes the ordering advantage at fixed total budget, timing scarcity is causal. If the advantage remains, total budget allocation itself is the dominant factor.

## Bounds

Counterfactual only. No live activation, future schedule input, adaptive budget, adaptive throughput, capacity change, extra training, semantic priority, or production authority.
