# Wingless UP-189C — correction policy transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-188C a64adbb6d9098d810eb908a0aa0987d25db96216.

## Question

Does the frozen two-stage correction rule remain beneficial when the external pressure policy changes, using the same native critical warning?

## Frozen policies

- no_refresh
- alternating_shield
- fixed_offset_refresh
- hostile_shield

## Frozen cadences

- 2 writes
- 4 writes

## Frozen population

- four cohorts
- all initial hand positions 0..15

## Frozen rule

Native warning:
- critical iff adversarial_shield_horizon <= cadence

Actions:
- maximum two targeted refreshes;
- stage one acts immediately on the endangered resident at the first critical warning;
- stage two requires a later fresh critical warning;
- after that second warning, stage two waits exactly one monitoring interval before acting;
- both actions target the endangered resident.

No policy-specific thresholds or timing changes.

## Comparator

Paired no-action baseline under the same external policy.

## Measurements

Per policy × cadence:
- arms
- baseline losses by write 64
- treated losses by write 64
- prevented losses
- accelerated losses
- stage-one actions
- stage-two actions
- prevented per total action
- mean treated-minus-baseline loss-step change among arms that still fail

## Interpretation

Positive prevention with bounded action cost across multiple policies supports transfer of the frozen self-maintenance rule. Failure under specific policies identifies the intervention-policy boundary without retuning.

## Bounds

Counterfactual only. No threshold tuning, policy-specific action logic, adaptive delay, action-budget change, future-policy schedule used by warning, capacity change, live activation, or production authority.
