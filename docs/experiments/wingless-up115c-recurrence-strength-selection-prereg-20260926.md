# Wingless UP-115C — local recurrence-strength selection

Status: preregistered scientific bounded-memory selection experiment.

Scientific parent: sealed UP-114C 775d35badbbd0a320f85739598575f1655fecb2b.

## Question

UP-114C showed that exact probation admits all legitimate recurring candidates, so when demand exceeds the fixed 16-entry cap, newly recurring items displace established hot residents. Can a purely local recurrence-strength signal select stronger candidates without semantic priority labels?

## Frozen exact memory

- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- generation interval 32;
- phase-free slot-reuse checkpoint;
- no query-derived admission evidence;
- no future oracle;
- no priority labels.

## Equal-memory probation arms

Admission sidecar remains exactly 32 x uint32 = 128 bytes.

1. exact2_control
   - exact UP-113C behavior;
   - admit on the second sighting within a generation.

2. exact3_strength
   - encode a two-bit occurrence count inside each uint32 sidecar entry;
   - admit only on the third sighting within the same generation;
   - no extra bytes.

All experimental key IDs fit below 2^30, leaving the top two bits available for the bounded count.

## Candidate population

Eight legitimate recurring candidates:
- four weak candidates receive exactly 2 sightings within 31 filtered writes;
- four strong candidates receive exactly 3 sightings within 31 filtered writes.

All candidates use the same value vocabulary and no semantic labels.

Target preference implied only by local recurrence strength:
- 12 established hot + 4 strong recurring = 16 desired retained items.

## Churn

After candidate presentation:
- 49,152 unique one-shot churn writes;
- query all established hot and all candidate keys every four churn writes;
- 32 episodes per seed;
- seeds 217M and 218M.

## Metrics

Per arm:
- weak candidate admission/retention;
- strong candidate admission/retention;
- 12-hot accuracy;
- target-16 accuracy and exactness;
- total candidate accuracy;
- one-shot false admissions;
- exact-memory occupancy;
- metadata bytes.

## Interpretation

If exact3 preserves the 12 hot residents while admitting the four stronger recurring candidates, a bounded local confidence mechanism can resolve the over-capacity selection problem without semantic priority labels.

## Bounds

No extra admission memory, no capacity increase, no query-derived admission, no semantic priority, no adaptive threshold, no phase labels, no future oracle, no result-informed retry, no live activation, no production authority.
