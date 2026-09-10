# Inference backends

The contract is ID, Capabilities, Health, EstimateCost, Invoke and Cancel. Implementations must honor parent cancellation and request deadline. Registration is deterministic by backend ID; the caller allowlist is mandatory. Mock provides immutable fixture configuration and supports cancellation and failure injection.

`NewLocalHTTP` accepts only a numeric loopback HTTP origin, e.g. http://127.0.0.1:8080. It never sets credentials, follows redirects or uses environment proxies. `/v1/models` must report the configured model. `/v1/chat/completions` receives a single bounded user message, max_tokens, stream=false. The response must contain exactly one choice for that model with finish_reason=stop. Token usage is nullable; provided output counts must fit the budget. Response bytes are bounded separately. This is a deliberately small compatibility subset; no streaming, tool calls or authentication support.

HTTP 401/403 => auth_failed; 429 => rate_limited; 404/503 => model_unavailable; deadline => worker_timeout; cancellation => canceled. Other errors => backend_failure. Raw HTTP error bodies are not retained. No automatic backend-error retry occurs. Cost currently exposes input bytes and requested output tokens, not fabricated token estimates.

Codex adapter, real fast/deep model runtime, model downloads, remote endpoints and sparse MoE execution are unavailable. Capabilities for an endpoint are trusted operator configuration, not inferred from a model claim. Lifecycle tests cover legal transitions; attaching process start/warm/unload/restart remains future work.
