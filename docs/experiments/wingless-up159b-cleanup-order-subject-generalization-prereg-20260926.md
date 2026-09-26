# Wingless UP-159B — cleanup-order subject generalization

Status: preregistered scientific lexical-stability generalization experiment.

Scientific parent: sealed UP-158B 50d4146aa51fca6289a8271ca1e92b9a8ce801bb.

## Question

UP-158B showed that cleanup class order affects repair at fixed composition and budget, but no single order was uniformly best for Mia and Noah. Does that subject-by-order interaction generalize across the full accepted six-subject new-family set?

## Frozen design

Reuse UP-158B unchanged:
- accepted 15-epoch common prefix;
- 20 non-REPORT new-family updates;
- old rehearsal indices 10,11,12 before REPORT;
- fixed four-example new REPORT block;
- fixed cleanup composition after REPORT: 5 STORE, 5 OBSERVE, 2 REPORT;
- total old updates = 15;
- total new updates = 24;
- learning rate = 0.08.

## Subjects

Use the complete accepted new-family subject set:
- mia
- noah
- opal
- pax
- quin
- rue

No new names are introduced.

## Cleanup orders

Test the same six class-block permutations from UP-158B:
- STORE_OBSERVE_REPORT
- STORE_REPORT_OBSERVE
- OBSERVE_STORE_REPORT
- OBSERVE_REPORT_STORE
- REPORT_STORE_OBSERVE
- REPORT_OBSERVE_STORE

This yields 36 fixed subject × order arms.

## Measurements

Per subject × order:
- mean new-prefix effect;
- mean pre-REPORT rehearsal recovery;
- mean REPORT damage;
- mean post-REPORT cleanup recovery;
- mean net epoch change;
- final mean old retention;
- final mean new accuracy.

## Interpretation

The full 6×6 response is the result. If the order preference changes across subjects, cleanup recency is context-sensitive and cannot yet justify an adaptive universal scheduler. If one ordering pattern generalizes across the added subjects, it becomes a stronger candidate for fixed scheduling.

No order is selected post hoc and no numeric threshold is introduced.

## Bounds

No extra updates, no rehearsal coverage change, no class-composition change, no memory change, no projector recomputation, no architecture change, no adaptive ordering, no result-informed retry, no live activation, no production authority.
