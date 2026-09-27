# Wingless UP-191C — critical-warning corrective-budget sweep

Status: preregistered scientific counterfactual intervention diagnostic.

Scientific parent: sealed UP-190C c77d611c31413454c4c8f0772fbec26ce495060c.

## Question

If two targeted refreshes are insufficient under several external pressure policies, does increasing the bounded corrective budget eventually turn delay into prevention?

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

## Frozen action caps

- 1
- 2
- 3
- 4
- 6
- 8
- 12
- 16

## Frozen rule

Native warning:
- critical iff adversarial_shield_horizon <= cadence

For cap = 1:
- stage one targeted refresh acts immediately at first critical warning.

For cap >= 2:
- stage one is unchanged;
- stage two requires a later fresh critical warning and waits exactly one monitoring interval before acting, unchanged from the accepted rule.

For caps > 2:
- after stage two has acted, each later fresh critical warning may trigger one immediate targeted refresh, until the preregistered cap is reached.

All actions target the endangered resident. The rule is identical across external policies.

## Comparator

Paired no-action baseline under the same policy and cadence.

## Measurements

Per policy × cadence × cap:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- actions taken
- prevented per action
- mean treated-minus-baseline loss-step change among arms that still fail

## Interpretation

A cap at which prevention first appears estimates the corrective-capacity boundary for this fixed action class. If prevention never appears by cap 16, targeted refresh alone is insufficient under that policy in the tested regime.

## Bounds

Counterfactual diagnostic only. No threshold tuning, policy-specific rule, adaptive cap, new action type, capacity growth, future-policy input to warning, live activation, or production authority.
