# Wingless UP-132C — staggered item-phase selectivity

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-131C 0a778b520a912ab099a054c5b103489478b67fee.

## Question

UP-131C showed item-selective consolidation when two target pairs followed clean same-generation versus cross-boundary paths. Does the same rule hold when four targets have overlapping, staggered recurrence phases at different positions in the same two-generation stream?

## Frozen mechanism and pressure

Reuse UP-131C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement;
- lowest-index tie break;
- first overflow wave = 16;
- Q2 = 4;
- no semantic/query/future priority.

## Frozen 64-write recurrence schedule

Each target receives exactly two recurrence sightings. Generation A is positions 1–32; generation B is positions 33–64.

Four temporal paths:
- same_early: positions 1 and 16;
- same_late: positions 17 and 32;
- cross_edge: positions 31 and 33;
- cross_wide: positions 20 and 45.

All unoccupied positions are unique one-shot writes. Total pre-overflow writes remain exactly 64.

Two mirrored assignments:
- assignment_a: targets 100,101,102,103 map respectively to same_early, same_late, cross_edge, cross_wide.
- assignment_b: targets 102,103,100,101 map respectively to same_early, same_late, cross_edge, cross_wide.

Then both arms receive the identical Q2=4 second-overflow wave, one full one-shot generation, unchanged churn, and final evaluation.

Seeds:
- 251000000;
- 252000000.

32 episodes per seed.

## Measurements

Per arm × temporal path:
- admission rate;
- final target accuracy.

Per arm:
- hot accuracy;
- target16 accuracy/exactness;
- nonpersistent admission/retention;
- false admissions;
- max current/history occupancy;
- recall entries used;
- panic rate.

## Interpretation

If same_early and same_late consolidate while cross_edge and cross_wide fail in both mirrored assignments, the generation-local rule remains item-selective under overlapping staggered phases and is insensitive to absolute within-generation position.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
