# Wingless UP-LM2U — downstream obligation completion

Status: preregistered counterfactual-only response-cost experiment.

Scientific parent: sealed UP-LM2T 73b5d5a6d7025a4e97c4eed06b50d9dc502d3417.

## Question

After warning-guided early closure prevents the endangered unreported eviction, does it preserve total downstream obligation completion, or merely move the failure to another pending dependency?

## Frozen design

Reuse LM2T exactly:
- recall cap 16;
- deferred levels 4,5,6;
- rotations 0,7;
- value shifts 0,1,2,3;
- twelve unique pressure writes;
- frozen guided intervention indices;
- matched sham intervention indices and action counts.

Compare baseline, guided, and sham.

## Outcome

After all 12 pressure writes, evaluate all twelve original first-chunk dependencies.

A dependency counts as completed if:
- its REPORT was already closed earlier in that arm; or
- its value is still present in recall and can be reported at end of stream.

Measure:
- completed obligations out of 12;
- failed obligations;
- remaining pending reports still recallable;
- earlier-closed reports;
- action counts.

No further writes occur during end-of-stream evaluation.

## Bounds

Counterfactual only. Same action budget for guided and sham. No live activation, capacity increase, extra training, future oracle, adaptive retuning, result-informed retry, or production authority.
