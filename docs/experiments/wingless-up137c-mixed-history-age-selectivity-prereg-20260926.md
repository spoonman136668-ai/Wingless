# Wingless UP-137C — mixed history-age selectivity

Status: preregistered scientific bounded-memory temporal-consolidation integration experiment.

Scientific parent: sealed UP-136C 94a7f56273f44d5a9300c34f2689b83a69b7a55f.

## Question

UP-136C established the rebuilt-history lifetime boundary in isolated single-target arms: a G3 rebuilt trace supports admission in G4, G5, or G6 but expires before G7. Do targets at different history ages obey those boundaries independently when they coexist in the same counter-exact stream?

## Frozen mechanism

Reuse UP-136C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement with maxAge=2;
- lowest-index tie break;
- no semantic/query/future priority.

## Frozen seven-generation mixed schedule

Seven exact 32-write generations, 224 counter-advancing positions total. Four recurrence targets coexist. No target is touched again after it becomes durable.

Temporal paths:
- next_g4: isolated G1/G2 sightings, rebuild pair in G3, second pair in G4.
- after_two_blank_g6: isolated G1/G2 sightings, rebuild pair in G3, second pair in G6.
- after_three_blank_g7: isolated G1/G2 sightings, rebuild pair in G3, second pair in G7.
- rebuild_only: isolated G1/G2 sightings and rebuild pair in G3, no second pair.

Assignment A:
- target100 → next_g4
- target101 → after_two_blank_g6
- target102 → after_three_blank_g7
- target103 → rebuild_only

Assignment B mirrors identities:
- target102 → next_g4
- target103 → after_two_blank_g6
- target100 → after_three_blank_g7
- target101 → rebuild_only

All other positions are unique one-shot writes.

Seeds:
- 261000000;
- 262000000.

32 episodes per seed.

## Interpretation

If G4 and G6 targets admit while G7 and rebuild-only targets fail in both mirrored assignments, history aging remains selective per item under simultaneous temporal states rather than collapsing into a global stream-age effect.

No post-result threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
