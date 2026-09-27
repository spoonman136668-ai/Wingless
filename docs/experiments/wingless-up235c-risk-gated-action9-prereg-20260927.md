# Wingless UP-235C — native-risk-gated action 9

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-234C ec5e25131dc99cd85a2bccec16b324a3a6118e7c.

## Question

Can the frozen Hamming native-risk classifier selectively authorize action 9, preserving most rescue benefit while reducing unnecessary ninth actions?

## Frozen heldout topology

Cohorts:
- 0,8,9,15
- 1,4,10,14
- 2,5,11,12
- 3,6,7,13

Phase advances:
- 2 writes
- 20 writes

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Horizon:
- 80 writes

Monitoring cadence:
- 2 writes

512 arms total.

## Frozen common policy through action 8

All arms use the same qualified sequence:
- initial action on adversarial critical warning;
- actions 2-4 on committed spacing;
- action 5 may advance on fresh adversarial critical warning;
- actions 6-7 retain original due boundaries;
- action 8 fires on the next adversarial critical warning;
- no schedule compression;
- same target and action type.

## Frozen action-9 arms

1. cap8:
- no action 9.

2. single_critical:
- at the first post-action8 adversarial-critical event, fire action 9.

3. hamming_risk:
- at that same first post-action8 adversarial-critical event, compute the frozen UP-233C Hamming classifier from the 12 failure and 6 survivor prototypes;
- fire action 9 only when class = high_risk;
- low_risk or unknown means no action 9;
- do not retry at later events.

Maximum action count remains 9.

## Measurements

Per policy pair × phase advance and overall:
- cap8 / single-critical / risk-gated losses
- ninth actions
- rescues relative to cap8
- unnecessary ninth actions on cap8 survivors
- ineffective ninth attempts
- rescue per unnecessary action where defined.

## Interpretation

Retaining most rescue with sharply fewer unnecessary ninth actions would demonstrate a closed-loop chain from native state → failure-risk discrimination → selective bounded correction. Loss of rescue would bound the classifier's control utility even if its diagnostic precision is high.

## Bounds

Counterfactual only. No threshold fitting, prototype mutation, learned weights, new native fields, later retry after low/unknown first event, action budget >9, new action type, cohort-specific tuning, cadence change, future-policy input, live activation, or production authority.
