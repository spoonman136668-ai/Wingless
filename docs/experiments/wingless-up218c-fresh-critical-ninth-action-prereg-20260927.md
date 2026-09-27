# Wingless UP-218C — fresh-critical ninth-action tail repair

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-217C fd3336e688a2e0a61a4150e1cfbe59a540d6568c.

## Question

Can exactly one additional bounded tail correction, triggered only by a fresh later critical warning after action 8, recover the late-horizon failures exposed in UP-217C without creating a large unnecessary-action burden?

## Frozen environment

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Phase advances:
- 8 writes
- 12 writes

Horizons:
- 74
- 76
- 78
- 80 writes

Population:
- heldout contiguous quartets
- initial hands 0..15

## Frozen action sequence

Actions 1-7:
- unchanged committed schedule.

Action 8:
- fire at first monitoring boundary after action 7 where adversarial_shield_horizon <= 2.

Action 9:
- may fire only after action 8 has already fired;
- requires a later monitoring boundary with a fresh adversarial_shield_horizon <= 2;
- cannot fire on the same monitoring boundary as action 8;
- maximum total actions = 9.

Comparator:
- frozen critical-tail cap-8 rule from UP-217C.

No threshold/timing changes.

## Measurements

Per schedule × phase advance × horizon:
- cap8 losses
- cap9 losses
- ninth actions
- ninth rescues
- unnecessary ninth actions
- necessary ninth attempts
- ineffective ninth attempts
- no-ninth survivors
- rescue / unnecessary ratio where defined

## Interpretation

If fresh-critical action 9 removes late cap-8 failures with bounded unnecessary action, that supports iterative bounded repair: intervene, re-evaluate native state, and intervene once more only on renewed danger.

## Bounds

Counterfactual only. No threshold fitting, adaptive selection, future-policy input, action budget >9, schedule retuning, live activation, or production authority.
