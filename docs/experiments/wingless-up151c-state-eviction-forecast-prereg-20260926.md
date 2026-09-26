# Wingless UP-151C — state-derived eviction forecast

Status: preregistered scientific durable-memory forecast experiment.

Scientific parent: sealed UP-150C e758b9094ce0a113dee157424b3e5e8415c582a6.

## Question

After sparse refreshes alter two-bit ages, are current replacement-hand position and the 16 age counters sufficient to predict the very next eviction exactly, without semantic identity or future information?

## Frozen design

Reuse the accepted up81cAging memory unchanged:
- exact recall cap 16;
- same 16 initial durable facts;
- starting hands 0, 4, 8, 12;
- 16 unique admissions per arm.

Before admissions 3, 6, 9, 12, and 15, issue one fixed refresh query to original keys:
- admission 3 -> key 3;
- admission 6 -> key 7;
- admission 9 -> key 11;
- admission 12 -> key 15;
- admission 15 -> key 0.

If a scheduled key has already been evicted, the query simply misses; no substitute is chosen.

## Frozen predictor

Immediately before each admission, read only:
- current hand;
- the 16 two-bit age counters.

Do not read key identity, lexical class, values, or future refresh schedule.

Predict the next replacement slot by applying the documented two-bit aging scan to a copy of those age counters. The predictor does not mutate Wingless memory.

Then perform the actual write and locate the new key's slot.

## Measurements

Per hand × admission:
- scheduled refresh key and whether it hit;
- current hand before write;
- predicted slot;
- actual slot;
- prediction correctness;
- evicted key;
- original facts remaining.

Primary metric: exact next-eviction prediction accuracy across all 64 admissions.

## Bounds

No semantic priority, adaptive refresh, capacity increase, policy change, future oracle, result-informed retry, live activation, or production authority.
