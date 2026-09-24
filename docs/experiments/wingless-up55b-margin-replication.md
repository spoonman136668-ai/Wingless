# Wingless UP-55B — margin calibration replication

Status: preregistered scientific robustness confirmation.

Scientific parent: UP-54B seal `d5964e1a7d07d747c764220a48d1101b49cec814`.

UP-54B found that the existing decoder margin is associated with commit correctness. UP-55B tests whether that association generalizes to untouched deterministic noise schedules.

The frozen UP-54B margin bins are reused unchanged. Three new deterministic noise prefixes (56M, 57M, 58M) are evaluated at noise 0.065 and 0.07 and write counts 32 and 64. Training, geometry, decoder, commit behavior, and thresholds are unchanged.

This experiment is observational only: margin never gates, retries, rejects, or changes a commit. It reports calibration under independent noise realizations. No intervention is authorized by a positive result.


## Authoritative Windows result

Workflow run: `36034174867`

Runner: `WINGLESS-UP-B`

Source head: `d768a2ba8a701053d7ae853d3c6a26a8769b39fe`

Artifact: `10822794533`

Artifact digest: `sha256:bd5b02bffc3154bc832de6d58af9eef71d70ddfdf5cdeedcf432eae3cdde506a`

All three untouched deterministic noise schedules reproduced the confidence separation at both noise levels and both write counts.

Examples:

- schedule 56M / noise 0.065 / 64 writes: low-margin accuracy `0.45`, high-margin accuracy `0.9383561643835616`;
- schedule 57M / noise 0.07 / 64 writes: low `0.4246575342465753`, high `0.9209716045159083`;
- schedule 58M / noise 0.07 / 64 writes: low `0.5`, high `0.9393112853733379`.

For every one of the 12 frozen schedule/noise/write conditions, high-margin accuracy exceeded low-margin accuracy.

## Scientific classification

The decoder margin is a replicated endogenous confidence signal rather than a schedule-specific artifact. It remains imperfect, so this result does not authorize confidence-gated behavior.

The next B-lane experiment remains observational: measure the frozen accuracy-versus-coverage/selective-risk curve on new held-out noise schedules for several preregistered margin cutoffs. No commit is actually rejected or retried.
