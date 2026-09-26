# Wingless UP-136B — matched-order four-update refresh confirmation

Status: preregistered scientific adaptation-retention confirmation experiment.

Scientific parent: sealed UP-135B 6b0d70a5552403c4af843bbc5fe59d59829c3c67.

## Question

UP-135B suggested that four final new-family updates can improve old-family retention without reducing new-family accuracy, but its preregistered epoch rotation meant the refresh_0 arm did not exactly reproduce UP-134B terminal_15. Does the four-update refresh effect persist when the control uses the exact UP-134B fixed new-example order?

## Frozen mechanism

Exact UP-134B:
- dual-view 128-D three-way lexical gate;
- frozen STORE projector;
- new grounding subjects: mia, noah;
- balanced old rehearsal subjects: ada, ben;
- 20 epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family updates per epoch;
- no architecture/projector changes.

## Arms

1. exact_terminal15_control
   - exact fixed UP-133B/UP-134B new block of all 24 examples;
   - then exact balanced 15-example old block.
   - no rotation.

2. fixed_tail4_refresh
   - exact same fixed 24-example new order;
   - execute the first 20 new examples;
   - execute the exact same 15-example old block;
   - execute the final 4 new examples.
   - the final four are the fixed tail of the original new block; no example is added or duplicated.

Both arms use identical examples and counts.

## Evaluation

Same B134/B135 metrics:
- primary new-family accuracy;
- secondary new-family accuracy;
- worst new-surface accuracy;
- old held-out class accuracy;
- old unseen class accuracy;
- old STORE precision/recall.

## Interpretation

If fixed_tail4_refresh improves the retention/new-family balance relative to the exact control, the B135 result survives matched-order confirmation. If not, B135's apparent advantage was tied to the preregistered rotation rather than refresh placement itself.

## Bounds

No extra updates, no rotation, no adaptive split, no threshold search, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
