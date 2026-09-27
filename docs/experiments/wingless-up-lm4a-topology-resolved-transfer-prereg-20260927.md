# Wingless UP-LM4A — topology-resolved resource transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3Z 10578d7c383bd64cfd75afa9e873db3d6e9f6126.

## Question

Is the exact minimal resource-state collapse genuine within each individual topology, or an artifact of summing outcomes across rotations and permutations?

## Frozen conditions

Deadline profiles:
- deferred_only
- layout_only
- hybrid_min

Pressure levels:
- 0 through 8

Resource grid:
- budgets 3,4,5,6,7
- action start rounds 2,3,4,5
- throughput 1,2,3,4

Topologies evaluated separately:
- rotations 5 and 13
- permutations identity, reverse, rotate2

Scheduling:
- six global rounds
- earliest_deadline policy

## Frozen coordinates

Prior:
- reachable_budget + capacity_by_round3 + capacity_by_round4

Minimal:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4

Full:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4 + capacity_by_round5

## Measurements

For each deadline profile × pressure × rotation × permutation × coordinate:
- groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread
- global outcome range
- distinct outcome count

There are 162 topology-condition cells and 486 coordinate summaries.

## Interpretation

Zero spread for the minimal coordinate within every individual topology would rule out aggregate cancellation as the explanation for prior exact collapse.

## Bounds

Diagnostic only. No adaptive topology/profile/pressure/coordinate selection, resource tuning, capacity change, extra training, live activation, or production authority.
