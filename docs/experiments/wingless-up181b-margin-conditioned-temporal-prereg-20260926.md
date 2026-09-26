# Wingless UP-181B — margin-conditioned temporal hazard

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-180B 4373af4553fa7a86c03a705f9ee85a1e62d68da5.

## Question

Does consecutive vulnerability history add residual next-update failure information after conditioning directly on continuous absolute decision margin?

## Frozen design

Held-out terminal phase: 26.

Reuse the same six subjects, two 12-step cleanup paths, 120 old examples per state, and accepted class-specific vulnerability bands.

Before every frozen update:
- compute absolute current decision margin;
- update the consecutive in-band streak;
- assign one fixed absolute-margin band:
  - m0_005: [0, 0.05);
  - m005_010: [0.05, 0.10);
  - m010_020: [0.10, 0.20);
  - m020_040: [0.20, 0.40);
  - m040_plus: [0.40, +inf).
- assign temporal state:
  - current_only = streak 1;
  - transition_to_persistent = streak 2;
  - established_persistent = streak 3+.

These margin bands are fixed before phase-26 execution and are not derived from phase-26 outcomes.

## Measurements

For ALL and each event class, per margin band × temporal state:
- slots;
- correctness crossings;
- crossing density.

Primary comparison: within the same margin band, does transition_to_persistent exceed established_persistent and/or current_only?

## Bounds

Shadow only. No intervention, threshold fitting, model fitting, extra training, adaptive margin bands, capacity change, result-informed retry, live activation, or production authority.
