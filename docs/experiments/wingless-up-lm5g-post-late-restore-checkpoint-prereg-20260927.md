# Wingless UP-LM5G — post-late-restore checkpoint

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5F 1ccefa94b81ce0ed4e0eab711ee066e5f985721a.

## Question

Is the remaining round-2 restoration ambiguity explained by delivery immediately after the later restoration transition?

## Frozen environment

Exactly the LM5E/LM5F heldout regime:
- deadline profiles: deferred_only, layout_only, hybrid_min
- windows:
  - budget 0/2 ; throughput 1/3
  - budget 1/3 ; throughput 0/2
  - budget 0/2 ; throughput 2/4
  - budget 2/4 ; throughput 0/2
  - budget 0/2 ; throughput 3/5
  - budget 3/5 ; throughput 0/2
- budget reductions: 1,2
- throughput reductions: 1,2
- rotations: 5,13
- permutations: identity, reverse, rotate2
- budgets: 4,5,6,7
- action starts: 2,3,4,5
- throughput: 2,3,4,5
- six global rounds
- earliest_deadline scheduling

Six schedules are pooled into 384 points per condition cell.
72 pooled condition cells.

## Frozen coordinates

Comparator — timing-aware + early-resource identity:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before earlier restoration
- cumulative actions after earlier restoration
- actions before later restoration
- earlier restore round
- later restore round
- early-restored resource identity

Primary — comparator + post-late checkpoint:
- all comparator fields
- cumulative actions immediately after the later restoration round

No cut-round labels, raw action-start field, raw budget field, or full schedule identity.

## Measurements

Per pooled condition cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

Zero primary spread would identify post-transition delivery at the later restore as the missing path-memory state. Residual spread would reject this targeted transition-state explanation and require a deeper within-segment diagnostic.

## Bounds

Diagnostic only. No adaptive checkpoint search, cut-round labels, raw scheduler parameters, full schedule identity, coordinate search, resource tuning, topology/profile changes, capacity changes after results, extra training, live activation, or production authority.
