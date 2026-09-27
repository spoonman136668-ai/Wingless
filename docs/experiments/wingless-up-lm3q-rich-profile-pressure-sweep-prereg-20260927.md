# Wingless UP-LM3Q — rich-profile pressure sweep

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3P 9edc932ad23072bc0cd5bf99586a50718f2a9e26.

## Question

How far does the richer deployment profile remain sufficient as static pre-window pressure increases beyond the first boundary?

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
- 4
- 5
- 6
- 7
- 8

## Frozen coordinates

Primary:
- reachable_budget + capacity_by_round3 + capacity_by_round4

Comparator:
- reachable_budget + capacity_by_round4

Definitions are unchanged from UP-LM3P.

## Measurements

Per pressure level × coordinate:
- groups
- groups with nonzero earliest-failure spread
- maximum failure spread
- mean failure spread

## Interpretation

The highest pressure level retaining zero spread for the richer profile defines its validated range. The first nonzero level identifies its next boundary without retuning.

## Bounds

Diagnostic only. No adaptive pressure levels, coordinate search, resource tuning, new intervention, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
