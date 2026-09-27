# Wingless UP-LM4J — mid-window throughput-drop transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4I 43523c3750ccd0a99df30564c0505b589e61467e.

## Question

Does the static compact resource law remain sufficient when resource delivery changes during the action window, and can an explicitly dynamic deployability coordinate recover exact outcome collapse?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen resource change

Initial throughput:
- 2,3,4 actions per active round

At one fixed global drop round:
- 2
- 3
- 4
- 5

throughput decreases by exactly one for that round and all later rounds:
- 2 -> 1
- 3 -> 2
- 4 -> 3

No throughput value is allowed below 1.

## Frozen resource grid
- budgets 3,4,5,6,7
- action starts 2,3,4,5
- initial throughput 2,3,4
- six global rounds
- earliest_deadline scheduling

60 resource configurations per topology-condition cell.

## Frozen topology
- rotations 5 and 13
- permutations identity, reverse, rotate2

## Frozen coordinates

Static comparator:
- original reachable_budget computed from initial throughput
- initial throughput
- original coverage_end_round computed from initial throughput

Dynamic primary:
- dynamic reachable budget under the actual throughput schedule
- initial throughput
- post-drop throughput
- dynamic coverage_end_round under the actual throughput schedule

Drop round is fixed within each condition cell and therefore is not part of the grouping key.

## Measurements

Per deadline profile × drop round × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

72 topology-condition cells; 144 coordinate summaries.

## Interpretation

If the static coordinate breaks but the dynamic coordinate restores zero spread, resource risk depends on actual temporal deployability rather than nominal initial capacity. If both remain exact, the prior compact law is more robust than expected. If the dynamic coordinate also breaks, additional resource-schedule state is required.

## Bounds

Diagnostic only. No adaptive drop timing, throughput recovery, resource tuning, coordinate modification, topology/profile selection, capacity change after results, extra training, live activation, or production authority.
