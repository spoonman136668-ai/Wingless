# Wingless UP-LM4C — coverage-horizon resource coordinate

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4B 37300f5cf16ac729d5ef9ff135d09116c36a6f66.

## Question

Can the exact temporal resource state be expressed more interpretabily as early deployability plus action-delivery duration and rate?

## Frozen matrix

Exactly the UP-LM4A/4B topology-resolved matrix:
- three deadline profiles
- pressures 0 through 8
- rotations 5 and 13
- permutations identity, reverse, rotate2
- budgets 3 through 7
- start rounds 2 through 5
- throughputs 1 through 4
- six rounds

## Frozen derived state

coverage_end_round:
- the last scheduler round on which the reachable action budget can be delivered at the configured throughput.

## Frozen coordinates

A:
- reachable_budget + capacity_by_round2 + coverage_end_round

B:
- reachable_budget + capacity_by_round2 + throughput + coverage_end_round

Control:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4

## Measurements

For each of 162 topology-condition cells × coordinate:
- groups
- nonzero-spread groups
- maximum spread
- mean spread

## Interpretation

If B is exact, the checkpoint vector can be replaced by a more interpretable total/early-capacity/rate/duration state. If B fails while the control remains exact, intermediate deployability checkpoints contain information not captured by simple coverage duration and rate.

## Bounds

Diagnostic only. No adaptive coordinate selection, resource tuning, topology selection, capacity change, extra training, live activation, or production authority.
