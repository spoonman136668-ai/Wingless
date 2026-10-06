# Free Fanout Sidecar R1

This directory defines a dormant, external-evidence-only screening lane for Wingless.

## Isolation contract

The sidecar is not part of the accepted research loop. It MUST NOT:

- use self-hosted runners;
- consume CKB-plane or KTRADE queues;
- write accepted research state;
- dispatch successors;
- touch production authority;
- use sealed unseen inputs for candidate selection;
- receive API/model secrets;
- promote branches or modify the North Star.

The workflow uses GitHub-hosted `ubuntu-latest` only, has `contents: read` permission, and uploads artifacts only.

## Activation boundary

While this implementation remains on `infra/free-fanout-sidecar-r1`, it is dormant. Do not merge it to `main` or dispatch it until the active research lanes reach an explicit safe boundary.

After activation, a frozen experiment package may opt in by adding:

1. a manifest under `research/fanout/*.json`;
2. an entrypoint under `research/fanout/`;
3. 1-32 preregistered candidate IDs.

The manifest must bind an exact 40-hex package SHA and declare:

- `authority = external-evidence-only`;
- historical/replay-only data class;
- no sealed inputs;
- no accepted-state mutation;
- no successor dispatch;
- no queue writes;
- no production access.

Each candidate runs independently on a GitHub-hosted Linux runner. Results are aggregated as external evidence only. The aggregator deliberately emits no scientific classification and has no promotion authority.

## Scientific use

The intended use is cheap breadth before a sealed experiment:

`many historical/replay candidates -> frozen shortlist -> official sealed evaluation`

Shortlisting and any later scientific claim remain under the existing preregistration, qualification, and CKB-plane governance rules.
