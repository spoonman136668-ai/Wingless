# Wingless UP-150C — eviction-horizon map

Status: preregistered scientific durable-memory replacement experiment.

Scientific parent: sealed UP-149C d1e1d0794501e7820d24fed83589ed26f47baa8d.

## Question

UP-149C established that durable-memory vulnerability follows the replacement hand. What is the exact deterministic eviction horizon across a full 16-admission cycle?

## Frozen design

Use the accepted up81cAging memory unchanged:
- exact recall cap 16;
- two-bit aging;
- same 16 durable lexical facts in their original insertion order;
- no refresh/query of durable entries after initialization.

Test four starting hand positions:
- 0;
- 4;
- 8;
- 12.

For each hand, inject exactly 16 novel direct writes with unique keys 1000..1015 and deterministic class values j mod 3.

After each admission, inspect presence with find only; do not query, because query refreshes age and would alter the mechanism.

## Measurements

Per hand × admission number:
- newly evicted original key;
- newly evicted original slot;
- original facts remaining;
- current hand;
- total entries.

## Preregistered expectation

Exactly one original fact is evicted per admission, following cyclic slot order from the starting hand. After 16 admissions, zero original facts remain.

This expectation does not control harness acceptance.

## Bounds

No refresh, no semantic priority, no query priority, no capacity increase, no adaptive policy, no result-informed retry, no live activation, no production authority.
