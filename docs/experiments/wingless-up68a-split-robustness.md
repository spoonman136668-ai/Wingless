# Wingless UP-68A — representation split robustness

Status: preregistered scientific coverage/robustness experiment.

Scientific parent: UP-66A seal `a4a00334e101f95ef1b25218058834183526ce69`.

This experiment is defined independently of UP-67A outcomes.

## Question

Does the factorized-versus-entangled generalization contrast remain when the deterministic 81/162 training partition is rotated across all three residue classes?

## Frozen design

The logical domain, observer classes, learner, optimizer, and gate are unchanged from UP-65A/UP-66A:

- 243 five-role ternary states;
- five independent 3-class linear-softmax role heads;
- 800 optimization steps at learning rate 1.0 from zero initialization;
- factorized 64-dimensional Hadamard representation;
- entangled 64-dimensional joint Fourier representation with 32 frequency pairs;
- unchanged per-role held-out gate of 0.98.

Three training partitions are fixed before execution:

- sum(role values) mod 3 == 0;
- sum(role values) mod 3 == 1;
- sum(role values) mod 3 == 2.

Each partition contains 81 training states and 162 held-out states. Both representations are trained independently on each partition.

No representation modification, extra training state, adaptive split selection, threshold change, nonlinear observer, promotion, production authority, or activation is permitted.

## Interpretation

The experiment measures partition robustness only. Variation across residue classes is itself a valid result. Passing or failing does not authorize mechanism promotion.
