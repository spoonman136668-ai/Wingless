# Stage 3.5 isolated ckb-plane integration proof

State: **IMPLEMENTED — ISOLATED HARNESS PASS — WINDOWS ACCEPTANCE PENDING**

## Frozen authority input

- ckb-plane source: `spoonman136668-ai/ckb-plane@539d56fec273c8851aea131cf4d31425a62f7250`
- exact-source Windows closeout: `acceptance/frozen-contract-closeout-20260911.json`
- Wingless plane adapter: `integration/ckbplane`
- live activation: forbidden

The proof uses normalized recorded identity from the frozen `ckb-r4a-source-proof` profile. It does not connect to a live plane, consume a queue, launch a plane worker, mutate accepted refs, launch CKB, access Coinbase/broker APIs, read credentials or touch production state.

## Proof paths

`integration/ckbplane/proof_test.go` exercises the real Wingless mock/broker/internal-service stack through the inactive adapter:

1. **candidate flow** — recorded plane work/workspace identity translates to one bounded Wingless request; `inference.Mock` returns one candidate; result remains `external_required`.
2. **authorized deep fallback** — an unhealthy authorized deep mock falls back to the authorized fast mock and preserves route rejection evidence.
3. **backend failure / no duplicate retry** — backend failure remains `operator_blocked`, produces exactly one attempt, and does not consume a Wingless repair loop because plane owns work-order retry.
4. **cancellation propagation** — plane work-order ID + attempt deterministically maps to the active Wingless request; cancellation propagates through the real internal `worker.Service`; authoritative plane lifecycle remains external.

The existing adapter replay tests additionally cover stale workspace revision, work-order/workspace mismatch, required plane protection gates, trusted backend policy, resource policy translation, unknown input fields, acceptance-claim rejection and unknown terminal states.

## Current validation

A narrow isolated compile/replay harness using the fetched Wingless API shapes passed all four proof tests.

GitHub `native.yml` cannot currently provide a repository verdict: hosted `windows-latest` and `ubuntu-latest` jobs repeatedly terminate before step 1 with `runner_id=0` and empty step lists. This is runner infrastructure unavailability, not a source/test result.

The authoritative acceptance gate is therefore native Windows PowerShell using:

```powershell
.\scripts\Verify-Isolated-Plane-Integration.ps1
```

The verifier requires exactly `go version go1.25.5 windows/amd64` and performs:

- focused `integration/ckbplane` tests;
- reconciled broker/inference/resources/worker/contextbroker package tests;
- full `go test ./...`;
- source-clean Wingless command build to `%TEMP%`;
- sanitized evidence emission under ignored `evidence/`.

A passing verifier must end with `WINGLESS STAGE 3.5 ISOLATED INTEGRATION PROOF PASS` and record zero live-plane transport, zero live queue mutation, zero accepted-ref mutation, zero runtime/broker/credential/production use, candidate acceptance `external_required`, and Wingless work-order retry budget `0`.

## Gate

Until native Windows passes, the isolated proof is implemented but not authoritatively accepted. No plane-side `wingless` worker contract or live transport may be added before this gate closes.
