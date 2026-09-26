# Wingless UP-138C — candidate-counter time versus raw event time

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-137C 85c0cd809c80ef71dd5b4475b8506f52feef0026.

## Question

UP-136C/137C established a finite history-age boundary. Is that age driven by the admission machine's candidate-write counter, or merely by the number of raw stream operations that occur between evidence bursts?

## Frozen mechanism

Reuse the accepted age-eviction machine unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maxAge = 2;
- lowest-index tie break;
- no semantic/query/future priority.

Each arm uses one recurrence target only.

## Frozen setup

As in the accepted lifetime experiments:
- initialize hot exact recall and bounded history;
- allow the original target history to expire via isolated sightings in G1/G2;
- rebuild fresh qualified history with a same-generation pair at positions 74 and 75 in G3;
- finish G3, so the rebuilt trace is stored in history at age 0.

## Matched raw-operation delay

After G3 boundary, execute exactly 96 raw operations in every arm.

Arms differ only in how many of those 96 are candidate writes:
- write_0_query_96
- write_32_query_64
- write_64_query_32
- write_96_query_0

Candidate writes are unique one-shot keys and are distributed deterministically across the 96 operation slots. Query operations are read-only exact-memory queries to frozen hot keys and do not invoke the admission candidate process.

Then present the same target twice consecutively.

Two target identities:
- 100
- 103.

Seeds:
- 263000000;
- 264000000.

32 episodes per seed.

## Measurements

Per target × arm:
- raw operations;
- candidate writes;
- read-only queries;
- candidate-write boundaries crossed;
- target-history presence and encoded age immediately before the second pair;
- admission rate;
- final target accuracy;
- hot accuracy;
- false admissions;
- max current/history occupancy;
- recall entries used;
- panic rate.

## Interpretation

If arms with 0/32/64 candidate writes retain usable history while 96 candidate writes expires it despite identical 96 raw operations, temporal lifetime is defined by candidate-counter time rather than raw event count.

No numeric threshold is introduced after observation.

## Bounds

No memory increase, no replacement-rule change, no query-side mutation, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
