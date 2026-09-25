# Wingless UP-95B — learned local event classifier for bounded-memory admission

Status: preregistered scientific memory-routing experiment.

Scientific parent: sealed UP-94B `83fc5d7dac7b4e6bebbec8c93aee0e6d73a111b6`.

## Question

UP-94B showed that explicit local event structure is sufficient for perfect bounded-memory admission, while generic change/confidence signals are not. Can a small learned classifier infer store-vs-observation event type from local byte context and preserve the same routing behavior without receiving an explicit event-type bit at inference time?

## Frozen local event language

Keys: `k00` through `k31`.

Values: `v00` through `v31`.

Surface forms:

Store event:
`<key> stores <value>.`

Observation event:
`<key> observes <value>.`

The classifier receives only the recurrent encoding of the event bytes.

## Frozen event encoder

- 64 float64 recurrent values;
- same deterministic signed-permutation transport and deterministic bipolar byte embedding used by the UP-LM line;
- recurrent state reset at each event boundary;
- no trainable recurrent parameters.

## Frozen classifier

- binary linear logistic head over the final 64-D event state;
- zero-initialized weights and bias;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- deterministic lexical order;
- decision threshold exactly 0.5;
- no validation tuning.

Training examples:
- keys 0..23;
- all values 0..31;
- both event types.

Held-out event classification:
- keys 24..31;
- all values 0..31;
- both event types.

## Frozen downstream routing task

Same sparse event stream structure as UP-94B:
- target keys: 4, 8, 16;
- total non-query writes: 32, 64, 128, 256;
- target keys receive store events near the beginning and later store rewrites;
- distractors are observation events to unique keys;
- final queries target all target keys;
- value vocabulary 32;
- 64 episodes per setting;
- seeds 139M and 140M.

## Frozen routing arms

1. `fifo_all_writes`
2. `explicit_store_type` — upper-bound control using the known event type.
3. `learned_local_classifier` — admit only events classified as store by the frozen 0.5 classifier.

All writes still update the same recurrent carrier.

## Metrics

Classifier:
- train accuracy;
- held-out-key accuracy;
- precision;
- recall.

Routing:
- final target query accuracy;
- whole-target-set exact accuracy;
- recall entries used;
- admission precision/recall;
- false-positive admissions.

## Bounds

No explicit event-type bit in the learned-classifier routing arm, no future-query oracle, no attention, no recurrent-state training, no classifier threshold search, no capacity increase beyond 16 entries, no result-informed retry, no production authority, no live activation.
