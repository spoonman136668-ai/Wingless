# Wingless UP-49A — role-binding composition

Status: Windows-qualified scientific result.

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


## Authoritative Windows result

Workflow run: `35987217037`

Runner: `WINGLESS-LINKDEADKB`

Artifact: `10801709799`

Artifact digest: `sha256:a1a85c30892c23b10fe1f27a8bb1381fb08ec2e1b0577e9b7df023dfe1e46496`

The harness, deterministic double probe, focused tests, full repository regression, and host-priority guard passed.

Results:

- unitary train accuracy: `1.0`
- unitary held-out accuracy: `1.0`
- unseen conjugated-right accuracy: `1.0`
- long mixed-program accuracy: `1.0`
- frozen role-binding gate: `PASS`
- matched non-unitary held-out accuracy: `0.20703125`

## Scientific classification

The learned unitary primitives composed into a right-role operation that was never directly trained, and the same system remained perfect on long noncommuting mixed programs.

The next cognition experiment should increase binding complexity rather than sequence length: multiple simultaneous roles and transfer of an operation across more than one role position.

## Plain speak

We taught it how to change the left variable and how to swap the two variables.

It correctly created the never-taught right-variable operation by composing those learned pieces, and it stayed perfect through long mixed programs.

That clears the first role-binding test cleanly.
