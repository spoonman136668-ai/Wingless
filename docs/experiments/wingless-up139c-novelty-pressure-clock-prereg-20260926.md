# Wingless UP-139C — novelty-pressure consolidation clock

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-138C 5ce5a9dd0776e3054dc16c90a322db9e5ce16caf.

## Question

UP-138C showed that read-only stream activity does not age qualified history, while novel candidate writes do. Is the clock specifically driven by novel/unadmitted candidate pressure, or does any call through the admission process age history?

## Frozen mechanism

Reuse the accepted age-eviction machine unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maxAge = 2;
- lowest-index tie break;
- no semantic/query/future priority.

Each arm uses one recurrence target only.

## Frozen setup

As in UP-138C:
- initialize hot exact recall and bounded history;
- allow original target history to expire via isolated G1/G2 sightings;
- rebuild fresh qualified history with a G3 pair at positions 74 and 75;
- finish G3 so rebuilt history is stored at age 0.

## Matched process-call delay

After G3 boundary, execute exactly 96 calls to the machine's process method in every arm.

Arms differ only in how many calls use a novel one-shot candidate key:
- novel_0_exact_96
- novel_32_exact_64
- novel_64_exact_32
- novel_96_exact_0

The remaining calls update the same already-durable hot key 0 with its correct value. Reusing one exact key prevents the 16-key reuse checkpoint from firing.

Novel candidate calls are spread deterministically across the 96 slots.

Then present the target twice consecutively.

Two target identities:
- 100
- 103.

Seeds:
- 265000000;
- 266000000.

32 episodes per seed.

## Measurements

Per target × arm:
- process calls;
- novel candidate calls;
- exact-memory update calls;
- candidate boundaries crossed;
- checkpoint count;
- target-history presence and age before the second pair;
- admission rate;
- final target accuracy;
- hot accuracy;
- false admissions;
- max current/history occupancy;
- recall entries used;
- panic rate.

## Interpretation

If 0/32/64 novel candidates preserve history at ages 0/1/2 while 96 novel candidates expire it despite all arms making 96 process calls, the consolidation clock is specifically novelty/candidate-pressure based rather than generic process-call time.

No post-result threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
