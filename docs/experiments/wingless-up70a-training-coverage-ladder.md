# Wingless UP-70A — representation training-coverage ladder

Status: preregistered scientific coverage/robustness experiment.

Scientific parent: UP-68A seal `49f8deb62b063ca5fd8f059867c0e82d1a5e9313`.

This experiment is defined independently of UP-69A outcomes.

## Question

UP-68A showed that the 64-dimensional factorized representation generalizes perfectly across all three 81/162 residue partitions while the equally wide entangled Fourier control does not.

UP-70A asks how much of the original 81-state training residue is required before either representation generalizes to the same fixed 162-state held-out set.

## Frozen design

The fixed evaluation set is unchanged:

- held out iff sum(role values) mod 3 != 0;
- 162 held-out states.

Training always comes only from the original residue-0 pool. A second deterministic modular coordinate is frozen before execution:

`secondary = (r0 + 2*r1 + r2 + 2*r3 + r4) mod 3`.

The three nested training levels are:

- 27 states: secondary == 0;
- 54 states: secondary in {0,1};
- 81 states: secondary in {0,1,2}.

At each level, two independently trained observer arms are compared:

- dense Hadamard factorized, 64 features;
- joint Fourier entangled, 64 features / 32 frequency pairs.

The five independent 3-class linear-softmax heads, zero initialization, 800 optimization steps, learning rate 1.0, and unchanged per-role held-out gate of 0.98 are preserved.

No extra training state, adaptive subset selection, nonlinear observer, threshold change, result-informed stopping, promotion, production authority, or activation is permitted.

## Interpretation

This is a fixed coverage ladder. It measures sample efficiency under two unchanged representation classes. Any non-monotonicity or failure is valid and must be sealed exactly as observed.
