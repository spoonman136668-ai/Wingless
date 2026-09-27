# Wingless UP-209C — repeated-switch committed tail-cap transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-208C 8f852d58ab5393cbc44dcdfa0f26a068cc9c17af.

## Question

Can the successful repeated-switch committed rescue schedule be shortened by removing only its late tail while preserving the causally important early/middle maintenance sequence?

## Frozen pressure timing

Exactly the four shifted schedules from UP-207C/UP-208C:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Switch phase:
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

## Frozen rescue timing
- cadence 2
- first native critical warning triggers action 1
- later committed opportunities every 4 monitoring intervals = every 8 writes
- target: endangered resident
- horizon: 56 writes

## Frozen action-opportunity caps
- 6 total opportunities
- 7 total opportunities
- 8 total opportunities

No earlier slot moves or warning gating.

## Measurements
Per pressure schedule × cap:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- total actions
- arms with action
- mean treated-minus-baseline loss-step change among failures

## Interpretation

Cap 6 or 7 preserving full rescue would identify a true tail-efficiency reduction. Any loss relative to cap 8 bounds the minimal committed depth for the tested repeated-switch schedules.

## Bounds

Counterfactual only. No threshold/cadence/spacing/target/cohort/horizon tuning, no adaptive cap selection, no replacement or pull-forward action, no live activation or production authority.
