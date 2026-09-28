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
