# Wingless UP-54B — endogenous margin calibration

Status: preregistered scientific robustness diagnostic.

Scientific parent: UP-53B seal `04f996031423faa596b07b150d8e30fb4b02b52b`.

UP-53B established both residual per-step decoder noise and recurrent amplification. UP-54B changes neither training nor commit behavior. It records the decoder's existing value margin for every closed-loop commit at noise 0.065 and 0.07 over 16, 32, and 64 writes.

Margins are frozen into bins [0,.001), [.001,.01), [.01,.05), [.05,.1), [.1,.25), [.25,1]. Accuracy per bin and mean margins for correct versus incorrect decodes are reported.

The question is whether the system already exposes a usable endogenous confidence signal. No thresholding, rejection, retry, or corrective action occurs in this experiment.
