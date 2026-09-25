# Wingless UP-81C — two-bit aging eviction

Status: preregistered scientific bounded-memory policy experiment.

Scientific parent: sealed UP-80C `b100e17cbac6354c44b42e9dad25101dd97c40be`.

## Question

UP-80C showed that one-bit CLOCK costs only 3 metadata bytes but behaves like FIFO rather than preserving LRU's locality advantage. Can a tiny two-bit aging signal preserve rehearsal locality while remaining far cheaper than LRU?

## Frozen memory budget

- exact recall capacity exactly 16 entries;
- key/value payload accounting: 16 bytes per occupied entry;
- no capacity increase;
- no future-query oracle.

## Frozen policies

Controls:
- FIFO unchanged from UP-79C;
- LRU unchanged from UP-79C;
- one-bit CLOCK unchanged from UP-80C.

Experimental `two_bit_aging`:
- 2-bit age per entry, values 0..3;
- one deterministic hand index;
- new insertions begin at age 1;
- update of an existing key sets age 3;
- successful query/rehearsal sets age 3;
- on eviction, inspect from the hand:
  - age 0: evict;
  - age >0: decrement by one and advance;
  - continue until an age-0 slot is reached;
- replacement entry begins at age 1.

## Frozen metadata accounting

- FIFO: 16 bytes;
- LRU: 136 bytes;
- CLOCK: 3 bytes;
- two-bit aging: 4 bytes for 16 two-bit ages + 1 byte hand = 5 bytes.

## Frozen workload

Exactly the UP-80C workload:
- active keys 24 and 32;
- hot set = even original keys;
- cold set = odd original keys;
- one rehearsal query of every hot key;
- churn writes 8 and 16 unseen keys;
- final complete hot and cold queries;
- value vocabulary 32;
- 64 episodes per setting;
- two fixed seed bases.

## Metrics

Hot/cold query accuracy, hot/cold set exact accuracy, entries used, metadata bytes, and total bounded-memory bytes for all four policies.

## Interpretation

The scientific question is whether two-bit aging recovers any of LRU's rehearsal-locality advantage at a bounded 5-byte policy cost. Full-oversubscription failure remains a valid result.

## Bounds

No oracle, no adaptive policy selection, no budget expansion, no result-informed retry, no attention, no live activation, no production authority.
