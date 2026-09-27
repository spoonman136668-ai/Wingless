# Wingless UP-LM4Y — post-early-restore transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4X acacd69acdea3a684b23a8f8403fdf68211cfcf7.

## Question

Does cumulative delivery immediately after the earlier resource restoration resolve asynchronous restore ambiguity across all tested budget/throughput timing orders?

## Frozen asynchronous windows

Budget cut/restore; throughput cut/restore:
- 1/3 ; 2/4
- 2/4 ; 1/3
- 1/4 ; 2/5
- 2/5 ; 1/4
- 1/3 ; 3/5
- 3/5 ; 1/3

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen reductions
Budget reduction:
- 1
- 2

Throughput reduction:
- 1
- 2

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
- final + both restoration checkpoints
  - final reachable actions
  - nominal throughput
  - final coverage end round
  - actions before budget restore
  - actions before throughput restore

Primary:
- comparator + cumulative actions immediately after the earlier restoration round

No additional path state.

## Measurements

Per profile × asynchronous window × budget reduction × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

432 topology-condition cells; 864 coordinate summaries.

## Interpretation

Zero primary spread would generalize the transition checkpoint from the localized UP-LM4X boundary into a compact asynchronous-resource path description. Persistent breaks would identify timing structures needing additional path state.

## Bounds

Diagnostic only. No adaptive checkpoint search, resource tuning, coordinate modification, topology/profile selection, capacity change after results, extra training, live activation, or production authority.
