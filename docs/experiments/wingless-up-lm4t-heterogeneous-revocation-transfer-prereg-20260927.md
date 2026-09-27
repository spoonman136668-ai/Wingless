# Wingless UP-LM4T — heterogeneous revocation-depth transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4S e2012d7693df8699972d89392e4433e098916cca.

## Question

Does final resource state plus cumulative deliverability before the latest restoration still summarize the path when successive authority-loss events have different severities?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen triple-window schedules
A: (cut0/restore1), (cut2/restore3), (cut4/restore5)
B: (cut0/restore2), (cut2/restore3), (cut4/restore5)
C: (cut0/restore1), (cut2/restore4), (cut4/restore5)

## Frozen per-window revocation patterns
- 1,2,1 actions
- 2,1,2 actions
- 1,1,2 actions
- 2,2,1 actions

At each restoration the original total budget returns. Actions already spent remain spent.

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
3. final+allrestores: final variables + cumulative deliverability strictly before all three restorations

216 condition cells; 648 summaries.

## Interpretation

Zero spread for latest-restoration state would show that earlier heterogeneous authority-loss history is still compressible. A break repaired only by all-restoration state would identify the boundary where older event memory becomes necessary.

## Bounds

Diagnostic only. No adaptive schedule/severity/resource selection, checkpoint search, coordinate modification, extra training, live activation, or production authority.
