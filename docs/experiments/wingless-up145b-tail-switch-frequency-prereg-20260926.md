# Wingless UP-145B — tail switch-frequency dose at fixed coverage

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-144B bdd4b4e30dcee2403d611e50f5da3b3af061df39.

## Question

UP-144B found that blocked recurrence outperformed fully interleaved recurrence when tail identities and counts were identical. Is there a graded switching cost as the number of tail-identity transitions increases?

## Frozen starting point and budget

Use the exact UP-144B original-family starting point and training budget:
- 20 epochs;
- learning rate 0.08;
- 24 new-family updates per epoch;
- 15 old-family rehearsal updates per epoch;
- four final new-family updates;
- shifts {0,6,12,18};
- every shift used exactly five times.

## Arms

The only difference is the number of between-epoch shift transitions.

1. low_switch — 3 switches
   - 0×5, 6×5, 12×5, 18×5.

2. medium_switch — 11 switches
   - 0,0,6,6,12,12,18,18,
   - 0,0,6,6,12,12,18,18,
   - 0,6,12,18.

3. high_switch — 19 switches
   - 0,6,12,18 repeated five times.

Every arm has exactly the same shift multiset and update count.

## Measurements

Same fixed within-epoch measurements as UP-144B:
- prefix damage;
- anchor recovery;
- tail damage;
- net epoch change;
- final old held-out/unseen retention;
- final primary/secondary new-family accuracy.

## Interpretation

A monotonic degradation with switch count would support transition frequency as a causal sequencing cost. A non-monotonic pattern would imply longer-range ordering or block-boundary geometry rather than switch count alone.

No numeric success threshold is introduced.

## Bounds

No adaptive ordering, no extra updates, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
