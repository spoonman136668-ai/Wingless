# Wingless UP-LM5K — restore-round1 extrapolation

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5J bd3c13a67abe238062a3a2be93369617c3c0b864.

## Question

Does the frozen deployment-onset-enriched resource coordinate extrapolate to asynchronous schedules whose earlier restoration occurs at round 1?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Heldout schedules

Six ordered restoration pairs introducing restore round 1:
- budget 0/1 ; throughput 1/3
- budget 1/3 ; throughput 0/1
- budget 0/1 ; throughput 2/4
- budget 2/4 ; throughput 0/1
- budget 0/1 ; throughput 3/5
- budget 3/5 ; throughput 0/1

## Frozen reductions
- budget reduction 1,2
- throughput reduction 1,2

## Frozen topology
- rotations 5,13
- permutations identity, reverse, rotate2

## Frozen resource grid
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

384 points per pooled condition cell; 72 pooled cells.

## Frozen coordinates

Comparator:
- LM5J frozen timing-aware path state without deployment onset.

Primary:
- comparator + first_deliverable_round.

No raw action-start, raw budget, cut-round label, or new scheduler field is added.

## Measurements

Per profile × budget reduction × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

Zero primary spread would show that the enriched deployability representation extrapolates beyond the round-2 restoration domain. Residual spread would expose a new temporal boundary.

## Bounds

Diagnostic only. No adaptive schedule selection, coordinate modification, resource tuning, topology/profile changes, capacity change after results, extra training, live activation, or production authority.
