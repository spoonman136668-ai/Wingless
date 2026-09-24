# Wingless UP-48A — size-transfer compositional cognition

Status: Windows-qualified scientific result.

Scientific parent: Windows-qualified UP-47 source head `474ebb42082d1a9dd022517edb3c48c4f16dbd1d`.

## Question

Can a tiny learned reversible program generalize beyond the state-space sizes seen during training and compose correctly over much longer unseen instruction sequences?

## Frozen design

Two shared local operator parameters are learned from single-instruction supervision only.

Training state-space dimensions: 4 and 6.

Held-out state-space dimensions: 8 and 12.

The two instruction families act on local spans 1 and 2. The same learned parameter for each instruction family is reused at every position and every dimension.

Training contains no multi-instruction programs.

Held-out programs have lengths 4, 16, 64, and 128 and use deterministic unseen instruction sequences.

A matched non-unitary path gets the same two trainable scalars, training samples, finite-difference optimizer, and optimization budget.

Frozen diagnostic gate for the unitary path:

- training accuracy >= 0.99;
- aggregate held-out accuracy >= 0.98;
- each unseen dimension accuracy >= 0.95;
- max norm drift <= 1e-9;
- max round-trip error <= 1e-8.

A scientific negative does not fail the harness.

## Interpretation

A pass would show algorithmic transfer across both composition length and state-space size: the learned local rule would not be tied to a fixed four- or six-state training substrate.

A failure with good training accuracy but weak unseen-dimension accuracy would identify representation/topology transfer as the cognition bottleneck.

A failure already on long programs would keep the bottleneck at compositional stability rather than abstraction across size.

## Plain speak

The system learns two tiny local rules on small worlds.

Then we drop the same learned rules into bigger worlds it never trained on and ask them to execute long programs correctly.

That is a stronger test of reusable cognition than memorizing one fixed-size task.


## Authoritative Windows result

Workflow run: `35986520028`

Runner: `WINGLESS-LINKDEADKB`

Artifact: `10802437520`

Artifact digest: `sha256:64dc19a732b7a3ea74f4a6d5edbd75cf4e5e62163f7c4f4592bb482857357ba5`

The harness, deterministic double probe, focused tests, full repository regression, and host-priority guard passed.

Results:

- unitary training accuracy: `1.0`
- unitary held-out accuracy: `1.0`
- minimum unseen-dimension accuracy: `1.0`
- frozen algorithmic-transfer gate: `PASS`
- matched non-unitary held-out accuracy: `0.675`
- unitary held-out advantage: `+0.325`

## Scientific classification

The learned reversible local program generalized perfectly from training dimensions 4/6 to unseen dimensions 8/12 and composed correctly over unseen programs up to length 128.

This answers the initial size-transfer question positively. The next cognition experiment should move to a different axis rather than extending sequence length again: role/binding composition with noncommuting operations and held-out combinations.

## Plain speak

The learned rule was not stuck to the size of the world it trained in.

We trained it only on tiny one-step examples in small state spaces, then moved it into larger spaces and ran long programs. It stayed perfect.

The next useful question is whether it can bind and manipulate multiple roles/variables, not whether we can make the same sequence longer.
