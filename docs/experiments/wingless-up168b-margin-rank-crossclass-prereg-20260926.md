# Wingless UP-168B — cross-class margin-rank generalization

Status: preregistered scientific lexical-stability diagnostic.

Scientific parent: sealed UP-167B 246eb162750f0e72122936502ee00d99b8a49c8a.

## Question

UP-167B showed that all crossings caused by five independent REPORT probes originated in the four smallest absolute pre-step margins out of 120 old examples. Is the same vulnerability-rank signal predictive for independent STORE and OBSERVE updates, or is it REPORT-specific?

## Frozen path states

For each accepted new-family subject:
- reconstruct the exact post-interference state;
- build path A = STORE→OBSERVE;
- build path B = OBSERVE→STORE.

## Independent cross-class probes

Probe old rehearsal surfaces independently from fresh clones:
- STORE indices 0..4;
- OBSERVE indices 5..9.

Each probe is one update only. Probes are not chained.

## Ranking

Before each probe:
- evaluate all 120 old examples;
- rank by absolute probability margin ascending;
- apply the frozen one-step surface;
- identify examples crossing margin zero;
- record rank and rank fraction.

No fitted threshold, classifier, or adaptive selection is used.

## Measurements

Per subject × path × probe:
- probe index and class;
- crossing count;
- mean crossing absolute margin;
- mean non-crossing absolute margin;
- mean/max crossing rank;
- mean/max crossing rank fraction.

Per crossing:
- subject/path/probe;
- split/name/verb/class;
- signed pre/post margin;
- vulnerability rank;
- crossing direction.

## Interpretation

If STORE/OBSERVE crossings also concentrate at the lowest pre-step ranks, absolute margin rank is a cross-class state-vulnerability signal rather than a REPORT-specific effect. High-rank crossings define its generalization boundary.

## Bounds

Diagnostic only. No retained diagnostic updates, no adaptive ordering, no extra training, no memory or architecture change, no post-result threshold/classifier, no result-informed retry, no live activation, no production authority.
