# Wingless UP-103C — hot-repeat reinforcement order

Status: preregistered scientific bounded-memory timing diagnostic.

Scientific parent: sealed UP-102C 1eaba94eccfcb3cc0d1fd441d907077b00e5587d.

## Question

UP-102C showed that reverse initial order loses two hot entries during the promotion phase unless the initially resident hot entries are reinforced by an early query sweep. Is the failure caused by the order in which hot second writes reinforce resident versus missing hot keys?

## Frozen memory mechanism

Use the exact UP-100C trigger8 system:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- trigger after 8 successful repeat admissions;
- post-trigger counter 8;
- no extra diagnostic query sweeps;
- no query-derived admission labels;
- no future oracle.

## Frozen workload

- 32 original keys;
- all 16 even keys are hot;
- initial orders:
  1. ascending
  2. reverse
  3. evens_then_odds
- then one second write to every hot key;
- then 12,288 unique one-shot churn writes;
- ordinary hot-query cadence after every four promotion writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 195M and 196M.

## Frozen hot second-write orders

1. ascending
   - 0,2,4,...,30

2. reverse
   - 30,28,26,...,0

3. alternating_ends
   - 0,30,2,28,4,26,...,14,16

No order is selected adaptively from memory state.

## Metrics

For each initial-order x second-write-order cell:
- resident hot-key count immediately before churn, measured without a query;
- trigger-fired episode count;
- mean trigger index;
- false-positive churn admissions;
- final hot query accuracy;
- final hot-set exact accuracy;
- final recall entries used.

## Interpretation

If reverse initial order is repaired by reinforcing high resident hot keys before admitting missing low hot keys, promotion-time reinforcement ordering is causal. If order does not matter, another aging-state property must explain the UP-102C effect.

## Bounds

No trigger change, no capacity expansion, no filter change, no extra diagnostic reads, no adaptive ordering, no result-informed retry, no live activation, no production authority.
