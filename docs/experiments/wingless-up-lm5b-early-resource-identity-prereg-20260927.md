# Wingless UP-LM5B — early-restored resource identity

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5A 610210b0fe8c4ac1e6af2ef1f2e160376c629107.

## Question

Is the residual ambiguity across pooled asynchronous schedules explained by which resource restores first?

## Frozen environment

Deadline profiles:
- deferred_only
- layout_only
- hybrid_min

Asynchronous windows:
- budget 1/3 ; throughput 2/4
- budget 2/4 ; throughput 1/3
- budget 1/4 ; throughput 2/5
- budget 2/5 ; throughput 1/4
- budget 1/3 ; throughput 3/5
- budget 3/5 ; throughput 1/3

Reductions:
- budget 1,2
- throughput 1,2

Topology:
- rotations 5,13
- permutations identity, reverse, rotate2

Resource grid:
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

Pooling:
- pool all six asynchronous schedules together per profile × budget reduction × throughput reduction × rotation × permutation
- 384 points per pooled cell
- 72 pooled cells total

## Frozen coordinates

Comparator — ordered path:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before earlier restoration
- cumulative actions after earlier restoration
- actions before later restoration

Primary — ordered path + early resource identity:
- comparator fields
- one categorical bit:
  - budget_first
  - throughput_first

No restore-round label or full schedule identity is included.

## Measurements

Per pooled cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

Zero spread after adding only early-resource identity would show that schedule identity matters only through which resource recovers first. Residual spread would mean restore timing itself carries additional information beyond ordered path counts and resource ordering.

## Bounds

Diagnostic only. No adaptive identity encoding, restore-round labels, full schedule identity, coordinate search, resource tuning, topology/profile changes, capacity change after results, extra training, live activation, or production authority.
