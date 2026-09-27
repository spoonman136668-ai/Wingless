# Wingless UP-169C — one-shot delay benefit

Status: preregistered scientific counterfactual intervention diagnostic.

Scientific parent: sealed UP-168C 71ba926d3024ffac81b27092a2a30a98aa90fba7.

## Question

Although one first-warning refresh did not prevent eventual loss, does it buy meaningful write-level time before failure relative to no action and a same-budget sham refresh?

## Frozen trigger

Use confirmation 1 only: the first interval satisfying the unchanged robust-horizon warning rule.

Cadences:
- 2;
- 4.

## Frozen interventions

- targeted_refresh: query the endangered resident once;
- sham_refresh: query a non-endangered resident once.

Exactly one action maximum per arm.

## Frozen policies

- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

## Measurements

Per policy × cadence × intervention:
- arms;
- actions;
- baseline loss step;
- treated loss step;
- mean/min/max loss-step delta;
- number of arms delayed, unchanged, or accelerated.

## Interpretation

A positive targeted delay that exceeds sham would show bounded corrective value even though permanent survival failed. No delay would close the one-shot action family and justify testing a different corrective mechanism rather than simply more of the same action.

## Bounds

Counterfactual only. No repeated maintenance, adaptive triggers, threshold changes, future-policy schedule input to the warning model, capacity change, live activation, or production authority.
