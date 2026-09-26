# Wingless UP-LM2O — unreported-eviction warning

Status: preregistered scientific fixed-capacity shadow-warning experiment.

Scientific parent: sealed UP-LM2N 7bda1b312da8b0958481b9f03f08a496b1476a75.

## Question

Can Wingless predict, before an incoming unique STORE is committed, whether that write will evict an entry that has never yet been REPORTed, using only current native recall state and past event history?

## Frozen streams

Use:
- exact recall cap 16;
- 24 entities in two chunks of 12;
- deferred first-chunk report levels D=4,5,6;
- identity rotations 0 and 7;
- value shifts 0,1,2,3.

No language-model training is needed.

## Frozen shadow rule

Maintain only:
- current recall FIFO order;
- current recall occupancy;
- set of entity names already REPORTed in the past.

Immediately before each second-chunk unique STORE:
- if recall is not full, predict no eviction;
- if recall is full, inspect the current FIFO front;
- warn iff that would-be victim has not yet been REPORTed.

The predictor does not inspect any future report schedule or future outcome.

## Preregistered expectation

Per arm:
- D=4: 0 unreported evictions;
- D=5: 1 unreported eviction;
- D=6: 2 unreported evictions.

The warning should exactly identify these events before the responsible write.

## Measurements

TP/FP/FN/TN for unreported-entry eviction, precision, recall, and counts by D.

## Bounds

Shadow only. No intervention, no future oracle, no semantic priority, no extra memory, no adaptive rule, no result-informed retry, no live activation, or production authority.
