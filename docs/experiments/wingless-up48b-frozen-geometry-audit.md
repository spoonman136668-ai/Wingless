# Wingless UP-48B — frozen geometry causal audit

Status: Windows-qualified scientific result.

Scientific parent: Windows-qualified UP-47 source head `474ebb42082d1a9dd022517edb3c48c4f16dbd1d`.

## Question

UP-47 found that step 42 maximizes the frozen smooth objective, while exact-fused step 29 has much stronger hard recurrent capability than step 42.

UP-48B asks whether exact equality itself explains the step-29 advantage, or whether another geometric feature is responsible.

## Frozen arms

1. UP-45 step 24.
2. UP-45 exact-fused step 29.
3. Step 29 with only fused pair [0,3] split symmetrically by the exact pair gap observed at frozen step 28; pair mean and all other coordinates remain unchanged.
4. UP-47 dense-selected step 42.
5. Step 42 with only its nearest frozen pair [4,5] set exactly equal at their mean.

No optimizer is run. No lambda, seed, task objective, held-out pool, or capability threshold is changed.

## Interpretation

- If step 29 materially outperforms its matched split control on hard recurrent capability, exact equality itself has causal support at the step-29 geometry.
- If the matched split retains the step-29 advantage, the useful feature is broader geometry rather than equality alone.
- If fusing step 42's nearest pair materially improves its hard capability, exact equality may repair part of the dense-selected geometry.
- If it does not, the step-29 advantage depends on other coordinates/interactions and the next mechanism experiment should isolate those features rather than increase fusion pressure.

Scientific negatives are valid.

## Plain speak

We have one state with the best smooth score and another state that is much better at the real memory task.

This experiment changes only the equality relationship in tightly matched controls. It asks whether making two channels exactly equal is actually what helped, or whether step 29 was good for some other reason.


## Authoritative Windows result

Workflow run: `35986102843`

Runner: `WINGLESS-UP-B`

Artifact: `10802307659`

Artifact digest: `sha256:846fae656e0caf3e7aa17fadfe651d6320e05e4cb6ef7d215e75035767ed0284`

The harness, deterministic double probe, focused tests, full repository regression, and host-priority guard passed.

Exact-fused step 29 minus its matched split control:

- objective: `+0.030081334046351493`
- held-out accuracy: `+0.042724609375`
- commit accuracy: `+0.2799479166666667`
- final-table accuracy: `+0.25`
- relation accuracy: `+0.1875`

Dense-selected step 42 minus step 29:

- objective: `+0.03270116784855848`
- held-out accuracy: `-0.0380859375`
- commit accuracy: `-0.2109375`
- final-table accuracy: `-0.14583333333333331`
- relation accuracy: `-0.16666666666666663`

Fusing only step 42's nearest pair [4,5] minus unchanged step 42:

- objective: `+0.04674098383119574`
- held-out accuracy: `+0.07373046875`
- commit accuracy: `+0.44140625`
- final-table accuracy: `+0.3333333333333333`
- relation accuracy: `+0.33333333333333326`

## Scientific classification

Exact equality itself has causal support at the frozen step-29 geometry: breaking only the fused pair materially worsened both the smooth score and all hard recurrent measures.

The UP-47 smooth objective remains misaligned with hard capability across whole geometries, because step 42 scores higher than step 29 while performing substantially worse on the hard recurrent task.

However, adding exact equality to step 42's nearest pair improved both the smooth objective and hard capability sharply. The next discriminating mechanism experiment is therefore an exhaustive frozen single-pair fusion screen over all 15 step-42 pairs, selected by training-side smooth objective only.

## Plain speak

The equality was not a coincidence. When we broke only that equality, performance got worse.

Even more importantly, forcing one clean equality into the otherwise best-scoring step-42 state made its real memory performance jump sharply. The next job is to test every possible pair fairly instead of assuming the nearest pair is the best one.
