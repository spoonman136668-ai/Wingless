# Wingless UP-233C — nearest-signature native risk transfer

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-232C 5054470dab9f5b157c907da40ec0da9b08a4ff20.

## Question

Can a frozen similarity rule over the existing C231 native geometry generalize failure risk beyond exact signature identity?

## Frozen training exemplars

Use only the 12 failure signatures and 6 survivor signatures frozen in UP-232C, all originating from sealed UP-231C.

No UP-232C heldout outcomes are used as exemplars.

## Frozen representation

Nine existing native fields:
- endangered age
- hand distance
- predicted replacement slot is endangered
- age-0 / age-1 / age-2 / age-3 resident counts
- adversarial horizon
- no-query horizon

## Frozen distance/classifier

Equal-weight Hamming distance across the nine discrete fields.

For the first post-action8 adversarial-critical event:
- compute distance to every frozen failure exemplar;
- compute distance to every frozen survivor exemplar;
- high_risk if minimum failure distance < minimum survivor distance;
- low_risk if minimum survivor distance < minimum failure distance;
- unknown on a tie;
- no_event if no post-action8 critical event occurs.

No weights, thresholds, or distance tuning.

## New heldout evaluation

Cohort topology:
- 0,7,10,13
- 1,4,11,14
- 2,5,8,15
- 3,6,9,12

Phase advances:
- 6 writes
- 16 writes

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Horizon:
- 80 writes

Monitoring cadence:
- 2 writes

This topology/timing was not used by UP-231C or UP-232C.

## Measurements

Across 512 arms:
- failures / survivors;
- event-bearing failures / survivors;
- high-risk / low-risk / unknown / no-event;
- high-risk failure recall;
- high-risk precision;
- low-risk failure count;
- mean minimum failure/survivor Hamming distance by outcome.

## Interpretation

Useful recall and precision on this disjoint evaluation would show that the native geometry supports an abstract risk neighborhood rather than only memorized exact states. Failure would bound this representation and require a different abstraction.

## Bounds

Shadow diagnostic only. No heldout fitting, no new native fields, no adaptive weights, no threshold fitting, no corrective action after action 8, no cohort-specific tuning, no cadence change, no future-policy input, no live activation, or production authority.
