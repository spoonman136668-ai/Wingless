# Free Fanout Sidecar R1

This directory defines a dormant, external-evidence-only screening lane for Wingless. It expands historical/replay candidate search breadth without changing scientific or execution authority.

## Authority and isolation

The sidecar is not the accepted research loop and cannot classify, promote, dispatch successors, mutate accepted state, write CKB-plane/KTRADE queues, touch production, or use sealed unseen outcomes for candidate selection. The workflow is manual `workflow_dispatch` only, uses GitHub-hosted `ubuntu-latest`, has read-only repository permission, and runs candidates with a minimal scrubbed environment rather than inherited job credentials.

CKB-plane remains the sole official execution/governance authority. Sidecar artifacts are advisory external evidence only.

## Two immutable identities

The screening manifest and sidecar implementation live at the workflow-dispatch revision. The manifest binds a separate exact 40-hex **frozen package commit** through `package_sha`.

This separation is intentional: a file inside a Git commit cannot contain that same commit's SHA without a self-reference problem. The workflow therefore:

1. checks out the sidecar/manifest revision;
2. checks out `package_sha` separately as `frozen-package`;
3. proves `git rev-parse HEAD == package_sha`;
4. runs only the manifest-declared entrypoint from `frozen-package/research/fanout/`.

The package checkout uses `persist-credentials: false`.

## Candidate interface

Target interface: `state_former(raw_input, history) -> fixed_size_state`

For Wingless the manifest freezes the state dimension and readout capacity. Candidate results must echo those fixed values exactly.

The manifest also freezes a resource envelope against a preregistered baseline: parameters, context bytes, model calls, peak RSS, and elapsed runtime. A candidate that exceeds any bound is rejected rather than counted as an improvement.

Candidate configuration is frozen in the manifest and passed as a JSON file. The entrypoint contract is:

`entrypoint --candidate <id> --candidate-config <json> --output <result.json>`

Entrypoints and manifests must remain under `research/fanout/`; traversal or symlink escape is rejected.

## Data and result contract

Allowed search data classes are `historical`, `replay`, and `historical-replay`. `sealed_inputs` and sealed-outcome exposure must be false. Candidate results must include a deterministic metric and bounded resource usage. The runner emits a SHA256-bound envelope with measured elapsed/CPU/RSS data.

The aggregate fails closed unless exactly one valid envelope exists for every preregistered candidate. It always retains:

- `scientific_classification = null`
- `promotion_authority = false`
- `successor_authority = false`

## Smoke package

`smoke_candidate.py` and `smoke_fixture.json` are deterministic synthetic/replay-only integration fixtures. `smoke.manifest.json` binds the exact package commit containing those fixtures. The smoke lane is not a scientific experiment and must not use the currently active experiment as its integration target.

## Future opt-in flow

`frozen historical/replay package -> bounded candidate fanout -> resource/complexity filtering -> frozen shortlist -> official sealed evaluation through existing authority`

Do not scale immediately to 32 candidates. Start with 3-5 when a real experiment opts in, preserve failed candidates as search history, and keep all preregistration/freeze/no-post-result-tuning rules intact.
