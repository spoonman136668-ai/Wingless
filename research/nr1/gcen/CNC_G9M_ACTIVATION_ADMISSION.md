# CNC-G9M — original-Qwen activation evidence and runtime admission boundary

**Reason:** Earlier G9E/F/I/L real-weight comparisons authenticate original BF16
weight windows, but inputs are synthetic Gaussian probes. The existing G21
`activation_replay.py` checks user-supplied trace declarations; it does not
capture or independently authenticate real original-model forward activations.
Existing G25 tiny-model disk-offload smoke is not a full Qwen3-30B-A3B proof.
Do **not** misrepresent either.

G9M adds a bounded, offline, **read-only** capture integrity preflight for
a future original-forward trace bundle. It requires pinned exact Qwen source
checkpoint revision, a 40-hex source SHA, workflow run ID, one named expert
tensor, and separately identified calibration/evaluation records. Two
distinct on-disk trace files must be local, non-symlink, each <=1GiB, with
matching SHA-256 bytes. It refuses overlapped records, protected heldout,
malformed scope/identities, and path traversal. This does not read private
raw tensors into repository artifacts or make any network/model calls.

Critically, a self-declared `native_model_forward_hook` flag and correct
checksums do **not** prove original-Qwen provenance. Every positive custody
result retains `original_forward_authenticity=UNVERIFIED_REQUIRES_INDEPENDENT_ATTESTATION`,
`source_model_activation_evaluation_admitted=false`, and
`scientific_acceptance=NOT_GRANTED`. True authenticity requires separately
authorized full-model forward-hook capture with independent runner, exact
binary/model/shard identities, record custody and authentic source-log
attestation; this module never grants it.

Before any authentic task quality or full-model 12-GiB claims:
1. Independently verify an authorized execution path with the required
   original model and runtime checkpoint (NOT G21 or tiny G25).
2. Capture disjoint development calibration and evaluation input IDs,
   original expert forward activations, source SHA, immutable workflow
   run, original tensor identity, and bounded trace byte checksums.
3. Verify original forward-hook logs and provenance against an independently
   retained receipt; keep raw prompts/tokens/activations private.
4. Measure GPU peak, workspace, tokens/context, latency and task-quality
   parity against frozen reference under independently authorized execution.
5. Admit heldout and scientific claims only through separate acceptance.

No GPU activation, full model download, dependency install, new controller,
protected-main change or promotion is authorized by G9M. G9J rejection of
direct uniform sparse corrections and G9L low-rank limited gain remain frozen.
