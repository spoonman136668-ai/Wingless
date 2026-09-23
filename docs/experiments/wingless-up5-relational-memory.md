# Wingless UP-5: relational read/write memory

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-4 Windows qualification sealed at `5b1224cb55ec431e196e0e42b12cbfdd1b431734`.

## Question

Can the unitary transport substrate support a mutable memory workload when true overwrite is handled by an explicit bounded irreversible commit boundary?

This experiment deliberately does **not** pretend that overwrite itself is unitary.

## Why an irreversible boundary is required

A true write such as:

`entity B := value 0`

must map several possible old B values to the same new value. That is many-to-one. A closed unitary transformation is reversible and therefore cannot erase the displaced value unless that information is retained somewhere else.

UP-5 therefore tests the hybrid architecture proposed earlier:

`unitary transport -> observe/commit -> unitary transport`

The write boundary is identical for both the unitary and matched non-unitary paths.

## Memory

The memory contains:

- four entities;
- four possible values per entity;
- sixteen complex latent dimensions;
- one active basis coordinate per entity;
- 256 possible complete memory tables.

A canonical memory table has four equal amplitudes and unit norm.

## Workload

There are 32 deterministic scenarios.

Each scenario contains:

- an initial four-entity table;
- 12 writes, including repeated overwrites of earlier entities;
- distractor/transport gaps of 8, 32, 128, or 512 propagation blocks;
- deterministic bounded complex perturbation with amplitude 0.10 before each transport segment;
- a final query involving two entities.

That yields 384 intermediate commit decodes per model plus 32 final memory queries.

## Commit/write boundary

At each write boundary:

1. the transported latent state is compared against all 256 clean transported memory prototypes;
2. the most likely complete table is observed;
3. exactly one entity is overwritten;
4. the resulting table is re-encoded into canonical latent form.

This operation is explicitly irreversible and is shared by both paths.

## Final query

The final state is decoded back to a complete table.

Two scores are reported:

- **exact final-table accuracy** — all four entity values must be correct;
- **relational query accuracy** — answer `(value[A] - value[B]) mod 4` for a deterministic pair of entities.

The relational arithmetic is a deterministic evaluation rule, not a learned language/reasoning head.

## Compared transports

Both paths use the same sixteen-dimensional state, the same 32-coupling butterfly topology, the same coupling scalars, the same perturbations, the same write boundary, and the same workload.

Only the transport differs:

- unitary: complex Givens rotations;
- matched control: unconstrained pair mixers.

## Pass conditions

For the unitary path:

- all 384 commit decodes are correct;
- exact final-table accuracy is 1.0;
- relational query accuracy is 1.0;
- minimum decode margin is greater than 0.1;
- maximum norm drift is <= 1e-12.

For the matched control:

- all reported metrics must remain finite;
- it is not required to fail.

Additionally:

- the overwrite test must prove that two distinct prior memories can map to the same committed memory;
- complete probe output must be deterministic;
- the full existing Wingless regression must remain green after advisory ICE rebuild.

## Interpretation boundary

A pass would establish a bounded hybrid memory primitive: reversible/norm-preserving transport between explicit irreversible write/commit boundaries, with repeated overwrite, long distractor intervals, perturbation, associative recovery, and a simple relational query.

This is closer to working memory than UP-4, but it is still not language reasoning.

A positive UP-5 permits UP-6: learn the observation/query interface rather than using exhaustive prototype matching and deterministic relational arithmetic, then test generalization to unseen entity/value combinations and longer write/query programs.

## Authority boundary

UP-5 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.
