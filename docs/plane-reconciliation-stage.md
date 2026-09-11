# Plane reconciliation — Stage 3.5 intake

State: **NOT_READY_FOR_ISOLATED_INTEGRATION_PROOF**

## Frozen inputs

Wingless repository: spoonman136668-ai/Wingless, branch main.
Reviewed documentation baseline: `584c92a04b9a9944fee491cc59a4025ea4143ffc`.
Review date: 2026-09-11.

Plane source location, branch, commit, snapshot timestamp and linked acceptance evidence: **not supplied for this stage**. No frozen stabilized plane implementation has been inspected. The handoff describes intended ownership; it does not establish actual request fields or runtime behavior. Historical pre-DO source and DN/DT acceptance summaries cannot substitute for the stabilized source.

This is an intake and blocked-seam inventory, not completed reconciliation. No adapter has been created.

## Contracts inventoried

Reviewed integration-seams.md/json, architecture.md, worker-protocol.md and known-gaps.md at the Wingless baseline. All 11 existing integration seams are classified BLOCKED for comparison with the actual plane. Existing descriptions remain provisional, including local worker terminal labels and resource policy.

The machine-readable inventory preserves those descriptions and adds status, source-reference availability, disposition, translation uncertainty and replay-test status. No CONFIRMED_MATCH is claimed. Documentation references are not a new code audit.

## Required ownership and duplicate disposition

These are requirements from the operator handoff, pending implementation verification:

| Concern | Required owner | Wingless disposition |
| --- | --- | --- |
| Work orders, queue, dispatch, leases if any, duplicate suppression and crash reconciliation | Plane | PLANE_OWNED; no new implementation |
| Materialization, workspace identity/revision, lifetime and cleanup | Plane | PLANE_OWNED; consume verified identity only |
| Work-order retry count/timing, repair dispatch and terminal failure | Plane | PLANE_OWNED |
| Existing local mock worker IDs, records and cancellation | Internal mock only | ADAPTER_ONLY; never plane lifecycle |
| Fixture verifier and bounded inference repairs | Internal research | KEEP_INTERNAL; no adapter retry budget until explicitly mapped |
| ICE/context, routing, inference and candidate telemetry | Wingless intelligence layer | KEEP_INTERNAL |
| Read primitive and future structured tool proposals | Adapter boundary | ADAPTER_ONLY; no write/process authority |
| Acceptance, completion and authoritative evidence export | Plane | PLANE_OWNED; candidate remains external_required |

No useful internal code was deleted based on unverified plane assumptions.

## Resource ownership questions

| Concern | Plane | Wingless | Shared adapter requirement |
| --- | --- | --- | --- |
| RAM / VRAM / disk floors | Budget/reservation contract unknown | Existing measurements and local admission checks | Translate floors without implying reservations |
| CPU ceiling | Scheduling ownership unknown | Existing nullable sample and ceiling check | Preserve unknown-value refusal |
| Model residency | Runtime policy unknown | Experimental supervisor | Establish one lifetime owner |
| Model process lifecycle | Worker containment contract unknown | Owns configured independent experiment child | Prevent competing termination owners |
| GPU allocation | Reservation contract unknown | Device selection and measurements | Measurement is not an allocation lease |
| Timeout | Work-order deadline ownership required | Bounded inference/runtime deadlines | Determine remaining-budget rule |
| Cancellation | Work-order cancel ownership required | Cooperative inference cancellation and supervisor stop | Define signal, grace, forced cleanup and result semantics |

## Adapter and replay patch plan after source arrives

1. Freeze plane commit or hash-identified archive and acceptance references. Inspect worker, workspace, cancellation, retry, resources, acceptance and evidence definitions through narrow source reads.
2. Replace each BLOCKED classification with the supported classification, retaining precise plane and Wingless source references. Record unresolved mismatches explicitly.
3. Add a non-default inactive adapter only against those real definitions. Reject unsupported fields/states. Translate identity, deadline, authorized fallback and bounded result evidence without broadening authority.
4. Bind context/ICE to actual authorized workspace identity plus revision; reject stale indexes or mismatches.
5. Add mock/replay tests for request translation, cancellation, stale context, resource policy, backend failure, authorized fast/deep fallback, unknown fields/states, candidate-only results and no duplicate work-order retry.
6. Record test results and remaining activation prerequisites. No live queue or worker is involved.

No Go tests were added or rerun for this documentation-only intake. Existing Stage 3 tests do not count as plane adapter tests.

## Needed source bundle

Provide an existing stabilized plane source ZIP, or a repository URL with immutable commit SHA, plus the corresponding DN/DO/DP/DT FixA acceptance evidence identifying that source. Include worker interfaces and implementations, controller/ownership code, workspace/materialization, cancellation, retry/reconciliation, result/evidence and acceptance definitions, focused tests and sanitized example envelopes.

Do not include credentials, runtime state, live queue contents or unrelated CKB product source. A repository snapshot already available remotely requires no local-machine execution. If source identity cannot be linked to the acceptance evidence, reconciliation remains blocked.

## Separate qualification gates

Native Windows execution, NVIDIA collection, CUDA inference/residency, Windows Job Object runtime containment, process GPU attribution and energy telemetry remain open as documented in Stage 3. Contract reconciliation cannot close them. No Stage 4 model research or live activation is authorized by this report.

NOT_READY_FOR_ISOLATED_INTEGRATION_PROOF
