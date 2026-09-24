# Wingless UP-48C — joint dimension/depth stress

Status: preregistered scientific stress test.

Scientific parent: Windows-qualified UP-47 source head `474ebb42082d1a9dd022517edb3c48c4f16dbd1d`.

## Question

Do the unitary substrate's strongest stability properties survive when state dimension and propagation depth increase together?

## Frozen sweep

Dimensions: 16, 32, 64, 128.

Depths: 0, 128, 512, 2048 repeated butterfly blocks.

Each dimension uses four deterministic signal classes, dimension-minus-four distractor coordinates, 32 perturbed samples, and a matched non-unitary control with the same topology and coupling values.

The unitary path measures classification, norm drift, pairwise-fidelity geometry drift, perturbation gain, and maximum-depth round-trip error.

Frozen unitary scale gate per dimension:

- accuracy = 1 at every depth;
- max norm drift <= 1e-9;
- max Gram/fidelity error <= 1e-9;
- perturbation gain within 1e-9 of 1;
- maximum-depth round-trip error <= 1e-8.

The gate is diagnostic only. A scientific negative does not fail the harness.

## Interpretation

A pass across all dimensions would show that the basic preservation properties survive a much larger joint dimension/depth envelope than UP-3.

A failure identifies the first tested scale where numerical or architectural preservation begins to break and becomes the boundary for the next stress experiment.

## Plain speak

Instead of making the task cleverer, this lane makes the same core substrate much bigger and much deeper and tries to break it.

If it stays stable, we can stop worrying as much about small-system-only behavior. If it breaks, we learn exactly where.
