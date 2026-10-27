# ITD-02 randomized-source matched controls — 2026-10-27

**Status:** AUTHENTIC SOURCE ENGINEERING EVIDENCE; NOT SCIENTIFICALLY ACCEPTED. No held-out learning or transfer is demonstrated.

**Provenance:** [GitHub Actions run 40507122466](https://github.com/spoonman136668-ai/Wingless/actions/runs/40507122466); `itd02-source-controls` artifact ID `12734858785`, archive SHA-256 `7143cf01b24d533f8f635eb31afb7365239332390672f9a8e616f08a43afef5f`. Exact evaluated commit `45cb287b8fdfb5a16511cbfe4a356470d043efb0`. Pinned `MiniGrid-Empty-Random-5x5-v0`, MiniGrid 2.3.1, Gymnasium 0.29.1. Independent new native processes per seed, 32-step maximum for each policy.

| Seed | Public start (x,y,dir) | Fixed steps/return/goal | Native R8 steps/return/goal | Contextual outcome steps/return/goal |
|---|---|---|---|---|
| 7 | (2,3,3) | 11 / 0.901 / yes | 32 / 0 / no | 32 / 0 / no |
| 11 | (3,2,2) | 2 / 0.982 / yes | 32 / 0 / no | 32 / 0 / no |
| 19 | (2,2,1) | 19 / 0.829 / yes | 32 / 0 / no | 7 / 0.937 / yes |

**Aggregate:** fixed 3/3 goal completions; R8 0/3; contextual learner 1/3. Mean returns: fixed 0.904, R8 0, contextual 0.31233. The contextual learner reaches a real source goal on seed 19, but it does not outperform the fixed control overall. Three distinct public initial states were verified.

**Interpretation:** This is a comparison of fresh learner instances in each bounded episode, **not** evidence of cumulative cross-episode reward learning. In the native contextual policy, initial per-context coverage drives action choice before sufficient reward samples exist. A single success can arise from its exploration trajectory without proving that reward-dependent policy improvement occurred. No cold restart, matched reward-blind contextual control, unfamiliar-family test, robust sample-size analysis, or resource accounting has been completed.

**Decision:** Preserve all failures and the single success. Do not claim broad generalization, native learning superiority, transfer, or scientific first-use. The next discriminator must separate reward-aware policy updating from matched reward-blind exploration across authentic episodes, with checkpoint retention and disjoint held-out seeds. Preserve immutable source/public observation custody and equal resource budgets.
