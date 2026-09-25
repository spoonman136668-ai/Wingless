# Wingless UP-98B — decomposed three-way classifier

Status: preregistered scientific event-routing experiment.

Scientific parent: sealed UP-97B `aca2c8ab1c465c3978f5bcb5078c663e789f4d02`.

## Question

UP-97B used a class-balanced corpus, yet the three-class softmax head showed a STORE recall failure and poor downstream routing. Is that failure caused by multiclass competition in the decision head rather than lack of separability in the frozen event representation?

## Frozen data

Use the exact UP-97B names, values, verbs, surfaces, structural train/held-out split, deterministic 64-D encoder, lexical order, 20 epochs, and learning rate 0.08.

No examples are added or removed and no class weighting is applied.

## Frozen arms

1. `softmax_baseline`
   - exact UP-97B three-class linear softmax head.

2. `one_vs_rest`
   - three independent linear logistic heads;
   - each head trains against its own binary target on every training example;
   - prediction is the class with maximum sigmoid probability;
   - no per-class threshold tuning.

Both arms start from zero weights and use identical training order, epochs, and learning rate.

## Evaluation

For both arms:
- train and held-out overall accuracy;
- held-out precision/recall for STORE, OBSERVE, REPORT;
- held-out per-verb accuracy;
- held-out 3x3 confusion matrix.

Primary diagnostic: whether STORE recall and the `stores` / `keeps` verb failures improve without degrading the other classes.

## Interpretation

If one-vs-rest materially improves STORE recall on the same frozen representation and data, classify UP-97B primarily as decision-head competition. If it does not, classify the boundary as representation/lexical separability and move the next experiment upstream rather than tuning thresholds.

## Bounds

No threshold search, no class weights, no encoder change, no recurrent training, no result-informed retry, no attention, no capacity expansion, no live activation, no production authority.
