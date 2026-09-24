# Wingless UP-53A — three-basis identifiability control

Status: preregistered scientific cognition experiment.

Scientific parent: UP-52A seal `c2cbcfa72c1169e52b3873032a731058beef2dfb`.

## Question

UP-52A found that six observed basis configurations pass the unchanged cognition gate while the frozen three-state subset `000/111/222` fails.

However, every diagonal state is unchanged by role swap 0/1, so that subset provides no training signal for the role-swap parameter. The failure may therefore be an identifiability defect in the chosen examples rather than a true minimum-coverage boundary.

UP-53A holds training count at exactly three states and changes only whether those states make the role-swap operation observable.

## Frozen triads

- invariant diagonal: `000, 111, 222`;
- informative cycle A: `012, 120, 201`;
- informative cycle B: `021, 102, 210`.

All three conditions use the unchanged UP-52A/UP-51A optimizer, single-step primitive supervision, held-out categories, long-program lengths, matched non-unitary control, and numerical/accuracy gates.

## Interpretation

If both informative triads pass while the invariant diagonal fails, the UP-52A three-state failure is attributable to training-example identifiability rather than basis-count insufficiency.

If informative triads also fail, three-state coverage itself remains a live bottleneck.

No post-result parameter or threshold changes are permitted.


## Authoritative Windows result

Workflow run: `36030967147`

Runner: `WINGLESS-LINKDEADKB`

Source head: `791ead9caf3cc3745a0449696ac84958e8d375b7`

Artifact: `10822205622`

Artifact digest: `sha256:dc5f667a55974d59eb8857495eb9680ead25077162b3b433091b8c618599427b`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Results:

- invariant diagonal `000/111/222`: gate FAIL; held-out `0.28125`; role-swap parameter remained `0.1`; derived role-1 `0.08333333333333333`; derived role-2 `0.08333333333333333`; long-program `0.125`;
- informative cycle A `012/120/201`: gate PASS; held-out `1.0`; derived role-1/role-2/long-program all `1.0`; learned parameters approximately `[1.5707961289495511, 1.5701448070025685]`;
- informative cycle B `021/102/210`: gate PASS; held-out `1.0`; derived role-1/role-2/long-program all `1.0`; learned parameters approximately `[1.5707961289495511, 1.5701448070025685]`.

## Scientific classification

The UP-52A three-state failure was an **identifiability failure**, not evidence that six basis examples are intrinsically required.

Exactly three examples are sufficient in this task when they make both learned primitive operations observable. The diagonal triad failed because role swap has no effect on those states and therefore supplies no gradient for the role-swap parameter.

This closes the immediate basis-count ambiguity. The A lane should now return to broader cognition: increase structural/task complexity rather than spend more budget narrowing a toy minimum-example boundary.
