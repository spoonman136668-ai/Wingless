# Wingless UP-LM2Z — budget-onset frontier

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM2Y e955419fdb6a38573cf6007555ecd4583f80889b.

## Question

At what total action budget does native earliest-deadline scheduling first begin to outperform fixed-order scheduling reproducibly?

## Frozen scenarios

Reuse UP-LM2Y exactly:
- identity rotations 3 and 11;
- arm-order permutations identity, reverse, rotate2;
- six arms per scenario;
- deferred levels 4,5,6;
- suffix_reported and alternating_reported layouts;
- exact recall cap 16;
- 12 global pressure rounds;
- at most one action per round.

## Frozen budgets

Test total action budgets:
- 4;
- 5;
- 6;
- 7;
- 8.

## Policies

- baseline;
- earliest_deadline;
- fixed_order.

Active policies have identical total budget and one-action-per-round throughput.

## Measurements

Per budget aggregate across both identity rotations and all three permutations:
- baseline failed;
- earliest failed;
- fixed failed;
- actions used;
- failures prevented;
- prevention per action;
- earliest-deadline advantage = fixed_failed - earliest_failed.

Primary onset criterion:
- first budget with aggregate earliest-deadline advantage > 0.

Secondary robustness:
- number of the six rotation×permutation scenarios at that budget where earliest-deadline strictly beats fixed-order.

## Bounds

Counterfactual only. No live activation, future schedule input, adaptive budget, capacity change, extra training, semantic priority, or production authority.
