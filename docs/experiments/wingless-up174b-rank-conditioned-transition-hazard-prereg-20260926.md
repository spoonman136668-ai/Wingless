# Wingless UP-174B — rank-conditioned transition hazard

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-173B d3c7ab99f3facc5b7c84a69db20bf8baefbd40ec.

## Question

UP-173B reproduced a higher hazard when an old example transitions from one in-band step to exactly two consecutive in-band steps. Does that temporal state add predictive information beyond static vulnerability rank itself?

## Frozen design

Use the accepted vulnerability bands unchanged:
- STORE rank <= 10;
- OBSERVE rank <= 10;
- REPORT rank <= 4.

Use the same six subjects, two 12-step cleanup paths, and 120 old examples per state, but build the frozen post-REPORT state at terminal phase 19.

Before every update:
1. compute each example's static absolute-margin rank;
2. update its consecutive in-band streak using the class-specific accepted band;
3. assign one frozen rank stratum:
   - rank_1_4;
   - rank_5_10;
4. for examples inside one of those strata and currently in-band, assign temporal state:
   - current_only = streak 1;
   - transition_to_persistent = streak 2;
   - established_persistent = streak 3+.

Measure whether the subsequent frozen update crosses correctness.

## Primary comparison

Within each fixed rank stratum, compare crossing density for transition_to_persistent against current_only and established_persistent.

If transition hazard remains elevated within the same rank stratum, temporal history contributes information beyond raw static rank.

## Bounds

Shadow only. No maintenance trigger, update suppression, adaptive threshold, added model calls, training, capacity change, result-informed retry, live activation, or production authority.
