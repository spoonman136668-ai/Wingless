# Nemotron experiment reasoner

Status: **qualified, not activated** on `research/nemotron-openrouter-reasoner-r1`. It is not wired into ckb-plane, KTRADE, accepted refs, queues, or any automatic research execution path.

## Purpose

Use NVIDIA Nemotron 3 Ultra through OpenRouter as an advisory scientific reasoning service for:

- planning the next preregistered experiment,
- critiquing a proposed experiment before freeze,
- interpreting a sealed result after deterministic qualification.

Model output never grants execution, acceptance, promotion, retry, threshold-change, or scope-expansion authority.

## Safety and scientific boundary

The caller must freeze the experiment contract before execution. A reasoner response may influence only a future preregistration or an interpretation record. It must never mutate the criteria of an experiment already in flight or already observed.

The adapter fails closed on:

- missing request identity,
- missing immutable frontier or qualification-contract SHA-256 identities,
- unknown reasoning role,
- unbounded context/output/time,
- endpoint drift,
- missing API key,
- any request not classified as public-repository,
- model-route drift away from the pinned free endpoint,
- model identity drift,
- incomplete generation,
- oversized or malformed responses.

No automatic retries are performed.

## Default model and routing

Default model:

`nvidia/nemotron-3-ultra-550b-a55b:free`

The free model route is the only staged route for this qualification. Every transmitted request must be classified `public-repository` and derived only from public Wingless or Yggdrasil repository material. Mind-Palace records, local-only evidence, credentials, personal data, and unpublished/private context are forbidden from this route.

The free NVIDIA endpoint is logged by NVIDIA for security and product-improvement purposes. That is acceptable for this staging path only because the transmitted scientific context is already public. Provider fallbacks remain disabled and supported request parameters are required.

## Activation ceremony

Activation is intentionally separate from setup.

Required environment variables on a dedicated research runner:

- `WINGLESS_REMOTE_REASONER_ENABLE=1`
- `WINGLESS_REASONER_OPENROUTER_API_KEY=<secret from runner secret store>`

The research API key must be distinct from any CKB repair/orchestration OpenRouter credential and must never be committed, printed, embedded in evidence, or passed through ckb-plane work-order text.

Example invocation after qualification and explicit activation authority:

`go run ./cmd/research-reasoner -request <request.json> -out <result.json>`

## Provenance

Every result records:

- requested and returned model identities,
- returned provider identity when supplied,
- seed and reasoning effort,
- frontier SHA-256,
- qualification-contract SHA-256,
- request SHA-256,
- configuration SHA-256,
- response SHA-256,
- a composed qualification-identity SHA-256,
- token usage and latency.

Provider/model/config identity changes therefore produce different evidence identity. Cached evidence must not be reused across identity drift.

## Qualification result

Qualification run `36483722629` at source head `212ca3031977da808a1416065599f65992f3f17f` passed:

- focused reasoner tests: PASS;
- full Wingless regression: PASS;
- seven blind historical decision fixtures across Wingless and Yggdrasil;
- two identical repetitions per fixture;
- 14/14 historical next-experiment selections correct;
- 0 scientific-boundary violations;
- replay artifact `10998885130`;
- artifact digest `sha256:19f5943efde3c2169a0b3d3fa1dc8def581dce41afa9e3266f2e335887471073`.

This qualifies the pinned free Nemotron route for bounded advisory reasoning over public-repository evidence. It does not grant execution or acceptance authority.

## Activation before live research

Before the remote reasoner is allowed into a live research loop:

1. Restore and verify the Mind-Palace mailbox/durable-context retrieval path.
2. Reconstruct every Nemotron packet exclusively from public repository evidence selected with that context.
3. Preserve the existing deterministic preregistration and experiment-execution authority boundaries.
4. Integrate through the existing research boundary without altering KTRADE scheduling or accepted refs.
5. Requalify if model, provider route, prompt contract, fixture contract, or relevant adapter identity changes. **Still pending intentionally.**

KTRADE and current ckb-plane production work must not be paused or modified for this qualification.


## Mind-Palace boundary

Mind-Palace remains the durable-context source for historical project context, but its retrieved records are not transmitted to the free Nemotron endpoint. When Mind-Palace is healthy, it may help identify which public repository records should be fetched. The actual reasoner packet must then be reconstructed exclusively from public repository content and bound to its public-source hashes.

This preserves the architecture:

Mind-Palace durable context -> select relevant public evidence -> public-repository reasoner packet -> Nemotron advisory reasoning -> deterministic preregistration/execution authority -> sealed public result -> durable context ingestion.


## Reduced-context Mind-Palace integration

Mailbox recovery was reverified end-to-end on 2026-09-28: a fresh ping succeeded, the bridge returned a current authority generation, and one generation-bound `mind-palace-context` local-exec completed successfully under the existing advisory-only contract.

The token-reduction path is deliberately split into two layers:

1. Mind-Palace performs durable historical retrieval and selects relevant durable record IDs.
2. Private Mind-Palace statement text is not forwarded to the free Nemotron route.
3. An ID-only sidecar resolves selected records to public experiment identifiers such as `UP-*`, `YGG-*`, or Yggdrasil lineage IDs.
4. Wingless reconstructs a canonical context pack from public Wingless/Yggdrasil experiment digests only.
5. The pack is content-addressed and bound into `FrontierSHA256`.
6. The pack admits at most 8 experiment digests and 65,536 encoded bytes. If the evidence does not fit, selection must narrow; scientific claims are never silently truncated.
7. Nemo receives that bounded public pack, not the full research history and not Mind-Palace private prose.

Implemented on this branch:

- `reasoner/contextpack.go`: canonical sealed-result digest and bounded context-pack contract;
- `reasoner/mindpalace.go`: parses the existing Mind-Palace context envelope, consumes durable IDs, and resolves them through a public-evidence index without using `model_context` as prompt material;
- deterministic tests for ordering, provenance identity, byte bounds, private-source rejection, unavailable-selector failure, and request binding.

The private Mind-Palace sidecar is staged independently on `spoonman136668-ai/Mind-Palace:research/nemotron-public-context-r1`. It extracts experiment IDs and approved public source references from selected durable records while emitting no statements or private model context.

### Current activation blocker

The mailbox itself is no longer the blocker. Live activation remains held until:

- the Mind-Palace public-evidence sidecar receives authoritative Windows qualification;
- its ckb-plane local-exec exposure is qualified without changing the existing `mind-palace-context` operation;
- an end-to-end dry run proves Mind-Palace IDs -> public experiment evidence -> reduced context pack -> advisory Nemotron result;
- final integration/full regression is green under exact immutable identities.

Private-repo hosted Actions are presently failing before runner step 1, so that failure is treated as infrastructure evidence, not as a scientific or sidecar-code result.
