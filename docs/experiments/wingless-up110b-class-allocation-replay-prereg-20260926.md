# Wingless UP-110B — class-allocation replay

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-109B 734f49c0da903da1505353d44e41b9311d0fc550.

## Question

UP-109B showed that replay timing does not remove the stage-5 loss at a fixed six-example replay budget. Is the remaining interference caused by how replay examples are distributed across event classes?

## Frozen representation and acquisition

Identical to UP-109B:
- deterministic 64-D prefix-through-verb encoder;
- same base lexicon;
- same six-verb acquisition sequence;
- four grounding subjects for the current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- base replay remains one fixed subject (ada) for every base verb every epoch;
- exactly six previous-acquired replay examples per epoch;
- replay occurs after current grounding and base replay.

## Frozen arms

1. pair_major
   - exact UP-109B after arm / UP-107B pair-major fixed-six control.

2. equal_prior_class
   - replay class slots are STORE, STORE, OBSERVE, OBSERVE, REPORT, REPORT.
   - for each slot, choose a previously acquired verb of that class in deterministic acquisition-order round robin.
   - subject alternates deterministically between ada and ben by epoch plus slot.
   - if no prior verb exists for a requested class, fill that slot from all previous verbs in deterministic round robin.

3. exposure_balanced
   - when prior verbs exist in all three classes, allocate three replay examples to each class other than the current acquisition class and zero to the current class.
   - subjects and verbs rotate deterministically within each allocated class using acquisition order and ada/ben.
   - before all three classes are represented, fall back to equal_prior_class.
   - total prior replay remains exactly six.

The exposure_balanced arm compensates for the four current-class grounding examples and the fixed base replay rather than merely balancing prior replay internally.

## Evaluation

For every arm, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy;
- aggregate acquired-verb accuracy;
- per-acquired-verb accuracy.

## Interpretation

Improvement from class-aware allocation at identical replay count would identify class competition as the B109 boundary. Failure would move the mechanism toward representation-level interference between specific lexical surfaces.

## Bounds

No encoder change, no replay-budget change, no adaptive replay, no threshold search, no class weighting in the loss, no extra hidden state, no pretrained semantics, no result-informed retry, no live activation, no production authority.
