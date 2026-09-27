# Wingless UP-166C — warning precision and lead-time calibration

Status: preregistered scientific shadow-monitoring experiment.

Scientific parent: sealed UP-165C f93cd9cdc9e4b371dcd46f88ef6aa987ee9cb792.

## Question

With recall already perfect in UP-165C, how precise is the frozen current-state warning rule, how many warning intervals expire without loss, and how much write-level lead time does the first warning provide before loss?

## Frozen warning model

Reuse unchanged:
- current-state adversarial-shield horizon;
- cadence 2 warns iff horizon <= 2;
- cadence 4 warns iff horizon <= 4.

No threshold, horizon rule, cadence, or pressure policy is changed.

## Frozen real pressure policies

- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

## Measurements

Per policy × cadence:
- arms;
- eventual losses;
- warned and missed losses;
- total warning intervals;
- expired warning intervals;
- interval precision = warned-loss intervals / all warning intervals;
- warning intervals per eventual loss;
- mean first-warning lead writes;
- minimum and maximum first-warning lead writes.

## Interpretation

This experiment does not ask whether warnings can catch failure; UP-165C already established that in the tested policies. It asks whether the warning burden is calibrated tightly enough to support a later bounded correction rule without excessive unnecessary interventions.

## Bounds

Shadow only. No corrective action, no adaptive cadence, no threshold tuning, no future policy schedule supplied to the warning model, no capacity change, no live activation, or production authority.
