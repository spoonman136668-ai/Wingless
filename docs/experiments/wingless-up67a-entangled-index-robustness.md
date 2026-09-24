# Wingless UP-67A — entangled index robustness

Status: preregistered scientific representation/robustness experiment.

Scientific parent: UP-66A seal `a4a00334e101f95ef1b25218058834183526ce69`.

## Question

UP-66A showed that increasing harmonic bandwidth from 32 to the full non-DC real Fourier basis did not recover compositional generalization from the entangled joint-state representation.

UP-67A asks whether the observed failure is specific to one arbitrary joint-state enumeration or remains under multiple deterministic bijective renumberings.

## Frozen design

The logical state domain, training split, learner, and scientific gate are unchanged:

- 243 five-role logical states;
- train iff sum(role values) mod 3 == 0, giving 81 train and 162 held-out states;
- five independent 3-class linear-softmax role heads;
- 800 optimization steps at learning rate 1.0 from zero initialization;
- sine/cosine joint-state features with exactly 32 frequency pairs / 64 features;
- unchanged per-role held-out gate of 0.98.

Four index maps are fixed before execution:

- lexicographic: i;
- affine-2-17: (2*i + 17) mod 243;
- affine-80-31: (80*i + 31) mod 243;
- affine-241-7: (241*i + 7) mod 243.

All multipliers are coprime to 243, so every map is a bijection.

No nonlinear observer, factorized representation, additional training state, result-informed map selection, threshold change, or activation authority is permitted.

## Interpretation

Strong dependence on the bijection would show that the entangled observer's apparent generalization is enumeration-sensitive. Failure across all fixed bijections would strengthen the conclusion that arbitrary joint-state entanglement lacks the compositional inductive bias needed for unseen role combinations under this learner. No mechanism is promoted by this experiment.
