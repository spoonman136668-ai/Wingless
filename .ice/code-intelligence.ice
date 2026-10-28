TITLE: Wingless Human/Machine Symbol and Relationship Atlas
STATUS: ADVISORY SOURCE-BOUND NAVIGATION
REPOSITORY: spoonman136668-ai/Wingless
SOURCE_COMMIT: 950dab6eb5e5079034f8196e97e1ef9c03187053
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

UNIFORM SYMBOL-LEVEL PROVENANCE CLOSEOUT — 2026-10-18
All 352/352 symbols now carry exact pinned source commit, source-file SHA-256, declaration line ranges, and literal within-declaration evidence tokens. The 27 historically source-reviewed descriptions were independently rechecked and brought into the same proof contract as the remaining 325.
All 352 symbols now record their SOURCE-FILE's local Wingless import dependencies from the frozen ICE dependency index. A file-level import is not evidence of a direct function call or runtime transition.
Original 9 source-reviewed typed call edges, 51 indexed package imports and 95 provisional syntax-call candidates retain their separate evidence classifications.
No application behavior, scientific experiment, runtime controller, or other repository is modified by this documentation/provenance update.

TAMPER-RESISTANCE NEGATIVE QUALIFICATION — 2026-10-18
The qualified Windows workflow now runs a separate fail-closed negative suite in a disposable clone of the exact PR head. Eight synthetic attacks are required to be rejected: source-hash drift, forged review anchors, forged local dependency, forged test classification, changed symbol line span, changed offline viewer data, erased type-resolved call evidence, and stale canonical ICE source. Original checkout and source remain untouched. A source anchor supports provenance, not a machine proof of human-written behavioral semantics; independent reviewer approval remains a separate gate.

INDEPENDENT REVIEW P1 — SOURCE-ONLY PR FRESHNESS (2026-10-18)
Review found that the prior workflow trigger observed atlas changes but could miss changes to indexed source on another pull request. The PR workflow now watches every source type accepted by ice.Build (.go/.md/.json/.ps1/.yml/.yaml, go.mod, .gitattributes and .gitignore) and both standalone validators. A read-only contract checks every file in the frozen 138-file ICE inventory and representative future new files against trigger patterns before performing the eight negative mutations. Source-only PRs must not bypass ICE freshness validation; the existing fail-closed no-merge decision is preserved. Broad patterns may schedule more validation runs, but the workflow remains read-only and does not execute scientific research.

G26 SOURCE-AWARE ATLAS MAINTENANCE
TARGET_SOURCE_COMMIT: d312f2ada99cce4d994f271bbed54cbe7d52027e
PREVIOUS_SOURCE_COMMIT: b87899b50d2c1558e26d74a1921623786fde1282
EXACT_DECLARATION_REUSE_COUNT: 352
FRESH_RECEIPT_COUNT: 0
FRESH_RECEIPTS_ARE_SELF_ASSERTED: TRUE
INDEPENDENT_HUMAN_REVIEW_REQUIRED_FOR_PROMOTION: TRUE
Historical evidence is retained as lineage, never silently renewed for changed declarations. Syntax-only edges remain provisional.

G26 SOURCE-AWARE ATLAS MAINTENANCE
TARGET_SOURCE_COMMIT: f49548c28c3c2aa8434c93b8a26546b8dcba4db5
PREVIOUS_SOURCE_COMMIT: d312f2ada99cce4d994f271bbed54cbe7d52027e
EXACT_DECLARATION_REUSE_COUNT: 352
FRESH_RECEIPT_COUNT: 0
FRESH_RECEIPTS_ARE_SELF_ASSERTED: TRUE
INDEPENDENT_HUMAN_REVIEW_REQUIRED_FOR_PROMOTION: TRUE
Historical evidence is retained as lineage, never silently renewed for changed declarations. Syntax-only edges remain provisional.

G26 SOURCE-AWARE ATLAS MAINTENANCE
TARGET_SOURCE_COMMIT: c206a9f4b3bd452d5ab1be293f361f8a704858e9
PREVIOUS_SOURCE_COMMIT: f49548c28c3c2aa8434c93b8a26546b8dcba4db5
EXACT_DECLARATION_REUSE_COUNT: 352
FRESH_RECEIPT_COUNT: 0
FRESH_RECEIPTS_ARE_SELF_ASSERTED: TRUE
INDEPENDENT_HUMAN_REVIEW_REQUIRED_FOR_PROMOTION: TRUE
Historical evidence is retained as lineage, never silently renewed for changed declarations. Syntax-only edges remain provisional.

G26 SOURCE-AWARE ATLAS MAINTENANCE
TARGET_SOURCE_COMMIT: a6f825ad6420e769db8fadee560fd2a379b94ce9
PREVIOUS_SOURCE_COMMIT: c206a9f4b3bd452d5ab1be293f361f8a704858e9
EXACT_DECLARATION_REUSE_COUNT: 352
FRESH_RECEIPT_COUNT: 0
FRESH_RECEIPTS_ARE_SELF_ASSERTED: TRUE
INDEPENDENT_HUMAN_REVIEW_REQUIRED_FOR_PROMOTION: TRUE
Historical evidence is retained as lineage, never silently renewed for changed declarations. Syntax-only edges remain provisional.

PROVENANCE REANCHOR — 2026-10-20
PREVIOUS_SOURCE_COMMIT: a6f825ad6420e769db8fadee560fd2a379b94ce9
PROTECTED_MAIN_SOURCE_COMMIT: fd00165e92b4fc21d93d33292c32038e6e96bb22
INDEXED_SOURCE_BLOB_IDENTITY: 138_OF_138_EQUAL
REVIEWED_DECLARATIONS_REUSED: 352_OF_352
HISTORICAL_PROVENANCE: retained per-symbol in code-intelligence.json
SCIENTIFIC_ACCEPTANCE: NONE; advisory navigation only

G26 SOURCE-AWARE ATLAS MAINTENANCE
TARGET_SOURCE_COMMIT: de8b41ffe157a79b82193661dcab4f53bc1d7204
PREVIOUS_SOURCE_COMMIT: fd00165e92b4fc21d93d33292c32038e6e96bb22
EXACT_DECLARATION_REUSE_COUNT: 352
FRESH_RECEIPT_COUNT: 0
FRESH_RECEIPTS_ARE_SELF_ASSERTED: TRUE
INDEPENDENT_HUMAN_REVIEW_REQUIRED_FOR_PROMOTION: TRUE
Historical evidence is retained as lineage, never silently renewed for changed declarations. Syntax-only edges remain provisional.

G26 SOURCE-AWARE ATLAS MAINTENANCE
TARGET_SOURCE_COMMIT: 950dab6eb5e5079034f8196e97e1ef9c03187053
PREVIOUS_SOURCE_COMMIT: de8b41ffe157a79b82193661dcab4f53bc1d7204
EXACT_DECLARATION_REUSE_COUNT: 352
FRESH_RECEIPT_COUNT: 0
FRESH_RECEIPTS_ARE_SELF_ASSERTED: TRUE
INDEPENDENT_HUMAN_REVIEW_REQUIRED_FOR_PROMOTION: TRUE
Historical evidence is retained as lineage, never silently renewed for changed declarations. Syntax-only edges remain provisional.
