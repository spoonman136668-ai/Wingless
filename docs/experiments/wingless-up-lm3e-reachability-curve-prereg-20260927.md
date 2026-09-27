# Wingless UP-LM3E — budget reachability curve

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3D 7b6759c78d0faceb7ba59f5f3adddf9f3b24186d.

## Question

As the intervention window opens progressively later, what minimum per-round throughput is required for Wingless to make its fixed six-action budget reachable before the pressure window closes?

## Frozen scenario

Reuse UP-LM3D:
- exact recall cap 16;
- six arms;
- deadline profiles by_deferred_level and by_layout;
- rotations 3 and 11;
- permutations identity and reverse;
- one action per arm per round;
- total budget 6;
- six global pressure rounds;
- baseline, earliest_deadline, fixed_order.

## Frozen reachability grid

Action start rounds:
- 0;
- 1;
- 2;
- 3;
- 4.

Throughput:
- 1;
- 2;
- 3.

For each cell, theoretical maximum deployable actions is min(6, throughput × remaining rounds).

## Measurements

Per profile × onset × throughput:
- theoretical reachable action budget;
- actual actions used;
- failures under each policy;
- failures prevented;
- earliest-deadline advantage.

## Interpretation

If earliest-deadline outcomes track reachable action budget, this identifies a compact resource law linking timing, throughput, and total corrective budget.

## Bounds

Counterfactual only. No adaptive onset/budget/throughput, future schedule oracle, capacity change, semantic priority, extra training, live activation, or production authority.
