# Wingless UP-133C — post-boundary recovery

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-132C db7f157a362455000e084299637572c9a3704eff.

## Question

UP-132C completed the boundary map: two sightings consolidate if they coexist within one generation and fail if split across the reset. After a split failure, can a third sighting recover the item once a later generation contains two sightings?

## Frozen mechanism and pressure

Reuse UP-132C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement;
- lowest-index tie break;
- first overflow wave = 16;
- Q2 = 4;
- no semantic/query/future priority.

## Frozen 64-write schedules

Generation A is positions 1–32; generation B is positions 33–64.

Four temporal paths:

1. split_control
   - positions 5 and 37
   - one sighting per generation; two total.

2. pre2_post1
   - positions 10, 20, 40
   - two sightings in generation A, one in generation B.

3. pre1_post2
   - positions 15, 45, 55
   - one sighting in generation A, then two in generation B.

4. edge_post2
   - positions 30, 33, 34
   - one sighting near the end of generation A, then two adjacent sightings at the start of generation B.

All other positions are unique one-shot writes. Total pre-overflow writes remain exactly 64.

Two mirrored assignments rotate target identities:
- assignment_a: targets 100,101,102,103 map respectively to the four paths above.
- assignment_b: targets 102,103,100,101 map respectively to the four paths above.

Then both arms receive the identical Q2=4 second-overflow wave, one full one-shot generation, unchanged churn, and final evaluation.

Seeds:
- 253000000;
- 254000000.

32 episodes per seed.

## Measurements

Per arm × temporal path:
- total sightings;
- admission rate;
- final accuracy.

Per arm:
- hot accuracy;
- target16 accuracy/exactness;
- nonpersistent admission/retention;
- false admissions;
- max current/history occupancy;
- recall entries used;
- panic rate.

## Interpretation

If split_control remains unadmitted while pre1_post2 and edge_post2 recover once generation B contains two sightings, boundary-split failure is reversible and evidence qualification restarts cleanly after the reset. pre2_post1 verifies the symmetric case where consolidation occurs before the boundary.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
