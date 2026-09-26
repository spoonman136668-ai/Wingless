# Wingless UP-126C — second overflow wave after stale-pool exhaustion

Status: preregistered scientific bounded-memory boundary experiment.

Scientific parent: sealed UP-125C 9e15a003ab99adf28ebe28777c52d5bf5e86931b.

## Question

UP-125C showed that age-only eviction safely absorbs a first 16-candidate overflow wave by reclaiming all 16 stale age-2 set-A history slots. What happens when another qualified overflow wave arrives before the four valid targets recur, after that stale reservoir has been exhausted?

## Frozen mechanism

Use the exact UP-125C mechanism:
- exact recall cap 16;
- 32-entry current table;
- 32-entry history table;
- 128-byte admission sidecar;
- age-only full-table replacement;
- maximum age selected;
- lowest table index tie-break;
- no semantic/query/future priority.

No rule changes are permitted.

## Sequence

1. Reproduce exact 32/32 saturation from UP-125C.
2. First overflow wave:
   - exactly 16 fresh qualified nonpersistent candidates;
   - two sightings each;
   - consumes the complete 32-write generation;
   - expected to reclaim the 16 age-2 stale set-A slots under the unchanged rule.
3. Second overflow wave arrives immediately, before valid-target recurrence.
4. Second-wave qualified-candidate count Q2 ∈ {0,1,4,8}.
   - each gets exactly two sightings;
   - remaining generation writes are unique one-shot candidates.
5. Only after the second wave do the four valid targets recur twice each.
6. Apply unchanged 12,288-write churn and final evaluation.

Seeds:
- 239000000;
- 240000000.

32 episodes per seed.

## Metrics

Per Q2:
- panic rate;
- first-wave and second-wave evictions;
- evictions by age;
- valid-target admission rate;
- final valid-target accuracy;
- hot accuracy;
- target16 accuracy and exact rate;
- nonpersistent retention/admission;
- false admissions;
- max history/current occupancy;
- recall entries used.

## Interpretation

If Q2>0 begins reducing target admission while Q2=0 remains exact, the rule's safe capacity is not “overflow in general”; it depends on the presence of genuinely stale older evidence. This would locate the boundary between safe forgetting and destructive forgetting under unchanged memory.

A negative result is valid and the rule will not be tuned after observation.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
