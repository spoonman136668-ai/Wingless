# Wingless UP-LM3R — extended rich-profile pressure sweep

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3Q 351cfd5fd87c5dd086a4ee5ecc0da8121f02fb9d.

## Question

Where is the next pressure boundary for the frozen richer deployment profile?

## Frozen scenario

Reuse:
- hybrid_min deadline profile
- rotations 5 and 13
- permutations identity, reverse, rotate2
- budgets 4,5,6
- action start rounds 3,4
- throughput 1,2,3
- six global rounds

## Frozen pressure sweep

Extra ordinary writes per arm after prepressure and before round 0:
- 9
- 10
- 11
- 12
- 13
- 14
- 15
- 16

## Frozen coordinates

Primary:
- reachable_budget + capacity_by_round3 + capacity_by_round4

Comparator:
- reachable_budget + capacity_by_round4

Definitions are unchanged from UP-LM3Q.

## Measurements

Per pressure level × coordinate:
- groups
- groups with nonzero earliest-failure spread
- maximum failure spread
- mean failure spread

## Bounds

Diagnostic only. No adaptive pressure levels, coordinate search, resource tuning, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
