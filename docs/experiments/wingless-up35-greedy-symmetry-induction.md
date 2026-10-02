# Wingless UP-35: greedy task-driven symmetry induction

Status: Windows-qualified scientific positive; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-34 qualified scientific positive sealed at `46403f5dbbbc73327f37ab392ea9db849c5eda4b`.

## Question

UP-34 showed that a task-only allocator can choose a lower-cost symmetry organization from a bounded menu and preserve full true held-out capability.

UP-35 removes the finished partition menu.

It asks:

**Can the task create useful repeated-spectrum structure incrementally, starting from minimum commutant capacity, by accepting only pairwise symmetry merges whose validation benefit justifies their added resource cost?**

## Starting point

The system begins with six distinct spectral directions:

`1 + 1 + 1 + 1 + 1 + 1`

This has minimum copy-axis commutant capacity:

`C = 6`.

No higher-capacity partition is supplied as the answer.

## Allowed developmental move

The only architecture mutation is:

**merge two currently distinct spectral groups.**

If groups of sizes `a` and `b` are merged, commutant capacity increases by:

`2ab`.

At each step every possible pairwise merge is evaluated using training-only evidence.

The best candidate is accepted only if:

`new_score - current_score >= 0.005`.

Otherwise induction stops.

## Frozen objective

UP-35 uses the same resource-price objective introduced in UP-34:

`score = min(validation held, validation commit) - 0.02 * capacity / 36`.

The capacity price and minimum accepted gain are frozen before results.

## Data discipline

The existing 128-table training pool is deterministically divided using the balanced observer hash ordering into:

- 96 inner-fit tables;
- 32 inner-validation tables.

The true 128-table held-out pool is unavailable during the entire merge path.

Every accepted merge is determined solely from the inner-fit/validation data.

Only after induction stops is the selected structure retrained on the complete original training pool and exposed once to the true held-out pool.

## Observer protocol

Every merge candidate receives the same observer machinery:

- full 96-dimensional coordinate mixing;
- transport-derived commuting-observable discovery;
- 128 candidate observables;
- interaction-aware training-label-only selection;
- 64 runtime observables;
- phase-code supervision;
- no held-out selection;
- no explicit depth;
- no runtime unmix;
- no known factorization.

## Scientific gates

Preflight:

- inner fit/validation split remains balanced;
- no true held-out leakage.

Developmental path:

- at least one merge is accepted;
- selected capacity must be > 6 and < 36;
- induction must stop because the frozen gain rule rejects the next merge, not because the search fails.

Final task performance:

- held-out accuracy >= 0.985;
- mutable commit accuracy >= 0.95;
- exact final table accuracy >= 0.90;
- relational query accuracy >= 0.95.

Retention relative to full capacity:

- held-out loss <= 0.02;
- mutable-commit loss <= 0.05.

## Interpretation boundary

A positive result would be stronger than UP-34.

UP-34 selected among complete candidate organizations supplied in advance.

UP-35 starts from the minimum-symmetry state and constructs larger symmetry groups step by step under task pressure.

That would demonstrate a bounded form of **task-driven symmetry induction**.

It would still not be continuous self-organization: the allowed mutation operator is a discrete pairwise merge and the six starting directions are fixed.

A negative result would tell us whether menu-level selection worked only because globally useful partitions had to be proposed in advance.

UP-29 already established exact real-orthogonal equivalence, so no result here supports uniquely complex or quantum-computation claims.

## Plain speak

UP-34 let the system choose a house from a catalog.

UP-35 gives it six separate rooms and one construction tool: it may knock down a wall between any two rooms.

Every time it knocks down a wall, it has to prove that the extra shared space helps the task enough to pay for it.

When no wall removal is worth the cost, construction stops.

If the resulting layout still performs almost perfectly on unseen data, the system has begun to **build its own useful internal symmetry** rather than merely choosing one we designed.

## Authority boundary

UP-35 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35899713393` completed successfully at source head `99ae568ad891577f702793e258948570daa66ce8`.

The system started from six singleton spectral groups, capacity 6, with no finished partition menu. Using only the balanced training-side fit/validation split and the frozen objective, it accepted three pairwise merges and stopped because the next best merge failed the frozen minimum-gain rule.

The selected groups were `[0,1,2,5] + [3] + [4]`, with commutant capacity 18. Its inner-validation static accuracy was 0.9873046875 and mutable-commit accuracy was 0.9700520833333334.

After retraining on the complete original training pool, the induced structure reached 1.0 held-out accuracy, 1.0 mutable-commit accuracy, 1.0 exact-final-table accuracy, and 1.0 relational-query accuracy on the untouched true held-out pool. The capacity-36 control also scored 1.0 on the two primary metrics, leaving both retention deltas at zero.

Interpretation: this is a bounded positive for task-driven symmetry induction. Useful repeated-spectrum structure can be constructed incrementally from the minimum-symmetry state under a frozen local task-minus-capacity rule, without a finished partition catalog or true held-out guidance. The next experiment should remove the discrete merge operator and test continuous symmetry formation from independently parameterized spectral directions.
