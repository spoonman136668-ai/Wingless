# Wingless UP-LM4E — out-of-range resource transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4D 707f192a9cd6838a467161a1b2967d75562fa7e0.

## Question

Does the minimal interpretable resource coordinate transfer to scheduler values outside the resource ranges used to derive it?

## Frozen deadline profiles

- deferred_only
- layout_only
- hybrid_min

## Frozen pressure levels

- 0
- 4
- 8

## Frozen topologies

- rotations 5 and 13
- permutations identity, reverse, rotate2

## Expanded preregistered resource grid

Budgets:
- 2,3,4,5,6,7,8

Action start rounds:
- 1,2,3,4,5,6

Throughput:
- 1,2,3,4,5

The original derivation range was budgets 3–7, starts 2–5, throughput 1–4. Values 2/8, 1/6, and 5 are therefore extrapolative.

210 resource configurations per topology-condition cell.

## Frozen coordinate

Minimal interpretable coordinate:
- reachable_budget
- throughput
- coverage_end_round

No refitting or additional state variable is allowed.

Comparator:
- prior checkpoint coordinate reachable + capacity_by_round2 + capacity_by_round3 + capacity_by_round4.

## Measurements

For each deadline profile × pressure × rotation × permutation × coordinate:
- groups
- groups with multiple underlying resource configurations
- nonzero-spread groups
- maximum failure spread
- mean failure spread

Also report aggregate number of nonzero topology-condition cells per coordinate.

## Interpretation

Zero spread for the frozen three-variable coordinate on this expanded grid supports true resource-state transfer beyond the derivation range. Nonzero spread identifies an extrapolation boundary.

## Bounds

Diagnostic only. No adaptive resource selection, coordinate modification, pressure/profile/topology tuning, capacity change after results, extra training, live activation, or production authority.
