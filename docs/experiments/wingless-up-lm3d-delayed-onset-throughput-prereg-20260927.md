# Wingless UP-LM3D — delayed-onset throughput

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3C bb29b01bf50424f94cc40412fd395059d3b0cc48.

## Question

UP-LM3C showed that distinct-arm distribution makes throughput useful for fixed-order, while earliest-deadline remained limited by the fixed six-action budget. Does throughput become causal for earliest-deadline when the intervention window is compressed so serial spending cannot deploy all six actions in time?

## Frozen scenarios

Reuse UP-LM3C unchanged:
- exact recall cap 16;
- six arms;
- identity rotations 3 and 11;
- permutations identity and reverse;
- deadline profiles by_deferred_level and by_layout;
- total action budget 6;
- global pressure rounds 6;
- throughputs 1, 2, and 3;
- one action per arm per round.

## Frozen delayed onset

No corrective action is allowed during pressure rounds 1, 2, or 3.

Corrective action becomes available only before pressure rounds 4, 5, and 6.

With this fixed window:
- throughput 1 can deploy at most 3 of the 6 available actions;
- throughput 2 can deploy all 6 actions in 3 rounds;
- throughput 3 can deploy all 6 actions in 2 rounds.

The onset is identical across policies, rotations, permutations, profiles, and throughput arms.

## Policies

- baseline;
- earliest_deadline;
- fixed_order.

Within an active round, the UP-LM3C distinct-arm constraint remains frozen:
- earliest_deadline selects the most urgent still-unused arm;
- fixed_order selects the first still-unused arm with a pending dependency;
- at most one action may be spent on any arm in the same round.

## Measurements

Per profile × throughput:
- actions used;
- completed originals;
- failed originals;
- failures prevented versus baseline;
- earliest-deadline versus fixed-order.

## Interpretation

A throughput-dependent reduction in earliest-deadline failures would show that throughput matters once the same finite action budget must be deployed inside a short rescue window.

No improvement would imply that either the delayed window begins after the relevant rescue opportunities have already closed or that total-budget/dependency structure still dominates.

## Bounds

Counterfactual only. No live activation, adaptive onset, adaptive budget, adaptive throughput, capacity change, extra training, semantic priority, or production authority.
