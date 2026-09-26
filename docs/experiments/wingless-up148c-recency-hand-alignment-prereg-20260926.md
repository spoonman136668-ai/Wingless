# Wingless UP-148C — durable recency × hand alignment

Status: preregistered scientific durable-memory replacement experiment.

Scientific parent: sealed UP-147C ac6f4f4e0aff80939ced6251a2b6461ad67ff2c7.

## Question

UP-147C rejected a purely scalar recency model: delayed refresh harmed only the durable cohort already sitting in the replacement path, while the other table regions survived every refresh arm. Does that vulnerability move with slot/hand alignment when the same lexical identities are inserted into different table regions?

## Frozen target facts

The target cohort is always the same four accepted lexical facts:
- key 0: quin lodges → STORE;
- key 1: quin stashes → STORE;
- key 2: quin caches → STORE;
- key 3: quin files → STORE.

No target identity or class changes across arms.

## Frozen initial durable contents

Use the same 16 durable lexical facts and unchanged two-bit aging memory from UP-145C..UP-147C.

Four deterministic cyclic insertion rotations:
- rotation_0: original insertion order; target keys 0..3 occupy the first slot region;
- rotation_4;
- rotation_8;
- rotation_12.

The rotation changes only which durable facts occupy each table slot. Capacity, hand logic, ages, fact multiset, and write count are unchanged.

## Frozen refresh arms

Exactly as UP-147C:
- before_1: refresh target immediately before admission 1;
- before_2: after admission 1, before admission 2;
- before_3: after admissions 1–2, before admission 3;
- before_4: after admissions 1–3, before admission 4.

If a target fact has already been evicted, its failed refresh query is recorded and does not restore it.

## Frozen replacement pressure

Use the same four direct durable admissions, in fixed order:
- key 100;
- key 101;
- key 102;
- key 103.

All 16 arms receive all four admissions.

## Measurements

Per rotation × refresh arm:
- initial table slots of target keys 0..3;
- hand position before admissions;
- target facts present immediately before refresh;
- successful refresh queries;
- final target retention and accuracy;
- non-target durable retention;
- new-candidate retention;
- final hand;
- exact recall entries used.

## Interpretation

If the delayed-refresh vulnerability follows the target cohort when it is placed in the hand's early replacement region, rather than following keys 0..3 regardless of slot placement, slot/hand alignment is causal. If keys 0..3 remain uniquely vulnerable in all rotations, identity-linked effects remain.

No memory policy is changed.

## Bounds

Exact recall cap 16. Same two-bit aging policy, same 16 durable facts, same four admissions, exactly four refresh attempts per arm, no semantic priority, no adaptive refresh, no capacity increase, no result-informed retry, no live activation, no production authority.
