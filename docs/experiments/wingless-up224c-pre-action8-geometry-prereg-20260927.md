# Wingless UP-224C — pre-action8 replacement-geometry diagnostic

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-223C bd93d6f27aa8c2899e646d70e426135c13804d4f.

## Question

What existing native replacement geometry precedes the one failure that occurs after action 7 but before action 8 ever triggers?

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
- actions 1-7 use the accepted committed schedule
- action 8 fires only at a normal monitoring boundary where adversarial_shield_horizon <= 2
- no action 9

No behavior changes.

## Observation window

Start:
- immediately when action 7 fires.

Stop per arm:
- action 8 fires, or
- endangered capability is lost, or
- horizon 80 is reached.

At every monitoring boundary before the action-8 decision, and before the second write of an interval when action 8 has not fired, record existing native state:
- clock hand
- endangered slot index
- endangered age
- accepted next-replacement prediction from hand+ages
- predicted-is-endangered
- circular hand-to-endangered distance
- age histogram
- adversarial_shield_horizon
- no_query_horizon

## Per-arm summaries

- action-7 step
- action-8 step
- loss step
- outcome
- observation count
- minimum adversarial horizon
- minimum no-query horizon
- minimum endangered age
- minimum hand distance
- count of predicted-is-endangered snapshots

Full trace is retained for any arm that loses before action 8.

## Aggregate groups

- pre_action8_loss
- action8_reached
- no_action8_survive

Report ranges of the frozen geometry summaries per group.

## Interpretation

The experiment may identify a replacement-trajectory hypothesis for later disjoint testing, but it does not select or activate a warning rule.

## Bounds

Diagnostic only. No intervention change, threshold fitting, adaptive feature selection, future-policy input to geometry, action-budget change, live activation, or production authority.
