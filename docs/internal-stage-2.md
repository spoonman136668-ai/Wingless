# Wingless internal stage 2

## Implemented

- Native Windows/Linux RAM, free disk and sampled CPU collection. GPU/VRAM, energy and per-process peaks remain null. Host metrics are not container allocations or reservations.
- Bounded numeric-loopback SSE inference with byte/token-usage checks, complete-stream enforcement, cancellation, first-content latency and generation interval. Ordinary JSON responses also reject malformed UTF-8 and oversize health/output data.
- One-shot local llama-server supervisor: hash-pinned executable and GGUF, measured RAM preflight, fixed loopback/CPU-only arguments, bounded startup/runtime, explicit stop and bounded logs. No arbitrary model-proposed processes or file edits.
- Real endpoint and supervised small-model microbenchmarks using fixed public prompts and deterministic JSON/Go AST checks. Generated code is never executed.
- ICE now includes constants, variables, operational PowerShell/YAML and Git attributes. Bodyless declarations no longer panic. Source-first validation remains unchanged.

## Measured result

See [raw measured rows and provenance](experiments/qwen-0.5b-linux-20260910.json).

Qwen2.5-Coder-0.5B-Instruct Q4_K_M (491,400,064 bytes), llama.cpp b10809, CPU-only, two inference threads, 1024 context tokens. One repetition per fixture on this Linux environment:

| Fixture | Strict correctness | Total latency | First content | Reported output tokens | Approximate generation rate |
| --- | --- | --- | --- | --- | --- |
| JSON arithmetic | Fail | 1269 ms | 609 ms | 14 | 23.22 tokens/s |
| Go Add shape | Fail | 1626 ms | 513 ms | 24 | 22.79 tokens/s |

Both outputs included Markdown fences despite the exact output instructions. Do not remove these failures or relax the verifier to turn this run green. Inference startup, streaming, usage collection and owned-process shutdown succeeded. No claim of coding competence or Codex parity follows from these tiny fixtures. Rate uses reported completion tokens divided by the observed first-to-last content interval; it is an approximate transport-observed rate, not a calibrated kernel benchmark. Peak/process/GPU/energy data was not measured.

## Native Windows status

[Hosted CI run 34542006779](https://github.com/spoonman136668-ai/Wingless/actions/runs/34542006779) failed before executing any steps for both Windows and Linux jobs. The connector exposed no usable cause; its annotation endpoint was unavailable. Native Windows acceptance remains blocked, not failed source acceptance and not proven by cross-compilation. No self-hosted CKB runner was started or substituted.

The committed `scripts/Test-Native.ps1` runs tests, ICE validation, host sampling and a native build while retaining evidence. Final verification on 2026-09-11: `go test -buildvcs=false -race ./... -count=1 -timeout 90s` passed all nine tested packages; `GOOS=windows GOARCH=amd64 go build -buildvcs=false -o bin/wingless.exe ./cmd/wingless` passed. Windows execution still requires a functioning independent Windows runner.

## Internal boundaries and remaining work

The supervisor is a trusted local inference experiment, not an execution sandbox. Its directory and linked runtime libraries must be trusted and stable; the executable pin alone does not authenticate dependencies. It kills only its owned direct child; descendant containment, graceful model unload and owner-crash cleanup are not guaranteed. No automatic restart occurs. Create a new supervisor after terminal state. Readiness uses a configured loopback port and model ID; this is not an adversarial localhost identity proof.

Do not add plane queues, leases, retries, acceptance or tools through this API. All plane-facing integration seams remain provisional.
