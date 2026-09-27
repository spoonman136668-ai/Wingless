# Wingless UP-LM5M — allocation-policy transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5L 52a30bf88ae98e36fd668e7203968fcd93c31e41.

## Question

Does the frozen enriched deployability state remain sufficient when the same resource path is allocated across arms by a different scheduler policy?

## Frozen scheduler policies
- earliest_deadline
- fixed_order

Both are existing accepted harness policies.

## Frozen environment
- six LM5K restore-round-1 schedules
- profiles: deferred_only, layout_only, hybrid_min
- budget reductions 1,2
- throughput reductions 1,2
- rotations 5,13
- permutations identity, reverse, rotate2
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds.

## Frozen resource representation

Enriched coordinate:
- final reachable actions
- nominal throughput
- final coverage end round
- pre/post early restoration actions
- pre/post late restoration actions
- early/late restore rounds
- early-restored resource identity
- first_deliverable_round.

This coordinate does not include scheduler policy.

## Compared keys

1. enriched:
- pool earliest_deadline and fixed_order together under the frozen enriched coordinate.

2. enriched+policy:
- same key plus scheduler policy identity.

384 points per policy and 768 pooled points per condition cell; 72 pooled cells; 144 summaries.

## Interpretation

If enriched remains exact while policies are pooled, the current path state is policy-invariant. If enriched breaks but enriched+policy is exact, action allocation order is causal state not represented by aggregate deployability. If both break, deeper per-arm allocation state is required.

## Bounds

Diagnostic only. No adaptive policy selection, no per-arm delivery features, no coordinate search, no resource tuning, no topology/profile changes, no live activation or production authority.
