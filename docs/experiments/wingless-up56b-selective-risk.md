# Wingless UP-56B — observational selective-risk curve

Status: preregistered scientific self-evaluation diagnostic.

Scientific parent: UP-55B seal `7d95194326285268ba65b87a079436145dc1cc6c`.

UP-54B and UP-55B established that the existing decoder margin replicably predicts commit correctness. UP-56B asks how much error risk could be identified at different confidence cutoffs without changing system behavior.

Three new deterministic noise schedules are evaluated at noise 0.065 and 0.07 and 32/64-write chains. Every commit still occurs exactly as before. Afterward, the recorded events are scored observationally at frozen margin cutoffs 0, 0.05, 0.1, 0.25, 0.5, and 0.75, reporting coverage and accuracy among events above each cutoff.

No rejection, retry, oracle reset, or confidence-gated action occurs. This experiment measures self-evaluation usefulness only.


## Authoritative Windows result

Workflow run: `36034628624`

Runner: `WINGLESS-UP-B`

Source head: `8ead119c5397ab16d5e7381fb73fc1540abc0470`

Artifact: `10824250425`

Artifact digest: `sha256:1c00c5aaa68d02173e8cb1db47663e52b5b52c6bb897994c93d87e9ac2d245de`

Across all untouched schedules, increasing the preregistered margin cutoff consistently reduced coverage while increasing accuracy among retained decisions.

Representative noise 0.07 / 64-write results:

- schedule 59M: baseline `0.9215494791666666`; margin >=0.75 coverage `0.8011067708333334`, selected accuracy `0.9475822836245429`;
- schedule 60M: baseline `0.8909505208333334`; margin >=0.75 coverage `0.7965494791666666`, selected accuracy `0.9293011851246424`;
- schedule 64M: baseline `0.9137369791666666`; margin >=0.75 coverage `0.7815755208333334`, selected accuracy `0.9458558933777592`.

No commit behavior was changed.

## Scientific classification

The endogenous decoder margin supports a repeatable selective-risk tradeoff on untouched schedules. This establishes usefulness for self-evaluation, but not yet usefulness for correction.

The next B-lane experiment preregisters a sandboxed confidence-gated carry-forward policy on new schedules. At low confidence, the system retains its pre-decode internal table and applies the requested write rather than overwriting that table with the uncertain decoded table. Baseline closed-loop and several frozen cutoffs are compared; no oracle state is used.
