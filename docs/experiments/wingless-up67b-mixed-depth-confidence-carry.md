# Wingless UP-67B — mixed-depth confidence carry

Status: preregistered scientific sequence-facing control.

Scientific parent: sealed UP-66B Windows evidence `4d4c814b04f69bcf5128f0d1a12729f62ae8bd5e`.

## Question

UP-66B showed that confidence-carry improves commit accuracy at 512 writes across fixed depths 32, 128, 512, and 1024 with no final or relation regression. UP-67B asks whether that benefit survives an irregular mixed-depth stream rather than a fixed propagation depth.

## Frozen design

Unchanged substrate, observer bank, decoder training, pair [0,5], training depth 0, relation head, closed-loop semantics, and confidence-carry semantics.

Depth pattern, repeated in order:

- 32;
- 128;
- 512;
- 1024.

Schedules:

- 102M;
- 103M.

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

No oracle, retry, retraining, adaptive threshold, threshold selection, result-informed stopping, promotion, production authority, or activation is permitted.

## Interpretation

The experiment asks whether confidence-gated correction remains useful when temporal spacing varies within a sequence. Improvement, neutrality, or degradation are all valid outcomes.
