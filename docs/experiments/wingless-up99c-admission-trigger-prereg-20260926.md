# Wingless UP-99C — admission-event triggered role separation

Status: preregistered scientific bounded-memory mechanism experiment.

Scientific parent: sealed UP-98C 89f6fd97bff95723c248a1a34a07ca43e163d864.

## Question

UP-98C achieved exact 16-item hot retention through 12,288 churn writes by clearing the filter after the explicit promotion phase and restarting at phase offset 8. Can the same separation be triggered using only successful repeat-admission events, without a workload phase label?

## Frozen memory and workload

Identical to UP-98C:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit legacy two-hash seen-once filter;
- generation interval 32;
- 32 original keys;
- 16 hot even keys;
- each hot key receives one second write before churn;
- 12,288 unique one-shot churn writes;
- hot queries after every four hot second writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 187M and 188M.

## Repeat-admission event

A repeat admission is counted only when:
- the key is not currently in exact memory;
- exact memory is full;
- the seen-once filter reports the key as previously seen;
- the write is admitted.

Existing exact-memory updates do not count.

## Frozen arms

1. continuous_control
   - no one-shot role-separation clear.

2. trigger4
3. trigger8
4. trigger12

For trigger arms:
- count successful repeat admissions from the first filtered write onward;
- when the fixed count is reached for the first time, clear only the seen-once filter;
- set the generation counter to 8;
- never perform this special trigger again in the episode;
- ordinary 32-write generation clearing continues afterward.

No phase label, query label, key identity, or future information is used.

## Metrics

For every arm:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- rejected churn writes;
- false-positive churn admissions;
- first false-positive churn index;
- episodes with false positives;
- trigger-fired episodes;
- mean filtered-write index at trigger;
- automatic generation reset count;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

If trigger8 reproduces the UP-98C exact result without a phase label, the diagnostic separation can be converted into an event-driven online rule. Trigger4/12 identify sensitivity to firing too early or too late.

## Bounds

No explicit promotion/churn phase label, no query-derived trigger, no exact-memory expansion, no filter-width/hash/interval change, no adaptive threshold, no future oracle, no result-informed retry, no live activation, no production authority.
