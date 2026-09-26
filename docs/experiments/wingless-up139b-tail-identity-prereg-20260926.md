# Wingless UP-139B — rotating-tail identity ablation

Status: preregistered scientific adaptation-order mechanism experiment.

Scientific parent: sealed UP-138B 661d0328475ccf459a9acb19f8e35f6d79a08338.

## Question

UP-138B established a positive rotation × final-refresh interaction for old-family retention. Is the gain caused by moving which new-family examples occupy the four post-anchor tail slots, or by rotating the pre-anchor new-example order itself?

## Frozen mechanism and budget

- exact dual-view 128-D lexical gate;
- frozen STORE projector;
- new grounding subjects: mia, noah;
- balanced old rehearsal subjects: ada, ben;
- 20 epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family updates per epoch;
- exactly four new-family updates after the old anchor in every arm;
- every new example used exactly once per epoch;
- no extra updates.

The canonical 24-example new-family list is the exact UP-138B list.

## Arms

1. moving_tail_control
   - exact UP-138B rotated_refresh4.
   - rotate the full 24-example list left by epoch modulo 24.
   - first 20 rotated examples;
   - all 15 old-anchor examples;
   - final 4 rotated examples.
   - tail identity therefore changes with epoch.

2. fixed_tail_last4
   - freeze the canonical examples at indices 20..23 as the four post-anchor tail examples every epoch.
   - rotate only the remaining 20 examples left by epoch modulo 20 before the old anchor.
   - then all 15 old-anchor examples;
   - then the same four frozen tail examples.

3. fixed_tail_first4
   - freeze canonical indices 0..3 as the post-anchor tail every epoch.
   - rotate the remaining 20 examples before the old anchor.
   - then the same full old anchor;
   - then the four frozen tail examples.

4. alternating_fixed_tails
   - even epochs use canonical indices 20..23 as tail;
   - odd epochs use canonical indices 0..3 as tail.
   - rotate the corresponding remaining 20 examples before the old anchor.
   - this changes tail identity across epochs, but only between two frozen sets rather than sweeping the entire 24-example family.

## Evaluation

Same UP-138B metrics:
- primary and secondary new-family accuracy;
- worst new-surface accuracy;
- old held-out and old unseen class accuracy;
- old STORE precision/recall.

Also report:
- mean old-retention score = mean(old held-out, old unseen);
- mean new score = mean(primary, secondary).

## Interpretation

If only moving_tail_control preserves the B138 retention gain, broad tail-identity coverage is the likely mechanism. If alternating fixed tails is sufficient, only coarse tail diversity is needed. If a fixed-tail arm matches the control, pre-anchor rotation rather than tail mobility is the main driver.

## Bounds

No extra updates, no adaptive tail choice, no threshold search, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
