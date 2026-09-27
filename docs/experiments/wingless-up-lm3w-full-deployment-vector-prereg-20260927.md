# Wingless UP-LM3W — full deployment-vector transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3V 66a6e40d82e43973f7e715ef5a4ccee4fa7cf886.

## Question

Does the broader resource-grid failure boundary disappear when the resource state includes deployable capacity across the full start-round range?

## Frozen scenario

- deadline profile: hybrid_min
- pressure: 4 extra pre-window writes
- rotations: 5 and 13
- permutations: identity, reverse, rotate2
- six global rounds
- earliest_deadline policy

## Frozen resource grid

- budgets: 3,4,5,6,7
- action start rounds: 2,3,4,5
- throughput: 1,2,3,4

80 configurations.

## Frozen coordinates

Prior rich comparator:
- reachable_budget + capacity_by_round3 + capacity_by_round4

Full deployment vector:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4 + capacity_by_round5

No checkpoint is added after results.

## Measurements

Per coordinate:
- groups
- nonzero-spread groups
- max spread
- mean spread

Global:
- outcome range
- distinct outcomes

## Interpretation

If the full vector restores zero spread while the prior rich coordinate does not, the missing state is temporal deployability at the omitted early/late checkpoints rather than an unrelated scheduler variable.

## Bounds

Diagnostic only. No adaptive coordinate search, resource tuning, pressure tuning, topology search, capacity change, extra training, live activation, or production authority.
