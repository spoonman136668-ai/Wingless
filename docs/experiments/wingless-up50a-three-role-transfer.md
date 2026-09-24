# Wingless UP-50A — three-role transfer

Status: preregistered scientific cognition experiment.

Scientific parent: UP-49A seal `d5a709e13ead3a7c2fef4df913a1b3e3acd3139b`.

## Question

Can the learned reversible rules transfer across a third bound role, including a role-swap position and role-specific operations never shown during training?

## Frozen design

Three roles, three values each, joint state dimension 27.

Training is single-step only and contains:

- a value-swap primitive acting only on role 0;
- a role-swap primitive trained only at position 0 (roles 0<->1).

The role-swap parameter is shared across positions, but position 1 (roles 1<->2) is never trained.

Held-out tests include:

- direct use of the never-trained position-1 role swap;
- derived mutation of role 1 by conjugation;
- derived mutation of role 2 by moving it across two role positions, mutating it with the role-0 rule, then restoring ordering;
- long mixed programs of lengths 12 and 36.

A matched non-unitary path receives the same parameters, training set, optimizer, and budget.

The frozen unitary gate requires >=0.98 on aggregate held-out, unseen swap12, derived role1, derived role2, and long-program categories, plus the existing norm and round-trip bounds.

Scientific negatives are valid.

## Plain speak

The previous test moved a learned rule from left to right.

This one adds a third variable. We never teach the system the second role-swap position or how to mutate roles 1 or 2. It has to reuse and compose what it learned to make those operations work.
