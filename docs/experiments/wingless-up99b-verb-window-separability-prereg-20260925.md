# Wingless UP-99B — verb-window event separability

Status: preregistered scientific representation-locality experiment.

Scientific parent: sealed UP-98B `dde36e56943a4e7af2c39d0aa33c7f6b9a9bb23b`.

## Question

UP-98B showed that changing the decision head from three-way softmax to one-vs-rest does not repair STORE recognition. Does the failure arise because the frozen recurrent event representation loses verb/class information after the value suffix is processed?

## Frozen data and training

Use the exact UP-97B:
- six names;
- six values;
- nine verbs across STORE / OBSERVE / REPORT;
- train/held-out recombination split;
- deterministic 64-D event encoder;
- three-class linear softmax head;
- zero initialization;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- deterministic order.

## Frozen arms

1. `full_clause`
   - exact UP-97B representation of `<name> <verb> <value>.`.

2. `prefix_through_verb`
   - encode only `<name> <verb>`;
   - the value and punctuation are never presented to the event classifier;
   - no explicit event-class label is supplied.

Both arms use the same training examples, split, head, epochs, and learning rate.

## Evaluation

For both arms:
- train and held-out overall three-way accuracy;
- held-out per-class precision/recall;
- held-out per-verb accuracy;
- held-out confusion matrix.

Primary diagnostic: whether STORE recall, especially `stores` and `keeps`, recovers when the value suffix is excluded.

## Interpretation

If prefix-through-verb materially improves class separation, treat readout timing / suffix interference as the primary event-class representation boundary. If not, move upstream to the recurrent encoder itself.

## Bounds

No threshold tuning, no class weighting, no attention, no trainable recurrence, no capacity expansion, no result-informed retry, no live activation, no production authority.
