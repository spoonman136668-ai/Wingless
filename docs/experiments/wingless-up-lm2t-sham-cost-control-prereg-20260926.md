# Wingless UP-LM2T — matched sham specificity control

Status: preregistered counterfactual-only response-control experiment.

Scientific parent: sealed UP-LM2S 1b755475d311f70774af1d40dea695a97e9a96d9.

## Question

Did LM2S succeed because the response targeted the endangered pending dependency, or would moving any pending REPORT at the same times provide the same benefit?

## Frozen design

Reuse:
- exact recall cap 16;
- deferred levels 4,5,6;
- identity rotations 0,7;
- value shifts 0,1,2,3;
- twelve unique pressure writes.

For each arm compute the exact intervention store indices from the frozen LM2S warning-guided policy.

Compare:

1. baseline — no early closure.
2. guided — at each frozen intervention index, mark the actual FIFO-head unreported dependency as reported before the write.
3. sham — at the same intervention indices and with the same number of moved REPORTs, mark a different currently-unreported pending dependency, choosing the last such FIFO entry that is not the endangered head.

The recall store itself is unchanged by REPORT closure in all arms.

## Measurements

Per arm:
- baseline harmful unreported evictions;
- guided harmful evictions;
- sham harmful evictions;
- guided and sham interventions;
- guided and sham prevented evictions;
- moved-report counts.

Aggregate intervention specificity:
- guided prevention fraction;
- sham prevention fraction;
- guided advantage over sham.

## Bounds

Counterfactual only. Same report-action budget in guided and sham. No capacity increase, extra training, live activation, adaptive threshold, future oracle, result-informed retry, or production authority.
