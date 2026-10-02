# Wingless UP-50A — three-role transfer

Status: Windows-qualified scientific result.

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


## Authoritative Windows result

Workflow run: `35987591398`

Runner: `WINGLESS-LINKDEADKB`

Artifact: `10803020811`

Artifact digest: `sha256:1b8fc0e4f2df05f8cff5e7ba60ee8eccde2e025233de700ec1ce096ccdd6ffc9`

Results:

- unitary train accuracy: `1.0`
- aggregate unitary held-out accuracy: `1.0`
- never-trained swap(1,2): `1.0`
- derived role-1 mutation: `1.0`
- derived role-2 mutation: `1.0`
- long mixed programs: `1.0`
- frozen three-role transfer gate: `PASS`
- matched non-unitary held-out accuracy: `0.4074074074074074`

## Scientific classification

The learned unitary rules transferred across an untrained role-swap position and supported multi-hop transport of a learned operation across three bound roles without direct training of those derived primitives.

The next A-lane experiment should reduce supervision rather than add a fourth role: train the same structural rules on a sparse subset of basis configurations, then test held-out basis states and coherent superpositions to determine whether the rule is genuinely structural rather than dependent on exhaustive basis coverage.

## Plain speak

The system did not need us to teach every role position separately.

It reused one learned swap rule in a new position and carried a learned value operation across two role boundaries perfectly.

The next harder question is whether it can do that without seeing every possible training state first.
