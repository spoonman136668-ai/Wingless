# Wingless UP-LM5E — heldout restore-time transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5D ea723abb0acd6a8bffdb9ff8470df3fc8c493e3b.

## Question

Does the frozen timing-aware ordered-path coordinate transfer to restoration-time pairs that introduce an earlier restore round never used in LM5A-LM5D?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Preregistered heldout restore-time pairs

All six ordered pairs below introduce restore round 2:
- budget restore 2; throughput restore 3
- budget restore 3; throughput restore 2
- budget restore 2; throughput restore 4
- budget restore 4; throughput restore 2
- budget restore 2; throughput restore 5
- budget restore 5; throughput restore 2

Concrete windows:
- budget 0/2 ; throughput 1/3
- budget 1/3 ; throughput 0/2
- budget 0/2 ; throughput 2/4
- budget 2/4 ; throughput 0/2
- budget 0/2 ; throughput 3/5
- budget 3/5 ; throughput 0/2

These restoration pairs were not present in the pooled LM5A-LM5D schedule sets.

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

64 resource configurations per schedule; six schedules pooled into 384 points per condition cell.

## Frozen coordinates

Comparator — ordered path:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before earlier restoration
- cumulative actions after earlier restoration
- actions before later restoration

Primary — ordered path + restore rounds:
- comparator fields
- earlier restore round
- later restore round

No cut-round labels, resource identity, or full schedule identity.

## Measurements

Per profile × budget reduction × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

72 pooled cells; 144 summaries.

## Interpretation

Zero primary spread would support transfer of the timing-aware path representation to unseen earlier restore-time pairs. Residual spread would identify restoration-time extrapolation as a boundary of the current resource-state law.

## Bounds

Diagnostic only. No adaptive schedule selection, cut-round labels, resource identity, full schedule identity, coordinate modification, resource tuning, topology/profile changes, capacity change after results, extra training, live activation, or production authority.
