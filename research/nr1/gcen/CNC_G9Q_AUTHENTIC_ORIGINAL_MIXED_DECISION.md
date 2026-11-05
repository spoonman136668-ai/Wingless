# CNC-G9Q - Existing original-weight mixed-precision result

**Scientific NO_GO on the existing weight-only mixed INT4/INT2 selector. No new experiment.**

An authentic original Qwen3-30B-A3B BF16 expert comparison already ran: Wingless SHA `b2956f2aa35fa04ca04079d963326990420fe071`, [Actions 41927892073](https://github.com/spoonman136668-ai/Wingless/actions/runs/41927892073), [artifact 13370003471](https://github.com/spoonman136668-ai/Wingless/actions/runs/41927892073/artifacts/13370003471). Artifact ZIP SHA256 `80e6f7b332c0c2cf9f6faa9655a92a3419579bea0b2127c9b4d1bf4f89f7c95b`; internal `g9n-original-expert.json` SHA256 `5cc746e62cac844e901c97c1463d3cac4c81d124fd03187ca2cac6b68b5ce578`.

Model revision `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`, first-shard SHA256 `454e77b346a61bfb201d54df60e15158838cf930617ee135113556204f2802b5`, 3,999,417,504 bytes authenticated. Original expert 2048x768; frozen synthetic Gaussian input seeds 751/811, batch 8, weight-only ranking, NumPy FP32 CPU execution.

| Original expert representation | BF16 storage reduction | Mean relative squared output error |
|---|---:|---:|
| Full INT4 block64 | 3.5555x | 0.01201301 |
| Full companded INT2 block256 | 7.7571x | 0.11897125 |
| Mixed uniform eighth INT4 | 6.6991x | 0.10552902 |
| Mixed ranked eighth INT4 | 6.6991x | 0.10275411 |
| Mixed ranked quarter INT4 | 5.9412x | 0.08752880 |

Ranked one-eighth reduces mean squared error by just 2.63% relative to the equal-storage uniform control. Ranked one-quarter remains about 7.29x pure INT4 error. Pure uniform signed INT2 is worse (~0.384695 error). **No inference-ready fidelity demonstrated.** These are synthetic probes, not original-model activations; no full model, GPU speed, real task-quality or deployment claims.

**Decision:** Do not repeat the same original-weight experiment or build another kernel for this selector. Next science would require separate authorization for original-forward activation custody, frozen real task controls and matching resource budgets; not granted here.
