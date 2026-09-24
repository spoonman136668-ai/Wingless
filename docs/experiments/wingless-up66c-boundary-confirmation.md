# Wingless UP-66C — six-bank boundary confirmation

Status: preregistered scientific robustness confirmation experiment.

Scientific parent: UP-65C seal `1895beda79f067ee1bf6646138445eed828fd3d8`.

## Question

UP-65C localized the tested all-schedule passing boundary to 0.0108 for golden rotation and 0.0086 for fixed irregular, with the next tested noise point producing at least one schedule failure.

UP-66C asks whether those boundary locations replicate on untouched schedules.

## Frozen design

Families:

- golden rotation;
- fixed irregular.

Untouched deterministic schedule bases:

- 99M;
- 100M;
- 101M;
- 102M;
- 103M;
- 104M.

Golden-rotation noise points:

- 0.0106;
- 0.0108;
- 0.0110.

Fixed-irregular noise points:

- 0.0085;
- 0.0086;
- 0.0087.

All other conditions are unchanged from UP-65C: six banks, state dimension 16, 256 scenarios, depth 64, coherence decoder, reversible transport, and the existing 0.99 value / 0.95 exact-scenario gate.

No geometry modification, decoder change, threshold change, family selection, adaptive stopping, production authority, or activation is permitted.

## Interpretation

This is untouched-schedule boundary confirmation. Replication, widening, narrowing, or schedule sensitivity are all valid scientific results and must be sealed exactly as observed. No mechanism promotion is authorized.
