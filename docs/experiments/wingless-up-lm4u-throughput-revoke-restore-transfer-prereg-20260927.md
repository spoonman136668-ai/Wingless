# Wingless UP-LM4U — temporary throughput-loss transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4T 3b6332e7c7e9cfb3e447ac24c899fd9b1f474ef6.

## Question

Does the minimal path law transfer from temporary total-budget loss to temporary throughput loss?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen throughput-loss schedules
- cut round 1, restore round 3
- cut round 2, restore round 4
- cut round 3, restore round 5

Temporary throughput reductions:
- 1 action/round
- 2 actions/round

Outside the revoked window, nominal throughput is restored.

## Frozen resource grid
- budgets 4,5,6,7
- action starts 2,3,4,5
- nominal throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

64 configurations per topology-condition cell.

## Frozen topology
- rotations 5,13
- permutations identity, reverse, rotate2

## Frozen coordinates
Comparator:
- actual final reachable actions
- nominal restored throughput
- actual final coverage end

Primary:
- comparator variables
- cumulative deliverability strictly before restoration

## Measurements
Per profile × schedule × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- max/mean failure spread

108 cells; 216 coordinate summaries.

## Interpretation

Zero primary spread would show that the pre-restore deliverability checkpoint captures path dependence across a different resource dimension, not just total-budget revocation.

## Bounds

Diagnostic only. No adaptive schedule/reduction/resource selection, checkpoint search, coordinate modification, extra training, live activation, or production authority.
