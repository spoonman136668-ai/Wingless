# Wingless UP-LM2W — bounded response budget frontier

Status: preregistered scientific counterfactual response experiment.

Scientific parent: sealed UP-LM2V d639f6e56f3cb923daa3c5d91da173b8a3eea9f0.

## Question

How does warning-guided prevention degrade when the corrective action budget is smaller than the number of endangered dependencies?

## Frozen workload

Reuse the disjoint LM2V conditions:
- deferred levels 4,5,6;
- identity rotations 3,11;
- value shifts 0,2;
- pending layouts suffix_reported and alternating_reported;
- exact recall cap 16;
- 12 unique pressure writes.

## Frozen action budgets

Per arm:
- 0 actions;
- 1 action;
- 2 actions;
- 3 actions.

Guided policy: when the current FIFO victim is unreported and budget remains, close that exact endangered dependency.

Sham policy: at the same opportunity and same budget, close a different unreported dependency.

No action occurs after the budget is exhausted.

## Measurements

Per budget:
- baseline failures;
- guided failures;
- sham failures;
- guided/sham action counts;
- failures prevented;
- prevented failures per action.

## Interpretation

A useful bounded response should retain target specificity and provide positive prevention per action even when budget-constrained. This experiment maps the cost/benefit frontier; it does not authorize live use.

## Bounds

Counterfactual only. No future schedule input, capacity increase, extra training, adaptive budget, live activation, or production authority.
