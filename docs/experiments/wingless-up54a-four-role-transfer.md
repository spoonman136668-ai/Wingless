# Wingless UP-54A — four-role structural transfer

Status: preregistered scientific cognition experiment.

Scientific parent: UP-53A seal `b02660bf643d98f70468ed6296e3547be4ea422c`.

## Question

UP-53A showed that three informative examples are enough to identify and structurally transfer the learned three-role operators. The next end-goal question is whether the same shared reversible rules continue to compose when the bound state grows from three roles to four.

## Frozen design

- four ordered roles;
- three values per role;
- joint basis dimension 81;
- two learned scalar parameters exactly as in the three-role family: role-0 value swap and adjacent-role swap;
- training uses single-step supervision only;
- value swaps are trained only on role 0;
- adjacent-role swap is trained only at position 0/1;
- held-out tests include never-trained swap positions 1/2 and 2/3, derived mutations of roles 1, 2, and 3, and mixed programs of lengths 16 and 48;
- matched non-unitary control gets the same training set, optimizer, and budget;
- the existing cognition thresholds remain unchanged: >=0.99 train, >=0.98 held/category accuracy, <=1e-9 norm drift, <=1e-8 round-trip error.

Scientific negatives are valid. No threshold or optimizer change is permitted after results are observed.


## Authoritative Windows result

Workflow run: `36031780863`

Runner: `WINGLESS-LINKDEADKB`

Source head: `e23b7cab7653f4e596b2e1ba23c96d523540af94`

Artifact: `10823120568`

Artifact digest: `sha256:9b466069db29d1eb0e5f8dcc446e7425bcf75bd1d5cc5ccf2f88999355d802be`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Unitary results:

- train accuracy: `1.0`;
- aggregate held-out accuracy: `1.0`;
- unseen swap 1/2: `1.0`;
- unseen swap 2/3: `1.0`;
- derived role-1/2/3 mutation: all `1.0`;
- long-program accuracy at lengths 16 and 48: `1.0`;
- maximum norm drift: `1.5543122344752192e-15`;
- maximum round-trip error: `1.554723214883133e-15`;
- four-role transfer gate: PASS.

Matched non-unitary held-out accuracy: `0.4139433551198257`, with long-program accuracy `0.018518518518518517`.

## Scientific classification

The shared reversible primitives continue to transfer structurally when the bound state grows from three roles / dimension 27 to four roles / dimension 81. The model learned only role-0 value mutation and swap 0/1, yet generalized perfectly to never-trained adjacent swap positions, mutation of roles 1–3 by conjugation, and much longer mixed programs.

The next A-lane experiment should increase **computational structure**, not merely add another role: test context-dependent/controlled reversible operations and their composition.
