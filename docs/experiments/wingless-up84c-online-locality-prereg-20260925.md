# Wingless UP-84C — online locality under bounded memory

Status: preregistered scientific bounded-memory robustness experiment.

Scientific parent: sealed UP-83C `3f6899cdbe85a666377e0724475b52205078a743`.

## Question

UP-83C quantified a retention-decay advantage for two-bit aging after one batch rehearsal. Does that advantage persist when useful keys are revisited naturally throughout loading and churn rather than receiving one special rehearsal phase?

## Frozen memory budget

- exact recall capacity exactly 16 entries;
- 16 payload bytes per occupied entry;
- no capacity increase;
- no future-query oracle.

## Frozen policies

1. FIFO — unchanged.
2. LRU — unchanged.
3. two_bit_aging — exact UP-81C policy, unchanged:
   - 2-bit age per entry;
   - insertion age 1;
   - update/query age 3;
   - deterministic clock hand;
   - 5 policy metadata bytes.

## Frozen workload

Active original keys:
- 24
- 32

Hot set:
- even-numbered original keys.

Cold control:
- odd-numbered original keys.

Initial loading:
- write original keys in deterministic order;
- after every 4 original-key writes, query every currently loaded hot key once.

Churn:
- 16 or 24 unseen writes;
- after every 4 churn writes, query every hot original key once;
- missing keys remain missing and are not reconstructed.

Final evaluation:
- query complete hot set;
- query complete cold control.

Other controls:
- value vocabulary 32;
- 64 episodes per setting;
- seeds 161M and 162M.

## Metrics

For every policy x active-key count x churn cell:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- cold-set exact accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

A two-bit-aging advantage under this access pattern would support a genuine online-locality mechanism rather than a batch-rehearsal artifact. If it disappears, the earlier gains are schedule-specific.

## Bounds

No future-query labels, no adaptive policy selection, no capacity expansion, no policy-parameter tuning, no attention, no result-informed retry, no live activation, no production authority.
