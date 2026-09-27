# Wingless UP-LM3B — bursty throughput stress

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3A 3798bd40e2ad514f0be6b62d7236b3ecbcaf178b.

## Question

UP-LM3A found little throughput effect because deadlines did not collide strongly. Does per-round corrective throughput become causal when several pending dependencies approach eviction in the same round?

## Frozen scenarios

Reuse six UP-LM2X arms per scenario and exact recall cap 16.

Identity rotations:
- 3;
- 11.

Arm permutations:
- identity;
- reverse.

Before global rounds begin, apply unique pre-pressure writes independently to each arm until its first pending dependency has a frozen target countdown. Stop before eviction.

Two deadline profiles:

1. **by_deferred_level**
   - d=4 arms -> countdown 1;
   - d=5 arms -> countdown 2;
   - d=6 arms -> countdown 3.
   This creates two imminent dependencies at each of the first three deadlines.

2. **by_layout**
   - suffix_reported arms -> countdown 1;
   - alternating_reported arms -> countdown 2.
   This creates three imminent dependencies at each of the first two deadlines.

## Frozen action resources

Total action budget:
- 6.

Per-round throughput:
- 1;
- 2;
- 3.

Global pressure rounds:
- 6.

## Policies

- baseline;
- earliest_deadline;
- fixed_order.

Before each round's pressure write, active policies may perform up to throughput actions while budget remains. Each action marks one pending dependency completed.

## Measurements

Per profile × rotation × permutation × throughput × policy:
- actions used;
- completed originals;
- failed originals.

Aggregate by profile × throughput:
- failures prevented versus baseline;
- earliest-deadline versus fixed-order;
- throughput gain relative to throughput 1.

## Interpretation

A monotonic reduction in failures as throughput increases establishes timing/throughput scarcity under synchronized deadlines. Failure to improve means total budget remains dominant even under this stress.

## Bounds

Counterfactual only. No live activation, no future adaptive budget, no capacity change, no extra training, no semantic priority, no production authority.
