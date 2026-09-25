# Wingless UP-82A — six-role one-step coverage boundary

Status: preregistered scientific coverage-boundary experiment.

Scientific parent: sealed UP-81A Windows evidence `8acef3abd83d0a40fe8ecbd2dbb5627bb3bace91`.

## Question

UP-81A transferred perfect one-step compact factorized held-out generalization to six roles at 108 training states out of the 243-state primary-residue-zero pool. UP-82A localizes the six-role structural-coverage boundary without changing optimizer discipline.

## Frozen design

Unchanged: six ternary roles; 729 joint states; primary residue sum mod 3; 486 held-out states with primary residue != 0; secondary coefficients [1,2,1,2,1,2]; tertiary coefficients [1,1,2,2,1,0]; raw factorized one-hot (18 features); role-factorized simplex (12 features); six independent 3-class linear-softmax heads; zero initialization; learning rate 1.0; exactly one optimization step; per-role held-out gate 0.98.

Training levels are exactly:
- 54: primary==0, secondary==0, tertiary<2;
- 81: primary==0, secondary==0;
- 108: primary==0 and (secondary==0 or (secondary==1 and tertiary==0)).

No adaptive stopping, optimizer change, learning-rate change, selector change, threshold change, seed search, extra states, representation change, promotion, or activation is permitted.

## Interpretation

All six level/representation points are sealed exactly as observed. The result localizes whether the six-role one-step boundary lies at or below 54, between 54 and 81, between 81 and 108, or above 108 under the frozen selectors.
