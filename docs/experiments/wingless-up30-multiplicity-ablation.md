# Wingless UP-30: orthogonal transport multiplicity ablation

Status: Windows-qualified causal positive with alignment caveat; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-29 qualified real-orthogonal equivalence pass sealed at `e787e506f6a8b31d703bb8e735c82bf8fb09334c`.

## Question

UP-29 established that the accepted capability is exactly reproducible as ordinary real orthogonal dynamics. Complex numbers are therefore not the causal explanation.

UP-23 through UP-28 instead point to a more specific mechanism: repeated transport eigenspaces support a rich commuting observable algebra that can preserve task-relevant internal orientation.

UP-30 directly tests that mechanism:

**If the repeated six-way transport spectrum is split while norm preservation, reversibility, state dimension, encoder, full mixing, data, discovery, selector, and decoder budgets remain fixed, does the memory capability materially degrade?**

## Base control

The base arm is the accepted UP-28 geometry:

`T = I_6 tensor U`

followed by the same full 96-dimensional coordinate mixing.

It uses the frozen UP-28 64-observable bank and decoder pipeline and must reproduce the accepted held-out/mutable performance.

## Multiplicity ablation

The ablation arm preserves six independent 16-dimensional unitary blocks but adds a distinct fixed global rotation to each block per transport step:

- channel 0: 0.000 rad;
- channel 1: +0.017 rad;
- channel 2: -0.029 rad;
- channel 3: +0.043 rad;
- channel 4: -0.061 rad;
- channel 5: +0.079 rad.

Each block remains unitary.

The full 96-dimensional transport is still conjugated into the same fully mixed latent basis.

By UP-29 realification, both the base and split transports have exact real orthogonal equivalents.

## Symmetry diagnostic

A one-channel cyclic shift operator commutes with the repeated base transport.

It should no longer commute after the per-channel phase split.

Preregistered symmetry-break gate:

- base cross-channel shift commutator <= 1e-10;
- split cross-channel shift commutator >= 1e-4.

## Observer isolation

The ablation construction uses the known six-way factorization only to create the diagnostic intervention.

The split observer does **not** receive:

- channel identities;
- split offsets;
- hidden mixer;
- semantic roles;
- the base observable bank.

Instead it reruns the same accepted discovery/selection pipeline:

- 128 deterministic dense transport-discovered candidates;
- 20 binary conjugation averaging rounds;
- interaction-aware training-only selection;
- 64 runtime observables;
- complete degree-2 feature map;
- UP-13 2-D phase-code supervision;
- ridge lambda 1e-6.

## Training and evaluation

Both arms use:

- depth-0 training only;
- all 128 balanced training tables;
- two deterministic noisy trials;
- memory noise 0.05;
- memory-only global-phase nuisance.

Held-out:

- all 128 disjoint tables;
- depths 32, 128, 512, 1024;
- two trials per table/depth.

Mutable:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- same irreversible overwrite/re-encode boundary.

## Scientific gates and diagnosis

**Base control**

- held-out >= 0.99;
- mutable commit >= 0.99;
- exact final >= 0.95;
- relation >= 0.95.

**Norm-preserving geometry**

- base realified orthogonality error <= 1e-10;
- split realified orthogonality error <= 1e-10.

**Symmetry break**

- base cross-shift commutator <= 1e-10;
- split cross-shift commutator >= 1e-4.

**Split discovery quality**

- selected-observable commutator <= 1e-5;
- feature drift <= 5e-3.

**Material ablation effect**

Supported if either:

- held-out accuracy drops by at least 0.10, or
- mutable commit accuracy drops by at least 0.10,

while the base, orthogonality, symmetry-break, and split discovery/invariance controls pass.

The split arm's training, unseen-depth, and mutable gates are reported separately rather than used to force the expected result.

## Interpretation boundary

If the material ablation criterion passes, the evidence supports dependence of the current architecture on repeated-spectrum/multiplicity structure, not on complex arithmetic itself.

This would not establish that multiplicity is universally necessary for memory. It would establish that it is causally important in this tested architecture and task family.

If the split arm remains strong, the multiplicity diagnosis is incomplete and the next experiment must search for the actual invariant mechanism rather than preserving the hypothesis.

No post-result change to offsets, discovery, data, or thresholds is authorized.

## Plain speak

UP-29 showed that complex numbers were only a convenient way to write the system.

UP-30 now asks what part of the geometry actually matters.

The working system has six copies of the same internal rotation. That repetition creates stable ways for different parts of the state to compare themselves.

UP-30 slightly changes the rotation speed of each copy while keeping every rotation perfectly reversible.

If memory breaks even though nothing loses norm or reversibility, we have strong evidence that the repeated internal rhythm—not “quantum phase”—is the important mechanism.

## Authority boundary

UP-30 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification completed successfully on workflow run `35873650226` at source head `fb0ce8161e2fa16afbaebd43dee5cb516d5033f1`.

The repeated-spectrum control reproduced the UP-28 frontier exactly: 1.0 held-out accuracy, 1.0 commit decode, 1.0 exact-final accuracy, and 1.0 relational-query accuracy.

The split transport remained norm-preserving and real-orthogonally equivalent, and transport-only observable discovery still produced approximately commuting, depth-stable observables. However, held-out accuracy fell to 0.461181640625 and mutable commit accuracy fell to 0.0078125. The preregistered material-effect and multiplicity-dependence gates passed.

Interpretation boundary: the split ablation breaks both repeated eigenvalue multiplicity **and** exact shared co-evolution/alignment across the six historical copies. Therefore UP-30 supports dependence on the repeated/aligned transport symmetry, but it does not yet isolate which part is causal.

The next experiment must preserve the repeated spectrum exactly while conjugating each copy into a different internal basis. That separates "same eigenvalues/multiplicity" from "same coordinate-aligned dynamics."
