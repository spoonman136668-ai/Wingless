# Wingless UP-LM4V — combined budget+throughput revoke/restore transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4U 8200c546c26e4b35b6e963df49f7adf08e907b07.

## Question

Does the frozen pre-restore deliverability law remain sufficient when total budget and per-round throughput are reduced simultaneously and later restored?

## Frozen profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen loss windows
- cut1/restore3
- cut2/restore4
- cut3/restore5

## Frozen combined reductions
- budget -1, throughput -1
- budget -1, throughput -2
- budget -2, throughput -1
- budget -2, throughput -2

## Frozen resource grid
- budgets 4,5,6,7
- starts 2,3,4,5
- nominal throughput 2,3,4,5
- six rounds
- earliest_deadline scheduling

## Frozen topology
- rotations 5,13
- permutations identity, reverse, rotate2

## Frozen coordinates
Comparator:
- actual final reachable actions + nominal restored throughput + actual final coverage end

Primary:
- comparator + cumulative deliverability strictly before restoration

216 condition cells; 432 summaries.

## Bounds
Diagnostic only. No adaptive schedule/reduction/resource selection, checkpoint search, coordinate modification, live activation, or production authority.
