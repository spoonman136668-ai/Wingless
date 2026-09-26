# Wingless UP-167B — margin-rank REPORT generalization

Status: preregistered scientific lexical-stability diagnostic.

Scientific parent: sealed UP-166B 46d8e1d00cf6cde7c0eb8cee6f5fb610466ae902.

## Question

UP-166B found that all observed crossings under REPORT items 13 and 14 came from the four smallest absolute pre-step margins out of 120 old examples. Does that vulnerability-rank signal generalize across all five frozen old REPORT surfaces?

## Frozen path states

For each accepted new-family subject:
- reconstruct the exact frozen post-interference state;
- build path A = STORE→OBSERVE;
- build path B = OBSERVE→STORE.

## Independent REPORT probes

Probe old rehearsal REPORT indices:
- 10;
- 11;
- 12;
- 13;
- 14.

Each REPORT surface is applied independently from a fresh clone of the same path state. REPORT probes are not chained.

## Ranking

Reuse UP-166B unchanged:
- evaluate all 120 old examples;
- rank by absolute pre-step probability margin ascending;
- apply one REPORT probe;
- identify margin-zero crossings;
- record crossing vulnerability ranks and rank fractions.

No threshold or classifier is fitted.

## Measurements

Per subject × path × REPORT surface:
- crossing count;
- mean crossing absolute margin;
- mean non-crossing absolute margin;
- mean/max crossing rank;
- mean crossing rank fraction.

Per crossing:
- subject/path/report;
- split/name/verb/class;
- signed pre/post margin;
- vulnerability rank and rank fraction;
- crossing direction.

## Interpretation

Low-rank concentration across REPORT identities supports pre-margin rank as an interference-general vulnerability signal. REPORT identities producing high-rank crossings define its boundary.

## Bounds

Diagnostic only. No retained diagnostic updates, no adaptive ordering, no extra training, no memory/architecture change, no post-result threshold, no result-informed retry, no live activation, no production authority.
