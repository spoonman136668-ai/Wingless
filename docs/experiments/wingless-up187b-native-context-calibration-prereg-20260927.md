# Wingless UP-187B — native-context calibration

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-186B b95c60d4a41f91ab089b0983fda82157a797885d.

## Question

Can the repeating calibration residual be reduced by conditioning risk on Wingless's own pre-cleanup context state rather than on external phase labels?

## Frozen native context rule

UP-186B training phases 26..30 showed exactly two pre-cleanup correct-count states:
- 620 old examples correct;
- 623 old examples correct.

Freeze the native context rule before held-out evaluation:
- **context_A** if pre-cleanup correct count <= 621;
- **context_B** otherwise.

The threshold 621 is the midpoint between the only two training-state values. Phase number/parity is not an input.

## Frozen probability models

1. **pooled_26_30** — identical to UP-184B.
2. **native_context_26_30** — same UP-182B feature key plus the frozen native context class.

Both use Laplace cell estimates and smoothed global fallback.

## Held-out phases

31, 32, 33.

For each held-out phase:
1. compute the pre-cleanup correct count from native post-REPORT state only;
2. assign frozen context_A/B;
3. evaluate both probability models untouched.

## Measurements

Per model × phase:
- actual crossings;
- predicted crossing mass;
- predicted/actual ratio;
- absolute mass-ratio error;
- Brier score;
- ECE;
- AUROC;
- average precision.

Aggregate mean absolute mass-ratio error across held-out phases.

## Bounds

Shadow only. No evaluation fitting, phase/parity input, adaptive threshold, intervention, maintenance action, capacity change, extra model calls, live activation, or production authority.
