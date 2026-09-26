# Wingless UP-129B — generic dual-view three-way class gate

Status: preregistered scientific representation-routing experiment.

Scientific parent: sealed UP-128B 34240b7e3da4c1de6e8fb9011bb635c6aa52a8dc.

## Question

UP-128B removed the binary false-positive problem, but the raw three-way gate still misclassifies the surface "stores" as OBSERVE. Can a generic gate that sees both the raw representation and the candidate representation after the frozen STORE-class competitor-nullspace projection resolve that ambiguity without explicit STORE-surface membership?

## Frozen representation

- raw deterministic 64-D prefix-through-verb encoder;
- exact UP-124B frozen STORE-class competitor-nullspace projector;
- projector coefficients frozen before this experiment;
- projector is computed for every surface as a candidate view;
- no surface-specific correction;
- no adaptive geometry.

## Gate inputs

For every surface, concatenate:
1. raw 64-D representation;
2. projected 64-D candidate representation.

Total gate input dimension: 128.

The gate never receives the surface string, explicit class label, or a Boolean saying whether projection should apply at inference.

## Gate

- three-way STORE / OBSERVE / REPORT linear softmax;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- balanced exposure across the three classes;
- training subjects ada, ben, cy, dee;
- held-out eli, fay;
- unseen gia, hal, ivy, jon, kia, leo.

Training surfaces are the exact UP-128B surface set.

## Downstream use

- if gate argmax is STORE: downstream classifier receives the projected 64-D view;
- otherwise: downstream classifier receives the raw 64-D view.
- downstream training/replay/class-order matrix is frozen from UP-124B/UP-128B.

## Metrics

Gate:
- class accuracy;
- STORE precision/recall;
- per-surface confusion probabilities;
- raw-vs-projected logit contribution diagnostics.

Downstream:
- base and acquired accuracy;
- stores/archives/keepsafe margins;
- non-STORE corruption rate;
- all six class orders.

## Interpretation

If the dual-view gate resolves "stores" while retaining zero or negligible non-STORE projection errors, the learned router can use the frozen correction as evidence rather than requiring explicit surface membership.

## Bounds

No threshold search, no projector retraining, no surface-specific feature, no adaptive geometry, no downstream replay change, no result-informed retry, no live activation, no production authority.
