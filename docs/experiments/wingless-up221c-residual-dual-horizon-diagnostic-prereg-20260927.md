# Wingless UP-221C — residual dual-horizon diagnostic

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-220C b06af67a6c5ef511c168d03dfdefd87f687000c5.

## Question

On the cap-8 residual boundary, does the existing no-query durability horizon expose danger that the adversarial-shield horizon misses?

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

## Frozen intervention

Cap-8 rule only:
- actions 1-7 unchanged committed schedule
- action 8 fires at the first post-action-7 normal monitoring boundary where adversarial_shield_horizon <= 2
- no action 9

The diagnostic does not alter behavior.

## Frozen observation window

For each arm, after action 8 fires and before loss or horizon 80:
- observe before every write
- compute adversarial_shield_horizon
- compute no_query_horizon

Both horizons are existing native counterfactual state measures and are evaluated on copies of current state.

## Per-arm measurements

- cap8 loss / survival
- action-8 step
- loss step
- minimum post-action8 adversarial_shield_horizon
- minimum post-action8 no_query_horizon
- whether adversarial <=2 occurred
- whether no_query <=2 occurred
- whether no_query <=4 occurred

## Aggregate measurements

For cap8 failures and cap8 survivors:
- arms
- adversarial<=2 count/rate
- no_query<=2 count/rate
- no_query<=4 count/rate

Also:
- cap8 failures with no adversarial<=2 warning
- among those silent failures, counts caught by no_query<=2 and no_query<=4

## Interpretation

If no-query horizon catches the silent residual with acceptable survivor-warning burden, Wingless already contains a complementary native danger signal. If it does not, the residual is not explained by either accepted horizon.

## Bounds

Diagnostic only. No intervention, no threshold fitting, no adaptive signal selection, no future-policy input to the horizons, no action-budget change, no live activation or production authority.
