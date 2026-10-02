# Wingless UP-16: single composite recurrent state

Status: Windows-qualified single-composite-state result; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-15 Windows learned-anchor qualification sealed at `33162a84dfd91e405564f3812303e0d49bf9ed8b`.

## Question

Can the complete learned memory/reference system be carried and evolved as **one recurrent complex state object** instead of one memory state plus five separately materialized reference vectors?

This is a structural-fold experiment.

It does **not** yet claim that the six logical subspaces have disappeared. Their locations inside the composite state remain explicit.

## Qualified geometry being frozen

UP-16 freezes all learned geometry from the prior lineage:

- UP-13 learned value-phase alphabet;
- UP-14 learned pilot support matrix;
- UP-15 learned dense complex anchor;
- rank 4;
- four pilot references plus one anchor.

No geometry is retrained in UP-16.

## Composite state

The runtime state has dimension:

`6 × 16 = 96 complex coordinates`.

The six logical inputs are:

1. mutable memory;
2. learned anchor;
3. learned pilot 0;
4. learned pilot 1;
5. learned pilot 2;
6. learned pilot 3.

Every input state is normalized and scaled by `1/sqrt(6)`, then packed into one unit-norm 96-dimensional state.

After packing:

- runtime state-object count = 1;
- separate frame vectors are not materialized during transport/readout;
- the six channel subspaces remain predefined.

## Evolution

The single state evolves under the direct-sum operator:

`I_6 tensor U`

for the unitary path, and the matched direct-sum non-unitary operator for the control.

Implementation applies the same 16-D transport block to each internal subspace, which is mathematically one block-diagonal operator over the 96-D state.

## Readout

The observer reads gauge-invariant correlations directly from internal subspaces of the single composite state.

For entity `e`, it uses:

- memory subspace;
- anchor subspace;
- pilot-`e` subspace.

The same two-real-feature observable qualified in UP-15 is reconstructed from the packed state.

No separate runtime frame object, explicit depth, inverse transport, or prototype lookup is supplied.

## Equivalence gate

UP-16 directly compares the observable computed from:

- the historical separate memory + frame states;
- the new single composite state.

Depths checked:

`1, 8, 128, 1024`.

Maximum observable error must be <= `1e-12`.

This prevents a false structural pass caused by silently changing the semantics of the reference measurement.

## Static evaluation

Fresh entity decoders train using:

- all 128 balanced training tables;
- depths 8, 24, 72, 216, 432, 648;
- memory noise 0.05;
- memory-only global-phase nuisance;
- four deterministic noisy trials per table/depth;
- 1,200 linear-softmax steps.

Evaluation uses:

- all 128 disjoint held-out tables;
- unseen depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary path gets the same composite packing and decoder budget.

## Mutable integration

The unitary composite state then runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- memory noise 0.05;
- memory-only global-phase nuisance;
- explicit irreversible overwrite/re-encode boundary;
- learned relation head.

Each segment begins by packing the current canonical memory and frozen learned frame into one composite state. Only that packed state is evolved and observed.

## Scientific gates

**Observable-equivalence pass**

- maximum separate-vs-composite feature error <= 1e-12.

**Static-fold pass**

- observable equivalence passes;
- unitary unseen-depth held-out accuracy >= 0.99;
- maximum composite norm drift <= 1e-12.

**Mutable-integration pass**

- commit accuracy >= 0.99;
- exact final-table accuracy >= 0.95;
- relational-query accuracy >= 0.95;
- maximum composite norm drift <= 1e-12.

## Interpretation boundary

A positive UP-16 result means the qualified memory/reference mechanism can be represented and evolved as one recurrent latent-state object without changing its observable behavior.

It does **not** mean the frame topology has been eliminated. The six internal subspaces are still known and the readout knows which subspace corresponds to memory, anchor, and each pilot.

That is the next boundary. A positive UP-16 permits UP-17 to apply a dense channel mixing / learned composite encoder so the runtime state no longer exposes fixed memory/anchor/pilot channel boundaries to the observer.

## Authority boundary

UP-16 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-23 against source head `91700cdc25e2ac538742ebd95b31010ee6dc1e2a` passed the harness and complete Wingless regression.

The packed 96-dimensional unitary state matched the historical separate-state observables with maximum error 1.176836406102666e-14 across depths 1, 8, 128, and 1024.

Fresh unitary decoders reached 1.0 train and 1.0 unseen-depth held-out accuracy with maximum composite norm drift 3.175237850427948e-14. The matched non-unitary composite path reached 0.763427734375 held-out accuracy and maximum norm drift 2211700831922.4956.

Mutable integration used 48 scenarios × 16 writes and achieved 1.0 commit decode, 1.0 exact final-table, and 1.0 relational-query accuracy. Minimum value margin was 0.17822707126915677.

Interpretation: memory plus the learned frame can be carried as one recurrent state object without changing observable behavior. The remaining explicit structure is the six known internal channel subspaces.
