# Wingless UP-186B — native-context diagnostic

Status: preregistered scientific shadow diagnostic.

Scientific parent: sealed UP-185B 0316209ffd00c60d12a89066935e7bf1e26e2234.

## Question

UP-183B/184B/185B established a repeating cross-phase calibration residual that is not explained by external odd/even parity. Is there a corresponding difference in Wingless's own frozen pre-cleanup state?

## Frozen phases

Inspect phases 26 through 33 inclusive. No phase is used to fit a predictor.

For each phase, build the same six-subject post-REPORT states used by the B-lane cleanup experiments, before either cleanup path begins.

## Native summaries

Aggregate only quantities already present in the gate state or its native old-example decision surface:
- mean signed old-example margin;
- mean absolute old-example margin;
- minimum absolute old-example margin;
- count of old examples with |margin| <= 0.05;
- count of currently correct old examples;
- gate weight L2 norm;
- the three gate biases.

Also execute the already accepted two 12-step cleanup paths and record the total number of next-update correctness crossings for that phase.

## Diagnostic comparison

No threshold or classifier is fitted.

Report whether phases with equal crossing totals also have identical or systematically separated native summaries, and whether any summary exhibits the same repeating phase-class structure.

## Interpretation

A native summary that tracks the repeating crossing class becomes a candidate context feature for later calibration. Failure to find such a summary is valid and blocks adding one by guesswork.

## Bounds

Shadow only. No predictor fitting, threshold tuning, maintenance action, capacity change, extra model calls, live activation, or production authority.
