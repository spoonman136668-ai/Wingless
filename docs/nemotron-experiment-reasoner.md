# Nemotron experiment reasoner

Status: staged only on `research/nemotron-openrouter-reasoner-r1`. It is not wired into ckb-plane, KTRADE, accepted refs, queues, or any automatic research execution path.

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

## Qualification before activation

This branch is intentionally unqualified. Before the remote reasoner is allowed into a research loop:

1. Run unit tests for request bounds, privacy controls, model identity and provenance.
2. Add fixture-only replay tests over historical Wingless/Yggdrasil decision points.
3. Compare Nemotron proposals blind to future results.
4. Measure constraint-violation rate, redundant-experiment rate, hypothesis coverage and discriminating value.
5. Perform final full Wingless regression.
6. Integrate through the existing ckb-plane research boundary only after explicit authority.

KTRADE and current ckb-plane production work must not be paused or modified for this qualification.


## Mind-Palace boundary

Mind-Palace remains the durable-context source for historical project context, but its retrieved records are not transmitted to the free Nemotron endpoint. When Mind-Palace is healthy, it may help identify which public repository records should be fetched. The actual reasoner packet must then be reconstructed exclusively from public repository content and bound to its public-source hashes.

This preserves the architecture:

Mind-Palace durable context -> select relevant public evidence -> public-repository reasoner packet -> Nemotron advisory reasoning -> deterministic preregistration/execution authority -> sealed public result -> durable context ingestion.
