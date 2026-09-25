# Wingless UP-89A — multi-query latest-value-wins sequence recall

Status: preregistered scientific sequence-intelligence experiment.

Scientific parent: sealed UP-88A Windows evidence `e97939a2ede84e22349c0ca2658be15fa48ceda7`.

## Question

UP-88A preserved perfect delayed role recall through 1,024 distractor writes at every frozen perturbation level. UP-89A moves from one protected target to an MQAR-like setting: after a long stream of repeated key/value overwrites, can the same compact factorized state answer multiple role-key queries using latest-value-wins semantics?

## Frozen design

Training is unchanged:

- six ternary role keys;
- 27-state structural training cell: primary==0, secondary==0, tertiary==0;
- raw factorized one-hot representation, 18 features;
- role-factorized simplex representation, 12 features;
- six independent 3-class linear-softmax readout heads;
- zero initialization;
- exactly one optimization step;
- learning rate 1.0;
- no retraining under sequence noise.

Evaluation:

- 64 deterministic scenarios;
- each sequence begins from a deterministic six-role tuple;
- every sequence step overwrites one role key with one ternary value;
- repeated key writes are allowed and latest-value-wins is the truth rule;
- only the written role slice is refreshed from exact current truth;
- deterministic additive perturbation is applied to the complete retained feature state after each write;
- four deterministic role-key queries are issued at the end of each episode.

Sequence lengths:

- 64;
- 256;
- 1024.

Per-step perturbation amplitudes:

- 0;
- 0.001;
- 0.002.

Metrics:

- query-value accuracy over all four queries;
- whole-query-set exact accuracy.

Diagnostic gate only:

- query accuracy >= 0.99;
- whole-query-set exact accuracy >= 0.95.

No denoising, clipping, adaptive normalization, retraining, seed search, threshold change, exact-recall side channel, representation change, result-informed retry, promotion, or activation is permitted.

## Interpretation

This is a bridge into MQAR-style sequence intelligence. A pass supports carrying factorized state into the serialized sequence qualification gate. Failure localizes where repeated overwrite/query interference exceeds the static readout mechanism.
