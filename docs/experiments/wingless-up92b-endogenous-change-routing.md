# Wingless UP-92B — endogenous change-based bounded recall routing

Status: preregistered scientific memory-routing experiment.

Scientific parent: sealed UP-91B evidence `9c2b5c8195bba9589afb488236c5ce6ec9271939`.

## Question

UP-91B showed that a fixed 16-entry recall store can preserve exact sparse information across 256 writes when an explicit salience bit routes only relevant writes. UP-92B asks whether a purely local endogenous signal—detected value change for an already represented key—can route those writes without an externally supplied salience bit.

## Frozen substrate

- same 64-float recurrent carrier;
- same gated-correction threshold 0.15;
- exact-recall capacity exactly 16;
- no attention;
- no state-size increase.

## Frozen arms

1. `fifo_all_writes`
   - admit every write to bounded exact recall.

2. `explicit_target_rewrite`
   - upper-bound control;
   - admit only the known target-key rewrite event.

3. `endogenous_change`
   - before each write, decode that key from the recurrent carrier;
   - declare a change when prior-value presence is true under the existing 0.25 presence threshold and decoded old value differs from the incoming value;
   - admit the new value only when that local change condition is true.

The endogenous arm receives no target/salience bit.

## Frozen task

Sparse target-change stream:
- total writes: 32, 64, 128, 256;
- target key counts: 4, 8, 16;
- every target key is written once near the beginning and rewritten once later to a different value;
- all other writes are unique distractor keys written once;
- only target keys are queried at the end;
- value vocabulary: 32;
- 64 episodes per setting;
- seeds: 127M and 128M.

## Metrics

- target query accuracy;
- whole-target-set exact accuracy;
- recall entries used;
- endogenous admission precision;
- endogenous admission recall relative to true target rewrites;
- false-positive admissions.

## Bounds

No future-query oracle, no external salience bit in the endogenous arm, no capacity increase, no threshold tuning, no result-informed retry, no production authority, no live activation.
