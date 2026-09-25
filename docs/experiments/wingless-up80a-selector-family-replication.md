# Wingless UP-80A — selector-family replication

Status: preregistered scientific replication experiment.

Scientific parent: sealed UP-79A Windows evidence `145f5730c4e86bced6f31e9a6bdb29eb25b03c56`.

## Question

UP-79A showed that all three tertiary slices reach perfect held-out accuracy with 36 training states, ruling out dependence on the original added nine-state slice. UP-80A asks whether the observed 36-state sufficiency is robust to changing the modular selector orientation itself rather than only changing the tertiary slice inside one selector family.

## Frozen design

Unchanged:

- the same 243 five-role ternary states;
- the same primary residue split;
- the same 162 held-out states with primary residue != 0;
- raw factorized one-hot representation, 15 features;
- role-factorized simplex representation, 10 features;
- five independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- exactly one optimization step;
- per-role held-out gate 0.98.

Three selector families are frozen before execution:

- family A secondary [1,2,1,2,1], tertiary [1,1,2,2,0];
- family B secondary [2,1,2,1,2], tertiary [1,2,2,1,0];
- family C secondary [1,1,2,1,2], tertiary [2,1,1,2,0].

For each family, training includes all primary-residue-zero states with secondary value 0 plus secondary value 1 / tertiary value 0. Each family has exactly 36 training states: a 27-state structural base plus one nine-state increment.

No adaptive stopping, optimizer change, learning-rate change, threshold change, seed search, post-result selector choice, representation change, promotion, or activation is permitted.

## Interpretation

Every family/representation cell is sealed exactly as observed. If all three families pass, the evidence for a generic 36-state structural-coverage threshold strengthens materially. If one or more fail, the threshold is selector-family dependent and the previously observed sufficiency must be interpreted as partition-specific. Harness acceptance is independent of scientific outcome.
