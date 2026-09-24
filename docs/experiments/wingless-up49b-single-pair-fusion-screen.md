# Wingless UP-49B — frozen step-42 single-pair fusion screen

Status: preregistered scientific mechanism experiment.

Scientific parent: UP-48B seal `3f9789035ac9324364ae4718bf900b07189cafce`.

## Question

UP-48B showed that exact equality itself is useful and that fusing step 42's nearest pair [4,5] sharply improves hard recurrent capability.

UP-49B asks whether the frozen training-side smooth objective can identify the best single exact pair fusion when every one of the 15 possible pairs is treated equally.

## Frozen design

Starting only from frozen UP-47 step 42:

- enumerate all 15 unordered channel pairs;
- for each pair, replace the two offsets by their mean and change nothing else;
- evaluate the unchanged training-side smooth objective for all candidates;
- select the maximum smooth objective with the existing lower-soft-capacity exact-tie rule;
- only after selection is frozen, evaluate true held-out hard capability.

For diagnostic correlation only, all 15 hard results are also reported after selection. Those hard results do not influence selection.

No optimizer, lambda, objective component, threshold, seed, or held-out split changes.

## Interpretation

If the smooth-selected pair is also the post-hoc best hard pair, the surrogate is locally aligned once exact fusion is part of the candidate space.

If a different pair is post-hoc best, the remaining bottleneck is surrogate ranking rather than fusion reachability.

If any single-pair candidate passes the frozen hard gates, the next mechanism experiment should test persistence/generalization of that structure before introducing multi-pair fusion.

If none passes, the next step is a bounded two-pair screen seeded only by this preregistered result.

## Plain speak

We now know one equality helps. Instead of guessing which equality to try next, this experiment gives every possible pair the same chance.

The score chooses one without looking at the real test set, then we check whether that choice was actually the best for memory and relations.
