# Wingless UP-LM3T — deadline-profile transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3S 26b995d12ffcb92afbcbcf637ac08a4ca846d277.

## Question

Does the frozen rich temporal resource profile transfer across distinct deadline structures in the informative, resource-sensitive pressure regime?

## Frozen resource grid

- budgets 4,5,6
- action start rounds 3,4
- throughput 1,2,3
- six global rounds
- rotations 5 and 13
- permutations identity, reverse, rotate2
- scheduling policy earliest_deadline

## Frozen pressure

Exactly four extra ordinary writes per arm after deadline prepressure and before round 0.

Pressure 4 is chosen because prior work showed it is resource-sensitive and exposed the simpler coordinate's boundary.

## Frozen deadline profiles

- deferred_only: prepressure target depends only on deferred level
- layout_only: prepressure target depends only on layout
- hybrid_min: minimum of deferred and layout targets, unchanged control

## Frozen coordinates

Primary:
- reachable_budget + capacity_by_round3 + capacity_by_round4

Comparator:
- reachable_budget + capacity_by_round4

## Measurements

Per deadline profile × coordinate:
- coordinate groups
- groups with nonzero earliest-failure spread
- maximum spread
- mean spread

Per deadline profile:
- global minimum and maximum earliest-failed outcome across resource configurations
- global outcome range
- distinct outcome count

## Interpretation

Zero within-group spread with nonzero global resource range supports transfer of the rich resource profile across deadline structures.

## Bounds

Diagnostic only. No profile selection after results, coordinate search, pressure tuning, resource tuning, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
