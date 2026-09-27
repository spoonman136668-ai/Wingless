# Wingless UP-188C — native risk policy transfer

Status: preregistered scientific diagnostic.

Scientific parent: sealed UP-187C d575783d52fac813325204d884dfde0573ea648c.

## Question

Do the frozen native failure-risk buckets remain calibrated when the external pressure policy changes?

## Frozen risk model

At each monitoring checkpoint:
- native state signal = adversarial_shield_horizon
- critical: horizon <= cadence
- near: cadence < horizon <= 2 × cadence
- far: horizon > 2 × cadence

No threshold changes.

## Frozen policies

- no_refresh
- alternating_shield
- fixed_offset_refresh
- hostile_shield

## Frozen cadences

- 2 writes
- 4 writes

## Frozen population

- four cohorts
- all initial hand positions 0..15
- no intervention
- evaluate through write 64

## Outcomes

Per policy × cadence × risk bucket:
- checkpoint count
- failures within one cadence
- failures within two cadences
- empirical one-cadence failure rate
- empirical two-cadence failure rate
- eventual losses by write 64 among contributing trajectories

The warning/risk computation uses only current native state. Future policy execution is used only to label outcomes.

## Interpretation

Monotonic separation that persists across policies supports a policy-robust native failure-risk estimate. Loss of separation identifies the external-policy boundary of the current risk model.

## Bounds

Diagnostic only. No intervention, threshold fitting, adaptive policy selection, cadence tuning, future-policy schedule used by the warning, live activation, or production authority.
