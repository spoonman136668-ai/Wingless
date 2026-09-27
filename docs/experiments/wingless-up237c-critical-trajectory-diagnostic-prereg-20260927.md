# Wingless UP-237C — first-critical native trajectory diagnostic

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-236C 397f7a0b3e75d7c8519e9575b5459d4d7eeec9d3.

## Question

When the frozen static Hamming classifier misses a failing arm, does the one-monitoring-step native trajectory into the first post-action8 critical event contain separable information?

## Frozen environment

Use exactly the UP-236C heldout conditions:
- cohorts 0,9,10,15 / 1,4,11,14 / 2,5,8,13 / 3,6,7,12
- phase advances 0 and 22
- four frozen policy pairs
- horizon 80
- monitoring cadence 2
- exact cap-8 policy with action-5 retiming
- no action 9.

512 arms.

## Frozen observations

At action 8:
- capture the post-action8 native snapshot.

At every later monitoring boundary:
- capture the same nine native fields used by C231-C236.

At the first later adversarial-critical event:
- retain current snapshot;
- retain the immediately previous monitoring snapshot;
- compute signed one-step deltas for all nine fields.

No intervention occurs.

## Signatures

Current-state signature:
- exact nine-field C231 signature.

Trajectory signature:
- current-state signature plus signed deltas from the previous monitoring snapshot.

Report distinct failure/survivor signatures and shared-signature counts for both representations.

## Interpretation

Fewer or zero shared trajectory signatures would show that native dynamics contain failure information absent from the static snapshot and justify a preregistered trajectory-risk classifier. Persistent overlap would reject this one-step trajectory as sufficient.

## Bounds

Shadow diagnostic only. No corrective action after action 8, no threshold fitting, no classifier training, no adaptive feature selection, no new native field beyond one-step deltas of existing fields, no cadence change, no future-policy input, no live activation or production authority.
