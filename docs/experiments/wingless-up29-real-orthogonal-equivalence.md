# Wingless UP-29: exact real orthogonal equivalence

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-28 qualified scientific pass sealed at `f0c4f92fda1d8e6ef44273edf061fbb3b11beac9`.

## Question

UP-28 closed the fully mixed, dynamics-discovered observer pipeline with perfect held-out and mutable performance.

Before making any stronger statement about complex phase or unitary computation, UP-29 performs the required real orthogonal comparator:

**Is the accepted 96-dimensional complex unitary runtime exactly reproducible as an ordinary 192-dimensional real orthogonal system with the same number of real scalar degrees of freedom?**

## Realification

A complex vector

`z = x + i y`

is represented as interleaved real coordinates:

`[x0, y0, x1, y1, ...]`.

Each complex matrix entry

`a + i b`

is translated to the real 2 × 2 block:

`[[a, -b], [b, a]]`.

Therefore a 96-dimensional complex state becomes a 192-dimensional real state with exactly the same 192 real scalar degrees of freedom.

A complex unitary matrix becomes a real orthogonal matrix under this mapping.

## Frozen accepted model

UP-29 does not search for a new model.

It freezes:

- the accepted UP-28 64-observable bank;
- the same depth-0 training data;
- the same ridge/phase-code decoder;
- the same held-out tables/depths;
- the same mutable workload.

The 64 observable indices are frozen from the accepted UP-28 result.

## Runtime real arithmetic

For the comparator arm:

- recurrent state evolution uses 192-dimensional real matrices/vectors;
- norm calculation is real;
- complex observable expectations are reproduced from real and imaginary coefficient arrays using only real arithmetic;
- the same quadratic feature map is then applied;
- the same accepted decoder weights are used.

No complex arithmetic is used by the real comparator after the complex state/transport/observable structures are translated.

## Explicit limitation

The real encoder is **not independently learned**.

The established complex encoder produces the accepted latent state and that state is realified.

Likewise the real transport matrices and observable coefficient arrays are exact translations of the accepted complex structures.

This experiment tests mathematical/runtime representational equivalence, not whether an independently trained real architecture rediscovers the same solution.

## Scientific gates

**Complex baseline reproduction**

- accepted frozen complex held-out accuracy = 1.0.

**Real orthogonality**

- maximum entry of `R(V)^T R(V) - I` <= 1e-10.

**State equivalence**

- maximum L2 error between complex evolution and complexified real evolution <= 1e-9.

**Feature equivalence**

- maximum absolute error between complex and real quadratic features <= 1e-8.

**Decision equivalence**

- real held-out accuracy equals complex held-out accuracy;
- zero class-decision disagreements.

**Real mutable integration**

- commit decode = 1.0;
- exact final table = 1.0;
- relational query = 1.0;
- maximum real norm drift <= 1e-10.

## Interpretation boundary

A positive result means the accepted complex/unitary runtime is exactly representable as a classical real orthogonal system at equal real-scalar degrees of freedom.

That would rule out any claim that the observed capability requires uniquely complex-number or quantum computation. The useful property would instead be consistent with reversible norm-preserving rotational dynamics plus the learned invariant observer structure.

A negative result would indicate an implementation mismatch in the proposed realification unless a specific mathematical nonequivalence is demonstrated; it would not by itself establish a complex advantage.

A later experiment may train a real orthogonal architecture independently, but that is a different question from exact representational equivalence.

## Plain speak

Complex numbers package two real numbers together.

UP-29 takes Wingless's accepted 96 complex coordinates and unpacks them into 192 ordinary real coordinates.

It then rewrites every rotation and every measurement in ordinary real arithmetic.

If the behavior is identical, the important discovery is not “quantum magic” or complex numbers themselves. It is the stable reversible geometry and the way the system learned to measure that geometry.

## Authority boundary

UP-29 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
