# Wingless UP-236C — repeated native-risk recheck for action 9

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-235C 0868cccf7a9d0ab0d8242251eef8de5f1fd7c7aa.

## Question

If the first post-action8 critical event is not high-risk, can Wingless keep observing its frozen native Hamming risk state and still authorize the same single action 9 later—recovering missed rescues without restoring broad overintervention?

## Frozen policy / classifier

Exactly UP-235C through action 8, including action-5 retiming.

Frozen Hamming classifier:
- same 12 failure prototypes;
- same 6 survivor prototypes;
- same nine native fields;
- high_risk only when nearest failure Hamming distance < nearest survivor Hamming distance.

No threshold or prototype change.

## Arms

1. cap8 — no action 9.
2. single_critical — fire action 9 at first post-action8 adversarial-critical event.
3. hamming_first — evaluate only first such event; fire only if high-risk.
4. hamming_recheck — after a rejected first event, evaluate each later adversarial-critical event; fire action 9 at the first event classified high-risk.

At most one action 9 may ever fire.

## Fresh heldout conditions

Cohorts:
- 0,9,10,15
- 1,4,11,14
- 2,5,8,13
- 3,6,7,12

Phase advances:
- 0
- 22

Four frozen policy pairs; horizon 80; monitoring cadence 2.

512 arms.

## Measurements

Per condition and overall:
- losses by arm;
- ninth actions;
- rescues versus cap8;
- unnecessary ninth actions on cap8 survivors;
- ineffective ninth attempts;
- delayed high-risk authorizations after an initially rejected event.

## Bounds

Counterfactual only. No classifier changes, no prototype mutation, no new native fields, no threshold fitting, no action budget >9, no second ninth action, no cadence change, no future-policy input, no live activation or production authority.
