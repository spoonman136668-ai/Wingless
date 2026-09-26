# Wingless UP-143B — within-epoch lexical interference trajectory

Status: preregistered scientific sequential-lexical interference experiment.

Scientific parent: sealed UP-142B 110d063dee8ecbdadb5996de9afd1fa295f0d2d3.

## Question

UP-142B established a lexical-family × tail-diversity interaction: broad tail diversity improved old-family retention for the original new family but not for the disjoint replication family. At what fixed stage inside each adaptation epoch is that family-specific interaction created?

## Frozen starting point

Use the exact UP-142B shared pre-adaptation gate:
- 64-D state;
- 128-D dual-view gate input;
- frozen STORE competitor-nullspace projector;
- exact old-family training history.

Families:
1. original UP-140B new lexical family;
2. disjoint UP-141B replication family.

## Frozen adaptation budget

For each family, run:
- diversity_1;
- diversity_20.

Every arm:
- 20 epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 balanced old-family rehearsal updates per epoch;
- exactly 4 new-family updates after the old anchor;
- every new-family example exactly once per epoch.

No extra training updates are introduced by measurement.

## Fixed within-epoch measurements

For every epoch, measure old-family retention at exactly four points:

1. before_epoch
2. after_new_prefix — after first 20 new-family updates
3. after_old_anchor — after all 15 old-family rehearsal updates
4. after_new_tail — after final 4 new-family updates

For every epoch derive:
- prefix damage = after_new_prefix - before_epoch;
- anchor recovery = after_old_anchor - after_new_prefix;
- tail damage = after_new_tail - after_old_anchor;
- net epoch change = after_new_tail - before_epoch.

Old retention is the exact mean of old held-out and old unseen class accuracy used in UP-142B.

After epoch 20 also report:
- final old held-out / unseen retention;
- final primary / secondary new-family accuracy;
- mean new-family accuracy.

## Interpretation

A family-specific difference in prefix damage, anchor recovery, or tail damage identifies where the UP-142B interaction is generated. If stage trajectories are nearly identical despite different final effects, the mechanism lies in cumulative geometry across epochs rather than one local phase.

No post-result numeric success threshold is used; the exact trajectories and derived differences are the scientific result.

## Bounds

No adaptive ordering, no extra updates, no memory increase, no projector recomputation, no architecture change, no threshold search, no result-informed retry, no live activation, no production authority.
