# Wingless UP-171B — shadow vulnerability risk through chained cleanup

Status: preregistered scientific internal-deficiency shadow diagnostic.

Scientific parent: sealed UP-170B 7447e0cf2a50b0fd76b63dea442c3c0a3b2962e2.

## Question

UP-170B established a monotonic relationship between low absolute-margin rank and one-step crossing risk on a disjoint phase-17 context. Does the frozen class-calibrated risk band continue to cover actual crossings when cleanup updates are chained and state changes accumulate?

## Frozen context

Reuse the phase-17 post-REPORT context and six accepted subjects.

For each subject, trace both fixed cleanup paths:
- STORE_OBSERVE_REPORT:
  - STORE indices 0..4;
  - OBSERVE indices 5..9;
  - REPORT indices 13..14.
- OBSERVE_STORE_REPORT:
  - OBSERVE indices 5..9;
  - STORE indices 0..4;
  - REPORT indices 13..14.

Each path contains exactly 12 real cleanup updates. No clone reset occurs between steps within a path.

## Frozen shadow bands

Carried from the sealed preregistered class calibration:
- STORE: vulnerability rank <= 10;
- OBSERVE: vulnerability rank <= 10;
- REPORT: vulnerability rank <= 4.

Bands are fixed before this run.

Before every actual cleanup step:
1. compute all 120 old-example probability margins;
2. rank by absolute margin;
3. mark examples inside that step class's frozen band;
4. apply the already-scheduled cleanup update;
5. record actual correctness crossings.

The shadow signal never changes or suppresses an update.

## Measurements

Per subject × path × step:
- class and old-surface index;
- band maximum rank;
- flagged examples;
- actual crossings;
- crossings covered by the band;
- false-warning slots;
- coverage;
- crossing density among flagged examples;
- minimum/mean absolute pre-step margin.

Aggregate by class and path.

## Interpretation

High coverage under chained state changes supports use of the rank signal as a persistent shadow warning. Coverage collapse would show that the independent-probe calibration does not survive evolving state. False-warning rate is recorded explicitly; this experiment does not authorize a trigger.

## Bounds

No maintenance action, no schedule change, no threshold fitting, no adaptive band, no retained diagnostic clone, no extra training, no memory change, no live activation, no production authority.
