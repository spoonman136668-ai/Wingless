# Wingless UP-LM5A — ordered path pooling across asynchronous schedules

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4Z c38dc038c49a397906a4cd6482783760dd010d9c.

## Question

Can asynchronous resource outcomes be predicted from ordered path landmarks alone, without encoding which resource restored first?

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

## Pooling rule

For each profile × budget reduction × throughput reduction × rotation × permutation:
- pool all six asynchronous schedules together;
- 384 points per pooled cell;
- schedule identity is not included in any coordinate key.

72 pooled cells total.

## Frozen coordinates

Resource-specific full:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before budget restoration
- actions before throughput restoration
- cumulative actions after the earlier restoration

Ordered path:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before the earlier restoration
- cumulative actions after the earlier restoration
- actions before the later restoration

Neither key includes budget/throughput restore identity or restore-round labels.

## Measurements

Per pooled cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

Zero ordered-path spread would show that the temporal path can be represented independently of resource identity. Residual spread would show that path counts alone are insufficient across schedule identities.

## Bounds

Diagnostic only. No adaptive pooling, coordinate search, schedule labels in keys, resource tuning, topology/profile changes, capacity change after results, extra training, live activation, or production authority.
