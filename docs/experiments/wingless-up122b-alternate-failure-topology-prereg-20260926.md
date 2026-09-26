# Wingless UP-122B — alternate lexical failure topology

Status: preregistered scientific diagnostic experiment.

Scientific parent: sealed UP-121B b887eb467212f16540ebc1af20feec74bfa263ec.

## Question

UP-121B showed that the frozen competitor-orthogonal stores correction remains perfect on unseen subjects, while the alternate six-verb acquisition set exhibits acquired-verb interference. Which alternate verbs/classes fail, and is the failure tied to lexical surface, unseen subject transfer, class order, or replay policy?

## Frozen training

Use the exact UP-121B competitor-orthogonal arm only:
- frozen stores correction from UP-120B;
- 64-D encoder;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- fixed replay budget 6;
- same six class orders;
- same current_class_excluded and stage3_anchor2 policies;
- no geometry or training changes.

Alternate acquisition verbs:
- STORE: keepsafe, archives
- OBSERVE: notices, spots
- REPORT: recounts, retells

Training subjects remain ada, ben, cy, dee.

## Evaluation

At every stage and for every class-order/policy cell, record each alternate acquired verb separately on:
- original held-out subjects: eli, fay
- unseen held-out subjects: kia, leo

Also record, for each verb and subject family:
- accuracy;
- mean correct-class minus strongest-competitor logit margin;
- minimum correct-class minus strongest-competitor logit margin;
- dominant wrong predicted class count.

Retain base stores/keeps and base-lexicon accuracy as controls.

## Interpretation

This is diagnostic only. It does not choose or tune a correction. A class-local failure would motivate a representation intervention for that class; a subject-specific failure would instead identify cross-subject generalization as the boundary.

## Bounds

No training change, no adaptive geometry, no replay change, no threshold tuning, no result-informed retry, no live activation, no production authority.
