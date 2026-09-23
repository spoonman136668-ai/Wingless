# Wingless UP-34: task-allocated commutant capacity

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-33 qualified causal positive sealed at `b83e09b45b0c3c887778fa6c907444a36e6b3f55`.

## Question

UP-33 showed that commutant capacity is a better causal descriptor than maximum multiplicity alone in the tested architecture.

UP-34 asks the first architecture-selection question:

**Can the task itself choose a lower-cost internal symmetry organization before true held-out evaluation, using training data only and paying an explicit price for commutant capacity?**

This is a bounded bridge toward self-organized dynamics. The candidate menu is still hand-specified; the system is not yet inventing arbitrary transport operators.

## Training-only allocation split

The existing 128-table training pool is deterministically divided into:

- an inner fit set;
- an inner validation set.

The true 128-table held-out pool remains inaccessible to architecture selection.

The harness verifies that no table index appears in both the allocator pool and the true held-out pool.

## Candidate menu

The allocator can choose among these spectral partitions:

- 1+1+1+1+1+1, capacity 6;
- 2+1+1+1+1, capacity 8;
- 2+2+1+1, capacity 10;
- 3+1+1+1, capacity 12;
- 2+2+2, capacity 12;
- 4+1+1, capacity 18;
- 3+3, capacity 18;
- 6, capacity 36.

Every nontrivial partition is tested under the same three anonymous-axis assignments used in UP-33:

- identity;
- rotate-two;
- interleave.

The fully repeated capacity-36 control has only one distinct assignment.

## Fixed resource price

For each candidate, the allocator measures on inner validation:

- static held-out accuracy;
- mutable commit accuracy.

Its performance term is the weaker of those two metrics.

The frozen allocation score is:

`score = min(validation held, validation commit) - 0.02 * capacity / 36`

The winning structurally valid candidate is the highest score. Ties are broken toward lower capacity, then lexicographically by name.

This resource price is fixed before results and is not tuned after seeing the selected candidate.

## Final evaluation

After the winning candidate is locked:

1. it is retrained from scratch on the complete original training pool;
2. it is evaluated on the untouched true held-out pool;
3. the capacity-36 repeated-spectrum control is retrained/evaluated under the same final protocol.

No true held-out result can influence the architecture selection.

## Scientific gates

Inner split:

- fit/validation pools must be disjoint from true held-out;
- minimum fit marginal count >= 12;
- minimum validation marginal count >= 4.

Resource allocation:

- selected commutant capacity must be below 36.

Final task performance:

- selected held-out accuracy >= 0.985;
- selected mutable commit accuracy >= 0.95;
- selected exact-final accuracy >= 0.90;
- selected relational-query accuracy >= 0.95.

Retention relative to full capacity:

- held-out loss <= 0.02;
- mutable-commit loss <= 0.05.

The task-allocation gate is positive only if all conditions pass.

## Interpretation boundary

A positive result would show that, given a bounded menu of internal symmetry organizations, training evidence plus a fixed capacity cost can select a smaller commutant that retains most or all task capability without consulting the true held-out set.

That would be the first step from "we specify the useful symmetry" toward "the task allocates how much stable symmetry it needs."

It would not yet show continuous self-organization or invention of new transport structure.

A negative result remains informative. It would indicate that either the bounded selector is insufficient, the chosen resource price is inappropriate under the frozen protocol, or useful symmetry cannot yet be compressed without hidden cost.

UP-29 already established exact real-orthogonal equivalence, so no outcome here supports uniquely complex or quantum-computation claims.

## Plain speak

Until now, we have been deciding how much stable internal symmetry the system gets.

UP-34 gives it several possible internal organizations and charges it for using a bigger one.

It can look only at training performance.

Then it has to choose.

Only after that choice is frozen do we test it on data it has never seen.

If it chooses a smaller internal symmetry space and still performs almost as well as the full one, that is our first real step toward the system allocating its own internal memory geometry.

## Authority boundary

UP-34 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
