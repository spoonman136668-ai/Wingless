# Wingless UP-54A — four-role structural transfer

Status: preregistered scientific cognition experiment.

Scientific parent: UP-53A seal `b02660bf643d98f70468ed6296e3547be4ea422c`.

## Question

UP-53A showed that three informative examples are enough to identify and structurally transfer the learned three-role operators. The next end-goal question is whether the same shared reversible rules continue to compose when the bound state grows from three roles to four.

## Frozen design

- four ordered roles;
- three values per role;
- joint basis dimension 81;
- two learned scalar parameters exactly as in the three-role family: role-0 value swap and adjacent-role swap;
- training uses single-step supervision only;
- value swaps are trained only on role 0;
- adjacent-role swap is trained only at position 0/1;
- held-out tests include never-trained swap positions 1/2 and 2/3, derived mutations of roles 1, 2, and 3, and mixed programs of lengths 16 and 48;
- matched non-unitary control gets the same training set, optimizer, and budget;
- the existing cognition thresholds remain unchanged: >=0.99 train, >=0.98 held/category accuracy, <=1e-9 norm drift, <=1e-8 round-trip error.

Scientific negatives are valid. No threshold or optimizer change is permitted after results are observed.
