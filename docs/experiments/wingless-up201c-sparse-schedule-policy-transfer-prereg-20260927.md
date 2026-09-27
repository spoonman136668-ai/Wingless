# Wingless UP-201C — sparse rescue schedule policy transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-200C 70915a195815f578ebae2f62ed2e2e1823ee1195.

## Question

Does the hostile-derived sparse rescue schedule transfer unchanged across distinct external pressure policies?

## Frozen policies
- no_refresh
- alternating_shield
- fixed_offset_refresh
- hostile_shield

## Frozen population

Heldout contiguous cohorts:
- 0,1,2,3
- 4,5,6,7
- 8,9,10,11
- 12,13,14,15

Initial hand positions:
- 0..15

## Frozen monitoring and rescue rule

Monitoring cadence:
- 2 writes

Native trigger:
- first monitoring boundary where adversarial_shield_horizon <= cadence

After trigger:
- immediate targeted refresh of endangered resident
- then one targeted refresh every 4 monitoring intervals = every 8 writes
- maximum 8 actions

Evaluation horizon:
- 56 writes

The rule is identical for all policies.

## Comparator

Paired no-action baseline under the same policy.

## Measurements

Per policy:
- arms
- baseline losses by horizon
- treated losses by horizon
- prevented losses
- accelerated losses
- total actions
- arms with at least one action
- mean loss-step extension among arms that still fail

## Interpretation

Broad prevention across policies would show the sparse temporal-coverage rule generalizes beyond the hostile condition where it was derived. Policy-specific failures would localize remaining intervention dependence despite the policy-robust critical warning.

## Bounds

Counterfactual only. No policy-specific threshold, cadence, spacing, cap, target, cohort, or horizon tuning; no adaptive action logic; no live activation or production authority.
