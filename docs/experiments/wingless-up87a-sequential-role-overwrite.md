# Wingless UP-87A — sequential role/filler overwrite bridge

Status: preregistered scientific sequence-state experiment.

Scientific parent: sealed UP-86A Windows evidence `c4b554ccacee9d0a7807603be070e6ed0d6857ae`.

## Question

UP-86A showed that the six-role compact factorized readouts remain perfect under static feature perturbations through amplitude 0.10. UP-87A moves from static states to an actual sequence/state-tracking task: can the same one-step-trained readouts preserve latest-value-wins semantics across repeated role overwrites while perturbation accumulates between writes?

## Frozen design

Training remains unchanged:

- six ternary roles;
- 27-state structural training cell: primary==0, secondary==0, tertiary==0;
- raw factorized one-hot representation, 18 features;
- role-factorized simplex representation, 12 features;
- six independent 3-class linear-softmax heads;
- zero initialization;
- exactly one optimization step;
- learning rate 1.0;
- no retraining under sequence noise.

Sequence evaluation:

- 64 deterministic scenarios;
- initial role tuples generated deterministically from the six-role state space;
- one role overwritten per sequence step;
- the written role slice is refreshed from the exact current role/filler representation;
- all other state is retained;
- deterministic additive perturbation is applied to the complete retained feature state after every write;
- latest-value-wins truth is scored after every write.

Sequence lengths:

- 16;
- 64;
- 256.

Per-step perturbation amplitudes:

- 0;
- 0.0005;
- 0.001;
- 0.002.

Metrics:

- role-value accuracy over all role/step events;
- exact six-role state accuracy after each write.

Diagnostic gate only:

- value accuracy >= 0.99;
- exact state accuracy >= 0.95.

No denoising, adaptive normalization, clipping, retraining, seed search, threshold change, representation change, result-informed retry, promotion, or activation is permitted.

## Interpretation

This is a bridge from static compositional representation to sequential state tracking. A pass supports carrying the representation into the serialized sequence qualification work. Failure localizes accumulated-state fragility and is sealed without tuning.
