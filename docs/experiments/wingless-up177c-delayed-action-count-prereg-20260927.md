# Wingless UP-177C — delayed action-count ablation

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-176C 8b79af488a6d5864078c1e6bc8ab065fd64ad7b3.

## Question

How much of the verified delayed correction benefit comes from the first warning-triggered action, and how much requires the bounded second warning-triggered action?

## Frozen state space and environment

- four cohorts
- all initial hand positions 0..15
- no_refresh pressure
- cadences 2 and 4
- 64-write horizon

## Frozen timing rule

- warning iff current adversarial-shield horizon <= cadence
- wait exactly one monitoring interval
- act
- every later action requires a new warning-positive interval and the same one-interval delay

## Frozen action caps

Evaluate independently:
- maximum 1 action per arm
- maximum 2 actions per arm

Target controls:
- targeted_refresh
- sham_refresh

## Measurements

Per cadence × action cap × intervention:
- arms
- baseline losses
- actions
- treated losses
- prevented losses
- prevented losses per action
- mean loss delay among treated failures

## Interpretation

If one delayed action already prevents most failures, the second action is incremental. If prevention appears mainly or only with two actions, the bounded recurrence of warning-action cycles is mechanistically important.

## Bounds

Counterfactual only. Delay fixed at one monitoring interval. No adaptive action cap, delay, threshold, future-policy schedule input, capacity change, live activation, or production authority.
