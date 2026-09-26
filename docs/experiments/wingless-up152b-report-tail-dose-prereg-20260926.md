# Wingless UP-152B — REPORT-tail dose

Status: preregistered scientific lexical-sequencing mechanism experiment.

Scientific parent: sealed UP-151B 069a1fbc5c0db1b16c1a8b51acf004f8f8685229.

## Question

UP-151B showed that terminal class composition is the main cause of the safe/damaging split: REPORT-only tails damage old retention for both subjects, while STORE+OBSERVE mixtures are non-damaging. Does interference scale with the number of REPORT examples in the fixed four-example post-anchor tail?

## Frozen common prefix

Reconstruct the exact accepted 15-epoch prefix:
- epochs 1–5: shift 0;
- epochs 6–10: shift 6;
- epochs 11–15: shift 12.

All arms start from the identical resulting gate.

## Frozen nested REPORT-dose tails

Tail size is always four examples. For each subject independently (Mia, Noah), start with the accepted protective mixture and replace one non-REPORT example at a time:

- report_0: caches, files, scans, checks
- report_1: caches, files, scans, relays
- report_2: caches, files, relays, announces
- report_3: caches, relays, announces, cites
- report_4: relays, announces, cites, summarizes

Thus REPORT count is exactly 0,1,2,3,4 while total tail size remains four. Every epoch still uses all 24 new-family examples exactly once, with the remaining 20 before the unchanged old anchor.

## Measurements

For each subject × REPORT count:
- common-prefix old retention;
- terminal mean prefix damage;
- terminal mean anchor recovery;
- terminal mean tail damage;
- terminal mean net epoch change;
- final old held-out/unseen retention;
- final mean old retention;
- final primary/secondary new-family accuracy;
- final mean new-family accuracy.

## Interpretation

The exact five-point response is the result. A graded decrease in old retention or increasingly negative tail damage as REPORT count rises would establish a dose relationship. A step or non-monotonic response would locate a sharper composition boundary.

No post-result threshold is introduced.

## Bounds

No adaptive tail choice, no extra updates, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
