# Wingless UP-132B — old-family rehearsal coverage at fixed update count

Status: preregistered scientific lexical-adaptation retention experiment.

Scientific parent: sealed UP-131B 72638010e1a03e3b32249f6c84ee05b38b21f0e5.

## Question

UP-131B showed that two new-family grounding subjects yield high new-family accuracy but old-family retention drops to about 87–90%. Can broader old-subject rehearsal recover retention without increasing rehearsal update count or weakening new-family acquisition?

## Frozen gate / projector

Exact UP-131B:
- raw 64-D representation;
- frozen STORE competitor-nullspace projected view;
- 128-D three-way gate;
- original 15-surface gate training first;
- projector never recomputed;
- new-family grounding surfaces unchanged.

## Frozen new-family grounding

Use exactly two grounding subjects:
- mia
- noah

For 20 epochs at lr 0.08, train all 12 new surfaces on both subjects exactly as UP-131B two_subjects.

## Old-family rehearsal arms

Every arm receives exactly 15 old-family rehearsal examples per epoch: one example for each original trained surface. Only subject coverage changes.

1. ada_only
   - all 15 surfaces use ada.
   - exact UP-131B two-subject control.

2. two_subject_balanced
   - old subjects: ada, ben.
   - for surface index s and epoch e, subject=(s+e) mod 2.

3. four_subject_balanced
   - old subjects: ada, ben, cy, dee.
   - for surface index s and epoch e, subject=(s+e) mod 4.

No arm receives extra rehearsal updates.

## Evaluation

New family:
- primary held-out: quin, rue;
- secondary: eli, fay, gia, hal, ivy, jon, kia, leo.

Old family:
- held-out original surfaces on eli/fay;
- prior unseen subjects gia, hal, ivy, jon, kia, leo.

Metrics per arm:
- new-family overall/per-class accuracy;
- worst new-surface accuracy;
- old-family retained class accuracy;
- STORE precision/recall on old family.

## Interpretation

If broader subject coverage improves old retention at fixed rehearsal count without hurting new-family accuracy, forgetting is partly a rehearsal-coverage problem rather than an unavoidable capacity tradeoff.

## Bounds

No extra rehearsal examples, no projector change, no architecture change, no threshold search, no adaptive subject choice, no result-informed retry, no live activation, no production authority.
