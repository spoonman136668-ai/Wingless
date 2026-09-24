# Wingless UP-57B — confidence-gated carry-forward

Status: preregistered sandboxed robustness/self-evaluation experiment.

Scientific parent: UP-56B seal `4e56b9fb943e82b72b8c5172c8a211f81b30665d`.

UP-56B established a repeatable accuracy-versus-coverage tradeoff from the decoder's endogenous margin. UP-57B asks whether that signal can reduce closed-loop error accumulation without oracle information.

At each commit the normal decoder still runs. Baseline always replaces the internal table with the decoded table before applying the requested write. Confidence-carry arms use frozen margin cutoffs 0.25, 0.5, and 0.75; when margin is below the cutoff, the arm retains its pre-decode internal table and applies the requested write to that retained state instead of overwriting it with an uncertain decode.

Three untouched deterministic schedules are tested at noise 0.065 and 0.07 with 64 writes. No oracle state, retry, retraining, threshold adaptation, or production authority is used. Scientific negatives are valid.
