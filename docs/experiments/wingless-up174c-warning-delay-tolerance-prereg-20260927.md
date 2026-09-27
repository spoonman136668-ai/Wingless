# Wingless UP-174C — warning-delay tolerance

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-173C e5a89e17a3f717d2f0d917479fa0417c35ffe850.

## Question

How quickly must Wingless act after its native warning for the bounded two-shot correction to retain benefit?

## Frozen state space and environment

Reuse UP-173C:
- four cohorts;
- all initial hand positions 0..15;
- no_refresh pressure;
- cadence 2 monitoring;
- 64-write horizon.

## Frozen warning and action rule

Reuse the same current-state robust-horizon warning and endangered-resident target.

At most two targeted refreshes per arm, no more than one action per monitoring interval.

## Frozen delay arms

After a warning-positive monitoring interval:
- delay 0: act in the same interval;
- delay 1: act at the next monitoring interval;
- delay 2: act two monitoring intervals later.

After the first action, a second action requires a later warning-positive interval and uses the same fixed delay.

No-action is the baseline.

## Measurements

Per delay:
- arms;
- actions;
- eventual losses;
- prevented losses;
- prevented losses per action;
- mean loss delay among arms that still fail.

## Interpretation

A sharp drop with one- or two-interval delay would define a narrow corrective window. Persistence under delay would show the warning provides useful lead time rather than an immediate last-moment alarm.

## Bounds

Counterfactual only. Maximum two actions per arm. No adaptive delay, action-budget tuning, threshold tuning, repeated maintenance beyond two actions, future-policy schedule input, capacity change, live activation, or production authority.
