# Wingless UP-56C — golden-angle untouched confirmation

Status: preregistered scientific scale confirmation.

Scientific parent: UP-55C seal `986d4f350c985b753f89fa535e4040c54f3e7c21`.

UP-55C's exploratory family comparison found that the preregistered golden-angle phase tags uniquely passed the existing four-bank gate at noise 0.04. That result is not promoted.

UP-56C freezes the exact golden tags and tests them on three untouched deterministic noise schedules at noise amplitudes 0.035, 0.04, 0.0425, 0.045, 0.0475, and 0.05. Dimension, bank count, scenarios, transport depth, decoder, global-phase nuisance, magnitude control, and gate remain unchanged.

No new tag family is searched or selected. The purpose is untouched replication and boundary mapping.


## Authoritative Windows result

Workflow run: `36034486491`

Runner: `WINGLESS-UP-C`

Source head: `474a4de77b172421a49d2a52da01e34c950da359`

Artifact: `10822879738`

Artifact digest: `sha256:e1f26c15b0a91725b16986992d07fa05ae9de3e4ee53d93bf2d102fd021af982`

The first attempt failed during concurrent shared Go-cache cleanup. The harness-only FIXA isolated `GOCACHE/GOTMPDIR`; the frozen scientific experiment then passed without changing tags, noise levels, scenarios, decoder, or gates.

Across all three untouched deterministic schedules:

- noise 0.035: PASS;
- noise 0.04: PASS;
- noise 0.0425: FAIL;
- noise 0.045: FAIL;
- noise 0.0475: FAIL;
- noise 0.05: FAIL.

At noise 0.04 exact-scenario accuracy ranged from `0.95703125` to `0.9765625`; at 0.0425 it fell to `0.88671875–0.92578125`.

## Scientific classification

The golden-angle four-bank robustness improvement replicates on independent noise schedules. Its frozen gate boundary lies between 0.04 and 0.0425 under this fixed-state configuration.

The next C-lane experiment uses the confirmed golden construction to test whether bank-count capacity extends beyond four across a preregistered bank-count/noise grid.
