# Wingless UP-50C — powered million-depth propagation

Status: preregistered scientific stress experiment.

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
