# Wingless UP-67C — six-bank boundary depth replication

Status: preregistered scientific scale/robustness experiment.

Scientific parent: UP-65C seal `1895beda79f067ee1bf6646138445eed828fd3d8`.

This experiment is defined independently of UP-66C outcomes.

## Question

UP-65C localized the depth-64 schedule-robust boundary near 0.0108 for golden rotation and 0.0086 for fixed irregular.

UP-67C asks whether that boundary neighborhood remains stable when reversible transport depth changes.

## Frozen design

Untouched deterministic schedule bases:

- 105M;
- 106M.

Families and noise neighborhoods:

- golden rotation: 0.0106, 0.0108, 0.0110;
- fixed irregular: 0.0085, 0.0086, 0.0087.

Reversible transport depths:

- 32;
- 64;
- 128.

All other conditions are unchanged: six banks, state dimension 16, 256 scenarios, coherence decoder, reversible transport, and the existing 0.99 value / 0.95 exact-scenario gate.

Prototype transport and scenario transport use the same preregistered depth for each point.

No geometry modification, decoder change, threshold change, family selection, adaptive stopping, production authority, or activation is permitted.

## Interpretation

This is a boundary-local depth replication. Depth invariance, shifts in the apparent boundary, or failures are all valid results and must be sealed exactly as observed. No mechanism promotion is authorized.
