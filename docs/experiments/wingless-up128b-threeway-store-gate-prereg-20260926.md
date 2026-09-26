# Wingless UP-128B — balanced three-way STORE gate

Status: preregistered scientific representation-routing experiment.

Scientific parent: sealed UP-127B d40cf39935ae917e3d1440032c0d5923c860c2c1.

## Question

UP-127B recovered STORE recall with balanced binary weighting, but introduced false STORE decisions concentrated in OBSERVE "observes" and REPORT "tells". Can a balanced three-way STORE/OBSERVE/REPORT classifier separate the competitor classes and provide a STORE projector gate without any binary threshold?

## Frozen encoder

- deterministic raw 64-D prefix-through-verb encoder;
- no representation projector used during gate training/evaluation.

## Three-way gate

Classes and surfaces:
- STORE: stores, keeps, holds, saves, archives
- OBSERVE: observes, sees, notes, watches, notices
- REPORT: reports, recalls, tells, remembers, recounts

Each class has exactly five surfaces.

Training subjects:
- ada, ben, cy, dee

Evaluation:
- heldout: eli, fay
- unseen: gia, hal, ivy, jon, kia, leo

Classifier:
- linear three-class softmax;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- deterministic per-surface interleaving STORE -> OBSERVE -> REPORT;
- no explicit event class at inference.

STORE gate is true iff STORE is argmax.

## Control

Report the UP-127B balanced binary gate metrics alongside the new three-way gate, but do not retrain or retune that control.

## Metrics

For each split:
- three-class accuracy;
- STORE-gate precision;
- STORE-gate recall.

For each surface:
- class accuracy;
- predicted-class distribution;
- mean STORE/OBSERVE/REPORT probabilities.

## Interpretation

If three-way argmax restores high STORE precision and recall simultaneously, the binary-gate overlap was caused by collapsing two distinct competitor classes. Remaining surface-local errors would identify the next representation boundary.

## Bounds

No threshold tuning, no encoder change, no projector use, no class weighting, no adaptive sampling, no result-informed retry, no live activation, no production authority.
