# Wingless UP-116C — exact3 recurrence timing and stronger-distractor boundary

Status: preregistered scientific bounded-memory confidence experiment.

Scientific parent: sealed UP-115C b1623072a56658734db37d149a0b50824ae17fc2.

## Question

UP-115C showed that exact3 can select four three-hit recurring candidates over four two-hit candidates while preserving 12 established hot residents. Does that result hold across mixed recurrence timing, and what happens when distractors themselves recur four times?

## Frozen mechanism

Exact UP-115C exact3:
- exact recall cap 16;
- 12 established hot residents;
- 32 x uint32 probation sidecar = 128 bytes;
- bounded two-bit occurrence count encoded in each sidecar entry;
- admit on third sighting within a 32-write generation;
- two-bit-aging replacement;
- phase-free slot-reuse checkpoint;
- no query-derived admission evidence;
- no semantic priority labels;
- no future oracle.

## Scenarios

### mixed_gap_clean

Four weak candidates:
- exactly 2 sightings;
- recurrence gaps: 4, 12, 20, 28 filtered writes.

Four strong candidates:
- exactly 3 sightings;
- first-to-second gaps: 5, 9, 13, 17;
- second-to-third gaps: 6, 10, 14, 12;
- all three sightings remain within one generation.

Desired retained target:
- 12 established hot + 4 strong = 16.

### four_hit_distractor

Use the same four three-hit strong candidates.

Add four distractor candidates:
- exactly 4 sightings within one generation;
- deterministic spacing 4/4/4 filtered writes.

No semantic signal distinguishes target strong candidates from stronger distractors.

This deliberately creates legitimate recurrence demand above the 16-entry cap.

## Churn

After candidate presentation:
- 49,152 unique one-shot churn writes;
- query hot and all candidate keys every four churn writes;
- 32 episodes per seed;
- seeds 219M and 220M.

## Metrics

Per scenario:
- weak / strong / distractor admission and retention;
- hot accuracy;
- target16 accuracy and exactness;
- all-candidate accuracy;
- one-shot false admissions;
- exact-memory occupancy;
- checkpoint count.

## Interpretation

Success in mixed_gap_clean would show exact3 is insensitive to benign recurrence timing inside the generation. Failure under four_hit_distractor would establish the expected boundary: recurrence strength alone cannot distinguish semantically important three-hit items from even stronger repeated distractors.

## Bounds

No admission-memory increase, no capacity increase, no semantic labels, no query-derived admission, no adaptive threshold, no phase labels, no future oracle, no result-informed retry, no live activation, no production authority.


## Frozen presentation order addendum

This ordering is fixed before execution:

- mixed_gap_clean: weak keys 100..103 first, then strong keys 200..203.
- four_hit_distractor: target strong keys 200..203 first, then distractor keys 300..303.
- candidate keys are presented in ascending numeric order within each group.
- each candidate begins at a fresh probation generation boundary.

The ordering is deterministic and is not changed after observing results.
