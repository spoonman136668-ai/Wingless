# Wingless UP-182C — horizon-selective correction burden

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-181C a8c6b5bb2df4761e0fb00f1e0666543304835970.

## Question

With the best bounded correction rule frozen, does Wingless prevent failures more often than it intervenes unnecessarily for a given operational horizon?

## Frozen environment

- no_refresh pressure
- four cohorts
- all initial hand positions 0..15
- cadences 2 and 4

## Frozen best rule

- warning iff current adversarial-shield horizon <= cadence
- stage one acts immediately on the endangered resident
- action two requires a later fresh warning
- stage two waits exactly one monitoring interval, then acts on the endangered resident
- maximum two actions per arm

## Frozen evaluation horizons

- 16 writes
- 24 writes
- 32 writes
- 40 writes
- 48 writes
- 56 writes
- 64 writes

Each horizon is evaluated independently from the same frozen initial state.

## Paired baseline classification

For each arm and horizon:
- baseline_failure: no-action loss occurs on or before the horizon.
- baseline_survivor: no-action loss occurs after the horizon.

For the intervention arm:
- prevented_failure: baseline_failure but intervention survives beyond the horizon.
- unnecessary_action_arm: baseline_survivor but the correction rule executes at least one action before the horizon.
- necessary_action_arm: baseline_failure and the rule executes at least one action.
- no_action_survivor: baseline_survivor and no action fires.

## Measurements

Per cadence × horizon:
- arms
- baseline failures
- baseline survivors
- intervention actions
- arms with any action
- prevented failures
- unnecessary-action arms
- necessary-action arms
- no-action survivors
- prevented / unnecessary-action-arm ratio where defined

## Interpretation

The frozen rule is selective at a horizon if it preserves failures with bounded intervention burden and does not systematically act on arms that would have survived that horizon without correction.

## Bounds

Counterfactual only. No threshold tuning, adaptive horizon, action-budget change, policy input, capacity change, live activation, or production authority.
