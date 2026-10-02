# Wingless UP-39: rich-task direct spectral-offset optimization

Status: Windows-qualified scientific negative; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-38 qualified scientific negative sealed at `54dc8619149bbd0db8f6949f21453f52d5707ca4`.

## Question

UP-38 showed that direct continuous optimization of the six spectral offsets is technically viable, but mean phase cosine plus smooth symmetry cost is not a sufficient task objective. The optimizer improved phase alignment while remaining at commutant capacity 6 and performing poorly on mutable and relational behavior.

UP-39 asks:

**Can the same direct optimizer form useful exact symmetry when its smooth training-side objective explicitly contains mutable value confidence and correct relation probability in addition to phase alignment?**

## Controlled change from UP-38

UP-39 preserves the corrected UP-38 optimization mechanics:

- six independent anonymous spectral offsets;
- zero mean and RMS 0.05;
- no partition menu;
- no hard merge candidate evaluation;
- no pair-attraction field;
- deterministic full-rank tangent directions;
- projected central-difference directional estimator;
- perturbation 0.002;
- learning rate 0.002;
- maximum coordinate update 0.004;
- 12 optimization steps;
- smooth symmetry cost with tau 0.0025;
- resource price 0.02;
- sticky exact-fusion projection at 0.0025;
- same 96/32 balanced training-side split;
- true 128-table held-out pool inaccessible during optimization.

The only scientific change is the task objective.

## Rich smooth task signals

For each candidate geometry the observer is trained on the 96-table fit pool and evaluated on the 32-table inner-validation pool.

Three smooth training-side signals are measured:

1. **Normalized phase alignment**
   - mean phase cosine mapped from [-1,1] to [0,1].

2. **Teacher-forced mutable value probability**
   - 32 deterministic validation scenarios;
   - 8 writes per scenario;
   - before each write, the canonical true current table is transported through the candidate geometry;
   - the observer's probability assigned to each correct current value is accumulated;
   - the true write is then applied directly.
   - This is teacher-forced by design: optimization remains smooth and does not feed argmax errors back into the architecture search.

3. **Correct relation probability**
   - after the final teacher-forced state, the shared relation head receives the two entity distributions;
   - the probability assigned to the correct modulo-4 relation is accumulated.

No discrete accuracy is used in the optimization objective.

## Bottleneck objective

The three signals are combined with an equal-weight harmonic mean:

`H = 3 / (1/phase + 1/value_probability + 1/relation_probability)`

This intentionally makes the weakest capability dominate the pressure rather than allowing strong phase alignment to hide poor relational behavior.

The final optimization objective is:

`H - 0.02 * C_soft / 36`

where `C_soft` is the same smooth commutant-capacity approximation used by UP-38.

## Final evaluation

After 12 steps:

1. the best inner-validation checkpoint is frozen;
2. that geometry is retrained on the complete original training pool;
3. it is evaluated once on the untouched true held-out pool using the actual recurrent mutable loop with decoded states fed back into subsequent writes;
4. the capacity-36 control is evaluated under the same final protocol.

Teacher forcing exists only inside architecture optimization. It is not used for final scientific gates.

## Scientific gates

Preflight:

- balanced inner split;
- no true held-out selection.

Optimization:

- selected step >= 1;
- exact symmetry must emerge;
- selected capacity > 6 and < 36.

Final capability:

- held-out accuracy >= 0.985;
- mutable commit accuracy >= 0.95;
- exact final table accuracy >= 0.90;
- relational query accuracy >= 0.95.

Retention versus full capacity:

- held-out loss <= 0.02;
- mutable-commit loss <= 0.05.

## Interpretation boundary

A positive result would show that direct transport-parameter optimization can create useful repeated-spectrum symmetry when the objective reflects the actual mutable and relational capabilities the architecture is meant to preserve.

That would remove both the partition menu and the hand-designed pair-attraction flow from the successful path.

A negative result would distinguish objective insufficiency from optimizer insufficiency:

- if the rich smooth objective improves strongly but no exact symmetry forms, the continuous optimizer/projection mechanics remain the likely bottleneck;
- if the rich objective itself remains flat or misleading, teacher-forced confidence is not an adequate surrogate for recurrent mutable behavior.

UP-39 still retains important scaffolding:

- six predefined spectral parameters;
- fixed RMS normalization;
- fixed projected central-difference optimizer;
- sticky exact-fusion projection;
- phase-code supervision;
- teacher-forced smooth mutable surrogate.

A positive-result successor should remove the sticky exact-fusion projection and test whether learned near-degeneracy alone is sufficient or converges to exact equality without snapping.

## Plain speak

UP-38 asked the system to reorganize itself mostly by making its phase representation look better.

It did that, but it did not learn the structure needed for robust memory.

UP-39 tells it what we actually care about in a smooth way:

- "How confident are you about the right values after transport?"
- "How confident are you about the right relationship between them?"
- "Are your phase codes still clean?"

The weakest of those pressures matters most.

If this produces the symmetry that UP-38 failed to discover, we will know the missing ingredient was not manual architecture search. It was **the system having the right internal self-evaluation signal**.

## Authority boundary

UP-39 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

The authoritative Windows qualification on workflow run `35930763303` completed successfully at source head `62ad70fc579fe8fc780c5ca2f05f838a5868ac2c` on runner `WINGLESS-LINKDEADKB`.

The harness and full repository regression passed. The rich smooth architecture objective improved deterministically from `0.6867884391276392` at initialization to `0.6972457303665964` at selected step 7. Its harmonic task score improved from `0.6901217724609725` to `0.7005790636999297`. The main contributing signal was correct-relation probability, which rose from `0.5260644564224123` to `0.5489698582362194`; normalized phase score and correct-value probability moved slightly downward.

Despite that smooth-objective improvement, no exact fusion occurred and selected commutant capacity remained 6. On the untouched true held-out pool, the selected checkpoint reached:

- static held-out accuracy: `0.901123046875`;
- mutable commit accuracy: `0.34375`;
- exact final table accuracy: `0.3541666666666667`;
- relational query accuracy: `0.5416666666666666`;
- held-out retention loss versus capacity 36: `0.098876953125`;
- mutable-commit retention loss versus capacity 36: `0.65625`.

The full-capacity control remained perfect.

Scientific classification: **qualified negative**. Because the smooth objective improved while the actual free-running recurrent loop remained poor, the frozen interpretation fork points to **teacher-forcing exposure mismatch** rather than permission to retune thresholds or gates.

The successor experiment should preserve the rich objective and corrected direct optimizer while replacing the teacher-forced mutable surrogate with a smooth free-running probabilistic memory rollout. The soft recurrent state should be reconstructed directly from the four decoder distributions in the existing 16-dimensional one-hot-per-entity memory representation, apply the commanded write exactly to the written entity, renormalize, and transport again. Final scientific evaluation remains the unchanged hard recurrent loop on the untouched held-out pool.
