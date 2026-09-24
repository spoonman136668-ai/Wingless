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
