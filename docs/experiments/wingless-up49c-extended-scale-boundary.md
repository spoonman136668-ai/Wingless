# Wingless UP-49C — extended scale boundary screen

Status: preregistered exploratory stress screen.

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
