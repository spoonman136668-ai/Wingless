# Wingless UP-64C — six-bank reversible-depth ladder

Status: preregistered scientific scale/robustness experiment.

Scientific parent: UP-62C seal `e2da0e6e6ee17b2eae765e07105ea2049e2a99e5`.

This experiment is defined independently of UP-63C outcomes.

## Question

UP-62C showed schedule-stable separation between golden rotation and fixed irregular at noise 0.01, while both passed at 0.0075.

UP-64C asks whether that robustness profile is stable as reversible transport depth changes substantially around the existing depth-64 condition.

## Frozen design

Families:

- golden rotation;
- fixed irregular.

Untouched deterministic schedule bases:

- 91M;
- 92M.

Noise levels:

- 0.0075;
- 0.0100.

Reversible transport depths:

- 16;
- 32;
- 64;
- 128;
- 256.

All other conditions are unchanged: six banks, state dimension 16, 256 scenarios, coherence decoder, reversible transport, and the existing 0.99 value / 0.95 exact-scenario gate.

Prototype transport and scenario transport always use the same preregistered depth for a point.

No geometry modification, decoder change, threshold change, family selection, adaptive stopping, production authority, or activation is permitted.

## Interpretation

This is a fixed depth-scale ladder. Depth dependence, invariance, or failure are all valid results and must be sealed exactly as observed. No mechanism promotion is authorized.
