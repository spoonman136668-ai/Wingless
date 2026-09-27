# Wingless UP-LM4W — asynchronous restore checkpoints

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4V 7a0bf7dc75c59be1d2c83cda189ba17d263a80c5.

## Question

When budget and throughput recover at different times, is one pre-full-restoration path checkpoint still sufficient, or are both resource-specific restore checkpoints required?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen asynchronous windows

Budget cut/restore; throughput cut/restore:
- 1/3 ; 2/4
- 2/4 ; 1/3
- 1/4 ; 2/5
- 2/5 ; 1/4
- 1/3 ; 3/5
- 3/5 ; 1/3

Budget and throughput restore rounds are distinct in every schedule.

## Frozen temporary reductions
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

A. final:
- final reachable actions
- nominal throughput
- final coverage end round

B. final+pre-full-restore:
- final coordinate
- actions delivered immediately before the later of the two restoration rounds

C. final+both-restores:
- final coordinate
- actions delivered immediately before budget restoration
- actions delivered immediately before throughput restoration

No coordinate modification after results.

## Measurements

Per profile × asynchronous schedule × budget reduction × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

432 topology-condition cells; 1296 coordinate summaries.

## Interpretation

If pre-full-restore alone remains exact, one checkpoint continues to summarize the relevant resource path. If it breaks while both-restores is exact, separate restore-path memory is necessary.

## Bounds

Diagnostic only. No adaptive schedule/reduction selection, resource tuning, coordinate modification, topology/profile selection, capacity change after results, extra training, live activation, or production authority.
