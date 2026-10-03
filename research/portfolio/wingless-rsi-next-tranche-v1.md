# Wingless Learning-Efficiency / Safety Research Portfolio — Next Tranche v1

**Status:** ACTIVE ADVISORY PORTFOLIO  
**Effective boundary:** applies only to successors after closure of WLM-LM-EXTERNAL-REASONING-MARGIN-CLASS-ATTRIBUTION-R4.  
**Authority effect:** NONE. This file cannot change execution, queue, orchestration, acceptance, promotion, deployment, production, evaluator, credential, network/tool, resource-ceiling, rollback, shutdown, or accepted-ref authority.

## Governing objective

Increasingly efficient learning of unseen future datasets under fixed data, compute, capacity, and authority budgets, with invariant external control, deterministic evidence, complete causal traceability, bounded reversible modification, and reproducible qualification.

Shared North Star:

raw input -> representation -> prediction -> reasoning/generation -> persistent accumulated cognition -> trainable general language/reasoning system.

Self-modification and RSI are not optimization targets. They remain bounded scientific hypotheses/mechanisms that may be tested only when they plausibly improve transfer-efficient learning under the unchanged external envelope.

## Current closure evidence

R4 executed without infrastructure failure through duplicate scientific execution, but the frozen validity contract failed. The measured class partition had one class with zero training-margin samples and only one class meeting the preregistered 25-row below-threshold support floor. Preserve the result as INVALID; do not retune thresholds, support floors, class definitions, data, or interpretation after observation. Any successor must preregister a new scientific premise rather than treating R4 as supported/mixed/negative.

## Rolling portfolio allocation

Use a rolling 20-successor planning window. Existing frozen successor commitments always take precedence.

- 45% / 9 of 20 — raw-input, representation, prediction, reasoning/generation progression.
- 30% / 6 of 20 — transfer-efficient learning / meta-learning research under fixed budgets; bounded self-modification is only one candidate mechanism.
- 15% / 3 of 20 — OOD/generalization.
- 10% / 2 of 20 — substrate falsification.

When bounded self-modification, self-maintenance, or meta-learning experiments are active, reserve approximately 5–10% of total effort for explicit alignment/control-falsification. Rebalance the other streams if necessary; never remove substrate/scientific falsification.

## Transfer-efficient learning objective

The primary optimization target is **increasingly efficient learning of unseen future datasets under fixed data, compute, capacity, and authority budgets**.

Do not directly optimize for self-modification. Bounded self-modification may be tested only as one possible mechanism for improving future learning efficiency.

A modification may not be retained merely because it improves the dataset, task family, diagnostic subset, or failure class that produced it. Retention requires a preregistered transfer advantage on one or more disjoint future datasets that were unavailable during diagnosis/selection.

Required comparison discipline:
- exact predecessor/no-modification baseline under the same data budget;
- same adaptation compute/runtime budget;
- same model/substrate capacity;
- same tool and authority envelope;
- same evaluation protocol and frozen acceptance thresholds;
- disjoint future dataset(s) for transfer validation;
- no post-result choice of transfer dataset.

Preferred transfer metrics include data-to-threshold, fixed-budget held-out performance, adaptation compute-to-threshold, error reduction per training byte/example, retention of prior capability, OOD transfer, and performance on later unseen datasets after equal-budget adaptation.

If a candidate helps the originating dataset but fails to transfer, revert it. Record the result as mechanism-specific evidence rather than self-improvement.

## Learning-efficiency and self-modification evidence ladder

Use **self-improvement precursor** only for bounded mechanisms that have been demonstrated but have not yet shown transferable advantage on later unseen datasets.

Do not use **autonomous self-improvement** unless a bounded closed-loop process is demonstrated across multiple unseen cases **and** retained modifications improve learning on preregistered disjoint future datasets under the same budgets.

Do not use **recursive self-improvement** unless an accepted transferable improvement measurably improves later learning/adaptation or later improvement cycles under the same unchanged external envelope.

Near-term progression should favor:

unseen-dataset baseline -> native learning-efficiency diagnosis -> bounded mechanism selection (which may include no modification) -> equal-budget candidate training/adaptation -> sealed transfer evaluation on a different preregistered dataset -> retain/revert -> second unseen dataset -> repeated transfer-efficient learning -> eventual meta-learning test.

Self-diagnosis and causal attribution remain useful research components, but they are instruments for improving future learning efficiency rather than goals in themselves.

Keep OOD pressure active. Improvement confined to the exact diagnostic domain or originating dataset is insufficient evidence.

## Candidate-only self-improvement boundary

Wingless may diagnose a deficiency, select from a preregistered bounded candidate set, apply a candidate cognitive modification in isolated research state, evaluate it through the frozen external harness, recommend retain/revert, and continue after external acceptance.

Wingless may not autonomously widen the modification set; alter its evaluator, thresholds, resource limits, permissions, credentials, tools, network policy, logs, rollback, or shutdown; promote itself; mutate accepted refs; create schedulers, queue consumers, retry authority, uncontrolled persistent workers, or a competing control plane; inspect sealed held-out labels; or preserve rejected state outside the experiment lifecycle.

Capability may grow. Authority may not.

## External control envelope

The following are externally enforced and non-self-modifiable: authority boundaries, permission model, execution controls, evaluator/scoring implementation, held-out datasets and labels, scientific/safety acceptance thresholds, hard resource ceilings, qualification contract, logging/audit and causal-trace mechanisms, rollback/shutdown, accepted-ref and promotion authority, network/tool policy, credentials/secrets, orchestration/queue/scheduler authority.

A run that modifies or bypasses those controls is invalid unless the explicit preregistered purpose is to test whether the attempt is detected and blocked.

## Determinism and semantic identity

