# Wingless UP-171C — two-shot correction state-space replication

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-170C 81cc048a2a5a54fcc5c064b632859b2ad5e65c7c.

## Question

Does the bounded two-shot targeted correction observed in UP-170C generalize beyond the four initial hand positions used there?

## Frozen trigger and intervention

Reuse UP-170C unchanged:
- current-state robust-horizon warning;
- cadences 2 and 4;
- at most two actions per arm;
- at most one action per warning interval;
- second action only at a later warning-positive interval;
- targeted_refresh versus same-budget sham_refresh.

## Frozen state-space expansion

Cohorts remain:
- 0,4,8,12;
- 1,5,9,13;
- 2,6,10,14;
- 3,7,11,15.

Initial hand positions expand from 0,4,8,12 to every position 0..15.

Only arms for which the frozen trigger-state construction succeeds are evaluated.

## Frozen future pressure policies

- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

## Measurements

Per policy × cadence × intervention:
- available trigger arms;
- baseline losses;
- actions used;
- treated losses;
- prevented losses;
- prevented losses per action;
- mean loss-step delay among treated failures.

Also report the original four-hand subset and the expanded non-original hands separately for the no_refresh cadence-2 targeted arm.

## Interpretation

Replication on the expanded hands would support a genuine state-space correction effect. Collapse outside the original four positions would reveal geometric over-specialization.

## Bounds

Counterfactual only. Maximum two actions per arm. No adaptive hand selection, unbounded maintenance, action-budget tuning, threshold tuning, future-policy schedule input, capacity change, live activation, or production authority.
