# Wingless UP-125B — learned raw-surface STORE projector gate

Status: preregistered scientific representation-routing experiment.

Scientific parent: sealed UP-124B adf0ead18133ee4390a309ed114440e27847ddd4.

## Question

UP-124B showed that one frozen STORE-class competitor-nullspace projector can replace per-verb corrections, but the experiment still explicitly knew which surfaces belonged to STORE. Can a learned gate decide from raw prefix-through-verb representation whether to apply the frozen projector, without explicit event-class labels at inference?

## Frozen representation intervention

- exact UP-124B class-wide STORE projector;
- projector coefficients frozen before this experiment;
- no recomputation or adaptive geometry;
- raw 64-D prefix-through-verb encoder unchanged.

## Learned gate

Binary gate: STORE-like vs non-STORE.

Training surfaces:
- STORE: stores, keeps, holds, saves, archives
- non-STORE: observes, sees, notes, watches, notices, reports, recalls, tells, remembers, recounts

Training subjects:
- ada, ben, cy, dee

Evaluation subjects:
- eli, fay
- unseen: gia, hal, ivy, jon, kia, leo

Gate:
- linear logistic readout on the raw 64-D encoding;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- threshold 0.5 frozen;
- no explicit class at inference.

## Downstream classifier stress

Use the exact UP-124B classifier/training/replay matrix.

Arms:
1. explicit_store_gate
2. learned_store_gate

When gate=true, apply the frozen classwide STORE projector before the downstream classifier.
When false, leave the representation unchanged.

## Metrics

Gate:
- accuracy, precision, recall on held-out and unseen subjects;
- per-surface accuracy.

Downstream:
- base/acquired accuracy;
- stores/archives/keepsafe margins;
- any non-STORE corruption introduced by false-positive projection;
- all six class orders.

## Interpretation

Matching explicit gating would replace hand-coded surface membership with learned raw-surface routing while preserving the successful classwide representation correction.

## Bounds

No gate threshold search, no projector retraining, no adaptive geometry, no downstream replay change, no state expansion, no result-informed retry, no live activation, no production authority.
