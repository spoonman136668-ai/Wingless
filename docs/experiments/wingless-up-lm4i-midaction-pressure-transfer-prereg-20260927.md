# Wingless UP-LM4I — mid-action dynamic-pressure transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4H 83e3c6dfa0f15e3c9b465ee1f3f6729ba2aee0ee.

## Question

Does the frozen minimal resource law remain sufficient when pressure arrives after the first scheduler action slot but before the remaining throughput slots in the same round?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen burst rounds
- 0,1,2,3,4,5

Each condition injects exactly one extra ordinary write per arm at the selected round.

## Frozen intra-round ordering

For each round:
1. if the scheduler is active and budget remains, execute at most the first action slot;
2. if this is the selected burst round, inject one extra ordinary write per arm;
3. execute any remaining throughput slots for that round;
4. apply the normal round write.

For throughput = 1 the burst is effectively post-action. For throughput > 1 it lands between action slots.

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

Zero primary spread would show the compact resource law survives pressure inserted inside a multi-action round. Nonzero spread would identify a micro-ordering state absent from round-level coverage summaries.

## Bounds

Diagnostic only. No adaptive burst timing, resource tuning, coordinate modification, topology/profile selection, capacity change, extra training, live activation, or production authority.
