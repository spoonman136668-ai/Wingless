# Wingless UP-LM4B — interior checkpoint ablation

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4A ddeb4b19bbb9a904832242ad7c8a976187787829.

## Question

With reachable budget and early deployability at round 2 retained, are both interior capacity checkpoints still necessary?

## Frozen conditions

Exactly the UP-LM4A topology-resolved matrix:
- deadline profiles deferred_only, layout_only, hybrid_min
- pressure levels 0 through 8
- rotations 5 and 13
- permutations identity, reverse, rotate2
- budgets 3,4,5,6,7
- action starts 2,3,4,5
- throughput 1,2,3,4
- six rounds, earliest_deadline policy

## Frozen coordinates

A:
- reachable_budget + capacity_by_round2 + capacity_by_round3

B:
- reachable_budget + capacity_by_round2 + capacity_by_round4

C control:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4

## Measurements

For each of 162 topology-condition cells × coordinate:
- groups
- nonzero-spread groups
- maximum spread
- mean spread

## Interpretation

If A or B retains zero spread in all cells, the omitted interior checkpoint is redundant. If only C remains exact, both cap3 and cap4 are necessary.

## Bounds

Diagnostic only. No adaptive coordinate selection, resource tuning, topology selection, capacity change, extra training, live activation, or production authority.
