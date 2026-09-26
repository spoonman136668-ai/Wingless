# Wingless UP-147B — tail identity × curriculum-position permutation

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-146B fe75f2ca45f4566b71afa14e18b439812514d9e3.

## Question

UP-146B rejected simple monotonic habituation, but each of the four tail identities was always tested in one fixed five-epoch block position. Are the different tail trajectories properties of the tail identities themselves, or consequences of where the block occurs in the 20-epoch curriculum?

## Frozen starting point and budget

Use the exact UP-146B starting gate and low-switch training budget:
- 20 epochs;
- learning rate 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family rehearsal updates per epoch;
- exactly four final new-family updates;
- shifts {0,6,12,18};
- each shift used exactly five consecutive epochs;
- exactly three block transitions per schedule.

No training budget, family, projector, or memory change is permitted.

## Preregistered schedules

Run four independent schedules from the same frozen starting gate:

- rotation_0: blocks 0 → 6 → 12 → 18
- rotation_1: blocks 6 → 12 → 18 → 0
- rotation_2: blocks 12 → 18 → 0 → 6
- rotation_3: blocks 18 → 0 → 6 → 12

Each block lasts exactly five consecutive epochs.

Across the four schedules:
- every tail identity appears exactly once in block positions 1,2,3,4;
- every tail identity is used for exactly five epochs per schedule in which it appears;
- every schedule contains the same tail multiset and total updates.

## Measurements

Per epoch:
- schedule;
- block position 1..4;
- tail shift;
- repetition 1..5;
- prefix damage;
- anchor recovery;
- tail damage;
- net epoch change.

Per tail identity × block position cell:
- mean prefix damage over the five repetitions;
- mean anchor recovery;
- mean tail damage;
- mean net epoch change.

Per schedule:
- final mean old retention;
- final mean new-family accuracy.

## Interpretation

If the same tail identity shows similar behavior across block positions, tail identity is the dominant mechanism. If behavior follows block position regardless of identity, curriculum state is dominant. A mixed pattern indicates an identity × position interaction.

No numeric success threshold is introduced; the exact 4×4 cell structure is the result.

## Bounds

No adaptive block order, no extra updates, no new memory, no projector recomputation, no architecture change, no threshold search, no result-informed retry, no live activation, no production authority.
