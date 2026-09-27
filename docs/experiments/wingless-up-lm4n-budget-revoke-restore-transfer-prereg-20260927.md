# Wingless UP-LM4N — budget revoke-then-restore transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4M 45ffd2f3cec82ef3f034020798434845610ad1c1.

## Question

After temporary resource loss followed by restoration, are actual reachable actions and actual coverage end sufficient, or does the path of budget availability itself affect outcomes?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen budget schedules

Initial budgets:
- 4,5,6,7 actions

Temporary revocation windows:
- cut round 2, restore round 4
- cut round 2, restore round 5
- cut round 3, restore round 5

During the revoked interval:
- total action cap is reduced by 1 or 2 actions.

At restore round and thereafter:
- original total action cap returns.

Actions already spent are never undone.

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

Static comparator:
- original reachable_budget
- nominal throughput
- original coverage_end_round

Dynamic primary:
- actual reachable actions under the revoke/restore cap schedule
- nominal throughput
- actual coverage_end_round under that schedule

The cut, restore, and revoke amount are fixed within each condition cell and are not grouping-key inputs.

## Measurements

Per deadline profile × schedule × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- max failure spread
- mean failure spread

108 topology-condition cells; 216 coordinate summaries.

## Interpretation

Zero dynamic spread would show that total actual deployability and coverage summarize this non-monotonic resource schedule. Dynamic breaks would demonstrate that two schedules with equal reachable amount and coverage can still differ because the timing path of temporary authority matters.

## Bounds

Diagnostic only. No adaptive cut/restore timing or amount, resource tuning, coordinate modification, topology/profile selection, extra training, live activation, or production authority.
