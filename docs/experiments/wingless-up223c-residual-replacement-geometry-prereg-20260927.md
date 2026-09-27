# Wingless UP-223C — residual replacement-geometry diagnostic

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-222C ef2fe261ce81565591f981182efd5c8a25575c7f.

## Question

Does the single cap-8 failure missed by both accepted horizon signals exhibit a distinct native replacement trajectory in Wingless's existing hand/age state?

## Frozen residual environment

Policy sequence:
- alternating_shield / fixed_offset_refresh

Phase advance:
- 8 writes

Horizon:
- 80 writes

Population:
- heldout contiguous quartets
- initial hands 0..15

## Frozen behavior

Cap-8 rule only:
- actions 1-7 unchanged committed schedule
- action 8 fires at first post-action-7 normal monitoring boundary with adversarial_shield_horizon <= 2
- no action 9

This experiment does not change intervention behavior.

## Existing native geometry observed after action 8

Immediately before every subsequent write:
- clock hand
- endangered key slot index
- endangered slot age
- next replacement slot predicted by the accepted hand+age predictor
- whether predicted slot is endangered
- circular hand-to-endangered distance
- age histogram over used slots
- adversarial_shield_horizon
- no_query_horizon

## Per-arm summaries

- loss / survival
- action-8 step
- loss step
- whether an adversarial critical warning occurs post-action8
- minimum endangered age
- minimum hand-to-endangered distance
- number of snapshots where predicted replacement slot is endangered
- final pre-loss/pre-horizon geometry

For the silent failing arm only:
- retain the full post-action8 snapshot trace.

## Aggregate groups

- silent_loss: cap8 loss with no post-action8 adversarial critical warning
- warned_loss
- survive

Per group report ranges for the frozen per-arm geometry summaries.

## Interpretation

A geometric distinction in the silent arm would identify a native replacement-trajectory hypothesis for a future disjoint preregistered warning test. No trigger is selected here.

## Bounds

Diagnostic only. No intervention, threshold fitting, adaptive feature selection, future-policy input to geometry, action-budget change, live activation, or production authority.
