# Wingless UP-LM3J — burst-timing reachability sweep

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3I 4fdaccbd98457f3cf9543ad50b65710ad8335e4b.

## Question

At which hazard-transition times does the simple reachable-budget law stop being sufficient?

## Frozen initial condition

Reuse the UP-LM3I hybrid_min prepressure, held-out rotations, and arm permutations unchanged.

## Frozen dynamic perturbation

Inject exactly one additional ordinary write into every arm at one globally fixed burst round.

Burst rounds are swept exhaustively:
- 0
- 1
- 2
- 3
- 4
- 5

Each arm is evaluated independently for each preregistered burst round. Burst magnitude is always one extra write per arm.

## Frozen resource grid

Budgets:
- 4
- 5
- 6

Action start rounds:
- 3
- 4

Throughput:
- 1
- 2
- 3

Reachable budget remains:
min(total budget, throughput × remaining action rounds).

## Measurements

Per burst round × budget × onset × throughput:
- reachable budget
- baseline failures
- earliest-deadline failures
- fixed-order failures
- actions used
- failures prevented

For each burst round × reachable budget:
- minimum earliest-deadline failures
- maximum earliest-deadline failures
- failure spread

## Interpretation

Burst rounds with zero spread preserve the one-number resource law. Nonzero spread identifies hazard timing regimes where deployment timing and reachable capacity interact.

## Bounds

Counterfactual only. No adaptive burst timing or magnitude, no resource tuning, no topology search, no future schedule oracle, no capacity change, no semantic priority, no extra training, no live activation, or production authority.
