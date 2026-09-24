# Wingless UP-52A — basis coverage ladder

Status: preregistered scientific cognition experiment.

Scientific parent: UP-51A seal `d7eef7212ad0e309246fcfc85bfd5c255e97e034`.

## Question

UP-51A learned the three-role reversible operators from only 9 of 27 basis configurations and generalized perfectly to the other 18.

UP-52A asks how far training-basis coverage can be reduced before that structural generalization fails.

## Frozen design

The UP-51A architecture, optimizer, primitive set, long-program lengths, matched non-unitary control, numerical gates, and accuracy gates are unchanged.

Four independent training conditions are evaluated from the same initialization:

- 9 observed basis states;
- 6 observed basis states;
- 3 observed basis states;
- 1 observed basis state.

The deterministic training-basis order is frozen before execution:

`000, 111, 222, 012, 120, 201, 021, 102, 210`.

Each condition trains only single-step role-0 value swaps and role swap 0/1 on its selected basis states. Every remaining basis configuration is held out completely.

Held-out evaluation retains the UP-51A categories: unseen primitives, unseen swap 1/2, derived role-1 mutation, derived role-2 mutation, and mixed programs of lengths 12 and 36.

No threshold changes are permitted after observing results.

## Interpretation

The minimum passing basis count is the smallest tested supervision set that still satisfies the unchanged UP-51A gate.

A failure is a valid scientific boundary, not a reason to change optimizer settings or gates.


## Authoritative Windows result

Workflow run: `36030288060`

Runner: `WINGLESS-LINKDEADKB`

Source head: `50f5335b1dadc4f47f0ab06f3876651490a719d7`

Artifact: `10820679947`

Artifact digest: `sha256:98518d03bafa54ab63b5e522fe0805cb4b9307e13591d760241e67072386cfc1`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Results:

- 9 basis states: gate PASS; held-out 1.0; primitive 1.0; derived role-1 1.0; derived role-2 1.0; long-program 1.0.
- 6 basis states: gate PASS; held-out 1.0; primitive 1.0; derived role-1 1.0; derived role-2 1.0; long-program 1.0.
- 3 basis states: gate FAIL; held-out 0.28125; primitive 1.0; derived role-1 0.08333333333333333; derived role-2 0.08333333333333333; long-program 0.125.
- 1 basis state: gate FAIL; held-out 0.30288461538461536; primitive 1.0; derived role-1 0.10256410256410256; derived role-2 0.10256410256410256; long-program 0.15384615384615385.

The smallest passing condition in this frozen ladder is 6 observed basis states.

## Scientific classification

The raw coverage boundary is between 3 and 6 states, but the 3-state subset `000/111/222` is invariant under the trained role-swap operation. The role-swap parameter remained at its 0.1 initialization in the failing low-coverage condition while value-swap transfer stayed perfect.

Therefore the immediate causal question is **identifiability versus coverage**. The next A-lane experiment must hold training count at three and compare the invariant diagonal triad against preregistered role-swap-informative triads. No optimizer or gate changes are permitted.
