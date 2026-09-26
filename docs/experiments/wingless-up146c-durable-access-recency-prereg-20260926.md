# Wingless UP-146C — durable lexical access recency

Status: preregistered scientific lexical-memory integration experiment.

Scientific parent: sealed UP-145C 365893cb14058ecded93d22719185878ffd08316.

## Question

UP-145C preserved all 12 repeatedly queried durable lexical facts while two of four cold, unqueried facts were displaced by the same two new durable admissions. Is that retention difference governed by access recency rather than lexical class, subject identity, or fixed table position?

## Frozen consolidation schedule

Reuse UP-145C unchanged:
- exact recall cap 16;
- history table 32;
- admission sidecar 128 bytes;
- three 32-candidate generations;
- identical two-stage, cross-boundary, delayed-two-stage, and isolated candidate schedules;
- identical mirrored assignments A and B;
- identical two successful new durable admissions;
- identical one-shot novelty;
- no semantic/query priority.

## Frozen durable lexical facts

The same 16 initial facts from UP-145C are used in the same table order.

Four cold-block arms rotate which four facts receive no durability-refresh queries:

1. cold_0_3
   - quin STORE: lodges, stashes, caches, files.

2. cold_4_7
   - quin OBSERVE: scans, checks, views, monitors.

3. cold_8_11
   - quin REPORT: relays, announces, cites, summarizes.

4. cold_12_15
   - rue STORE: lodges, stashes, caches, files.

After every four candidate calls, all 12 non-cold durable facts are queried. Exactly four facts are cold in every arm.

## Measurements

Per assignment × cold block:
- hot queried durable accuracy;
- cold unqueried durable accuracy;
- per-durable-fact retention;
- candidate temporal-role admission/retention;
- recall entries used;
- max current/history occupancy;
- false one-shot admissions;
- panic rate.

## Interpretation

If whichever four facts are left unqueried preferentially absorb the two durable-store evictions while the 12 queried facts remain intact, access recency—not lexical class or fixed table region—is the controlling durable-retention signal.

No memory-policy change is made.

## Bounds

No capacity increase, no new replacement policy, no semantic priority, no query priority, no adaptive refresh, no result-informed retry, no live activation, no production authority.
