# Wingless UP-118C — equal-memory cross-generation age-2 persistence

Status: preregistered scientific bounded-memory persistence-horizon experiment.

Scientific parent: sealed UP-117C c2df738f6ea8974036f7138696fddedfa1a39f78.

## Question

UP-117C showed that adjacent-generation recurrence cleanly rejects one-generation burst distractors, but its previous-generation table can only recognize recurrence in the immediately following generation. Can the same 128-byte admission budget support one skipped generation of persistence without admitting one-generation bursts?

## Frozen exact memory

- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- phase-free slot-reuse checkpoint;
- no semantic priority labels;
- no query-derived admission evidence;
- no future oracle;
- no explicit phase labels.

## Equal-memory arms

Both arms use exactly 128 bytes of admission sidecar memory.

1. previous_generation_control
   - exact UP-117C crossgen_2x2:
   - 32-entry uint16 current table + 32-entry uint16 previous-qualifier table;
   - candidate must receive at least two sightings in adjacent generations.

2. age2_history
   - 32-entry uint16 current table = 64 bytes;
   - 32-entry uint16 history table = 64 bytes;
   - all experimental keys remain below 2^14;
   - history entry stores key plus a two-bit generation age;
   - a key that qualifies with >=2 sightings enters history at age 0;
   - at each generation boundary, history ages by one and entries older than age 1 expire;
   - a nonresident candidate is admitted when it qualifies with >=2 sightings in the current generation and exists in history at age 0 or age 1.
   - current-generation qualification refreshes its history age to 0.

## Workload cells

A. adjacent_recurrence
- four target candidates: 2 sightings in generation 1 and 2 sightings in generation 2.
- four burst distractors: 4 sightings in generation 1 only.

B. skip_one_generation
- four target candidates: 2 sightings in generation 1, absent in generation 2, 2 sightings in generation 3.
- four burst distractors: 4 sightings in generation 1 only.

Each generation is completed to exactly 32 filtered writes with one-shot fillers.

## Churn

After the target recurrence generation:
- 12,288 unique one-shot churn writes;
- hot, target, and burst keys queried every four churn writes;
- 32 episodes per seed;
- seeds 223M and 224M.

## Metrics

Per arm x workload:
- target admission/retention;
- burst admission/retention;
- hot accuracy;
- target16 accuracy/exactness;
- one-shot false admissions;
- exact-memory occupancy;
- checkpoint count;
- admission metadata bytes.

## Interpretation

If age2_history preserves skip-one-generation targets while rejecting one-generation bursts at the same 128-byte admission budget, cross-generation persistence can be extended without capacity growth. Failure would establish that this compact history representation is insufficient.

## Bounds

No admission-memory increase, no exact-memory capacity increase, no semantic labels, no query-derived admission, no adaptive horizon, no future oracle, no result-informed retry, no live activation, no production authority.

## Pre-execution feasibility correction

Before any qualification run, the churn horizon was corrected from 24,576 to 12,288 unique writes because the preregistered 14-bit exact key encoding can represent at most 16,383 nonzero key codes. The recurrence mechanism, equal 128-byte budget, generation timing, seeds, candidate workload, and interpretation criteria are unchanged.
