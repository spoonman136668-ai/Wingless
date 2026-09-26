# Wingless UP-170B — monotonic margin-rank risk

Status: preregistered scientific lexical-stability diagnostic.

Scientific parent: sealed UP-169B 938afd01c7da1ff384e1a3198c268df5e042271f.

## Question

Hard vulnerability cutoffs did not transfer perfectly in UP-169B. Does crossing risk nevertheless decrease monotonically as pre-step absolute-margin rank moves away from the fragile tail in a new fixed context?

## Holdout context

Use the same accepted subjects, matched pairs, report block, and update budget, but change the frozen context before measuring risk:

- terminal rehearsal phase = 17;
- the same 20 non-report new-family prefix examples are cyclically rotated by exactly 7 positions;
- no examples are added or removed;
- the same three pre-REPORT old updates are applied using phase-17 subject assignment;
- the same target-subject REPORT interference block is applied.

Rotation 7 is preregistered and not selected from results.

Build STORE→OBSERVE and OBSERVE→STORE cleanup path states from this frozen context.

## Independent interference probes

Probe all 15 old surfaces independently from fresh clones:
- five STORE;
- five OBSERVE;
- five REPORT.

Before each probe, rank all 120 old examples by absolute probability margin ascending.

## Fixed rank bins

Preregistered before this run:
- ranks 1..4;
- ranks 5..10;
- ranks 11..20;
- ranks 21..40;
- ranks 41..120.

For each interference class and rank bin, measure crossing density = observed boundary crossings / example-slots in that bin.

## Interpretation

A broadly decreasing crossing-density curve supports continuous margin rank as a risk signal even though hard class cutoffs are unstable. Non-monotonic reversals define the signal's boundary.

No threshold is fitted and no maintenance action is triggered.

## Bounds

Diagnostic only. Same update multiset and count, no adaptive rotation, no fitted threshold/classifier, no retained diagnostic updates, no extra training, no memory/architecture change, no result-informed retry, no live activation, no production authority.
