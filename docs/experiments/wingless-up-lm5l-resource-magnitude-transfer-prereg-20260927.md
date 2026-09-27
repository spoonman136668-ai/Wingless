# Wingless UP-LM5L — out-of-range resource magnitude transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5K 72fa0ef77f4c3e482bc3308aa533f62ec41c85ea.

## Question

Does the frozen deployment-onset-enriched resource coordinate remain exact when resource magnitudes move outside the previously tested budget/start/throughput ranges?

## Frozen schedules

Exactly the six UP-LM5K restore-round-1 schedules:
- budget 0/1 ; throughput 1/3
- budget 1/3 ; throughput 0/1
- budget 0/1 ; throughput 2/4
- budget 2/4 ; throughput 0/1
- budget 0/1 ; throughput 3/5
- budget 3/5 ; throughput 0/1

## Frozen profiles / reductions / topology

Profiles:
- deferred_only
- layout_only
- hybrid_min

Reductions:
- budget 1,2
- throughput 1,2

Topology:
- rotations 5,13
- permutations identity, reverse, rotate2

Scheduler:
- earliest_deadline
- six global rounds.

## Heldout resource grid

Budgets:
- 2,3,8,9

Action starts:
- 1,2,5,6

Throughput:
- 1,2,5,6

These include lower and higher magnitudes than the previously tested budget 4..7, start 2..5, throughput 2..5 grid.

64 resource configurations per schedule; 384 pooled points per condition cell; 72 pooled cells.

## Frozen coordinates

Comparator:
- timing-aware path state without deployment onset.

Primary:
- the same state + first_deliverable_round.

No raw action-start, raw budget, cut-round label, or new scheduler field is added.

## Measurements

Per profile × budget reduction × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread.

144 summaries.

## Interpretation

Zero primary spread would support the enriched deployability state across substantial resource-magnitude extrapolation. Residual spread would identify a magnitude boundary and justify a targeted missing-state diagnostic.

## Bounds

Diagnostic only. No adaptive grid selection, coordinate modification, resource tuning after results, topology/profile changes, extra training, live activation, or production authority.
