# Wingless UP-56B — observational selective-risk curve

Status: preregistered scientific self-evaluation diagnostic.

Scientific parent: UP-55B seal `7d95194326285268ba65b87a079436145dc1cc6c`.

UP-54B and UP-55B established that the existing decoder margin replicably predicts commit correctness. UP-56B asks how much error risk could be identified at different confidence cutoffs without changing system behavior.

Three new deterministic noise schedules are evaluated at noise 0.065 and 0.07 and 32/64-write chains. Every commit still occurs exactly as before. Afterward, the recorded events are scored observationally at frozen margin cutoffs 0, 0.05, 0.1, 0.25, 0.5, and 0.75, reporting coverage and accuracy among events above each cutoff.

No rejection, retry, oracle reset, or confidence-gated action occurs. This experiment measures self-evaluation usefulness only.
