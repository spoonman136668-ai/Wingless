# Wingless UP-102B — few-shot lexical grounding curve

Status: preregistered scientific lexical-learning experiment.

Scientific parent: sealed UP-101B `397840ebed9f170741980ee2d4c2f393bbe3a748`.

## Question

UP-101B showed that prefix timing preserves class information for known verbs but does not infer the semantics of unseen verbs. How many labeled examples of a new verb are required before the same frozen representation can ground that verb class and transfer it to different subjects?

## Frozen base training

Use the UP-101B prefix-through-verb representation and classifier.

Base training verbs:
- STORE: stores, keeps
- OBSERVE: observes, sees
- REPORT: reports, recalls

Base training subjects:
- ada, ben, cy, dee, eli, fay

Held-out verbs:
- STORE: holds
- OBSERVE: notes
- REPORT: tells

## Frozen grounding ladder

Grounding examples per held-out verb:
- 0
- 1
- 2
- 4

Grounding subjects are the first K subjects from:
- ada
- ben
- cy
- dee

Evaluation subjects for held-out verbs are always:
- eli
- fay

Thus held-out-verb evaluation subjects are never used as grounding examples.

For each K:
- start from an identical copy of the base classifier;
- train only on the K grounding examples for each held-out verb;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- deterministic class/subject order.

## Evaluation

For every K:
- held-out-verb overall accuracy on eli/fay;
- held-out STORE/OBSERVE/REPORT precision and recall;
- seen-verb accuracy after grounding to detect catastrophic forgetting.

## Interpretation

This is a sample-efficiency boundary, not a threshold search. Recovery at small K would support rapid lexical grounding on top of the fixed event representation. Persistent failure would move the problem upstream to representation learning.

## Bounds

No pretrained semantics, no encoder change, no threshold tuning, no class weighting, no attention, no memory-cap change, no result-informed retry, no live activation, no production authority.
