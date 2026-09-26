# Wingless UP-153C — hand-state ablation

Status: preregistered scientific durable-memory forecast ablation.

Scientific parent: sealed UP-152C 45772d81c0ea096e21fa0f1c33b7cfe43cc60051.

## Question
Is replacement-hand position necessary for exact next-eviction forecasting once age state is known, or can ages alone recover the victim?

## Frozen held-out schedule
Use the accepted memory unchanged, exact recall cap 16, starting hands 0/4/8/12, and 16 unique admissions. Use a new sparse-refresh schedule:
- admission 3 -> key 4
- admission 7 -> key 8
- admission 10 -> key 12
- admission 13 -> key 0
- admission 16 -> key 6

Missed refreshes are not substituted.

## Frozen predictors
1. full_state: current hand + 16 two-bit ages using the accepted predictor.
2. age_only: use the same age-transition scan but always start at slot 0; it cannot inspect hand, identity, class, or future schedule.

## Measurements
Prediction accuracy for both predictors over the same 64 admissions.

## Bounds
No policy change, semantic information, adaptive refresh, future oracle, capacity increase, result-informed retry, live activation, or production authority.
