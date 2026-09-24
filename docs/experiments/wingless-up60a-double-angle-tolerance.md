# Wingless UP-60A — doubly-controlled angle tolerance

Status: preregistered scientific cognition diagnostic.

Scientific parent: UP-59A seal `1bc1ea41cb01f97f37c516f4745d77b31f112bf8`.

UP-59A isolated the doubly-controlled gate angle as the source of long-horizon composition error. UP-60A freezes the learned role-swap angle and perturbs only the doubly-controlled angle away from exact π/2 by a preregistered ladder: 0, -0.001, -0.0025, -0.005, -0.01, -0.02, -0.04, -0.08, -0.12, -0.16 radians.

Each angle is tested at mixed-program lengths 24, 48, and 128 under the unchanged cognition and numerical gates. No training or post-result search occurs.


## Authoritative Windows result

Workflow run: `36034399404`

Runner: `WINGLESS-LINKDEADKB`

Source head: `fab31257bfb05735fd67e3fe53001b4501a4914a`

Artifact: `10823069492`

Artifact digest: `sha256:3d0a786151c28f824b253347d1f137b2047d01dbbb029e0da1ea3340547d7b3f`

At length 128 the frozen cognition gate remains PASS for double-control angle offsets from exact π/2 through `-0.02` radians, then fails at `-0.04` (`0.9506172839506173`) and degrades further for larger errors. At length 48, `-0.04` still passes while `-0.08` fails. At length 24, errors through `-0.12` pass and `-0.16` fails.

Norm drift and round-trip error remain essentially machine precision throughout.

## Scientific classification

Required gate-angle precision tightens with composition horizon. The current learned double-control angle error (~-0.1315 rad) is far outside the tolerance needed for long programs.

The next A-lane experiment freezes optimizer hyperparameters and varies only training budget to determine whether the existing learner converges toward the required angle precision or plateaus away from it.
