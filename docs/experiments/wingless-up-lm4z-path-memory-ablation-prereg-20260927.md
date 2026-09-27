# Wingless UP-LM4Z — asynchronous path-memory ablation

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4Y 0724b89e0630e14e832d56d3754736449610e2cd.

## Question

What is the smallest temporal resource-memory coordinate that preserves exact outcome collapse across the accepted asynchronous restore grid?

## Frozen environment

Exactly UP-LM4Y:
- deadline profiles deferred_only, layout_only, hybrid_min
- six asynchronous budget/throughput cut-restore windows
- budget reductions 1,2
- throughput reductions 1,2
- rotations 5,13
- permutations identity, reverse, rotate2
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

432 topology-condition cells; 64 resource configurations per cell.

## Frozen coordinates

A. final+post-early
- final reachable actions
- nominal throughput
- final coverage end round
- cumulative actions immediately after the earlier restoration round

B. final+post-early+pre-late
- coordinate A
- cumulative actions immediately before the later restoration round

C. full accepted coordinate
- final reachable actions
- nominal throughput
- final coverage end round
- actions before budget restoration
- actions before throughput restoration
- cumulative actions immediately after the earlier restoration round

No coordinate changes after results.

## Measurements

Per cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

1296 coordinate summaries.

## Interpretation

If coordinate A is exact, one transition checkpoint is sufficient. If A fails but B is exact, two ordered transition checkpoints are sufficient. If only C is exact, resource-specific restoration identity remains necessary.

## Bounds

Diagnostic only. No adaptive coordinate search, resource tuning, schedule changes, topology/profile changes, capacity change after results, extra training, live activation, or production authority.
