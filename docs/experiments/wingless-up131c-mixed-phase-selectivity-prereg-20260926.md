# Wingless UP-131C — mixed-phase consolidation selectivity

Status: preregistered scientific bounded-memory integration experiment.

Scientific parent: sealed UP-130C 639b5607ea32406c092b7731ee565752bb322025.

## Question

UP-130C established that a generation boundary causally breaks two-sighting consolidation. In one mixed stream, can some targets consolidate while other equally valid targets fail solely because their two sightings straddle the boundary?

## Frozen mechanism and pressure

Reuse UP-130C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement;
- lowest-index tie break;
- first overflow wave = 16;
- Q2 = 4;
- no semantic/query/future priority.

## Arms

Four valid targets are split into two pairs. Every target receives exactly two pre-overflow recurrence sightings.

1. assignment_a
   - targets 100/101: both sightings inside generation A;
   - targets 102/103: first sighting at end of generation A, second at start of generation B.

2. assignment_b
   - targets 102/103: both sightings inside generation A;
   - targets 100/101: split across the boundary.

Each arm executes exactly two complete 32-write pre-overflow generations.

Then both arms receive the identical Q2=4 second-overflow wave, one full one-shot generation, unchanged churn, and final evaluation.

Seeds:
- 249000000;
- 250000000.

32 episodes per seed.

## Measurements

Per arm:
- same-generation pair admission/final accuracy;
- cross-boundary pair admission/final accuracy;
- hot accuracy;
- target16 accuracy/exactness;
- nonpersistent admission/retention;
- false admissions;
- max current/history occupancy;
- recall entries used.

## Interpretation

If the same-generation pair consolidates while the cross-boundary pair fails in both mirrored assignments, consolidation is selective per item and follows each item's temporal evidence path rather than a global stream state.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
