# Wingless UP-186C — native warning lead calibration

Status: preregistered scientific diagnostic.

Scientific parent: sealed UP-185C 97fd5882ceb29b815fdca2a5480822a70df9875c.

## Question

How far before actual no-action capability loss does the frozen native warning first trigger?

## Frozen environment

- no_refresh pressure
- four cohorts
- all initial hand positions 0..15
- cadences 2 and 4
- no intervention
- warning iff current adversarial-shield horizon <= cadence
- evaluate through write 64

## Measurements

Per cadence:
- arms
- arms with a warning before loss
- arms lost by 64
- min/max/mean warning-to-loss lead in writes
- lead <= cadence count
- lead <= 2×cadence count
- warnings occurring after loss (must remain zero)

## Bounds

Diagnostic only. No intervention, threshold tuning, adaptive cadence, policy change, capacity change, live activation, or production authority.
