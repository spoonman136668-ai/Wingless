# ITD-02 frozen/native/conventional linear Q cross-family discrimination — 2026-10-27

**Authentic source observations; SCIENTIFIC SUPERIORITY NOT ACCEPTED.**

Source: [Actions run 40577082353](https://github.com/spoonman136668-ai/Wingless/actions/runs/40577082353), `itd02-linear-q-crossfamily` artifact ID `12768415639`, ZIP SHA-256 `2ab27d06264313e05e63773383a5503462aa5583b58a80e1a33ecbd5dc6b7b47`. Evaluated commit `48f6a73a98fb41f64505f779edbcb1a9f17c744e`. Environment versions pinned: MiniGrid 2.3.1 / Gymnasium 0.29.1. Training `MiniGrid-Empty-Random-5x5-v0` seeds 101, 103, 107, 109, 113, 127, 131, 137. Evaluation `MiniGrid-FourRooms-v0` seeds 139, 149, 151, 157, 163, 167, 173, 179.

| Arm | FourRooms goals / 8 | Successful source seed | Episode steps | Source return |
|---|---:|---|---:|---:|
| Reward-aware native **online** | 1 | 139 | 85 | 0.235 |
| Reward-aware native **frozen preference** | 1 | 179 | 28 | 0.748 |
| Reward-blind native online | 0 | none | — | 0 |
| Reward-blind native frozen preference | 0 | none | — | 0 |
| Conventional 27-weight linear Q, **frozen** | 0 | none | — | 0 |

The two reward-aware actor evaluation modes succeed on **different** seeds. Frozen preference is implemented by snapshot-cloning the native learner for every selected action; this intentionally removes its full route-state continuity. The result demonstrates reward-influenced source action choice under two modes, not stable or deployment-ready cross-family transfer.

Training interaction counts: reward-aware native **519**, reward-blind native **631**, linear Q **431** actual source steps. Although each arm trained on the same 8 environments and had the same per-episode 100-step limit, the *total interactions were not matched*. Exact training budget parity is **FAIL**; the workflow explicitly reports `training_interactions_matched=false`. Native online evaluation also differs from the frozen conventional Q policy, and must not be compared as if inference/training permissions were identical.

The conventional linear Q uses fixed, unoptimized hyperparameters (alpha 0.1, gamma 0.9, exploration 0.2). It is a valid baseline mechanism, **not** a strong established optimized RL control. There is no demonstration of statistical superiority, robust unfamiliar transfer, independently matched acquisition budget, or retail consumer-hardware efficiency. These eight seeds have now been examined and cannot be used as blind holdout seeds in later experiments.

**Research conclusion:** Useful source-based evidence exists that rewarded persistent Wingless cognition can solve some unfamiliar-family instances; the native learner still fails 7/8 in each mode. Preserve every negative. The next scientific acceptance gate is a qualified *equal-source-step* comparison with an independently strong conventional learner, exact evaluation-update parity, full hardware accounting, and a distinct blinded evaluation split. No merge/scientific first-use authorization follows.
