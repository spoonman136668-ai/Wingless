# Wingless UP-123C — exact age2 history saturation

Status: preregistered scientific bounded-memory sidecar-capacity experiment.

Scientific parent: sealed UP-122C 6cac73d570c9f795bcb7879cc6d57d1f4f8211dd.

## Question

UP-122C remained exact with 28 of 32 history slots occupied. Does the unchanged 128-byte exact history still admit valid recurring targets when the history table reaches exactly 32/32 occupied entries?

## Frozen mechanism

Use exact UP-120C/UP-122C age2_extended:
- current table: 32 × uint16 = 64 bytes;
- history table: 32 × uint16 = 64 bytes;
- total admission sidecar = 128 bytes;
- 14-bit exact candidate keys;
- history ages 0/1/2 valid;
- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- generation interval 32;
- phase-free slot-reuse checkpoint;
- no semantic priority labels;
- no query-derived admission evidence;
- no future oracle.

## Pressure construction

Four valid recurring targets.

Generation 1:
- 16 qualified nonpersistent set-A keys;
- each receives exactly 2 sightings;
- 32 filtered writes total;
- no valid target appears yet.

Generation 2:
- four valid target keys each receive exactly 2 sightings;
- set-B qualified nonpersistent keys each receive exactly 2 sightings;
- set-B count P ∈ {4, 8, 12};
- fill with unique one-shot writes when needed.

Immediately after generation 2:
- 16 aged set-A records remain in history;
- 4 valid target records are fresh;
- P set-B records are fresh.

Expected history occupancy:
- P=4  => 24 entries;
- P=8  => 28 entries;
- P=12 => 32 entries.

Generation 3:
- the four valid targets recur with exactly 2 sightings each and should be admitted;
- no set-A or set-B key recurs.

## Churn

After generation 3:
- 12,288 unique one-shot writes;
- query hot, valid targets, and all nonpersistent keys every four churn writes;
- 32 episodes per seed;
- seeds 233M and 234M.

## Metrics

Per pressure:
- valid target admission/final accuracy;
- 12-hot accuracy;
- target16 accuracy/exactness;
- nonpersistent admission/final retention;
- one-shot false admissions;
- maximum current-table occupancy;
- maximum history-table occupancy;
- whether the history table reached exactly 32 entries;
- exact-memory occupancy;
- checkpoint count.

## Interpretation

Exact target admission at 32/32 history occupancy would establish that the present sidecar remains correct through its representable saturation point. Any failure at exactly full occupancy would expose an insertion/aging boundary before overflow.

## Bounds

No admission-memory increase, no exact-memory capacity increase, no replacement rule for history entries, no adaptive horizon, no semantic labels, no query-derived admission, no future oracle, no result-informed retry, no live activation, no production authority.
