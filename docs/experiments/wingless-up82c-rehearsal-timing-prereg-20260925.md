# Wingless UP-82C — two-bit aging rehearsal-timing robustness

Status: preregistered scientific bounded-memory robustness experiment.

Scientific parent: sealed UP-81C `54d90fc5c2a95a37b9afb12898005afa7d7f6614`.

## Question

UP-81C showed that a 5-byte two-bit aging policy matches LRU under moderate churn and retains some hot information under full churn. Is that gain robust to when rehearsal occurs, or is it an artifact of the single rehearse-then-churn schedule?

## Frozen memory system

- exact recall cap 16;
- 16-byte payload accounting per occupied entry;
- no capacity increase;
- no future-query oracle.

Policies:
- FIFO control;
- LRU control;
- two-bit aging exactly as frozen in UP-81C.

## Frozen workloads

- active original keys: 24 and 32;
- churn writes: exactly 16 unseen keys;
- hot set: even original keys;
- cold set: odd original keys;
- exactly one attempted rehearsal pass over the hot set;
- value vocabulary 32;
- 64 episodes per cell;
- two fixed seed bases.

Rehearsal timing:
1. `early`: before any churn write;
2. `mid`: after 8 of 16 churn writes;
3. `late`: after all 16 churn writes, immediately before final evaluation.

Rehearsal only refreshes keys that still exist; missing keys are not reconstructed.

## Metrics

For every policy x active-key count x timing cell:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- cold-set exact accuracy;
- entries used;
- metadata bytes;
- total bounded-memory bytes.

## Interpretation

A robust aging advantage should persist across more than the single early-rehearsal schedule. Strong timing dependence is a valid boundary and should remain visible.

## Bounds

No oracle, no adaptive policy selection, no capacity expansion, no result-informed retry, no attention, no live activation, no production authority.
