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


## Authoritative Windows result

Workflow run: `36032468246`

Runner: `WINGLESS-UP-B`

Source head: `e265ad3ea3b380af7b01b818ed32d6fc4d5ba045`

Artifact: `10823426908`

Artifact digest: `sha256:d3620f62c52e310c5409e8c7449e095ad6e8853e169348c38d814dc1b9e85149`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Representative results:

- noise 0.065 / 16 writes: closed commit `0.9518229166666666`, oracle-reset `0.9830729166666666`;
- noise 0.065 / 64 writes: closed `0.923828125`, oracle-reset `0.9762369791666666`;
- noise 0.07 / 16 writes: closed `0.9453125`, oracle-reset `0.9817708333333334`;
- noise 0.07 / 64 writes: closed `0.90625`, oracle-reset `0.970703125`.

The full-capacity control remains perfect.

## Scientific classification

Closed-loop error amplification is a major part of the pair-[0,5] long-chain failure: resetting each step to the true post-write state recovers roughly 3–6.4 percentage points of commit accuracy and substantially improves final-table accuracy.

Oracle reset does not restore perfect decoding, so a smaller independent per-step noise/decoder error remains. The next B-lane diagnostic will measure the decoder margin/error relationship without changing commit behavior, to test whether the system already contains a usable endogenous confidence signal.
