# Wingless UP-203B — order-signal persistence

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-202B fc4846afbf855323c2c90e160c2bccf545a1661a.

## Question

Does each history family retain a stable relative permutation signature on later phases after removing phase-wide calibration shifts?

## Frozen compositions

- mixed4: 0,5,1,6
- observe4: 5,6,7,8

All 24 permutations per composition.

## Frozen phases

Training/template phases:
- 55
- 56
- 57

Untouched diagnostic phases:
- 76
- 77
- 78

## Frozen diagnostic

For each composition:
1. Compute the training-template residual for each exact permutation:
   mean factor for that permutation across template phases minus the composition-wide template mean.
2. For each untouched phase:
   center the 24 permutation factors by that phase's composition mean.
3. Measure Pearson correlation between the frozen training residual template and the untouched centered residual vector.
4. Measure within-phase factor spread across the 24 permutations.

No prediction is applied and no held-out parameter is fit.

## Interpretation

Stable positive residual correlation supports persistent order geometry even if absolute calibration moves. Near-zero correlation with nonzero spread indicates order effects exist but reorganize across phases. Near-zero spread indicates little order dependence in that history family.

## Bounds

Diagnostic only. No held-out fitting, feature search, correction activation, maintenance action, capacity change, live activation, or production authority.
