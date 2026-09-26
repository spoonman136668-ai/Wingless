# Wingless UP-83C — bounded-memory retention decay curve

Status: preregistered scientific bounded-memory characterization experiment.

Scientific parent: sealed UP-82C `362fc1ef42363c3cc9da74c0638a4973bdc778e4`.

## Question

UP-81C and UP-82C show that two-bit aging preserves rehearsed hot information longer than FIFO/CLOCK and, under early rehearsal, longer than LRU. What is the fixed-budget retention decay curve as churn increases?

## Frozen policies

- FIFO control;
- LRU control;
- two-bit aging exactly as frozen in UP-81C.

Memory:
- exact recall capacity 16;
- 16-byte entry payload accounting;
- FIFO metadata 16 bytes;
- LRU metadata 136 bytes;
- two-bit aging metadata 5 bytes;
- no future-query oracle.

## Frozen workload

- active original keys: 24 and 32;
- hot set: even original keys;
- cold set: odd original keys;
- one rehearsal pass over the hot set immediately before churn;
- churn writes: 0, 4, 8, 12, 16, 20, 24 unseen keys;
- value vocabulary 32;
- 64 episodes per cell;
- two fixed seed bases.

## Metrics

For every policy x active-key count x churn-length cell:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- cold-set exact accuracy;
- entries used;
- metadata bytes;
- total bounded-memory bytes.

## Interpretation

The objective is a retention curve, not a pass/fail threshold. The experiment should identify where two-bit aging begins to lose its rehearsal advantage and how that boundary compares with LRU at fixed capacity.

## Bounds

No policy adaptation, no capacity expansion, no oracle, no result-informed retry, no attention, no live activation, no production authority.
