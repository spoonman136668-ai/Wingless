# Wingless UP-54B — endogenous margin calibration

Status: preregistered scientific robustness diagnostic.

Scientific parent: UP-53B seal `04f996031423faa596b07b150d8e30fb4b02b52b`.

UP-53B established both residual per-step decoder noise and recurrent amplification. UP-54B changes neither training nor commit behavior. It records the decoder's existing value margin for every closed-loop commit at noise 0.065 and 0.07 over 16, 32, and 64 writes.

Margins are frozen into bins [0,.001), [.001,.01), [.01,.05), [.05,.1), [.1,.25), [.25,1]. Accuracy per bin and mean margins for correct versus incorrect decodes are reported.

The question is whether the system already exposes a usable endogenous confidence signal. No thresholding, rejection, retry, or corrective action occurs in this experiment.


## Authoritative Windows result

Workflow run: `36033248133`

Runner: `WINGLESS-UP-B`

Source head: `50c42a8e06ed0758aed5b1e2c87a7f77ba422114`

Artifact: `10823387770`

Artifact digest: `sha256:d5174b4b215448e9023d5923d105a1a63cccb715719f6fc871cbb9debebd5526`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Across noise 0.065/0.07 and 16/32/64-write closed-loop runs, correct commits have mean decoder margins around `0.856–0.865`; incorrect commits have lower mean margins around `0.597–0.643`.

At noise 0.07 / 64 writes, frozen-bin accuracy was:

- margin [0,0.001): 0/1;
- [0.001,0.01): 1/2;
- [0.01,0.05): 5/22;
- [0.05,0.1): 10/30;
- [0.1,0.25): 42/70;
- [0.25,1]: 2726/2947 = `0.9250084832032576`.

## Scientific classification

The existing decoder margin contains a substantial endogenous confidence signal without adding a confidence model or changing behavior. It is informative but not perfectly calibrated; high-margin errors still occur.

The next B-lane experiment freezes these bins and repeats the analysis under independent deterministic noise schedules. No confidence threshold or intervention will be introduced until calibration generalization is confirmed.
