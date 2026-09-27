# Wingless UP-LM2Y — budget-allocation replication

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM2X 79a4c00a527e037ac4209c81ecb7b9f49eb44cdc.

## Question

Does the UP-LM2X advantage of native earliest-deadline scheduling persist across different total action budgets and different arm index orders?

## Frozen scenarios

Identity rotations:
- 3;
- 11.

Each scenario contains the same six independent arms:
- deferred levels 4,5,6;
- suffix_reported and alternating_reported layouts;
- value shift 0;
- exact recall cap 16 per arm;
- 12 global pressure rounds.

## Frozen arm-order permutations

- identity;
- reverse;
- rotate2.

The arm content is unchanged; only array/tie order is permuted.

## Frozen total action budgets

Across the full 12-round scenario:
- 4;
- 8;
- 12 actions maximum.

At most one action may occur per round.

## Policies

- baseline: no action;
- earliest_deadline: among arms with a pending dependency, choose the smallest current native countdown;
- fixed_order: choose the first pending arm in the current permuted order.

Both active policies share the same total budget and per-round throughput.

## Measurements

Per rotation × permutation × budget × policy:
- actions used;
- completed originals;
- failed originals.

Aggregate for each budget:
- failures prevented versus baseline;
- prevention per action;
- earliest-deadline advantage over fixed-order.

## Interpretation

A persistent earliest-deadline advantage across budgets and order permutations supports genuine scheduling value in the native countdown rather than one fortunate tie/order arrangement.

## Bounds

Counterfactual only. No live activation, future schedule input, capacity change, extra training, adaptive action budget, semantic priority, or production authority.
