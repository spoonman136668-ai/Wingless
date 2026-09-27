# Wingless UP-LM3E — temporal reachability envelope

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3D 7b6759c78d0faceb7ba59f5f3adddf9f3b24186d.

## Question

Can the throughput effect found under compressed intervention windows be characterized as a simple temporal reachability boundary under fixed total resources?

## Frozen scenario

Reuse UP-LM3D unchanged:
- exact recall cap 16;
- six arms;
- deadline profiles by_deferred_level and by_layout;
- identity rotations 3 and 11;
- permutations identity and reverse;
- one action per arm per round;
- total action budget 6;
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

For each cell, theoretical reachable actions per scenario are min(6, throughput × remaining action rounds).

## Measurements

Per profile × start round × throughput:
- theoretical reachable actions per scenario;
- actual actions;
- completed and failed originals;
- failures prevented versus baseline;
- earliest-deadline advantage over fixed-order.

## Interpretation

If failure outcomes organize primarily by reachable action capacity, throughput is a temporal deployment constraint. Deviations at equal reachable capacity identify additional deadline-order or topology effects.

## Bounds

Counterfactual only. No adaptive onset, budget, throughput, or policy; no capacity change; no future schedule oracle; no semantic priority; no extra training; no live activation; no production authority.
