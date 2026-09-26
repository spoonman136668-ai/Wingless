# Wingless UP-114B — acquisition-order generalization

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-113B bea5f6847499e5287788309bb66d85f73a6efbdd.

## Question

UP-113B closed the fixed-six replay instability by using exactly two current-class base anchors only at stage 3. Does that rule generalize when the semantic-class acquisition order changes, or is it specific to STORE -> OBSERVE -> REPORT?

## Frozen system

Identical to the successful UP-113B stage3_anchor2 arm:
- deterministic 64-D prefix-through-verb encoder;
- same six base verbs;
- same six acquired verbs;
- four grounding subjects for the current verb: ada, ben, cy, dee;
- evaluation on eli and fay;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- standard base replay remains one ada example for every base verb each epoch;
- exactly six additional replay examples per epoch whenever previous acquired verbs exist;
- only stage 3 uses the two-current-class-base-anchor allocation;
- every other stage uses current-class-excluded prior replay.

## Frozen acquisition orders

All three orders use the exact same six acquired verbs, each exactly once, and preserve the first/second verb identity within each semantic class.

1. STORE_OBSERVE_REPORT
   - holds, notes, tells, saves, watches, remembers.

2. OBSERVE_REPORT_STORE
   - notes, tells, holds, watches, remembers, saves.

3. REPORT_STORE_OBSERVE
   - tells, holds, notes, remembers, saves, watches.

Each first three stages therefore contains exactly one new verb from every semantic class, but the class completing stage 3 changes.

## Evaluation

For every order, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy across all six subjects;
- aggregate acquired-verb accuracy on eli/fay;
- per-acquired-verb accuracy.

## Interpretation

Perfect base and acquired accuracy at every stage across all three class orders would support the stage-3 anchor rule as an order-general mechanism. Any order-specific failure is sealed as observed and localizes the remaining dependence.

## Bounds

No encoder change, no replay-budget change, no adaptive replay, no threshold search, no class weighting, no extra state, no pretrained semantics, no result-informed retry, no live activation, no production authority.
