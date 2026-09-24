# Wingless UP-72A — compact factorized 54-state split replication

Status: preregistered scientific coverage/robustness experiment.

Scientific parent: UP-71A seal `b9738d3d63e591d9d6194fdba43315344dc9dc88`.

## Question

UP-71A found that both compact factorized representations failed the unchanged held-out gate at 27 or fewer training states and passed perfectly at 54 and 81 states on residue-0 training.

UP-72A asks whether the 54-state passing point is robust to rotating the primary train/held-out residue rather than being specific to the original residue-0 partition.

## Frozen design

The logical domain, learner, optimizer, and gate are unchanged:

- 243 five-role ternary states;
- five independent 3-class linear-softmax role heads;
- zero initialization;
- 800 optimization steps at learning rate 1.0;
- unchanged per-role held-out gate of 0.98.

Three primary partitions are fixed before execution:

- sum(role values) mod 3 == 0;
- sum(role values) mod 3 == 1;
- sum(role values) mod 3 == 2.

Within each 81-state training residue, the unchanged UP-71A secondary coordinate is used and exactly 54 states are retained by requiring secondary < 2. The held-out set is always the other 162 states.

Two independently trained compact factorized representations are tested at each residue:

- raw factorized one-hot, 15 features;
- role-factorized simplex, 10 features.

No adaptive subset selection, threshold change, extra training state, nonlinear observer, promotion, production authority, or activation is permitted.

## Interpretation

This is a split-replication experiment at the observed 54-state passing point. Failure on any residue is valid evidence that the UP-71A sample-efficiency result is partition-sensitive. Passing across all three residues supports split robustness but does not authorize mechanism promotion.
