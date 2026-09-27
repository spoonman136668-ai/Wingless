# Wingless UP-LM2X — global budget allocation

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM2W 4941bb92181dd47a4917e67efa1ebc92180ae3d1.

## Question

When several independent Wingless dependency streams are simultaneously under pressure and only one corrective action can be taken per global round, does native time-to-failure improve action allocation over a fixed arm order?

## Frozen concurrent scenarios

Run two scenarios, identity rotations 3 and 11.

Each scenario contains six independent recall arms:
- deferred levels 4,5,6;
- pending layouts suffix_reported and alternating_reported;
- value shift 0;
- exact recall cap 16 per arm.

Each global round:
1. optionally perform at most one early-closure action globally;
2. then every arm receives one unique pressure write.

Twelve global rounds total.

## Policies

- **baseline:** no actions.
- **earliest_deadline:** choose the arm with the smallest current native countdown to its first unreported eviction; close that first pending dependency. Ties use fixed arm index.
- **fixed_order:** choose the lowest-index arm that still has any pending dependency, regardless of countdown.

Both active policies have the same maximum throughput: one action per round.

Countdown is computed only from current recall order, free slots, and reported state. No future schedule or outcome is used.

## Measurements

Per scenario/policy:
- actions;
- final completed originals;
- final failed originals.

Aggregate:
- failures prevented versus baseline;
- prevention per action.

## Interpretation

If earliest-deadline beats fixed-order under identical throughput, native countdown has scheduling value beyond merely identifying endangered dependencies.

## Bounds

Counterfactual only. No live activation, future schedule input, capacity growth, extra training, adaptive throughput, or production authority.
