# Wingless autonomous Windows research qualification

This automation is intentionally separate from KTRADE/CKB production authority.

## Trigger

A research branch does **not** qualify on every push.

Qualification starts only when the branch receives:

`.wingless/qualification-request.json`

with schema:

`wingless.research-qualification-request.v1`

This lets an experiment be assembled across multiple commits without launching incomplete runs.

## Runner isolation

The workflow requires a dedicated self-hosted runner label:

`wingless-research`

Recommended install root:

`C:\actions-runner-wingless`

The existing production runner at `C:\actions-runner` is not modified.

## Production priority

Before any Wingless qualification begins, `scripts/Test-WinglessHostGuard.ps1` checks for an active `Runner.Worker.exe` under the production runner root.

If production is active, Wingless waits.

After admission, the Wingless qualification PowerShell process is lowered to `BelowNormal` priority and `GOMAXPROCS` defaults to 4. Child Go processes inherit the lower scheduling priority.

This does not give Wingless authority over KTRADE. It only makes Wingless yield host resources.

## Evidence

Each run uploads an artifact containing:

- `transcript.txt`
- `probe.json` when probe output is parseable
- `summary.json`
- `git-status.txt`

Scientific negative results do not fail the workflow when the experiment harness itself passed. Harness/build/test defects, unexpected worktree mutation, missing harness markers, or malformed probe output do fail the workflow.

## Acceptance boundary

The runner does not:

- declare scientific acceptance;
- rewrite thresholds;
- retry until green;
- seal research evidence into Git;
- modify accepted refs;
- invoke ckb-plane;
- access brokers or credentials;
- promote/deploy production code.

Interpretation and sealing remain outside the execution runner.
