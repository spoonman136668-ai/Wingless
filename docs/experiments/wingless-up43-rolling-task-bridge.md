# Wingless UP-43: rolling-gradient free-running task bridge

Status: Windows-qualified scientific negative; not activated, promoted, or connected to ckb-plane.

Parent calibration: UP-42 sealed at `a76463ccff664fdb392ead5d024fb1cb8e7778e6`.

## Question

UP-42 showed that seven-observation rolling directional memory reaches an independently known-good symmetry basin on the same minimum tested horizon as an exact-gradient reference, while the old sequential estimator does not enter the basin even by 96 updates.

UP-43 returns that optimizer repair to the real task objective:

**Does the UP-40 free-running rich task objective, with no known symmetry target injected, drive independently parameterized spectral offsets into a reversible near-symmetry basin when optimized with rolling full-rank directional memory over the calibrated 48-step horizon?**

## Controlled inheritance from UP-40

UP-43 keeps:

- six independent spectral offsets;
- zero mean;
- RMS 0.05;
- no partition menu;
- no hard merge search;
- no pair-attraction field;
- perturbation 0.002;
- learning rate 0.002;
- maximum coordinate update 0.004;
- smooth capacity tau 0.0025;
- resource price 0.02;
- balanced training-only 96/32 split;
- untouched true held-out pool of 128;
- free-running probabilistic mutable surrogate;
- normalized phase, correct-value probability, and correct-relation probability;
- equal-weight harmonic bottleneck objective;
- unchanged final hard recurrent evaluation.

No UP-35 target offsets, groups, labels, or geometry enter the task optimizer.

## Calibrated optimizer repair

Three optimizer mechanics change based on UP-41/UP-42 calibration evidence:

1. the one-direction update is replaced by rolling full-rank directional memory;
2. maximum horizon becomes 48 updates;
3. irreversible sticky projection is disabled.

The rolling estimator still uses one plus/minus directional probe pair per update. It stores at most seven recent directional derivatives and reconstructs the five-dimensional zero-mean tangent gradient once full rank is available.

## Reversible near-symmetry measurement

Because UP-41 demonstrated that irreversible threshold snapping can fuse the wrong structure, UP-43 does **not** mutate offsets when they become close.

After every update it only measures a provisional clustering at the existing 0.0025 distance:

- the measurement is reversible;
- it cannot alter the next optimizer state;
- it is not used to choose a known target;
- it records provisional capacity and nearest pairwise gap.

Thus UP-43 can answer whether the task gradient naturally drives coordinates into a fusion basin without yet deciding how exact equality should be imposed.

## Frozen training-only checkpoint selection

Exact smooth-objective evaluations are performed at:

`0, 12, 24, 48`

Step 0 is the initial state. The best objective checkpoint is selected using only the 96/32 training-side split.

The true held-out pool is touched only after that checkpoint is frozen.

## Final evaluation

The selected raw, unfused offsets are evaluated with the unchanged hard recurrent mutable loop on the untouched true held-out pool, alongside the capacity-36 control.

The original capability gates are reported unchanged:

- held-out >= 0.985;
- mutable commit >= 0.95;
- exact final table >= 0.90;
- relation >= 0.95;
- held-out retention loss <= 0.02;
- commit retention loss <= 0.05.

Because exact fusion is deliberately disabled, UP-43 is a **bridge experiment**, not a final symmetry-induction pass. It separately reports whether all capability/retention gates other than exact equality are met.

## Frozen interpretation

- If the task gradient enters a reversible fusion basin and the selected checkpoint materially improves task behavior, proceed to a principled proximal/convex exact-fusion operator; do not restore sticky snapping.
- If the task gradient enters a basin but capability remains poor, the geometry being approached is not yet task-useful; inspect the training-side objective/representation before exact fusion.
- If the smooth objective improves but the task gradient never enters a fusion basin, the corrected optimizer is no longer the main blocker. The next suspect is the free-running soft-state representation/objective, with the already identified norm-preserving `sqrt(p)` encoding as the clean representation test.
- If the objective itself fails to improve, seal the negative and test the smooth state representation rather than increasing horizon or loosening gates.

No threshold or task gate is changed after the result.

## Plain speak

We fixed the steering system in calibration. Now we put it back on the real road without telling it where UP-35 went.

It gets 48 moves, remembers the recent directional evidence, and is not allowed to glue anything together.

We only watch whether the task itself pulls any spectral directions close enough that a future principled fusion rule would have something real to fuse. Then we test the chosen checkpoint on the untouched hard recurrent task.

That separates “can the task discover near-symmetry?” from “how should exact equality be enforced?”

## Authority boundary

UP-43 is isolated Wingless research. It does not invoke a language model, activate Wingless, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

The authoritative Windows qualification on workflow run `35935011865` completed successfully at source head `47fac701b3e73e19fa8826a959776ce65e8c01ea` on runner `WINGLESS-LINKDEADKB`. The production-priority guard, focused tests, deterministic double probe, and full repository regression passed.

The repaired rolling optimizer did enter reversible symmetry neighborhoods:

- first provisional fusion step: `10`;
- maximum provisional capacity: `8`;
- step-10 nearest gap: `0.0011419367693766488`;
- step-24 provisional capacity: `8`;
- step-24 nearest gap: `0.0007047869659162478`;
- another capacity-8 visit occurred later with nearest gap `0.00010190873120542211`.

However, the training-side objective rejected every preregistered checkpoint:

- initial objective: `0.6498734172582427`;
- step 12 objective: `0.6257176169903816`;
- step 24 objective: `0.6164862228413023`;
- step 48 objective: `0.6222920471220585`.

Therefore the selected checkpoint remained step 0, with objective gain exactly `0`, provisional capacity 6, and nearest gap `0.029277002188455997`.

On the untouched true held-out pool, the selected step-0 checkpoint reached:

- held-out accuracy: `0.90185546875`;
- mutable commit accuracy: `0.34375`;
- exact final table accuracy: `0.375`;
- relation accuracy: `0.5625`;
- held-out retention loss versus capacity 36: `0.09814453125`;
- mutable-commit retention loss versus capacity 36: `0.65625`.

Capability gates excluding exact fusion did not pass.

Scientific classification: **qualified negative with optimizer-transfer evidence**.

Interpretation: UP-42 repaired reachability sufficiently for the task gradient to enter symmetry neighborhoods, but the UP-40 free-running soft-state objective did not prefer those states. Under the frozen UP-43 interpretation, the next controlled test changes the soft-state representation rather than the horizon, gates, optimizer, or fusion rule. The successor should replace amplitude `p` followed by global renormalization with a norm-faithful `sqrt(p)` block encoding while preserving the rolling optimizer, 48-step horizon, no sticky snapping, checkpoint schedule, task objective weights, and held-out discipline.
