# Wingless UP-LM2V — response generalization

Status: preregistered scientific counterfactual response-generalization experiment.

Scientific parent: sealed UP-LM2U e5a45d44e7be97dabfc21700bdcc95c567dc57e1.

## Question

Does the frozen online warning-guided early-closure rule still prevent dependency loss under disjoint identity rotations and different pending-dependency layouts?

## Frozen capacity and response rule

- exact recall cap: 16;
- 12 original dependencies followed by 12 unique pressure writes;
- no extra memory or training;
- before each pressure write:
  - if the current FIFO victim is still unreported, guided arm closes that exact dependency;
  - sham arm closes a different unreported dependency at the same opportunity;
  - baseline does nothing.

The rule reads only current recall order and current reported/unreported state.

## Disjoint controls

Deferred levels: 4,5,6.

Identity rotations not used in LM2S-U:
- 3
- 11

Value shifts:
- 0
- 2

Pending layouts:
- suffix_reported: the newest 12-d originals are already reported, leaving the oldest d pending;
- alternating_reported: a fixed alternating permutation determines which 12-d originals are already reported.

These layouts alter which original dependency is at risk without changing capacity.

## Measurements

Per arm:
- baseline/guided/sham completed originals;
- baseline/guided/sham failed originals;
- guided/sham actions.

Aggregate:
- failures prevented by guided response;
- failures prevented by sham;
- unnecessary guided actions in safe states.

## Interpretation

Generalization requires warning-guided benefit to survive the new identity/layout conditions and remain target-specific relative to sham.

## Bounds

Counterfactual only. No live activation, future schedule input, capacity growth, extra training, semantic oracle, or post-result policy change.
