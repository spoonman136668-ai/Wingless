# Wingless UP-56A — doubly-controlled reversible composition

Status: preregistered scientific cognition experiment.

Scientific parent: UP-55A seal `6a0f346422c273998953b8a84793614456369dc3`.

UP-55A established perfect transfer for single-control context-dependent operations. UP-56A adds a second simultaneous condition: a target value changes only when two other roles match their specified control values.

Training uses only target role 0, controls roles 1/2, control tuple (0,0), plus swap 0/1. Held-out evaluation covers all other control-value tuples, moved target roles, and mixed programs of lengths 16 and 48.

The same cognition gates remain frozen. Scientific negatives are valid; no post-result topology, optimizer, or threshold change is allowed.


## Authoritative Windows result

Workflow run: `36032761782`

Runner: `WINGLESS-LINKDEADKB`

Source head: `ceb6fa089b4955c143a9e20ea554f7139f8c9def`

Artifact: `10822897389`

Artifact digest: `sha256:901bf6def2a15c7626cc122b022f448da4bc7e0950b147bbbc481db45bc91213`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Unitary results:

- train accuracy: `1.0`;
- aggregate held-out accuracy: `0.996031746031746`;
- unseen control-tuple accuracy: `1.0`;
- moved-target accuracy: `1.0`;
- long-program accuracy: `0.9444444444444444`;
- maximum norm drift: `1.5543122344752192e-15`;
- maximum round-trip error: `2.043470511879551e-15`;
- frozen double-control gate: FAIL because long-program accuracy is below 0.98.

## Scientific classification

The doubly-controlled operation itself transfers perfectly across unseen control tuples and target locations, but repeated mixed composition accumulates logical error before any meaningful norm or reversibility failure.

This is a scientific boundary, not an optimizer failure. The next A-lane experiment freezes the learned task and maps accuracy versus program length to locate the composition-depth boundary without changing training, gates, or parameters.
