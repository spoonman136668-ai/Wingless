# Wingless UP-75A — compact factorized optimization-step floor

Status: preregistered scientific efficiency-boundary experiment.

Scientific parent: sealed UP-74A Windows evidence `04cd3f9b2e77f706e624790902fcd8bc022d104b`.

## Question

UP-74A passed the frozen 0.98 per-role held-out gate for both compact factorized representations at every tested budget down through 25 optimization steps. UP-75A asks where that optimization floor actually begins.

## Frozen design

All UP-74A data, representation, optimizer, and interpretation rules remain unchanged:

- the same 54 residue-0 training states selected by the unchanged secondary coordinate;
- the same 162 held-out states;
- raw factorized one-hot representation, 15 features;
- role-factorized simplex representation, 10 features;
- five independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- per-role held-out gate 0.98.

The only preregistered change is the training-step ladder:

- 1;
- 2;
- 4;
- 8;
- 16.

No adaptive stopping, learning-rate change, threshold change, extra training state, representation change, result-informed retry, promotion, or activation is permitted.

## Interpretation

Every point is sealed exactly as observed. The first passing budget, if any, is descriptive only; no post-result tuning is allowed. A complete failure is a valid scientific result.
