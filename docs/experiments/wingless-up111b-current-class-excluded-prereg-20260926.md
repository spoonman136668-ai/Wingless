# Wingless UP-111B — current-class-excluded replay

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-110B 892e78c7ea1e97b2a7c1cd7e635cfa47d049d471.

## Question

UP-110B removed the stage-5 acquired-verb loss by allocating replay away from the current acquisition class, but retained the stage-3 base dip because the arm fell back before all three prior classes existed. Does applying current-class exclusion whenever any non-current prior class exists remove the remaining instability at the same six-example budget?

## Frozen representation and acquisition

Identical to UP-110B:
- deterministic 64-D prefix-through-verb encoder;
- same base lexicon and six-verb acquisition sequence;
- four grounding subjects for the current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- base replay: one fixed ada example for every base verb after every epoch;
- previous-acquired replay: exactly six examples per epoch whenever previous verbs exist;
- replay occurs after current grounding and base replay.

## Frozen arms

1. pair_major
   - exact pair-major fixed-six control.

2. exposure_balanced_control
   - exact UP-110B exposure_balanced arm, including its fallback before all three prior classes exist.

3. current_class_excluded
   - identify prior classes different from the current acquisition class that have at least one acquired verb;
   - distribute all six replay slots round-robin across those eligible non-current classes;
   - within each class, rotate verbs in acquisition order and alternate subjects ada/ben deterministically by epoch and slot;
   - if only one non-current prior class exists, all six examples come from that class;
   - if no previous verb exists outside the current class, fall back to deterministic pair-major replay over all previous verbs;
   - never use performance feedback.

## Evaluation

For every arm, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy;
- aggregate acquired-verb accuracy;
- per-acquired-verb accuracy.

## Interpretation

Perfect base and acquired accuracy at every stage would show that the remaining fixed-budget instability was class-exposure competition. Residual loss would move the B lane toward lexical-surface-specific interference or parameter protection.

## Bounds

No encoder change, no replay-budget change, no adaptive replay, no threshold search, no class-weighted loss, no extra state, no pretrained semantics, no result-informed retry, no live activation, no production authority.
