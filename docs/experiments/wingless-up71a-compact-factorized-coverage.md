# Wingless UP-71A — compact factorized training-coverage ladder

Status: preregistered scientific coverage/sample-efficiency experiment.

Scientific parent: UP-69A seal `d0b5ccd3566141b090b88c4a033f4c924a9a11b2`.

This experiment is defined independently of UP-70A outcomes.

## Question

UP-69A showed that compact 10-feature role-factorized simplex and 15-feature role-factorized one-hot observers both retain perfect held-out accuracy at the original 81-state training coverage.

UP-71A asks how far that result survives as the frozen training subset shrinks while the 162-state held-out set remains fixed.

## Frozen design

The fixed held-out set is unchanged:

- held out iff sum(role values) mod 3 != 0;
- 162 held-out states.

Training is selected only from the residue-0 pool using two fixed modular coordinates:

`secondary = (r0 + 2*r1 + r2 + 2*r3 + r4) mod 3`

`tertiary = (r0 + r1 + 2*r2 + 2*r3) mod 3`

The nested training levels are fixed before execution:

- 9 states: secondary == 0 and tertiary == 0;
- 18 states: secondary == 0 and tertiary in {0,1};
- 27 states: secondary == 0;
- 54 states: secondary in {0,1};
- 81 states: secondary in {0,1,2}.

At every level, two observer arms are trained independently:

- raw role-factorized one-hot, 15 features;
- role-factorized simplex, 10 features.

The five independent 3-class linear-softmax heads, zero initialization, 800 optimization steps, learning rate 1.0, and unchanged per-role held-out gate of 0.98 are preserved.

No extra training state, adaptive subset selection, nonlinear observer, threshold change, result-informed stopping, promotion, production authority, or activation is permitted.

## Interpretation

This is a fixed sample-efficiency ladder for already-qualified compact factorized representations. Any degradation, non-monotonicity, or class-specific failure is valid and must be sealed exactly as observed.
