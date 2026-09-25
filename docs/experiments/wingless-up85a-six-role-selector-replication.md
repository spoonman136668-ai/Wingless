# Wingless UP-85A — six-role fine-coverage selector replication

Status: preregistered scientific replication experiment.

Scientific parent: sealed UP-84A Windows evidence `fb50c21756f2a1f6a1cb472fd382c74760cb52e0`.

## Question

UP-84A localized the six-role one-step coverage transition to 18–27 training states under quaternary coordinate role0 mod 3: both representations scored 0.944444 at 18 and 1.0 at 27. UP-85A asks whether that transition replicates under a distinct fourth coordinate rather than being specific to the role0 slice.

## Frozen design

Unchanged: six ternary roles; 729 joint states; existing primary/secondary/tertiary coordinates; 486 held-out states with primary residue != 0; raw factorized one-hot (18 features); role-factorized simplex (12 features); six independent 3-class linear-softmax heads; zero initialization; learning rate 1.0; exactly one optimization step; per-role held-out gate 0.98.

The preregistered alternate fourth coordinate is `quaternary_b = role1 mod 3`. Within primary==0, secondary==0, tertiary==0:

- 18 states: quaternary_b<2;
- 27 states: all states in the cell, retained as the sealed control.

No adaptive stopping, optimizer change, learning-rate change, selector change, threshold change, seed search, extra states, representation change, promotion, or activation is permitted.

## Interpretation

Replication of 18-state failure with 27-state success supports a selector-robust 18–27 transition. An 18-state pass indicates selector sensitivity. Any other pattern is sealed without tuning.
