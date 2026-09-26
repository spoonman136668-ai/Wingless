# Wingless UP-141B — tail-diversity replication on a disjoint lexical family

Status: preregistered scientific adaptation-order generalization experiment.

Scientific parent: sealed UP-140B 3dfc55991f8c2835d1e93114ef7faeddea99e74d.

## Question

UP-140B showed that broad post-anchor tail diversity improves old-family retention for one newly grounded lexical family. Does the same scheduling mechanism replicate when the entire new verb surface family is replaced by a disjoint set while subjects, update counts, old-family rehearsal, and gate architecture remain fixed?

## Frozen gate / projector

- exact pre-adaptation dual-view 128-D three-way gate from the B lineage;
- frozen STORE competitor-nullspace projector;
- original old lexical family training unchanged;
- no weights from the prior UP-130B new-family adaptation are reused.

## Replication lexical family

Completely disjoint from the UP-130B/UP-140B new surfaces:

STORE:
- stows
- shelves
- tucks
- deposits

OBSERVE:
- audits
- reviews
- tracks
- samples

REPORT:
- briefs
- mentions
- narrates
- explains

Grounding subjects:
- mia
- noah

Primary evaluation subjects:
- quin
- rue

Secondary evaluation subjects:
- eli, fay, gia, hal, ivy, jon, kia, leo

No replication surface appears in the old-family gate training or the first new lexical family.

## Frozen training budget

All adaptation arms:
- 20 epochs;
- lr 0.08;
- exactly 24 replication-family updates per epoch;
- exactly 15 balanced old-family rehearsal updates per epoch;
- exactly 4 replication-family updates after the old anchor;
- every replication example used exactly once per epoch.

## Arms

Use the exact UP-140B shift construction:

1. diversity_1
   - shift set {0}.

2. diversity_8
   - shift set {0,3,6,9,12,15,18,21}.

3. diversity_20
   - shift=epoch for epochs 0..19.

For each shift:
- rotate the 24-example replication list left by shift;
- first 20 new;
- all 15 old-anchor updates;
- last 4 new.

## Evaluation

Before adaptation:
- zero-shot primary and secondary replication-family accuracy.

After each arm:
- primary and secondary replication-family overall/per-class accuracy;
- worst replication-surface accuracy;
- old held-out and old unseen class accuracy;
- old STORE precision/recall;
- mean old-retention score;
- mean new-family score.

## Interpretation

A retention improvement from diversity_1 to diversity_8/diversity_20 on this disjoint verb family would show the scheduling mechanism generalizes beyond a single lexical surface set. Failure to replicate would imply the previous effect depends on specific surface geometry.

## Bounds

No extra updates, no adaptive shift choice, no projector recomputation, no threshold search, no architecture change, no result-informed retry, no live activation, no production authority.
