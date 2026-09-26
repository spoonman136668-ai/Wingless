# Wingless UP-163B — cleanup path curvature

Status: preregistered scientific lexical-stability mechanism diagnostic.

Scientific parent: sealed UP-162B a705c3abb67a77c4575f43e2f872424fe881c899.

## Question

UP-162B showed that frozen STORE-vs-OBSERVE old-memory alignment predicts the sealed cleanup-order sign for five of six subjects, with pax as a counterexample. At what point along the actual chained cleanup paths do the subject-specific retention differences emerge?

## Frozen state

For each accepted subject:
- reconstruct the exact UP-160B frozen post-REPORT state at terminal epoch 16;
- use the same matched two-subject corpus;
- use the same 20 non-tail new-family updates;
- use the same three pre-REPORT old rehearsals;
- use the same four-example target-subject REPORT interference block.

## Frozen cleanup paths

Trace two fixed 12-update cleanup paths from the identical post-REPORT gate:

A. STORE → OBSERVE → REPORT
- STORE indices 0..4;
- OBSERVE indices 5..9;
- REPORT indices 13..14.

B. OBSERVE → STORE → REPORT
- OBSERVE indices 5..9;
- STORE indices 0..4;
- REPORT indices 13..14.

## Checkpoints

For each path record:
- checkpoint 0: frozen post-REPORT state;
- checkpoint 1: after first cleanup block;
- checkpoint 2: after second cleanup block;
- checkpoint 3: after final REPORT cleanup block.

At every checkpoint:
- old retention;
- displacement norm from checkpoint 0;
- cosine of displacement with the frozen old-reference vector.

At matched checkpoints 1–3:
- gate distance between A and B;
- old-retention difference A-B.

## Interpretation

The exact six-subject path table is the result. If pax begins with the same first-order tendency as the five aligned subjects but reverses after checkpoint 2 or checkpoint 3, later nonlinear composition explains the UP-162B exception. If the sign is already opposite at checkpoint 1, the simple aggregate geometry was insufficient even locally.

No threshold, classifier, feature selection, or adaptive schedule is allowed.

## Bounds

Diagnostic only. No retained diagnostic updates, no adaptive ordering, no extra training epochs, no memory change, no architecture change, no result-informed retry, no live activation, no production authority.
