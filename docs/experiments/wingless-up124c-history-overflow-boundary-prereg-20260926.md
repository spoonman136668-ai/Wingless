# Wingless UP-124C — exact history overflow boundary

Status: preregistered scientific bounded-memory overflow experiment.

Scientific parent: sealed UP-123C dfa1f6e4176b2a35213b91c97b82cfa8ae7af6f4.

## Question

UP-123C established exact behavior at 32/32 history occupancy. What does the unchanged mechanism do when at least one additional qualified, nonpersistent candidate must be inserted while all 32 history slots remain occupied?

## Frozen mechanism

Use the exact UP-120C/UP-123C age2_extended mechanism:
- current table: 32 × uint16;
- history table: 32 × uint16;
- admission sidecar: 128 bytes;
- exact recall cap: 16;
- 14-bit candidate keys;
- history ages 0/1/2;
- 12 established hot residents;
- two-bit-aging replacement;
- generation interval 32;
- no semantic priority;
- no query-derived admission;
- no future oracle.

The existing history insertion behavior is not modified. If the mechanism raises its existing full-history panic, the probe records that event as the scientific outcome for the episode and does not repair or replace the table.

## Saturation construction

Reproduce the exact UP-123C P=12 path:

Generation 1:
- 16 qualified nonpersistent set-A candidates, exactly 2 sightings each.

Generation 2:
- four valid recurring targets, exactly 2 sightings each;
- 12 qualified nonpersistent set-B candidates, exactly 2 sightings each.

At the end of generation 2 the history table is exactly 32/32 occupied.

## Overflow generation

Arms use overflow qualified-candidate counts Q ∈ {0, 1, 4, 8}.

For each Q:
- Q fresh nonpersistent candidates receive exactly 2 sightings each;
- remaining writes in the generation are unique one-shot candidates;
- the generation always reaches the exact 32-write boundary.

Q=0 is the saturated control.
Q>0 tests the first and larger over-capacity insertion pressures.

If the generation completes without overflow:
- the four valid targets recur with exactly 2 sightings each;
- unchanged 12,288-write churn follows;
- hot and valid target accuracy are measured.

## Seeds and episodes

- seeds: 235000000 and 236000000;
- 32 episodes per seed;
- 64 episodes per arm.

## Metrics

Per Q:
- overflow episode rate;
- exact overflow message;
- maximum history/current occupancy observed before overflow;
- valid-target evaluated episode count;
- valid admission/final accuracy on evaluable episodes;
- hot accuracy;
- target16 exact accuracy;
- nonpersistent admission/retention;
- one-shot false admissions;
- exact recall entries used;
- checkpoint count.

## Interpretation

Overflow at Q=1 establishes a hard representable boundary immediately beyond 32 history entries. Survival at Q=1 but failure at larger Q identifies bounded slack in the current insertion/aging schedule. Survival across all Q would show that aging frees enough space before insertion despite nominal saturation.

A capacity panic is a valid scientific negative and must be sealed exactly as observed.

## Bounds

No history replacement rule, no memory increase, no exact-recall increase, no adaptive horizon, no semantic labels, no query-derived evidence, no future oracle, no panic suppression inside the mechanism, no result-informed retry, no live activation, no production authority.
