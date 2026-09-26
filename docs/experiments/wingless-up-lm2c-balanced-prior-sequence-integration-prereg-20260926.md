# Wingless UP-LM2C — balanced-prior sequence integration

Status: preregistered scientific five-family language integration experiment.

Scientific parent: sealed UP-LM2B b90ad4353a318b49ea8e446da8cbad7997cf87b1.

## Question

UP-LM2B showed that every prior family benefits from later update position while the fifth family retains a stable terminal-position advantage. Does position-balanced prior rotation preserve or improve family balance across the accepted sequence-order transformations, rather than only on block-heldout evaluation?

## Frozen starting point and budget

Reuse the exact accepted five-family starting state and router:
- 64-D recurrent transition;
- recurrent parameters frozen;
- accepted four-family rotating-palindromic prehistory;
- fifth-family alphabet extension;
- exact recall cap 16;
- total adaptation mass exactly 0.40 per matched example;
- four adaptation epochs;
- fifth family always updated last.

Allocation arms remain the five accepted fixed-mass settings from UP-LM2B.

## Training schedules

1. canonical_prior
   - base → paraphrase → third → fourth → fifth.

2. balanced_prior
   - prior start position = (epoch + matched-example index) mod 4;
   - all four prior families updated once in cyclic order;
   - fifth always last.

No schedule is selected adaptively.

## Evaluation sequence orders

Evaluate all five accepted sequence organizations:
- block;
- per_name;
- paired_names;
- stores_then_local_reports;
- reverse_report_tail.

For every allocation × training schedule × evaluation order × family, evaluate stream depths 1 and 4.

## Measurements

Per cell:
- integrated top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- query-set exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Per allocation × training schedule × evaluation order:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

If balanced prior rotation preserves structural exactness and narrows prior-family imbalance across evaluation orders without disrupting fifth terminal performance, recency equalization survives integrated sequence organization. If sequence order interacts strongly with the balanced schedule, that interaction becomes the next language boundary.

No schedule or allocation is declared a winner; the full factorial response is the result.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
