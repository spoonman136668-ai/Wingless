# Wingless UP-211C — committed tail-cap horizon transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-210C 9b6227c8ddf0d5ef28e248173b929b1af8b95da1.

## Question

Does the seven-opportunity committed rescue remain equivalent to the eight-opportunity comparator when the observation horizon extends beyond 56 writes?

## Frozen policy pairs
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

## Frozen switch phases
- phase advance 0: initial A length 8 writes
- phase advance 4: initial A length 4 writes

After the initial segment, policies alternate every 8 writes.

## Frozen horizons
- 56 writes
- 64 writes
- 72 writes

Each horizon is evaluated independently from the same initial state.

## Frozen population
- heldout contiguous quartets
- initial hands 0..15

## Frozen rescue rule
- cadence 2
- first native critical warning triggers action 1
- later committed opportunities every 8 writes
- target endangered resident

## Frozen caps
- 7 total opportunities
- 8 total opportunities

## Measurements
Per policy pair × phase advance × horizon × cap:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- actions taken
- mean loss-step extension among failures

## Interpretation

If cap 7 matches cap 8 through horizon 72, the eighth opportunity is genuinely redundant over a longer durability window. Divergence identifies the horizon at which the eighth action becomes necessary.

## Bounds

Counterfactual only. No threshold/cadence/spacing/target/cohort/phase/horizon tuning after results, no adaptive cap selection, no pull-forward/replacement action, no live activation or production authority.
