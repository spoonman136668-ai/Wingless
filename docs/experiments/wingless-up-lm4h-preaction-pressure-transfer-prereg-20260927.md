# Wingless UP-LM4H — pre-action dynamic-pressure transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4G b3842ea2b08f31ac05cec5ae659464a1ecfb7fc2.

## Question

Does the frozen minimal resource law remain sufficient when an extra pressure write arrives before the scheduler's action opportunity within a round, rather than after it?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen burst rounds
- 0
- 1
- 2
- 3
- 4
- 5

Each condition injects exactly one extra ordinary write per arm at the selected round.

Critical ordering:
- the burst is applied first;
- then the scheduler may act for that round;
- then the ordinary round write occurs.

This differs only in intra-round ordering from UP-LM4F.

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

Zero primary spread would show the compact resource law is insensitive to whether pressure arrives immediately before or after the scheduler's action opportunity. Breaks localized to particular rounds would expose an intra-round ordering state absent from the current coordinate.

## Bounds

Diagnostic only. No adaptive burst timing, resource tuning, coordinate modification, topology/profile selection, capacity change, extra training, live activation, or production authority.
