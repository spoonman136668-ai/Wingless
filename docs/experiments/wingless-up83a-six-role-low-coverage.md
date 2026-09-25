# Wingless UP-83A — six-role one-step low-coverage challenge

Status: preregistered scientific coverage-boundary experiment.

Scientific parent: sealed UP-82A Windows evidence `fa2f6e41aa2871643f23556bde2ac4dfc875e5a0`.

## Question

UP-82A passed both compact factorized representations at 54, 81, and 108 training states with one optimization step over the frozen six-role structural split. UP-83A asks whether the six-role boundary extends to a single 27-state secondary/tertiary cell.

## Frozen design

Unchanged: six ternary roles; 729 joint states; primary residue sum mod 3; 486 held-out states with primary residue != 0; secondary coefficients [1,2,1,2,1,2]; tertiary coefficients [1,1,2,2,1,0]; raw factorized one-hot (18 features); role-factorized simplex (12 features); six independent 3-class linear-softmax heads; zero initialization; learning rate 1.0; exactly one optimization step; per-role held-out gate 0.98.

Training levels are exactly:
- 27: primary==0, secondary==0, tertiary==0;
- 54: primary==0, secondary==0, tertiary<2, retained as the sealed UP-82A boundary control.

No adaptive stopping, optimizer change, learning-rate change, selector change, threshold change, seed search, extra states, representation change, promotion, or activation is permitted.

## Interpretation

All four level/representation points are sealed exactly as observed. A 27-state pass moves the boundary to at most 27 under this selector; a failure with the frozen 54-state control passing localizes the transition to 27–54. Any other pattern is retained without tuning.
