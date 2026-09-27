# Wingless UP-187C — native-state failure-risk calibration

Status: preregistered scientific diagnostic.

Scientific parent: sealed UP-186C 0565a60b0388bf927f138ee2c467440352d6d6dc.

## Question

Does the native adversarial-shield horizon define calibrated failure-risk strata before observable loss?

## Frozen environment

- no_refresh pressure
- four cohorts
- all initial hand positions 0..15
- cadences 2 and 4
- no intervention
- monitor at the same cadence boundaries used by accepted warning experiments
- evaluate through write 64

## Frozen native risk buckets

At each monitoring checkpoint before loss:
- critical: adversarial_shield_horizon <= cadence
- near: cadence < horizon <= 2 × cadence
- far: horizon > 2 × cadence

No bucket threshold is chosen after results.

## Frozen outcomes

For each checkpoint:
- failure within one cadence
- failure within two cadences

The no-action path uses the accepted baseline write sequence.

## Measurements

Per cadence × risk bucket:
- checkpoint count
- one-cadence failures and empirical rate
- two-cadence failures and empirical rate
- mean adversarial-shield horizon
- mean actual writes remaining to baseline loss

## Interpretation

Monotonic empirical failure rates from far → near → critical would support a calibrated ordinal risk estimate derived entirely from native state.

## Bounds

Diagnostic only. No intervention, threshold fitting, adaptive bucketing, cadence tuning, future-policy input, live activation, or production authority.
