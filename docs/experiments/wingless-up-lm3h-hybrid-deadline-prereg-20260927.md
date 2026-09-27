# Wingless UP-LM3H — hybrid-deadline reachability transfer

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3G afe5ac7f16beb36eac779ca77654c1cba5b8c911.

## Question

Does the reachable-budget law survive a genuinely new deadline shape that combines the two previously tested urgency structures?

## Frozen hybrid profile

For each arm, compute the two already-defined prepressure targets:
- deferred-level target: d=4 -> 1, d=5 -> 2, otherwise -> 3;
- layout target: suffix_reported -> 1, otherwise -> 2.

The new hybrid_min target is the smaller (more urgent) of those two targets.

This rule is fixed before the run and uses no observed outcome.

## Frozen topology

Held-out topology from UP-LM3G:
- rotations 5 and 13;
- permutations identity, reverse, rotate2.

## Frozen resource grid

Budgets:
- 4;
- 5;
- 6.

Action start rounds:
- 3;
- 4.

Throughput:
- 1;
- 2;
- 3.

Reachable budget remains:
min(total budget, throughput × remaining action rounds).

## Measurements

Per budget × onset × throughput:
- reachable budget;
- baseline failures;
- earliest-deadline failures;
- fixed-order failures;
- actions used;
- failures prevented.

Group earliest-deadline outcomes by reachable budget and report within-reachable-budget failure spread.

## Interpretation

Zero spread again would show that reachable budget survives a new composite deadline geometry. Nonzero spread would identify interaction between resource reachability and deadline shape.

## Bounds

Counterfactual only. No adaptive deadline construction, resource tuning, topology search, future schedule oracle, capacity change, semantic priority, extra training, live activation, or production authority.
