# Wingless UP-LM3M — compressed deployment profile

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3L a5baad4ec2d7d89b270a0c888c4fee4e909fe841.

## Question

Can the full reachable-budget + onset + throughput coordinate be compressed into a smaller derived deployment profile that preserves deterministic outcome collapse under heavier pressure?

## Frozen scenario

Reuse UP-LM3L:
- static_pre pressure;
- hybrid_min deadline profile;
- rotations 5 and 13;
- permutations identity, reverse, rotate2;
- budgets 4,5,6;
- action starts 3,4;
- throughput 1,2,3;
- six global rounds.

## Frozen derived quantities

For each resource cell:
- reachable_budget = total actions deployable by the end of round 5.
- capacity_by_round3 = cumulative actions deployable through round 3.
- capacity_by_round4 = cumulative actions deployable through round 4.

Each cumulative capacity is clipped by the total budget.

## Preregistered candidate coordinates

1. reachable_budget
2. reachable_budget + capacity_by_round3
3. reachable_budget + capacity_by_round4
4. reachable_budget + capacity_by_round3 + capacity_by_round4

No coordinate is added after results.

## Measurements

For each candidate:
- group count;
- groups with nonzero earliest-failure spread;
- maximum spread;
- mean spread.

## Interpretation

A zero-spread derived coordinate would compress raw onset and throughput into the action-delivery profile that actually matters. If only the two-checkpoint profile works, both early and mid-window deployability are required.

## Bounds

Diagnostic only. No adaptive coordinate search, resource tuning, new intervention, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
