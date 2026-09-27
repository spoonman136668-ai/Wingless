# Wingless UP-212C — dense committed tail-cap durability boundary

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-211C 25e0427860cd8532401f40cecf3f6f4559da3778.

## Question

At what horizon between 64 and 72 writes does the seven-opportunity committed rescue first diverge from the eight-opportunity comparator?

## Frozen policy pairs
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

## Frozen switch phases
- phase advance 0
- phase advance 4

Policies alternate every 8 writes after the initial A segment.

## Frozen horizons
- 64, 66, 68, 70, 72 writes

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
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- actions taken
- mean loss-step extension among failures

## Interpretation

The first horizon where cap 7 prevents fewer failures than cap 8 is the durability boundary of the tail reduction under that condition.

## Bounds

Counterfactual only. No threshold/cadence/spacing/target/cohort/phase/horizon tuning after results, no adaptive cap selection, no replacement/pull-forward action, no live activation or production authority.
