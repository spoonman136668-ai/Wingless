# Wingless UP-113C — exact probation mixed legitimate-recurrence working set

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-112C 11d807704fc0eab679530497bfceba1d0e0f2c86.

## Question

UP-112C established that exact32 has the intended 31-write recurrence window and eliminates Bloom false admissions at equal admission-memory cost. Can exact32 admit legitimate recurring non-hot candidates while preserving a mixed 16-entry working set under long one-shot churn?

## Frozen memory system

Use the exact UP-111C/UP-112C exact32 mechanism:
- exact recall cap 16;
- 32-entry uint32 probation table = 128 bytes;
- generation interval 32;
- one-byte generation counter;
- phase-free slot-reuse checkpoint;
- no query-derived admission evidence;
- no future oracle.

## Working set

Exact-memory target occupancy:
- 12 persistent hot keys;
- 4 legitimate recurring candidate keys;
- total target working set = 16.

Initial phase:
- load 32 original keys deterministically;
- promote the 12 hot keys using two writes within the ordinary recurrence window;
- present each of 4 recurring candidate keys once to probation and then repeat it after a fixed in-window gap so it becomes legitimately admissible.

Recurring candidate gaps:
- 8
- 16
- 24
- 31 filtered writes.

Each cell uses one fixed gap for all four recurring candidates.

## Churn

After the mixed working set is formed:
- unique one-shot churn ladder:
  - 24,576
  - 49,152
  - 98,304 writes
- ordinary hot/recurring queries every four churn writes;
- 32 episodes per seed;
- seeds 213M and 214M.

## Metrics

- admission success for all four recurring candidates;
- recurring-candidate value accuracy;
- 12-hot accuracy;
- whole 16-entry mixed-set exact accuracy;
- one-shot false admissions;
- recall entries used;
- checkpoint count;
- policy metadata bytes.

## Interpretation

Exact mixed-set retention with legitimate recurring admissions would establish that exact32 is a selective recurrence detector rather than a one-shot rejection filter.

## Bounds

No capacity increase, no adaptive gap, no query-derived admission, no future oracle, no phase labels, no result-informed retry, no live activation, no production authority.
