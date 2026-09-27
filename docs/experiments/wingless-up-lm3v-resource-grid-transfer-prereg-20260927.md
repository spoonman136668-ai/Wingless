# Wingless UP-LM3V — broader resource-grid transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3U e037a27f4419c0ac8ac43c4b0ff3fee2ec998cfe.

## Question

Does the frozen rich temporal resource coordinate remain sufficient on resource values not used in its original validation grid?

## Frozen scenario

- deadline profile: hybrid_min
- pressure: 4 extra pre-window writes
- rotations: 5 and 13
- permutations: identity, reverse, rotate2
- six global rounds
- earliest_deadline policy

Hybrid pressure 4 is retained because it is resource-sensitive and requires the rich coordinate in prior experiments.

## Broader preregistered resource grid

- budgets: 3,4,5,6,7
- action start rounds: 2,3,4,5
- throughput: 1,2,3,4

80 resource configurations total.

## Frozen coordinates

Simple:
- reachable_budget + capacity_by_round4

Rich:
- reachable_budget + capacity_by_round3 + capacity_by_round4

No new checkpoint or coordinate is added after results.

## Measurements

Across the 80 configurations:
- global outcome range and distinct outcome count.

Per coordinate:
- groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

## Interpretation

Zero rich-coordinate spread with a nonzero global resource range supports transfer beyond the original resource grid.

## Bounds

Diagnostic only. No adaptive resource selection, coordinate search, pressure tuning, topology search, capacity change, extra training, live activation, or production authority.
