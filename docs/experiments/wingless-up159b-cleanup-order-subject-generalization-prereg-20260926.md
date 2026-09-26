# Wingless UP-159B — cleanup-order subject generalization

Status: preregistered scientific lexical-stability generalization experiment.

Scientific parent: sealed UP-158B 50d4146aa51fca6289a8271ca1e92b9a8ce801bb.

## Question

UP-158B showed that cleanup class order affects repair at fixed composition and budget, but no single order was uniformly best for Mia and Noah. Does that subject-by-order interaction generalize across the accepted six-name new-family set when every arm preserves the exact 24-update new-family budget?

## Frozen matched subject pairs

Use three preregistered two-subject corpora, each containing the same 12 verbs/classes per subject:
- mia / noah
- opal / pax
- quin / rue

For a subject arm, its matched pair supplies exactly 24 new-family examples per epoch. The target subject's fixed four-example REPORT block is removed from the 20-example prefix and placed at the interference point, so total new-family updates remain exactly 24.

The other four accepted new-family names are the held-out new-name evaluation set for that arm.

## Frozen epoch design

Reuse UP-158B:
- accepted 15-epoch common prefix;
- 20 matched-pair non-tail new-family updates;
- old rehearsal indices 10,11,12 before REPORT;
- fixed target-subject REPORT block: relays, announces, cites, summarizes;
- fixed cleanup composition after REPORT: 5 STORE, 5 OBSERVE, 2 REPORT;
- total old updates = 15;
- total new updates = 24;
- learning rate = 0.08.

## Subjects

- mia
- noah
- opal
- pax
- quin
- rue

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
- matched training pair;
- four held-out new-family names;
- mean new-prefix effect;
- mean pre-REPORT rehearsal recovery;
- mean REPORT damage;
- mean post-REPORT cleanup recovery;
- mean net epoch change;
- final mean old retention;
- final held-out new-family accuracy.

## Interpretation

The full 6×6 response is the result. If the cleanup-order preference changes across the added subject pairs, repair ordering is context-sensitive and cannot justify a universal scheduler. If a stable ordering pattern generalizes, it becomes a stronger fixed-schedule candidate.

No order is selected post hoc and no numeric threshold is introduced.

## Bounds

No extra updates, no rehearsal coverage change, no class-composition change, no memory change, no projector recomputation, no architecture change, no adaptive ordering, no result-informed retry, no live activation, no production authority.
