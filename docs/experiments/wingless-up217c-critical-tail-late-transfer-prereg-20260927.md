# Wingless UP-217C — critical tail-trigger late transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-216C 14b37dbfb429e4264bb3d279f1ff3488f79efbd9.

## Question

Does the frozen critical action-8 trigger retain its rescue/selectivity behavior at later phase advances and longer horizons?

## Frozen environment

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Heldout phase advances:
- 8 writes
- 12 writes

Heldout horizons:
- 74
- 76
- 78
- 80 writes

Population:
- heldout contiguous quartets
- initial hands 0..15

## Frozen correction

Actions 1-7:
- unchanged committed schedule.

Action 8:
- fire at first monitoring boundary after action 7 where adversarial_shield_horizon <= 2.

Comparator:
- cap 7, no eighth action.

No threshold/timing changes from UP-216C.

## Measurements

Per schedule × phase advance × horizon:
- cap7 losses
- critical-tail losses
- eighth actions
- tail rescues
- unnecessary eighth actions
- necessary eighth attempts
- ineffective eighth attempts
- no-eighth survivors
- rescue / unnecessary ratio where defined

## Interpretation

Transfer of low treated losses with bounded unnecessary action supports the critical trigger beyond the original phase/horizon window. Any loss of rescue is a genuine transfer boundary because threshold and schedule remain frozen.

## Bounds

Counterfactual only. No threshold fitting, adaptive selection, future-policy input, action budget >8, schedule retuning, live activation, or production authority.
