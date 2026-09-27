# Wingless UP-164C — robust-horizon sparse warning

Status: preregistered scientific shadow-monitoring experiment.

Scientific parent: sealed UP-163C 382c766c147ea3a2e91e6874b7e296d8b52b5261.

## Question

UP-163C recovered perfect sparse-check recall by simulating the exact known pressure policy across the next interval. Can current native state support the same sparse monitoring without being given the exact future schedule?

## Frozen starting states

Reuse the 16 UP-161C post-guided-refresh starting states.

## Frozen real pressure

Use the same hostile local pressure process as UP-161C/162C:
- before each unique write, if the predicted next victim is not the endangered item, refresh that resident item;
- then issue the unique write.

## Current-state robust horizon

At each scheduled check, compute **adversarial shield horizon** from a copy of the current state using UP-161C's frozen current-state rule:
- repeatedly protect the currently predicted non-endangered victim;
- issue hypothetical unique writes;
- stop when the endangered item would be evicted.

This computation sees only current hand, ages, entries, and the endangered item. It receives no real future write keys, refresh schedule, or externally supplied future actions.

## Sparse monitoring arms

- cadence 2: warn for the next interval iff robust horizon <= 2;
- cadence 4: warn for the next interval iff robust horizon <= 4.

A warning remains active until the next scheduled check. Warnings do not alter real memory.

## Measurements

Per cadence:
- eventual losses;
- losses preceded by an active warning;
- missed losses;
- warning intervals that expired without loss;
- warning recall.

## Bounds

Shadow only. No corrective action, adaptive cadence, exact future schedule, semantic priority, capacity increase, extra training, live activation, or production authority.
