# Wingless UP-LM4L — consecutive two-round throughput outage transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4K 3ae95dcc94af0a6399fd5a933de277b12a169dc1.

## Question

Does the compact resource law remain sufficient when action delivery is completely unavailable for two consecutive rounds before recovering?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen outage windows
- rounds 2–3
- rounds 3–4
- rounds 4–5

At both outage rounds:
- scheduler throughput is exactly zero.
- nominal throughput resumes after the outage window when a later round remains.

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
- reachable budget under the actual two-round outage schedule
- nominal throughput
- dynamic coverage_end_round under the actual outage schedule

The outage window is fixed within each condition cell and is not part of the grouping key.

## Measurements

Per deadline profile × outage window × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

54 topology-condition cells; 108 coordinate summaries.

## Interpretation

A static-coordinate break with dynamic recovery would directly show that actual deployability must be represented when resource delivery is interrupted for a material fraction of the action window. If both remain exact, the current outcome law is remarkably invariant to short resource outages.

## Bounds

Diagnostic only. No adaptive outage timing/duration, resource tuning, coordinate modification, topology/profile selection, capacity change after results, extra training, live activation, or production authority.
