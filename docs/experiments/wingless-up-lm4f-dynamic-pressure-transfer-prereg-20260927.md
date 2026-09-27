# Wingless UP-LM4F — dynamic-pressure resource-law transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4E 81ca89b4c8eb786a513391e6dca377e17ecc1c3e.

## Question

Does the frozen minimal resource law remain sufficient when extra pressure arrives during the action window instead of before it?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen dynamic pressure
Exactly one extra ordinary write per arm is injected at one fixed global round:
- 0
- 1
- 2
- 3
- 4
- 5

Each burst round is a separate preregistered condition.

## Frozen topology
- rotations 5 and 13
- permutations identity, reverse, rotate2

## Frozen resource grid
- budgets 3,4,5,6,7
- action starts 2,3,4,5
- throughput 1,2,3,4
- six global rounds
- earliest_deadline scheduling

80 resource configurations per topology-condition cell.

## Frozen coordinates

Primary:
- reachable_budget + throughput + coverage_end_round

Comparator:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4

No coordinate changes after results.

## Measurements

Per deadline profile × burst round × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

108 topology-condition cells; 216 coordinate summaries.

## Interpretation

Zero primary spread under all burst timings supports transfer of the three-variable resource law to nonstationary pressure. Breaks at specific burst rounds identify where pressure timing adds state not represented by resource coverage alone.

## Bounds

Diagnostic only. No adaptive burst timing, resource tuning, coordinate modification, topology/profile selection, capacity change, extra training, live activation, or production authority.
