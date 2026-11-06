# ITD-02 selected-outcome learner — source-world experiment gate

Status: **PREREGISTRATION / NOT EXECUTED / NOT SCIENTIFICALLY ACCEPTED**.
This is an opt-in successor experiment. Do not modify frozen R1–R8, earlier evidence, or the read-only Atlas authority. Synthetic Go tests are engineering tests, not world-learning evidence.

## Source and custody
- Use a version-pinned MiniGrid environment with `FullObsWrapper`, explicit dependency lock, environment IDs, seed lists, installed versions and immutable source revision recorded *before* running evaluation. The native learner must receive only the six `r1.PublicGrid` fields derived from source-public observations. No private map, oracle path, task-family label or held-out reward map reaches the learner.
- The environment process owns reset, step, actual reward and termination/truncation. Only its completed selected-action transition may feed `outcomebridge`'s `feedback` operation, with monotonically sequential `event_id`. Record raw public before/after, selected action, environment reward, done flags and checkpoint digest.
- Map legal native actions 0/1/2 to the source environment's left/right/forward. Reject every other action. Verify the public-state decoder against observed source frames and record any ambiguity as a failure; do not fill missing source fields with privileged data.
- Refuse fabricated, replayed, misordered, mismatched or unselected feedback; stop on unresolved custody errors. A fresh process must restore the exact training checkpoint before held-out testing.

## Predeclared comparisons
- Training: predetermined seeds from one bounded source-world family. Held-out: disjoint seeds **and** an unfamiliar environment family declared before inspecting outcomes. No retraining on held-out episodes for the primary transfer claim.
- Matched budgets: identical observations, legal actions, steps, resets, CPU/memory limits and source interactions. Compare the actual contextual learner with (a) native R8, (b) fixed deterministic action policy, and (c) a reward-blind control whose selected public transitions are retained without reward differentiation. Log all failed episodes and rejected operations.
- Counterbalance assignments of contexts, action rewards, seed order and run order where supported by the authentic source. No after-the-fact removal of poor seeds or cherry-picked best checkpoints.
- Primary: held-out environment return and completion rate relative to controls. Secondary: training learning curve, unseen-seed return, cold-restart retention, selected-action validity, CPU wall time, peak RSS and interaction count. Report per-seed results and uncertainty intervals; do not replace performance with the number of passing unit tests.

## Fail-closed acceptance
- Evidence must contain: pinned environment/source versions, predeclared train/held-out splits, raw source episode transcript hashes, checkpoint-before/after hashes, exact native binary revision, control baselines, reproducible commands and complete compute accounting.
- No verified source execution => **SOURCE_NOT_RUN**. No positive held-out improvement over matched controls => **TRANSFER_NOT_DEMONSTRATED**. Invalid source provenance, hidden-information leakage, sample mismatch, checkpoint mutation during held-out testing or omitted failure => **REJECT**.
- Qualifying integration/CI can report engineering PASS while scientific status remains **NOT ACCEPTED**. Preserve negative outcomes. No merge, architectural replacement or scientific first-use authorization follows from this file.

## Immediate implementation boundary
The existing `outcomebridge` is an IPC learner only and does not execute MiniGrid. Implement a separate, opt-in source adapter and preregistered runner next. It must validate every public observation and return authentic selected transitions; the controlled evaluation must run before making a native-intelligence progress claim.
