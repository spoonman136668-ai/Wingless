# Wingless UP-93C — generation-salted seen-once hashing

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-92C ce595308c7c181f5a43a1f8e6e7a5b9c8014a810.

## Question

UP-92C showed that continuous 32-write generations prevent cumulative saturation through 384 churn writes, but static hash geometry produces the first false admission by 768. Can a deterministic generation-dependent hash salt decorrelate later key ranges without changing exact-memory capacity or filter width?

## Frozen memory and admission system

Common to all arms:
- exact recall cap 16;
- two-bit-aging replacement;
- 256-bit two-hash seen-once filter;
- continuous generation interval 32 filtered unseen writes;
- existing exact-memory updates do not consume generation count;
- first-seen filtered writes are rejected and set both bits;
- filtered writes whose two bits are already set are admitted;
- after the 32nd filtered write, clear only the seen-once filter and begin the next generation;
- exact memory and aging state are never reset;
- no query labels or future information.

## Frozen arms

1. static_hash
   - exact UP-92C hash geometry.
   - metadata: 38 bytes.

2. generation_salted_hash
   - maintain an unsigned one-byte generation index in addition to the existing interval counter.
   - before the two hashes, mix the key with a deterministic salt derived from the generation index using the existing sq0Mix64 primitive.
   - both first and repeated writes within a generation use the same salt.
   - increment generation index only when the 32-write generation clears.
   - metadata: 39 bytes.
   - generation index cannot wrap in the tested horizons.

## Frozen workload

- 32 original keys;
- 16 hot even keys;
- each hot key receives one second write before churn;
- unique one-shot churn depths:
  - 768
  - 1536
  - 3072 writes
- hot queries after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm/cell;
- seeds 175M and 176M.

## Metrics

For every arm x churn depth:
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

A salted arm that postpones or removes false admissions would show that the UP-92C horizon is a deterministic hash-geometry boundary rather than unavoidable filter saturation.

## Bounds

No filter-width increase, no exact-memory expansion, no adaptive salt, no query-derived salt, no future oracle, no interval tuning, no threshold tuning, no result-informed retry, no live activation, no production authority.
