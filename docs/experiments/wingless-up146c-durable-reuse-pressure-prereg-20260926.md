# Wingless UP-146C — durable reuse pressure

Status: preregistered scientific lexical durable-memory retention experiment.

Scientific parent: sealed UP-145C 365893cb14058ecded93d22719185878ffd08316.

## Question

UP-145C preserved repeatedly queried durable lexical facts at 100% while four unqueried durable controls fell to 50% under the same full 16-entry durable store plus new admissions. Is that retention difference caused by access/reuse pressure rather than lexical identity?

## Frozen durable groups

Use the same 16 initial durable facts and exact recall policy as UP-145C.

Two identity-matched four-fact STORE groups:
- quin_store4: quin × {lodges, stashes, caches, files}, keys 0..3;
- rue_store4: rue × {lodges, stashes, caches, files}, keys 12..15.

Both groups have the same class and verb set.

The remaining quin facts, keys 4..11, stay as background durable controls.

## Fixed protection arms

Every four candidate calls, issue exactly four read-only durable-memory queries:

1. protect_quin_store4
   - query keys 0..3.

2. protect_rue_store4
   - query keys 12..15.

The query budget is identical. Queries do not advance the admission counter and do not alter values.

## Candidate stream

Reuse UP-145C unchanged:
- 96 candidate calls;
- three 32-candidate generations;
- same four lexical candidates;
- same two mirrored temporal-role assignments A and B;
- same unique one-shot novelty filler;
- same two-stage, cross-boundary, delayed-two-stage, and isolated schedules.

## Measurements

Per assignment × protection arm:
- quin_store4 final accuracy;
- rue_store4 final accuracy;
- protected-group accuracy;
- unprotected matched-group accuracy;
- background quin key 4..11 accuracy;
- candidate admission/retention by temporal role;
- exact recall entries used;
- false one-shot admissions;
- current/history maxima;
- panic rate.

## Interpretation

If whichever matched group receives the fixed query budget is preferentially retained, access/reuse pressure is causal and the UP-145C hot/cold difference is not lexical identity. If rue remains weak even when protected, identity/layout effects remain.

No memory policy is changed and no priority signal other than the existing query-driven aging behavior is added.

## Bounds

Exact recall cap remains 16. Same two-bit aging policy, same admission mechanism, same candidate schedule, no semantic/query priority inside consolidation, no memory increase, no future oracle, no adaptive query selection, no result-informed retry, no live activation, no production authority.
