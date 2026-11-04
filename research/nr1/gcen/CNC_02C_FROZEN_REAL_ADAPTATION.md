# CNC-02C — frozen authentic Qwen expert-weight adaptation
Status: development-only, no scientific acceptance or production activation.

- Checkpoint: Qwen/Qwen3-30B-A3B revision ad44e777bcd18fa416d9da3bd8f70d33ebb85d39.
- Source: one G31/G32 SHA256-verified original BF16 shard, first complete down_proj 16-expert layer, first FOUR authentic experts.
- Frozen original-window seeds: 11, 23, 37, 53, 71.
- Frozen original 128x256 window coordinates as G32 for each seed.
- Independent synthetic Gaussian FIT seed 6011, TRANSFER seed 17201, 32 examples each.
- Candidate: G8 block64 FP32-scale INT4 direct code updates, two coordinate updates per expert window, no retained master-weight or optimizer moments.
- Control: frozen INT4 representation without updates. NO matched full-precision training control yet; do not claim superiority to FP32/QLoRA or task transfer.
- Evaluate FIT and TRANSFER squared-output-error independently, logging edits, persistent compressed state and process RSS. TRANSFER labels never affect optimization.
- Exactly 20 expert-window comparisons (4 experts x 5 seeds). Reject missing or modified source identity.
- No original Qwen forward pass, activation capture, real-language task quality, GPU kernel, scientific acceptance, or deployment.
- One-shot exact-head CI cnc02c_real_weight_fit_transfer runs after mandatory Linux/Windows native regression and only on the CNC02 branch when head commit includes [cnc02c-real].
- Independent compressed-inference PR #443 untouched.
