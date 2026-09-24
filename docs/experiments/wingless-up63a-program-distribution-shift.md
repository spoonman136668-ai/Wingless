# Wingless UP-63A — held-out program-distribution shift

Status: preregistered scientific cognition experiment.

Scientific parent: UP-61A seal `e3b319607b677dd7e9537d9651767c013e50d07e`.

This experiment is independent of UP-62A. It keeps the confirmed 1440-step learner and tests whether its learned reversible primitives generalize beyond the program generator used by the earlier long-horizon tests.

Frozen held-out families are control-heavy bursts, swap-heavy programs, shifted-control-value programs, and self-inverse palindromic programs at lengths 32, 64, 128, and 256. Training remains the original single-step primitive supervision.

The existing 0.98 cognition/category gate and numerical gates remain unchanged. No held-out program family appears in training.


## Authoritative Windows result

Workflow run: `36035702828`

Runner: `WINGLESS-LINKDEADKB`

Source head: `537436dfb379c2d1bbf5bea3a9c33e7557f1f2e9`

Artifact: `10825280865`

Artifact digest: `sha256:c73db677437e3231b79c93062cdf4829c6c8725d69251d49d31abdf39243cace`

Results:

- aggregate held-out accuracy: `1.0`;
- every frozen program family passed at lengths 32, 64, 128, and 256;
- control-burst accuracy: `1.0`;
- swap-heavy accuracy: `1.0`;
- shifted-control accuracy: `1.0`;
- palindromic-program accuracy: `1.0`;
- maximum norm drift: `5.551115123125783e-15`;
- maximum round-trip error: `5.551115123126598e-15`;
- frozen gate: `PASS`.

## Scientific classification

The learned reversible primitives generalize across substantial held-out program-distribution shifts through length 256. The previous long-horizon success is therefore not specific to one program generator.

The next A-lane question should move beyond additional symbolic program-length ladders and test whether the same substrate can acquire and preserve an additional capability without losing the already-qualified reversible program behavior.
