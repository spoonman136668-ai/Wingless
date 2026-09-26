# Wingless UP-134B — terminal old-family anchor allocation

Status: preregistered scientific adaptation-retention experiment.

Scientific parent: sealed UP-133B 9c724753be853be217da23f2355b2841464d0b40.

## Question

UP-133B showed that ending every epoch with old-family rehearsal materially improves retention. How many of the fixed 15 old-family rehearsal updates must be reserved as a terminal anchor block to obtain that benefit?

## Frozen data and mechanism

Exact UP-133B:
- dual-view 128-D three-way lexical gate;
- frozen STORE projector;
- two new grounding subjects: mia, noah;
- two balanced old rehearsal subjects: ada, ben;
- 20 epochs;
- lr 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family updates per epoch;
- no projector or architecture change.

## Old rehearsal ordering

At each epoch, construct the same 15 old rehearsal examples.
- surface order is rotated left by epoch modulo 15 before splitting, so terminal status is not permanently tied to a particular lexical surface;
- subject for rotated-position example j uses the original surface index s and subject=(s+epoch) mod 2.

Arms:

1. terminal_0
   - all 15 old updates first;
   - then all 24 new updates.
   - exact old_then_new boundary.

2. terminal_5
   - first 10 rotated old updates;
   - all 24 new updates;
   - final 5 rotated old updates.

3. terminal_10
   - first 5 rotated old updates;
   - all 24 new updates;
   - final 10 rotated old updates.

4. terminal_15
   - all 24 new updates;
   - all 15 rotated old updates.
   - exact new_then_old boundary modulo harmless old-order rotation.

## Evaluation

Same UP-133B new/old subject sets and metrics:
- primary and secondary new-family accuracy;
- worst new-surface accuracy;
- old held-out and old unseen class accuracy;
- old STORE precision/recall.

## Interpretation

A monotonic or threshold-like retention response would quantify the amount of terminal old-family evidence required to counter recency forgetting at fixed update budget.

## Bounds

No extra updates, no adaptive split, no threshold search, no architecture change, no projector recomputation, no result-informed retry, no live activation, no production authority.
