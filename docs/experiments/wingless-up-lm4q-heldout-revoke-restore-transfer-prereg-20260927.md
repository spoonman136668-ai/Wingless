# Wingless UP-LM4Q — heldout revoke/restore schedule transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4P 32880aedcc636afebbb7d305b50ee5d37e717efb.

## Question

Does the minimal revoke/restore resource coordinate transfer unchanged to authority-loss schedules that were not used to derive it?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Heldout revoke/restore schedules

Not used in UP-LM4N–LM4P:
- cut round 1, restore round 3
- cut round 1, restore round 4
- cut round 1, restore round 5
- cut round 2, restore round 3
- cut round 3, restore round 4
- cut round 4, restore round 5

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
- actual final reachable actions
- throughput
- actual final coverage end

Primary minimal coordinate:
- all comparator variables
- cumulative deliverability strictly before restore round

No pre-cut checkpoint and no new state variable.

## Measurements

Per deadline profile × heldout schedule × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- max failure spread
- mean failure spread

216 topology-condition cells; 432 coordinate summaries.

## Interpretation

Zero primary spread supports true transfer of the pre-restore checkpoint law beyond the derivation schedules. Residual spread identifies the schedule boundary of the current compact path state.

## Bounds

Diagnostic only. No adaptive schedule/revocation/resource selection, coordinate modification, topology/profile tuning, extra training, live activation, or production authority.
