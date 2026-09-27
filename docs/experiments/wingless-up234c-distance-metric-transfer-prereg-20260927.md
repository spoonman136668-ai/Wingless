# Wingless UP-234C — native distance-metric transfer

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-233C 016fe588ae331a4456b35bc094b8c56c5d147d17.

## Question

Does magnitude-aware distance over the same native geometry improve transferable failure-risk separation where equal-weight Hamming leaves ties?

## Frozen prototypes

Exactly the 12 failure and 6 survivor prototypes inherited from sealed UP-231C through UP-233C.

## Frozen native fields

The same nine fields only:
- endangered age
- hand distance
- predicted replacement slot is endangered
- age-0 / age-1 / age-2 / age-3 resident counts
- adversarial horizon
- no-query horizon

## Frozen classifiers

Hamming comparator:
- one point per unequal field.

Manhattan primary:
- absolute numeric difference per integer field;
- boolean mismatch contributes one;
- all fields equal weight;
- no learned scaling or weights.

For each metric:
- high_risk if nearest failure distance < nearest survivor distance;
- low_risk if nearest survivor distance < nearest failure distance;
- unknown on a tie;
- no_event if no post-action8 critical event occurs.

## New heldout evaluation

Cohort topology:
- 0,4,10,14
- 1,5,11,15
- 2,6,8,12
- 3,7,9,13

Phase advances:
- 4 writes
- 18 writes

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Horizon 80; monitoring cadence 2; exact cap-8 policy; no action 9.

512 arms total.

## Measurements

Per metric:
- failures / survivors
- event-bearing failures / survivors
- high / low / unknown / no-event
- high-risk failure recall
- high-risk precision
- low-risk failures
- unknown failures / survivors

## Interpretation

Improved recall at preserved selectivity would show that the native geometry transfers through magnitude structure rather than categorical identity alone. Failure would reject simple unweighted Manhattan abstraction.

## Bounds

Shadow diagnostic only. No heldout fitting, no new native fields, no learned scaling/weights, no threshold fitting, no corrective action after action 8, no cohort-specific tuning, no cadence change, no future-policy input, no live activation, or production authority.
