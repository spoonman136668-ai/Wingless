# Wingless UP-170C — two-shot bounded correction

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-169C dc3131bf32e6860f34a53ab9f76d688e2ffefba6.

## Question

Can a minimal increase from one to at most two warning-triggered targeted refreshes convert the verified delay benefit into actual failure prevention?

## Frozen trigger

Use the unchanged first-warning robust-horizon rule at cadences 2 and 4.

## Frozen intervention budget

At most two actions per arm.
- no more than one action in any monitoring interval;
- after the first action, a second action may occur only at a later monitoring interval that is again warning-positive.

## Interventions

- targeted_refresh: refresh the endangered resident;
- sham_refresh: spend the same action count on a non-endangered resident.

## Policies

- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

## Measurements

Per policy × cadence × intervention:
- baseline losses;
- actions used;
- treated losses;
- prevented losses;
- prevention per action;
- mean loss-step delay for arms that still fail.

## Bounds

Counterfactual only. Maximum two actions per arm. No unbounded maintenance, adaptive action budget, threshold tuning, future-policy schedule input, capacity change, live activation, or production authority.
