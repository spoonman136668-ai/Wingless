# Wingless UP-68C — six-bank deep propagation extension

Status: preregistered scientific scale/robustness experiment.

Scientific parent: UP-67C seal `7aa5f2cf56a5d61aef33807ee86ae25f0644c7f6`.

## Question

UP-67C found identical per-condition outcomes across depths 32, 64, and 128 for the tested boundary neighborhoods.

UP-68C asks whether that depth invariance persists substantially farther into the reversible propagation regime.

## Frozen design

Untouched deterministic schedule bases:

- 107M;
- 108M.

Families and boundary-local noise points are unchanged:

- golden rotation: 0.0106, 0.0108, 0.0110;
- fixed irregular: 0.0085, 0.0086, 0.0087.

Reversible transport depths:

- 256;
- 512;
- 1024.

All other conditions are unchanged: six banks, state dimension 16, 256 scenarios, coherence decoder, reversible transport, and the existing 0.99 value / 0.95 exact-scenario gate.

Prototype transport and scenario transport use the same preregistered depth for each point.

No geometry modification, decoder change, threshold change, family selection, adaptive stopping, production authority, or activation is permitted.

## Interpretation

This is a fixed deep-propagation scale extension. Invariance, drift, abrupt failure, or schedule sensitivity are all valid scientific results and must be sealed exactly as observed.
