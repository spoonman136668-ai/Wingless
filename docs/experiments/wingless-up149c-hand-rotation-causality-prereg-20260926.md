# Wingless UP-149C — replacement-hand rotation causality

Status: preregistered scientific durable-memory replacement experiment.

Scientific parent: sealed UP-148C 84bfd1dde558aa5efdecb13f4e9f17e290c48c08.

## Question

UP-148C showed that delayed-refresh vulnerability follows the slot region occupied by the target cohort. If target facts remain fixed in slots 0..3 while only the initial replacement hand is rotated, does the vulnerable region move with the hand?

## Frozen state

Use the same 16 durable lexical facts in the original insertion order.
Target facts remain keys 0..3 and therefore slots 0..3 in every arm.

Rotate only the initial two-bit-aging replacement hand:
- hand_0;
- hand_4;
- hand_8;
- hand_12.

No ages, entries, identities, classes, or capacity are changed.

## Refresh arms

Exactly as UP-147C/148C:
- before admission 1;
- before admission 2;
- before admission 3;
- before admission 4.

Use the same four direct durable admissions in the same order.

## Measurements

Per hand × refresh arm:
- target slots;
- starting hand;
- target facts present before refresh;
- successful target refreshes;
- final target retention;
- non-target durable retention;
- new-candidate retention;
- final hand.

## Interpretation

If vulnerability occurs only when the rotated hand reaches target slots before refresh, the hand position is causal. If slot 0..3 remains vulnerable independent of starting hand, a fixed slot effect remains.

## Bounds

Exact recall cap 16; unchanged two-bit-aging policy; same 16 durable facts; same four admissions; no semantic priority, adaptive refresh, capacity increase, result-informed retry, live activation, or production authority.
