# Wingless UP-LM5F — restore-time × resource-identity interaction

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5E a41f2bb21374392e21fdc2c1aa37f0346307a243.

## Question

Does the LM5E round-2 restoration ambiguity disappear when the timing-aware ordered-path representation is told only which resource owns the earlier restoration landmark?

## Frozen environment

Exactly the LM5E heldout restoration-time regime:
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

Comparator — timing-aware ordered path:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before earlier restoration
- cumulative actions after earlier restoration
- actions before later restoration
- earlier restore round
- later restore round

Primary — comparator + early-resource identity:
- all comparator fields
- one categorical bit:
  - budget_first
  - throughput_first

No cut-round label or full schedule identity.

## Measurements

Per pooled condition cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

If the primary becomes exact, the LM5E failure is specifically an interaction between absolute restoration time and which resource recovered at that time. Persistent spread means resource identity remains insufficient and a deeper pre-action resource-state variable is missing.

## Bounds

Diagnostic only. No adaptive identity encoding, cut-round labels, full schedule identity, coordinate search, resource tuning, topology/profile changes, capacity changes after results, extra training, live activation, or production authority.
