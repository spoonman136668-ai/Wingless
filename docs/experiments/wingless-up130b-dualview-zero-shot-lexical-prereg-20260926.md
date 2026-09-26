# Wingless UP-130B — dual-view zero-shot lexical breadth

Status: preregistered scientific semantic-routing generalization experiment.

Scientific parent: sealed UP-129B 3a69812d126a79dc509b2eaf5825173012472e35.

## Question

UP-129B achieved exact three-way routing and exact downstream behavior on all trained lexical surfaces and unseen subjects. Does the frozen 128-D dual-view gate generalize zero-shot to entirely unseen event verbs that were never present in gate training?

## Frozen gate and projector

Use the exact UP-129B learned gate:
- raw deterministic 64-D prefix-through-verb encoder;
- frozen STORE-class competitor-nullspace projector;
- gate input = raw 64-D + projected candidate 64-D;
- three-way STORE / OBSERVE / REPORT softmax;
- 20 training epochs at lr 0.08 on the original UP-128B 15-surface training set;
- no retraining on new surfaces;
- no threshold;
- no explicit class or surface-membership bit at inference.

## Zero-shot lexical families

Completely unseen surfaces:
- STORE: lodges, stashes, caches, files
- OBSERVE: scans, checks, views, monitors
- REPORT: relays, announces, cites, summarizes

These 12 surfaces are excluded from gate training.

## Subject families

Evaluate each surface on:
- original held-out: eli, fay
- prior unseen: gia, hal, ivy, jon, kia, leo
- new unseen: mia, noah, opal, pax, quin, rue

No new subject is used for gate training.

## Metrics

For each subject split and each surface:
- three-way class accuracy;
- STORE precision/recall;
- predicted class rates;
- mean class probabilities;
- raw-vs-projected STORE logit contribution.

Aggregate:
- accuracy per class;
- overall accuracy;
- worst-surface accuracy.

## Interpretation

Strong zero-shot accuracy would show that the learned dual-view gate has generalized beyond memorized verb identities. Failure localized to particular surfaces would identify the next representation/grounding boundary without changing the successful mechanism.

## Bounds

No gate retraining, no new lexical grounding, no projector recomputation, no threshold search, no surface-specific feature, no downstream replay change, no result-informed retry, no live activation, no production authority.
