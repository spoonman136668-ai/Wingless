# Wingless UP-188B — native-state mass correction

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-187B 5647d1c20194b0e5d133055c2e25f90c7bd73518.

## Question

UP-187B showed that duplicating the entire probability table by native context does not improve held-out calibration. Can the native context instead support a small multiplicative correction to the already-strong pooled risk model?

## Frozen training data

Use only phases 26..30.

Build the same pooled UP-184B cell model over all five phases.

Freeze the native context rule from UP-187B:
- context_A if pre-cleanup native correct count <= 621;
- context_B otherwise.

For each context class in phases 26..30:
1. score every row with the pooled model;
2. sum predicted crossing mass;
3. count actual crossings;
4. freeze one multiplicative factor = actual / predicted.

No per-cell context split is created.

## Untouched held-out phases

Evaluate on new phases:
- 34;
- 35;
- 36.

Compare:
1. pooled_26_30 unchanged;
2. pooled_plus_native_mass_correction.

For the corrected model, multiply each pooled probability by its context factor and clamp to [0,1].

## Measurements

Per model × phase:
- native context and correct count;
- actual crossings;
- predicted mass;
- predicted/actual ratio;
- absolute mass-ratio error;
- Brier;
- ECE;
- AUROC;
- average precision.

Aggregate mean absolute mass-ratio error.

## Interpretation

Improvement would show that native context is useful as a low-dimensional calibration correction even though full context-specific tables overfragmented the data. Failure is valid and blocks further calibration complexity without new evidence.

## Bounds

Shadow only. No held-out fitting, phase/parity input, adaptive scaling, intervention, maintenance action, capacity change, extra model calls, live activation, or production authority.
