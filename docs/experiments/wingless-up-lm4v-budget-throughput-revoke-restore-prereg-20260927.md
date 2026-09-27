# Wingless UP-LM4V — combined budget + throughput revoke/restore transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4U.

## Question

Does the frozen path-memory coordinate remain sufficient when total action budget and delivery throughput are both temporarily reduced and later restored?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen revoke/restore windows
- cut round 1, restore round 3
- cut round 2, restore round 4
- cut round 3, restore round 5

## Frozen temporary reductions
Budget reduction:
- 1
- 2

Throughput reduction:
- 1
- 2

Both reductions apply over the same cut/restore window.

## Frozen topology
- rotations 5 and 13
- permutations identity, reverse, rotate2

## Frozen resource grid
- budgets 4,5,6,7
- action starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

64 resource configurations per topology-condition cell.

## Frozen coordinates

Comparator:
- final = final reachable actions + nominal throughput + final coverage end round

Primary:
- final+prerestore = final coordinate + actions delivered before restoration

No new state is added after results.

## Measurements

Per profile × revoke/restore window × budget reduction × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

216 topology-condition cells; 432 coordinate summaries.

## Interpretation

Zero primary spread under simultaneous temporary budget and throughput loss would support pre-restoration delivered capacity as a general path-memory variable across combined resource perturbations.

## Bounds

Diagnostic only. No adaptive schedule/reduction selection, resource tuning, coordinate modification, topology/profile selection, capacity change after results, extra training, live activation, or production authority.
