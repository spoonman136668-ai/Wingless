# Wingless UP-64B — depth-stratified confidence-carry intervention

Status: preregistered scientific control experiment.

Scientific parent: sealed UP-63B Windows evidence `e0b9d3ed351bdd20ee336655c0daf5dce0e9ea55`.

## Question

UP-63B found the frozen decoder-margin bins monotonically ordered by empirical correctness in every tested 256-write schedule/noise/depth condition. UP-64B tests whether the already-established confidence-carry intervention retains behavioral utility when each condition is restricted to one propagation depth.

## Frozen design

Unchanged substrate, observer bank, decoder training, pair [0,5], training depth 0, relation head, closed-loop semantics, and confidence-carry update semantics.

Held-out propagation depths: 32, 128, 512, 1024.
Schedules: 96M, 97M.
Memory noise: 0.085, 0.09.
Writes per scenario: 256.
Policies per condition:
- baseline closed loop;
- confidence carry threshold 0.25;
- confidence carry threshold 0.50;
- confidence carry threshold 0.75.

These thresholds are inherited unchanged from the preregistered UP-57B/UP-60B policy family; UP-63B does not select a threshold.

No oracle, retry, retraining, adaptive threshold, threshold selection, result-informed stopping, promotion, production authority, or activation is permitted.

## Interpretation

For each fixed depth, confidence-carry may improve, match, or degrade commit/final/relation accuracy. Every outcome is valid and is sealed without post-result tuning.
