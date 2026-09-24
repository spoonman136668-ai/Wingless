# Wingless research-harness performance equivalence R1

Status: staging only. This branch is intentionally outside `research/wingless-up**` and therefore does not auto-run the Windows scientific qualification workflow.

## Goal

Reduce wall-clock time of future Wingless UP experiments without changing scientific results.

## Non-negotiable acceptance rule

A performance change is accepted only if the optimized path reproduces the reference path for the same fixed candidate offsets, seeds, tables, depths, and hyperparameters.

The proof must compare all scientific fields that can affect trajectory or classification, including:

- scalar objective;
- component phase/value/relation signals;
- selected observable indices;
- soft capacity and resource penalty;
- gradient directional derivative;
- proposed update and resulting offsets;
- provisional symmetry groups/capacity/gaps;
- checkpoint selection;
- final held-out capability metrics;
- scientific diagnosis and classification.

Where the computation order is unchanged, require exact equality. If parallel execution changes only floating-point reduction order, it is not accepted in R1; parallelism is deferred to a separately qualified step.

## R1 optimization scope

Only two deterministic, candidate-offset-invariant computations are eligible:

1. `taskSelectedTrainingStates(fitTables, mixer, memoryNoise, 2)`
2. the relation head produced by `relationHeadTrainingSamples()` followed by `trainLinearSoftmax(..., 4, 16, 600, 1.0)`

Both are currently reconstructed inside repeated candidate evaluations even though neither depends on candidate spectral offsets.

R1 may build these once per scientific probe and pass immutable copies/references into candidate evaluations.

R1 must not cache or reuse:

- candidate commuting observables;
- candidate selected observables;
- candidate depth operators;
- candidate decoder/regressor training;
- candidate rollout distributions;
- candidate objective values;
- candidate gradients;
- held-out results.

## Qualification design

The reference evaluator remains untouched.

The optimized evaluator is implemented alongside it.

Use a frozen set of representative candidate offsets drawn from existing preregistered trajectories, including:

- initial singleton offsets;
- a perturbation pair;
- a non-fusion checkpoint;
- a reversible fusion-basin checkpoint.

For every candidate, run reference then optimized evaluation and compare the complete scientific evaluation object after normalizing only non-scientific labels/names.

Then run one complete rolling trajectory through both implementations and require identical:

- 48 directional derivatives;
- 48 update vectors;
- 48 output offset vectors;
- checkpoint evaluations at 0/12/24/48;
- selected checkpoint;
- final hard held-out metrics.

Benchmark only after equivalence passes.

## Rejection rules

Reject the optimization if:

- any scientific scalar differs;
- any selected observable index differs;
- any offset/update/gradient differs;
- checkpoint choice differs;
- held-out metrics differ;
- deterministic replay differs;
- performance gain is absent or negligible enough not to justify maintenance risk.

No sample reduction, gate relaxation, changed seed, changed horizon, changed checkpoint schedule, or changed objective is permitted.

## Plain speak

This pass is allowed to stop recomputing facts that are literally the same every time. It is not allowed to approximate, skip evidence, or change the order of scientific decisions.

The old evaluator remains the referee. The faster evaluator only becomes usable after proving it gives the same experiment back.
