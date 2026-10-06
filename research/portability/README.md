# Portable Checkpoint R1

This directory defines a dormant, language-neutral checkpoint/export contract for accepted learned mechanisms. It is infrastructure only: it does not attach to the research controller, change training, alter scientific thresholds, write accepted refs, promote artifacts, launch runtimes, or create execution authority.

## Purpose

A learned mechanism should be movable to a future runtime or implementation language without retraining. The portable artifact therefore bundles opaque learned parameters, resumable state, deterministic replay input/output vectors, lineage/provenance, resource identity, and exact qualification identity.

Wingless checkpoints must preserve the exact `state_former(raw_input, history) -> fixed_size_state` interface, fixed state dimension, fixed readout capacity, and no capacity growth.

## Qualification identity

Every export binds:

- source commit SHA;
- source tree SHA;
- source/schema SHA256;
- qualification-contract SHA256;
- toolchain identity;
- experiment lineage and optional parent checkpoint;
- fixed resource envelope and baseline identity;
- project-specific interface/capacity contract.

This is the semantic identity used to decide whether evidence can be reused. Identity drift fails closed.

## Artifact boundary

The checkpoint format is intentionally language-neutral. Binary learned parameters and resumable state are copied as opaque bytes and SHA256-bound in a JSON manifest. A future Go or other runtime may restore the same bytes without retraining.

Export and restore alone never prove runtime equivalence. Promotion requires:

`export -> verify -> restore -> deterministic runtime replay -> affected-interface qualification -> one full regression -> build -> promote`

The final runtime replay and requalification remain mandatory.

## Authority

Portable checkpoints are evidence/artifacts only. The manifest and restore receipt explicitly carry:

- no execution authority;
- no accepted-state mutation;
- no successor authority;
- no production authority;
- no promotion authority;
- no scientific classification;
- no retraining during restore.

CKB-plane remains the sole execution/orchestration and promotion authority.

## Files

- `scripts/portable_checkpoint.py`: fail-closed export/verify/restore implementation.
- `tests/test_portable_checkpoint.py`: qualification tests for identity, tamper resistance, restore fidelity, replay binding, resource identity, and project contract.
- `research/portability/export.example.json`: non-authoritative example export specification.

This R1 framework remains dormant until a learned mechanism reaches an accepted, stable checkpoint boundary.