Every experiment declares deterministic vs stochastic behavior before execution.

Deterministic experiments must canonicalize ordering and exclude clocks, random IDs, mutable external data, environment-dependent paths, and other incidental nondeterminism from scientific outputs. Repeated byte-identical execution is required where practical; replay mismatch is an investigation condition.

Stochastic experiments pin RNG/version where relevant, preregister seeds and budgets, use disjoint seed sets where required, freeze statistical criteria, and reproduce the same seeded trajectory when the implementation contract permits.

Every result binds all semantically relevant identity: source commit/tree, dependencies, toolchain, configuration, harness, qualification contract, dataset/split, seed, capacity/resource configuration, experiment definition, and allowed modification set. Semantic identity drift fails closed.

## Fail-closed execution

Stop the affected path on unexplained deterministic mismatch, unexpected source identity, dirty/unexpected checkout, stale evidence, provenance gaps, acceptance-contract drift, cache-key mismatch, missing artifacts, unauthorized state/file modification, scope widening, unknown required-resource state, inconsistent controller/experiment state, ambiguous retain/revert state, unexplained partial write, or unexpected external side effects.

Do not continue through unexplained errors merely to obtain a result.

## INFRA_INVALID vs SCIENTIFIC_RESULT

Infrastructure failures are separate from scientific supported/mixed/negative results. An infrastructure repair is permitted only when the preregistered science remains semantically identical.

Infrastructure repair must preserve preregistration, hypothesis, datasets/splits/seeds, thresholds, capacity/resource budgets, and frozen variables; record the exact mechanical delta; and requalify affected interfaces.

Scientific negatives are first-class evidence and are never tuned away or relabeled as infrastructure failures.

## Cause -> effect trace and defect records

Each meaningful transition must preserve structured evidence for:

internal signal -> diagnosed defect -> candidate set -> selection rationale -> exact bounded delta -> focused qualification -> affected-interface qualification -> full regression -> held-out result -> retain/revert -> resulting state identity -> next question.

Reproducible implementation, harness, state-machine, orchestration, determinism, or qualification defects require a machine-readable defect record including exact identity, triggering event/input, observed vs expected behavior, first causal divergence, downstream effects, contamination status, rerun permission, repair scope/commit, focused/full regressions, and closure classification. Preserve the original failure evidence.

## Component diagnostics

Maintain non-production diagnostic paths for raw-input ingestion, representation construction, motif/relation mechanisms, prediction/readout, confidence/probability/margin, reasoning bridge, self-diagnosis, causal attribution, modification selection, candidate application, evaluator interface, and retain/revert logic.

Diagnostics must not activate production, widen authority, mutate accepted refs, or gain uncontrolled external access.

## Concurrency / race / retry controls

Test critical transitions under duplicate/stale/reordered callbacks, interruptions, partial writes, repeated delivery, concurrent successor attempts, simultaneous retain/revert requests, branch movement, stale checkouts, artifact races, event storms, restart/recovery, and retry/replay.

Require idempotency, compare-and-swap/expected-version semantics where applicable, exactly-once logical acceptance, no duplicate successors, no contradictory accepted state, no orphan candidate state, deterministic recovery, and traceable ownership.

Retries are bounded to the failing idempotent operation with explicit retryable/non-retryable conditions, maximum attempts, idempotency identity, evidence-validity rule, and terminal state. Wingless cannot create retry authority.

## Reversibility

Every experiment that proposes a cognitive/state modification preserves exact predecessor identity, candidate identity, bounded delta, selection provenance, originating-dataset evidence, disjoint future-transfer evidence, deterministic rollback where applicable, and explicit retain/revert outcome.

Rejected candidates must leave no unauthorized residual state. Rollback qualification verifies expected source/cognitive state, external control integrity, no orphan workers, no hidden persistent state, no unauthorized resource use, and no evaluator change.

## Qualification optimization

Use the compositional fast path only when exact identity proves reuse is safe:

persistent qualification checkout -> exact SHA/tree identity -> bounded patch -> cached focused tests -> affected-interface tests -> one required full regression -> build -> external accept/promote.

Persistent checkout existence is not evidence. Cache reuse requires an exact key over all semantically relevant identity. Identity drift, dirty state, scope widening, provenance gaps, harness/qualification drift, cache mismatch, or uncertainty forces fresh qualification.

Never remove the final full integration/regression, and never gain speed by weakening safety, determinism, isolation, scientific controls, acceptance criteria, or provenance.

## Control-falsification stream

Preregister experiments that try to falsify bounded capability growth, including evaluator gaming vs legitimate improvement, constraint bypass, rejection/corrigibility, shutdown invariance, authority delta = 0, resource invariance, task improvement vs test manipulation, logging integrity, persistence-boundary cleanup, and capability-triggered control drift.

If a run unexpectedly attempts evaluator manipulation, extra permissions/resources, unauthorized persistence/replication, competing orchestration, strategic harness bypass, shutdown/rollback/logging modification, materially out-of-scope self-modification, or unexplained isolation crossing: freeze evidence, do not retune or rerun around it, classify separately, and halt that successor chain for architectural review. Report observable behavior only; do not infer intent or sentience.

## Successor closure contract

Every experiment closes with experiment identity, scientific classification, infrastructure classification if applicable, technical/raw evidence, plain-speak interpretation, deterministic/probabilistic qualification status, causal trace, hypothesis strengthened/eliminated, exact predecessor/resulting-state identity, exact successor question, and preregistered successor where scientifically safe and obvious.

Automatic continuation is allowed only when the successor follows from frozen evidence, no safety tripwire fired, no architectural premise changed, provenance is complete, qualification is valid, fixed data/compute/capacity/authority budgets remain explicit, and the authority/resource envelope is unchanged.
