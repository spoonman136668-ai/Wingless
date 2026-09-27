# Wingless UP-LM4X — overlap transition checkpoint

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4W 417dbcd932142c34f8e3a9f4f17ca67d25807636.

## Question

Is the residual asynchronous restore ambiguity explained by how much action capacity is delivered immediately after throughput restoration, while the strong budget cut is active?

## Frozen residual boundary

Deadline profiles:
- layout_only
- hybrid_min

Budget window:
- cut round 3
- restore round 5
- reduction 2

Throughput window:
- cut round 1
- restore round 3
- reductions 1 and 2

Topology:
- rotations 5 and 13
- permutations identity, reverse, rotate2

Resource grid:
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

64 resource configurations per cell.

## Frozen coordinates

Comparator:
- final + both restore checkpoints
  - final reachable actions
  - nominal throughput
  - final coverage end round
  - actions before budget restore
  - actions before throughput restore

Primary:
- comparator + cumulative actions after round 3
  - equivalently actions available before round 4

No other variable is added.

## Measurements

Per profile × throughput reduction × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

24 topology-condition cells; 48 coordinate summaries.

## Interpretation

Zero spread after adding the round-3 transition checkpoint would localize the missing state to within-window delivery timing. Persistent spread would reject this targeted explanation.

## Bounds

Diagnostic only. No adaptive checkpoint search, resource tuning, coordinate modification, topology/profile expansion, capacity change after results, extra training, live activation, or production authority.
