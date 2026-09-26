# Wingless UP-118B — stores centroid-shift causal diagnostic

Status: preregistered scientific representation intervention.

Scientific parent: sealed UP-117B fa5add5fe53f609913041f6a12b50c545dabc2cd.

## Question

UP-117B showed that the frozen encoding of `stores` is much closer to the OBSERVE-family centroid than to its own STORE-family centroid, while `keeps` is not. Is that anomalous lexical geometry causally responsible for the order-dependent `stores` collapse?

## Frozen training matrix

Repeat the exact UP-117B / UP-114B scientific matrix:
- all six semantic-class acquisition orders;
- policies:
  - current_class_excluded
  - stage3_anchor2
- same six subjects;
- same base/acquired verbs;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- four grounding subjects;
- fixed replay budget 6;
- no threshold changes.

## Frozen representation arms

1. baseline
   - exact existing deterministic 64-D encoder.

2. stores_centroid_shift
   - all verb encodings except `stores` are bit-for-bit identical to baseline;
   - compute the subject-averaged baseline vector for `stores`;
   - compute the subject-averaged centroid of the other STORE verbs:
     - keeps
     - holds
     - saves
   - define one fixed 64-D offset:
     `delta = STORE_sibling_centroid - stores_mean`;
   - whenever the lexical surface is exactly `stores`, use:
     `h_shifted = h_baseline + delta`;
   - apply that same shifted encoding during base training, replay, and evaluation;
   - no normalization, scaling, clipping, or learned correction.

Thus the intervention preserves each subject-specific deviation around the `stores` mean while moving only that lexical surface's mean to the sibling STORE centroid.

## Evaluation

For every representation arm x class order x replay policy x stage:
- `stores` accuracy;
- `keeps` accuracy;
- base-lexicon aggregate accuracy;
- acquired-verb aggregate accuracy;
- mean and minimum STORE-minus-OBSERVE margin for `stores`.

Also report the shifted static centroid cosines for `stores`.

## Interpretation

If the centroid shift removes the `stores` failures across class orders without damaging acquired or sibling lexical accuracy, the UP-117B geometry is causally implicated. Persistent failures would show that learned trajectory/interference remains necessary even after correcting the static mean geometry.

## Bounds

Diagnostic lexical intervention only. No new hidden state, no replay-budget change, no class weighting, no adaptive geometry, no learned offset, no threshold tuning, no result-informed retry, no live activation, no production authority.
