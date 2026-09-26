# Wingless UP-120C — equal-memory cross-generation age-3 horizon

Status: preregistered scientific bounded-memory persistence-horizon experiment.

Scientific parent: sealed UP-119C b49b8ecc045b71c2788fe3935bbbd3dc761576f4.

## Question

UP-119C showed that the current equal-memory age2 history supports adjacent and one-skipped-generation recurrence but fails after two skipped generations. The frozen uint16 history encoding already has two age bits. Can extending the allowed history age by one generation recover two-skip persistence at the same 128-byte admission-memory budget?

## Frozen exact memory

- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- phase-free slot-reuse checkpoint;
- generation interval 32;
- no semantic priority labels;
- no query-derived admission evidence;
- no future oracle;
- no explicit phase labels.

## Equal-memory arms

Both arms:
- 32-entry uint16 current table = 64 bytes;
- 32-entry uint16 history table = 64 bytes;
- total admission sidecar = 128 bytes;
- 14-bit exact candidate keys;
- current generation qualifies a key after at least two sightings.

1. age1_control
   - exact UP-118C/UP-119C age2_history behavior;
   - history ages 0 and 1 are valid;
   - age-1 entries expire at the next boundary.

2. age2_extended
   - same representation and tables;
   - history ages 0, 1, and 2 are valid;
   - age-2 entries expire at the next boundary.
   - qualifying again refreshes age to 0.

No additional table entries or bytes are added.

## Workloads

Each cell has four target candidates and four burst distractors.

Burst distractors:
- 4 sightings in generation 1 only.

Target recurrence cells:
1. adjacent
   - 2 sightings generation 1;
   - 2 sightings generation 2.

2. skip_one
   - generation 1, absent generation 2, recur generation 3.

3. skip_two
   - generation 1, absent generations 2 and 3, recur generation 4.

4. skip_three_boundary
   - generation 1, absent generations 2/3/4, recur generation 5.
   - deliberately beyond the extended age2 horizon.

Each target occurrence generation gives exactly 2 sightings.
Each generation is completed to 32 filtered writes with one-shot fillers.

## Churn

After recurrence:
- 12,288 unique one-shot writes;
- query hot, target, and burst keys every four churn writes;
- 32 episodes per seed;
- seeds 227M and 228M.

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

If age2_extended supports adjacent, skip-one, and skip-two while rejecting bursts and failing only at skip-three, the temporal horizon can be extended one generation at no memory cost. If longer retention increases false admission or disrupts hot retention, temporal persistence has a selection cost even without extra bytes.

## Bounds

No admission-memory increase, no exact-memory capacity increase, no adaptive horizon, no semantic labels, no query-derived admission, no future oracle, no result-informed retry, no live activation, no production authority.
