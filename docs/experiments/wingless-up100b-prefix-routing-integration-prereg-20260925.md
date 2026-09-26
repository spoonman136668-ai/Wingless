# Wingless UP-100B — prefix-timed learned routing integration

Status: preregistered scientific routing-integration experiment.

Scientific parent: sealed UP-99B `b54496cb0cd508acfb0a6fac65de96a57ed49a1f`.

## Question

UP-99B showed perfect three-way STORE / OBSERVE / REPORT separation when the event classifier reads only the prefix through the verb. Does that representation repair the downstream bounded-routing failure observed in UP-97B?

## Frozen classifier

Use the exact UP-99B `prefix_through_verb` arm:
- deterministic 64-D event encoder;
- input `<name> <verb>` only;
- three-class linear softmax head;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- exact UP-97B train split and lexical domain;
- no explicit event-class bit.

## Frozen routing workload

Reuse UP-97B semantics:
- STORE writes/overwrites a binding;
- OBSERVE has no exact-memory write side effect;
- REPORT queries a binding;
- exact recall capacity 16;
- target set sizes 4, 8, 16;
- total event counts 32, 64, 128, 256;
- deterministic two-seed qualification.

Arms:
1. `explicit_event_class` control.
2. `learned_prefix_threeway`.

For the learned arm, event class is predicted from `<name> <verb>` before the value suffix is considered.

## Metrics

For all 24 routing cells:
- final query accuracy;
- whole target-set exact accuracy;
- maximum recall entries used.

## Interpretation

If learned-prefix routing matches the explicit-event control, the UP-97B routing deficit is attributable to readout timing rather than a fundamental inability to route natural lexical events. Any remaining deficit is sealed without threshold tuning.

## Bounds

No explicit event class in learned inference, no value suffix in the classifier input, no attention, no recurrent training, no exact-recall capacity increase, no threshold search, no result-informed retry, no live activation, no production authority.
