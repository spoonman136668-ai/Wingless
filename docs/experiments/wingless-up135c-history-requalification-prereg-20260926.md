# Wingless UP-135C — clean history-expiry requalification

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-134C e0678ea72a932beed1e1fcd24db2bd88e43af27a.

## Question

UP-134C showed that a later within-generation pair can fail after the earlier qualified history trace has aged out, refining consolidation into a two-stage current+history rule. Can an expired item rebuild a fresh history trace with one pair and then become durable when a second pair appears in the next generation?

## Isolation correction

Each arm contains exactly one recurrence target. This removes a counter-progression confound from mixed streams: once an item is already in exact recall, later hits return before incrementing the admission generation counter.

Thus every arm below executes exactly 128 counter-advancing pre-overflow writes, preserving four exact 32-write generations.

## Frozen mechanism and pressure

Reuse UP-134C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement;
- lowest-index tie break;
- first overflow wave = 16;
- Q2 = 4;
- no semantic/query/future priority.

Two target identities are tested independently: 100 and 103.

## Temporal-path arms

Generation A = positions 1–32, B = 33–64, C = 65–96, D = 97–128.

1. live_pair_g1
   - positions 15,16.
   - positive control while the preexisting qualified history trace is live.

2. expired_pair_g3
   - positions 5,37,74,75.
   - two isolated singles age the old trace out; the G3 pair should only rebuild history.

3. rebuild_g3_admit_g4
   - positions 5,37,74,75,106,107.
   - same G3 rebuilding pair, then a second pair in G4.

4. expired_pair_g4_only
   - positions 5,37,69,106,107.
   - isolated singles through G3, then only one pair in G4.

All unoccupied positions are unique one-shot writes.

After the 128-write schedule, every arm receives the identical Q2=4 pressure, one full one-shot generation, unchanged churn, and final evaluation.

Seeds:
- 257000000;
- 258000000.

32 episodes per seed.

## Measurements

Per target identity × temporal path:
- total sightings;
- admission rate;
- final target accuracy;
- hot accuracy;
- nonpersistent admission/retention;
- false admissions;
- max current/history occupancy;
- recall entries used;
- panic rate.

## Interpretation

If expired_pair_g3 and expired_pair_g4_only remain unadmitted but rebuild_g3_admit_g4 admits, requalification after expiry requires a fresh two-stage pair→history→pair cycle. live_pair_g1 verifies the original live-history admission path.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
