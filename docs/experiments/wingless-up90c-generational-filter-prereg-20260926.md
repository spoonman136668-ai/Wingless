# Wingless UP-90C — generational seen-once filter

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-89C 4834165acfdd79552d12517014d59cf30a2e75e0.

## Question

UP-89C showed that a persistent 256-bit seen-once filter preserves all 16 hot items through 24 unique churn writes but saturates by 48. Can deterministic generational clearing prevent stale one-shot evidence from accumulating while retaining the same admission rule?

## Frozen exact memory and admission mechanism

Identical to UP-89C:
- exact recall cap 16;
- two-bit-aging replacement;
- 256-bit two-hash seen-once filter;
- existing exact-memory keys always update;
- first-seen unseen writes are rejected and set their two bits;
- unseen writes whose two bits are already set are admitted;
- no query labels or future information.

## Frozen churn workload

- 32 original keys;
- 16 hot even keys;
- every hot key receives one second write before churn;
- 192 unique one-shot churn writes;
- query every hot key after each group of four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 169M and 170M.

## Frozen arms

1. no_reset
   - exact UP-89C persistent filter control.
   - policy metadata = 37 bytes.

2. reset_32
   - clear only the 256-bit seen-once filter after every 32 churn writes.
   - exact memory and aging state are not reset.

3. reset_64
   - clear filter after every 64 churn writes.

4. reset_96
   - clear filter after every 96 churn writes.

Reset arms add one byte of deterministic writes-since-reset state:
- policy metadata = 38 bytes.

The first reset occurs only during the churn phase; the pre-churn hot second-write phase is unchanged.

## Metrics

For every arm:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- rejected one-shot writes;
- false-positive churn admissions;
- filter reset count;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

If finite reset intervals restore exact retention, stale admission evidence—not the exact-memory capacity—is the limiting mechanism under long unique churn.

## Bounds

No exact-memory capacity expansion, no reset triggered by queries or accuracy, no future oracle, no adaptive interval, no threshold tuning, no result-informed retry, no attention, no live activation, no production authority.
