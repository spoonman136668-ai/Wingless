# Wingless UP-117B — stores representation/projection diagnostic

Status: preregistered scientific lexical-geometry diagnostic.

Scientific parent: sealed UP-116B 73e575a82532a0886309a112c7501467d7bc5ebe.

## Question

UP-116B showed that only the base verb stores collapses, always toward OBSERVE, while keeps remains robust. Is stores intrinsically closer to OBSERVE surfaces in the frozen 64-D encoder, or does the collapse arise primarily from learned classifier-weight trajectory?

## Frozen system

No learning rule changes. Repeat the exact UP-116B training matrix:
- all six semantic-class acquisition orders;
- current_class_excluded and stage3_anchor2 policies;
- same encoder, base/acquired verbs, subjects, epochs, learning rate, and replay rules.

## Static encoder geometry

Before any classifier training, compute subject-averaged 64-D encoder vectors for:
- base STORE: stores, keeps;
- acquired STORE: holds, saves;
- base OBSERVE: observes, sees;
- acquired OBSERVE: notes, watches;
- base REPORT: reports, recalls;
- acquired REPORT: tells, remembers.

For stores and keeps report cosine similarity to every other verb vector.

Also report cosine to three frozen family centroids:
- STORE family centroid;
- OBSERVE family centroid;
- REPORT family centroid.

The verb being measured is excluded from its own family centroid when applicable.

## Dynamic classifier projections

At stage 0 and after every acquisition for every order x policy, evaluate stores and keeps over all six subjects.

Report for each:
- mean STORE logit;
- mean OBSERVE logit;
- mean REPORT logit;
- mean and minimum STORE-minus-OBSERVE logit margin;
- mean encoder dot-product contribution to STORE-minus-OBSERVE margin;
- classifier bias contribution to STORE-minus-OBSERVE margin;
- accuracy.

## Interpretation

If stores is already much closer to the OBSERVE centroid than keeps, encoder geometry is implicated. If static geometry is not exceptional but the learned STORE-minus-OBSERVE projection collapses dynamically, the failure is primarily classifier trajectory/interference.

## Bounds

No training modification, no replay change, no encoder change, no threshold tuning, no intervention selection, no result-informed retry, no live activation, no production authority.
