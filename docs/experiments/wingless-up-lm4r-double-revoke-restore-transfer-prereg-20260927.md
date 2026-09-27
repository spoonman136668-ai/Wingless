# Wingless UP-LM4R — double revoke/restore path transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4Q 8b3dfbc3253add2c2688f68dbaa9bc36ed945bfd.

## Question

When authority is revoked and restored twice, is one cumulative-deployability checkpoint per restoration event sufficient to recover deterministic resource-state prediction?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen two-window schedules

A:
- cut 1, restore 3
- cut 4, restore 5

B:
- cut 2, restore 3
- cut 4, restore 5

C:
- cut 1, restore 2
- cut 3, restore 5

During either revoked interval the total action cap is reduced by the same fixed amount:
- 1 action
- 2 actions

At each restoration the original total cap returns. Actions already spent are never undone.

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

A. final-only:
- actual final reachable actions
- throughput
- actual final coverage end

B. last-restore:
- final-only variables
- cumulative deliverability strictly before second restoration

C. both-restores:
- final-only variables
- cumulative deliverability strictly before first restoration
- cumulative deliverability strictly before second restoration

No cut checkpoints are included.

## Measurements

Per deadline profile × double-window schedule × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- max failure spread
- mean failure spread

108 topology-condition cells; 324 coordinate summaries.

## Interpretation

Zero spread for both-restores would support a bounded event-indexed resource memory: one checkpoint per restoration event. If last-restore alone is exact, earlier restoration history is compressible. Residual spread even with both checkpoints would show that finer path structure is required.

## Bounds

Diagnostic only. No adaptive window/revocation/resource selection, checkpoint search, coordinate modification, topology/profile tuning, extra training, live activation, or production authority.
