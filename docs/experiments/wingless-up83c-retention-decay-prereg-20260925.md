# Wingless UP-83C — bounded-memory retention decay ladder

Status: preregistered scientific bounded-memory robustness experiment.

Scientific parent: sealed UP-82C `362fc1ef42363c3cc9da74c0638a4973bdc778e4`.

## Question

UP-82C showed that two-bit aging preserves useful hot entries under early and mid rehearsal, but late rehearsal cannot recover entries already evicted. What is the retention-decay curve as churn increases under a single frozen early rehearsal?

## Frozen memory system

- exact recall cap 16;
- 16-byte payload accounting per occupied entry;
- no capacity increase;
- no future-query oracle.

Policies:
- FIFO control;
- LRU control;
- two-bit aging exactly as frozen in UP-81C.

## Frozen workload

- active original keys: 24 and 32;
- hot set = even original keys;
- cold set = odd original keys;
- exactly one rehearsal pass over the hot set immediately after initial loading;
- churn writes: 0, 4, 8, 12, 16, 20, 24, 32 unseen keys;
- value vocabulary 32;
- 64 episodes per cell;
- fixed seed bases 157M and 158M.

## Metrics

For each policy x active-key count x churn depth:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- cold-set exact accuracy;
- recall entries used;
- metadata bytes;
- total bounded-memory bytes.

## Interpretation

This experiment measures the survival curve rather than selecting a new policy. The two-bit aging parameters are frozen. A smooth advantage across multiple churn depths supports a real low-metadata locality mechanism; a narrow spike would indicate schedule-specific behavior.

## Bounds

No parameter search, no adaptive policy selection, no extra rehearsal, no oracle, no capacity expansion, no attention, no result-informed retry, no live activation, no production authority.
