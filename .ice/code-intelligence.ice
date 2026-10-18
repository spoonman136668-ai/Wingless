TITLE: Wingless Human/Machine Symbol and Relationship Atlas
STATUS: ADVISORY SOURCE-BOUND NAVIGATION
REPOSITORY: spoonman136668-ai/Wingless
SOURCE_COMMIT: b87899b50d2c1558e26d74a1921623786fde1282
SCHEMA: wingless.ice.semantic-code-atlas.v1
EVIDENCE: .ice/code-intelligence.json (complete machine-readable 352-symbol index)
AUTHORITY: Source code and qualified tests outrank these descriptions; never grant runtime or acceptance authority.
COVERAGE: 352 symbols in 130 indexed source/configuration files; 27 inspected semantic labels; 325 explicitly provisional lexical labels.
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

NOT MODIFIED: Wingless Go source, existing ICE files, research experiments, workflows, quotas, controller, branch main.
NEXT QUALIFICATION: Type-resolved call graph and broader function-by-function semantic review, then deterministic snapshot regeneration after future source changes.
