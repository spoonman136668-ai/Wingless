# Wingless UP-LM4G — two-burst dynamic-pressure transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4F 0f9c211e92c1eb784c4d93be06cef2bf89e718f1.

## Question

Does the frozen minimal resource law remain sufficient when two separate pressure bursts arrive during the action window?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen two-burst timing pairs

All distinct unordered pairs of global rounds 0 through 5:
- (0,1), (0,2), (0,3), (0,4), (0,5)
- (1,2), (1,3), (1,4), (1,5)
- (2,3), (2,4), (2,5)
- (3,4), (3,5)
- (4,5)

At each selected round, exactly one extra ordinary write per arm is injected.

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

Per deadline profile × burst pair × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

270 topology-condition cells; 540 coordinate summaries.

## Interpretation

Zero primary spread under all two-burst timings supports transfer of the compact resource law to repeated nonstationary pressure. Breaks at particular burst pairs would identify temporal interactions not represented by static resource coverage.

## Bounds

Diagnostic only. No adaptive burst-pair selection, resource tuning, coordinate modification, topology/profile selection, capacity change, extra training, live activation, or production authority.
