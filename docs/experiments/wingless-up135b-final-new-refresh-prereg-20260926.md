# Wingless UP-135B — final new-family refresh at fixed update count

Status: preregistered scientific adaptation-retention ordering experiment.

Scientific parent: sealed UP-134B 28d28abf2039726a84bff6b7197ad3472207b541.

## Question

UP-134B showed that placing all 15 old-family rehearsal updates after the new-family block gives the best old-family retention, but secondary new-family accuracy remains slightly below the no-terminal-anchor control. Can a small final new-family refresh recover that loss without increasing the total update budget?

## Frozen mechanism and data

Exact UP-134B:
- dual-view 128-D three-way lexical gate;
- frozen STORE projector;
- two new grounding subjects: mia, noah;
- two balanced old rehearsal subjects: ada, ben;
- 20 epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family updates per epoch;
- no architecture/projector changes.

## Fixed-count ordering arms

At each epoch:
- construct the same 24 new-family examples;
- rotate the new-example list left by epoch modulo 24 so refresh status is not permanently tied to one surface or subject;
- construct the same 15 old-family rehearsal examples using the UP-134B balanced subject rule;
- every example is used exactly once per epoch.

Arms:

1. refresh_0
   - 24 new;
   - 15 old.
   - exact UP-134B terminal_15 control.

2. refresh_4
   - first 20 rotated new;
   - all 15 old;
   - final 4 rotated new.

3. refresh_8
   - first 16 rotated new;
   - all 15 old;
   - final 8 rotated new.

4. refresh_12
   - first 12 rotated new;
   - all 15 old;
   - final 12 rotated new.

Total new and old update counts are identical across arms.

## Evaluation

Same UP-134B evaluation:
- primary new-family accuracy on quin/rue;
- secondary new-family accuracy;
- worst new-surface accuracy;
- old held-out class accuracy;
- old unseen class accuracy;
- old STORE precision/recall.

## Interpretation

An intermediate refresh allocation that restores secondary new-family accuracy while preserving most of the terminal-anchor retention gain would show that the adaptation-retention tradeoff can be managed by update placement alone at fixed budget.

## Bounds

No extra updates, no adaptive split, no threshold search, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
