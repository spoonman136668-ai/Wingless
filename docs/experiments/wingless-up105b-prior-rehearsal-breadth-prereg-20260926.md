# Wingless UP-105B — prior-acquired rehearsal breadth

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-104B 9a54896c6903f5c06f86263ec050af53cb8985a6.

## Question

UP-104B showed backward interference: after sequentially grounding holds, notes, then tells, the earliest acquired verb (holds) was lost even though the base lexicon recovered. How much rehearsal breadth for previously acquired verbs is required to prevent that interference?

## Frozen base and acquisition

Identical to UP-104B:
- deterministic 64-D prefix-through-verb representation;
- base verbs: stores, keeps, observes, sees, reports, recalls;
- acquisition sequence: holds -> STORE, notes -> OBSERVE, tells -> REPORT;
- four grounding subjects for the current verb: ada, ben, cy, dee;
- evaluation subjects: eli, fay;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- after each epoch, base seen verbs rehearse exactly one fixed anchor subject: ada.

## Frozen arms

Previous-acquired rehearsal subjects per verb:
- 1: ada
- 2: ada, ben
- 4: ada, ben, cy, dee

For each arm:
- start from the identical base classifier;
- run all three acquisition batches in the same order;
- during a batch, previously acquired verbs rehearse on exactly the arm's fixed subject set after every epoch;
- the current batch still receives its normal four grounding subjects.

## Evaluation

After each acquisition stage:
- base seen-verb accuracy across all six subjects;
- accuracy for holds, notes, tells on eli/fay;
- aggregate accuracy over verbs acquired so far.

## Interpretation

This maps replay breadth rather than choosing a post-hoc winning threshold. If broader fixed replay preserves early verbs, the B104 boundary is rehearsal coverage. If not, the interference requires a different continual-learning mechanism.

## Bounds

No encoder change, no extra hidden state, no class weighting, no threshold search, no adaptive replay selection, no pretrained semantics, no result-informed retry, no live activation, no production authority.
