# Wingless UP-179C — stage-target ablation

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-178C 1f456fdca38cba2f5a4535af5241105304b57da0.

## Question

In the two-stage warning/action sequence, which stages must target the endangered resident?

## Frozen environment and timing

- no_refresh pressure
- cadences 2 and 4
- four cohorts
- all initial hand positions 0..15
- warning iff adversarial-shield horizon <= cadence
- each action waits exactly one monitoring interval after its warning
- action two requires a later fresh warning
- maximum two actions

## Preregistered target arms

- targeted_targeted
- targeted_sham
- sham_targeted
- sham_sham

A sham action refreshes a non-endangered resident using the existing frozen sham-selection rule.

## Measurements

Per cadence × target arm:
- arms
- stage-one actions
- stage-two actions
- eventual losses
- prevented losses
- mean loss delay among treated failures

## Interpretation

If targeted_sham loses prevention, the second targeted action is necessary. If sham_targeted cannot reach or reproduce prevention, the first targeted action is a causal bridge. Only targeted_targeted succeeding would support a genuinely sequential two-stage correction.

## Bounds

Counterfactual only. No adaptive target selection, timing, action budget, warning threshold, policy input, capacity change, live activation, or production authority.
