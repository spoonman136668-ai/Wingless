# Wingless UP-232C — frozen critical-event signature risk transfer

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-231C 3b6edd53b9cf2bfcb70f413fcba4e0e82cd8deb3.

## Question

Do the exact native critical-event signatures discovered in UP-231C transfer as a failure-risk classifier to disjoint cohort topology and phase timing without refitting?

## Frozen classifier

At the first post-action8 adversarial-critical event, compute the exact UP-231C native signature:

- endangered age
- hand distance
- predicted replacement slot is endangered
- age-0 / age-1 / age-2 / age-3 resident counts
- adversarial horizon
- no-query horizon

Classification is fixed before evaluation:

- high_risk: signature is one of the 12 UP-231C failure signatures
- low_risk: signature is one of the 6 UP-231C survivor signatures
- unknown: event signature is in neither frozen set
- no_event: no post-action8 critical event occurs

The 12 failure signatures and 6 survivor signatures are copied exactly from the sealed UP-231C result.

No heldout outcome may alter those sets.

## Heldout cohort topology

A new partition not used by UP-231C:

- 0,6,9,15
- 1,7,8,14
- 2,4,11,13
- 3,5,10,12

## Heldout phase advances

- 10 writes
- 14 writes

## Frozen policy pairs

- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

## Frozen horizon and policy

- horizon 80 writes
- monitoring cadence 2 writes
- exact UP-231C cap-8 policy
- no action 9

Eight condition cells; 512 arms total.

## Measurements

Across all arms:
- failures / survivors
- event-bearing failures / survivors
- high-risk / low-risk / unknown / no-event decisions

For high-risk:
- failures classified high-risk
- survivors classified high-risk
- failure recall
- precision

For low-risk:
- failures classified low-risk
- survivors classified low-risk

For unknown:
- failures / survivors

Also report exact-signature coverage among event-bearing arms.

## Interpretation

High failure recall with low survivor false-positive burden would show that native replacement geometry generalizes as a transferable risk classifier. Large unknown or misclassified fractions would bound exact-signature transfer and motivate a more abstract representation.

## Bounds

Shadow diagnostic only. No heldout fitting, no signature additions/removals, no adaptive feature selection, no threshold fitting, no corrective action after action 8, no cohort-specific tuning, no cadence change, no future-policy input, no live activation, or production authority.
