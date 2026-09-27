# Wingless UP-210C — committed tail-cap switch-phase transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-209C f30b098ea5914c3edd2f513633b2084d3afe1a70.

## Question

Does the frozen seven-opportunity committed rescue remain equivalent to the eight-opportunity comparator when the repeated-switch phase is changed?

## Frozen policy pairs
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

## Frozen phase advances

Relative to the original 8-write A block:
- advance 0: initial A length 8 writes
- advance 2: initial A length 6 writes
- advance 4: initial A length 4 writes
- advance 6: initial A length 2 writes

After the initial A segment:
- B runs for 8 writes;
- then A/B alternate every 8 writes through horizon 56.

## Frozen population
- heldout contiguous quartets
- initial hands 0..15

## Frozen rescue rule
- cadence 2
- first native critical warning triggers action 1
- later committed opportunities every 4 monitoring intervals = every 8 writes
- target endangered resident
- horizon 56

## Frozen caps
- 7 total opportunities: primary
- 8 total opportunities: comparator

## Measurements
Per policy pair × phase advance × cap:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- actions taken
- mean loss-step extension among failures

## Interpretation

Matching cap-8 prevention across all phase advances supports the seventh-opportunity cap as a robust tail reduction. Any gap bounds the reduction to a narrower switch phase.

## Bounds

Counterfactual only. No threshold/cadence/spacing/target/cohort/horizon tuning, no adaptive phase or cap selection, no pull-forward/replacement action, no live activation or production authority.
