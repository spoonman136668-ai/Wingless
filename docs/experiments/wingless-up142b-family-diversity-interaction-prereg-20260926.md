# Wingless UP-142B — paired lexical-family × tail-diversity interaction

Status: preregistered scientific sequential-lexical generalization experiment.

Scientific parent: sealed UP-141B 7040fd48fd135ab0d0e68966582b877f8d8d79f1.

## Question

UP-140B showed an old-retention gain from broad tail diversity on the first new lexical family, while UP-141B did not replicate that retention gain on a disjoint second family. Under one identical harness, is there a genuine lexical-family × tail-diversity interaction?

## Frozen design

Use the exact shared pre-adaptation gate, projector, learning rate, subjects, old-family rehearsal block, and 20-epoch budget from UP-140B/UP-141B.

Families:
1. original new lexical family from UP-140B;
2. disjoint replication family from UP-141B.

Arms for each family:
- diversity_1;
- diversity_8;
- diversity_20.

Every arm keeps:
- 24 new-family updates per epoch;
- 15 balanced old-family rehearsal updates per epoch;
- 4 new-family updates after the old anchor;
- every new-family example exactly once per epoch;
- no extra updates.

## Metrics

For every family × arm:
- primary and secondary new-family accuracy;
- worst surface accuracy;
- old held-out and old unseen accuracy;
- old STORE precision/recall;
- mean old-retention score;
- mean new-family score.

Derived, frozen comparisons:
- diversity20 minus diversity1 old-retention delta for each family;
- diversity20 minus diversity1 new-family delta for each family;
- difference between the two families' old-retention deltas.

## Interpretation

If the two families show materially different old-retention deltas under the same harness, the B140 mechanism is family/geometry dependent rather than universal. If both now show the same direction and similar magnitude, the B141 non-replication was more likely due to run-context or implementation differences.

No numeric success threshold is introduced; the exact measured interaction is the result.

## Bounds

No adaptive family choice, no threshold tuning, no extra memory, no extra updates, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
