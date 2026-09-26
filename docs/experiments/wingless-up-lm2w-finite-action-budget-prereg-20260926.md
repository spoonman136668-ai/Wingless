# Wingless UP-LM2W — finite corrective-action budget

Status: preregistered scientific counterfactual response experiment.

Scientific parent: sealed UP-LM2V d639f6e56f3cb923daa3c5d91da173b8a3eea9f0.

## Question

When endangered dependencies outnumber available corrective moves, does the frozen native warning-guided rule spend a finite action budget only on genuine hazards, and how much failure remains after budget exhaustion?

## Frozen system

Reuse UP-LM2V unchanged:
- exact recall cap 16;
- deferred levels 4, 5, 6;
- identity rotations 3 and 11;
- value shifts 0 and 2;
- pending layouts suffix_reported and alternating_reported;
- same 12 incoming unique stores;
- no future schedule input.

## Action budgets

Test hard per-arm budgets:
- 1 corrective move;
- 2 corrective moves;
- 3 corrective moves.

Guided arm:
- before each write, detect whether the FIFO victim is still unreported;
- if harmful and budget remains, close exactly that endangered dependency;
- once budget is exhausted, take no further action.

Sham arm:
- same action budget and same harmful trigger times;
- move a different pending report using the frozen UP-LM2V sham target.

Baseline:
- no action.

## Measurements

Per stream × budget:
- baseline/guided/sham failed dependencies;
- actions spent;
- prevented failures;
- residual failures after budget exhaustion;
- prevention per guided action;
- whether any guided action fired without a genuine current hazard.

## Interpretation

A bounded response is useful only if finite actions buy predictable reductions without unnecessary spending. Failure reduction that scales with available budget supports a cost-accountable response; nonlinearity or wasted actions defines the boundary.

## Bounds

Counterfactual only. No live activation, capacity increase, extra training, adaptive budget, semantic priority, future oracle, or result-informed retry.
