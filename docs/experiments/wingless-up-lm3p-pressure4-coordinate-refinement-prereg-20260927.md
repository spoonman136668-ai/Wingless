# Wingless UP-LM3P — pressure-4 coordinate refinement

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3O fe054ae3e21b001c722caaf13a694559104e2e3a.

## Question

At the first pressure level where reachable budget + capacity by round 4 stops being sufficient, does adding earlier deployability restore deterministic collapse?

## Frozen scenario

Reuse UP-LM3O with:
- pressure level = 4 extra pre-window writes per arm
- hybrid_min deadline profile
- rotations 5 and 13
- permutations identity, reverse, rotate2
- budgets 4,5,6
- action start rounds 3,4
- throughput 1,2,3
- six global rounds

## Frozen candidate coordinates

1. reachable_budget + capacity_by_round4
2. reachable_budget + capacity_by_round3 + capacity_by_round4
3. reachable_budget + action_start_round + throughput

No coordinate is added after results.

## Measurements

For each coordinate:
- number of groups
- groups with nonzero earliest-failure spread
- maximum spread
- mean spread

## Interpretation

If adding capacity by round 3 restores zero spread, pressure 4 requires a richer deployment profile but not raw scheduler parameters. If only the raw timing coordinate works, the compact checkpoint representation is insufficient at this boundary.

## Bounds

Diagnostic only. No adaptive coordinate search, pressure tuning, resource tuning, new intervention, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
