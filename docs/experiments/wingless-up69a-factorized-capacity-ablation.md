# Wingless UP-69A — factorized capacity ablation

Status: preregistered scientific representation/mechanism-ablation experiment.

Scientific parent: UP-67A seal `1c055b437255a40aa0158a97f981d0f7a094e9d4`.

This experiment is defined independently of UP-68A outcomes.

## Question

UP-67A showed that all four frozen bijective renumberings of the 64-feature joint Fourier representation failed the unchanged 0.98 per-role held-out gate.

UP-69A asks whether the successful generalization seen in the factorized representation depends on its 64-feature distributed Hadamard width, or whether lower-dimensional role-separated encodings retain the same inductive bias.

## Frozen design

The logical domain, train/held-out split, learner, optimizer, and gate are unchanged:

- 243 five-role ternary states;
- training iff sum(role values) mod 3 == 0;
- 81 training and 162 held-out states;
- five independent 3-class linear-softmax role heads;
- 800 optimization steps at learning rate 1.0 from zero initialization;
- unchanged per-role held-out gate of 0.98.

Four observer arms are fixed before execution:

- raw role-factorized one-hot: 15 features;
- role-factorized simplex: 10 features, two coordinates per role;
- distributed Hadamard factorized: 64 features;
- joint Fourier entangled control: 64 features / 32 frequency pairs.

No nonlinear observer, additional training state, adaptive arm selection, threshold change, mechanism promotion, production authority, or activation is permitted.

## Interpretation

If the compact role-factorized arms preserve held-out accuracy while the joint entangled control fails, factorization rather than observer width is the stronger explanatory variable. Any contrary result is equally valid and must be sealed unchanged.
