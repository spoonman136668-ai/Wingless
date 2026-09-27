# Wingless UP-231C — first post-action8 critical geometry

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-230C 869367b4c7b66ced2d757f3d1b0a4be12e0eeedd.

## Question

At the first post-action8 adversarial-critical event, do existing native replacement-geometry variables distinguish cap-8 failures from cap-8 survivors on the interleaved topology?

## Frozen environment

Interleaved cohorts:
- 0,5,10,15
- 1,6,11,12
- 2,7,8,13
- 3,4,9,14

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Phase advances:
- 8 writes
- 12 writes

Horizon:
- 80 writes

Monitoring cadence:
- 2 writes

Population:
- initial hands 0..15

Eight condition cells; 512 arms total.

## Frozen cap-8 policy

Exactly the UP-230C cap-8 arm:
- first action at initial adversarial critical warning;
- actions 2-4 on committed 8-write spacing;
- action 5 may advance before its original due boundary only when adversarial_shield_horizon <= 2;
- actions 6-7 retain original due boundaries;
- no schedule compression;
- action 8 at first later adversarial critical warning after action 7;
- no action 9.

## Frozen diagnostic event

After action 8, capture the first later monitoring boundary where:
- adversarial_shield_horizon <= 2.

No intervention occurs at that event.

## Frozen native snapshot fields

Exactly existing state variables:
- endangered age
- hand distance to endangered slot
- predicted replacement slot is endangered
- age-0 / age-1 / age-2 / age-3 resident counts
- adversarial horizon
- no-query horizon

Also retain hand, endangered slot, predicted slot, action8 step, event step, and eventual loss step for traceability.

## Measurements

Across all arms:
- cap-8 failures
- cap-8 survivors
- event-bearing failures
- event-bearing survivors
- silent failures with no post-action8 critical event
- survivors with no event

For event-bearing arms:
- raw frozen native snapshot
- full native-geometry signature over the frozen snapshot fields

Report:
- distinct failure signatures
- distinct survivor signatures
- shared signatures between failure and survivor groups

## Interpretation

Zero shared signatures would show that existing native replacement geometry contains enough information to distinguish useful versus unnecessary action-9 opportunities, even though simple persistence did not. Shared signatures would show the current native snapshot is still aliased and another state dimension is required.

## Bounds

Shadow diagnostic only. No corrective action after action 8, no threshold fitting, no adaptive feature selection, no new native feature, no cohort-specific tuning, no cadence change, no future-policy input, no live activation, no production authority.
