# Wingless UP-91C — continuous seen-once generations

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-90C 68d924c4ddc53e4bcc3145cc5fa1eb5ef11d1539.

## Question

UP-90C reduced long-churn collisions sharply, but reset_32 still carried stale pre-churn seen-once evidence into the churn phase and produced one false admission per episode. Does counting filter generations continuously from the first filtered write remove that residual collision without using phase labels?

## Frozen memory and workload

Identical to UP-90C:
- exact recall cap 16;
- two-bit-aging replacement;
- 256-bit two-hash seen-once filter;
- 32 original keys;
- 16 hot even keys;
- one second write for every hot key;
- 192 unique one-shot churn writes;
- hot queries after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 171M and 172M.

## Frozen arms

1. phase_local_32
   - exact UP-90C reset_32 control.
   - reset counter begins at the first churn write.

2. continuous_16
3. continuous_32
4. continuous_64

For continuous arms:
- the generation counter begins on the first unseen write presented to the seen-once filter after exact memory is full;
- existing exact-memory updates do not increment the counter;
- every filtered unseen write increments it, whether admitted or rejected;
- after processing the write that reaches the fixed interval, clear only the seen-once filter and reset the counter to zero;
- exact memory and aging state are never reset.

Metadata:
- 256-bit filter + two-bit aging + one-byte generation counter = 38 bytes.

## Metrics

For every arm:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- rejected one-shot writes;
- false-positive churn admissions;
- filter reset count;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

Exact retention under a continuous interval would show that stale evidence at phase boundaries was the residual C90 problem. Degradation at very short intervals would reveal the opposing cost: forgetting legitimate second-write evidence before it can be reused.

## Bounds

No phase label in continuous arms, no query-derived reset, no future oracle, no exact-memory expansion, no adaptive interval, no threshold tuning, no result-informed retry, no live activation, no production authority.
