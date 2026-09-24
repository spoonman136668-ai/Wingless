# Wingless UP-48C — joint dimension/depth stress

Status: Windows-qualified scientific result.

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


## Authoritative Windows result

Workflow run: `35986325168`

Runner: `WINGLESS-UP-C`

Artifact: `10801934438`

Artifact digest: `sha256:aa0861f311c196c8329fe1537026969f4f12d69f46578b7a7293ff86a4ce35b5`

All four frozen unitary scale gates passed.

At depth 2048:

- dimension 16: accuracy `1`, norm drift `3.97e-13`, round trip `3.95e-13`
- dimension 32: accuracy `1`, norm drift `7.81e-13`, round trip `7.83e-13`
- dimension 64: accuracy `1`, norm drift `1.29e-12`, round trip `1.30e-12`
- dimension 128: accuracy `1`, norm drift `2.47e-12`, round trip `2.47e-12`

Pairwise-fidelity errors remained at roughly machine precision and perturbation gain remained effectively 1.

## Scientific classification

No tested dimension/depth boundary was reached. The unitary path retained classification, norm, relational geometry, perturbation magnitude, and reversibility across the full frozen sweep.

The next stress experiment should therefore increase scale aggressively rather than refine the already-green region.

## Plain speak

We made the substrate eight times wider than the original 16-dimensional stress case and ran it thousands of layers deep. It still behaved almost perfectly reversibly.

We did not find the wall. The next C-lane experiment needs to look much farther out.
