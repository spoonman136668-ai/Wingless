# Wingless UP-98B — three-way decision geometry

Status: preregistered scientific semantic-routing experiment.

Scientific parent: sealed UP-97B `aca2c8ab1c465c3978f5bcb5078c663e789f4d02`.

## Question

UP-97B used a class-balanced three-way corpus but still showed a STORE-specific recall collapse, especially on the lexical family `stores`. Is that failure caused primarily by softmax class competition rather than corpus imbalance?

## Frozen corpus and split

Use the exact UP-97B lexical domain and split:
- six names;
- six values;
- STORE: stores / keeps / holds;
- OBSERVE: observes / sees / notes;
- REPORT: reports / recalls / tells;
- train iff `(name_index + 2*value_index + verb_global_index) mod 4 != 3`;
- held-out recombination otherwise.

No resampling and no class-frequency changes are allowed.

## Frozen arms

1. `softmax_threeway`
   - exact UP-97B classifier.

2. `ovr_logistic_threeway`
   - three independent one-vs-rest logistic heads;
   - same deterministic 64-D event encoder;
   - zero initialization;
   - 20 epochs;
   - online SGD;
   - learning rate 0.08;
   - deterministic order;
   - prediction = highest of the three logistic probabilities.

This is a loss/decision-geometry ablation, not an imbalance correction.

## Evaluation

For both arms:
- train and held-out three-way accuracy;
- per-class precision and recall;
- per-verb held-out accuracy;
- STORE/OBSERVE/REPORT confusion counts;
- identical downstream routing probe from UP-97B;
- target set sizes 4, 8, 16;
- write loads 32, 64, 128, 256;
- seeds 151M and 152M;
- exact recall cap 16.

## Interpretation

If one-vs-rest improves STORE recall and downstream exact routing without degrading the other classes, softmax competition is implicated. If the same lexical failure persists, move next to representation/lexical-family separation rather than further loss tuning.

## Bounds

No threshold search, no class resampling, no encoder change, no capacity increase, no attention, no explicit event class at learned inference, no result-informed retry, no live activation, no production authority.
