# Wingless UP-51A — sparse binding generalization

Status: preregistered scientific cognition experiment.

Scientific parent: UP-50A seal `8aa5ed497de5a152c634211117e9001ed445b092`.

## Question

UP-50A generalized perfectly across three roles, but training still covered every one of the 27 basis configurations.

UP-51A asks whether the learned reversible rules are genuinely structural: can they be learned from only one-third of the basis states and then transfer to the other two-thirds?

## Frozen design

The three-role / three-value architecture and UP-50A optimizer are unchanged.

Training uses only basis states satisfying:

`(role0 + role1 + role2) mod 3 == 0`

That is exactly 9 of the 27 possible basis configurations.

Training remains single-step only and includes the same two primitives as UP-50A:

- role-0 value swap;
- role swap at position 0.

All other 18 basis configurations are held out completely.

Held-out evaluation on those unseen basis states includes:

- the directly trained primitive types applied to unseen states;
- the never-trained role-swap position 1;
- derived role-1 mutation;
- derived role-2 mutation;
- long mixed programs of lengths 12 and 36.

The matched non-unitary path receives exactly the same sparse training set, optimizer, and budget.

The UP-50A numerical and accuracy thresholds remain unchanged.

## Interpretation

A pass means the learned operators are not dependent on exhaustive enumeration of the state basis; they transfer structurally across unseen bound configurations.

A failure on unseen primitive applications would identify basis coverage as the bottleneck. A failure only on derived or long-program categories would identify composition under sparse supervision as the bottleneck.

Scientific negatives are valid.

## Plain speak

Until now, we taught the rule while showing it every possible starting arrangement.

This time it sees only one-third of them.

If it still works perfectly on the other two-thirds—including operations we never directly taught—then the rule is behaving much more like a reusable algorithm than a memorized table.
