# Wingless UP-LM3L — temporal resource coordinate

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3K d077e4f8ee66b0d5a0b9c55f7b3108ec0290569d.

## Question

Under the heavier static pressure where reachable budget alone no longer determines outcome, what is the smallest preregistered temporal resource coordinate that restores deterministic collapse?

## Frozen scenario

Reuse UP-LM3K unchanged with:
- pressure mode = static_pre
- hybrid_min deadline profile
- held-out rotations 5 and 13
- permutations identity, reverse, rotate2
- six global rounds
- one action per arm per round

Resource grid:
- budgets 4, 5, 6
- action start rounds 3, 4
- throughput 1, 2, 3

## Preregistered candidate coordinates

Group earliest-deadline failure outcomes by:
1. reachable_budget
2. reachable_budget + action_start_round
3. reachable_budget + throughput
4. reachable_budget + action_start_round + throughput

No coordinate is added after results.

## Measurements

For every resource cell:
- reachable budget
- onset
- throughput
- earliest-deadline failure count

For each candidate coordinate:
- number of groups
- groups with nonzero failure spread
- maximum within-group failure spread
- mean within-group failure spread

## Interpretation

The lowest-dimensional coordinate with zero within-group spread is a candidate compact temporal resource state under heavier pressure. If only the full coordinate resolves outcomes, no simpler collapse is supported.

## Bounds

Diagnostic only. No adaptive grouping, resource tuning, new intervention, topology search, future schedule oracle, capacity change, semantic priority, extra training, live activation, or production authority.
