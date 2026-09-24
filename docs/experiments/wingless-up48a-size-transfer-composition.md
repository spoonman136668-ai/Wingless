# Wingless UP-48A — size-transfer compositional cognition

Status: preregistered scientific cognition experiment.

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
