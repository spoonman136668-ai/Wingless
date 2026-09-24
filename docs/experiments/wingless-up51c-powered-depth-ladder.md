# Wingless UP-51C — dimension-128 powered depth ladder

Status: preregistered scientific numerical-boundary diagnostic.

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
