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
- disabled ZDR,
- allowed provider data collection,
- model identity drift,
- incomplete generation,
- oversized or malformed responses.

No automatic retries are performed.

## Default model and routing

Default model:

`nvidia/nemotron-3-ultra-550b-a55b`

The paid model route is the default because actual research frontier packets may be proprietary. The request pins `data_collection=deny`, `zdr=true`, `require_parameters=true`, and provider fallbacks off.

The free model can be selected explicitly with:

`WINGLESS_REASONER_MODEL=nvidia/nemotron-3-ultra-550b-a55b:free`

Only use the free route with a sanitized/non-confidential context packet. Do not weaken the request privacy controls in code merely to make a provider eligible.

## Activation ceremony

Activation is intentionally separate from setup.

Required environment variables on a dedicated research runner:

- `WINGLESS_REMOTE_REASONER_ENABLE=1`
- `OPENROUTER_API_KEY=<secret from runner secret store>`

Optional model override:

- `WINGLESS_REASONER_MODEL=nvidia/nemotron-3-ultra-550b-a55b:free`

The API key must never be committed, printed, embedded in evidence, or passed through ckb-plane work-order text.

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
