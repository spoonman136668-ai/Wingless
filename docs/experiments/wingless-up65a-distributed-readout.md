# Wingless UP-65A — sparse learned readout from distributed representations

Status: preregistered scientific cognition/representation experiment.

Scientific parent: UP-64A seal `4bc2a1ec7418b1be091fe54d6ad0fc952f46325c`.

## Question

UP-64A showed zero-shot operator transfer into a 243-state five-role substrate. UP-65A isolates the next representation question: can a learned observer recover the five role values from a distributed code when trained on only one-third of the logical state combinations?

## Frozen design

The logical domain is the same five-role / three-value / 243-state space.

Training states satisfy:

`sum(role values) mod 3 == 0`

This yields exactly 81 training states and 162 completely held-out combinations.

Two 64-feature distributed encodings are frozen before execution:

1. **dense Hadamard factorized** — a factorized role/value code is spread across all 64 coordinates by a fixed orthogonal Hadamard transform;
2. **joint Fourier entangled** — the complete joint-state index is encoded by 32 sine/cosine frequency pairs, with no explicit role factorization.

Five independent three-class linear-softmax heads are trained per arm, one per role, with identical optimization budgets.

The factorized arm uses the existing 0.98 per-role held-out gate. The entangled arm is diagnostic and is not required to pass.

This experiment trains only the observer. It does not alter the unitary operator learner, topology, or activation authority.

## Interpretation

If the dense factorized arm generalizes, hand-aligned coordinate readout is not required: compositional state can be densely mixed and recovered from sparse combinations by a learned observer.

The entangled arm tests whether arbitrary joint-state encoding preserves that property. A factorized-pass / entangled-fail outcome would identify representational factorization as a necessary inductive bias rather than a limitation of the unitary operators themselves.


## Authoritative Windows result

Workflow run: `36049598653`

Runner: `WINGLESS-LINKDEADKB`

Source head: `0f18d99dff1959fcf0fed8cfbe3869ecb000cc39`

Artifact: `10830385546`

Artifact digest: `sha256:8ef1079b0d94d09661d6f0e562f4bb2838e43a9ff9320e7244e52f0adec79b9b`

Results:

- dense Hadamard factorized arm: mean held-out accuracy `1.0`; every role `1.0`; gate PASS;
- joint Fourier entangled arm: mean held-out accuracy `0.7876543209876543`; gate FAIL;
- entangled per-role held-out: role0 `1.0`, role1 `0.9753086419753086`, role2 `0.8888888888888888`, role3 `1.0`, role4 `0.07407407407407407`;
- both arms reached `1.0` training accuracy on the same 81 training states.

## Scientific classification

A learned linear observer can recover all five compositional role values from a dense orthogonally mixed representation without hand-aligned coordinates, provided the representation preserves factorized structure.

The fully entangled joint code can be memorized on training states but generalizes poorly to unseen combinations. Representation factorization is therefore a real inductive-bias requirement in this observer regime, not an operator-learning failure.

The next A-lane experiment holds the split and observer class fixed while increasing the Fourier bandwidth of the entangled representation to distinguish insufficient harmonic support from a deeper compositional-generalization failure.
