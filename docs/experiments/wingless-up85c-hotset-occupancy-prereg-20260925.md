# Wingless UP-85C — online hot-working-set occupancy

Status: preregistered scientific bounded-memory capacity experiment.

Scientific parent: sealed UP-84C `51805b22b80b26e927e12dd8b3c66e1282e3cb55`.

## Question

UP-84C preserved all 12 hot items at active=24 but only 81.25% of 16 hot items at active=32. Where is the exact online hot-working-set boundary under the fixed 16-entry memory budget?

## Frozen memory system

- exact recall capacity exactly 16 entries;
- payload accounting 16 bytes per occupied entry;
- FIFO, LRU, and two-bit-aging policies unchanged;
- two-bit-aging metadata remains 5 bytes;
- no future-query oracle.

## Frozen workload

- total original active keys: 32;
- hot-set sizes: 4, 8, 12, 16;
- hot keys are the first N even-numbered original keys;
- all other original keys form the cold control;
- initial keys written in deterministic order;
- after every 4 original writes, query every currently loaded hot key once;
- 24 unseen churn writes;
- after every 4 churn writes, query every hot original key once;
- missing keys remain missing;
- final complete hot-set and cold-control queries;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 163M and 164M.

## Metrics

For every policy x hot-set-size cell:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- cold-set exact accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

The experiment maps the effective online working-set threshold separately from nominal capacity. It does not change capacity or policy parameters.

## Bounds

No adaptive policy selection, no capacity expansion, no future-query labels, no policy tuning, no attention, no result-informed retry, no live activation, no production authority.
