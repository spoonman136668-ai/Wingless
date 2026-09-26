# Wingless UP-LM0L — interleaved clause-order routing

Status: preregistered scientific language/memory structural-generalization experiment.

Scientific parent: sealed UP-LM0K 73584e206a72ae8d3b3c7ebd2c24596f92ee6662.

## Question

UP-LM0K remained exact through stream16 on block-ordered paragraphs. Does the learned prefix-timed event router remain correct when STORE, OBSERVE, update, and REPORT clauses are interleaved rather than grouped by event type?

## Frozen models

Use the exact UP-LM0K language model and UP-LM0J three-way event classifier:
- language model is trained only on the original block-ordered training corpus;
- classifier is unchanged;
- 64-D recurrent state;
- 512 recurrent-state bytes;
- exact recall cap 16;
- no attention;
- no future oracle.

## Interleaved held-out corpus

Use the same six names, six values, structural held-out split, and update-count schedule as UP-LM0F.

For each of four active names, emit:
1. initial STORE;
2. one OBSERVE distractor;
3. that name's update STORE if its update index is active;
4. REPORT for that name.

Then proceed to the next active name.

Thus queries occur throughout the paragraph instead of only after all stores/observes/updates.

No new lexical items are introduced.

## Frozen arms

1. explicit_event_routing
2. learned_prefix_threeway

Evaluate each at:
- stream1
- stream4

## Metrics

- top-1 byte accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- exact accuracy by update count;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- maximum recall entries.

Classifier metrics are reported unchanged.

## Interpretation

Matching explicit routing on the interleaved corpus would demonstrate structural event-order robustness of the learned router. Failure shared by both arms would identify generic language-order shift; learned-only failure would identify routing integration.

## Bounds

No retraining on the interleaved corpus, no state expansion, no memory-cap increase, no attention, no threshold tuning, no result-informed retry, no live activation, no production authority.
