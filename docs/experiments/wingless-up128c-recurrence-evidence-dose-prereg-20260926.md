# Wingless UP-128C — pre-overflow recurrence evidence dose

Status: preregistered scientific bounded-memory timing experiment.

Scientific parent: sealed UP-127C 521b2123f4cd64459255eb37c2d176a93513a1f2.

## Question

UP-127C showed that two target recurrences before overflow protect all four valid targets, while the same two recurrences after overflow are too late. Holding total recurrence evidence fixed at exactly two sightings per target, how many of those sightings must occur before overflow to protect the targets?

## Frozen mechanism and sequence

Use the exact UP-127C mechanism and saturation path:
- exact recall cap 16;
- 32-entry current table;
- 32-entry history table;
- 128-byte admission sidecar;
- maximum-history-age eviction;
- lowest-index tie-break;
- first overflow wave of 16 qualified nonpersistent candidates;
- second overflow wave of four qualified nonpersistent candidates;
- unchanged 12,288-write churn;
- no semantic priority;
- no query-derived priority;
- no future oracle.

Every valid target receives exactly two recurrence sightings total.

## Preregistered arms

- pre0_post2: zero target sightings before second overflow; two after.
- pre1_post1: one target sighting before second overflow; one after.
- pre2_post0: two target sightings before second overflow; zero after.

No arm receives extra sightings.

## Seeds and episodes

- seeds: 243000000 and 244000000;
- 32 episodes per seed;
- 64 episodes per arm.

## Metrics

Per arm:
- panic rate;
- first- and second-wave eviction counts;
- evictions by age;
- valid-target admission/final accuracy;
- hot accuracy;
- target16 accuracy and exact rate;
- nonpersistent admission/retention;
- false admissions;
- maximum current/history occupancy;
- exact recall entries used.

## Interpretation

If pre1_post1 remains destructive while pre2_post0 protects targets, the current mechanism requires a complete two-sighting recurrence event before eviction; one early sighting is insufficient. If pre1_post1 provides partial protection, the mechanism carries useful pre-eviction evidence across the overflow boundary.

No post-result threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive timing, no extra recurrence sightings, no result-informed retry, no live activation, no production authority.
