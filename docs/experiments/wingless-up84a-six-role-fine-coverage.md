# Wingless UP-84A — six-role one-step fine coverage boundary

Status: preregistered scientific coverage-boundary experiment.

Scientific parent: sealed UP-83A Windows evidence `9167f6909c9102a82fe81753e4c337cc22cebe1f`.

## Question

UP-83A passed both compact factorized representations at 27 training states with one optimization step on the frozen six-role structural split. UP-84A localizes the remaining data-coverage boundary inside that 27-state cell using a fourth structural coordinate frozen before execution.

## Frozen design

Unchanged:

- six ternary roles and 729 joint states;
- primary residue, secondary coordinate, and tertiary coordinate from UP-81A/UP-83A;
- 486 held-out states with primary residue != 0;
- raw factorized one-hot representation, 18 features;
- role-factorized simplex representation, 12 features;
- six independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- exactly one optimization step;
- per-role held-out gate 0.98.

The preregistered fourth coordinate is `quaternary = role0 mod 3`. Within the existing 27-state cell primary==0, secondary==0, tertiary==0, the training levels are exactly:

- 9 states: quaternary==0;
- 18 states: quaternary<2;
- 27 states: all states in the existing cell, retained as the sealed UP-83A control.

No adaptive stopping, optimizer change, learning-rate change, selector change, threshold change, seed search, extra states, representation change, promotion, or activation is permitted.

## Interpretation

All six level/representation points are sealed exactly as observed. A 9-state pass bounds this selector family at <=9; a 9-state failure and 18-state pass localizes the transition to 9–18; an 18-state failure with the 27-state control passing localizes it to 18–27. Any other pattern is retained without tuning.
