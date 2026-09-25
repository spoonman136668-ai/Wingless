# Wingless UP-81A — six-role coverage transfer

Status: preregistered scientific state-space scaling experiment.

Scientific parent: sealed UP-80A Windows evidence `83ed6b29dc439653505e2b0dbd9239f2fa0373fe`.

## Question

UP-77A through UP-80A established that, in the frozen five-role ternary task, one-step compact factorized readout requires more than 27 states and is sufficient at 36 states across tertiary slices and multiple modular selector orientations.

UP-81A asks whether that structural-coverage phenomenon scales to a larger role topology and joint state space without changing the optimizer discipline.

## Frozen design

- six roles, three values per role;
- full joint state space: 729 states;
- primary residue split is sum(roles) mod 3;
- held-out set: primary residue != 0, exactly 486 states;
- training pool: primary residue == 0;
- frozen secondary selector coefficients [1,2,1,2,1,2];
- frozen tertiary selector coefficients [1,1,2,2,1,0];
- training states: secondary == 0 plus secondary == 1 / tertiary == 0;
- expected training count: exactly 108 states, preserving the same 4/9 fraction of the primary-residue-zero pool as the five-role 36-of-81 condition;
- raw factorized one-hot representation, 18 features;
- role-factorized simplex representation, 12 features;
- six independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- exactly one optimization step;
- per-role held-out gate 0.98.

No adaptive stopping, optimizer change, learning-rate change, threshold change, seed search, selector change, representation change, promotion, or activation is permitted.

## Interpretation

Each representation is sealed exactly as observed. A pass would show that the one-step structural-coverage result survives a 3x increase in joint state count and an added role while preserving the same relative structural coverage. A failure would identify a role-count/state-space scaling boundary and must not trigger post-result tuning.
