# Wingless UP-LM5C — ordered restore-timing labels

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5B 14fa0ec1abcf1f5afa63057c714cff9147cfb10e.

## Question

Is the pooled asynchronous-schedule ambiguity explained by the absolute rounds at which the early and late restorations occur?

## Frozen environment

Exactly UP-LM5B:
- deadline profiles deferred_only, layout_only, hybrid_min
- six asynchronous budget/throughput cut-restore windows
- budget reductions 1,2
- throughput reductions 1,2
- rotations 5,13
- permutations identity, reverse, rotate2
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

Pooling:
- all six asynchronous schedules pooled together per profile × budget reduction × throughput reduction × rotation × permutation
- 384 points per pooled cell
- 72 pooled cells

## Frozen coordinates

Comparator — ordered path:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before earlier restoration
- cumulative actions immediately after earlier restoration
- actions before later restoration

Primary — ordered path + restore timing:
- comparator fields
- earlier restoration round
- later restoration round

The primary key does NOT include:
- which resource restores first
- budget restore identity
- throughput restore identity
- full schedule identity
- cut-round labels

## Measurements

Per pooled cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

Zero primary spread would show that absolute temporal placement of the restoration landmarks—not resource identity—is the missing state. Residual spread would require examining cut timing or another path feature without post-hoc coordinate search.

## Bounds

Diagnostic only. No adaptive timing labels, resource identity in primary key, full schedule identity, coordinate search, resource tuning, topology/profile changes, capacity change after results, extra training, live activation, or production authority.
