# Wingless UP-66B — 512-write depth-stratified confidence carry

Status: preregistered scientific control experiment.

Scientific parent: sealed UP-65B Windows evidence `957cef220597126b2ca9b99872f0a3b1b5a38a5b`.

## Question

UP-65B found the frozen decoder-margin bins monotonically ordered by correctness in all 16 tested fixed-depth conditions after 512 closed-loop writes. UP-66B asks whether the already-established confidence-carry intervention remains behaviorally useful at that doubled horizon.

## Frozen design

Unchanged substrate, observer bank, decoder training, pair [0,5], training depth 0, relation head, closed-loop semantics, and confidence-carry update semantics.

Held-out propagation depths:

- 32;
- 128;
- 512;
- 1024.

Schedules:

- 100M;
- 101M.

Memory noise:

- 0.085;
- 0.09.

Writes per scenario:

- exactly 512.

Policies per condition:

- baseline closed loop;
- confidence carry threshold 0.25;
- confidence carry threshold 0.50;
- confidence carry threshold 0.75.

Thresholds are inherited unchanged from the prior confidence-carry family. UP-65B does not select a threshold.

No oracle, retry, retraining, adaptive threshold, threshold selection, result-informed stopping, promotion, production authority, or activation is permitted.

## Interpretation

For each fixed depth, confidence-carry may improve, match, or degrade commit/final/relation accuracy. All outcomes are valid and sealed without post-result tuning.
