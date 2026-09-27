# Wingless Experiment Preregistration Template

Use this template **before** execution. Freeze it with the experiment branch/commit. Do not modify scientific acceptance criteria after observing evaluation results.

## Identity
- Experiment ID:
- Parent experiment / accepted baseline:
- North-Star gate:
- Local research family:
- Date:
- Frozen commit:
- Counterfactual only / live activation:
- External execution authority: none unless separately authorized.

## Scientific question
- Question:
- Hypothesis:
- Competing explanation:
- What result would be informative even if the hypothesis fails:

## Manipulation
- Independent variable(s):
- Dependent variable(s):
- Frozen controls:
- No-intervention control:
- Resource-matched sham/random control:
- Held-out split / transfer condition:

## Resource envelope
- Compute budget:
- Memory/state budget:
- Action/intervention budget:
- Persistent-state budget:
- Any other hard ceiling:
- Rule: no scientific PASS if the experiment exceeds its declared envelope.

## Allowed information
- Native signals available to the system:
- External challenge information allowed:
- Prohibited oracle/future information:
- Information visible only to evaluator:

## Determinism / reproducibility
- Mode: exact deterministic / controlled stochastic
- Seed policy:
- Model/version pinning:
- Policy/gate version:
- State/checkpoint identity:
- Replay requirements:
- Statistical aggregation and uncertainty rule, if applicable:

## Intervention / modification authority
- Allowed modification classes:
- Forbidden modification classes:
- May the actor write durable state?:
- Who/what authorizes application?:
- Who/what independently verifies acceptance?:
- Rollback checkpoint and procedure:

## Metrics
- Primary functional metric:
- Prospective warning metric:
- False-positive / unnecessary-intervention metric:
- False-negative / missed-failure metric:
- Lead-time metric:
- Collateral-damage metric:
- Prior-capability regression metric:
- Resource-cost metric:
- Retention metric:
- Transfer/generalization metric:

## Frozen success criterion
State exact thresholds or comparison rules here.

## Frozen falsification criterion
State the observation that would reject the proposed mechanism or block advancement.

## Guardrail invariants
All must remain true:
- no authority expansion;
- no resource-ceiling expansion;
- no hidden external execution;
- no guardrail modification by the cognitive actor;
- immutable baseline/checkpoint preserved;
- complete provenance/audit record;
- rollback remains available;
- OOD/uncertain cases follow the frozen abstention or fallback rule.

## Result handling
- On scientific PASS:
- On scientific FAIL:
- On infrastructure defect:
- Retain/revert rule:
- Next experiment allowed:
- Claims explicitly **not** authorized by this result:

## Plain-speak interpretation
- Positive result would mean:
- Negative result would mean:
