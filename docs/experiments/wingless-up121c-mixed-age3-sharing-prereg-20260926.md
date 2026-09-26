# Wingless UP-121C — heterogeneous persistence sharing at fixed admission memory

Status: preregistered scientific bounded-memory persistence experiment.

Scientific parent: sealed UP-120C d457b42196b2ee87cb701bed1da63f759701e840.

## Question

UP-120C showed that the equal-memory extended history admits a homogeneous candidate set through two skipped generations while rejecting one-generation bursts and failing at skip-three. Does that temporal capacity remain correct when adjacent, one-skip, and two-skip candidates share the same history table in one episode?

## Frozen mechanisms

Compare:
1. age1_control — exact UP-118C/UP-119C age1-history behavior.
2. age2_extended — exact UP-120C extended history.

Both:
- current table: 32 × uint16 = 64 bytes;
- history table: 32 × uint16 = 64 bytes;
- total admission sidecar: 128 bytes;
- 14-bit exact candidate keys;
- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- generation interval 32;
- phase-free slot-reuse checkpoint;
- no semantic priority labels;
- no query-derived admission evidence;
- no future oracle.

## Mixed workload

Desired valid recurring set: four candidates total.

- adjacent A: 2 sightings generation 1; 2 sightings generation 2.
- adjacent B: 2 sightings generation 1; 2 sightings generation 2.
- skip-one: 2 sightings generation 1; absent generation 2; 2 sightings generation 3.
- skip-two: 2 sightings generation 1; absent generations 2 and 3; 2 sightings generation 4.

Negative candidates:

- two burst distractors: 4 sightings generation 1 only.
- two skip-three negatives: 2 sightings generation 1; absent generations 2/3/4; 2 sightings generation 5.

Each generation is completed to exactly 32 filtered writes with unique one-shot fillers.

The desired exact-memory target remains:
- 12 hot residents + 4 valid recurring candidates = 16.

## Churn

After generation 5:
- 12,288 unique one-shot writes;
- query hot, all valid candidates, bursts, and skip-three negatives every four churn writes;
- 32 episodes per seed;
- seeds 229M and 230M.

## Metrics

Per arm:
- admission and final retention separately for adjacent, skip-one, skip-two, burst, and skip-three classes;
- 12-hot accuracy;
- desired target16 accuracy and exactness;
- negative retention rate;
- one-shot false admissions;
- exact-memory occupancy;
- checkpoint count.

## Interpretation

If age2_extended simultaneously preserves all four heterogeneous valid candidates while rejecting both negative classes, the fixed 128-byte history can share temporal persistence across mixed horizons. Any selective failure identifies competition within the compact history rather than a simple horizon limit.

## Bounds

No admission-memory increase, no exact-memory capacity increase, no adaptive horizon, no semantic labels, no query-derived admission, no future oracle, no result-informed retry, no live activation, no production authority.
