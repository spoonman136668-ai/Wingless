# Wingless UP-LM3U — deadline-profile × pressure map

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3T 50df407515a97ac6f8141a233a6381dc074e8dab.

## Question

Across the informative pressure range, which deadline structures require the richer temporal resource profile?

## Frozen resource grid

- budgets 4,5,6
- action start rounds 3,4
- throughput 1,2,3
- six global rounds
- rotations 5 and 13
- permutations identity, reverse, rotate2
- earliest_deadline policy

## Frozen deadline profiles

- deferred_only
- layout_only
- hybrid_min

## Frozen pressure levels

0 through 8 extra pre-window writes per arm.

## Frozen coordinates

Simple:
- reachable_budget + capacity_by_round4

Rich:
- reachable_budget + capacity_by_round3 + capacity_by_round4

## Measurements

Per profile × pressure:
- global outcome range and distinct outcome count.

Per profile × pressure × coordinate:
- groups
- nonzero-spread groups
- max spread
- mean spread.

## Bounds

Diagnostic only. No adaptive pressure/profile/coordinate selection, resource tuning, topology search, capacity change, extra training, live activation, or production authority.
