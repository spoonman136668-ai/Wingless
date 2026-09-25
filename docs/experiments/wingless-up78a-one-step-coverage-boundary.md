# Wingless UP-78A — one-step coverage-boundary refinement

Status: preregistered scientific boundary-refinement experiment.

Scientific parent: sealed UP-77A Windows evidence `e90596fcf15feef8ab13523e9a6bda5dc52feff9`.

## Question

UP-77A showed that both frozen compact factorized representations fail the 0.98 per-role held-out gate at 27 training states and pass perfectly at 54 training states after exactly one optimization step. UP-78A asks where that transition lies while preserving the same structural selector family.

## Frozen design

Unchanged from UP-77A:

- the same 243 five-role ternary states;
- the same primary residue split;
- the same 162 held-out states with primary residue != 0;
- raw factorized one-hot representation, 15 features;
- role-factorized simplex representation, 10 features;
- five independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- exactly one optimization step;
- per-role held-out gate 0.98.

Training-state levels are exactly:

- 27: secondary selector s == 0;
- 36: s == 0 plus s == 1, tertiary t == 0;
- 45: s == 0 plus s == 1, tertiary t < 2;
- 54: s < 2.

The secondary/tertiary cells contain exactly nine states each, so the ladder is nested in fixed nine-state increments and introduces no new random sampling.

No adaptive stopping, extra optimization step, learning-rate change, threshold change, seed search, result-informed retry, representation change, promotion, or activation is permitted.

## Interpretation

Every coverage/representation cell is sealed exactly as observed. The experiment may locate the transition at 36 or 45 states, retain only the 54-state pass, or produce representation-specific transitions. Harness acceptance is independent of the scientific result.
