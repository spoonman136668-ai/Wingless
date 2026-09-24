# Wingless UP-44: square-root probability soft-memory bridge

Status: preregistered research branch only; not activated, promoted, or connected to ckb-plane.

Parent result: UP-43 qualified scientific negative sealed at `85f9ab7837bacda5deed635c5e7bf1e5ed5b4c05`.

## Why this experiment exists

UP-43 repaired the optimizer and demonstrated that the real task gradient can enter reversible symmetry neighborhoods. It reached provisional capacity 8, including a nearest offset gap of approximately `0.0007048` at step 24. However, every preregistered smooth-objective checkpoint scored below the initial geometry, so selection remained at step 0.

That means optimizer reachability is no longer the cleanest explanation.

The UP-40/UP-43 soft recurrent state encodes each decoder probability as amplitude `0.5*p` and then globally normalizes the 16-D vector. Because the L2 norm of a probability vector depends on its concentration, that construction changes each entity block's effective weight as its uncertainty changes. Global normalization therefore couples confidence in one entity to the amplitudes of all four entities.

UP-44 tests only that representation choice.

## Question

**Does replacing probability amplitudes `p` with `sqrt(p)` make the free-running task objective align better with the symmetry neighborhoods that the repaired optimizer can already reach?**

## Controlled change

All UP-43 mechanics remain frozen:

- six independent spectral offsets;
- zero mean and RMS 0.05;
- seven-observation rolling full-rank directional memory;
- one plus/minus probe pair per update;
- 48 updates;
- checkpoint evaluations at 0, 12, 24, and 48;
- perturbation 0.002;
- learning rate 0.002;
- maximum coordinate update 0.004;
- no sticky snapping;
- no known symmetry target;
- no partition menu or hard merge search;
- same 96/32 training-only split;
- untouched 128-table true held-out pool;
- same phase/value/relation harmonic objective;
- resource price 0.02;
- soft-capacity tau 0.0025;
- same final hard recurrent evaluation and capability gates.

The only scientific change is the recurrent soft-memory amplitude encoding.

## Square-root probability representation

For each unwritten entity/value coordinate:

`amplitude = 0.5 * sqrt(p)`

Because each entity distribution sums to one:

`sum_v amplitude_v^2 = 0.25 * sum_v p_v = 0.25`.

Four entity blocks therefore have total squared norm exactly one before noise:

`4 * 0.25 = 1`.

The commanded-write entity is still replaced by an exact one-hot block with amplitude 0.5.

No global soft-memory normalization is applied after constructing the state.

This preserves the decoder probability distribution in squared amplitudes and prevents confidence in one entity from rescaling all other entity blocks.

## Selection and held-out discipline

The smooth objective selects only among checkpoints 0, 12, 24, and 48 using the same 96/32 training-side split.

The true held-out pool is touched only after selection is frozen.

Provisional symmetry measurements never select a checkpoint.

## Frozen interpretation

- If the square-root representation produces a positive selected objective gain and the selected checkpoint lies in a provisional fusion basin, the representation mismatch was a material blocker. Proceed next to a principled proximal exact-fusion operator; do not restore sticky snapping.
- If the objective improves but the selected checkpoint remains outside a fusion basin, the representation helps task alignment but symmetry pressure remains insufficient. Isolate the resource/symmetry term without changing final gates.
- If the trajectory enters a basin but selection still prefers step 0 or another non-basin checkpoint, the task objective still rejects emergent symmetry; next isolate the objective composition rather than optimizer reachability.
- If the objective does not improve at all, the square-root representation is not sufficient. Seal the negative exactly and examine whether the decoder-probability surrogate itself is the wrong state variable.
- If capability gates excluding exact equality become green, report that separately but do not claim exact symmetry induction until equality is produced by a principled operator or emerges spontaneously.

No post-result threshold or gate change is allowed.

## Plain speak

UP-43 showed the steering system can reach symmetry, but its internal score says those places are worse.

One possible reason is that our soft memory representation was distorting uncertainty. A blurry prediction for one entity could change the scale of every other entity after global normalization.

UP-44 removes that distortion. Each entity always owns exactly one quarter of the state's probability mass, while uncertainty only changes how that quarter is distributed inside the entity.

If the task score now starts preferring the symmetry neighborhoods, we will have isolated a real representation bug rather than forcing the architecture toward a desired answer.

## Authority boundary

UP-44 is isolated mathematical research. It does not invoke a language model, activate Wingless, alter CKB/ckb-plane/KTRADE/Nemotron authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
