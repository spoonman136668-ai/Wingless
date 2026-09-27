# Wingless UP-195B — adjacent-swap history localization

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-194B 4e2636ed3871b6ec2adc08318ede8567c302a320.

## Question

Is Wingless's order-sensitive calibration effect distributed across local temporal transitions, or concentrated in one particular ordering boundary?

## Frozen compositions

Canonical histories:
- store4: 0,1,2,3;
- observe4: 5,6,7,8;
- mixed4: 0,5,1,6;
- cross3: 0,5,13.

For each composition, evaluate every single adjacent transposition of the canonical sequence. Composition is unchanged; only one neighboring pair is swapped.

## Untouched evaluation phases

- 52;
- 53;
- 54.

## Frozen predictor

Reuse the pooled phase-26..30 probability model unchanged. No correction is fit or applied.

## Measurements

Per composition × swap position × phase:
- canonical native correct count;
- swapped native correct count;
- canonical required mass factor;
- swapped required mass factor;
- absolute factor difference;
- absolute native-count difference;
- whether native correct count is unchanged.

Aggregate:
- mean factor difference by composition;
- mean factor difference by swap position;
- count of same-count swaps with nonzero calibration difference.

## Interpretation

Nonzero effects across several adjacent positions support distributed temporal sensitivity. A single dominant boundary would localize the history dependence more narrowly.

## Bounds

Diagnostic only. No adaptive swap search, fitted correction, held-out tuning, phase/parity predictor input, maintenance action, capacity change, live activation, or production authority.
