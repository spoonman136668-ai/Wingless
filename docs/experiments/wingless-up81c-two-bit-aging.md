# Wingless UP-81C — two-bit bounded-memory aging

Status: preregistered scientific bounded-memory policy experiment.

Scientific parent: sealed UP-80C `b100e17cbac6354c44b42e9dad25101dd97c40be`.

## Question

UP-80C showed that one-bit CLOCK uses only 3 metadata bytes but behaves like FIFO rather than LRU under locality. Can a tiny two-bit saturating aging signal recover useful locality while retaining much lower metadata cost than LRU?

## Frozen memory budget

- exact recall capacity exactly 16 entries;
- 16 payload bytes per occupied entry;
- no capacity increase;
- no future-query oracle.

## Frozen policy

`aging2_clock`:
- two-bit age per entry, values 0..3;
- successful insertion/update sets age to 3;
- successful query/rehearsal sets age to 3;
- deterministic circular hand;
- on eviction scan:
  - if age > 0, decrement age and advance;
  - evict first entry encountered with age 0;
- new entry receives age 3.

Metadata:
- 32 age bits = 4 bytes;
- clock hand = 1 byte;
- total policy metadata = 5 bytes.

## Controls

Compare unchanged FIFO, LRU, and one-bit CLOCK from UP-80C against `aging2_clock`.

## Frozen workloads

Identical to UP-79C/UP-80C:
- active keys 24 and 32;
- even keys hot, odd keys cold;
- one rehearsal of every hot key;
- churn writes 8 and 16;
- complete final hot and cold queries;
- value vocabulary 32;
- 64 episodes per setting;
- seeds 153M and 154M.

## Metrics

- hot query accuracy and hot-set exact accuracy;
- cold query accuracy and cold-set exact accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

If two-bit aging materially improves hot retention over one-bit CLOCK while staying far below LRU metadata, keep it as the next bounded-memory candidate. If it still behaves like FIFO/CLOCK, the locality signal needs a different mechanism rather than more counter bits.

## Bounds

No adaptive policy selection, no future oracle, no capacity increase, no attention, no result-informed retry, no live activation, no production authority.
