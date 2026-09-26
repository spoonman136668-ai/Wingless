# Wingless UP-31: isospectral hidden per-copy misalignment

Status: Windows-qualified causal positive; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-30 qualified causal positive with alignment caveat sealed at `7a8343cd2056033183ff0d374c6eeccc3fcc24e9`.

## Question

UP-30 showed a very large capability loss when the six repeated transport copies were split by different per-copy phase offsets. That ablation broke both:

- exact repeated-spectrum multiplicity; and
- exact coordinate-aligned co-evolution across the six historical copies.

UP-31 isolates those explanations.

**Can the task survive when exact six-fold spectral multiplicity is preserved but the six copies no longer share the same internal coordinate basis?**

## Hidden per-copy basis change

UP-31 constructs a fixed block-diagonal unitary permutation `S`.

Each of the six 16-dimensional historical copies receives a different deterministic internal permutation.

The encoder and transport are transformed consistently:

- encoder mixer: `Q S`;
- transport: `Q S T S^H Q^H`.

Because this is a unitary conjugacy of the original repeated transport:

- the complete spectrum is unchanged;
- six-fold multiplicity is preserved exactly;
- norm preservation is unchanged;
- an exact 192-real-dimensional orthogonal representation still exists.

At runtime, the state remains one fully mixed 96-dimensional object.

The observer does **not** receive `S`, the six-copy factorization, or the hidden intertwiners.

## Alignment diagnostics

Two cross-copy shift operators are measured only for diagnosis:

1. **naive aligned shift** — the shift that commuted in the original coordinate-aligned basis;
2. **correct hidden intertwiner** — the shift conjugated by the hidden per-copy basis change.

Required diagnostics:

- naive aligned shift commutator >= 1e-4;
- correct hidden intertwiner commutator <= 1e-10.

This proves coordinate alignment is genuinely broken while the repeated-spectrum intertwiner structure still exists.

## Observer

The learned observer is the same transport-only / training-only construction that passed UP-28:

- discover 128 approximately commuting observables from the transformed transport;
- use only depth-0 training states and task labels;
- rank candidates using the full quadratic interaction model;
- select 64 runtime observables;
- regress to the already learned UP-13 2-D phase code;
- fresh 2-D softmax calibration heads.

No held-out data, explicit depth, hidden basis, hidden mixer, prototype, inverse transport, or runtime unmix is used.

## Training and evaluation

Training:

- all 128 balanced training tables;
- depth 0 only;
- two deterministic noisy trials per table;
- memory noise 0.05.

Held-out evaluation:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two deterministic noisy trials per table/depth.

Mutable integration:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only.

## Scientific gates

**Base control**

- aligned UP-28 control held-out >= 0.99;
- mutable commit >= 0.99;
- exact final >= 0.95;
- relational query >= 0.95.

**Misaligned transport**

- real orthogonality error <= 1e-10;
- exact conjugacy error <= 1e-10;
- naive aligned-shift commutator >= 1e-4;
- correct hidden-intertwiner commutator <= 1e-10.

**Discovered observer**

- discovered-bank commutator error <= 1e-5;
- feature drift <= 5e-3;
- training accuracy >= 0.99;
- held-out accuracy >= 0.99;
- mutable commit >= 0.99;
- exact final >= 0.95;
- relational query >= 0.95.

## Interpretation boundary

If UP-31 passes, exact coordinate-aligned co-evolution is not required. The stronger surviving explanation is the repeated-spectrum / commutant structure itself, because the six copies may live in different hidden internal bases and transport-only discovery still recovers the task.

If UP-31 fails while the exact-conjugacy and hidden-intertwiner diagnostics pass, repeated spectrum is representationally preserved but the current transport-only discovery/selection mechanism is not invariant enough to hidden copy-basis changes. The next experiment should improve dynamics-informed subspace discovery rather than revert to visible channel labels.

This still does not imply a uniquely complex or quantum mechanism. UP-29 already established exact real-orthogonal equivalence.

## Plain speak

UP-30 told us that making the six internal copies evolve differently destroys the memory system.

But that test changed two things at once.

UP-31 keeps the same six repeated resonances and only rotates the internal coordinate system of each copy differently.

If Wingless still works, the copies do not need to line up coordinate-for-coordinate. What matters is the deeper repeated dynamical structure and the stable measurements it creates.

## Authority boundary

UP-31 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35874862731` completed successfully at source head `95d0c7e587cf7b09134af7850983b1aeba8852ba`.

The hidden per-copy basis change was genuine: the old aligned cross-copy shift commutator rose to 0.13094732886183125, while the correctly conjugated hidden intertwiner still commuted to 9.626179587282793e-15. Exact isospectral conjugacy error was 1.3100721734176645e-14.

Despite the coordinate misalignment, the transport-only discovered observer reached 1.0 training accuracy, 0.999755859375 held-out accuracy, 1.0 mutable commit accuracy, 1.0 exact-final accuracy, and 1.0 relational-query accuracy.

Interpretation: exact coordinate-aligned co-evolution is not required. Together with UP-30, the evidence now points more specifically at repeated-spectrum multiplicity / its commutant structure as the important tested ingredient.

Next: preregister a multiplicity dose-response that progressively partitions the six repeated copies while preserving norm, state dimension, and observer protocol. A monotonic or threshold-like loss would materially strengthen that causal interpretation.
