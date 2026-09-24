# Wingless UP-61A — frozen training-budget convergence

Status: preregistered scientific cognition experiment.

Scientific parent: UP-60A seal `c6d8494bbbb2ddc1bf2895c4405790558de78051`.

UP-60A quantified the double-control angle precision required for long composition. UP-61A asks whether the existing learner naturally converges into that precision band when given more optimization steps.

Five independent runs begin from the same [0.1,0.1] initialization with unchanged learning rate 0.10 and finite-difference epsilon 1e-6. Only the training-step budget varies: 90, 180, 360, 720, 1440. Each learned parameter pair is evaluated on the unchanged length-128 mixed-program task and numerical gates.

No learning-rate, epsilon, topology, loss, initialization, or acceptance threshold changes are permitted.
