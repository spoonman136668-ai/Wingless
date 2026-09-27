# Wingless UP-LM3O — static-pressure boundary sweep

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3N a99020004253810d937b5e95e12c558fb5643824.

## Question

How far does the frozen compressed temporal resource coordinate transfer as static pre-window pressure increases?

## Frozen scenario

Reuse:
- hybrid_min deadline profile
- rotations 5 and 13
- permutations identity, reverse, rotate2
- budgets 4,5,6
- action start rounds 3,4
- throughput 1,2,3
- six global rounds
- one action per arm per round

## Frozen pressure sweep

Extra ordinary writes per arm after prepressure and before round 0:
- 0
- 1
- 2
- 3
- 4

## Frozen coordinates

Primary:
- reachable_budget + capacity_by_round4

Comparator:
- reachable_budget

Definitions are unchanged from UP-LM3M/UP-LM3N.

## Measurements

Per pressure level × coordinate:
- groups
- groups with nonzero earliest-failure spread
- maximum failure spread
- mean failure spread

## Interpretation

The highest pressure level retaining zero spread for the frozen primary coordinate defines the current validated pressure range. A first nonzero level identifies its boundary without retuning.

## Bounds

Diagnostic only. No adaptive pressure levels, coordinate search, resource tuning, new intervention, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
