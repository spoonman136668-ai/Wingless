# Plane reconciliation — Stage 3.5 intake

State: **NOT_READY_FOR_ISOLATED_INTEGRATION_PROOF**

## Frozen inputs

Wingless repository: `spoonman136668-ai/Wingless`, branch `main`.
Reviewed documentation baseline: `584c92a04b9a9944fee491cc59a4025ea4143ffc`.
Review date: 2026-09-11.

Frozen plane contract source is now available:

- repository: `spoonman136668-ai/ckb-plane`
- frozen source commit: `539d56fec273c8851aea131cf4d31425a62f7250`
- frozen source root: `snapshot/source/`
- source manifest: `snapshot/SOURCE-MANIFEST.sha256`
- snapshot timestamp: `2026-09-11T16:47:33.7746436Z`
- source files: 130
- Go files: 75
- source bytes: 758340
- live activation authorized by snapshot: false
- runtime state included: false
- acceptance evidence included: false

The ckb-plane repository subsequently received only advisory ICE metadata and a Windows closeout verifier; the frozen `snapshot/` tree remains the source identity established by the commit above.

Acceptance linkage for the exact frozen source is **pending**. Historical DN/DO/DP evidence and a passing post-DT controller-owned export proof exist, but they do not by themselves prove the later DU/DW/DX/DXA/DXB source contained in the frozen snapshot. The repository therefore includes `scripts/Verify-Frozen-Contract.ps1` to establish a direct Windows binding between the exact frozen source and the current accepted local plane without importing runtime state.

This remains an intake and blocked-seam inventory until that exact-source Windows closeout is recorded. No adapter has been activated or connected to a live plane.

## Contracts inventoried

Reviewed `integration-seams.md`/json, `architecture.md`, `worker-protocol.md` and `known-gaps.md` at the Wingless baseline. The existing 11 integration seams remain BLOCKED for final comparison until the frozen-source acceptance link is complete and source reconciliation is performed. Existing descriptions remain provisional, including local worker terminal labels and resource policy.

The machine-readable inventory preserves those descriptions and adds status, source-reference availability, disposition, translation uncertainty and replay-test status. No `CONFIRMED_MATCH` is claimed yet. The presence of frozen source removes the source-availability blocker but does not itself authorize assumptions about runtime ownership.

## Required ownership and duplicate disposition

These are operator requirements to validate against the frozen implementation:

| Concern | Required owner | Wingless disposition |
| --- | --- | --- |
| Work orders, queue, dispatch, leases, duplicate suppression and crash reconciliation | Plane | PLANE_OWNED; no new implementation |
| Materialization, workspace identity/revision, lifetime and cleanup | Plane | PLANE_OWNED; consume verified identity only |
| Work-order retry count/timing, repair dispatch and terminal failure | Plane | PLANE_OWNED |
| Existing local mock worker IDs, records and cancellation | Internal mock only | ADAPTER_ONLY; never plane lifecycle |
| Fixture verifier and bounded inference repairs | Internal research | KEEP_INTERNAL; no adapter retry budget until explicitly mapped |
| ICE/context, routing, inference and candidate telemetry | Wingless intelligence layer | KEEP_INTERNAL |
| Read primitive and future structured tool proposals | Adapter boundary | ADAPTER_ONLY; no write/process authority |
| Acceptance, completion and authoritative evidence export | Plane | PLANE_OWNED; candidate remains external_required |

No useful internal code is to be deleted based on unverified plane assumptions.

## Resource ownership questions

| Concern | Plane | Wingless | Shared adapter requirement |
| --- | --- | --- | --- |
| RAM / VRAM / disk floors | Inspect frozen budget/reservation behavior | Existing measurements and local admission checks | Translate floors without implying reservations |
| CPU ceiling | Inspect frozen scheduling ownership | Existing nullable sample and ceiling check | Preserve unknown-value refusal |
| Model residency | Inspect plane process/worker ownership boundary | Experimental supervisor | Establish one lifetime owner |
| Model process lifecycle | Inspect controller/worker cancellation and containment | Owns configured independent experiment child | Prevent competing termination owners |
| GPU allocation | Inspect whether plane expresses allocation/reservation at all | Device selection and measurements | Measurement is not an allocation lease |
| Timeout | Plane work-order deadline ownership required | Bounded inference/runtime deadlines | Determine remaining-budget rule |
| Cancellation | Plane work-order cancel ownership required | Cooperative inference cancellation and supervisor stop | Define signal, grace, forced cleanup and result semantics |

## Source reconciliation plan

After exact-source Windows closeout is linked:

1. Inspect the frozen plane worker, workspace, cancellation, retry, resources, acceptance and evidence definitions through narrow source reads from `539d56fec273c8851aea131cf4d31425a62f7250` only.
2. Replace each BLOCKED classification with the supported classification, retaining precise plane and Wingless source references. Record unresolved mismatches explicitly.
3. Add a non-default inactive adapter only against those real definitions. Reject unsupported fields/states. Translate identity, deadline, authorized fallback and bounded result evidence without broadening authority.
4. Bind context/ICE to actual authorized workspace identity plus revision; reject stale indexes or mismatches.
5. Add mock/replay tests for request translation, cancellation, stale context, resource policy, backend failure, authorized fast/deep fallback, unknown fields/states, candidate-only results and no duplicate work-order retry.
6. Record test results and remaining activation prerequisites. No live queue or worker is involved.

No live plane tests are authorized by this document. Existing Wingless Stage 3 tests do not count as plane adapter tests.

## Remaining intake gate

The frozen source requirement is satisfied. The remaining intake gate is a sanitized acceptance record binding the exact frozen source identity to authoritative Windows verification.

The ckb-plane closeout verifier must establish all of the following without live-state import or accepted-ref mutation:

- frozen manifest hashes verify;
- current local contract source exactly matches the frozen source set;
- post-DT focused regressions pass on the frozen source;
- the full `internal/plane` suite passes on the frozen source;
- the frozen `cmd/ckb-plane` builds on Windows;
- plane authority remains `plane`;
- controller remains paused with a fresh heartbeat;
- accepted-ref mutation remains disabled;
- the original GitHub runner remains stopped.

Until that record is linked, the final source-to-acceptance identity is not claimed.

## Separate qualification gates

Native Windows execution, NVIDIA collection, CUDA inference/residency, Windows Job Object runtime containment, process GPU attribution and energy telemetry remain separate as documented in Stage 3. Contract reconciliation cannot close them. No Stage 4 model research or live activation is authorized by this report.

NOT_READY_FOR_ISOLATED_INTEGRATION_PROOF
