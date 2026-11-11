# ITD-02 fresh-seed reward ablation — observed 2026-10-27

**Scope: authentic source engineering experiment; no scientific acceptance.**

Provenance: [GitHub Actions run 40516799061](https://github.com/spoonman136668-ai/Wingless/actions/runs/40516799061), artifact `itd02-reward-ablation-restart` ID `12738514018`, archive SHA-256 `87b3a179e39f434df9b0a3bd50c3165902db7752ac872717004b1ab3d1014638`. Exact source revision `0d972214644581c4e0a487c8ad1409f9615d1adb`. Native checkpoint training seed 19 is **previously examined**. Probe seeds 23, 29, 31, 37 were predeclared and had not been examined in previous PR #399 experiments. Environment `MiniGrid-Empty-Random-5x5-v0`; source-only public observations; 64 steps each.

Both arms followed the same seven training actions, reached the goal with return 0.937, and produced distinct saved checkpoints. Native restore preserved checkpoints byte-for-byte before each probe. Probe feedback was processed by each disposable evaluation learner and could change its policy, so these are **online-learning probes**, not frozen-policy tests.

| Probe seed | Reward-aware | Reward-blind |
|---|---|---|
| 23 | 32 steps / return 0.712 / goal | 14 / 0.874 / goal |
| 29 | 3 / 0.973 / goal | 64 / 0 / no goal |
| 31 | 1 / 0.991 / goal | 64 / 0 / no goal |
| 37 | 32 / 0.712 / goal | 64 / 0 / no goal |

Reward-aware: 4/4 goal completions. Reward-blind: 1/4. Returns are substantially different and the reward-aware arm is better on three seeds; the reward-blind arm is faster on seed 23. **Do not erase the seed-23 counterexample.**

This is direct evidence of **reward-dependent native adaptation with cold-restart continuity in a real source environment**, but only four probe seeds and one previously explored training seed. It does not demonstrate a statistically stable improvement, zero-update frozen inference, cross-family transfer, independence of training design, or retail-hardware resource efficiency. The claim remains bounded to this exploratory comparison.

Next discrimination: isolated read-only native `peek_frozen` selection from an exact checkpoint clone per action, with *no policy updates during evaluation*. This intentionally omits within-episode native route adaptation, so it measures frozen preference use rather than the full deployed native actor. Predeclared probe seeds: 41, 43, 47, 53; compare both reward-aware and reward-blind arms with exact checkpoint identity checks.
