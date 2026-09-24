# Wingless UP-74A — compact factorized low-budget boundary

Status: preregistered scientific optimization/efficiency experiment.

Scientific parent: UP-73A seal `6f9ea0d3f771e11831222e0742b7b799ca1e6892`.

## Question

UP-73A showed perfect held-out performance for both compact factorized representations at 54 training states for every tested budget from 200 through 1600 optimization steps.

UP-74A asks where the low-budget boundary begins while preserving the same representation, data split, optimizer, and gate.

## Frozen design

Training set and evaluation are unchanged from UP-73A:

- 54 residue-0 training states selected by the unchanged secondary coordinate;
- 162 held-out states;
- raw factorized one-hot, 15 features;
- role-factorized simplex, 10 features;
- five independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- unchanged per-role held-out gate of 0.98.

Training-step ladder, frozen before execution:

- 25;
- 50;
- 100;
- 200.

No adaptive stopping, learning-rate change, threshold change, extra training state, representation modification, promotion, production authority, or activation is permitted.

## Interpretation

This is a compute-efficiency boundary search. Passing, failure, or non-monotonicity at any budget is valid scientific evidence and must be sealed exactly as observed.
