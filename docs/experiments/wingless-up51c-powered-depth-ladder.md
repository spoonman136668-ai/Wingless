# Wingless UP-51C — dimension-128 powered depth ladder

Status: Windows-qualified scientific result.

Scientific parent: UP-50C seal `7eb1bfb8c13c139e732a72c08763e5a1043bbc4e`.

## Question

UP-50C validated powered propagation against literal propagation, then found the first tested frozen-gate failure at dimension 128 and depth 1,048,576 while task accuracy remained perfect.

UP-51C maps the depth ladder at dimension 128 to locate the first gate crossing without changing any threshold.

## Frozen depths

8192, 32768, 131072, 524288, 1048576.

The exact UP-50C powered metric and the unchanged UP-48C numerical gate are reused.

No precision change, renormalization, reunitarization, or threshold adjustment is allowed.

## Interpretation

The first failing depth becomes the numerical boundary to diagnose next. This remains a numerical implementation result unless task capability also degrades.

## Plain speak

We know one million steps is just over our extremely strict numerical line at dimension 128.

This run fills in the gap and tells us where that line is actually crossed.


## Authoritative Windows result

Workflow run: `35988082984`

Runner: `WINGLESS-UP-C`

Artifact: `10802034889`

Artifact digest: `sha256:51b66e01b6a6f41be99b8a0940f7e165de592626bebda7b5b1baaa7068b1967b`

Results at dimension 128:

- depth 8,192: gate PASS, norm drift `1.62e-11`
- depth 32,768: gate PASS, norm drift `6.48e-11`
- depth 131,072: gate PASS, norm drift `2.59e-10`
- depth 524,288: gate FAIL, norm drift `1.04e-9`
- depth 1,048,576: gate FAIL, norm drift `2.07e-9`

Last passing depth: `131072`.

First failing depth: `524288`.

## Scientific classification

The first tested crossing of the unchanged numerical preservation gate at dimension 128 lies between depth 131,072 and 524,288. The measured norm error grows approximately with propagation depth in this region, while prior UP-50C evidence showed task classification remained perfect at one million steps.

This closes the immediate numerical-boundary question sufficiently for the broader program. The C lane should now spend its experimental budget on a more consequential scaling risk: memory/load interference rather than narrowing a floating-point threshold that is already orders of magnitude beyond ordinary tested depth.

## Plain speak

We found the line.

The system still solves the task at a million steps, but our very strict numerical-cleanliness rule starts failing somewhere after 131 thousand steps and by 524 thousand steps.

That is enough to know where the floating-point boundary lives. The more important next scale question is how much information the substrate can hold before memories interfere with one another.
