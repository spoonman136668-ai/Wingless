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


## Windows qualification result

R1 was evaluated twice on the dedicated Windows runner.

### Attempt 1 — run 35940424574

The scientific equivalence proof passed completely:

- reference SHA-256: `370afca2ad1ae23235906b938604e6c7bdf9db5e2375e1d425d232c7f53ecc42`;
- optimized SHA-256: identical;
- optimized replay SHA-256: identical;
- exact full-result equivalence: passed;
- optimized deterministic replay: passed;
- reference duration: `537509 ms`;
- optimized duration: `511749 ms`;
- measured speedup: `1.0503360711063274x`;
- frozen minimum speedup: `1.05x`.

The test script emitted `WINGLESS_HARNESS_PERF_EQUIVALENCE_R1_PASS`, but the outer qualification wrapper rejected the run because its generic marker parser did not recognize the custom harness marker. That is an infrastructure classification defect, not a scientific mismatch.

### Attempt 2 — run 35950581567

The marker integration was repaired and the equivalence checks again progressed through schema, scope, exact-equivalence, determinism, and hash checks. The run then failed only at the frozen performance gate:

- measured speedup: `1.0011856315194536x`;
- required speedup: `1.05x`;
- failure: `HPERF_SPEEDUP_GATE_FAILED`.

Because R1 preregistered a minimum speedup of 1.05x and the repeated measurement did not reproduce that benefit, the optimization is **not adopted**.

## Final classification

**Scientific equivalence supported; performance benefit not reproducibly qualified.**

The invariant-reuse implementation demonstrated that the two targeted computations can be reused without changing the complete UP-44 scientific result. However, the measured wall-clock gain was not stable enough to clear the frozen maintenance threshold on repeat qualification.

Therefore future scientific experiments continue on the untouched reference evaluator. The cached evaluator remains research-only and must not be substituted into the scientific lineage unless a later independently preregistered performance pass demonstrates a reproducible gain.

## Plain speak

The shortcut was safe, but it was not consistently faster.

One run barely cleared the speed target; the repeat was essentially the same speed as the original. Since we promised not to trade reliability for complexity, we keep the old evaluator for science and move on.
