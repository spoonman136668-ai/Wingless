# Architecture

`worker.Service` provides a provisional, single-concurrency in-memory boundary. `broker.Runner` checks request bounds, obtains trusted resource metrics, chooses an explicitly allowed healthy/capable backend and invokes it within the shared deadline. A trusted optional verifier can request bounded repair; outputs remain candidates.

`broker.Select` prefers fast inference. Caller-classified hard failure, ambiguity, low confidence or explicit deep request selects deep when authorized. An explicit backend ID wins over automatic routing but never bypasses capability, health or resource checks. Deep-to-fast fallback requires explicit policy. `PlanThenImplement` uses deep planning and fast implementation with the plan as untrusted bounded context. It is a separate caller-selected two-phase operation, not an automatic policy override.

`inference.Mock` and `LocalHTTP` implement the backend contract. Backend class names permit later Codex and local model adapters, but registering a class name does not implement an adapter. `Lifecycle` validates state transitions independently; modelhost separately manages bounded pinned CPU-only llama-server experiments. It does not supply plane tool execution authority.

`contextbroker.ICE` retrieves current source ranges through a versioned hash-validated index. `toolboundary` provides a separate read-only authorized tool primitive. Neither is automatically granted to model output. Future model/tool loops must explicitly parse proposals and invoke the tool boundary; there is no generic shell executor.

No live-plane, CKB, broker, credential, deployment, accepted-ref or queue integration exists.

## Standing cognitive North Star

## Destination

Wingless is an **independent, experimentally falsifiable, non-transformer cognitive architecture research program**. Its goal is a functioning autonomous cognitive system eventually meeting or exceeding LLM-class language, reasoning and learning capabilities, ideally at materially lower compute and memory cost. A coding/reasoning worker, local model broker, or inference adapter is a current engineering artifact or application, **not** the project's ultimate definition.

**Strict native token-free invariant:** Wingless' target cognitive architecture must not use token IDs, a tokenizer vocabulary, BPE/subword segmentation, learned next-token prediction, token-keyed/sequence memory, or token-level autoregressive reasoning **anywhere along its learned perception-to-state-to-reasoning-to-action path**. The foundation must also not be a transformer. Its native learned state should be continuous or otherwise natively structured around real observations, dynamics, relationships and action consequences—not renamed linguistic tokens.

Raw language, audio, visual and action signals are permitted as modalities, and communication with software may require a lossless external serializer; that boundary cannot become a tokenizer-dependent cognitive input, internal state, or generation policy. Conventional token-based models may appear only as separately labeled external baselines or bounded research tools/teachers. Their tokenizers, hidden states and inference must not count as Wingless' own cognition. A conventional language-model wrapper, tokenizer removal by renaming embeddings, or token-based model hidden behind a continuous adapter does **not** satisfy the North Star.

Progress should test native persistent state, cumulative memory, online/continual learning, causal and compositional reasoning, generalization to unfamiliar settings, adaptive representation formation, self-monitoring and eventually safely bounded self-improvement. Integrate validated learning mechanisms into the working Wingless architecture, rather than claiming success solely from small synthetic benchmark gains, experiment throughput, model routing, or compression ratios.

## Independent architecture and scientific evidence

Yggdrasil shares the long-term capability target but **must remain a distinct developmental/regenerative architecture**, with independent hypotheses, experiments, resources, and evaluation. Mechanism proposals (including Learned Causal Relational Memory) are *candidates*, not mandated solutions.

Retain strong independently trained and matched-resource conventional/symbolic controls; frozen preregistrations; authentic novel task-family transfer; causal ablations; complete positives/negatives; CPU/memory accounting; one qualified final regression; accepted source/provenance/first-use boundaries. Preserve all prior evidence without retroactively reclassifying outcomes.

## Status and authority

This is a **prospective governance alignment**, not a scientific achievement or authorization for execution, deployment, live broker activity or accepted-ref mutation. The current Go broker and local LLM comparison experiments are not evidence of learned native cognition. Avoid disturbing ongoing qualified research.

Mind-Palace source record: `boundary/wingless-yggdrasil-nontransformer-nontoken-northstar-20261020` on `spoonman136668-ai/Mind-Palace` commit `5e7414966bdcb6e75776ea3422b98422b20c9025` (advisory provenance; operator instruction controls).

**Frozen-identity boundary:** the active branch `research/autonomy-wingless-active` uses `research/bootstrap/wingless-lm-north-star.txt` pinned by `.research-autonomy/config.json`. This document does **not** replace or alter that file, its hash, its controller behavior, or active research authority. Integrating this boundary into the hash-pinned active charter requires a separately qualified source-and-pin transition, never an unqualified direct write to the active branch.
