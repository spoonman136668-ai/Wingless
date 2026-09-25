# Wingless UP-77A — one-step compact-factorized coverage ladder

Status: preregistered scientific data-efficiency experiment.

Scientific parent: sealed UP-76A Windows evidence `5616cd22352031edd2da1c37a3cac81fd8cb4574`.

## Question

UP-76A showed perfect held-out performance for both compact factorized representations after exactly one optimization step across learning rates 0.125 through 2.0 when trained on the frozen 54-state structural subset. Earlier UP-71A evidence showed a coverage transition between 27 and 54 training states after 800 steps. UP-77A asks whether structural coverage, rather than optimization duration, remains the limiting variable when training is restricted to one step.

## Frozen design

Unchanged:

- the same 243 five-role ternary states;
- the same residue-based structural partition and UP-71A secondary/tertiary selectors;
- the same 162 held-out states with primary residue != 0;
- raw factorized one-hot representation, 15 features;
- role-factorized simplex representation, 10 features;
- five independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- exactly one optimization step;
- per-role held-out gate 0.98.

Training-state levels are exactly:

- 9;
- 18;
- 27;
- 54;
- 81.

No adaptive stopping, extra step, learning-rate change, threshold change, result-informed retry, representation change, promotion, or activation is permitted.

## Interpretation

Each coverage/representation point is sealed exactly as observed. The experiment may show the prior 27-to-54 coverage transition, a shifted transition under one-step training, or complete failure. No post-result data selection is allowed.
