# Wingless UP-153B — anchor-side block swap

Status: preregistered scientific lexical-sequencing mechanism experiment.

Scientific parent: sealed UP-152B 5f7ca686414f67cd56f0e126f87df9c90188e2c7.

## Question

UP-152B showed that post-anchor REPORT-tail damage increases monotonically with REPORT count, while final retention is non-monotonic because moving REPORTs into the tail also removes them from the pre-anchor prefix. Is the old-memory interference specifically caused by which class block appears after the anchor?

## Frozen common prefix

Reconstruct the exact accepted 15-epoch prefix:
- epochs 1–5: shift 0;
- epochs 6–10: shift 6;
- epochs 11–15: shift 12.

All arms start from the identical resulting gate.

## Frozen swapped blocks

For each subject independently (Mia, Noah), define two fixed four-example blocks:

REPORT block:
- relays
- announces
- cites
- summarizes

STORE+OBSERVE block:
- caches
- files
- scans
- checks

The other 16 new-family examples form a fixed common prefix within each terminal epoch.

Two arms:

1. report_after_anchor
   - 16 common examples
   - STORE+OBSERVE block
   - unchanged old anchor
   - REPORT block

2. report_before_anchor
   - 16 common examples
   - REPORT block
   - unchanged old anchor
   - STORE+OBSERVE block

Every epoch uses all 24 new-family examples exactly once, 15 old-anchor updates, and four post-anchor updates. Total update count and example multiset are identical.

## Measurements

Per subject × arm:
- old retention at the common-prefix state;
- pre-anchor new-block effect;
- old-anchor recovery;
- post-anchor block effect;
- net epoch change;
- final old held-out/unseen retention;
- final mean old retention;
- final primary/secondary new-family accuracy;
- final mean new-family accuracy.

## Interpretation

If swapping the blocks reverses the sign or materially changes the post-anchor effect and final retention, class interference is causally anchor-relative. If both orders converge, the UP-152B dose pattern depended on broader epoch ordering rather than the anchor boundary itself.

No numeric success threshold is introduced.

## Bounds

No adaptive ordering, no extra updates, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
