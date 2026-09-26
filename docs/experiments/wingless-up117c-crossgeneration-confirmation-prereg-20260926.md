# Wingless UP-117C — cross-generation recurrence confirmation

Status: preregistered scientific bounded-memory persistence experiment.

Scientific parent: sealed UP-116C f8669aaa748a610c85847512cb96d6af850b91a1.

## Question

UP-116C showed that within-generation recurrence strength cannot distinguish important three-hit candidates from even stronger four-hit burst distractors. Can persistence across two separate generations provide a local signal that rejects one-generation bursts while admitting genuinely recurring candidates?

## Frozen exact memory

- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- no semantic priority labels;
- no query-derived admission evidence;
- no future oracle;
- no explicit phase labels.

## Equal admission-memory arms

Both arms use exactly 128 bytes of admission sidecar memory.

1. exact3_control
   - exact UP-115C threshold-3 mechanism;
   - 32-entry uint32 table;
   - admit on third sighting within one 32-filtered-write generation.

2. crossgen_2x2
   - two 32-entry uint16 exact tables:
     - current generation: 64 bytes;
     - previous qualifying generation: 64 bytes.
   - all experimental keys remain below 2^15.
   - current entry uses one high bit to record whether the key has been seen at least twice this generation.
   - at a 32-filtered-write boundary, only keys seen at least twice become the previous-generation qualifying set; current table clears.
   - admit a nonresident key only when:
     a) it qualified in the immediately previous generation, and
     b) it is seen at least twice again in the current generation.
   - thus one-generation bursts cannot admit, regardless of receiving 3 or 4 sightings.

## Candidate workload

Generation 1:
- four target persistent candidates each receive exactly 2 sightings;
- four burst distractors each receive exactly 4 sightings;
- deterministic ascending key order;
- one-shot fillers complete the generation.

Generation 2:
- the four target candidates each receive exactly 2 sightings;
- burst distractors do not recur;
- one-shot fillers complete the generation.

Target working set:
- 12 established hot + 4 persistent candidates = 16.

## Churn

After generation 2:
- 24,576 unique one-shot churn writes;
- all synthetic keys remain below 2^15;
- query hot, target, and distractor keys every four churn writes;
- 32 episodes per seed;
- seeds 221M and 222M.

## Metrics

Per arm:
- target admission and retention;
- burst-distractor admission and retention;
- hot accuracy;
- target16 accuracy and exactness;
- one-shot false admissions;
- exact-memory occupancy;
- checkpoint count;
- admission metadata bytes.

## Interpretation

If crossgen_2x2 preserves the 12 hot + 4 persistent targets while rejecting four-hit one-generation bursts, cross-generation persistence supplies the missing local confidence dimension without semantic labels or more memory.

## Bounds

No admission-memory growth, no exact-memory capacity increase, no semantic labels, no query-derived admission, no adaptive threshold, no result-informed retry, no live activation, no production authority.
