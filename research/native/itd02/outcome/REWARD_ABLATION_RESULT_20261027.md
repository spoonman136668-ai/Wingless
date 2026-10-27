# ITD-02 authentic reward-ablation/cold-restart evidence — 2026-10-27

Status: **REWARD-SENSITIVE BEHAVIOR OBSERVED; MIXED PERFORMANCE; NOT SCIENTIFICALLY ACCEPTED**.

Source: [GitHub Actions run 40513127899](https://github.com/spoonman136668-ai/Wingless/actions/runs/40513127899); `itd02-reward-ablation-restart` artifact `12737798236` (archive SHA-256 `f6cf3d86c670b161f2d2ef8806b4be94dd76ef75547a785a6c66d2f6d3ed84ea`). Commit `3266224e8786378ec3e6f92e875cb56d2e1a9c4a`. MiniGrid 2.3.1, Gymnasium 0.29.1; `MiniGrid-Empty-Random-5x5-v0`. One previously examined training seed 19; previously examined probe seeds 7 and 11; 64-step cap per episode.

The native reward-aware learner and its reward-blind ablation followed **identical seven-action training trajectories** and both reached the source goal with return 0.937. The selected-feedback reward differed (actual versus 0), producing distinct checkpoint hashes:
- Reward-aware: `d85d00d0eb6d3a7bab6b18c29a88262b9121414105b05711a24d5703616c413d`
- Reward-blind: `95cd588f2d3253a1acd8406cea0325709dc5c1cdcb8066ffb24e0d0c6c5451e5`

Each evaluation episode was restored into a fresh native IPC process from the same arm's unmodified training checkpoint, verified byte-for-byte. Results:

| Probe seed | Reward-aware | Reward-blind |
|---|---|---|
| 7 | goal in 18 steps, return 0.838 | no goal in 64 steps, return 0 |
| 11 | goal in 13 steps, return 0.883 | goal in 4 steps, return 0.964 |

**Interpretation:** The two arms had matched training exposure and different reward inputs. Divergent checkpoint and later behavior provide direct evidence of reward-sensitive adaptation under these engineering conditions. They do **not** establish consistent advantage: the reward-aware arm wins on seed 7 while reward-blind wins on seed 11. Probe feedback also updates disposable copies during the probe, so this is a cold-restart **online-adaptation** comparison, not a strict frozen-policy zero-update test.

**Scientific limitations:** The seeds were already examined; one training episode, two probe episodes, no independent held-out family, no strong conventional baseline in this ablation, no uncertainty intervals, and no full resource accounting. Evidence is authentic but exploratory. No claims of broad generalization, autonomous intelligence, superiority or scientific first use are authorized.

**Next discrimination gate:** Predeclare previously unexamined seeds, keep matched reward-aware/reward-blind arms and cold restore, retain all seed-level outcomes, then add a true frozen-evaluation mode before interpreting held-out transfer.
