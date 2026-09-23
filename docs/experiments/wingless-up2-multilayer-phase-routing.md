# Wingless UP-2: multilayer phase routing

Status: Windows-qualified research primitive; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-1 Windows qualification sealed at `eb68ea7435e385b63c2f0bd427fd28476c9b7494`.

## Question

Can learned unitary propagation scale from one binary two-coordinate decision to a composed four-coordinate, four-class phase code while retaining exact bounded-state behavior?

A matched non-unitary control is trained beside it to prevent a misleading conclusion that merely solving the task demonstrates a unitary-specific advantage.

## Task

Each input is a normalized four-coordinate complex state. The four classes are Walsh sign patterns:

- class 0: + + + +
- class 1: + - + -
- class 2: + + - -
- class 3: + - - +

Every coordinate therefore has identical magnitude before propagation. Class information exists only in relative phase/sign relationships.

The train set contains all four classes at global phases 0 and pi/2. Held-out evaluation contains all four classes at unseen global phases 0.37, -1.11 and 2.22 radians.

## Unitary model

Four trainable complex Givens couplings are composed in two routing stages:

1. coordinates 0<->1 and 2<->3;
2. coordinates 0<->2 and 1<->3.

The model starts at identity and learns all four angles independently.

## Matched non-unitary control

The control receives:

- the same four-dimensional inputs;
- the same pair topology;
- exactly four trainable scalars;
- the same training examples;
- the same central-difference gradient estimator;
- the same learning rate;
- the same 60-step optimization budget.

Its pair update is an unconstrained mixer:

```
a' = a - g*b
b' = g*a + b
```

This control is capable of solving the routing task. It is not intentionally crippled. The difference under observation is the numerical/state behavior imposed by the unitary constraint.

## Pass conditions

- initial accuracy is exactly 0.25 for both paths;
- both paths reach 1.0 train accuracy;
- both paths reach 1.0 held-out accuracy;
- unitary final loss is <= 1e-6;
- each learned unitary angle is within 1e-3 radians of -pi/4;
- unitary maximum norm drift is <= 1e-12;
- unitary forward/inverse round-trip error is <= 1e-12;
- both control and unitary metrics remain finite;
- complete results are deterministic;
- existing Wingless regression remains green after advisory ICE rebuild.

## Interpretation boundary

A pass would demonstrate that the learned unitary primitive composes across multiple independently trained couplings and decodes a four-way relative-phase code on unseen global phases.

The matched control is expected to show whether ordinary unconstrained mixing can solve the same task under the same parameter/optimizer budget. No superiority claim is authorized merely from UP-2.

A positive result permits UP-3: repeated/deeper propagation with distractor dimensions, perturbation/noise tests, and measurements of information retention versus the matched non-unitary control.

## Authority boundary

UP-2 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-22 against source head `b07848b5ecbb2594f06646fe1487374f4e4326d1` passed the focused UP-2 suite and complete Wingless regression.

Observed results:

- both matched models: 1.0 train accuracy and 1.0 held-out accuracy;
- unitary final loss: 3.1754427471288606e-9;
- matched non-unitary final loss: 0.0006077264502525287;
- unitary maximum norm drift: 5.551115123125783e-16;
- matched non-unitary maximum norm drift: 2.735086374557516;
- unitary forward/inverse round-trip error: 2.603703785810335e-16.

Interpretation is deliberately bounded: both paths solved the routing task. The unitary path additionally preserved norm/reversibility and reached lower loss on this task under this optimizer. UP-2 alone does not establish a general unitary advantage.
