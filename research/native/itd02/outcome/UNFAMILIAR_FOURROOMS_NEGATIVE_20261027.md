# ITD-02 unfamiliar-family frozen-preference result — 2026-10-27

**SCIENTIFIC NEGATIVE: NO FOURROOMS GOAL COMPLETIONS.**

Source: [GitHub Actions run 40573264019](https://github.com/spoonman136668-ai/Wingless/actions/runs/40573264019), `itd02-unfamiliar-family` artifact ID `12765558641`, archive SHA-256 `90ed7f2a62d150966a47e96d6759645c5ace7e5e6de0ed6a92c81a80aae7c9cc`. Evaluated commit `2cbc61edb63d9adae961cf57da31a71d6788224b`. MiniGrid 2.3.1, Gymnasium 0.29.1. Training: `MiniGrid-Empty-Random-5x5-v0` seed 19 (previously studied); target: genuinely different `MiniGrid-FourRooms-v0` seeds 61, 67, 71, 73.

| Policy | Goal successes | Episode outcomes | Total return |
|---|---:|---|---:|
| Frozen reward-aware Wingless | **0/4** | all four source-truncated at 100 steps | 0 |
| Frozen reward-blind Wingless | **0/4** | all four source-truncated at 100 steps | 0 |
| Fixed deterministic baseline | **0/4** | all four source-truncated at 100 steps | 0 |

**Important:** The runner's generic `completed: true` field means the *source episode ended*, including truncation, **not** that a goal was reached. Inspecting each final event confirms `terminated: false` and `truncated: true` in all 12 runs. No source goal was reached.

This observation **fails to demonstrate unfamiliar-family transfer**. It does not prove transfer impossible; FourRooms is a harder, longer-horizon task, the fixed baseline also failed, and this study did not include an established strong trained conventional agent. Preserve the negative, do not relabel the engineering PASS as learning or cross-family success, and do not discard the positive but narrowly scoped prior same-family reward-ablation results.

Resource observations: process wall time 4.1989 seconds; Python-parent `ru_maxrss` 68,344 KiB on GitHub-hosted Ubuntu. Native subprocess peak RSS, CPU energy, and consumer-Windows timings were **not measured**. These readings cannot qualify a retail PC.

**Next scientific experiment:** compare a full native stateful actor with strong conventional learning and planning controls, diagnose navigation failures using *source-public* observations only, and track authentic goal termination separately from source truncation. Set data/compute budgets in advance; do not optimize retrospectively for FourRooms seeds 61/67/71/73 and call those seeds held out.
