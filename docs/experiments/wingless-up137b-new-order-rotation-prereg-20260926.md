# Wingless UP-137B — new-family example-order rotation under full old anchor

Status: preregistered scientific adaptation-order experiment.

Scientific parent: sealed UP-136B e86911ab08a6daef957c4c7873e60aa21e1e3695.

## Question

UP-136B showed that the fixed four-update refresh itself is not beneficial, while UP-135B's rotated schedule produced better retention. Is epoch-wise rotation of the new-family example order the actual mechanism improving the adaptation-retention balance when the full old-family terminal anchor is held fixed?

## Frozen mechanism and counts

Exact UP-136B / UP-134B:
- dual-view 128-D three-way lexical gate;
- frozen STORE projector;
- new grounding subjects: mia, noah;
- balanced old rehearsal subjects: ada, ben;
- 20 epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family updates per epoch;
- every epoch ends with the complete 15-update old-family anchor;
- no extra updates.

## Arms

1. fixed_new_order
   - exact UP-136B exact_terminal15_control.
   - the 24 new examples use the same fixed order every epoch.
   - then all 15 old rehearsal updates.

2. rotate_new_by_epoch
   - use the same 24 new examples exactly once each epoch;
   - rotate the new-example list left by epoch modulo 24;
   - then all 15 old rehearsal updates.
   - exact new-order rotation used by UP-135B refresh_0.

3. reverse_rotate_new_by_epoch
   - reverse the fixed 24-example new list once;
   - rotate that reversed list left by epoch modulo 24;
   - then all 15 old rehearsal updates.
   - tests whether coverage rotation rather than a particular cyclic direction is responsible.

Old-family rehearsal order and subject assignment remain identical across arms.

## Evaluation

Same B136 metrics:
- primary new-family accuracy;
- secondary new-family accuracy;
- worst new-surface accuracy;
- old held-out class accuracy;
- old unseen class accuracy;
- old STORE precision/recall.

## Interpretation

If both rotating arms improve retention relative to fixed order without harming new-family accuracy, epoch-wise coverage rotation is the supported mechanism. If only one direction works, the effect is order-specific rather than generic rotation.

## Bounds

No extra updates, no final new refresh, no adaptive ordering, no threshold search, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
