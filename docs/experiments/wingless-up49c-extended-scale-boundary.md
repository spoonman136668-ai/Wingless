# Wingless UP-49C — extended scale boundary screen

Status: Windows-qualified scientific result.

Scientific parent: UP-48C seal `d28cf67488c33509c7f42c6f523755e00031c467`.

## Question

UP-48C did not find a preservation boundary through dimension 128 at depth 2048.

UP-49C asks whether the same unitary stability gates survive a much harder diagonal width/depth envelope.

## Frozen screening cases

- dimension 128, depth 8192
- dimension 256, depth 4096
- dimension 512, depth 2048
- dimension 1024, depth 1024

This is a screening run: two deterministic perturbation trials per class and the unitary path only.

No matched non-unitary control is run in this screen because UP-48C already established the matched-control contrast in the lower envelope. If any case fails, the first failing case must be repeated in the next experiment with the full 32-sample battery and matched control before interpreting the boundary.

The exact UP-48C unitary gate is reused unchanged.

## Interpretation

- Any failure: freeze the first failing case and run a full confirmation at that exact point.
- All pass: no boundary found; move the stress lane toward either still larger sparse dimensions or much longer depth using a performance-optimized harness.

## Plain speak

This run is a fast scouting mission.

Instead of spending a long time proving every point we already expect to work, it jumps much farther out. If something cracks, the next experiment slows down and examines that exact crack carefully.


## Authoritative Windows result

Workflow run: `35987022666`

Runner: `WINGLESS-UP-C`

Artifact: `10802498095`

Artifact digest: `sha256:e65adf84a04fac773c4383bf1d0ad2e4312d9fd90604ad21c2386e65228cb7c0`

All four screening gates passed. No first failing case was found.

- dim 128 / depth 8192: norm drift `9.88e-12`, round trip `9.88e-12`
- dim 256 / depth 4096: norm drift `1.05e-12`, round trip `1.05e-12`
- dim 512 / depth 2048: norm drift `6.38e-13`, round trip `6.37e-13`
- dim 1024 / depth 1024: norm drift `1.50e-12`, round trip `1.50e-12`

Accuracy remained 1 and fidelity-geometry error stayed near machine precision in every case.

## Scientific classification

The expanded width/depth screen still did not expose a unitary preservation boundary.

The next C-lane question is whether a mathematically equivalent powered propagation method can reproduce iterative propagation at an overlapping depth and then probe orders-of-magnitude deeper without changing the underlying unitary block.

## Plain speak

We jumped much farther out and still did not find a break.

The next efficient move is not to make the runner repeat the same block millions of times literally. We first prove a faster powered method matches the ordinary method, then use it to look vastly deeper.
