# Wingless UP-147C — durable recency gradient

Status: preregistered scientific durable-memory replacement experiment.

Scientific parent: sealed controller UP-146C durable-access-recency.

## Question

UP-146C showed that queried durable facts survive while equally cold facts are displaced, and that equal-recency ties are resolved by table order. Does survival vary systematically with how many replacement admissions have occurred since the cohort's last access?

## Frozen durable store

Use the exact 16 lexical durable facts from UP-145C, inserted in the same order into the unchanged two-bit aging memory.

Four four-fact cohorts:
- keys 0..3;
- keys 4..7;
- keys 8..11;
- keys 12..15.

Each cohort is tested independently so every table region occupies every recency arm.

## Fixed replacement pressure

Use the exact four accepted lexical candidate associations as four direct durable-admission writes, in fixed key order 100, 101, 102, 103.

This experiment isolates the durable replacement policy; temporal consolidation is not under test. Every arm receives the same four admissions.

## Fixed refresh arms

Each arm has exactly one four-query refresh of its target cohort:
- before_1: immediately before admission 1;
- before_2: after admission 1, immediately before admission 2;
- before_3: after admissions 1–2, immediately before admission 3;
- before_4: after admissions 1–3, immediately before admission 4.

Thus the refreshed cohort faces respectively 4, 3, 2, or 1 subsequent replacement admissions.

If a target fact has already been displaced before its scheduled refresh, the failed query is recorded and does not restore it.

Total arms: 4 cohorts × 4 refresh positions = 16.

## Measurements

Per arm:
- target cohort and refresh position;
- target facts present immediately before refresh;
- successful refresh queries;
- target final retention;
- non-target durable final retention;
- newly admitted candidate retention;
- final exact-recall entries;
- final aging-memory hand position.

## Interpretation

A systematic survival gradient with fewer post-refresh admissions supports access recency as a graded durable-retention signal. Differences among cohorts at the same recency quantify the residual table-order/hand effect. Lack of recency ordering would falsify the simple graded interpretation.

No memory rule is changed.

## Bounds

Exact recall cap remains 16. Same two-bit aging policy, same initial durable insertion order, same four admissions, exactly four refresh attempts per arm, no adaptive query selection, no semantic priority, no capacity increase, no result-informed retry, no live activation, no production authority.
