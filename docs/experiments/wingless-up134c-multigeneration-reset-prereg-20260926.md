# Wingless UP-134C — multi-generation reset and later qualification

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-133C e84d69e7daa206134e909ce111ceeee5e9297a7c.

## Question

UP-133C showed that a failed cross-boundary pair can recover once two later sightings coexist in one generation. Across four consecutive generations, do isolated single sightings remain fully non-accumulating across resets while a later same-generation pair still qualifies normally?

## Frozen mechanism and pressure

Reuse UP-133C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement;
- lowest-index tie break;
- first overflow wave = 16;
- Q2 = 4;
- no semantic/query/future priority.

## Frozen 128-write recurrence schedule

Four consecutive generations:
- G1 positions 1–32;
- G2 positions 33–64;
- G3 positions 65–96;
- G4 positions 97–128.

Four temporal paths:

1. isolated_4
   - positions 5, 37, 69, 101
   - exactly one sighting in each generation.

2. late_pair_g3
   - positions 10, 42, 74, 75
   - one sighting in G1, one in G2, then two adjacent sightings in G3.

3. early_pair_then_singles
   - positions 15, 16, 50, 82
   - two sightings in G1, then one in G2 and one in G3.

4. pair_after_three_isolated
   - positions 20, 52, 84, 116, 117
   - one sighting in G1, G2, and G3, then two adjacent sightings in G4.

All other positions are unique one-shot writes.

Two mirrored assignments:
- assignment_a: targets 100,101,102,103 map respectively to the four paths.
- assignment_b: targets 102,103,100,101 map respectively to the four paths.

Then both arms receive the identical Q2=4 second-overflow wave, one full one-shot generation, unchanged churn, and final evaluation.

Seeds:
- 255000000;
- 256000000.

32 episodes per seed.

## Measurements

Per arm × path:
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

If isolated_4 remains unadmitted despite four total sightings while late_pair_g3 and pair_after_three_isolated qualify exactly when a same-generation pair appears, evidence does not accumulate across resets. early_pair_then_singles verifies that once qualified, later isolated sightings do not destabilize durable recall.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
