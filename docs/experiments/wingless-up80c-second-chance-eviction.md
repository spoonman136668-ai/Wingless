# Wingless UP-80C — low-metadata second-chance eviction

Status: preregistered scientific bounded-memory policy experiment.

Scientific parent: sealed UP-79C `525487c6011abc396948bee12f032398f201eb89`.

## Question

UP-79C showed that LRU can double hot-set retention under moderate churn, but costs 136 metadata bytes and still collapses under full churn. Can a low-metadata CLOCK/second-chance policy capture most of the locality benefit at a much lower metadata cost?

## Frozen memory budget

- exact recall capacity exactly 16 entries;
- key/value payload accounting: 16 bytes per occupied entry;
- no capacity increase;
- no future-query oracle.

## Frozen policies

1. `fifo`
   - same UP-79C behavior.

2. `lru`
   - same UP-79C behavior.

3. `clock_second_chance`
   - one reference bit per entry;
   - insertion/update sets the entry reference bit;
   - successful rehearsal/final query sets the reference bit;
   - on eviction, advance a deterministic clock hand:
     - if reference bit is 1, clear it and continue;
     - evict the first entry encountered with reference bit 0;
   - new entry receives reference bit 1.

## Frozen metadata accounting

- FIFO metadata: 16 bytes;
- LRU metadata: 136 bytes;
- CLOCK metadata: 2 bytes for 16 reference bits + 1 byte clock-hand index = 3 bytes;
- total bounded-memory bytes = occupied payload bytes + policy metadata bytes.

## Frozen workloads

Same UP-79C workload:
- active original keys: 24 and 32;
- hot set = even-numbered original keys;
- cold control = odd-numbered original keys;
- one rehearsal query to every hot key;
- churn writes: 8 and 16 previously unseen keys;
- final query complete hot set;
- final query complete cold control;
- value vocabulary 32;
- 64 episodes per setting;
- seeds 149M and 150M.

## Metrics

- hot final query accuracy;
- hot-set exact accuracy;
- cold final query accuracy;
- cold-set exact accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

If CLOCK preserves most of LRU's hot-set benefit at moderate churn with substantially lower metadata, prefer CLOCK for downstream bounded exact memory. If not, retain the simpler FIFO/LRU tradeoff characterization.

## Bounds

No future-query labels, no adaptive policy selection, no capacity increase, no attention, no result-informed retry, no production authority, no live activation.
