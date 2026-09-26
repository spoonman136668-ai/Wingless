# Wingless UP-130C — generation-edge causality

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-129C 608584f4c506cf8780a348fe00d4322690466e9f.

## Question

UP-129C showed that two target sightings consolidate when they occur in one admission generation, but fail when a generation boundary separates them. Is boundary crossing itself the causal cutoff, even when the sightings are placed immediately around the edge?

## Frozen mechanism and pressure

Reuse UP-129C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement;
- lowest-index tie break;
- first overflow wave = 16 qualified nonpersistent candidates;
- second overflow pressure Q2 = 4;
- no semantic/query/future priority.

Every arm gives each valid target exactly two pre-overflow sightings and exactly 64 pre-overflow writes.

## Arms

1. early_same_generation
   - four first sightings;
   - 24 one-shot writes;
   - four second sightings;
   - then one full 32-write one-shot generation.

2. edge_same_generation
   - 24 one-shot writes;
   - four first sightings;
   - four second sightings;
   - then one full 32-write one-shot generation.

3. edge_cross_boundary
   - 28 one-shot writes;
   - four first sightings as the final four writes of generation A;
   - four second sightings as the first four writes of generation B;
   - 28 one-shot writes to finish generation B.

Then all arms receive the identical Q2=4 second overflow wave, one further full one-shot generation, unchanged churn, and final evaluation.

Seeds:
- 247000000;
- 248000000.

32 episodes per seed.

## Interpretation

If both same-generation arms consolidate and edge_cross_boundary fails, the generation boundary itself is causal rather than absolute inter-sighting spacing or proximity to the edge.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
