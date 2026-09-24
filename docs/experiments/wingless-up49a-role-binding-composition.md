# Wingless UP-49A — role-binding composition

Status: preregistered scientific cognition experiment.

Scientific parent: UP-48A seal `9d76277d68b392597ac3a4e6cad2352a3c15b849`.

## Question

UP-48A showed perfect transfer of learned local reversible rules across unseen state-space sizes and long sequences.

UP-49A changes the cognitive axis: can the substrate learn operations over a two-role bound state and use noncommuting composition to create an operation on a role that was never directly trained?

## Frozen design

The state is a 4 x 4 joint role/value space (dimension 16).

Only two primitives are trained, and only with single-step supervision:

1. `left_swap(start)`: swaps adjacent values in the left role, using one shared learned parameter across all four positions.
2. `swap_roles`: swaps the left and right roles, using one shared learned parameter across all off-diagonal role pairs.

There is no directly trained right-role primitive.

Held-out tests include:

- `swap_roles -> left_swap -> swap_roles`, which implements a right-role update only by composition;
- deterministic mixed programs of lengths 9, 27, and 81.

A matched non-unitary path receives the same two parameters, examples, optimizer, and training budget.

Frozen unitary diagnostic gate:

- train accuracy >= 0.99;
- aggregate held-out accuracy >= 0.98;
- conjugated-right accuracy >= 0.98;
- long-program accuracy >= 0.98;
- max norm drift <= 1e-9;
- max round-trip error <= 1e-8.

Scientific negatives are valid.

## Interpretation

A pass shows that reversible composition can move an operation between bound roles without directly training that derived operation.

A failure specifically on conjugated-right cases identifies role transport/binding as the next cognition bottleneck.

A failure only on long mixed programs identifies depth/composition rather than role binding.

## Plain speak

We teach the system how to change the left variable and how to swap the two variables.

We never teach it how to change the right variable.

Then we ask whether it can discover that changing the right variable is just: swap the roles, use the left rule, swap them back.
