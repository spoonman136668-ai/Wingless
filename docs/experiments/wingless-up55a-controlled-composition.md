# Wingless UP-55A — controlled/context-dependent composition

Status: preregistered scientific cognition experiment.

Scientific parent: UP-54A seal `bddeaf6fc74cd4bfc509b09bdb948b7600a41088`.

## Question

UP-54A showed perfect structural transfer for unconditional reversible primitives across four roles. UP-55A increases computational structure by introducing a context-dependent controlled operation: mutate one role only when another role carries a specified value.

Can a shared unitary rule learned from one target/control placement and one control value generalize to unseen control values, moved control/target roles, and long mixed programs?

## Frozen design

- 3 roles, 3 values, joint dimension 27;
- two learned scalar parameters: adjacent-role swap and controlled value swap;
- training is single-step only;
- controlled mutation is trained only for target role 0, control role 1, control value 0;
- role swap is trained only at position 0/1;
- held-out evaluation includes control values 1 and 2, control moved to role 2, target moved to role 1, and mixed programs of lengths 16 and 48;
- matched non-unitary control receives identical supervision and optimization;
- unchanged cognition gates: >=0.99 train, >=0.98 held/category accuracy, <=1e-9 norm drift, <=1e-8 round-trip error.

No result-informed change to topology, optimizer, thresholds, or held-out categories is permitted.
