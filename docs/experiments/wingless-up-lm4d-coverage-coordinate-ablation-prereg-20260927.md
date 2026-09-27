# Wingless UP-LM4D — coverage-coordinate single-variable ablation

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4C a74d3337d57aa7da9fc36dc8e6ce7958b90c68c8.

## Question

Are all four variables in the exact interpretable coverage coordinate necessary?

## Frozen matrix

Exactly the UP-LM4A/4C topology-resolved matrix:
- deadline profiles: deferred_only, layout_only, hybrid_min
- pressures: 0 through 8
- rotations: 5 and 13
- permutations: identity, reverse, rotate2
- budgets: 3 through 7
- action starts: 2 through 5
- throughputs: 1 through 4
- six rounds

## Frozen exact coordinate

- reachable_budget
- capacity_by_round2
- throughput
- coverage_end_round

## Frozen single-variable ablations

A. drop reachable:
- cap2 + throughput + coverage_end

B. drop cap2:
- reachable + throughput + coverage_end

C. drop throughput:
- reachable + cap2 + coverage_end

D. drop coverage_end:
- reachable + cap2 + throughput

E. control:
- reachable + cap2 + throughput + coverage_end

## Measurements

For each of 162 topology-condition cells × coordinate:
- groups
- nonzero-spread groups
- maximum spread
- mean spread

## Interpretation

If only E remains exact, all four variables are independently necessary in the tested resource regime.

## Bounds

Diagnostic only. No adaptive coordinate selection, resource tuning, topology selection, capacity change, extra training, live activation, or production authority.
