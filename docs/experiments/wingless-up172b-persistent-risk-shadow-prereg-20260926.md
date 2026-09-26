# Wingless UP-172B — persistent risk shadow

Status: preregistered scientific internal-deficiency diagnostic.

Scientific parent: sealed UP-171B 4516ed4184145beebee4e368b19220f714187313.

## Question

UP-171B showed that class-calibrated low-margin bands cover most chained decision crossings but over-flag many examples. Does remaining inside the frozen vulnerability band for at least two consecutive chained steps predict an actual crossing more precisely than one-step membership?

## Frozen trajectory

Reuse UP-171B exactly:
- phase-17 post-REPORT starting state;
- six subjects;
- two fixed cleanup paths:
  - STORE_OBSERVE_REPORT;
  - OBSERVE_STORE_REPORT;
- twelve chained cleanup updates per path;
- 120 old examples evaluated before and after every step;
- class-specific frozen vulnerability bands:
  - STORE rank <= 10;
  - OBSERVE rank <= 10;
  - REPORT rank <= 4.

No update is changed by the diagnostic.

## Frozen persistence rule

For every old example independently, before each chained update:
- if its absolute-margin rank is inside the current step's frozen class band, increment its consecutive in-band streak;
- otherwise reset the streak to zero.

Preregistered bins:
- out_of_band = streak 0;
- current_only = streak 1;
- persistent_2 = streak 2;
- persistent_3plus = streak >= 3.

The primary comparison is:
- one-step warning = streak >= 1;
- persistent warning = streak >= 2.

## Measurements

Across all trajectories and separately by update class:
- example-slots per streak bin;
- crossings per streak bin;
- crossing density per streak bin;
- one-step flagged slots, covered crossings, crossing density, and crossing coverage;
- persistent flagged slots, covered crossings, crossing density, and crossing coverage.

A crossing remains the same zero-margin correctness change used by UP-165B..UP-171B.

## Interpretation

Higher persistent-warning crossing density than one-step-warning density would show that temporal persistence improves warning precision. Lower coverage is allowed and is part of the scientific result. Failure to improve density falsifies the persistence hypothesis.

No threshold is fitted and no maintenance action is permitted.

## Bounds

Shadow diagnostic only. No maintenance trigger, no update suppression, no adaptive band, no schedule change, no extra training, no capacity change, no result-informed retry, no live activation, no production authority.
