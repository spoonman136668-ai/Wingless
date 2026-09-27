# Wingless UP-208C — repeated-switch committed maintenance slot ablation

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-207C 9023c20d428e3f742b2f9840bf67e52164ec0f3b.

## Question

Which post-trigger committed maintenance positions are actually necessary for full rescue under repeated switching pressure?

## Frozen pressure timing

Exactly the four shifted schedules from UP-207C:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Schedule:
- writes 1–4: A
- 5–12: B
- 13–20: A
- 21–28: B
- 29–36: A
- 37–44: B
- 45–52: A
- 53–56: B

## Frozen population
- heldout contiguous quartets
- initial hands 0..15

## Frozen rescue schedule
- cadence 2
- first native critical warning triggers action 1
- nominal spacing: every 4 monitoring intervals = 8 writes
- maximum 8 scheduled action opportunities
- target: endangered resident
- horizon: 56 writes

## Frozen ablations

Comparator:
- omit ordinal 0 = full committed schedule.

Ablations:
- omit exactly one post-trigger scheduled action ordinal: 2,3,4,5,6,7,8.

When an ordinal is omitted:
- no refresh occurs at that scheduled opportunity;
- the original 8-write clock is preserved;
- later scheduled opportunities remain at their original positions;
- no replacement action is inserted;
- maximum actual actions therefore becomes 7 for that arm if all eight opportunities are reached.

## Measurements

Per pressure schedule × omitted ordinal:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- actions taken
- arms with action
- mean treated-minus-baseline loss-step change among failures

## Interpretation

Any omitted ordinal that causes failures is causally necessary under at least one tested switch schedule. Ordinals removable without loss identify real action-efficiency opportunities without introducing new warning or timing logic.

## Bounds

Counterfactual only. No post-result slot movement, threshold/cadence/spacing/target/cohort/horizon tuning, no replacement action, no adaptive omission, no extra action budget, no live activation or production authority.
