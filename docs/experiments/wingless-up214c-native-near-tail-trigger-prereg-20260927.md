# Wingless UP-214C — native near-risk eighth-action trigger

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-213C 3f042545dd0e2077aa54e1a6be06a8c5fee8b8cc.

## Question

After seven frozen committed refreshes, can Wingless use its own near-risk state to trigger the eighth refresh in time, replacing the fixed nominal tail clock?

## Frozen schedules
Four repeated-switch policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Phase advances:
- 0
- 4 writes

Horizons:
- 66,68,70,72 writes

## Frozen actions 1-7
- cadence 2
- first native critical warning triggers action 1
- actions 2-7 every four monitoring intervals (8 writes)
- target endangered resident

## Frozen tail modes
- none: no eighth action
- nominal: eighth action four monitoring intervals after action 7
- native_near: after action 7, fire the eighth action at the first monitoring boundary where adversarial_shield_horizon <= 4 (two cadences)

The native trigger receives no future policy schedule or horizon information.

## Measurements
Per schedule × phase × horizon × mode:
- baseline/treated losses
- prevented/accelerated losses
- total actions
- eighth actions fired

## Bounds
Counterfactual only. No threshold fitting, adaptive timing, future-policy input to trigger, action budget >8, live activation, or production authority.
