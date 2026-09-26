# Wingless UP-147B — tail identity × block-position rotation

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-146B fe75f2ca45f4566b71afa14e18b439812514d9e3.

## Question

UP-146B did not support a simple monotonic habituation curve. Individual five-epoch tail blocks behaved very differently, but tail identity was confounded with block position. Does the interference pattern follow the lexical tail identity, or when that tail appears in the 20-epoch curriculum?

## Frozen starting point and budget

Use the exact UP-146B starting gate and low-switch budget:
- 20 epochs;
- learning rate 0.08;
- 24 new-family updates per epoch;
- 15 old-family rehearsal updates per epoch;
- four final new-family updates;
- four tail shifts {0,6,12,18};
- each shift used exactly five times;
- exactly three between-block switches per arm.

## Arms

Four cyclic block orders:

1. order_0: 0 → 6 → 12 → 18
2. order_1: 6 → 12 → 18 → 0
3. order_2: 12 → 18 → 0 → 6
4. order_3: 18 → 0 → 6 → 12

Each block lasts exactly five consecutive epochs.

Across the four arms, every tail identity occupies every block position exactly once.

## Measurements

For each arm × block:
- block position;
- shift identity;
- mean prefix damage;
- mean anchor recovery;
- mean tail damage;
- mean net epoch change.

For each arm:
- final old held-out/unseen retention;
- final mean old retention;
- final primary/secondary new-family accuracy;
- final mean new-family accuracy.

## Interpretation

If a tail produces similar behavior across different block positions, tail geometry dominates. If behavior follows block position regardless of tail identity, curriculum stage dominates. Mixed behavior would indicate an interaction between lexical geometry and accumulated training state.

No numeric success threshold is introduced; the full frozen matrix is the result.

## Bounds

No adaptive ordering, no extra updates, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
