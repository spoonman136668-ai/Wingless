# ITD-02 authenticated-source matched-control observation — 2026-10-27

**Status: SCIENTIFIC NEGATIVE / SOURCE SMOKE ONLY / NO TRANSFER CLAIM.**

Source: GitHub Actions [run 40497199593](https://github.com/spoonman136668-ai/Wingless/actions/runs/40497199593), source-controls artifact `12731775329` (`itd02-source-controls.json`). Results were read from the actual archived JSON, not inferred from the green job status.

Source revision: `c26e0c35f9f8e6f92206905696ce7a8beb076082`; source environment `MiniGrid-Empty-5x5-v0`; seed `7`; maximum 32 selected environment steps each. MiniGrid 2.3.1, Gymnasium 0.29.1. Source adapter and both native IPC executables built in the hosted run. No hidden world state exposed to the policies.

| Policy | Steps | Environment return | Terminal goal | Truncated |
|---|---:|---:|---|---|
| Fixed deterministic | 15 | 0.865 | yes | no |
| Native R8 | 32 | 0 | no | no |
| Contextual reward learner | 32 | 0 | no | no |

The fixed deterministic baseline **outperformed both native learners** in this observed source episode. The contextual learner did **not** improve over native R8. All three source runs produced observations without authentic reward differentiation for either native learner (zero reward); this dataset cannot establish that reward-based adaptation works.

This is one seed only, not a statistically powered finding. No held-out family, independent training/evaluation partition, frozen checkpoint retention test, adaptive-learning control, resource accounting or scientific FSA is established. Preserve the observed losses; do not relabel CI green as scientific success.

**Next discriminator:** predeclare a bounded multi-seed comparison with identical source budget and three policies; preserve all individual seeds and source transcript digests. Only after that add training/cold-restart/held-out transfer controls. Do not tune based solely on this seed.
