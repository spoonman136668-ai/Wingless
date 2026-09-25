# Wingless UP-93B — confidence-qualified endogenous change routing

Status: preregistered scientific memory-routing experiment.

Scientific parent: sealed UP-92B `ece60e33e57c3bf830948f2a2cf4d1487b3af085`.

## Question

UP-92B found that local change detection has high recall but poor precision because interference creates false change detections. Can the existing endogenous decoder margin reject enough low-confidence false changes to make bounded admission useful without an external salience bit?

## Frozen substrate

- same 64-float recurrent carrier;
- same gated-correction update;
- same presence threshold 0.25;
- exact-recall capacity exactly 16;
- no attention;
- no state-size increase.

## Frozen arms

1. `explicit_target_rewrite` — upper-bound control.
2. `endogenous_change_raw` — UP-92B rule.
3. `endogenous_change_margin_025`
4. `endogenous_change_margin_050`
5. `endogenous_change_margin_075`

For margin-qualified arms:
- a write is eligible only when prior-value presence is true;
- decoded old value differs from incoming value;
- normalized pre-write decoder margin is at least the frozen threshold;
- normalized margin = (best_score - second_score) / max(1, abs(best_score)+abs(second_score)).

Thresholds are diagnostics only; none is selected post hoc.

## Frozen task

Same sparse target-change task as UP-92B:
- total writes: 32, 64, 128, 256;
- target keys: 4, 8, 16;
- value vocabulary 32;
- 64 episodes per setting;
- new seeds: 131M and 132M.

## Metrics

- target query accuracy;
- whole-target-set exact accuracy;
- recall entries used;
- admission precision;
- admission recall;
- false-positive admissions.

## Bounds

No external salience bit for endogenous arms, no future-query oracle, no threshold tuning or selection after results, no capacity increase, no result-informed retry, no production authority, no live activation.
