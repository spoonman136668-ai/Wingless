# Wingless UP-120B — stores competitor-span orthogonalization

Status: preregistered scientific representation intervention.

Scientific parent: sealed UP-119B fc178cbf08c5ab6cc223e11a723cf8c683ce3ec0.

## Question

UP-119B showed that the one-axis sibling-STORE centroid shift reduces OBSERVE collapse but exposes REPORT-directed failures. Can a fixed stores-only representation target that removes both OBSERVE- and REPORT-family components eliminate the three-class lexical collapse without changing training or replay?

## Frozen training matrix

Repeat the exact UP-118B/UP-119B matrix:
- all six semantic-class acquisition orders;
- replay policies:
  - current_class_excluded
  - stage3_anchor2
- same subjects, verbs, 20 epochs per acquisition batch, learning rate 0.08, grounding count, and fixed replay budget 6.

## Frozen representation arms

1. baseline
   - exact existing deterministic encoder.

2. stores_centroid_shift
   - exact sealed UP-118B shift to the sibling STORE centroid.

3. stores_competitor_orthogonal
   - compute baseline subject-averaged sibling STORE centroid S from:
     keeps, holds, saves;
   - compute baseline OBSERVE-family centroid O from:
     observes, sees, notes, watches;
   - compute baseline REPORT-family centroid R from:
     reports, recalls, tells, remembers;
   - project S onto span{O,R} using the exact 2x2 Gram-system solution;
   - define residual T = S - projection_span(O,R)(S);
   - scale T to preserve ||S||;
   - define one fixed stores-only offset:
     delta = T_scaled - mean_baseline(stores);
   - whenever the lexical surface is exactly stores, use h' = h + delta;
   - all other lexical surfaces are unchanged;
   - the same offset is used during base training, replay, and evaluation.

No coefficient or norm is tuned from outcomes.

## Evaluation

For every representation arm x class order x replay policy x stage:
- stores accuracy;
- keeps accuracy;
- base-lexicon aggregate accuracy;
- acquired-verb aggregate accuracy;
- STORE-minus-OBSERVE mean/min margin;
- STORE-minus-REPORT mean/min margin;
- wrong OBSERVE and wrong REPORT counts.

Also report target-vector cosine to STORE, OBSERVE, and REPORT centroids.

## Interpretation

If the competitor-orthogonal target removes the order-dependent stores failures while leaving acquired and sibling accuracy intact, the collapse is causally attributable to low-dimensional lexical aliasing against the two competing class directions. Failure would indicate that simple mean geometry is still insufficient.

## Bounds

Diagnostic stores-only representation intervention. No replay-budget change, no learned offset, no adaptive geometry, no extra hidden state, no threshold tuning, no result-informed retry, no live activation, no production authority.
