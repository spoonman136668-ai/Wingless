# Wingless UP-152C — age-state ablation

Status: preregistered scientific durable-memory forecast ablation.

Scientific parent: sealed UP-151C 8f59d4ee2494e5d641be5f1b2c76922b3d657033.

## Question
Are the two-bit age counters actually necessary for exact next-eviction prediction once refreshes perturb the replacement path, or is the hand position alone sufficient?

## Frozen held-out schedule
Use the accepted memory unchanged, exact recall cap 16, starting hands 0/4/8/12, and 16 unique admissions. Use a disjoint sparse-refresh schedule:
- admission 2 -> key 1
- admission 5 -> key 5
- admission 8 -> key 9
- admission 11 -> key 13
- admission 14 -> key 2

Missed refreshes are not substituted.

## Frozen predictors
1. full_state: current hand + 16 two-bit ages, using the accepted UP-151C predictor.
2. hand_only: predict the current hand slot directly; it cannot inspect ages, identity, class, or future schedule.

## Measurements
Prediction accuracy for both predictors over the same 64 admissions.

## Bounds
No policy change, semantic information, adaptive refresh, future oracle, capacity increase, result-informed retry, live activation, or production authority.
