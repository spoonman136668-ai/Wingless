# Wingless UP-LM4K — one-round throughput outage transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4J 232895f37c3a49ecb4b87a45d61ce2c9f9ce734e.

## Question

Does the compact resource law remain sufficient when action delivery becomes completely unavailable for one round and then recovers?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen outage rounds
- 2
- 3
- 4
- 5

At the selected outage round:
- scheduler throughput is exactly zero for that round.
- the original throughput resumes on the next round.

## Frozen resource grid
- budgets 3,4,5,6,7
- action starts 2,3,4,5
- nominal throughput 1,2,3,4
- six global rounds
- earliest_deadline scheduling

80 resource configurations per topology-condition cell.

## Frozen topology
- rotations 5 and 13
- permutations identity, reverse, rotate2

## Frozen coordinates

Static comparator:
- original reachable_budget
- nominal throughput
- original coverage_end_round

Dynamic primary:
- reachable budget under the actual outage schedule
- nominal throughput
- dynamic coverage_end_round under the actual outage schedule

Outage round is fixed within each condition cell and is not part of the grouping key.

## Measurements

Per deadline profile × outage round × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

72 topology-condition cells; 144 coordinate summaries.

## Interpretation

If the static coordinate breaks but the dynamic coordinate restores zero spread, exact risk depends on actual temporal deployability rather than nominal resource settings. If both remain exact, the compact law is robust even to temporary total unavailability.

## Bounds

Diagnostic only. No adaptive outage timing, outage duration, resource tuning, coordinate modification, topology/profile selection, capacity change after results, extra training, live activation, or production authority.
