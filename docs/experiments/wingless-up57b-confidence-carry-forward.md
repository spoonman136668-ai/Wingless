# Wingless UP-57B — confidence-gated carry-forward

Status: preregistered sandboxed robustness/self-evaluation experiment.

Scientific parent: UP-56B seal `4e56b9fb943e82b72b8c5172c8a211f81b30665d`.

UP-56B established a repeatable accuracy-versus-coverage tradeoff from the decoder's endogenous margin. UP-57B asks whether that signal can reduce closed-loop error accumulation without oracle information.

At each commit the normal decoder still runs. Baseline always replaces the internal table with the decoded table before applying the requested write. Confidence-carry arms use frozen margin cutoffs 0.25, 0.5, and 0.75; when margin is below the cutoff, the arm retains its pre-decode internal table and applies the requested write to that retained state instead of overwriting it with an uncertain decode.

Three untouched deterministic schedules are tested at noise 0.065 and 0.07 with 64 writes. No oracle state, retry, retraining, threshold adaptation, or production authority is used. Scientific negatives are valid.


## Authoritative Windows result

Workflow run: `36035159819`

Runner: `WINGLESS-UP-B`

Source head: `81f5962b2ac747d40479d711baf12dbd89d72b58`

Artifact: `10825040612`

Artifact digest: `sha256:7eb5f34d9d1b1dec7eb5f0a16b48e437df959a95748687ec92ec723c02a57e77`

Across all three untouched schedules at noise 0.065 and 0.07, confidence carry increased commit accuracy relative to the closed-loop baseline. At the preregistered 0.5 margin arm, representative changes were:

- noise 0.065 / schedule 67M: commit `0.9186 -> 0.9730`, final `0.8750 -> 0.9583`;
- noise 0.065 / schedule 68M: commit `0.9349 -> 0.9772`, final `0.8958 -> 0.9583`;
- noise 0.070 / schedule 67M: commit `0.8910 -> 0.9652`, final `0.8750 -> 0.9583`;
- noise 0.070 / schedule 68M: commit `0.9160 -> 0.9720`, final `0.8750 -> 0.9583`.

Relational accuracy was preserved.

## Scientific classification

The endogenous decoder margin is not merely correlated with error: using it causally to prevent low-confidence state replacement reduces closed-loop error amplification without oracle state.

The next B-lane experiment should test this policy family on new schedules, write horizons, and nearby noise levels without selecting a winning threshold from UP-57B.
