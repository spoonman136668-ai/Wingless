# Wingless UP-79C — bounded exact-memory eviction under oversubscription

Status: preregistered scientific memory-capacity experiment.

Scientific parent: sealed UP-78C `4644c0ba295515a9d9212b7184ea0b28ff2341c7`.

Dynamic two-choice carrier routing is closed. The downstream approximate carrier is frozen to fixed 8×8 geometry. This experiment pivots C to bounded exact-memory policy.

## Question

The exact store is capped at 16 entries, while useful language tasks can exceed 16 active bindings. Can a local access-based eviction policy preserve recently relevant bindings better than FIFO without future-query information?

## Frozen memory budget

- exact recall capacity exactly 16 entries;
- key/value payload unchanged;
- no capacity increase;
- no future-query oracle.

## Frozen policies

1. `fifo`
   - insertion order eviction;
   - updates overwrite in place without changing order;
   - queries do not change order.

2. `lru`
   - each admitted store/update touches the entry;
   - each successful rehearsal query touches the entry;
   - evict least-recently-touched entry;
   - deterministic monotonically increasing local counter only.

Metadata bytes are reported separately and are not hidden inside the 16-entry payload budget.

## Frozen workloads

Active keys:
- 24
- 32

For each episode:
1. store one value for every active key;
2. choose the deterministic hot set as the even-numbered original keys and the cold control as the odd-numbered original keys (12/12 when active_keys=24; 16/16 when active_keys=32);
3. issue one rehearsal query to every hot key;
4. append churn writes to 8 or 16 previously unseen keys, forcing eviction;
5. final-query the complete hot set;
6. final-query the equally sized cold control.

Value vocabulary: 32.
64 episodes per setting.
Seeds: 145M and 146M.

## Metrics

- hot-set final query accuracy;
- hot-set exact accuracy;
- cold-set final query accuracy;
- cold-set exact accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

This experiment does not seek one universally best eviction rule. It establishes the cost/benefit boundary of exploiting observed access locality under a hard 16-entry cap.

## Bounds

No future-query labels, no capacity increase, no adaptive policy selection, no attention, no result-informed retry, no production authority, no live activation.


## Frozen byte accounting

- exact entry payload accounting: 16 bytes per occupied entry, matching the prior bounded-recall accounting;
- full 16-entry payload budget: 256 bytes;
- FIFO policy metadata: 16 bytes total for fixed-array head/count bookkeeping;
- LRU policy metadata: 8-byte monotonic counter plus 8-byte last-touch stamp per capacity slot = 136 bytes;
- reported total bounded-memory bytes = occupied payload bytes + frozen policy metadata bytes.
