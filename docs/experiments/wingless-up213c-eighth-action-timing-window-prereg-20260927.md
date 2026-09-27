# Wingless UP-213C — eighth-action timing window

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-212C c4f349bf8341b725d435c3109c24c81a00feb438.

## Question

At the newly identified durability boundary, is the eighth committed refresh effective only within a narrow timing window?

## Frozen policy pairs
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

## Frozen switch phases
- phase advance 0
- phase advance 4

## Frozen horizons
- 66,68,70,72 writes

## Frozen actions 1-7
- cadence 2
- first native critical warning triggers action 1
- actions 2-7 remain committed every 4 monitoring intervals = every 8 writes
- target endangered resident

## Frozen eighth-action modes
- none: cap 7 comparator
- early: eighth action after 3 monitoring intervals following action 7
- nominal: eighth action after 4 monitoring intervals
- late: eighth action after 5 monitoring intervals

No other timing changes.

## Measurements
Per policy pair × phase advance × horizon × mode:
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- actions taken
- mean loss-step extension among failures

## Interpretation

A narrow best timing identifies a causal maintenance window for the eighth refresh. Similar protection across early/nominal/late would indicate action presence matters more than precise timing.

## Bounds

Counterfactual only. No threshold/cadence/spacing/target/cohort/phase/horizon tuning after results, no adaptive timing selection, no extra action budget beyond eight, no pull-forward rule, no live activation or production authority.
