# Wingless UP-40: free-running probabilistic direct objective

Status: preregistered research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-39 qualified scientific negative sealed at `8994882409133d2da22bf7f6d482aa09e6dd0a9f`.

## Question

UP-39 replaced phase-only self-evaluation with an equal-weight harmonic objective containing phase alignment, correct mutable-value probability, and correct relation probability. Its smooth objective improved deterministically, but the selected architecture remained at capacity 6 with no exact fusion and the actual hard recurrent mutable loop remained poor.

UP-40 asks:

**Can the same rich direct optimizer form useful exact symmetry when its mutable training-side surrogate must live with its own uncertainty across writes instead of being reset to the true table after every write?**

## Controlled change from UP-39

Everything below remains frozen:

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
- smooth commutant-capacity tau 0.0025;
- resource price 0.02;
- sticky exact-fusion projection at 0.0025;
- balanced 96/32 training-side split;
- untouched 128-table true held-out pool;
- normalized phase score;
- correct mutable-value probability;
- correct relation probability;
- equal-weight harmonic bottleneck objective;
- all final scientific gates.

The only scientific change is the mutable smooth surrogate.

## Smooth free-running probabilistic memory

The existing memory representation is 16-dimensional: four entities times four values.

UP-40 starts each training-side scenario from the canonical 16-D state. After each transport:

1. decode the four value probability distributions;
2. score probability assigned to the correct current values;
3. construct the next 16-D memory directly from those probabilities, using the same four entity blocks;
4. apply the commanded write exactly to the written entity;
5. renormalize the 16-D state;
6. feed that soft state into the next transport.

No argmax is fed back. The true table is retained only as the training-side scoring target and is updated by the commanded write; it is never used to reset the architecture state.

The same deterministic perturbation/noise and global-phase schedule used by UP-39 is preserved.

## Objective

The three smooth signals remain:

1. normalized phase score;
2. mean probability assigned to correct mutable values;
3. mean probability assigned to the correct relation.

They remain combined as:

`H = 3 / (1/phase + 1/value_probability + 1/relation_probability)`

and:

`objective = H - 0.02 * C_soft / 36`

with the same `C_soft` and tau 0.0025 used by UP-38/UP-39.

## Final evaluation

The soft probabilistic recurrence exists only inside training-side architecture optimization.

After model selection, final evaluation remains the unchanged actual recurrent mutable loop:

- decode a hard table;
- feed that decoded table back after the commanded write;
- evaluate once on the untouched true held-out pool.

Thus a positive result must transfer from the smooth exposure-aware surrogate to the real discrete recurrence.

## Frozen scientific gates

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

Retention versus capacity 36:

- held-out loss <= 0.02;
- mutable-commit loss <= 0.05.

These gates are frozen before the result.

## Interpretation fork

- If UP-40 passes, the next scaffold to remove is sticky tolerance-based snapping. Test a proximal / convex-clustering style fusion operator so equality arises from the optimization penalty itself.
- If the free-running smooth objective improves and forms useful symmetry but hard recurrence still misses final gates, isolate the remaining smooth-to-discrete decoding mismatch without changing the frozen final gates.
- If the free-running objective improves but still cannot move the offsets into symmetry, direct optimizer geometry becomes the main suspect; isolate optimizer adequacy without returning to manual architecture selection.
- If the free-running objective itself degrades or becomes unstable, the soft-state construction is not an adequate surrogate; seal the negative exactly as observed rather than tuning it after the result.

## Plain speak

UP-39 graded the system after every write, then quietly handed it the correct memory again before the next write.

UP-40 removes that crutch.

The smooth training loop now has to carry its own uncertainty forward. A confident mistake or ambiguous memory can affect later writes, which is much closer to the real recurrent test while remaining continuous enough to optimize.

If that pressure causes the spectral offsets to organize into useful symmetry, exposure mismatch was the missing piece. If not, the failure tells us where to look next without relaxing any gate.

## Authority boundary

UP-40 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
