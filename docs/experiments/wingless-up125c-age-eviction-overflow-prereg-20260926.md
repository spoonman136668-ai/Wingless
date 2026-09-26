# Wingless UP-125C — age-only history eviction under overflow

Status: preregistered scientific bounded-memory mechanism experiment.

Scientific parent: sealed UP-124C b78728a9f3247f5c09f88f8b5f42b20c33b465a6.

## Question

UP-124C established a hard overflow at the first qualified candidate beyond 32/32 history occupancy. Without adding memory, can one deterministic age-only replacement rule reclaim stale history slots while preserving valid recurring targets?

## Frozen memory and admission mechanism

Unchanged:
- exact recall cap: 16;
- current table: 32 × uint16;
- history table: 32 × uint16;
- total admission sidecar: 128 bytes;
- 14-bit candidate keys;
- history ages 0/1/2;
- generation interval: 32;
- 12 hot residents;
- existing current-table qualification and exact history match;
- two-bit exact-memory aging;
- no semantic priority;
- no query-derived admission;
- no future oracle.

## Preregistered replacement rule

Only history insertion when all 32 history slots are occupied changes.

When inserting a qualified history record into a full history table:
1. find the maximum encoded history age currently present;
2. evict the first/lowest-index entry having that maximum age;
3. insert the incoming record at age 0 into that slot.

No key identity, value, semantic class, query count, future recurrence, or target label is consulted.

The tie-break is fixed before execution and never adaptive.

## Saturation construction

Use the exact UP-124C saturation path:
- generation 1: 16 qualified nonpersistent set-A candidates;
- generation 2: four valid recurring targets + 12 qualified nonpersistent set-B candidates;
- history reaches exactly 32/32 occupied.

At the next boundary, set-A evidence is age 2 while target/set-B evidence is age 1.

## Overflow arms

Fresh qualified nonpersistent candidates Q ∈ {0, 1, 8, 16}.

Each candidate receives exactly two sightings in the overflow generation. Remaining writes are unique one-shot candidates.

Q=16 consumes all 16 stale age-2 slots available at that boundary if the preregistered rule behaves as intended.

After overflow:
- four valid targets recur exactly twice;
- unchanged 12,288 one-shot churn follows;
- all candidate groups are queried only for final evaluation / fixed churn observation, never admission.

Seeds:
- 237000000;
- 238000000.

32 episodes per seed.

## Metrics

Per Q:
- completion/panic rate;
- total evictions;
- evictions by age 0/1/2;
- valid target admission/final accuracy;
- 12-hot accuracy;
- target16 exact accuracy;
- nonpersistent admission/final retention;
- one-shot false admissions;
- max current/history occupancy;
- exact recall entries;
- checkpoint count.

## Interpretation

Successful Q=1..16 with evictions confined to age 2 and exact valid-target retention would establish that stale-evidence reclamation crosses the UP-124C hard boundary at unchanged memory. Loss of valid targets, younger-age eviction, or false admission would expose the cost or insufficiency of this rule.

A negative result is valid and will be sealed without tuning the rule.

## Bounds

No memory increase, no exact-recall increase, no semantic labels, no query-derived priority, no adaptive horizon, no future oracle, no alternate replacement rule after results, no result-informed retry, no live activation, no production authority.
