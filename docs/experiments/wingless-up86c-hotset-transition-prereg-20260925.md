# Wingless UP-86C — online hot-set transition refinement

Status: preregistered scientific bounded-memory capacity experiment.

Scientific parent: sealed UP-85C 8586ad159bde0d49ffad2e1f3e61e306670a4b3c.

## Question

UP-85C showed exact online hot-set retention at size 12 and failure at size 16. Where does the transition occur?

## Frozen system and workload

Identical to UP-85C:
- active keys 32;
- churn writes 24;
- exact recall cap 16;
- value vocabulary 32;
- online hot-key queries after every 4 original writes and every 4 churn writes;
- FIFO, LRU, and two-bit-aging unchanged;
- two-bit-aging metadata 5 bytes;
- seeds 163M and 164M.

Only the hot-set-size ladder changes to:
- 13
- 14
- 15
- 16

Hot keys remain the first N even-numbered original keys.

## Metrics

For every policy x hot-size cell:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- cold-set exact accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes.

## Bounds

No capacity expansion, no adaptive policy selection, no future-query labels, no policy tuning, no attention, no result-informed retry, no live activation, no production authority.
