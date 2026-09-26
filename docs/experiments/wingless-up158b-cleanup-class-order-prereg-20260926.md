# Wingless UP-158B — cleanup class-block order

Status: preregistered scientific lexical-stability attribution experiment.

Scientific parent: sealed UP-157B b969ff7c6b4c7ac8ba3eb83b073216780cabc392.

## Question

UP-157B showed that at a fixed 12-after / 3-before rehearsal split, repair is strongest when three old REPORT examples are moved before the damaging new REPORT block, leaving a cleanup tail composed of five STORE, five OBSERVE, and two REPORT examples. With that exact class composition frozen, does the order of the three cleanup class blocks change repair?

## Frozen epoch design

Reconstruct the accepted 15-epoch common prefix.

For each terminal epoch and each new-family subject (Mia, Noah):
- 20 non-REPORT new-family updates;
- old rehearsal original indices 10,11,12 before REPORT;
- fixed four-example new REPORT block;
- old rehearsal original indices {0..9,13,14} after REPORT;
- all 15 old-rehearsal items used exactly once;
- total old updates = 15, total new updates = 24;
- learning rate remains 0.08.

The cleanup composition is always:
- 5 STORE;
- 5 OBSERVE;
- 2 REPORT.

## Cleanup-order arms

Within each class, original-index order is canonical. Test all six class-block permutations:
- STORE_OBSERVE_REPORT
- STORE_REPORT_OBSERVE
- OBSERVE_STORE_REPORT
- OBSERVE_REPORT_STORE
- REPORT_STORE_OBSERVE
- REPORT_OBSERVE_STORE

No order is selected adaptively.

## Measurements

Per subject × cleanup order:
- mean new-prefix effect;
- mean pre-REPORT rehearsal recovery;
- mean REPORT damage;
- mean post-REPORT cleanup recovery;
- mean net epoch change;
- final mean old retention;
- final mean new accuracy.

## Interpretation

The six-point response separates fixed class composition from within-cleanup recency. If orders ending in STORE or OBSERVE consistently repair more than orders ending in REPORT, cleanup recency remains class-structured even after composition is held constant.

No numeric success threshold is introduced.

## Bounds

No extra updates, no rehearsal coverage change, no class-composition change, no memory change, no projector recomputation, no architecture change, no adaptive ordering, no result-informed retry, no live activation, no production authority.
