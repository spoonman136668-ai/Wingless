# Wingless UP-112B — current-class base-anchor allocation

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-111B c9d4ba5a1a53b300c5fc435a41c3d965a4e1ac43.

## Question

UP-111B improved the remaining stage-3 base dip by excluding the current acquisition class from prior replay, while keeping later stages perfect. Can the last base error be removed at the same six-example replay budget by replacing some prior replay slots with extra base-lexicon anchors from the current class?

## Frozen system

Identical to UP-111B:
- deterministic 64-D prefix-through-verb encoder;
- same base lexicon and six-verb acquisition sequence;
- four grounding subjects for the current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- standard base replay remains one ada example for every base verb each epoch;
- exactly six additional replay examples per epoch whenever previous acquired verbs exist;
- replay occurs after current grounding and standard base replay.

## Frozen arms

1. current_class_excluded
   - exact UP-111B current_class_excluded control: all six additional slots go to non-current acquired classes.

2. base_anchor_2
   - two of the six additional slots are current-class base anchors:
     - one ben example for each of the two base verbs in the current class.
   - remaining four slots go to non-current acquired classes, distributed round-robin.

3. base_anchor_4
   - four of the six additional slots are current-class base anchors:
     - ben and cy examples for each of the two base verbs in the current class.
   - remaining two slots go to non-current acquired classes, distributed round-robin.

If no previous acquired verb exists, no six-slot additional replay block is used, matching prior lineage behavior.

Within non-current classes, verb choice follows acquisition order and subjects alternate ada/ben deterministically. No performance feedback is used.

## Evaluation

For every arm, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy;
- aggregate acquired-verb accuracy;
- per-acquired-verb accuracy.

## Interpretation

A base-anchor arm that reaches perfect base and acquired accuracy at every stage would close the fixed-six replay instability without increasing replay count. Loss of acquired verbs at larger base-anchor allocation would expose the retention tradeoff directly.

## Bounds

No encoder change, no additional replay slots, no adaptive replay, no threshold search, no class-weighted loss, no extra state, no result-informed retry, no live activation, no production authority.
