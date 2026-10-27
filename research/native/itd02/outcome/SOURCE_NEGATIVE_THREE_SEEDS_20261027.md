# ITD-02 three-seed source controls — factual negative, 2026-10-27

**Status: REAL SOURCE OBSERVATION; THREE SEEDS BUT ONE OBSERVABLE TRAJECTORY PER POLICY; TRANSFER NOT DEMONSTRATED.**

Source: [GitHub Actions run 40501585983](https://github.com/spoonman136668-ai/Wingless/actions/runs/40501585983); exact source-controls artifact ID `12733615603`, archive SHA-256 `018f2c9316ef9e9ab26bbd880324574e4097bda65b87a9120831c79c30e13650`. Source commit `4a0df25d4bebcf412a441d28c4e664bd62288024`. Source world: `MiniGrid-Empty-5x5-v0`; MiniGrid 2.3.1, Gymnasium 0.29.1. Three predeclared seeds: 7, 11, 19; 32 steps per policy per seed. Native learning begins from a new process on each seed.

| Policy | Goal finishes | Mean return | Steps per seed | Per-policy transcript diversity |
|---|---:|---:|---:|---|
| Fixed deterministic policy | 3/3 | 0.865 | 15 | **1 distinct of 3** |
| Native R8 | 0/3 | 0 | 32 | **1 distinct of 3** |
| Contextual selected-reward learner | 0/3 | 0 | 32 | **1 distinct of 3** |

Raw event digests were identical for each policy across seeds:
- Fixed: `9b28695daf3bec0d2d73e5171c2dbad51834791188795af9cb918c51511dfaa1`
- Native R8: `21c25e3928bc786c2109f11d56fb8b2101b68d85ced7aca72c16ef19f8594e7e`
- Contextual outcome: `81f584a3631bd05b458156fbda98ca114f872ab6b41fdc1e54ee9bacd2b44b8f`

The test demonstrates three *repeat executions*, **not three independent observed situations**. The deterministic Empty-5x5 family provides no empirical diversity across these seeds in this experiment. The fixed baseline wins, R8 is motionless, and contextual reward learning shows no goal completion or environment reward. The contextual policy has no positive reward signal on these trajectories and cannot be credited with improving task performance.

**Decision:** Keep the negative intact. Do not claim multi-seed generalization, successful reinforcement learning, held-out transfer, statistical replication or scientific first use. Next gate must vary source-public initial situations or genuinely distinct family configurations **before** interpreting an expanded seed panel; predeclare the diversity measure and source versions. Preserve matched controls, source-only access, action budgets and resource accounting. No mechanism tuning to these three exact outcomes.
