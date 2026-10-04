# External contribution: R31 causal-state sufficiency

## Status

**Classification: `INVALID`**

This is external contributor evidence only. It is **not** an accepted Wingless result, does not promote any branch or mechanism, does not change accepted research state, and does not authorize a successor.

## Scientific question

Does a representation that preserves domain-specific temporal pre-intervention learning dynamics predict the sign of fixed-budget counterfactual marginal value on held-out and genuinely disjoint manifests better than the existing six-feature aggregate representation?

The experiment was designed after the valid negative R30 historical-sign-instability correction result. Its purpose was to distinguish a representation bottleneck from cross-manifest invariance/evidence and observability bottlenecks.

## Frozen design

The preregistration was committed before the primary workflow was activated.

- Parent: `ab1d7467edc021d70a96fcc3b749651e4bb46fea`
- Initial preregistration commit: `c9f8765bb4870834dc1667a1ce8ec0add0825754`
- Frozen runner SHA-256: `d7d0b4c1409026c9d378ed791acd4a9f2a7d659eab8706ee60a765b305795c75`
- Final pre-run identity commit: `08ed30192d73ca6b69ef5cea524a1ffc145c823c`
- Workflow activation head: `884e4e99512c5bdfe2394e45c579becc7dc64896`
- Primary run: GitHub Actions `37243786526`
- Job: `111557550102`

The treatment compared:

1. the R27/R30 six-feature aggregate control;
2. a frozen 72-dimensional domain-by-time raw trajectory;
3. a deterministic 6-dimensional PCA latent state learned from training manifests only.

The preregistration also froze a random-projection control, policy-only control, cyclic label permutation, and temporal-permutation falsification controls.

## Prospective sixth manifest

Before execution, the sixth manifest was frozen from three disjoint pinned repositories:

- code: Prometheus `promql/engine.go`;
- structured: OpenAPI Specification `yarn.lock`;
- technical prose: Flask `docs/quickstart.rst`.

Full-file and prefix SHA-256 identities are recorded in `prereg.json` and `provenance.json`.

## Result

The workflow successfully passed:

- checkout;
- frozen runner SHA-256 verification;
- Python syntax/compile verification;
- parent ancestry verification;
- download of all frozen sources;
- all full-source SHA-256 checks;
- all prefix SHA-256 checks.

The primary experiment then reached the preregistered prospective evaluator and encountered a **single-class sixth-manifest target**. The evaluator defined positive marginal value as `actual advantage > 0` and required both positive and non-positive examples in every required fold.

The frozen `class_metrics` routine therefore returned `valid=false` for the sixth manifest. The result serializer subsequently attempted to read the undefined `sign_accuracy` field and raised:

`KeyError: 'sign_accuracy'`

The preregistration explicitly declared a required fold with a missing positive/non-positive class to be **INVALID**. Therefore this run is classified `INVALID`, not `SCIENTIFIC_NEGATIVE`.

No representation-sufficiency conclusion is claimed.

## Why the run was not repaired and rerun

The single-class prospective outcome was discovered during the primary run. Changing the evaluator, thresholds, target definition, sixth manifest, or admission rule afterward would violate the no-post-result-tuning contract.

The serialization failure could be mechanically repaired, but doing so would not change the preregistered classification: the experiment was already invalid. The run was therefore stopped and preserved without a post-outcome rerun.

## Interpretation

This result does **not** weaken or support the latent-state hypothesis.

It instead exposes a design requirement for the next preregistered experiment: the prospective evaluation contract must remain defined under class degeneracy, or the manifest-admission rule must be frozen in advance in a way that guarantees an evaluable target distribution without inspecting the prospective outcomes.

The observed sixth-manifest outcome must not be reused as fresh confirmatory evidence in that replacement experiment.

## Reproduction

The exact frozen executable is:

`research/contrib/chatgpt-external/WLM-LM-EXTERNAL-FUTURE-DATA-CAUSAL-STATE-SUFFICIENCY-R31-CONTRIB-CGPT/run.py`

The original workflow is:

`.github/workflows/contrib-r31-causal-state-sufficiency.yml`

A reproduction of the frozen run is expected to reproduce the invalid condition rather than produce a scientific classification. The GitHub Actions run `37243786526` is the preserved primary execution record.

## Authority

- research only: yes
- accepted-state mutation: no
- production authority: no
- promotion authority: no
- scheduler/queue authority: no
- RSI success: no
- adoption candidate: external evidence only
