# Wingless UP-113B — stage-limited base-anchor allocation

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-112B 21c8645d492af224fdae7f734af62fce31af6f88.

## Question

UP-112B showed that two current-class base anchors remove the stage-3 base-lexicon error, but keeping those anchors at every later acquisition reintroduces acquired-verb loss. Is the useful effect specific to the early three-class boundary?

## Frozen system

Identical to UP-112B:
- deterministic 64-D prefix-through-verb encoder;
- same base lexicon and six-verb acquisition sequence;
- four grounding subjects for the current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- standard base replay: one ada example for every base verb each epoch;
- exactly six additional replay examples per epoch whenever previous acquired verbs exist;
- additional replay occurs after current grounding and standard base replay.

## Frozen arms

1. current_class_excluded
   - exact UP-111B control at every stage.

2. stage3_anchor2
   - use exact UP-112B base_anchor_2 allocation only while acquiring the third verb;
   - use current_class_excluded at all other stages.

3. stages2_3_anchor2
   - use base_anchor_2 while acquiring the second and third verbs;
   - use current_class_excluded at all other stages.

base_anchor_2 means:
- two of six additional slots are current-class base anchors: ben for each of the two base verbs in the current class;
- remaining four slots go to non-current acquired classes in the frozen deterministic order.

## Evaluation

For every arm, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy;
- aggregate acquired-verb accuracy;
- per-acquired-verb accuracy.

## Interpretation

Perfect base and acquired accuracy at every stage in a stage-limited arm would close the fixed-six instability without increasing replay cost. Failure would indicate that replay allocation must depend on a richer state than acquisition stage alone.

## Bounds

No replay-budget increase, no adaptive switching, no encoder change, no extra state, no threshold search, no class-weighted loss, no result-informed retry, no live activation, no production authority.
