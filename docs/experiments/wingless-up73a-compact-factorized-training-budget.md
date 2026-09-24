# Wingless UP-73A — compact factorized training-budget robustness

Status: preregistered scientific optimization/robustness experiment.

Scientific parent: UP-71A seal `b9738d3d63e591d9d6194fdba43315344dc9dc88`.

This experiment is defined independently of UP-72A outcomes.

## Question

UP-71A found perfect held-out performance for both compact factorized representations at 54 training states using the existing 800-step learner.

UP-73A asks whether that 54-state result is robust to optimization budget rather than being specific to 800 training steps.

## Frozen design

The original residue-0 54-state training subset is unchanged:

- sum(role values) mod 3 == 0;
- UP-71A secondary coordinate < 2;
- 54 training states;
- the same fixed 162-state held-out set.

Representations:

- raw factorized one-hot, 15 features;
- role-factorized simplex, 10 features.

Training-step ladder, frozen before execution:

- 200;
- 400;
- 800;
- 1600.

All arms retain zero initialization, learning rate 1.0, five independent 3-class linear-softmax heads, and the unchanged per-role held-out gate of 0.98.

No adaptive stopping, learning-rate change, threshold change, extra training state, representation modification, promotion, production authority, or activation is permitted.

## Interpretation

This is an optimization-budget robustness ladder. Faster convergence, budget sensitivity, non-monotonicity, or failure are all valid results and must be sealed exactly as observed. No mechanism promotion is authorized.
