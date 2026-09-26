# Wingless UP-177B — threshold frontier

Status: preregistered scientific held-out detector experiment.

Scientific parent: sealed UP-176B acfe77e44d2082f0db0d85ba9d7deba6f91f84f5.

## Question
What precision–recall tradeoff does the frozen rank-plus-temporal risk predictor achieve on a fresh phase when the warning threshold is changed without refitting the predictor?

## Frozen predictor
Use the exact UP-175B class × rank × temporal probabilities unchanged.

## Frozen thresholds
Evaluate three thresholds, fixed before phase-22 execution:
- 0.20
- 0.30
- 0.40

## Held-out evaluation
Use terminal phase 22 with the same six subjects, two 12-step paths, and 120 old examples per state.

## Measurements
Per threshold:
- TP/FP/FN/TN
- precision
- recall
- warning count
- false-warning count

No threshold is selected inside this experiment. The result is a frontier only.

## Bounds
Shadow only. No intervention, no adaptive threshold, no refitting, no extra training, no capacity change, no result-informed retry, live activation, or production authority.
