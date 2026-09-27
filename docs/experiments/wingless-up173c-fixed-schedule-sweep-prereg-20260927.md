# Wingless UP-173C — fixed two-shot schedule sweep

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-172C 0b8e5a7a6cf10d3b88301fc7e667b8ab2ef1c41f.

## Question

Can any simple global two-shot schedule reproduce the prevention achieved by Wingless's state-dependent warning timing?

## Frozen state space and environment

Reuse UP-172C:
- four cohorts;
- all initial hand positions 0..15;
- no_refresh pressure;
- cadence 2 monitoring;
- 64-write horizon;
- endangered resident is the same target in every targeted arm.

## Frozen intervention schedules

- warning_targeted: unchanged UP-170C/171C warning-triggered two-shot rule.
- fixed targeted schedules, each with at most two actions:
  - 1,3
  - 9,11
  - 17,19
  - 25,27
  - 33,35
  - 41,43
  - 49,51
  - 57,59
- no_action baseline.

Schedules are fixed globally before the run and are not adapted per arm.

## Measurements

Per intervention:
- arms;
- actions taken;
- eventual losses;
- prevented losses versus no action;
- prevented losses per action;
- mean loss delay among arms that still fail.

Also report the best fixed schedule by prevented losses as a descriptive result only. It does not alter the experiment or trigger any retuning.

## Interpretation

If warning_targeted still outperforms every fixed schedule, the timing signal is state-dependent rather than a proxy for one globally favorable window. If a fixed schedule matches or exceeds it, the current warning may be unnecessary for correction timing.

## Bounds

Counterfactual only. Maximum two targeted actions per arm. No adaptive schedule search, no post-result schedule insertion, no action-budget tuning, no threshold tuning, no future-policy schedule input, no capacity change, no live activation, or production authority.
