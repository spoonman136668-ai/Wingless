# Wingless UP-48B — frozen geometry causal audit

Status: preregistered scientific diagnostic.

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
