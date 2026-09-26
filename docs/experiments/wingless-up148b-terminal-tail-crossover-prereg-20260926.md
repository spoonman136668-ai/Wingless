# Wingless UP-148B — terminal-tail crossover under fixed coverage

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-147B 8a6c31cc4f3c3c439d8d297d953b847cfe8b53ae.

## Question

UP-147B found a tail-identity × curriculum-position interaction: shifts 0 and 12 become damaging late, while shifts 6 and 18 remain neutral-to-positive, and schedules ending in 6/18 finish better than schedules ending in 0/12. Does terminal tail identity causally drive that final retention gap when total tail coverage is held exactly fixed?

## Frozen starting point and budget

Use the exact UP-147B starting gate and block budget:
- 20 epochs;
- learning rate 0.08;
- 24 new-family updates per epoch;
- 15 old-family rehearsal updates per epoch;
- four final new-family updates;
- shifts {0,6,12,18};
- every shift used for exactly five consecutive epochs;
- exactly three block transitions;
- same projector and family.

## Preregistered matched crossover pairs

Pair A shares the first two blocks 0 → 6:
- safe_final_18: 0 → 6 → 12 → 18
- damaging_final_12: 0 → 6 → 18 → 12

Pair B shares the first two blocks 12 → 18:
- safe_final_6: 12 → 18 → 0 → 6
- damaging_final_0: 12 → 18 → 6 → 0

Within each pair the tail multiset, per-tail count, total updates, and first ten epochs are identical. Only the order of the last two five-epoch blocks is swapped.

## Measurements

Per arm:
- final old held-out and unseen retention;
- final primary and secondary new-family accuracy;
- mean final old retention;
- mean final new-family accuracy.

For block positions 3 and 4:
- mean prefix damage;
- mean anchor recovery;
- mean tail damage;
- mean net epoch change.

Derived pairwise deltas:
- safe-final minus damaging-final old retention;
- safe-final minus damaging-final new accuracy.

## Interpretation

If both matched pairs favor the safe-final identity, terminal-tail order is a causal contributor to the UP-147B final gap under fixed total coverage. If the paired effect is absent or reverses, the earlier schedule difference depends on a broader identity × history interaction rather than terminal identity alone.

No numeric success threshold is introduced.

## Bounds

No adaptive ordering, no extra updates, no new memory, no projector recomputation, no architecture change, no threshold search, no result-informed retry, no live activation, no production authority.
