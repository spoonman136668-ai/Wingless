# Wingless UP-LM3C — distinct-arm throughput

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3B eb17353a337835e67ff68a28ed5d4c3af7f7e9a3.

## Question

UP-LM3B synchronized deadlines but throughput 2/3 still did not improve outcomes. Is extra throughput being neutralized because multiple same-round actions collapse into the same arm?

## Frozen scenarios

Reuse UP-LM3B unchanged:
- exact recall cap 16;
- six arms;
- identity rotations 3 and 11;
- permutations identity and reverse;
- deadline profiles by_deferred_level and by_layout;
- total action budget 6;
- global pressure rounds 6;
- throughputs 1, 2, 3.

## Policy change under test

For earliest_deadline and fixed_order only, at most one action may be spent on a given arm in the same round.

Within that constraint:
- earliest_deadline selects the most urgent still-unused arm;
- fixed_order selects the first still-unused arm with a pending dependency.

The constraint is frozen before the run and identical across throughput levels.

## Measurements

Per profile × throughput:
- baseline failures;
- earliest-deadline failures;
- fixed-order failures;
- actions used;
- failures prevented;
- earliest-deadline advantage.

## Interpretation

If throughput now reduces failures, the LM3B null was a distribution/scheduler bottleneck rather than evidence that throughput is intrinsically irrelevant. If outcomes remain unchanged, total budget or deeper dependency structure is still dominant.

## Bounds

Counterfactual only. No live activation, adaptive budget, adaptive throughput, capacity change, extra training, semantic priority, or production authority.
