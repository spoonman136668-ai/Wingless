# Wingless UP-33: equal commutant-capacity partition control

Status: Windows-qualified causal positive; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-32 ordered graded multiplicity trend sealed at `e97ad901dec440d6ffed209504363f79ff2e70ff`.

## Question

UP-32 reduced repeated-spectrum multiplicity and observed an ordered degradation trend, but its partitions changed two quantities at once:

- maximum repeated multiplicity;
- copy-axis commutant capacity `C = sum(m_i^2)`.

UP-33 asks which quantity better explains capability.

## Matched-capacity pairs

Two equal-capacity pairs are tested:

- `4+1+1` versus `3+3`: both have `C=18`, but maximum multiplicity differs (4 vs 3);
- `3+1+1+1` versus `2+2+2`: both have `C=12`, but maximum multiplicity differs (3 vs 2).

A six-copy repeated-spectrum control is also reproduced.

Every non-control arm has:

- zero-mean offsets;
- RMS offset magnitude 0.05 radians per step;
- norm-preserving transport;
- an exact real-orthogonal equivalent;
- the same 96-complex / 192-real scalar state budget.

## Assignment robustness

Each non-control partition is evaluated under three fixed assignments of the phase groups to the anonymous copy axis:

- identity;
- rotate-two;
- interleave.

The observer is not given those assignments.

## Observer protocol

Every variant receives the same budget:

- discover 128 approximately commuting observables from that variant's transport;
- depth-0 training states only;
- interaction-aware training-label-only selection;
- 64 runtime observables;
- ridge regression to the previously learned 2-D phase code;
- fresh 2-D calibration heads.

No held-out data, explicit depth, runtime unmixing, hidden partition labels, or prototype lookup is used.

## Scientific gates

Structural validity requires every tested arm to satisfy:

- real orthogonality error <= 1e-10;
- selected-bank commutator error <= 1e-5;
- feature drift <= 5e-3;
- non-control RMS offset = 0.05 within 1e-12.

The six-copy control must satisfy the existing static and mutable acceptance gates.

Equal-capacity equivalence uses frozen tolerances:

- median held-out difference <= 0.03;
- median mutable-commit difference <= 0.05.

This is tested independently for the `C=18` and `C=12` pairs.

Capacity ordering requires the mean across both `C=18` partitions to exceed the mean across both `C=12` partitions for both held-out accuracy and mutable commit accuracy.

The capacity hypothesis is supported only if structural validity, base control, both equal-capacity equivalence gates, and capacity ordering all pass.

## Interpretation boundary

A positive result would indicate that the size of the transport commutant is a better causal description than maximum multiplicity alone.

A negative result would show that partition shape, maximum multiplicity, or assignment-specific geometry matters beyond `sum(m_i^2)`.

Raw per-assignment measurements remain authoritative even if the summary gate is negative.

UP-29 already established exact real-orthogonal equivalence. No outcome here supports uniquely complex or quantum-computation claims.

## Plain speak

Two different ways of grouping the six resonances can create the same amount of stable internal symmetry.

UP-33 compares those matched pairs.

If systems with the same symmetry capacity behave similarly even when their largest repeated group is different, then the important thing is likely the amount of stable internal room the dynamics leave available—not simply “how many copies match.”

## Authority boundary

UP-33 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35881192286` completed successfully at source head `5aaf6a54278e1744e35d919487af80e9458a8c79`.

Both equal-capacity equivalence controls passed. For capacity 18, the median held-out difference between `4+1+1` and `3+3` was 0.00244140625 and the median mutable-commit difference was 0.0. For capacity 12, the corresponding differences were 0.00537109375 and 0.00390625.

The higher-capacity pair also outperformed the lower-capacity pair on both aggregate metrics: mean held-out accuracy 0.9974365234375 versus 0.985107421875, and mean mutable-commit accuracy 0.9963107638888888 versus 0.9233940972222221.

Interpretation: in the tested architecture, commutant capacity is a better causal descriptor than maximum multiplicity alone. The next research frontier is to determine whether useful commutant capacity can be discovered, induced, or allocated by the system itself rather than hand-specified through repeated-spectrum structure.
