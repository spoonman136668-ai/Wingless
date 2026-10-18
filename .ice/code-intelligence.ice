TITLE: Wingless Human/Machine Symbol and Relationship Atlas
STATUS: ADVISORY SOURCE-BOUND NAVIGATION
REPOSITORY: spoonman136668-ai/Wingless
SOURCE_COMMIT: b87899b50d2c1558e26d74a1921623786fde1282
SCHEMA: wingless.ice.semantic-code-atlas.v1
EVIDENCE: .ice/code-intelligence.json (complete machine-readable 352-symbol index)
AUTHORITY: Source code and qualified tests outrank these descriptions; never grant runtime or acceptance authority.
COVERAGE: 352 symbols in 138 indexed source/configuration files; 352 inspected semantic labels; 0 explicitly provisional lexical labels.
RELATIONSHIPS: 9 inspected relationships, 51 local package imports, 95 unresolved syntax-only call candidates.

HOW TO USE
Search code-intelligence.json for a human-readable label, file, symbol, or state.
Every symbol includes a file content hash, exact source lines and a pinned-commit source link.
Source-reviewed labels explain behavior observed in the pinned source.
Lexical-provisional labels are navigational aids, not claims of verified behavior.
Syntactic call candidates are not type-checked or proven runtime paths. Do not use them to infer side effects.

SOURCE-REVIEWED WORKFLOW A — Bounded inference
1. worker.Service.SubmitWork validates and admits one in-memory request before starting Runner.Run.
2. broker.Runner.Run validates request/resources and invokes broker.Select.
3. broker.Select enforces allowlist, tier, capabilities, health and resources.Check.
4. The selected InferenceBackend produces a candidate. Optional deterministic verifier limits bounded repair.
5. worker.Service.Status, Cancel, and Result observe, request cancellation, or copy a terminal outcome; no durable queue is created.

SOURCE-REVIEWED WORKFLOW B — Validated contextual retrieval
1. contextbroker.ICE.Retrieve invokes ice.Query, which requires a fresh source-derived index.
2. The provider retrieves exact Go symbol line ranges, enforces a byte budget and revalidates freshness.
3. contextbroker.Collect applies a second provider-agnostic budget cap.

SOURCE-REVIEWED WORKFLOW C — Inactive external translation
1. integration/ckbplane.Translate validates identity, workspace SHA, protected-work flags and trusted routing.
2. It produces one bounded inference.Request with internal repair attempts fixed to zero.
3. CandidateFromOutcome projects only non-authoritative terminal candidate data and rejects attempts to claim acceptance.
4. No live controller, queue, dispatch, or acceptance path is activated by this atlas.

SOURCE-REVIEWED STATE MACHINE — inference.Lifecycle.Move
- STOPPED -> STARTING, UNAVAILABLE
- UNAVAILABLE -> STOPPED
- STARTING -> READY, FAILED
- READY -> BUSY, DRAINING, FAILED
- BUSY -> IDLE, FAILED
- IDLE -> BUSY, DRAINING, FAILED
- DRAINING -> STOPPED, FAILED
- FAILED -> STOPPED
These are explicit local transitions only; an external controller has a separate state machine.

NOT MODIFIED: Wingless application and scientific Go source, research experiments, existing runtime workflows/quotas, controller, or main branch.
NEXT QUALIFICATION: Type-resolved call graph and broader function-by-function semantic review, then deterministic snapshot regeneration after future source changes.

SOURCE-BODY EVIDENCE EXPANSION — 2026-10-18
New source-reviewed symbols: 77
Source review is anchored to frozen b87899b50d2c1558e26d74a1921623786fde1282; exact line ranges, source SHA-256 and literal body/field evidence are attached to each newly reviewed JSON node.
Reviewed functions and interfaces cover tool authorization, context limits, local inference transport, resource admission, model process supervision, type-resolved relationship extraction, strict candidate transport and bounded synthetic replay.
Provisional symbols remain explicitly unverified; no behavioral/source execution authority is implied.

SOURCE-BODY EVIDENCE EXPANSION — TRANCHE 2 (2026-10-18)
Additional source-reviewed symbols: 78; cumulative source-reviewed: 182 / 352; still provisional: 170.
New coverage includes: core inference request/result schema; scoped worker/broker records; local lifecycle states; inactive external work translation; transport guard constants; GPU collection/CSV validation; exact microfixture benchmarking; typed ICE index metadata; and CLI fixture modes.
Each of these 78 descriptions is anchored to an immutable source commit, matching file SHA-256 and exact declaration lines, plus source literal evidence verified during qualification.
Each newly reviewed symbol also records the file's imported local Wingless packages. These are file-level dependency facts, NOT evidence that the particular function calls those packages.
All preexisting 9 confirmed typed call relationships remain unchanged. Syntax-only relationship candidates remain provisional.
No execution, experimental authority, data mutation, deployment, production activation, or controller scheduling occurs as part of this documentary expansion.

SOURCE-BODY EVIDENCE EXPANSION — TRANCHE 3 (2026-10-18)
Added 54 source-reviewed production symbols, retaining exact source commit, SHA-256, definition lines, and validated literal source anchors.
Cumulative coverage: 236 / 352 source-reviewed; remaining 116 provisional declarations are test functions or test helpers only.
Source-reviewed coverage now includes all currently indexed non-test declarations across benchmark fixtures/provenance, diagnostic CLIs, platform-specific resource sampling, Windows/non-Windows process containment, and the isolated transport proof executable.
The 54 new source-bound entries also record source-derived file-level local imports, explicitly not treated as proof of direct runtime calls.
Existing scientific execution, application source, and protected acceptance/research workflow authority remain unchanged. Existing 9 confirmed typed relationship edges and 95 unverified syntactic candidates remain separately labeled.

SOURCE-REVIEWED TEST CONTRACT COMPLETION — 2026-10-18
Test source declarations and fixtures reviewed: 116 (66 test functions, 50 fixtures/helpers).
Cumulative descriptions reviewed: 352 of 352 on the frozen index; 0 provisional symbol labels.
New test contracts cover bounded benchmark scoring, routing, transport and workspace identity checks, cancellation, exact model-boundary input validation, source-index anti-tampering, platform-specific host/process sampling, owned model lifecycle and read-only tool boundary.
Every test record includes test-case versus fixture/helper classification and an immutable source-line witness; validation reopens the source file and matches exact proof tokens to the indexed declaration range.
Test descriptions explain existing assertions and negative cases; they are not evidence of behavior outside the test's platforms, mocks or other observed preconditions.
Untyped source-syntax call candidates remain separate from 9 confirmed typed reviewed edges; do not infer runtime reachability from syntax alone.

NEXT: preserve source identity and validate exact-head Windows full regression, then independently review and separately authorize PR merge. Any later source change requires rebuilding/requalifying the source-bound index.
