# Wingless UP-207C — warning-driven maintenance pull-forward transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-206C e8a269efb3bf247d954617416d1fd783bba5faa4.

## Question

Can a fresh native critical warning safely pull the next committed maintenance refresh forward when a spacing-phase miss emerges, without increasing the action cap?

## Fresh heldout pressure timing

Use the same four policy pairs as UP-204C, but shift the repeated-switch phase by 4 writes.

For each pair A/B:
- writes 1–4: A
- 5–12: B
- 13–20: A
- 21–28: B
- 29–36: A
- 37–44: B
- 45–52: A
- 53–56: B

Pairs:
1. no_refresh / hostile_shield
2. hostile_shield / no_refresh
3. alternating_shield / fixed_offset_refresh
4. fixed_offset_refresh / alternating_shield

These timing schedules were not used in UP-204C–206C.

## Frozen population
- heldout contiguous quartets
- initial hands 0..15

## Frozen shared parameters
- cadence 2
- initial critical trigger: adversarial_shield_horizon <= 2
- targeted endangered-resident refresh
- nominal spacing: 4 monitoring intervals = 8 writes
- action cap: 8
- horizon: 56 writes

## Comparator: committed schedule
After initial trigger:
- refresh every 4 monitoring intervals regardless of fresh warning.

## Primary: warning pull-forward
After initial trigger:
- maintain the same 4-interval committed schedule;
- at every monitoring boundary, compute the current native critical warning before action;
- if critical occurs before the next nominal refresh, consume that next maintenance action immediately;
- reset the 4-interval spacing clock from the pull-forward action;
- total actions remain capped at 8.

A pull-forward consumes future scheduled maintenance; it does not add action budget.

## Measurements
Per policy-pair schedule × mode:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- total actions
- pull-forward actions
- arms with action
- mean loss extension among treated failures

## Interpretation

Better prevention or equal prevention with no action-cap increase would support native warning as a bounded phase-correction mechanism layered on top of committed maintenance. Earlier cap exhaustion or worse survival would falsify the pull-forward rule.

## Bounds
Counterfactual only. No extra action budget, no threshold/cadence/spacing/cap/target/cohort/horizon tuning, no policy-specific rule, no future-schedule input to warning, no adaptive selection, no live activation or production authority.
