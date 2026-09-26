# Wingless UP-122C — age2 history occupancy pressure

Status: preregistered scientific bounded-memory sidecar-capacity experiment.

Scientific parent: sealed UP-121C 33516f1cc1a939e8859c74d1cef17c41fa6cf9a1.

## Question

UP-121C showed that heterogeneous adjacent/skip-one/skip-two persistence shares the 128-byte age2 history correctly at modest history occupancy. Does the same exact mechanism preserve valid targets when many qualified-but-nonpersistent candidates occupy the history table across consecutive generations?

## Frozen mechanism

Use exact UP-120C/UP-121C age2_extended:
- current table: 32 × uint16 = 64 bytes;
- history table: 32 × uint16 = 64 bytes;
- total admission sidecar = 128 bytes;
- 14-bit exact candidate keys;
- current-generation qualification after at least two sightings;
- history ages 0, 1, and 2 valid;
- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- generation interval 32;
- phase-free slot-reuse checkpoint;
- no semantic priority labels;
- no query-derived admission evidence;
- no future oracle.

## Valid targets

Four skip-one targets:
- each receives 2 sightings in generation 1;
- absent in generation 2;
- receives 2 sightings in generation 3.

Desired exact-memory target:
- 12 hot + 4 valid targets = 16.

## Qualified nonpersistent pressure

Pressure N ∈ {0, 4, 8, 12}.

Generation 1:
- N ephemeral set-A keys each receive exactly 2 sightings;
- they never recur.

Generation 2:
- N distinct ephemeral set-B keys each receive exactly 2 sightings;
- they never recur.

Generation 3:
- only the four valid targets recur.

Every generation is completed to exactly 32 filtered writes with unique one-shot fillers.

At N=12, the history immediately after generation 2 can contain 4 valid target records + 12 aged set-A records + 12 fresh set-B records = 28 occupied entries out of 32, without changing the memory budget.

## Churn

After generation 3:
- 12,288 unique one-shot writes;
- query hot, valid targets, and all ephemeral keys every four churn writes;
- 32 episodes per seed;
- seeds 231M and 232M.

## Metrics

Per pressure:
- valid-target admission and final accuracy;
- hot accuracy;
- target16 accuracy/exactness;
- ephemeral admission and final retention rate;
- one-shot false admissions;
- maximum observed current-table occupancy;
- maximum observed history-table occupancy;
- exact-memory occupancy;
- checkpoint count.

## Interpretation

If target exactness survives through near-full history occupancy with zero ephemeral admission, the fixed exact sidecar has substantial pressure headroom. Degradation before table saturation would reveal competition/aging interactions rather than raw entry capacity.

## Bounds

No admission-memory increase, no exact-memory capacity increase, no adaptive horizon, no semantic labels, no query-derived admission, no future oracle, no result-informed retry, no live activation, no production authority.
