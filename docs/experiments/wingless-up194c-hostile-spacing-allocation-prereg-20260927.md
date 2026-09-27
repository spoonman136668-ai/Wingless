# Wingless UP-194C — hostile corrective-budget spacing

Status: preregistered scientific counterfactual intervention diagnostic.

Scientific parent: sealed UP-193C 8fe1c3cdb7a645463f2654a985c0147fee152a71.

## Question

With the same cap-16 targeted-refresh budget, does spreading corrective actions over a longer duration restore hostile-pressure prevention?

## Frozen policy
- hostile_shield

## Frozen population
- four cohorts
- all initial hand positions 0..15

## Frozen trigger
- first native critical warning: adversarial_shield_horizon <= monitoring cadence

## Frozen cap
- maximum 16 targeted refreshes

## Frozen schedules

Primary monitoring cadence: 2 writes.

After the first critical-warning refresh, subsequent refreshes occur every:
- 1 monitoring interval = every 2 writes
- 2 monitoring intervals = every 4 writes
- 3 monitoring intervals = every 6 writes
- 4 monitoring intervals = every 8 writes

Reference:
- cadence 4, spacing 1 monitoring interval = every 4 writes

No fresh-warning recurrence is required after the first trigger.

## Measurements

Per schedule:
- arms
- baseline losses
- prevented losses
- actions used
- accelerated losses
- mean and maximum treated-minus-baseline loss-step extension among failures

## Interpretation

If cadence-2 spacing 2 reproduces the cadence-4 rescue, temporal coverage rather than monitoring cadence itself explains UP-193C. If no cadence-2 spacing rescues, the monitoring cadence changes the warning/intervention interaction in a deeper way.

## Bounds

Counterfactual only. No threshold tuning, policy change, cap change, new action type, adaptive spacing, future-policy input to warning, live activation, or production authority.
