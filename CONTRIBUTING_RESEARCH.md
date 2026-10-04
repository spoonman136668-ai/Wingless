# External Research Contributions

Wingless accepts independent research contributions that are run outside the official autonomous research loop.

The purpose of this lane is to increase replication, falsification, ablation, robustness, and disjoint-dataset evidence without giving external experiments authority over accepted research state.

## Core rule

A contribution may produce **evidence**. It does not automatically change accepted state, thresholds, datasets, promotion refs, successor selection, execution authority, or research-loop authority.

Official acceptance remains a separate qualification step.

## What contributors can run

Contributors may use their own workstation, GitHub Actions, university/organization compute, cloud compute, or other reproducible environments.

Useful contribution classes:

1. **Replication** — reproduce a frozen experiment independently.
2. **Disjoint replication** — apply a frozen mechanism to a genuinely new dataset/context.
3. **Falsification** — preregister a test intended to break a current hypothesis.
4. **Ablation** — remove one bounded mechanism to test causality.
5. **Control** — random, sham/no-op, equal-budget, fixed-heuristic, or other frozen control.
6. **Robustness** — frozen seed/order/runtime/environment variations.
7. **Successor proposal** — propose a new preregistered experiment without claiming accepted evidence until it is run and reviewed.

## Evidence classes

A submission is reviewed as one of:

- `REPLICATION_SUPPORTED`
- `REPLICATION_NEGATIVE`
- `SCIENTIFIC_SUPPORTED`
- `SCIENTIFIC_NEGATIVE`
- `MIXED`
- `INVALID`
- `UNPREREGISTERED`

Negative results are first-class evidence.

Exploratory work is welcome, but if the primary result was observed before the hypothesis, thresholds, split, budget, and classification rule were frozen, it must be marked `UNPREREGISTERED` and cannot be promoted as confirmatory evidence.

## Required workflow

### 1. Freeze a preregistration before the primary run

Create:

`research/contrib/<contributor>/<experiment-id>/prereg.json`

Use the template in `research/contrib/experiment-template.json`.

Commit this file **before** producing or inspecting the primary outcome.

For confirmatory evidence, the preregistration commit must precede the primary result artifact/logs.

### 2. Pin exact identity

Record:

- project
- exact parent commit SHA
- experiment ID
- dataset/context identities and hashes
- code/config identity
- permitted changed paths
- seeds/splits/order
- data and compute budgets
- capacity
- evaluator/scoring rule
- success/mixed/negative/invalid classification rules
- expected deterministic/stochastic behavior

Do not silently update the parent or inputs after running.

### 3. Run outside the official loop

Use any reproducible environment you control.

Do not use production credentials, production authority, accepted-ref mutation, hidden evaluator labels, or private held-out answers.

Do not increase capacity, data, compute, model calls, tools, or permissions unless that variable was explicitly preregistered.

### 4. Preserve provenance

Add:

- `result.json` — primary metrics and classification
- `provenance.json` — source SHA, tree/config/runtime/dependency identities, hardware/OS where relevant
- `README.md` — plain-language question, method, result, limitations, and exact reproduction command
- optional logs/artifacts needed to reproduce the claim

If the experiment is stochastic, include all preregistered seeds and aggregate statistics.

### 5. Submit a pull request

Use the Research Contribution PR template.

The PR must identify the preregistration commit and the later result commit separately.

Review must be able to establish that the result did not exist before the preregistration boundary.

## Scientific constraints

Unless the preregistration explicitly makes one of these the independent variable:

- no post-result threshold changes
- no post-result split changes
- no hidden capacity growth
- no extra model calls
- no evaluator/held-out leakage
- no post-result policy selection
- no silent compute/data budget increase
- no accepted-ref mutation
- no production authority
- no queue/scheduler authority
- no changing a negative into an infrastructure failure merely because the result was unfavorable

Infrastructure defects and scientific outcomes must be classified separately.

## Dataset contributions

New datasets/contexts are especially useful.

A dataset contribution must state:

- source/license
- acquisition method
- immutable hash
- size
- preprocessing
- whether the contributor has inspected evaluation labels/outcomes
- contamination risks
- relationship to existing manifests

Do not submit private personal data, credentials, proprietary data without redistribution rights, or data whose provenance cannot be documented.

## Isolation

Contributor results remain external evidence until qualified.

A contributor result must never automatically:

- promote a branch
- mutate accepted state
- change the research North Star
- alter the official evaluator
- authorize more resources
- launch an official successor
- modify the autonomous research controller

## High-value contributions right now

For Wingless, prioritize independent replication, falsification, and genuinely disjoint transfer/generalization tests over small parameter tweaks.

A clean negative on a new context is more valuable than a tuned positive on an already studied one.

## Directory layout

```
research/contrib/
  experiment-template.json
  <contributor>/
    <experiment-id>/
      prereg.json
      README.md
      result.json
      provenance.json
      ...bounded experiment code/artifacts...
```

This lane is designed so outside researchers can contribute meaningful evidence without needing access to the official runners or autonomous loop.
