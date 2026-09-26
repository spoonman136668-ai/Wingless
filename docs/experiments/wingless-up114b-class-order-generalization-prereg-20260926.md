# Wingless UP-114B — stage-3 anchor class-order generalization

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-113B bea5f6847499e5287788309bb66d85f73a6efbdd.

## Question

UP-113B achieved perfect base and acquired accuracy at every stage by using two current-class base anchors only at the third acquisition. Does that rule generalize across all semantic-class acquisition orders, or is it specific to STORE -> OBSERVE -> REPORT?

## Frozen representation and learning

Identical to UP-113B:
- deterministic 64-D prefix-through-verb encoder;
- same six base verbs;
- same six acquired verbs, two per semantic class;
- four grounding subjects per current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- standard base replay remains one ada example for every base verb every epoch;
- exactly six additional replay examples per epoch whenever previous acquired verbs exist.

Acquired verb pairs:
- STORE: holds, saves
- OBSERVE: notes, watches
- REPORT: tells, remembers

## Frozen class orders

Evaluate all six permutations:
- SOR
- SRO
- OSR
- ORS
- RSO
- ROS

For each class order:
- stages 1-3 acquire the first verb from each class in that order;
- stages 4-6 acquire the second verb from each class in the same order.

## Frozen policies

For every class order:

1. current_class_excluded
   - exact UP-111B control at every stage.

2. stage3_anchor2
   - exact UP-113B successful rule:
   - at acquisition stage 3 only, two of six additional replay slots are current-class base anchors (ben for each base verb of the current class);
   - remaining four slots use current-class-excluded prior replay;
   - all other stages use current_class_excluded.

## Evaluation

For every order x policy at stage 0 and after every acquisition:
- base-lexicon accuracy across all six subjects;
- aggregate accuracy over the verbs acquired so far on eli/fay;
- per-acquired-verb accuracy.

## Interpretation

If stage3_anchor2 remains perfect across all six class orders while matched controls show order-dependent dips, the rule generalizes to the first completion of semantic-class coverage rather than a specific lexical sequence.

## Bounds

No replay-budget increase, no adaptive rule selection, no encoder change, no extra state, no threshold tuning, no class weighting, no result-informed retry, no live activation, no production authority.
