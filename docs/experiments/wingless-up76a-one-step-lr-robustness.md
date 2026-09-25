# Wingless UP-76A — one-step learning-rate robustness

Status: preregistered scientific robustness experiment.

Scientific parent: sealed UP-75A Windows evidence `8b0773a354f04bc73e6960a0c7ec197349ae46a1`.

## Question

UP-75A passed the unchanged 0.98 per-role held-out gate for both compact factorized representations after a single optimization step. UP-76A asks whether that one-step result is robust to the scale of the single update or is narrowly tied to learning rate 1.0.

## Frozen design

Unchanged from UP-75A:

- the same 54 residue-0 training states selected by the unchanged secondary coordinate;
- the same 162 held-out states;
- raw factorized one-hot representation, 15 features;
- role-factorized simplex representation, 10 features;
- five independent 3-class linear-softmax heads;
- zero initialization;
- exactly one optimization step;
- per-role held-out gate 0.98.

The only changed variable is the preregistered learning-rate ladder:

- 0.125;
- 0.25;
- 0.5;
- 1.0;
- 2.0.

No adaptive stopping, extra step, threshold change, extra training state, representation change, result-informed retry, promotion, or activation is permitted.

## Interpretation

Every learning-rate/representation point is sealed exactly as observed. Robust success, a narrow success band, or complete failure are all valid results. No learning rate is selected for production by this experiment.
