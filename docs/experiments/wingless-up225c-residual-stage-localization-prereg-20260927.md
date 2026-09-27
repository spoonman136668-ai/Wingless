# Wingless UP-225C — residual correction-stage localization

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-224C 6b5bc0d6774374ea79d6594bd548101345af51ca.

## Question

At what exact correction stage does the one silent residual arm fail, and what existing native replacement geometry precedes that failure from the first correction onward?

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

Cap-8 rule unchanged:
- action 1 fires on the accepted initial adversarial critical warning
- actions 2-7 follow the accepted committed spacing
- action 8 fires only on the accepted adversarial critical warning
- no action 9

No action timing, threshold, or target changes.

## Recorded action progression

For every arm:
- exact write step of actions 1 through 8, using -1 when an action is never reached
- number of actions taken before loss/horizon
- loss step / survival

## Existing native geometry observed

After action 1 has fired, immediately before every write:
- clock hand
- endangered slot
- endangered age
- accepted next-replacement prediction
- predicted-is-endangered
- circular hand-to-endangered distance
- age histogram
- adversarial_shield_horizon
- no_query_horizon

Per arm summarize:
- minimum adversarial horizon
- minimum no-query horizon
- minimum endangered age
- minimum hand distance
- predicted-is-endangered count

For any loss that occurs without a post-action-1 adversarial critical warning, retain the full native snapshot trace.

## Aggregate measurements

For action counts 0 through 8:
- loss arms
- survivor arms

Also report all silent-loss arm summaries and traces.

## Interpretation

This diagnostic localizes the blind spot to a specific correction interval and exposes the existing native state immediately before failure. It does not select a new trigger or intervention.

## Bounds

Diagnostic only. No intervention change, threshold fitting, adaptive feature selection, future-policy input to geometry, action-budget change, live activation, or production authority.
