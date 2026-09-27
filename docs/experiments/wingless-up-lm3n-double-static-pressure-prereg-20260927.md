# Wingless UP-LM3N — double-static-pressure transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3M d067be9ceea97f83c61bec184767a56494e8fa33.

## Question

Does the compressed temporal resource coordinate found in UP-LM3M transfer when static pre-window pressure doubles?

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

## Frozen pressure change

After hybrid prepressure and before global round 0, add exactly two ordinary writes per arm instead of one.

No other parameter changes.

## Frozen coordinate

Primary coordinate:
- reachable_budget + capacity_by_round4

Comparator:
- reachable_budget alone

Definitions are unchanged from UP-LM3M.

## Measurements

For each coordinate:
- number of groups
- groups with nonzero earliest-failure spread
- maximum failure spread
- mean failure spread

## Interpretation

Zero spread for the frozen compressed coordinate would show transfer across a stronger static load. Nonzero spread would define its current pressure boundary.

## Bounds

Diagnostic only. No adaptive coordinate search, pressure tuning, resource tuning, new intervention, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
