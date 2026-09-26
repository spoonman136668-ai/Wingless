# Wingless UP-115B — base-lexicon error topology

Status: preregistered scientific continual lexical-learning diagnostic.

Scientific parent: sealed UP-114B 3c64686896dd22342a1d5f230c1292f6f9c91755.

## Question

UP-114B showed that the successful SOR stage-3 anchor rule does not generalize across semantic-class acquisition orders. Which specific base verbs are displaced at each failing stage, and do those errors track the current acquisition class or cross-class competition?

## Frozen training system

Reproduce the exact UP-114B training setup:
- deterministic 64-D prefix-through-verb encoder;
- same six base verbs;
- same six acquired verbs;
- all six semantic-class acquisition orders;
- four grounding subjects per current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- exactly six additional replay examples whenever previous acquired verbs exist.

Policies:
1. current_class_excluded
2. stage3_anchor2

No training behavior changes from UP-114B.

## Added diagnostics

At stage 0 and after every acquisition, report:
- aggregate base-lexicon accuracy;
- per-base-verb accuracy for:
  - stores
  - keeps
  - observes
  - sees
  - reports
  - recalls
- aggregate acquired-verb accuracy;
- per-acquired-verb accuracy;
- current acquisition semantic class.

Base-verb accuracy is evaluated across all six subjects.

## Interpretation

If base failures consistently localize to the current acquisition class, the next replay rule should protect same-class base surfaces. If failures occur in other classes, the mechanism is broader class competition or representation interference.

## Bounds

Diagnostic only: no replay-rule change, no budget increase, no encoder change, no extra state, no threshold tuning, no result-informed retry, no live activation, no production authority.
