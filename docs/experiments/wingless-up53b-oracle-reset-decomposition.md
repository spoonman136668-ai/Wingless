# Wingless UP-53B — closed-loop versus oracle-reset decomposition

Status: preregistered scientific robustness experiment.

Scientific parent: UP-52B seal `a18fafcbd02bfc7ee5f29eaef43a657f10146985`.

## Question

UP-52B showed that pair-[0,5] memory survives short chains at both boundary noise levels but degrades as decoded state is repeatedly committed back into the next step.

UP-53B separates **per-step decoder error** from **closed-loop error propagation** by comparing the unchanged recurrent path to an oracle-reset path that supplies the true post-write table as the next step's input. The decoder, geometry, noise, transport, training, and evaluation remain unchanged.

## Frozen design

- pair-[0,5] geometry;
- noise 0.065 and 0.07;
- write counts 16, 32, 64;
- held depths 32, 128, 512, 1024;
- same trained observer/decoder and full-capacity control as the UP-52B family;
- closed-loop path: the decoded table is committed and becomes the next input;
- oracle-reset path: accuracy is measured at every decode, but the true post-write table becomes the next input;
- no retraining difference between the two paths.

## Interpretation

If oracle-reset accuracy remains high where closed-loop accuracy falls, the dominant mechanism is recurrent error amplification. If oracle-reset also falls similarly, the boundary is primarily per-step decoding/noise sensitivity.

Scientific negatives are valid. No geometry, noise, threshold, decoder, or optimizer changes are permitted after execution.
