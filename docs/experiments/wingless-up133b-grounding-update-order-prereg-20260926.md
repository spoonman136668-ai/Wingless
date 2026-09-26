# Wingless UP-133B — lexical grounding update-order isolation

Status: preregistered scientific adaptation-retention experiment.

Scientific parent: sealed UP-132B e55b98cfb0449a090db4ac248896673235f4f9d3.

## Question

UP-132B showed that broader old-subject rehearsal improves retention at fixed rehearsal volume. With the successful two-subject new grounding and two-subject balanced old rehearsal held fixed, does new-vs-old update order explain the remaining adaptation-retention tradeoff?

## Frozen mechanism and data

- exact UP-129B dual-view three-way gate;
- frozen STORE projector;
- original 15-surface gate training first;
- new family: exact UP-130B 12 surfaces;
- new grounding subjects: mia, noah;
- old rehearsal subjects: ada, ben;
- 20 grounding epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family rehearsal updates per epoch;
- old subject assignment for surface index s and epoch e = (s+e) mod 2.

## Arms

1. new_then_old
   - all 24 new-family updates;
   - then all 15 old-family rehearsal updates.
   - exact UP-132B two_subject_balanced control.

2. old_then_new
   - all 15 old-family rehearsal updates;
   - then all 24 new-family updates.

3. alternating_epoch_order
   - even epochs: new_then_old;
   - odd epochs: old_then_new.

No arm changes update counts, examples, subjects, learning rate, projector, or architecture.

## Evaluation

New family:
- primary held-out: quin, rue;
- secondary: eli, fay, gia, hal, ivy, jon, kia, leo.

Old family:
- held-out original surfaces on eli/fay;
- prior unseen subjects gia, hal, ivy, jon, kia, leo.

Metrics:
- new-family overall/per-class accuracy;
- worst new-surface accuracy;
- old-family retained class accuracy;
- old-family STORE precision/recall.

## Interpretation

Differences isolate recency/order effects at identical lexical coverage and update budget.

## Bounds

No extra updates, no projector change, no threshold search, no adaptive ordering, no architecture change, no result-informed retry, no live activation, no production authority.
