# Wingless UP-55B — margin calibration replication

Status: preregistered scientific robustness confirmation.

Scientific parent: UP-54B seal `d5964e1a7d07d747c764220a48d1101b49cec814`.

UP-54B found that the existing decoder margin is associated with commit correctness. UP-55B tests whether that association generalizes to untouched deterministic noise schedules.

The frozen UP-54B margin bins are reused unchanged. Three new deterministic noise prefixes (56M, 57M, 58M) are evaluated at noise 0.065 and 0.07 and write counts 32 and 64. Training, geometry, decoder, commit behavior, and thresholds are unchanged.

This experiment is observational only: margin never gates, retries, rejects, or changes a commit. It reports calibration under independent noise realizations. No intervention is authorized by a positive result.
