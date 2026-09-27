# Wingless UP-LM4P — revoke/restore path-checkpoint ablation

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4O 65dcdb44921783624ac69875b999a8fb420bd671.

## Question

Which causal checkpoint is necessary to resolve revoke/restore path dependence: deliverability before the cut, deliverability before restoration, or both?

## Frozen scenario

Identical to UP-LM4O:
- deadline profiles: deferred_only, layout_only, hybrid_min
- schedules: cut2/restore4, cut2/restore5, cut3/restore5
- revoke amounts: 1,2
- budgets: 4,5,6,7
- starts: 2,3,4,5
- throughput: 1,2,3,4
- six global rounds
- rotations 5 and 13
- permutations identity, reverse, rotate2
- earliest_deadline scheduling

## Frozen coordinates

A. final-only comparator:
- actual reachable actions
- throughput
- actual coverage end

B. pre-cut:
- comparator
- cumulative deliverability strictly before cut

C. pre-restore:
- comparator
- cumulative deliverability strictly before restore

D. both:
- comparator
- cumulative deliverability strictly before cut
- cumulative deliverability strictly before restore

No other checkpoint is permitted.

## Measurements

Per deadline profile × schedule × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- max spread
- mean spread

108 topology-condition cells; 432 coordinate summaries.

## Interpretation

A single-checkpoint coordinate with zero spread is the minimal causal state representation. If neither single checkpoint is exact but the joint coordinate is, both transition-boundary states are jointly necessary.

## Bounds

Diagnostic only. No adaptive checkpoint search, schedule/revocation/resource tuning, coordinate modification, topology/profile selection, extra training, live activation, or production authority.
