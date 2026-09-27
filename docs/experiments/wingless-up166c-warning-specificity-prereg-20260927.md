# Wingless UP-166C — robust-warning specificity

Status: preregistered scientific shadow-monitoring experiment.

Scientific parent: sealed UP-165C f93cd9cdc9e4b371dcd46f88ef6aa987ee9cb792.

## Question

After perfect warning recall transferred across several future pressure policies, how costly is that warning rule in false or expired warning intervals, and can it remain quiet on a deliberately protected safe control?

## Frozen warning model

At each sparse check compute the unchanged UP-161C adversarial-shield horizon from current state only.

Warn for the next interval iff:
- horizon <= 2 for cadence 2;
- horizon <= 4 for cadence 4.

No future real-policy schedule is passed into the warning model.

## Real policies

Reuse:
- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

Add one evaluation-only safe control:
- protected_endangered: refresh the endangered resident before each unique pressure write.

The safe-control refresh is part of the exogenous real policy. It is not caused by the warning.

## Measurements

Per policy × cadence:
- arms;
- eventual losses;
- warning intervals;
- warned losses;
- missed losses;
- expired warning intervals;
- interval warning precision;
- warning recall;
- mean first-warning lead writes for warned losses.

## Interpretation

Recall without specificity is not yet a calibrated operational warning. The safe control tests whether the current-state robust horizon remains conservative even when the future environment prevents loss.

## Bounds

Shadow only. Warnings never change memory. No warning-triggered correction, adaptive cadence, threshold tuning, capacity change, semantic priority, live activation, or production authority.
