# Wingless UP-LM3X — endpoint checkpoint ablation

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3W dc7684ed002aaac557e955eac52ab443da67e8bb.

## Question

Which omitted endpoint checkpoint caused the broader resource-grid boundary: early capacity at round 2, late capacity at round 5, or both?

## Frozen scenario

- deadline profile: hybrid_min
- pressure: 4 extra pre-window writes
- rotations: 5 and 13
- permutations: identity, reverse, rotate2
- six global rounds
- earliest_deadline policy

## Frozen resource grid

- budgets: 3,4,5,6,7
- action start rounds: 2,3,4,5
- throughput: 1,2,3,4

80 configurations.

## Frozen coordinates

A. prior rich:
- reachable_budget + capacity_by_round3 + capacity_by_round4

B. early endpoint added:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4

C. late endpoint added:
- reachable_budget + capacity_by_round3 + capacity_by_round4 + capacity_by_round5

D. full:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4 + capacity_by_round5

No coordinate is added or removed after results.

## Measurements

Per coordinate:
- groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

Global:
- outcome range
- distinct outcomes

## Interpretation

If only one endpoint restores zero spread, that endpoint is the missing temporal state. If neither alone restores zero but the full vector does, both endpoint capacities are jointly necessary.

## Bounds

Diagnostic only. No adaptive coordinate search, resource tuning, pressure tuning, topology search, capacity change, extra training, live activation, or production authority.
