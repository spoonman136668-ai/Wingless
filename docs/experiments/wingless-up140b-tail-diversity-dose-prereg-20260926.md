# Wingless UP-140B — post-anchor tail-identity diversity dose

Status: preregistered scientific adaptation-order mechanism experiment.

Scientific parent: sealed UP-139B c9963417bf1ecd9ec4fdc3353dcb577801febd94.

## Question

UP-139B showed that broad epoch-to-epoch mobility of the four post-anchor new-family examples preserves substantially more old-family accuracy than fixed or two-tail schedules. Is that benefit a graded function of how many distinct tail identities are visited?

## Frozen mechanism and budget

- exact dual-view 128-D lexical gate;
- frozen STORE projector;
- new grounding subjects: mia, noah;
- balanced old rehearsal subjects: ada, ben;
- 20 epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family updates per epoch;
- exactly four new-family updates after the old anchor;
- every new example used exactly once per epoch;
- no extra updates.

The canonical 24-example new-family list is the exact UP-139B list.

For a shift s:
- rotate the full canonical list left by s;
- execute the first 20 rotated examples;
- execute all 15 old-anchor updates;
- execute the last 4 rotated examples.

Thus shift controls both pre-anchor order and tail identity exactly as in the successful moving-tail mechanism.

## Arms

1. diversity_1
   - shift set: {0}
   - every epoch uses shift 0.

2. diversity_2
   - shift set: {0, 12}
   - cycle through the two shifts by epoch.

3. diversity_4
   - shift set: {0, 6, 12, 18}
   - cycle by epoch.

4. diversity_8
   - shift set: {0, 3, 6, 9, 12, 15, 18, 21}
   - cycle by epoch.

5. diversity_20
   - shift = epoch for epochs 0..19.
   - exact 20-epoch moving-tail schedule used by UP-139B moving_tail_control.

No arm adapts its shifts to results.

## Evaluation

Same UP-139B metrics:
- primary new-family accuracy;
- secondary new-family accuracy;
- worst new-surface accuracy;
- old held-out and old unseen class accuracy;
- old STORE precision/recall;
- mean old-retention score;
- mean new-family score.

Also report the realized number of distinct tail sets in each arm.

## Interpretation

A monotonic or saturating retention improvement with tail diversity would establish identity coverage as a dose-dependent mechanism. A sharp threshold would identify the minimum diversity needed. Non-monotonic behavior would indicate specific order structure rather than generic diversity.

## Bounds

No extra updates, no adaptive shift choice, no threshold search, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
