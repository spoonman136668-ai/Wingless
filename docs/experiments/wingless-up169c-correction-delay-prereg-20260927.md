# Wingless UP-169C — bounded-correction delay

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-168C 71ba926d3024ffac81b27092a2a30a98aa90fba7.

## Question

UP-168C found that one targeted refresh prevents no eventual losses over the full horizon. Does the same frozen action nevertheless buy measurable time before failure, beyond a same-budget sham refresh?

## Frozen warning and triggers

Reuse unchanged:
- current-state adversarial-shield horizon;
- cadence 2 warns iff horizon <= 2;
- cadence 4 warns iff horizon <= 4;
- confirmation depths 1 and 2.

At most one intervention may occur per arm.

## Frozen interventions

- targeted_refresh: query the endangered resident once at trigger;
- sham_refresh: query a non-endangered resident once at trigger.

## Frozen future pressure policies

- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

## Measurements

Per policy × cadence × confirmation × intervention:
- arms and actions;
- baseline and treated eventual losses;
- baseline and treated loss step, capped at 65 for survival through 64 writes;
- total and mean delay in writes;
- arms delayed;
- arms advanced;
- maximum delay;
- actions producing no positive delay.

## Interpretation

Positive delay unique to targeted refresh would show bounded corrective value even when eventual survival is unchanged. No delay beyond sham would close this one-shot refresh mechanism more strongly.

## Bounds

Counterfactual only. One action maximum per arm. No repeated maintenance, adaptive trigger selection, threshold tuning, future-policy schedule input to warning, capacity change, live activation, or production authority.
