# Wingless UP-SQ0 — implementation freeze

Status: frozen before first execution of UP-SQ0.

This file adds only implementation details that were not specified in the original preregistration. It does not change any frozen task family, split, seed, arm, metric, exact-recall cap, or interpretation rule.

## Shared recurrent carrier

- recurrent state: 64 float64 scalars (512 bytes);
- transport between events: identity transform, which is reversible and unitary;
- associative binding: deterministic bipolar key/value vectors with dimension 64 and norm 1;
- each key/value vector is generated from the frozen generator seed, key id, and value id using a deterministic xorshift64 stream;
- decoding uses dot-product similarity across the task's value vocabulary;
- ties resolve to the numerically smallest value id;
- prior-value presence is declared only when the strongest absolute similarity is at least 0.25; otherwise the write is treated as first insertion and no old-value subtraction occurs;
- no learned parameter, oracle state, hidden state expansion, or post-result fit is used.

## Closed-loop write rule

For a write (key,value):

1. decode the carrier's current value for that key;
2. subtract one unit vector for the decoded old value when a prior value is present;
3. add one unit vector for the new value.

The write rule is identical across arms before optional correction / exact recall.

## Gated-correction arm

- endogenous confidence is the normalized decoder margin between the intended new value and the strongest competing value immediately after the base write;
- normalized margin = (intended_score - strongest_competing_score) / max(1, abs(intended_score) + abs(strongest_competing_score));
- frozen correction threshold: 0.15;
- if the intended margin is below 0.15, exactly one additional copy of the intended key/value vector is added;
- no repeated correction loop is permitted;
- the threshold is not selected or changed after execution.

## Bounded exact recall arm

- the third arm uses the same gated carrier plus a 16-entry exact key/value table;
- replacement policy: FIFO by first insertion; rewriting an existing key updates its value without increasing entry count or changing its FIFO position;
- queries consult exact recall first, then fall back to carrier decoding;
- capacity is exactly 16 entries;
- estimated exact-recall bytes are frozen as 16 bytes per occupied entry for measurement accounting;
- no attention matrix or global sequence cache is present.

## Episode counts

For every seed, arm, and setting:

- delayed copy: 32 episodes;
- selective copy: 32 episodes;
- MQAR: 64 episodes;
- overwrite/latest-value-wins: 64 episodes;
- state-machine composition: 128 episodes;
- role/filler recombination: all 54 structural train tuples and all 162 held-out tuples;
- multi-bank interference: 64 episodes.

## Deterministic sequence generators

- seed bases remain exactly 120M, 121M, 122M;
- each family derives per-episode streams only by deterministic arithmetic/xorshift expansion of those seed bases;
- no result-dependent seed replacement is permitted.

## State-machine operators

The six 8-state operators are fixed permutations, frozen here before execution:

- op0: [1,2,3,4,5,6,7,0]
- op1: [7,0,1,2,3,4,5,6]
- op2: [0,2,4,6,1,3,5,7]
- op3: [3,0,7,4,1,6,2,5]
- op4: [2,5,0,7,4,1,6,3]
- op5: [6,3,5,1,7,2,0,4]

Held-out depth-8 and depth-16 operator sequences are generated from the frozen seeds and are not present as complete compositions in the depth-1..4 reference set.

## Resource measurement

Per arm/family:

- recurrent-state bytes: fixed 512;
- exact-recall bytes: occupied entries × 16;
- peak process RAM proxy: maximum Go heap allocation observed before/after the family evaluation using runtime.MemStats;
- peak VRAM: 0 because this qualification is CPU-only;
- throughput: task events processed / measured wall time;
- wall-clock qualification time: measured with time.Since.

These measurements are diagnostics and do not control scientific pass/fail.

## Scheduling

UP-SQ0 remains serialized. It must not execute concurrently with A/B/C qualification runs.
