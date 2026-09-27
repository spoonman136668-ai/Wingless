# Wingless UP-196B — permutation-distance history geometry

Status: preregistered scientific shadow-diagnostic experiment.

Scientific parent: sealed UP-195B 44faf848f60913208fc6089932ba977f9f97e5c8.

## Question

Can Wingless's order-sensitive calibration effect be summarized by simple permutation distance from a canonical history?

## Frozen compositions

Two contrasting four-update compositions from UP-195B:
- mixed4: 0,5,1,6 — highest adjacent-swap sensitivity;
- observe4: 5,6,7,8 — lowest adjacent-swap sensitivity.

For each composition, enumerate all 24 permutations exactly once.

## Untouched evaluation phases

- 55;
- 56;
- 57.

## Frozen predictor

Reuse the pooled phase-26..30 model unchanged. No correction is fit or applied.

## Frozen order metric

Permutation distance = Kendall inversion count relative to the canonical sequence.

No alternative distance metric is selected after results.

## Measurements

Per composition × permutation × phase:
- Kendall distance;
- native correct count;
- required aggregate mass factor;
- absolute factor difference from canonical;
- whether native correct count matches canonical.

Per composition:
- Pearson association between Kendall distance and absolute factor difference;
- mean factor difference by distance;
- same-count nonzero-difference count.

## Interpretation

A strong monotonic distance relationship would support a compact order-distance coordinate. Weak or irregular structure would imply that specific transition identity, not generic permutation distance, carries the history signal.

## Bounds

Diagnostic only. No fitted correction, adaptive permutation selection, alternative distance search, held-out tuning, phase/parity predictor input, maintenance action, capacity change, live activation, or production authority.
