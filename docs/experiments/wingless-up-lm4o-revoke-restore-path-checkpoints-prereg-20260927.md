# Wingless UP-LM4O — revoke/restore path checkpoints

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4N on branch research/wingless-up-lm4n-budget-revoke-restore-transfer-r1.

## Question

Can two fixed cumulative-deployability checkpoints recover the resource information lost by final reachable amount and coverage under temporary budget revocation?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen budget schedules
- cut round 2, restore round 4
- cut round 2, restore round 5
- cut round 3, restore round 5

Revocation amounts:
- 1 action
- 2 actions

## Frozen resource grid
- initial budgets 4,5,6,7
- action starts 2,3,4,5
- throughput 1,2,3,4
- six global rounds
- earliest_deadline scheduling

64 resource configurations per topology-condition cell.

## Frozen topology
- rotations 5 and 13
- permutations identity, reverse, rotate2

## Frozen coordinates

Comparator:
- actual final reachable actions under revoke/restore schedule
- nominal throughput
- actual final coverage end

Primary path coordinate:
- all comparator variables
- cumulative actions deliverable strictly before cut round
- cumulative actions deliverable strictly before restore round

These two checkpoints are preregistered from the causal boundaries of the resource schedule; no arbitrary round search is allowed.

## Measurements

Per deadline profile × schedule × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- max failure spread
- mean failure spread

108 topology-condition cells; 216 coordinate summaries.

## Interpretation

Zero primary spread would show that the revoke/restore path dependence can be summarized by bounded causal checkpoints at the authority transitions. Residual spread would demonstrate that finer-grained availability history remains necessary.

## Bounds

Diagnostic only. No adaptive checkpoint search, schedule/revocation/resource tuning, coordinate modification, topology/profile selection, extra training, live activation, or production authority.
