# Wingless UP-109B — replay timing at fixed budget

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-108B 095c70f4ad4cafa1d3a92e3ebbe9d0dded53583d.

## Question

UP-108B showed that simple verb-stratified coverage does not remove intermediate instability at a fixed six-example replay budget. Does the timing of the same replay examples relative to current-verb grounding control that instability?

## Frozen representation and acquisition

Identical to UP-108B:
- deterministic 64-D prefix-through-verb encoder;
- same base lexicon;
- same six-verb acquisition sequence;
- four grounding subjects for the current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- base replay remains one fixed subject (ada) for every base verb every epoch;
- prior-acquired replay budget remains exactly six examples per epoch.

## Frozen prior-replay example set

For each epoch, construct the exact pair-major fixed-six example sequence used by UP-107B:
- previous verbs in acquisition order;
- each paired with ada then ben;
- deterministic rotating offset based only on epoch;
- exactly six examples, wrapping if required.

Every timing arm uses exactly the same six prior examples for a given stage and epoch.

## Frozen arms

1. after
   - current grounding examples;
   - base replay;
   - all six prior replay examples.
   - exact pair-major fixed-six timing control.

2. before
   - all six prior replay examples;
   - current grounding examples;
   - base replay.

3. split_3_3
   - first three prior replay examples;
   - current grounding examples;
   - base replay;
   - final three prior replay examples.

## Evaluation

For every arm, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy;
- aggregate acquired-verb accuracy;
- per-acquired-verb accuracy.

## Interpretation

Differences between arms at identical update counts and identical replay examples isolate recency/timing effects. Persistent instability in every arm would move the mechanism beyond simple replay placement.

## Bounds

No encoder change, no replay-budget change, no adaptive replay, no threshold search, no class weighting, no extra state, no pretrained semantics, no result-informed retry, no live activation, no production authority.
