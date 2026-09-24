# Wingless harness performance equivalence R1

Status: preregistered performance qualification only. This is not a scientific UP result and does not change any scientific threshold, sample, seed, horizon, checkpoint, or held-out boundary.

Baseline scientific seal: UP-44 `584373c3bd32fc305cffffeec16ad5e7ee07878c`.

## Purpose

Reduce repeated computation in the UP-44 evaluator without changing its scientific result.

## Allowed optimization

Only two candidate-offset-invariant computations may be reused across candidate evaluations:

1. task-selection states from `taskSelectedTrainingStates(fitTables, mixer, memoryNoise, 2)`;
2. the fixed relation head trained from `relationHeadTrainingSamples()` with the unchanged `trainLinearSoftmax(..., 4, 16, 600, 1.0)` call.

Everything dependent on candidate offsets remains recomputed for every evaluation, including observable discovery, selected observables, depth operators, decoder/regressor training, rollouts, objective components, gradients, updates, geometry measurements, checkpoint selection, and held-out evaluation.

## Equivalence protocol

Run:

1. one untouched `RunUP44()` reference trajectory;
2. one cached-invariant trajectory;
3. a second independently initialized cached-invariant trajectory.

Acceptance requires:

- `reflect.DeepEqual(reference, optimized)`;
- `reflect.DeepEqual(optimized, optimizedReplay)`;
- equal SHA-256 of the complete JSON scientific result objects;
- unchanged selected checkpoint, objective, trajectory, held-out metrics, diagnosis, and classification;
- full repository regression green;
- wall-clock speedup of at least `1.05x` for the first optimized trajectory versus reference.

The `1.05x` maintenance threshold is frozen before measurement.

## Rejection

Any scientific mismatch rejects the optimization even if it is faster.

No approximate equality is accepted in R1. No parallel floating-point reductions are introduced. No reduction in samples, seeds, tests, gates, horizon, checkpoints, or held-out scope is allowed.

## Plain speak

The fast path is only allowed to remember two facts that are identical every time. It still performs the actual experiment normally.

We run the old and new versions and compare the entire result. If even one scientific value changes, the speedup is thrown away.
