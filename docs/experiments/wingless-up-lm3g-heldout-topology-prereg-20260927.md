# Wingless UP-LM3G — reachable-budget held-out topology

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3F 062aa92cfa59ab09ce77b336ca788d96e76ed0dc.

## Question

Does the LM3F reachable-budget law transfer to scheduler geometries not used in its derivation?

## Frozen resource law

Reachable budget remains:
min(total budget, throughput × remaining action rounds).

No law parameters are fitted in this experiment.

## Held-out topology

Identity rotations:
- 5;
- 13.

Arm permutations:
- identity;
- reverse;
- rotate2.

The rotate2 permutation was not part of LM3F.

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

Deadline profiles remain:
- by_deferred_level;
- by_layout.

## Measurements

Per profile × budget × onset × throughput:
- reachable budget;
- baseline failures;
- earliest-deadline failures;
- fixed-order failures;
- actions used;
- failures prevented.

Group earliest-deadline outcomes by reachable budget to test whether nominal budget/onset/throughput still collapse to the same effective resource coordinate.

## Interpretation

Topology-invariant collapse supports reachable budget as a general fixed-resource law. Systematic deviations identify interaction with arm ordering or identity rotation.

## Bounds

Counterfactual only. No adaptive resources, topology search, future schedule oracle, capacity change, semantic priority, extra training, live activation, or production authority.
