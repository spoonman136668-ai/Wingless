# Wingless UP-91B — selective bounded-recall routing

Status: preregistered scientific memory-routing experiment.

Scientific parent: sealed UP-SQ0 evidence `8de9e11f26cc44c75f51d26cb606f97e35816479`.

## Question

SQ0 showed that a 16-entry exact-recall channel closes sparse associative tasks inside its capacity but fails when all sequence content is indiscriminately admitted beyond that cap. Can an explicit local salience bit route only relevant information into the same fixed 16-entry store and preserve exact recall across much longer distractor streams?

## Frozen substrate

- same 64-float recurrent carrier as SQ0;
- same gated-correction rule and threshold 0.15;
- exact-recall capacity remains exactly 16 entries;
- no state-size increase;
- no attention;
- no oracle future-query access beyond the explicit per-write salience bit supplied by the task.

## Frozen arms

1. `fifo_all_writes`: every write is admitted to bounded exact recall.
2. `salience_only`: only writes carrying the explicit salience bit are admitted.

Both arms update the same recurrent carrier on every event.

## Task 1 — selective copy

- lengths: 32, 64, 128, 256;
- exactly 25% of positions are marked salient by a deterministic schedule;
- query only marked positions at the end;
- alphabet: 16;
- 64 episodes per length;
- seeds: 123M and 124M.

## Task 2 — sparse-query distractor stream

- total writes: 32, 64, 128, 256;
- salient key counts: 4, 8, 16;
- salient keys are written once early and overwritten once later with latest-value-wins semantics;
- remaining events are unmarked distractor writes to distinct keys;
- only salient keys are queried at the end;
- value vocabulary: 32;
- 64 episodes per setting;
- seeds: 123M and 124M.

## Metrics

- query accuracy;
- whole-query-set exact accuracy;
- recall entries used;
- recurrent-state bytes;
- exact-recall bytes.

## Bounds

No capacity increase, no threshold tuning, no retrospective salience inference, no result-informed retry, no production authority, no live activation.
