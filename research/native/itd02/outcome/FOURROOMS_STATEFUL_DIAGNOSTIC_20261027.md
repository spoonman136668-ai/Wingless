# ITD-02 full-stateful FourRooms source diagnostic — 2026-10-27

**Status: AUTHENTIC EXPLORATORY SOURCE RESULT; 1/4 REWARD-AWARE GOALS; SCIENTIFIC ACCEPTANCE NOT GRANTED.**

**Custody:** GitHub Actions [run 40575034197](https://github.com/spoonman136668-ai/Wingless/actions/runs/40575034197), artifact `itd02-fourrooms-stateful-diagnostic` ID `12767645617`, ZIP SHA-256 `39961199037787297774c88a0a77c37bceb616d04611b65b4c5cd510d6e851a0`. Source commit `ea4f11080b20c0393bd65ab56c86f2027b4951c7`. Native binaries recorded as SHA-256 `d5f8388eea98c1a673edc3ce03c30ada674d32579d56bda48e0cbe00f24c15ac` (outcome) and `f570100b2ee7aa4b275c96c481bcf9ef142141930160386c85e736ad88c00a24` (R8). Pinned MiniGrid 2.3.1, Gymnasium 0.29.1.

**Training:** Existing previously explored `MiniGrid-Empty-Random-5x5-v0` seed 19, with reward-aware and zeroed-reward-blind arms undergoing matched selected-action exposure. Both training checkpoints cold-restored into fresh native processes. **Target:** `MiniGrid-FourRooms-v0` seeds 79, 83, 89, 97, maximum requested 128 steps. Actual source ends at its **100-step truncation**, which must never count as a goal.

| Policy | Goal successes | Positions visited, seeds 79/83/89/97 | Position advances, seeds 79/83/89/97 |
|---|---:|---|---|
| Native reward-aware **full stateful actor** | **1/4** (seed 83, 77 steps, return 0.307) | 47 / 46 / 32 / 41 | 70 / 54 / 65 / 68 |
| Native reward-blind full actor | 0/4 | 3 / 8 / 6 / 3 | 7 / 10 / 9 / 6 |
| Original R8 | 0/4 | 1 / 1 / 1 / 1 | 0 / 0 / 0 / 0 |
| Fixed left/forward policy | 0/4 | 15 / 14 / 12 / 15 | 72 / 73 / 72 / 74 |
| Conventional reactive compass/obstacle heuristic | **1/4** (seed 79, 12 steps, return 0.892) | 10 / 2 / 1 / 1 | 10 / 1 / 0 / 0 |

**Mechanistic diagnosis:** Reward-aware Wingless is *not simply stationary* in FourRooms: it explores many distinct public positions, and one source goal was reached. The reward-blind actor revisits few positions; R8 turns left on every action and never translates. Both the full native learner and the conventional reactive heuristic have one success on different seeds. Neither is a strong conventional trained learner; thus no scientific superiority claim follows. The previous frozen-checkpoint clone-per-action experiment had 0/4 successes, whereas the complete stateful actor yielded 1/4 on another seed panel. That is consistent with stateful within-episode dynamics being important, but seed panels were **different** and do not permit a direct causal contrast.

**Resource receipt:** process wall ~5.023s; subprocess CPU user ~1.895s, CPU system ~0.129s; child peak RSS ~69,912 KiB reported by Linux `RUSAGE_CHILDREN`. This is a process maximum, **not** sum-of-all-process peak footprint, and the source runner was GitHub-hosted Ubuntu, not a retail Windows PC. Energy not measured.

**Decision:** retain all losses and successes. No frozen scientific acceptance, no merge, no scientific first-use, no broad unfamiliar-family generalization claim. The next gate should compare full actor against a **matched strong trained baseline**, with resource and observation access parity, and ablate causal memory/route updates without tuning on these observed FourRooms seeds. Use a fresh preregistered split.
