# Wingless UP-50C — powered million-depth propagation

Status: Windows-qualified scientific result.

Scientific parent: UP-49C seal `c725f5b83a0f38c74acc3afe602b5468c9200c57`.

## Question

Can the already-qualified unitary propagation block be evaluated orders of magnitude deeper with the repository's existing exponentiation-by-squaring matrix path without changing the mathematical operator, and do the frozen preservation gates still hold at depth 1,048,576?

## Frozen design

First, an equivalence gate:

- dimension 64;
- depth 2048;
- compare powered propagation against literal repeated propagation on all four prototypes and the eight screening noisy/clean states;
- maximum state L2 difference must be <= 1e-8.

Only if that equivalence gate passes may deep results be interpreted.

Deep powered cases:

- dimension 64, depth 1,048,576;
- dimension 128, depth 1,048,576.

The one-step matrix is constructed by applying the exact existing unitary block to every basis vector. Powers use the repository's existing `latentMatrixPower` exponentiation-by-squaring implementation.

The same UP-48C preservation thresholds are reused unchanged.

## Interpretation

- Overlap failure: treat powered propagation as an implementation mismatch; do not infer anything about the substrate at million depth.
- Overlap pass + deep gate failure: a numerical powered-depth boundary has been found and must be separated from mathematical/substrate failure.
- Overlap pass + deep gates pass: the tested unitary preservation properties survive a million repeated applications under an independently validated powered execution path.

## Plain speak

We first make the fast method prove it matches the ordinary method where we can still afford to run both.

Only then do we use it to jump from thousands of steps to more than a million.


## Authoritative Windows result

Workflow run: `35987745058`

Runner: `WINGLESS-UP-C`

Artifact: `10802623810`

Artifact digest: `sha256:992fd78acfd4f29c62b89dc349d8a0a3b8bdc8ee087b24323c1ea23b6332d3f6`

The powered path matched literal depth-2048 propagation with maximum state error `1.285e-13`, comfortably inside the preregistered overlap gate.

At depth 1,048,576:

- dimension 64: accuracy `1`, gate `PASS`, norm drift `7.18e-10`
- dimension 128: accuracy `1`, gate `FAIL`, norm drift `2.07e-9`

The dimension-128 failure is only in the extremely tight numerical preservation thresholds; classification remains perfect, Gram error remains `9.03e-11`, and round-trip error remains `2.12e-9`.

## Scientific classification

The powered implementation is validated against literal propagation at the overlap point.

A numerical preservation boundary appears by dimension 128 at one million applications under the unchanged `1e-9` norm/perturbation gate. This is not evidence of task-capability failure: accuracy remains 1 and the round-trip criterion still passes.

The next C-lane experiment should map the powered depth ladder at dimension 128 with unchanged thresholds to identify the first failing depth before changing precision or numerical stabilization.

## Plain speak

At a million steps, the system still gets the answer right.

What finally broke was our extremely strict numerical-cleanliness threshold at dimension 128. The next experiment finds exactly where that tiny floating-point drift crosses the line.
