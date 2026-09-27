# Wingless UP-192C — bounded post-warning commitment

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-191C 7cd014c6e0eab81cd9fc2f73d4cd14941a2edcc6.

## Question

After the first native critical warning, is fresh-warning recurrence itself preventing Wingless from using enough corrective capacity under hostile pressure?

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
- 2
- 4
- 6
- 8

## Frozen trigger
- first critical warning iff adversarial_shield_horizon <= cadence

## Commitment rule

Before the first critical warning:
- no action.

At the first critical warning:
- targeted refresh of the endangered resident.

After that first action:
- on each subsequent monitoring interval, perform one targeted refresh of the endangered resident until the preregistered cap is reached;
- no fresh-warning recurrence is required.

This is compared against the UP-191C fresh-warning-gated rule at the same cap.

## Measurements

Per policy × cadence × cap:
- baseline losses
- fresh-gated prevented losses and actions
- commitment prevented losses and actions
- commitment accelerated losses
- mean treated-minus-baseline loss-step change among commitment arms that still fail

## Interpretation

If commitment rescues hostile-pressure arms where fresh-warning gating cannot, the post-action warning state is masking persistent vulnerability. If commitment also fails, targeted refresh itself is inadequate under hostile pressure.

## Bounds

Counterfactual only. No threshold tuning, policy-specific rule, adaptive cap, new action type, future-policy schedule used by the first warning, live activation, or production authority.
