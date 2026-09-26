# Wingless UP-112C — probation recurrence-window diagnostic

Status: preregistered scientific admission-memory diagnostic.

Scientific parent: sealed UP-111C bab1e8fdd41476f0d889e2e68dd7275c2b91d4d7.

## Question

UP-111C showed that a 32-entry exact uint32 generational probation table eliminates one-shot false admissions at the same 128-byte budget as the Bloom filter. Does it still admit genuinely repeated identities within the intended 32-write generation while rejecting repeats that cross the generation boundary?

## Frozen admission-memory arms

1. bloom1024
   - 1024-bit two-hash filter;
   - generation interval 32.

2. exact32
   - 32 uint32 identities;
   - generation interval 32.

Both:
- 128 bytes admission memory;
- one-byte generation counter;
- start every episode empty at counter 0;
- no exact-memory replacement layer in this diagnostic;
- no queries, phase labels, or future information.

## Repeat-distance ladder

A target identity is written once, then repeated after:
- 1
- 8
- 16
- 24
- 31
- 32
- 33
- 40 filtered writes

Definition:
"repeat after N" means first target write, then N-1 unique one-shot writes, then the target repeat.

Under the frozen 32-write generation:
- N <= 31 is expected to remain within generation;
- N >= 32 is expected to cross the clear boundary.

After the target repeat, continue with unique one-shot identities until 64 total filtered writes have been processed.

## Sampling

- 64 deterministic episodes per arm x distance;
- unique key ranges per episode;
- all keys fit uint32 exactly.

## Metrics

- target repeat admission rate;
- expected repeat decision accuracy;
- one-shot false admission count/rate;
- generation clears;
- admission-memory bytes.

## Interpretation

The desired exact32 behavior is a sharp local recurrence window: admit genuine repeats before the boundary, reject after it, and never admit one-shot churn. This validates that the collision-free mechanism retains useful admission semantics rather than merely rejecting everything.

## Bounds

Diagnostic only. No exact-memory layer, no interval tuning, no adaptive window, no queries, no future oracle, no result-informed retry, no live activation, no production authority.
