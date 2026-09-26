# Wingless UP-154B — fixed-budget rehearsal placement

Status: preregistered scientific lexical-stability intervention experiment.

Scientific parent: sealed UP-153B 00e708d49751c637e3bb2f82018a8f314bdf1ac6.

## Question

UP-153B established that REPORT interference is causally anchor-relative: the same REPORT block damages old retention when it follows rehearsal, while moving it before rehearsal substantially improves retention. Can the interference be mitigated by redistributing the same fixed 15 old-rehearsal updates around the REPORT block without adding any updates?

## Frozen starting state and new-family ordering

Reconstruct the exact accepted 15-epoch common prefix.

For each terminal epoch and each subject independently (Mia, Noah):
- use all 20 non-REPORT-tail new-family examples first, in canonical order;
- use the fixed four-example REPORT block: relays, announces, cites, summarizes;
- use exactly the same 15 old-rehearsal examples and subject alternation as the accepted old anchor;
- learning rate remains 0.08;
- total new updates = 24;
- total old updates = 15.

## Rehearsal-placement arms

1. all_before
   - 20 new examples
   - all 15 old-rehearsal updates
   - 4 REPORT updates

2. split_8_7
   - 20 new examples
   - first 8 old-rehearsal updates
   - 4 REPORT updates
   - remaining 7 old-rehearsal updates

3. all_after
   - 20 new examples
   - 4 REPORT updates
   - all 15 old-rehearsal updates

The old-rehearsal multiset and order are identical across arms; only the split point changes.

## Measurements

Per subject × arm:
- common-prefix old retention;
- effect of the 20-example new prefix;
- pre-REPORT rehearsal recovery;
- REPORT-block effect;
- post-REPORT rehearsal recovery;
- net epoch change;
- final old held-out/unseen retention;
- final mean old retention;
- final primary/secondary new-family accuracy;
- final mean new-family accuracy.

## Interpretation

If moving some or all rehearsal after REPORT reduces net forgetting without extra updates, rehearsal placement is a fixed-budget mitigation mechanism. The exact three-point response is the result; no arm is selected adaptively.

## Bounds

No extra updates, no memory change, no projector recomputation, no architecture change, no adaptive scheduling, no result-informed retry, no live activation, no production authority.
