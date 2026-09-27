# Wingless UP-LM3Z — broad profile × pressure transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3Y 9cef7e13cd1caa3faa6a475be3dbf7a432503fc3.

## Question

Does the minimal temporal resource coordinate remain sufficient across deadline structure and pressure variation on the broader resource grid?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen pressure levels
- 0 through 8 extra pre-window writes

## Frozen resource grid
- budgets: 3,4,5,6,7
- action start rounds: 2,3,4,5
- throughput: 1,2,3,4
- rotations: 5 and 13
- permutations: identity, reverse, rotate2
- six global rounds
- earliest_deadline policy

## Frozen coordinates

Prior:
- reachable_budget + capacity_by_round3 + capacity_by_round4

Minimal:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4

Full comparator:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4 + capacity_by_round5

## Measurements

Per deadline profile × pressure:
- global outcome range
- distinct outcome count

Per deadline profile × pressure × coordinate:
- groups
- nonzero-spread groups
- maximum spread
- mean spread

## Bounds

Diagnostic only. No adaptive coordinate/profile/pressure selection, resource tuning, topology search, capacity change, extra training, live activation, or production authority.
