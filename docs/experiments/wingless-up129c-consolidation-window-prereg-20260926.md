# Wingless UP-129C — consolidation window across a generation boundary

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-128C f0ee54626aee3fb99a48fcebfb276087bc31b9ca.

## Question

UP-128C found an exact two-sighting threshold for protecting pending targets before destructive Q2=4 overflow. Does that two-sighting rule depend on both sightings occurring inside the same 32-write admission generation, or can the same evidence consolidate when the sightings are separated across a generation boundary?

## Frozen mechanism

Use the exact UP-128C mechanism:
- exact recall cap 16;
- 32-entry current table;
- 32-entry history table;
- 128-byte admission sidecar;
- maximum-history-age replacement only;
- lowest-index tie break;
- no semantic priority;
- no query priority;
- no future oracle.

Reproduce exact saturation and the exact first 16-candidate overflow wave.

## Frozen evidence and timing arms

Every arm gives each of the four valid targets exactly two pre-overflow recurrence sightings and executes exactly 64 writes before the Q2=4 second overflow wave.

1. adjacent_same_generation
   - generation A starts with two adjacent sightings per target (8 writes total);
   - fill generation A with 24 unique one-shot writes;
   - generation B is 32 unique one-shot writes.

2. wide_same_generation
   - generation A: one sighting per target (4 writes);
   - 16 unique one-shot writes;
   - second sighting per target (4 writes);
   - 8 unique one-shot writes;
   - generation B is 32 unique one-shot writes.

3. boundary_split
   - generation A: one sighting per target (4 writes), then 28 unique one-shot writes;
   - generation B: second sighting per target (4 writes), then 28 unique one-shot writes.

Then all arms:
- present four Q2 candidates twice each;
- fill that generation with unique one-shot writes;
- execute one additional full 32-write one-shot generation;
- apply unchanged 12,288-write churn and final evaluation.

Seeds:
- 245000000;
- 246000000.

32 episodes per seed.

## Metrics

Per arm:
- panic rate;
- valid target admission/final accuracy;
- hot accuracy;
- target16 accuracy/exactness;
- nonpersistent admission/retention;
- false admissions;
- post-first-wave evictions and age distribution;
- max current/history occupancy;
- recall entries used.

## Interpretation

If both same-generation arms remain exact while boundary_split fails, the two-sighting consolidation rule has a hard generation-local evidence window. If boundary_split also succeeds, the qualifying evidence survives across the current-table boundary through the existing history mechanism.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
