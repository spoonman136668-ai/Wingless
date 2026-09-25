# Wingless UP-86A — factorized readout noise bridge

Status: preregistered scientific robustness experiment.

Scientific parent: sealed UP-85A Windows evidence `02d0e76c50e756cab685aec2a2b51f99254663d3`.

## Question

UP-84A and UP-85A localized and replicated a six-role one-step structural-coverage transition: 18 training states fail while 27 training states pass. Those experiments supplied exact engineered factorized features. UP-86A asks whether the validated readout remains usable when the upstream representation is imperfect, as it will be when later supplied by sequential/recurrent state rather than exact static features.

## Frozen design

Unchanged:

- six ternary roles and 729 joint states;
- frozen primary/secondary/tertiary structural coordinates;
- training set: exactly the 27-state cell primary==0, secondary==0, tertiary==0;
- held-out set: primary!=0, exactly 486 states;
- raw factorized one-hot representation, 18 features;
- role-factorized simplex representation, 12 features;
- six independent 3-class linear-softmax heads;
- zero initialization;
- exactly one optimization step;
- learning rate 1.0;
- no retraining under noise;
- per-role held-out gate 0.98.

Held-out feature-noise amplitudes are exactly:

- 0;
- 0.01;
- 0.025;
- 0.05;
- 0.10.

Noise is deterministic additive feature perturbation:
`x[j] += amplitude * sin((sample_index+1)*(feature_index+1)*17)`.

The same deterministic perturbation rule is applied independently to both representations.

No denoising, adaptive normalization, retraining, threshold change, seed search, feature clipping, representation change, result-informed retry, promotion, or activation is permitted.

## Interpretation

This is a robustness bridge, not an optimization experiment. Exact-state success with rapid noisy degradation would show that the structural representation is brittle and should not be treated as integration-ready. Stable accuracy under modest perturbation would support carrying the mechanism into sequence-state integration. Every point is sealed exactly as observed.
