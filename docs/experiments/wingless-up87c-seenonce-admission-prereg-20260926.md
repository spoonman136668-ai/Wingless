# Wingless UP-87C — low-metadata seen-once admission gate

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-86C 58cfc6b5179fdc23a972866250ab380943f7601f.

## Question

UP-86C showed that replacement policy alone cannot preserve 13 or more hot items exactly under one-shot churn. Can a low-metadata reuse gate reject one-shot writes while admitting keys that demonstrate a second write, without future-query labels?

## Frozen memory

Control:
- exact UP-81C two-bit-aging memory;
- exact recall cap 16;
- 5 bytes replacement metadata.

Experimental:
- same two-bit-aging memory;
- plus a fixed 64-bit two-hash seen-once filter (8 bytes);
- total policy metadata 13 bytes.

Admission rule when exact memory is full:
- existing keys always update normally;
- for an unseen key, test the two seen-once filter bits;
- if both bits were already set, admit the current write using ordinary two-bit-aging replacement;
- otherwise set both bits and reject this write from exact memory.

Before exact memory becomes full, writes admit normally.

The filter uses only prior writes; no query labels or future information.

## Frozen workload

- active original keys: 32;
- hot-set sizes: 13 and 16;
- hot keys are the first N even original keys;
- all original keys receive one initial write in deterministic order;
- after initial loading, every hot key receives one second write with its current value, in deterministic order;
- cold original keys receive no second write;
- then 24 unseen churn keys each receive exactly one write;
- after every 4 writes in the second-hot-write and churn phases, query every hot key once;
- missing keys are not reconstructed by queries;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 165M and 166M.

## Arms

1. two_bit_aging
2. seen_once_gate_plus_aging

## Metrics

- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- cold-set exact accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes;
- number of rejected one-shot writes;
- number of second-write admissions.

## Interpretation

Exact preservation above the UP-86C 12-item boundary would support admission control as the missing mechanism. Failure remains valid and must not be repaired by enlarging the memory.

## Bounds

No capacity expansion, no future oracle, no query-derived admission labels, no adaptive thresholds, no policy tuning after results, no attention, no live activation, no production authority.
