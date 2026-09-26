# Wingless UP-100B — prefix-timed three-way routing integration

Status: preregistered scientific event-routing integration experiment.

Scientific parent: sealed UP-99B `b54496cb0cd508acfb0a6fac65de96a57ed49a1f`.

## Question

UP-99B showed perfect STORE / OBSERVE / REPORT classification when the readout occurs at the verb boundary. Does that representation actually support bounded downstream routing, or is the gain limited to isolated classification?

## Frozen classifier

Use the exact UP-99B prefix-through-verb arm:
- deterministic 64-D event encoder;
- input only `<name> <verb>`;
- three-class linear softmax head;
- zero initialization;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- exact UP-97B train/held-out split;
- no explicit class label at learned inference.

## Frozen routing arms

1. `explicit_event_class`
   - upper-bound control.

2. `learned_prefix_threeway`
   - infer STORE / OBSERVE / REPORT from the frozen verb-prefix classifier;
   - STORE writes bounded exact recall;
   - OBSERVE does not write exact recall;
   - REPORT queries exact recall.

## Frozen downstream workload

- exact recall cap 16;
- target sets 4, 8, 16;
- total event-write loads 32, 64, 128, 256;
- same deterministic sequence construction as UP-97B;
- two fixed seed bases 155M and 156M.

## Metrics

- held-out classifier accuracy/precision/recall;
- per-verb held-out accuracy;
- downstream query accuracy;
- whole-target-set exact accuracy;
- maximum exact-recall entries.

## Interpretation

If learned prefix routing matches explicit routing, the B97/B98 failure is closed as a readout-timing problem on this domain. Any downstream gap despite perfect isolated classification becomes an integration boundary.

## Bounds

No threshold tuning, no class weighting, no attention, no trainable recurrence, no recall-cap increase, no explicit event class in learned inference, no result-informed retry, no live activation, no production authority.
