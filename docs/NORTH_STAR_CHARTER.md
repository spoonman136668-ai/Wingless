# Wingless North-Star Research Charter

**Status:** Proposed research-governance overlay  
**Effective scope:** Experiments preregistered after adoption  
**Non-retroactivity:** This charter does not alter, reinterpret, retune, or invalidate experiments already in flight or already frozen.

## North Star

Wingless contributes toward a future general-purpose LLM-class intelligent system that can communicate, reason, plan, code, and use tools while also supporting persistent cognition, continual capability retention, prospective self-monitoring, bounded self-maintenance, iterative refinement, deterministic or statistically reproducible operation as appropriate, and independent safety/authorization/recovery controls.

Wingless is not required to recreate a foundation model. The functional LLM/reasoning core is treated as a replaceable substrate behind a stable adapter boundary. Wingless research focuses on persistent cognitive machinery that can remain useful across model replacements.

## Project role

Wingless investigates:
- persistent reusable cognitive procedures and memory;
- bounded composition and routing;
- fixed-resource continual learning;
- degradation/failure-risk detection from native state;
- bounded maintenance and correction;
- retain/revert discipline;
- transfer and repeated refinement.

Wingless does **not** own external execution authority, deployment authority, credentials, policy authority, or the right to expand its own resource/permission envelope.

## Functional LLM requirement

The long-term system must include a strong general-purpose model capable of language, reasoning, planning, coding, tool interpretation, and broad task competence.

The LLM interface must be replaceable. A model may:
- consume approved context and persistent state;
- emit reasoning/results;
- propose bounded tool/action requests;
- emit confidence/diagnostic metadata.

The model may not unilaterally:
- authorize external execution;
- change guardrails;
- grant itself new permissions;
- accept its own durable cognitive modifications;
- overwrite immutable baselines.

## Determinism and reproducibility

Default research rule: use exact determinism wherever technically feasible.

For accepted deterministic paths:
**same accepted state + same inputs + same frozen policy + same resource budget -> same decision and state transition.**

Where stochasticity is scientifically necessary, the experiment must use controlled statistical reproducibility with:
- declared seed/distribution policy;
- preregistered sample count;
- frozen aggregation;
- declared uncertainty bounds;
- replayable provenance.

Every retained state transition must be attributable to its input state, policy version, evidence, authorization, and resulting state hash or equivalent stable identity.

## Defense in depth

A future autonomous loop must preserve independent layers:
1. capability boundary;
2. proposal generation;
3. action authorization;
4. runtime containment;
5. independent verification;
6. immutable audit/provenance;
7. rollback/checkpoint recovery;
8. external stop/shutdown authority.

No cognitive component may rewrite the authority envelope that constrains it.

## Research governance

Every mainline experiment must do at least one of the following:
- advance a North-Star gate;
- falsify an assumption required by a gate;
- reduce uncertainty about how to cross a gate.

Scientific negatives are valid results.

No threshold, metric, action budget, gate, or success condition may be retuned after viewing evaluation results unless a new experiment is preregistered.

External information may define a challenge or evaluation but must not encode the internal solution.

## Current strategic frontier

The current near-term bridge is:

**native state -> prospective risk -> frozen bounded intervention -> independent validation -> retain/revert -> held-out transfer -> repeat**

Current experiments remain authoritative for what has or has not been demonstrated. This charter is direction, not evidence.

## Claim discipline

Do not claim autonomous self-improvement until Wingless has demonstrated repeated, resource-bounded cycles in which it:
1. detects a genuine deficit or impending failure using allowed native signals;
2. selects a bounded response without oracle leakage;
3. applies only authorized changes;
4. independently validates the result;
5. retains beneficial change or reverts harmful change;
6. transfers benefit to held-out conditions;
7. repeats the loop without expanding authority or resource ceilings.
