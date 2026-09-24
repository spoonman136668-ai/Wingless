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
