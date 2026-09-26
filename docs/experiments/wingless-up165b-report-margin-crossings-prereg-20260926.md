# Wingless UP-165B — REPORT margin crossings

Status: preregistered scientific lexical-stability mechanism diagnostic.

Scientific parent: sealed UP-164B report-path sensitivity.

## Question

UP-164B showed that REPORT item 14 flips the STORE↔OBSERVE retention winner for Noah and Opal, while REPORT item 13 creates Quin's negative preference. Are these coarse retention changes caused by specific old-surface examples crossing the gate's own class decision boundary?

## Frozen state and paths

For each of the six accepted subjects:
- reconstruct the exact frozen post-REPORT state used by UP-164B;
- construct path A = STORE then OBSERVE cleanup blocks;
- construct path B = OBSERVE then STORE cleanup blocks;
- use the same REPORT cleanup items 13 and 14 in the same order.

## Frozen old evaluation set

Evaluate every old rehearsal surface on:
- held-out old names `up121bOriginalNames[4:6]`;
- unseen old names `up121bUnseenNames`.

For each example, define probability margin:

`p(correct class) - max(p(other two classes))`.

Margin >= 0 is classified correct; margin < 0 is classified incorrect. This is not a tuned threshold.

## Measurements

Per subject × REPORT step:
- path A/B old retention before and after;
- path A/B correct→incorrect crossing count;
- path A/B incorrect→correct crossing count;
- path A/B minimum signed margin before/after;
- path A/B mean signed margin before/after;
- between-path count of examples with different correctness after the step.

Also report the exact crossing examples with split, name, verb, target class, before margin, and after margin.

## Interpretation

If coarse retention flips coincide with asymmetric boundary crossings, REPORT-on-path-state sensitivity is explained at the example decision-margin level. If retention changes occur without corresponding crossings, the proposed mechanism is falsified.

## Bounds

Diagnostic only. No retained diagnostic updates, no adaptive ordering, no extra training, no memory change, no architecture change, no post-result threshold or classifier, no result-informed retry, no live activation, no production authority.
