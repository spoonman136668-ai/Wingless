# ITD-01 R7 — Wingless native online state prototype, prospective source-free DEV pilot

**PRE-IMPLEMENTATION FROZEN** in own research-only branch atop protected main `48ba343b403bc9b4389ce1777130e546b3c446bd`. This is a native learning prototype milestone **NOT the full ITD-01 integrated challenge**, not originally source-born H1/H2 science, and not accepted research promotion. R1–R6 source and negatives are unchanged. The common source-free observable-feedback contract below is shared **externally**; model internals, update rules, storage, weights and code are independent.

## Distinct native hypothesis

Native Wingless packed online error-corrected signed causal weights, two 4-feature contextual compressed predictive states, integer arithmetic only, source-owned Go code, no model-host/broker/transformer.

A learner receives public *context ID A/B* and four publicly sensed binary entities. It predicts a binary public terminal decision (action 0 or 1), then after the public episode's delayed feedback it may update its own learned state from that observed outcome only. A private DEV generator makes noisy causal contexts (A depends on public entity 2; B on inverted entity 1; independently flipped outcome probability 0.10). Hidden rule/cause never provided to learner; evaluator uses it solely for generating public outcomes. The model must not read private generator state, source-family H1/H2 or paid/broker inference. This is narrow sensory association evidence, NOT discovery of latent causes or multi-hop interactive transfer.

## Frozen cohorts and order

Fixed init seeds **17,31,47**. Each seed uses independent PRNG-specified source-free TRAIN A 320 public experiences, then public A VALIDATION 160 new experiences, then TRAIN B 320 intervening experiences, then post-interference A evaluation 160 disjoint public experiences, and B evaluation 160 disjoint experiences. All test labels withheld from model before prediction; no on-test training. One frozen fit per init, no tuning after test. Validate actual nonzero outcome-driven learned updates, independently seeded deterministic input and exactly serialized/restarted learned states before post-B A and B scoring. Public feedback arrives only after decision and a two-step delay marker; validate no learner update before feedback. No additional teacher labels, original data, external model inference or optimizer fitting.

## Frozen controls and thresholds

Prespecified controls: never-trained zero-state sham (same inference code), public current-only majority guess (no learned history), and erased learned state ablation after intervention B. One actual *learned-state pilot* success requires ALL:
1. A post-B test accuracy >= 0.75 and B test accuracy >= 0.75 on each init, despite intervening B episodes.
2. Post-B A accuracy <=0.05 below A prior validation accuracy; no catastrophic contextual overwrite.
3. Learned predictor advantage vs zero-state and constant-class control >= **0.15** absolute accuracy on A and B averaged across all 3 seeds.
4. Erasing stored contextual learned state decreases averaged post-B A or B task accuracy by >= 0.15.
5. Save/restore yields EXACT same predicted actions for all 160+160 post-B public test observations, confirms nonzero source-independent training updates and exact bytes/hash of saved state after cold process restart.
6. State serialized <= **32 KiB**; local peak RAM and total CPU recorded; no GPU, API, external model or hidden labels.
7. Tests on native GitHub-hosted Linux and Windows with no repository accepted-ref edits or scientific execution authority.

Report full score and negatives even on failure. Since generator law and two public context IDs are synthetic and pretrained by design, no claim of unfamiliar causal law, independently authored family transfer, autonomous architecture emergence, LLM-class intelligence or original first use is authorized. If threshold fails classify **SOURCE_FREE_NATIVE_LEARNING_NO_GO**. Missing reproducibility or privileged information is **INVALID**.

## External ITD-01 interface bridge

The source-owned native learner shall expose `Predict(context, public four-entity readings)`, `LearnFromCompletedPublicFeedback(...)`, `Save/Restore` and trained-update counter; future true ITD-01 observe-neutral/act adapter must be separate and cannot be called working until it has genuine full-episode run evidence. This experiment validates only a bounded experience→state→prediction→persistence subcycle. Wingless and Yggdrasil must not exchange architecture, parameters or memory state.

## Authority boundaries

R6 neutral external contract source SHA `d1c2eee8c7863349aa0721fe7a480473e052dde0` in FORGE; original P0F-C2 scientific cohorts and CKB controller R18 remain blocked/paused. No H1/H2 original first use, existing controller/scheduler change, KTRADE live/paper broker, upstream accepted-ref promotion, LLM broker delegation or paid APIs. PR review only after actual qualification.
