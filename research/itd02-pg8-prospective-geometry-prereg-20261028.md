# ITD-02-PG8 — prospective native geometry-shift falsification (preregistered; development only)

**Freeze chronology:** This preregistration is the FIRST commit on this isolated branch, based on Wingless PR #414 head `4f547b79708c7f3249a6cf2afb96d0f6f5b39919`. Record it before adding any executable experiment or running any of the below held-out seeds. Existing learner and original five-arm protocol are frozen. This is NOT a proposal to activate a research controller, merge a model, or assert scientific first use.

**Parent native source:** Wingless V4 source ancestor `21928da6e939b84cf71f70476215bfa348aecb85`; published native comparison reference: https://github.com/spoonman136668-ai/Wingless/actions/runs/40665581320. Known previous result: 0/8 held-out FourRooms goals for every arm with test-time learning; do not bury this negative.

## Frozen hypothesis and falsifier

Question: Does the existing selected-action goal-progress V4 policy transfer its learned navigation preference from random Empty-5x5 training to *new geometry* in a source-public Empty-8x8 world better than its predecessor, native comparators, conventional linear Q, and a goal-directed nonlearning reactive policy?

Predicted directional effect: V4 should reach at least one of eight new target goals AND should not lose on goal count to any of those controls. Classification:
- `PG8_DEVELOPMENT_POSITIVE_NOT_SCIENCE` only if V4 goals >= 1 AND V4 goals >= maximum across all four comparisons plus the public reactive reference.
- `PG8_NATIVE_TRANSFER_NO_GO` otherwise, including zero goals.
- Any source, identity, deterministic protocol, observation, seed reuse, or resource accounting violation is `INVALID`, not a negative experiment.

**Important interpretation:** Empty-8x8 tests an unobserved grid geometry, NOT unseen entity appearances, learned causal-role bindings, 24/48-hour retention, or independent new-world transfer. V4 does not implement the external causal-binding proposal. Do not elevate this development result to scientific acceptance even if positive.

## Locked protocol

- Frozen source interface: public `{x,y,goal_x,goal_y,direction,front_blocked}` only; no privileged MiniGrid or goal oracle, no new model capacity, no transformer/token cognitive substrate.
- Environment versions: `minigrid==2.3.1`, `gymnasium==0.29.1`. Train in `MiniGrid-Empty-Random-5x5-v0`. Test exclusively `MiniGrid-Empty-8x8-v0`.
- Train: reuse PREVIOUSLY PUBLISHED V4's exact 256 source-selected transitions per learning arm, seeds `7101,7103,...,7611`. This training corpus is *not* sealed or new; do not claim a fully fresh experiment.
- New target cohort locked now, before implementation or execution: `18041,18047,18061,18077,18089,18101,18119,18133`. These eight seeds are disjoint from V4 train and earlier V4 FourRooms test seeds; global provenance independence is unproven.
- Test: 100 maximum authentic selected steps per target seed; 8 episodes; native policies each restore their identical post-training checkpoint at each test seed. Test-time selected feedback remains **ON**, with exact update counters; this is online adaptation, not frozen zero-shot transfer. Linear Q receives the same 256 training transitions and test-time updates. A separate nonlearning public reactive policy receives the same initial observation, no training rewards, and an explicitly unmatched zero-training budget.
- Frozen five learning arms: V4 `progress_native`, its V3 original predecessor `progress_v3` from `0bd552d58b4b26f6189a6c62be67c30b7f22937f`, `motion_native`, `original_native`, conventional 9-feature `linear_q` as already defined in source V4 protocol. No hyperparameter tuning after seeing results.
- Additional strong navigation check: existing `PublicReactive` heuristic, same six source-public fields, no privileged world map. Interpret separately as a fixed nonlearner and not an equal-training learner.
- Primary outcome: count of real environment goal completions / 8, not position advances, episode closure, CI green, model plausibility, or diagnostic state changes. Report eight raw per-seed goal outcomes and full selected source trace hashes for every arm.
- Secondary: per-arm reward, steps, advances, checkpoint SHA and update accounting, Go runtime health, child/parent host peak RSS where available, wall and CPU. GitHub Ubuntu is not a retail-machine benchmark; energy not measured.
- Fail closed on sample mismatch, unequal public starting observations, missing authenticated source transitions, checkpoint drift, unexpected terminal response, duplicate output, nonfinite outcomes, new training budget or changed model/core.
- Run exactly once as a new source-bound read-only GitHub Actions scientific-development probe, with no provider key, write token, scheduler, promotion or native-model changes. No post-result retuning under this preregistration. Independently validate the resulting artifact in a credential-free CI job.

**Acceptance:** Development-only evidence, classification remains NO SCIENTIFIC FIRST USE. Preserve all positives and negatives, and require independently controlled unfamiliar object-family causal-role tests, preregistered 24/48h cold retention, strong conventional comparisons, and retail PC measurements before any North-Star progress claim. 
