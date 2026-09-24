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
