# Wingless UP-149B — four-tail terminal causal map

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-148B f5c8bb8c67260f266b74513a118ab4a09c7dce76.

## Question

UP-148B proved that terminal tail identity can causally change both old-family retention and new-family accuracy from an identical 15-epoch state. What is the complete causal response across all four preregistered tail identities?

## Frozen common prefix

Every arm starts from the same gate and runs exactly:
- epochs 1–5: shift 0;
- epochs 6–10: shift 6;
- epochs 11–15: shift 12.

Every epoch:
- learning rate 0.08;
- 24 new-family updates;
- 15 old-family rehearsal updates;
- four final new-family updates after the old anchor.

## Terminal arms

Only epochs 16–20 differ:

- terminal_0: shift 0 for five epochs;
- terminal_6: shift 6 for five epochs;
- terminal_12: shift 12 for five epochs;
- terminal_18: shift 18 for five epochs.

All arms have identical update counts and identical history through epoch 15.

## Measurements

For each terminal identity:
- old retention at the common-prefix boundary;
- terminal mean prefix damage;
- terminal mean anchor recovery;
- terminal mean tail damage;
- terminal mean net epoch change;
- final old held-out/unseen retention;
- final primary/secondary new-family accuracy;
- final mean old retention;
- final mean new-family accuracy.

## Interpretation

The exact four-point response establishes whether the terminal effect is a two-class safe/damaging split or a graded identity-specific geometry. No numeric success threshold and no post-result tail selection are introduced.

## Bounds

No adaptive tail selection, no extra updates, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
