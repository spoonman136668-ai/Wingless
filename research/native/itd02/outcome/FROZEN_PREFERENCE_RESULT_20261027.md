# ITD-02 frozen preference source probe — 2026-10-27

**Result: source-executed positive reward-ablation discrimination, NOT scientific acceptance.**

Source: [GitHub Actions run 40549946135](https://github.com/spoonman136668-ai/Wingless/actions/runs/40549946135), `itd02-frozen-source-probe` artifact ID `12754597500`, ZIP SHA-256 `de676e707d8aceea3f0e7c576309d03b2c83c9d22efd0618162f171376003d68`. Tested code commit `9e7d947221596ff168c760ed72e1e30afb4a1618`. MiniGrid 2.3.1 / Gymnasium 0.29.1; public-only `MiniGrid-Empty-Random-5x5-v0`.

Training seed 19 was **previously observed**, followed by four newly selected probe seeds 41, 43, 47, and 53. Each arm received seven identical selected actions in its source training episode and achieved source return 0.937. Reward-aware and zero-reward-blind checkpoint digests differed:
- Reward-aware `d85d00d0eb6d3a7bab6b18c29a88262b9121414105b05711a24d5703616c413d`
- Reward-blind `95cd588f2d3253a1acd8406cea0325709dc5c1cdcb8066ffb24e0d0c6c5451e5`

Each source probe restored its original training checkpoint, and every choice used `peek_frozen`, which calls native `Choose` on a separately restored clone. No probe source feedback was passed to the active learner. Checkpoint bytes stayed identical before and after each episode; qualification included explicit tests that frozen peeks cannot arm/modify the original native learner.

| New probe seed | Reward-aware | Reward-blind |
|---|---|---|
| 41 | goal, 7 steps, return 0.937 | failed, 64 steps, return 0 |
| 43 | goal, 11 steps, return 0.901 | failed, 64 steps, return 0 |
| 47 | goal, 11 steps, return 0.901 | failed, 64 steps, return 0 |
| 53 | goal, 4 steps, return 0.964 | failed, 64 steps, return 0 |

**Interpretation:** In this small fixed-checkpoint preference evaluation, the reward-aware native state achieved 4/4 real source goals and the matched reward-blind native state achieved 0/4. These are authentic selected environment actions on previously unexamined seeds, and no within-probe learning is possible in the native checkpoint being measured. This is meaningful controlled evidence that reward experience changed action preferences usefully in this narrow family.

**Limitations:** Source training seed 19 has been studied and selected for follow-up; it is not a preindependent scientific training choice. The four probe seeds come from the *same MiniGrid environment family*, not an unfamiliar family. The `peek_frozen` operation restores a fixed native checkpoint clone **for every action**, intentionally omitting the complete actor's normal within-episode route-state evolution. The experiment is limited to four probes, with no confidence interval, complete resource accounting, or strong conventional learning baseline. The result is **not** proof of broad native intelligence, compositional unfamiliar-family transfer, or scientific first use.

**Next discriminator:** run independent, predeclared novel-family frozen-preference and full-actor comparisons at matched information/interaction budgets, include a strong conventional learner, ensure consumer-PC CPU/memory/latency accounting, and preserve all failures. Do not tune against probes 41/43/47/53 and then rebrand them as held out.
