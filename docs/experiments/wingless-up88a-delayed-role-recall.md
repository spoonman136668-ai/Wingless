# Wingless UP-88A — delayed role recall under distractor overwrites

Status: preregistered scientific sequence-memory experiment.

Scientific parent: sealed UP-87A Windows evidence `eda467a422930058f7ef380f02aac7a240240560`.

## Question

UP-87A preserved perfect six-role latest-value-wins state tracking through 256 writes with small deterministic perturbation. UP-88A asks a stricter delayed-memory question: after writing a target role once, can the trained compact factorized readout recover that role after a long sequence of distractor overwrites to every other role while perturbation accumulates on the retained target slice?

## Frozen design

Training is unchanged:

- six ternary roles;
- 27-state structural training cell: primary==0, secondary==0, tertiary==0;
- raw factorized one-hot representation, 18 features;
- role-factorized simplex representation, 12 features;
- six independent 3-class linear-softmax heads;
- zero initialization;
- exactly one optimization step;
- learning rate 1.0;
- no retraining under sequence noise.

Evaluation:

- 64 deterministic scenarios;
- each scenario writes one target role/value at the start;
- the target role is never overwritten again;
- all later writes target only the other five roles;
- each distractor write refreshes only that written role slice from exact current truth;
- deterministic additive perturbation is then applied to the complete retained feature state;
- the target role and full six-role state are decoded only at the final query.

Delays:

- 16;
- 64;
- 256;
- 1024.

Per-step perturbation amplitudes:

- 0;
- 0.0005;
- 0.001;
- 0.002.

Metrics:

- target-role recall accuracy;
- exact six-role final-state accuracy.

Diagnostic gate only:

- target-role recall >= 0.99;
- exact final-state accuracy >= 0.95.

No denoising, clipping, adaptive normalization, retraining, seed search, threshold change, representation change, result-informed retry, promotion, or activation is permitted.

## Interpretation

This experiment tests delayed selective memory under structured interference, not benchmark optimization. A pass supports moving the factorized mechanism toward serialized sequence qualification. A delay/noise boundary identifies where persistent-state protection or exact recall support becomes necessary.
