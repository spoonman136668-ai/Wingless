# Wingless UP-131B — few-shot lexical grounding efficiency

Status: preregistered scientific semantic-routing adaptation experiment.

Scientific parent: sealed UP-130B 9d58620a0871e67ecfcec6f2f2f59c1d3e0fcb80.

## Question

UP-130B showed that the exact UP-129B dual-view gate does not zero-shot generalize to entirely new event verbs. How many grounding subjects are required before the same frozen architecture/projector can learn a new 12-surface lexical family?

## Frozen mechanism

Start from the exact UP-129B gate:
- raw 64-D representation;
- frozen STORE-class competitor-nullspace projected candidate view;
- 128-D concatenated gate input;
- three-way STORE / OBSERVE / REPORT softmax;
- original 15-surface training completed first exactly as UP-129B;
- projector never recomputed.

New surfaces are the exact UP-130B set:
- STORE: lodges, stashes, caches, files
- OBSERVE: scans, checks, views, monitors
- REPORT: relays, announces, cites, summarizes

## Grounding arms

After original gate training, clone the gate and perform new-surface grounding for 20 epochs at lr 0.08 using:

1. zero_shot
   - no new-surface grounding.

2. one_subject
   - mia only.

3. two_subjects
   - mia, noah.

4. four_subjects
   - mia, noah, opal, pax.

Every grounding epoch also rehearses the original 15 surfaces on one fixed original subject (ada) to measure new lexical acquisition without silently discarding prior routing.

## Evaluation subjects

New lexical family:
- primary held-out: quin, rue
- secondary cross-subject: eli, fay, gia, hal, ivy, jon, kia, leo

Old lexical family:
- held-out original surfaces on eli/fay and prior unseen subjects.

## Metrics

For every arm:
- new-family overall accuracy;
- per-class accuracy;
- STORE precision/recall;
- per-surface accuracy;
- worst-surface accuracy;
- old-family retained class accuracy.

## Interpretation

A sharp improvement with one or two grounding subjects would indicate that the current gate needs lexical grounding but can acquire it sample-efficiently. Persistent failure at four subjects would implicate representation capacity rather than zero-shot semantics alone.

## Bounds

No projector recomputation, no threshold search, no architecture change, no surface-specific feature, no downstream replay change, no result-informed retry, no live activation, no production authority.
