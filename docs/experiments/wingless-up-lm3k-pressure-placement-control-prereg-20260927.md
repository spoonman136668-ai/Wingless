# Wingless UP-LM3K — pressure-placement control

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3J 2a8db4e75c0bc5c49271eb7760cbb8cb085c78c0.

## Question

Did the reachable-budget law break because pressure was dynamic, or simply because one extra write of pressure arrived before corrective deployment was effectively over?

## Frozen initial condition

Reuse hybrid_min prepressure, held-out rotations, permutations, and resource grid unchanged.

## Frozen pressure-placement arms

Each arm receives one of four preregistered pressure modes:
- none: no extra write
- static_pre: one extra ordinary write per arm immediately after prepressure and before global round 0
- dynamic_r3: one extra ordinary write per arm after the standard write in global round 3
- dynamic_r5: one extra ordinary write per arm after the standard write in global round 5

Extra-write magnitude is exactly one in every non-none arm.

## Frozen resource grid

Budgets: 4, 5, 6.
Action start rounds: 3, 4.
Throughput: 1, 2, 3.

Reachable budget remains min(total budget, throughput × remaining action rounds).

## Measurements

Per pressure mode × budget × onset × throughput:
- reachable budget
- baseline failures
- earliest-deadline failures
- fixed-order failures
- actions used
- failures prevented

For each pressure mode × reachable budget:
- minimum and maximum earliest-deadline failures
- failure spread

## Interpretation

If static_pre reproduces dynamic_r3, total pre-effective-window pressure is sufficient to explain the break. If only dynamic_r3 breaks the collapse, nonstationary timing is causal. dynamic_r5 tests pressure arriving after useful corrective deployment.

## Bounds

Counterfactual only. No adaptive pressure placement or magnitude, no resource tuning, no topology search, no future schedule oracle, no capacity change, no semantic priority, no extra training, no live activation, or production authority.
