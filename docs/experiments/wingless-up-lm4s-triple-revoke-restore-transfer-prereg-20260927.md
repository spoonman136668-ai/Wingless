# Wingless UP-LM4S — triple revoke/restore transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4R e3df4b476ef273a679da519c0d7632ccda09e4d3.

## Question

After three authority-loss/restoration events, does final resource state plus cumulative deliverability before the latest restoration still summarize the path, or are earlier restoration checkpoints required?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen triple-window schedules
A:
- cut0/restore1
- cut2/restore3
- cut4/restore5

B:
- cut0/restore2
- cut2/restore3
- cut4/restore5

C:
- cut0/restore1
- cut2/restore4
- cut4/restore5

Revocation amount is fixed within an arm:
- 1 action
- 2 actions

At each restoration the original budget returns. Actions already spent are retained.

## Frozen resource grid
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 1,2,3,4
- six rounds
- earliest_deadline scheduling

## Frozen topology
- rotations 5,13
- permutations identity, reverse, rotate2

## Frozen coordinates
1. final: actual final reachable actions + throughput + actual final coverage end
2. final+latestrestore: final variables + cumulative deliverability strictly before the third restoration
3. final+allrestores: final variables + cumulative deliverability strictly before each of the three restorations

## Measurements
Per profile × schedule × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- max/mean failure spread

108 condition cells; 324 summaries.

## Bounds
Diagnostic only. No adaptive checkpoint search, schedule/resource tuning, coordinate changes, extra training, live activation, or production authority.
