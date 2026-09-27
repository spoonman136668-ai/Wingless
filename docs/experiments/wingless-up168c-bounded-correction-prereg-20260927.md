# Wingless UP-168C — bounded warning-triggered correction

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-167C 8acfc2c4e35967728992598246194efa38cd79e2.

## Question

Can a frozen one-action rule use Wingless's native warning to prevent failure more often than it wastes an intervention?

## Frozen warning model

Reuse unchanged:
- current-state adversarial-shield horizon;
- cadence 2 warns iff horizon <= 2;
- cadence 4 warns iff horizon <= 4.

## Frozen trigger rules

Evaluate independently:
- confirmation 1: first warning interval;
- confirmation 2: first interval with two consecutive warnings.

At most one intervention may occur per arm.

## Frozen interventions

- targeted_refresh: query the endangered resident once at the trigger;
- sham_refresh: spend the same one-query budget on a non-endangered resident.

No-action outcome is measured from the same initial state and future pressure policy.

## Frozen future pressure policies

- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

## Measurements

Per policy × cadence × confirmation × intervention:
- arms;
- baseline losses;
- actions taken;
- treated losses;
- prevented losses;
- harmed arms where baseline survived but treatment lost;
- wasted actions that did not prevent a baseline loss;
- prevented losses per action.

## Interpretation

A targeted rule is promising only if its prevention efficiency exceeds its waste and the sham action does not reproduce the effect. This experiment does not authorize activation.

## Bounds

Counterfactual only. One action maximum per arm. No repeated maintenance, adaptive trigger selection, threshold tuning, future-policy schedule input to the warning model, capacity change, live activation, or production authority.
