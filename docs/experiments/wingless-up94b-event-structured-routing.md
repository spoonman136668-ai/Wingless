# Wingless UP-94B — event-structured local bounded-memory routing

Status: preregistered scientific memory-routing experiment.

Scientific parents:
- sealed UP-93B `b667821eb97e31fc69db72d6113135aeebeaa918`;
- sealed UP-LM0C `6323443c61da2841f9ad6889d13abacd23590e61`.

## Question

UP-93B showed that generic change and confidence signals are insufficient for reliable bounded-memory admission. UP-LM0C showed that local event structure can route memory perfectly in one language dependency. How much structure is required? Can a small local finite-state event detector route memory across long distractor streams without an external salience bit or future-query oracle?

## Frozen event stream

Events are byte-like triples:

- `S,key,value` = structured store event;
- `O,key,value` = ordinary observation/distractor event;
- `Q,key,0` = query event.

The router sees the current event type and its fields as they arrive. It does not receive a target/salience flag.

For each episode:
- target keys: 4, 8, or 16;
- each target receives one `S` event near the beginning and one later `S` rewrite to a different value;
- distractors are unique `O` events;
- final `Q` events query all target keys;
- total non-query writes: 32, 64, 128, or 256;
- value vocabulary: 32;
- 64 episodes per setting;
- seeds: 135M and 136M.

## Frozen arms

1. `fifo_all_writes`
   - admit S and O writes.

2. `structured_store_only`
   - admit only S events.

3. `structure_plus_change`
   - admit S only when the key is locally decoded as present and the incoming value differs from the decoded old value;
   - first S for a key is therefore not admitted;
   - this arm tests whether change semantics add anything once event structure is available.

All arms update the same recurrent carrier on every S/O write.

## Metrics

- final target query accuracy;
- whole-target-set exact accuracy;
- recall entries used;
- admission precision relative to S events;
- admission recall relative to S events;
- false-positive admissions.

## Interpretation

This is a pivot away from generic salience threshold ladders. If structured-store-only succeeds, downstream Wingless should treat event recognition—not generic confidence—as the memory-admission problem.

## Bounds

No target/salience bit, no future-query oracle, no capacity increase beyond 16 entries, no attention, no threshold search, no result-informed retry, no production authority, no live activation.
