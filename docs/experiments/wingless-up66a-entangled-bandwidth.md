# Wingless UP-66A — entangled observer bandwidth ladder

Status: preregistered scientific cognition/representation experiment.

Scientific parent: UP-65A seal `bd76bbe20261cccb039c319ae36b0bc2ee34c6b9`.

## Question

UP-65A showed perfect sparse-combination generalization from a dense factorized representation, but only 0.7877 mean held-out accuracy from a 32-frequency-pair entangled joint representation.

UP-66A asks whether that failure is caused by insufficient harmonic bandwidth rather than entanglement itself.

## Frozen design

The logical state domain, 81/162 training split, five independent linear-softmax role heads, optimizer, and training budget are unchanged.

Only the number of sine/cosine frequency pairs in the entangled joint-state representation changes:

- 16 pairs / 32 features;
- 32 pairs / 64 features (UP-65A reproduction point);
- 64 pairs / 128 features;
- 121 pairs / 242 features, the full non-DC real Fourier basis for 243 logical states.

Each point is trained independently from the same zero initialization.

The existing 0.98 per-role held-out gate is unchanged. No nonlinear observer, factorized features, extra training states, or result-informed tuning is permitted.

## Interpretation

If held-out accuracy approaches the gate only at larger bandwidth, the UP-65A failure is primarily observer/feature bandwidth. If even the full Fourier basis fails to generalize from the sparse split, arbitrary joint-state entanglement itself lacks the compositional inductive bias needed for unseen combinations under this learner.
