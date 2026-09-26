# Wingless UP-111C — equal-memory exact generational probation table

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-110C 9954d1813002fbed244db9645572edce699cb256.

## Question

UP-110C showed that explicit and phase-free checkpoints fail identically at long churn because rare Bloom false admissions accumulate across many generations. Can the same 128-byte admission-memory budget eliminate that failure by storing exact 32-bit probationary identities for one 32-write generation instead of Bloom occupancy bits?

## Frozen exact memory and detector

- exact recall cap 16;
- two-bit-aging replacement unchanged;
- phase-free slot_reuse_confirmation detector unchanged;
- 2-byte reuse mask;
- checkpoint semantics unchanged;
- generation interval 32;
- no phase label;
- no query-derived detector evidence;
- no future oracle.

## Admission-memory arms

1. bloom1024
   - exact UP-110C 1024-bit two-hash filter.
   - filter memory: 128 bytes.

2. exact32
   - 32-slot table of uint32 key identities.
   - table memory: 32 × 4 = 128 bytes.
   - first unseen key in a generation is stored if absent;
   - a later matching key is admitted;
   - table clears at the same 32-filtered-write generation boundary;
   - checkpoint clears the table and sets generation counter to 8, exactly paralleling the Bloom arm.
   - workload key IDs are guaranteed to fit uint32 exactly.

Policy metadata is held at 136 bytes in both arms:
- replacement aging: 5 bytes;
- generation counter: 1 byte;
- slot reuse mask: 2 bytes;
- admission memory: 128 bytes.

## Workload

Initial orders:
- ascending
- reverse
- evens_then_odds
- odds_then_evens
- interleaved_low_high
- rotate8

Promotion:
- exact two-pass alternating-ends promotion.

Churn horizons:
- 24,576
- 49,152
- 98,304 unique one-shot writes

Sampling:
- seeds 209M and 210M;
- 16 episodes per seed per cell;
- 32 total episodes per cell.

## Metrics

- resident hot before churn;
- checkpoint count;
- false-positive admissions;
- final hot accuracy;
- final hot-set exact accuracy;
- recall entries;
- policy metadata bytes.

## Interpretation

Exact32 success at equal admission-memory cost would show that probabilistic occupancy—not bounded memory itself—is the long-horizon limit. Failure would implicate replacement/checkpoint dynamics beyond Bloom collisions.

## Bounds

No capacity increase, no interval tuning, no adaptive table size, no query evidence, no future oracle, no result-informed retry, no live activation, no production authority.
