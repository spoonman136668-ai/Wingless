# Wingless UP-146B — repeated-tail habituation trajectory

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-145B 4bf88b261a8179a6c484fc190e445c2892d210fc.

## Question

UP-145B showed that three switches with five consecutive uses per tail set outperform schedules with 11 or 19 switches, while medium and high switching plateau. Does the blocked advantage arise because interference changes as Wingless repeatedly sees the same post-anchor tail identity?

## Frozen starting point and budget

Use the exact UP-145B original-family starting gate and low-switch schedule:
- 20 epochs;
- learning rate 0.08;
- 24 new-family updates per epoch;
- 15 old-family rehearsal updates per epoch;
- four final new-family updates;
- shifts {0,6,12,18};
- each shift repeated for exactly five consecutive epochs;
- exactly three between-block switches.

No training count or order is changed.

## Measurements

For each epoch:
- tail identity / shift;
- repetition position within the current five-epoch block (1..5);
- prefix damage;
- anchor recovery;
- tail damage;
- net epoch change.

Aggregate across the four tail identities at each repetition position 1..5:
- mean prefix damage;
- mean anchor recovery;
- mean tail damage;
- mean net epoch change.

Also report final old held-out/unseen retention and primary/secondary new-family accuracy.

## Interpretation

A systematic reduction in post-anchor tail damage across repetition positions would support local habituation to a repeated lexical tail as the mechanism behind blocked scheduling. A flat or irregular trajectory would shift the explanation toward block-boundary geometry, cumulative training state, or another longer-range ordering effect.

No numeric success threshold is introduced; the exact frozen trajectory is the scientific result.

## Bounds

No adaptive ordering, no extra updates, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
