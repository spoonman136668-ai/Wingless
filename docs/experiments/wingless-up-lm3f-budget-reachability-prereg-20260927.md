# Wingless UP-LM3F — cross-budget reachability replication

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3E 5f3267b1b9451d4a785d232fcc05cfb4dc90a10f.

## Question

Does the LM3E reachable-budget law generalize when the total corrective budget itself changes?

## Frozen scenario

Reuse LM3E:
- exact recall cap 16;
- deadline profiles by_deferred_level and by_layout;
- rotations 3 and 11;
- permutations identity and reverse;
- one action per arm per round;
- six global rounds;
- earliest_deadline and fixed_order policies.

## Frozen grid

Total budgets:
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

Reachable budget is min(total budget, throughput × remaining rounds).

## Measurements

Per profile × budget × onset × throughput:
- reachable budget;
- actual actions used;
- earliest and fixed-order failures;
- failures prevented versus baseline.

## Interpretation

If failure count depends primarily on reachable budget rather than nominal budget, start time, or throughput separately, the resource law is general across the tested budgets.

## Bounds

Counterfactual only. No adaptive resource settings, future schedule oracle, capacity change, semantic priority, extra training, live activation, or production authority.
