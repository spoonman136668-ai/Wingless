# Wingless UP-61A — frozen training-budget convergence

Status: preregistered scientific cognition experiment.

Scientific parent: UP-60A seal `c6d8494bbbb2ddc1bf2895c4405790558de78051`.

UP-60A quantified the double-control angle precision required for long composition. UP-61A asks whether the existing learner naturally converges into that precision band when given more optimization steps.

Five independent runs begin from the same [0.1,0.1] initialization with unchanged learning rate 0.10 and finite-difference epsilon 1e-6. Only the training-step budget varies: 90, 180, 360, 720, 1440. Each learned parameter pair is evaluated on the unchanged length-128 mixed-program task and numerical gates.

No learning-rate, epsilon, topology, loss, initialization, or acceptance threshold changes are permitted.


## Authoritative Windows result

Workflow run: `36034883586`

Runner: `WINGLESS-LINKDEADKB`

Source head: `f2b913e4e7917769daf3773dbe8e1f6550705524`

Artifact: `10824621473`

Artifact digest: `sha256:8ed46a7e6f99ea1ec5e990c0c783397289553f1e21ca0679c924c5fb76993f9f`

Length-128 results by frozen training budget:

- 90 steps: accuracy `0.691358024691358`, double-angle error `-0.19753719917986023`;
- 180: `0.691358024691358`, error `-0.13151978940933357`;
- 360: `0.7407407407407407`, error `-0.07676879712987916`;
- 720: `0.9135802469135802`, error `-0.0333669835866508`;
- 1440: accuracy `1.0`, error `-0.006147038516198578`, PASS.

All other optimizer settings were unchanged.

## Scientific classification

The prior long-horizon failure is not an optimizer plateau. The existing learner converges into the required angle-precision band when given sufficient frozen optimization budget.

The next A-lane experiment tests whether that convergence is robust to preregistered initial parameter values rather than being specific to the original [0.1,0.1] initialization.
