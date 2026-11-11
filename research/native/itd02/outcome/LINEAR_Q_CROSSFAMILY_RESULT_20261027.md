# ITD-02 linear Q vs native learner: eight-seed source result — 2026-10-27

**EXPLORATORY EVIDENCE — NOT A QUALIFIED SCIENTIFIC SUPERIORITY OR TRANSFER CLAIM.**

Source: [GitHub Actions run 40576677966](https://github.com/spoonman136668-ai/Wingless/actions/runs/40576677966); archived `itd02-linear-q-crossfamily` artifact ID `12768215357`, ZIP SHA-256 `4647fe23440efbcd0567ffce79074198a50b71c42f127ebd12647473b503c356`. Evaluated commit `237e7574ce27a24e0356441105debfd1ff91049f`. All arms received source-public six-field observations and the same source environment. Training in `MiniGrid-Empty-Random-5x5-v0` on seeds 101, 103, 107, 109, 113, 127, 131, 137. Evaluation in `MiniGrid-FourRooms-v0` on seeds 139, 149, 151, 157, 163, 167, 173, 179. Each episode maximum: 100 authentic source steps. Source versions pinned by the CI job: MiniGrid 2.3.1 and Gymnasium 0.29.1.

| Arm | Training episodes | Actual source training interactions | Training goal successes | FourRooms evaluation goals | Mean evaluation return |
|---|---:|---:|---:|---:|---:|
| Native reward-aware | 8 | **519** | 3/8 | **1/8** (seed 139, 85 steps, return 0.235) | **0.029375** |
| Native reward-blind | 8 | **631** | 2/8 | 0/8 | 0 |
| Conventional linear Q (27 parameters, fixed hyperparameters) | 8 | **431** | 5/8 | 0/8 | 0 |

**Limitations that prevent any science-passing comparison:**
1. The same source seeds and per-episode maximum are not matched **total actual interaction budgets**. Arms used 519, 631, and 431 training steps; terminate early on real source success.
2. Native reward-aware and reward-blind learners continued to update on the evaluation source transition stream. Conventional linear Q was evaluated with frozen trained weights. This is not matched evaluation update authority.
3. The conventional control uses a legitimate simple linear Q-learning algorithm with epsilon-greedy training and fixed `alpha=.1`, `gamma=.9`, `epsilon=.2`; it is **not** an independently tuned or established high-performing conventional RL baseline.
4. Only eight evaluation seeds and a known MiniGrid family training regimen, no comprehensive uncertainty estimates or retail-Windows memory/latency/energy acceptance. The FourRooms transfer success rate remains low.
5. The evaluation run's engineer-facing CI PASS means it ran and recorded outcomes; it is not scientific first use.

**Decision:** Preserve the one real native success and every failure. Do not promote these results into strong conventional superiority, matched-information equivalence, verified cross-family generalization, or scientific acceptance. The next gate must equalize source step budgets and online/frozen evaluation semantics, introduce a properly qualified strong conventional baseline, and quantify source-only navigation/coverage and resource use. Do not tune on seeds 139–179 and then call those seeds unexamined.
